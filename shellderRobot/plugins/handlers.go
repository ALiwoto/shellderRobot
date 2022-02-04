package plugins

import (
	"github.com/AnimeKaizoku/shellderRobot/shellderRobot/plugins/shellPlugin"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

func LoadAllHandlers(d *ext.Dispatcher, triggers []rune) {
	shellPlugin.LoadAllHandlers(d, triggers)
}
