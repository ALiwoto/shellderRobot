package plugins

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/ALiwoto/mdparser/mdparser"
	"github.com/AnimeKaizoku/shellderRobot/shellderRobot/core/logging"
	"github.com/AnimeKaizoku/shellderRobot/shellderRobot/core/wotoConfig"
	wv "github.com/AnimeKaizoku/shellderRobot/shellderRobot/core/wotoValues"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

func StartTelegramBot() error {
	token := wotoConfig.GetBotToken()
	if len(token) == 0 {
		return errors.New("bot token is empty")
	}

	b, err := gotgbot.NewBot(token, &gotgbot.BotOpts{
		Client: http.Client{},
		RequestOpts: &gotgbot.RequestOpts{
			Timeout: 2 * gotgbot.DefaultTimeout,
		},
	})
	if err != nil {
		return err
	}

	mdparser.AddSecret(b.GetToken(), "$TOKEN")

	uOptions := &ext.UpdaterOpts{
		DispatcherOpts: ext.DispatcherOpts{
			MaxRoutines: -1,
		},
	}
	utmp := ext.NewUpdater(uOptions)
	updater := &utmp
	err = updater.StartPolling(b, &ext.PollingOpts{
		DropPendingUpdates: wotoConfig.DropUpdates(),
	})
	if err != nil {
		return err
	}

	logging.Info(fmt.Sprintf("%s has started | ID: %d", b.Username, b.Id))

	wv.HelperBot = b
	wv.BotUpdater = updater

	LoadAllHandlers(updater.Dispatcher, wotoConfig.GetCmdPrefixes())

	updater.Idle()
	return nil
}
