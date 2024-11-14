package send

import (
	"database/sql"
	"footballresult/logger"
	"footballresult/send/telegram"
	"time"
)

func Telegram(db *sql.DB) {

	logger.Send("sender is starting...")

	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()
	i := 1
	for {
		select {
		case <-ticker.C:
			result, err := telegram.SendProcessing(db)
			if err != nil {
				logger.Send(err.Error())
			} else if result != "" {
				logger.Send(result)
				i = 0
			} else {
				i = i + 1
				if i > 60 {
					logger.Send("no sent messages for the last hour")
					i = 0
				}
			}

		}
	}

}
