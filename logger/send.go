package logger

import (
	"footballresult/config"
	"footballresult/send/telegram"
	"log"
)

func Send(message string) {

	log.Printf(message)
	botToken := config.Load.TelegramBotToken
	logChannelID := config.Load.TelegramLogChannelID

	err := telegram.SendMessageToTelegram(botToken, logChannelID, message)
	if err != nil {
		log.Printf("ERROR: %v", err)
	}

}
