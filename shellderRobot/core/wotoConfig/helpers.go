package wotoConfig

import (
	"github.com/ALiwoto/StrongStringGo/strongStringGo"
)

func ParseConfig(filename string) (*BotConfig, error) {
	if ConfigSettings != nil {
		return ConfigSettings, nil
	}
	config := &BotConfig{}

	err := strongStringGo.ParseConfig(config, filename)
	if err != nil {
		return nil, err
	}

	ConfigSettings = config

	return ConfigSettings, nil
}

func LoadConfig() (*BotConfig, error) {
	return ParseConfig("config.ini")
}

func IsDebug() bool {
	if ConfigSettings != nil {
		return ConfigSettings.IsDebug
	}
	return true
}

func IsAllowed(id int64) bool {
	if ConfigSettings == nil {
		return false
	}

	for _, current := range ConfigSettings.OwnerIds {
		if current == id {
			return true
		}
	}

	return false
}

func GetBotToken() string {
	if ConfigSettings != nil {
		return ConfigSettings.BotToken
	}
	return ""
}

func DropUpdates() bool {
	if ConfigSettings != nil {
		return ConfigSettings.DropUpdates
	}
	return false
}

func GetCmdPrefixes() []rune {
	return []rune{'/', '!'}
}

func GetHandlerCommand() string {
	if ConfigSettings == nil {
		return "sh"
	}

	return ConfigSettings.HandlerCmd
}
