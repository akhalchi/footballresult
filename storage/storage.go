package storage

import (
	"database/sql"
	"fmt"
	"footballresult/types"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
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
