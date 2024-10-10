package storage

import (
	"database/sql"
	"fmt"
	"footballresult/types"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"os"
	"time"
)

func InitDB() (*sql.DB, error) {
	err := godotenv.Load(".env")
	if err != nil {
		return nil, fmt.Errorf("error loading .env file: %v", err)
	}

	host, err := LoadEnvVariable("DB_HOST")
	if err != nil {
		return nil, err
	}

	port, err := LoadEnvVariable("DB_PORT")
	if err != nil {
		return nil, err
	}

	user, err := LoadEnvVariable("DB_USER")
	if err != nil {
		return nil, err
	}

	password, err := LoadEnvVariable("DB_PASS")
	if err != nil {
		return nil, err
	}

	dbname, err := LoadEnvVariable("DB_NAME")
	if err != nil {
		return nil, err
	}

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

func LoadEnvVariable(key string) (string, error) {
	value := os.Getenv(key)
	if value == "" {
		return "", fmt.Errorf("%s not set in .env file", key)
	}
	return value, nil
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

func GetMinutesSinceLastAction(db *sql.DB, action string, status string) (int64, error) {
	query := `
		SELECT date
		FROM log
		WHERE action = $1 AND status = $2
		ORDER BY date DESC
		LIMIT 1;
	`

	var logDate time.Time

	err := db.QueryRow(query, action, status).Scan(&logDate)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("no log entries found for action: %s, status: %s", action, status)
		}

		return 0, fmt.Errorf("error executing query: %v", err)
	}

	now := time.Now()
	duration := now.Sub(logDate)
	minutesSinceLast := int64(duration.Minutes())

	return minutesSinceLast, nil
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
