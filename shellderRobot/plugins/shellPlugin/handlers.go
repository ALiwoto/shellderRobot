package shellPlugin

import (
	"bytes"
	"context"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/AnimeKaizoku/shellderRobot/shellderRobot/core/logging"
	"github.com/AnimeKaizoku/shellderRobot/shellderRobot/core/utils"
	"github.com/AnimeKaizoku/shellderRobot/shellderRobot/core/wotoConfig"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"

	"github.com/ALiwoto/mdparser/mdparser"
	ws "github.com/AnimeKaizoku/ssg/ssg"
)

func termHandlerBase(b *gotgbot.Bot, ctx *ext.Context) error {
	msg := ctx.EffectiveMessage
	wholeStrs := ws.SplitN(msg.Text, 2, " ", "\n", "\r", "\t")
	if len(wholeStrs) < 2 {
		// No command
		// TODO: show stats and... stuff
		return ext.EndGroups
	}
	whole := wholeStrs[1]
	whole = strings.TrimSpace(whole)

	finishChan := make(chan bool)

	result := ws.RunCommandAsyncWithChan(whole, finishChan)
	result.UniqueId = generateUniqueId()

	finishedFunc := func() {
		var errStr string
		err := result.Error
		output := result.Stdout
		errOut := result.Stderr
		if err != nil {
			errStr = err.Error()
		}

		if len(output+errOut+errStr) > 4080 {
			myAllStr := output + "\n\n" + errOut + "\n\n" + errStr
			namedFile := &gotgbot.NamedFile{
				File:     bytes.NewBuffer([]byte(myAllStr)),
				FileName: "output.txt",
			}
			_, _ = b.SendDocument(msg.Chat.Id, namedFile, &gotgbot.SendDocumentOpts{
				ReplyToMessageId:         msg.MessageId,
				AllowSendingWithoutReply: true,
			})
			return
		}

		if output == "" && errOut == "" && err == nil {
			_, _ = b.SendMessage(msg.Chat.Id, "No output", &gotgbot.SendMessageOpts{
				ParseMode:                utils.MarkDownV2,
				ReplyToMessageId:         msg.MessageId,
				AllowSendingWithoutReply: true,
			})
			return
		}

		if errStr != "" {
			md := mdparser.GetBold("Error:\n").Mono(errStr)
			if output != "" {
				md.Normal("\n\n").Mono(output)
			}

			md.Normal("\n\n").Mono(errOut)
			_, err = b.SendMessage(msg.Chat.Id, md.ToString(), &gotgbot.SendMessageOpts{
				ParseMode:                utils.MarkDownV2,
				ReplyToMessageId:         msg.MessageId,
				AllowSendingWithoutReply: true,
			})
			if err != nil {
				logging.Error(err)
			}

			return
		}

		md := mdparser.GetBold("Output:\n").Mono(output)
		if errOut != "" {
			md.Normal("\n\n").Bold("StdError:\n").Mono(errOut)
		}

		_, err = b.SendMessage(msg.Chat.Id, md.ToString(), &gotgbot.SendMessageOpts{
			ParseMode:                utils.MarkDownV2,
			ReplyToMessageId:         msg.MessageId,
			AllowSendingWithoutReply: true,
		})
		if err != nil {
			logging.Error(err)
		}
	}

	deadlineCtx, cancelContext := context.WithTimeout(context.Background(), time.Second*2)
	defer cancelContext()

	var container *commandContainer

	select {
	case <-deadlineCtx.Done():
		// command is taking longer than expected
		container = &commandContainer{
			result:      result,
			bot:         b,
			userContext: ctx,
		}

		uId := container.GetUniqueId()
		commandsMap.Add(uId, container)
		botMsg, err := msg.Reply(b, container.ParseAsMd().ToString(), &gotgbot.SendMessageOpts{
			ParseMode:             utils.MarkDownV2,
			ReplyMarkup:           generateCancelButton(uId),
			DisableWebPagePreview: true,
		})
		if err != nil {
			container.isRunningSilently = true
		}

		container.botMessage = botMsg
	case <-finishChan:
		// Finished, send the output directly
		finishedFunc()
		return ext.EndGroups
	}

	<-finishChan
	if container.isCanceled {
		// we assume that cancel callback handler has already handled
		// everything here, all we have to do here is to return and kill
		// the goroutine.
		return ext.EndGroups
	}

	if result.IsDone() {
		finishedFunc()
		if container.botMessage != nil {
			// this here needs a better design, idk maybe show execution time,
			// or put button there to paste it on pasty or something like that...
			// maybe in future.
			_, _ = container.botMessage.Delete(b, nil)
		}
		return nil
	} else {
		// impossible to reach, but needs more investigation...
		log.Println("reached ")
		finishedFunc()
	}

	return ext.EndGroups
}

