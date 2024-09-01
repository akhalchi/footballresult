package footballdata

import (
	"database/sql"
	"fmt"
	"footballresult/types"
)

func GetEventsFromDB(db *sql.DB) ([]types.Event, error) {
	// SQL-запрос для выборки всех событий в статусе TIMED
	query := `SELECT event_id, event_date, event_tournament, team_home, team_away, goals_home, goals_away, pen_home, pen_away, rc_home, rc_away, importance, event_status, published_status 
              FROM events 
              WHERE event_status = 'TIMED'`

	// Выполнение запроса
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %v", err)
	}
	defer rows.Close()

	var events []types.Event

	// Проход по всем строкам результата
	for rows.Next() {
		var event types.Event
		// Сканирование значений из строки в структуру Event
		if err := rows.Scan(
			&event.EventID,
			&event.EventDate,
			&event.Tournament,
			&event.TeamHome,
			&event.TeamAway,
			&event.GoalsHome,
			&event.GoalsAway,
			&event.PenHome,
			&event.PenAway,
			&event.RcHome,
			&event.RcAway,
			&event.Importance,
			&event.EventStatus,
			&event.PublishedStatus,
		); err != nil {
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}
		events = append(events, event)
	}

	// Проверка на наличие ошибок при итерации по строкам
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %v", err)
	}

	return events, nil
}
