package config

import (
	"fmt"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/spf13/viper"
)

var telegramConn *gotgbot.Bot

func (config AppConfig) SetupTelegramConnection() {
	telegramBot, err := gotgbot.NewBot(viper.GetString("TELEGRAM_KEY"))
	if err != nil {
		panic(fmt.Sprintf(
			"failed to create new bot: , cause: %s",
			err.Error(),
		))
	}

	setTelegramConnection(telegramBot)
}

func setTelegramConnection(telegramCurrentConn *gotgbot.Bot) {
	telegramConn = telegramCurrentConn
}

func (config AppConfig) GetTelegramConnection() *gotgbot.Bot {
	return telegramConn
}