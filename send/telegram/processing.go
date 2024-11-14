package telegram

import (
	"database/sql"
	"fmt"
	"footballresult/config"
	"footballresult/storage"
	"footballresult/types"
	"strconv"
	"time"
)

func eventExpired(eventDate time.Time) bool {

	now := time.Now()

	duration := now.Sub(eventDate)

	return duration > 4*time.Hour

}

func SendProcessing(db *sql.DB) (string, error) {

	botToken := config.Load.TelegramBotToken
	channelID := config.Load.TelegramChannelID

	finishedEvents, err := GetFinishedEventsFromDB(db)
	if err != nil {
		return "", err
	}

	var result string

	if len(finishedEvents) > 0 {
		result = "events: "
		var updateFinishedEvents []types.Event

		for _, event := range finishedEvents {

			if !eventExpired(event.EventDate) {
				message := event.TeamHome + " " + strconv.Itoa(event.GoalsHome) + " - " + strconv.Itoa(event.GoalsAway) + " " + event.TeamAway + "\n" + event.Tournament + "\n"

				if event.PenHome != 0 || event.PenAway != 0 {
					message = message + "Penalties: " + strconv.Itoa(event.PenHome) + ":" + strconv.Itoa(event.PenAway)
				}

				err = SendMessageToTelegram(botToken, channelID, message)
				if err != nil {
					return "", err
				}

				result = result + "sent - " + event.TeamHome + " vs " + event.TeamAway + " | "
				event.PublishedStatus = "SENT"

			} else {

				result = result + "expired - " + event.TeamHome + " vs " + event.TeamAway + " | "
				event.PublishedStatus = "EXPIRED"

			}
			updateFinishedEvents = append(updateFinishedEvents, event)

		}

		err = storage.UpdateEvents(db, updateFinishedEvents)
		if err != nil {
			return "", fmt.Errorf("failed to update events in DB: %v", err)
		}
	}

	return result, nil
}
