package shellPlugin

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/AnimeKaizoku/shellderRobot/shellderRobot/core/wotoConfig"
	wv "github.com/AnimeKaizoku/shellderRobot/shellderRobot/core/wotoValues"
	"github.com/AnimeKaizoku/ssg/ssg"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
)

func Shellout(command string) (string, string, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	var cmd *exec.Cmd
	if os.PathSeparator == '/' {
		cmd = exec.Command(ShellToUseUnix, "-c", command)
	} else {
		cmd = exec.Command(ShellToUseWin, "/C", command)
	}
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	return stdout.String(), stderr.String(), err
}

func generateCancelButton(uniqueId string) *gotgbot.InlineKeyboardMarkup {
	return &gotgbot.InlineKeyboardMarkup{
		InlineKeyboard: [][]gotgbot.InlineKeyboardButton{
			{
				{
					Text:         "Cancel",
					CallbackData: executeCancelDataPrefix + cbDataSep + uniqueId,
				},
			},
		},
	}
}

func generateUniqueId() string {
	idGeneratorMutex.Lock()
	defer idGeneratorMutex.Unlock()
	lastId++

	return strconv.Itoa(lastId) + "Z" + ssg.ToBase32(time.Now().Unix())
}

func DownloadFile(filePath string) ([]byte, error) {
	pre := fmt.Sprintf("%s/file/bot%s/", wv.HelperBot.GetAPIURL(), wv.HelperBot.GetToken())

	resp, err := http.Get(pre + filePath)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

func Cmdout(command string) (string, string, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command(ShellToUseWin, "/C", command)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	return stdout.String(), stderr.String(), err
}

func Bashout(command string) (string, string, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command(ShellToUseUnix, "-c", command)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	return stdout.String(), stderr.String(), err
}

func LoadAllHandlers(d *ext.Dispatcher, t []rune) {
	cmdPre := wotoConfig.GetHandlerCommand()
	shellCommand := handlers.NewCommand(cmdPre, shellHandler)
	vServersCommand := handlers.NewCommand(vServersCmd, shellHandler)
	downloadCommand := handlers.NewCommand(cmdPre+downloadCmd, downloadHandler)
	uploadCommand := handlers.NewCommand(cmdPre+uploadCmd, uploadHandler)
	dlCommand := handlers.NewCommand(cmdPre+dlCmd, downloadHandler)
	ulCommand := handlers.NewCommand(cmdPre+ulCmd, uploadHandler)
	exitCommand := handlers.NewCommand(cmdPre+exitCmd, exitHandler)
	cancelCallBack := handlers.NewCallback(cancelButtonFilter, cancelButtonCallBackQuery)

	shellCommand.Triggers = t
	vServersCommand.Triggers = t
	downloadCommand.Triggers = t
	uploadCommand.Triggers = t
	dlCommand.Triggers = t
	ulCommand.Triggers = t
	exitCommand.Triggers = t

	d.AddHandler(vServersCommand)
	d.AddHandler(shellCommand)
	d.AddHandler(downloadCommand)
	d.AddHandler(uploadCommand)
	d.AddHandler(dlCommand)
	d.AddHandler(ulCommand)
	d.AddHandler(exitCommand)
	d.AddHandler(cancelCallBack)
}
