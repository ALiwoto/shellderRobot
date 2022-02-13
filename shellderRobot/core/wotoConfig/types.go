package wotoConfig

type BotConfig struct {
	BotToken          string  `section:"general" key:"bot_token"`
	OwnerIds          []int64 `section:"general" key:"owner_ids"`
	HandlerCmd        string  `section:"general" key:"handler_command"`
	DownloadDirectory string  `section:"general" key:"download_directory"`
	DropUpdates       bool    `section:"general" key:"drop_updates"`
	IsDebug           bool    `section:"general" key:"debug"`
}
