package footballdata

import (
	"database/sql"
	"fmt"
	"footballresult/storage"
	"footballresult/types"
	"github.com/joho/godotenv"
	"log"
	"strconv"
)

func AddNewsEvents(db *sql.DB, url, token string) error {

	teamIDs, err := GetActiveTeamIDsFromDB(db)
	if err != nil {
		return err
	}

	var allParsedEvents []types.Event
	for _, teamID := range teamIDs {

		url := GetEventsURL(url, teamID)
		getJson, err := GetJSON(url, token)
		if err != nil {
			return err
		}
		teamParsedEvents, err := FilterTimedEvents(getJson)
		if err != nil {
			return err
		}

		newParsedEvents := CheckExistEvents(allParsedEvents, teamParsedEvents)
		allParsedEvents = append(allParsedEvents, newParsedEvents...)

	}

	if len(allParsedEvents) == 0 {

		err = storage.InsertLog(db, "ADD EVENTS", "FAILED", "0 PARSED EVENTS")
		if err != nil {
			return err
		}

		return fmt.Errorf("FAILED TO PARSE EVENTS")
	}

	eventsFromDB, err := GetTimedEventsFromDB(db)
	if err != nil {
		return err
	}

	newEventsAdd := CheckExistEvents(eventsFromDB, allParsedEvents)

	err = InsertEventsInDB(db, newEventsAdd)
	if err != nil {
		return err
	}

	detailes := "PARSED: " + strconv.Itoa(len(allParsedEvents)) + ", ADDED: " + strconv.Itoa(len(newEventsAdd))
	err = storage.InsertLog(db, "ADD EVENTS", "SUCCESS", detailes)
	if err != nil {
		return err
	}
	return nil

}

func UpdateActiveEvents(db *sql.DB, url, token string) error {

	activeEvents, err := GetActiveEventsFromDB(db)
	if err != nil {
		log.Fatal(err)
	}

	var updatedEvents []types.Event

	for _, event := range activeEvents {
		eventURL := GetOneEventURL(url, int(event.EventID))
		getJson, _ := GetJSON(eventURL, token)
		parsedEvent, _ := ParseFootballEvent(getJson)

		if !CompareEvents(parsedEvent, event) {
			updatedEvents = append(updatedEvents, parsedEvent)

		}

	}

	if len(updatedEvents) > 0 {
		if err := InsertEventsInDB(db, updatedEvents); err != nil {
			return fmt.Errorf("failed to insert events into DB: %w", err)
		}
	}

	return nil

}

func timeToAddEvents(db *sql.DB) (bool, error) {
	minutesSinceSuccess, err := storage.GetMinutesSinceLastAction(db, "ADD EVENTS", "SUCCESS")
	if err != nil {
		return false, fmt.Errorf("error checking last SUCCESS action: %v", err)
	}

	if minutesSinceSuccess > 1440 {
		return true, nil
	}

	minutesSinceFailed, err := storage.GetMinutesSinceLastAction(db, "ADD EVENTS", "FAILED")
	if err != nil {
		return false, fmt.Errorf("error checking last FAILED action: %v", err)
	}

	if minutesSinceFailed > 30 {
		return true, nil
	}

	return false, nil
}

func Start(db *sql.DB) {

	err := godotenv.Load(".env")
	if err != nil {
		log.Printf("error loading .env file: %v", err)
	}

	token, err := storage.LoadEnvVariable("FOOTBALL_DATA_TOKEN")
	if err != nil {
		log.Printf("error loading FOOTBALL_DATA_TOKEN from .env file: %v", err)
	}

	url, err := storage.LoadEnvVariable("FOOTBALL_DATA_URL")
	if err != nil {
		log.Printf("error loading FOOTBALL_DATA_URL from .env file: %v", err)
	}

	timeToAddEvents, err := timeToAddEvents(db)

	if timeToAddEvents {
		err = AddNewsEvents(db, url, token)
		if err != nil {
			log.Printf("Add events error: %v", err)
		}
	} else {

		err = UpdateActiveEvents(db, url, token)
		if err != nil {
			log.Printf("Update Active events error: %v", err)
		}

	}

}
