package send

import (
	"database/sql"
	"footballresult/config"
	"footballresult/send/telegram"
	"log"
	"time"
)

func Telegram(db *sql.DB) {

	botToken := config.Load.TelegramBotToken
	logChannelID := config.Load.TelegramLogChannelID

	log.Printf("sender is starting...")
	err := telegram.SendMessageToTelegram(botToken, logChannelID, "Sender is starting...")
	if err != nil {
		log.Printf("ERROR: %v", err)
	}

	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()
	i := 1
	for {
		select {
		case <-ticker.C:
			result, err := telegram.SendProcessing(db)
			if err != nil {
				log.Printf("ERROR: %v", err)
			} else if result != "" {
				log.Printf(result)
				err = telegram.SendMessageToTelegram(botToken, logChannelID, result)
				if err != nil {
					log.Printf("ERROR: %v", err)
				}
				i = 0
			} else {
				i = i + 1
				if i > 60 {
					log.Printf("no sent messages for last hour")
					err := telegram.SendMessageToTelegram(botToken, logChannelID, "no sent messages for last hour")
					if err != nil {
						log.Printf("ERROR: %v", err)
					}
					i = 0
				}
			}

		}
	}

}
