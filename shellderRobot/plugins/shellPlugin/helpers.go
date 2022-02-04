package shellPlugin

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"

	"github.com/AnimeKaizoku/shellderRobot/shellderRobot/core/wotoConfig"
	wv "github.com/AnimeKaizoku/shellderRobot/shellderRobot/core/wotoValues"
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

func DownloadFile(filePath string) ([]byte, error) {
	pre := fmt.Sprintf("%s/file/bot%s/", wv.HelperBot.GetAPIURL(), wv.HelperBot.Token)

	resp, err := http.Get(pre + filePath)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
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
	ShellCommand := handlers.NewCommand(cmdPre, shellHandler)
	downloadCommand := handlers.NewCommand(cmdPre+downloadCmd, downloadHandler)
	uploadCommand := handlers.NewCommand(cmdPre+uploadCmd, uploadHandler)
	exitCommand := handlers.NewCommand(cmdPre+uploadCmd, exitHandler)
	ShellCommand.Triggers = t
	downloadCommand.Triggers = t
	uploadCommand.Triggers = t
	exitCommand.Triggers = t
	d.AddHandler(ShellCommand)
	d.AddHandler(downloadCommand)
	d.AddHandler(uploadCommand)
	d.AddHandler(exitCommand)
}
