package storage

import (
	"database/sql"
	"fmt"
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
