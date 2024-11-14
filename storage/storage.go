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

func InsertEvents(db *sql.DB, events []types.Event) error {
	query := `
		INSERT INTO events (event_id, event_date, event_tournament, team_home, team_away, goals_home, goals_away, pen_home, pen_away, rc_home, rc_away, importance, event_status, published_status)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
		ON CONFLICT (event_id) DO NOTHING; -- Конфликты игнорируем, если запись уже существует
	`

	for _, event := range events {
		_, err := db.Exec(query,
			event.EventID,
			event.EventDate,
			event.Tournament,
			event.TeamHome,
			event.TeamAway,
			event.GoalsHome,
			event.GoalsAway,
			event.PenHome,
			event.PenAway,
			event.RcHome,
			event.RcAway,
			event.Importance,
			event.EventStatus,
			event.PublishedStatus,
		)
		if err != nil {
			return fmt.Errorf("failed to insert event %d: %v", event.EventID, err)
		}
	}

	return nil
}

func UpdateEvents(db *sql.DB, events []types.Event) error {
	// SQL-запрос для обновления всех полей
	query := `
		UPDATE events
		SET event_date = $2,
			event_tournament = $3,
			team_home = $4,
			team_away = $5,
			goals_home = $6,
			goals_away = $7,
			pen_home = $8,
			pen_away = $9,
			rc_home = $10,
			rc_away = $11,
			importance = $12,
			event_status = $13,
			published_status = $14
		WHERE event_id = $1;
	`

	for _, event := range events {
		// Выполнение SQL-запроса с передачей всех значений из структуры
		_, err := db.Exec(query,
			event.EventID,         // $1
			event.EventDate,       // $2
			event.Tournament,      // $3
			event.TeamHome,        // $4
			event.TeamAway,        // $5
			event.GoalsHome,       // $6
			event.GoalsAway,       // $7
			event.PenHome,         // $8
			event.PenAway,         // $9
			event.RcHome,          // $10
			event.RcAway,          // $11
			event.Importance,      // $12
			event.EventStatus,     // $13
			event.PublishedStatus, // $14
		)
		if err != nil {
			return fmt.Errorf("failed to update event %d: %v", event.EventID, err)
		}
	}

	return nil
}