func cancelButtonFilter(cq *gotgbot.CallbackQuery) bool {
	return strings.HasPrefix(cq.Data, executeCancelDataPrefix+cbDataSep)
}

func cancelButtonCallBackQuery(b *gotgbot.Bot, ctx *ext.Context) error {
	user := ctx.EffectiveUser
	query := ctx.CallbackQuery
	if !wotoConfig.IsAllowed(user.Id) {
		_, _ = query.Answer(b, &gotgbot.AnswerCallbackQueryOpts{
			Text:      "This button is not for you...",
			CacheTime: 5500,
		})
		return ext.EndGroups
	}

	// format is like: caEx_UNIQUE-ID
	myStrs := strings.Split(query.Data, cbDataSep)
	if len(myStrs) < 2 {
		// impossible to happen, this condition is here only to prevent
		// panics
		return ext.EndGroups
	}

	uniqueId := myStrs[1]
	container := commandsMap.Get(uniqueId)
	if container == nil {
		// data is either too old, or it has been removed
		// from our memory... in any case, this shouldn't
		// happen here, we should make sure the moment data
		// is deleted from memory, bot's message is edited as well.
		// (unless the bot is rebooted)
		_, _ = query.Answer(b, &gotgbot.AnswerCallbackQueryOpts{
			Text:      "This command is either too old, or I've removed it from my memory...",
			CacheTime: 5500,
		})
		return ext.EndGroups
	}

	container.killRequestedBy = user
	_, _ = query.Answer(b, &gotgbot.AnswerCallbackQueryOpts{
		Text:      "Killing the process, please wait...",
		CacheTime: 5500,
	})

	container.Kill()
	commandsMap.Delete(uniqueId)

	return ext.EndGroups
}

func exitHandler(b *gotgbot.Bot, ctx *ext.Context) error {
	user := ctx.EffectiveUser
	if user == nil || !wotoConfig.IsAllowed(user.Id) {
		return ext.EndGroups
	}

	msg := ctx.EffectiveMessage
	myStr := ws.SplitN(msg.Text, 2, " ", "\n", "\r", "\t")
	var exitCode int
	if len(myStr) > 1 {
		whole := myStr[1]
		whole = strings.TrimSpace(whole)
		exitCode, _ = strconv.Atoi(whole)
	}

	md := mdparser.GetNormal("Exiting with code " + strconv.Itoa(exitCode))
	_, _ = msg.Reply(b, md.ToString(), &gotgbot.SendMessageOpts{
		ParseMode: utils.MarkDownV2,
	})
	os.Exit(exitCode)

	return ext.EndGroups
}

func uploadHandler(b *gotgbot.Bot, ctx *ext.Context) error {
	user := ctx.EffectiveUser
	if user == nil || !wotoConfig.IsAllowed(user.Id) {
		return ext.EndGroups
	}

	msg := ctx.EffectiveMessage
	myStrs := ws.SplitN(msg.Text, 2, " ", "\n", "\r", "\t")
	if len(myStrs) < 2 {
		md := mdparser.GetBold("You need to specify a local file name/path...")
		_, _ = msg.Reply(b, md.ToString(), &gotgbot.SendMessageOpts{
			ParseMode:                utils.MarkDownV2,
			AllowSendingWithoutReply: false,
			DisableWebPagePreview:    true,
		})
		return ext.EndGroups
	}

	whole := myStrs[1]
	whole = strings.TrimSpace(whole)
	myFile, err := os.Open(whole)
	if err != nil {
		errMd := mdparser.GetBold("Error:\n").Mono(err.Error())
		_, _ = msg.Reply(b, errMd.ToString(), &gotgbot.SendMessageOpts{
			ParseMode:                utils.MarkDownV2,
			AllowSendingWithoutReply: true,
		})
		return ext.EndGroups
	}

	md := mdparser.GetBold("Uploading ").Mono(whole).Normal("...")
	topMsg, _ := msg.Reply(b, md.ToString(), &gotgbot.SendMessageOpts{
		ParseMode:                utils.MarkDownV2,
		AllowSendingWithoutReply: true,
	})

	f := gotgbot.NamedFile{
		FileName: path.Base(whole),
		File:     myFile,
	}

	if len(whole) > 2040 {
		whole = whole[:2040]
	}

	_, err = b.SendDocument(msg.Chat.Id, f, &gotgbot.SendDocumentOpts{
		ParseMode:        utils.MarkDownV2,
		ReplyToMessageId: msg.MessageId,
		Caption:          mdparser.GetMono(whole).ToString(),
	})

	if topMsg != nil {
		if err != nil {
			md := mdparser.GetBold("Error:\n").Mono(err.Error())
			_, _, _ = topMsg.EditText(b, md.ToString(), &gotgbot.EditMessageTextOpts{
				ParseMode: utils.MarkDownV2,
			})
			return ext.EndGroups
		}

		_, _ = topMsg.Delete(b, nil)
	}

	return ext.EndGroups
}

