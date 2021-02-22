package handlers

import (
	"fmt"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters"
	"github.com/aasumitro/svc-telegram-notify/config"
)

type telegramCommandHandler struct {
	config       *config.AppConfig
}

func NewTelegramCommandHandler(
	appConfig *config.AppConfig,
)  {
	commandHandler := &telegramCommandHandler{config: appConfig}

	fmt.Println(fmt.Sprintf(
		"%s has been started. . .",
		appConfig.GetTelegramConnection().User.Username))
	fmt.Println("waiting request from user. . .")

	updater := ext.NewUpdater(appConfig.GetTelegramConnection(), nil)
	dispatcher := updater.Dispatcher

	// Add handler to reply to all messages.
	dispatcher.AddHandler(handlers.NewCommand(
		"start",
		commandHandler.startCommand,
	))
	dispatcher.AddHandler(handlers.NewCallback(
		filters.Equal("chat_id_command"),
		commandHandler.chatIDCallback,
	))
	dispatcher.AddHandler(handlers.NewCallback(
		filters.Equal("help_command"),
		commandHandler.helpCallback,
	))

	// Start receiving updates.
	err := updater.StartPolling(
		appConfig.GetTelegramConnection(),
		&ext.PollingOpts{Clean: true})
	if err != nil {
		fmt.Println("failed to start polling: " + err.Error())
	}

	// Idle, to keep updates coming in, and avoid bot stopping.
	updater.Idle()
}

func (handler telegramCommandHandler) startCommand(ctx *ext.Context) error {
	_, err := ctx.EffectiveMessage.Reply(
		ctx.Bot,
		fmt.Sprintf(
			"Hello, I'm @%s. Your <b>personal assistant</b>, how can i help you?",
			ctx.Bot.User.Username,
		),
		&gotgbot.SendMessageOpts{
			ParseMode: "html",
			ReplyMarkup: gotgbot.InlineKeyboardMarkup{
				InlineKeyboard: [][]gotgbot.InlineKeyboardButton{{
					{Text: "My Telegram ID", CallbackData: "chat_id_command"},
					{Text: "Please Help Me!", CallbackData: "help_command"},
				}},
			},
		},
	)

	if err != nil {
		fmt.Println("failed to send: " + err.Error())
	}

	return nil
}

func (handler telegramCommandHandler) chatIDCallback(ctx *ext.Context) error {
	callback := ctx.Update.CallbackQuery

	_, err := callback.Answer(ctx.Bot, nil)
	if err != nil {
		fmt.Println("failed initialize callback answer: " + err.Error())
	}

	_, err = callback.Message.EditText(
		ctx.Bot,
		fmt.Sprintf(
			"Your Telegram ID : %d, \nYour Chat ID : %d",
			ctx.Bot.User.Id,
			ctx.EffectiveChat.Id,
		),
	nil)
	if err != nil {
		fmt.Println("failed send callback to user: " + err.Error())
	}

	return nil
}

func (handler telegramCommandHandler) helpCallback(ctx *ext.Context) error {
	callback := ctx.Update.CallbackQuery

	_, err := callback.Answer(ctx.Bot, nil)
	if err != nil {
		fmt.Println("failed initialize callback answer: " + err.Error())
	}

	_, err = callback.Message.EditText(
		ctx.Bot,
		"Hello dear, i will help you to access your registered account to use our service",
		nil)
	if err != nil {
		fmt.Println("failed send callback to user: " + err.Error())
	}

	return nil
}