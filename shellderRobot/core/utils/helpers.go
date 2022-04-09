package utils

import (
	"github.com/ALiwoto/mdparser/mdparser"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

func SendAlert(b *gotgbot.Bot, m *gotgbot.Message, md mdparser.WMarkDown) error {
	md.Replace(b.Token, "$TOKEN")
	_, _ = m.Reply(b, md.ToString(), &gotgbot.SendMessageOpts{
		ParseMode:                MarkDownV2,
		DisableWebPagePreview:    true,
		AllowSendingWithoutReply: true,
	})

	return ext.EndGroups
}

func SendAlertErr(b *gotgbot.Bot, m *gotgbot.Message, e error) error {
	if e == nil {
		return ext.EndGroups
	}

	md := mdparser.GetBold("Error: \n").Mono(e.Error())
	return SendAlert(b, m, md)
}
