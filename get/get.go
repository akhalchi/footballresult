package get

import (
	"database/sql"
	"fmt"
	"footballresult/get/footballdata"
	"footballresult/storage"
	"github.com/joho/godotenv"
	"time"
)

func Events(db *sql.DB) error {

	err := godotenv.Load(".env")
	if err != nil {
		return fmt.Errorf("error loading .env file: %v", err)
	}

	token, err := storage.LoadEnvVariable("FOOTBALL_DATA_TOKEN")
	if err != nil {
		return err
	}

	url, err := storage.LoadEnvVariable("FOOTBALL_DATA_URL")
	if err != nil {
		return err
	}

	err = footballdata.GetEvents(db, url, token)
	if err != nil {
		return err
	}

	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			err := footballdata.UpdateActiveEvents(db, url, token)
			if err != nil {
				return err
			}
		}
	}

}