func downloadHandler(b *gotgbot.Bot, ctx *ext.Context) error {
	user := ctx.EffectiveUser
	if user == nil || !wotoConfig.IsAllowed(user.Id) {
		return ext.EndGroups
	}

	msg := ctx.EffectiveMessage
	if msg.ReplyToMessage == nil {
		_, _ = msg.Reply(b, "Reply to something...", &gotgbot.SendMessageOpts{
			ReplyToMessageId:         msg.MessageId,
			AllowSendingWithoutReply: false,
		})
		return ext.EndGroups
	}
	replied := msg.ReplyToMessage

	var fileId string
	var fileType string
	switch {
	case replied.Animation != nil:
		fileType = "animation"
		fileId = replied.Animation.FileId
	case replied.Audio != nil:
		fileType = "audio"
		fileId = replied.Audio.FileId
	case replied.Document != nil:
		fileType = "document"
		fileId = replied.Document.FileId
	case replied.Photo != nil:
		fileType = "photo"
		fileId = replied.Photo[len(replied.Photo)-1].FileId
	case replied.Sticker != nil:
		fileType = "sticker"
		fileId = replied.Sticker.FileId
	case replied.Video != nil:
		fileType = "video"
		fileId = replied.Video.FileId
	case replied.Voice != nil:
		fileType = "voice"
		fileId = replied.Voice.FileId
	case replied.VideoNote != nil:
		fileType = "video note"
		fileId = replied.VideoNote.FileId
	default:
		_, _ = msg.Reply(b, "No media specified...", &gotgbot.SendMessageOpts{
			ReplyToMessageId:         msg.MessageId,
			AllowSendingWithoutReply: false,
		})
		return ext.EndGroups
	}

	allStrs := ws.SplitN(msg.Text, 2, " ", "\n", "\r", "\t")
	myPath := ""
	if len(allStrs) > 1 {
		myPath = allStrs[1]
	}

	if !strings.Contains(myPath, string(os.PathSeparator)) {
		myPath = wotoConfig.GetDownloadDirectory() + myPath
	}

	md := mdparser.GetMono("Downloading ").Bold(fileType)
	if myPath != "" {
		md.Normal(" to ").Mono(myPath)
	}

	topMsg, _ := msg.Reply(b, md.ToString(), &gotgbot.SendMessageOpts{
		ParseMode:             utils.MarkDownV2,
		DisableWebPagePreview: true,
	})

	f, err := b.GetFile(fileId, nil)
	if err != nil {
		return utils.SendAlertErr(b, msg, err)
	}

	bytes, err := DownloadFile(f.FilePath)
	if err != nil {
		return utils.SendAlertErr(b, msg, err)
	}

	if myPath == "" {
		myPath = f.FilePath
	}

	err = os.WriteFile(myPath, bytes, 0644)
	if err != nil {
		return utils.SendAlertErr(b, msg, err)
	}

	if topMsg != nil {
		_, _ = topMsg.Delete(b, nil)
	}

	md = mdparser.GetBold("Downloaded to ").Mono(myPath)

	_, _ = msg.Reply(b, md.ToString(), &gotgbot.SendMessageOpts{
		ParseMode: utils.MarkDownV2,
	})

	return ext.EndGroups
}

func shellHandler(b *gotgbot.Bot, ctx *ext.Context) error {
	user := ctx.EffectiveUser
	if user == nil || !wotoConfig.IsAllowed(user.Id) {
		return ext.EndGroups
	}

	go termHandlerBase(b, ctx)

	return ext.EndGroups
}
