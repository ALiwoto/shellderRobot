package shellPlugin

import (
	"io/ioutil"
	"os"
	"path"
	"runtime"
	"strconv"
	"strings"

	"github.com/AnimeKaizoku/shellderRobot/shellderRobot/core/logging"
	"github.com/AnimeKaizoku/shellderRobot/shellderRobot/core/utils"
	"github.com/AnimeKaizoku/shellderRobot/shellderRobot/core/wotoConfig"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"

	ws "github.com/ALiwoto/StrongStringGo/strongStringGo"
	"github.com/ALiwoto/mdparser/mdparser"
)

func termHandlerBase(b *gotgbot.Bot, ctx *ext.Context, getter outputGetter) error {
	msg := ctx.EffectiveMessage
	whole := strings.Join(ws.SplitN(msg.Text, 2, " ", "\n", "\r", "\t")[1:], "")
	whole = strings.TrimSpace(whole)

	output, errOut, err := getter(whole)

	var errStr string
	if err != nil {
		errStr = err.Error()
	}
	if len(output+errOut+errStr) > 4080 {
		myAllStr := output + "\n\n" + errOut + "\n\n" + errStr
		_, _ = b.SendDocument(msg.Chat.Id, []byte(myAllStr), &gotgbot.SendDocumentOpts{
			ReplyToMessageId:         msg.MessageId,
			AllowSendingWithoutReply: true,
		})
		return ext.EndGroups
	}

	if output == "" && errOut == "" && err == nil {
		_, _ = b.SendMessage(msg.Chat.Id, "No output", &gotgbot.SendMessageOpts{
			ParseMode:                utils.MarkDownV2,
			ReplyToMessageId:         msg.MessageId,
			AllowSendingWithoutReply: true,
		})
		return ext.EndGroups
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

		return ext.EndGroups
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

		_, _ = topMsg.Delete(b)
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

	f, err := b.GetFile(fileId)
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

	err = ioutil.WriteFile(myPath, bytes, 0644)
	if err != nil {
		return utils.SendAlertErr(b, msg, err)
	}

	if topMsg != nil {
		_, _ = topMsg.Delete(b)
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

	switch runtime.GOOS {
	case "linux":
		return termHandlerBase(b, ctx, Shellout)
	case "windows":
		return termHandlerBase(b, ctx, Cmdout)
	}

	_, _ = ctx.EffectiveMessage.Reply(b, unsupportedMessage, nil)

	return ext.EndGroups
}
