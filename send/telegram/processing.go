package telegram

import (
	"database/sql"
	"fmt"
	"footballresult/config"
	"footballresult/storage"
	"footballresult/types"
	"strconv"
	"strings"
	"time"
)

func eventExpired(eventDate time.Time) bool {

	now := time.Now()

	duration := now.Sub(eventDate)

	return duration > 4*time.Hour

}

func GetLeagueEmoji(league string) string {
	// Приводим строку к нижнему регистру для удобного сравнения
	league = strings.ToLower(league)

	// Карта соответствий чемпионатов и эмодзи
	emojis := map[string]string{
		"serie a":          "🇮🇹",
		"premier league":   "🏴󠁧󠁢󠁥󠁮󠁧󠁿", // Флаг Англии
		"primera division": "🇪🇸",
	}

	// Если чемпионат найден — возвращаем эмодзи, иначе 🏆
	if emoji, found := emojis[league]; found {
		return emoji
	}
	return "🏆"
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
				emoji := GetLeagueEmoji(event.Tournament)
				score := "⚽️" + " <b>" + event.TeamHome + "</b> " + strconv.Itoa(event.GoalsHome) + " - " + strconv.Itoa(event.GoalsAway) + " <b>" + event.TeamAway + "</b> \n"
				tournament := emoji + " <b>" + event.Tournament + "</b> \n"
				message := score + tournament

				if event.PenHome != 0 || event.PenAway != 0 {
					penalties := message + "🥅 <b>Penalties:</b> " + strconv.Itoa(event.PenHome) + ":" + strconv.Itoa(event.PenAway) + "\n"
					message = score + penalties + tournament
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
