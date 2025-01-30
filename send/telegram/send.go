package telegram

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func SendMessageToTelegram(botToken, channelID, message string) error {
	// Создаем нового бота с использованием переданного токена
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		return fmt.Errorf("failed to create bot: %v", err)
	}

	// Создаем сообщение для отправки в указанный канал
	msg := tgbotapi.NewMessageToChannel(channelID, message)
	msg.ParseMode = "HTML"

	// Отправляем сообщение через API
	_, err = bot.Send(msg)
	if err != nil {
		return fmt.Errorf("failed to send message to telegram: %v", err)
	}

	return nil
}
