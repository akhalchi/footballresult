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

	// Получаем завершенные события из БД
	finishedEvents, err := GetFinishedEventsFromDB(db)
	if err != nil {
		return "", err
	}

	var updateFinishedEvents []types.Event
	var result string

	if len(finishedEvents) > 0 {
		for _, event := range finishedEvents {
			// Проверяем, не истекло ли событие
			if !eventExpired(event.EventDate) {
				// Формируем сообщение
				message := event.TeamHome + " " + strconv.Itoa(event.GoalsHome) + " - " + strconv.Itoa(event.GoalsAway) + " " + event.TeamAway + "\n" + event.Tournament + "\n"

				// Отправляем сообщение в Telegram
				err = SendMessageToTelegram(botToken, channelID, message)
				if err != nil {
					return "", err
				}
				result = "event sent to telegram: " + strconv.FormatInt(event.EventID, 10)

				// Если отправка успешна, обновляем статус события
				updateFinishedEvents = append(updateFinishedEvents, types.Event{
					EventID:         event.EventID,
					PublishedStatus: "SENT",
				})
			} else {
				// Если событие истекло, обновляем статус как "EXPIRED"
				updateFinishedEvents = append(updateFinishedEvents, types.Event{
					EventID:         event.EventID,
					PublishedStatus: "EXPIRED",
				})
				result = "event expired: " + strconv.FormatInt(event.EventID, 10)
			}

		}

		// Вставляем обновленные данные в БД
		err = storage.InsertUpdateEventsInDB(db, updateFinishedEvents)
		if err != nil {
			return "", fmt.Errorf("failed to update events in DB: %v", err)
		}
	}

	// Возвращаем результат
	return result, nil
}
