package listeners

import (
	"fmt"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/aasumitro/svc-telegram-notify/config"
	"strconv"
)

// MessagingListener represent the data-struct for configuration
type MessagingListener struct {
	chatID string
	message string
	appConfig *config.AppConfig
}

func (listener MessagingListener) SendNotify()  {
	fmt.Println(fmt.Sprintf(
		"Trying send message to: %s",
		listener.chatID,
	))
	fmt.Println("=====================================================")

	chatID, _ := strconv.ParseInt(listener.chatID, 10, 64)
	message, err := listener.appConfig.GetTelegramConnection().SendMessage(
		chatID, listener.message, &gotgbot.SendMessageOpts{},
	)
	if err != nil {
		fmt.Println(fmt.Sprintf("failed send notificaion cause: %v\n", err))
	} else {
		// Response is a message ID string.
		fmt.Println("Successfully sent message to:", message.Chat.Username)
	}
	fmt.Println("=====================================================")
}

// InitMessagingListener initialize the app configuration
func InitMessagingListener(
	chatID string,
	message string,
	appConfig *config.AppConfig,
) *MessagingListener {
	return &MessagingListener{
		chatID: chatID,
		message: message,
		appConfig: appConfig,
	}
}