package shellPlugin

import (
	"github.com/AnimeKaizoku/ssg/ssg/shellUtils"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

type commandContainer struct {
	result            *shellUtils.ExecuteCommandResult
	bot               *gotgbot.Bot
	userContext       *ext.Context
	botMessage        *gotgbot.Message
	killRequestedBy   *gotgbot.User
	isCanceled        bool
	isRunningSilently bool
}
