package storage

import (
	"database/sql"
	"fmt"
	"footballresult/config"
	"footballresult/types"
	_ "github.com/lib/pq"
	"time"
)

func InitDB() (*sql.DB, error) {

	host := config.Load.DBHost
	dbname := config.Load.DBName
	port := config.Load.DBPort
	user := config.Load.DBUser
	password := config.Load.DBPass

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}

func GetEventsFromDB(db *sql.DB, query string) ([]types.Event, error) {

	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %v", err)
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

			fmt.Println("error closing rows:", err)
		}
	}(rows)

	var events []types.Event

	for rows.Next() {
		var event types.Event

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

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %v", err)
	}

	return events, nil
}

func InsertLog(db *sql.DB, action, status, details string) error {

	query := `INSERT INTO log (action, status, details) VALUES ($1, $2, $3)`

	_, err := db.Exec(query, action, status, details)
	if err != nil {
		return fmt.Errorf("failed to insert log entry: %v", err)
	}

	return nil
}

func GetLastActionResult(db *sql.DB, action string) (string, int64, error) {
	query := `
		SELECT status, date
		FROM log
		WHERE action = $1
		ORDER BY date DESC
		LIMIT 1;
	`

	var status string
	var logDate time.Time

	err := db.QueryRow(query, action).Scan(&status, &logDate)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", 0, fmt.Errorf("no log entries found for action: %s", action)
		}
		return "", 0, fmt.Errorf("error executing query: %v", err)
	}

	now := time.Now()
	duration := now.Sub(logDate)
	minutesSinceLast := int64(duration.Minutes())

	return status, minutesSinceLast, nil
}

func InsertUpdateEventsInDB(db *sql.DB, events []types.Event) error {

	for _, event := range events {
		query := "INSERT" + " INTO events (event_id"
		values := "VALUES ($1"
		params := []interface{}{event.EventID}
		paramIdx := 2

		// Динамически добавляем только непустые или значимые поля в запрос
		if !event.EventDate.IsZero() {
			query += ", event_date"
			values += fmt.Sprintf(", $%d", paramIdx)
			params = append(params, event.EventDate)
			paramIdx++
		}
		if event.Tournament != "" {
			query += ", event_tournament"
			values += fmt.Sprintf(", $%d", paramIdx)
			params = append(params, event.Tournament)
			paramIdx++
		}
		if event.TeamHome != "" {
			query += ", team_home"
			values += fmt.Sprintf(", $%d", paramIdx)
			params = append(params, event.TeamHome)
			paramIdx++
		}
		if event.TeamAway != "" {
			query += ", team_away"
			values += fmt.Sprintf(", $%d", paramIdx)
			params = append(params, event.TeamAway)
			paramIdx++
		}
		if event.GoalsHome != 0 {
			query += ", goals_home"
			values += fmt.Sprintf(", $%d", paramIdx)
			params = append(params, event.GoalsHome)
			paramIdx++
		}
		if event.GoalsAway != 0 {
			query += ", goals_away"
			values += fmt.Sprintf(", $%d", paramIdx)
			params = append(params, event.GoalsAway)
			paramIdx++
		}
		if event.PenHome != 0 {
			query += ", pen_home"
			values += fmt.Sprintf(", $%d", paramIdx)
			params = append(params, event.PenHome)
			paramIdx++
		}
		if event.PenAway != 0 {
			query += ", pen_away"
			values += fmt.Sprintf(", $%d", paramIdx)
			params = append(params, event.PenAway)
			paramIdx++
		}
		if event.RcHome != 0 {
			query += ", rc_home"
			values += fmt.Sprintf(", $%d", paramIdx)
			params = append(params, event.RcHome)
			paramIdx++
		}
		if event.RcAway != 0 {
			query += ", rc_away"
			values += fmt.Sprintf(", $%d", paramIdx)
			params = append(params, event.RcAway)
			paramIdx++
		}
		// Importance - булевое поле, его всегда добавляем
		query += "," + " importance"
		values += fmt.Sprintf(", $%d", paramIdx)
		params = append(params, event.Importance)
		paramIdx++

		if event.EventStatus != "" {
			query += ", event_status"
			values += fmt.Sprintf(", $%d", paramIdx)
			params = append(params, event.EventStatus)
			paramIdx++
		}
		if event.PublishedStatus != "" {
			query += ", published_status"
			values += fmt.Sprintf(", $%d", paramIdx)
			params = append(params, event.PublishedStatus)
			paramIdx++
		}

		query += ") " + values + ") ON CONFLICT (event_id) DO UPDATE SET "

		// Динамически добавляем только непустые поля в часть SET
		updateFields := ""
		if !event.EventDate.IsZero() {
			updateFields += "event_date = EXCLUDED.event_date, "
		}
		if event.Tournament != "" {
			updateFields += "event_tournament = EXCLUDED.event_tournament, "
		}
		if event.TeamHome != "" {
			updateFields += "team_home = EXCLUDED.team_home, "
		}
		if event.TeamAway != "" {
			updateFields += "team_away = EXCLUDED.team_away, "
		}
		if event.GoalsHome != 0 {
			updateFields += "goals_home = EXCLUDED.goals_home, "
		}
		if event.GoalsAway != 0 {
			updateFields += "goals_away = EXCLUDED.goals_away, "
		}
		if event.PenHome != 0 {
			updateFields += "pen_home = EXCLUDED.pen_home, "
		}
		if event.PenAway != 0 {
			updateFields += "pen_away = EXCLUDED.pen_away, "
		}
		if event.RcHome != 0 {
			updateFields += "rc_home = EXCLUDED.rc_home, "
		}
		if event.RcAway != 0 {
			updateFields += "rc_away = EXCLUDED.rc_away, "
		}
		updateFields += "importance = EXCLUDED.importance, " // always update importance
		if event.EventStatus != "" {
			updateFields += "event_status = EXCLUDED.event_status, "
		}
		if event.PublishedStatus != "" {
			updateFields += "published_status = EXCLUDED.published_status, "
		}

		// Удаляем последнюю запятую и пробел
		if len(updateFields) > 0 {
			query += updateFields[:len(updateFields)-2]
		}

		// Выполнение запроса
		_, err := db.Exec(query, params...)
		if err != nil {
			return fmt.Errorf("failed to insert/update event %d: %v", event.EventID, err)
		}
	}

	return nil
}
