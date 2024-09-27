package telegram

import (
	"fmt"
	"footballresult/storage"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

func SendMessageToTelegram(message string) error {

	err := godotenv.Load(".env")
	if err != nil {
		return fmt.Errorf("error loading .env file: %v", err)
	}

	botToken, err := storage.LoadEnvVariable("TELEGRAM_BOT_TOKEN")
	if err != nil {
		return err
	}

	channelID, err := storage.LoadEnvVariable("TELEGRAM_CHANNEL_ID")
	if err != nil {
		return err
	}

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		return err
	}

	msg := tgbotapi.NewMessageToChannel(channelID, message)

	_, err = bot.Send(msg)
	if err != nil {
		return err
	}

	return nil
}
