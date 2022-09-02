package shellPlugin

import (
	"github.com/ALiwoto/mdparser/mdparser"
	"github.com/AnimeKaizoku/shellderRobot/shellderRobot/core/utils"
	"github.com/PaulSonOfLars/gotgbot/v2"
)

// GetUniqueId returns the unique-id of this command-container
// variable.
func (c *commandContainer) GetUniqueId() string {
	return c.result.UniqueId
}

// SetUniqueId method sets the specified string value as unique-id of
// this commandContainer variable.
func (c *commandContainer) SetUniqueId(value string) {
	c.result.UniqueId = value
}

func (c *commandContainer) ParseAsMd() mdparser.WMarkDown {
	if !c.isCanceled {
		return mdparser.GetBold("Executing ").Mono("#" + c.GetUniqueId()).Normal("...")
	}

	md := mdparser.GetNormal("Command ").Mono("#" + c.GetUniqueId())
	if c.killRequestedBy != nil {
		md.Normal(" has been canceled by ")
		md.Mention(c.killRequestedBy.FirstName, c.killRequestedBy.Id).Normal(".")
	} else {
		md.Normal(" has been canceled.")
	}
	return md
}

// Kill  method kills the container, edits the message, etc...
func (c *commandContainer) Kill() {
	if c.isCanceled {
		return
	}
	c.isCanceled = true

	if c.isRunningSilently {
		// everything is being done silently, no need to edit anything...
		_ = c.result.Kill()
		return
	}

	// if we are not running silently (e.g bot has sent a message replying to
	// the user that we are executing command with a cancel button under it),
	// we have to first edit the message and then kill the process.

	_, _, _ = c.botMessage.EditText(c.bot, c.ParseAsMd().ToString(), &gotgbot.EditMessageTextOpts{
		ParseMode: utils.MarkDownV2,
	})
}
