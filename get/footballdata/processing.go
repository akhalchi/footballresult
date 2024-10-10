package footballdata

import (
	"database/sql"
	"fmt"
	"footballresult/storage"
	"footballresult/types"
	"github.com/joho/godotenv"
	"log"
	"strconv"
	"time"
)

func AddNewsEvents(db *sql.DB, url, token string) (error error, result string) {

	teamIDs, err := GetActiveTeamIDsFromDB(db)
	if err != nil {
		return err, ""
	}

	var allParsedEvents []types.Event
	for _, teamID := range teamIDs {

		url := GetEventsURL(url, teamID)
		getJson, err := GetJSON(url, token)
		if err != nil {
			return err, ""
		}
		teamParsedEvents, err := FilterTimedEvents(getJson)
		if err != nil {
			return err, ""
		}

		newParsedEvents := CheckExistEvents(allParsedEvents, teamParsedEvents)
		allParsedEvents = append(allParsedEvents, newParsedEvents...)

	}

	if len(allParsedEvents) == 0 {

		err = storage.InsertLog(db, "ADD EVENTS", "FAILED", "0 PARSED EVENTS")
		if err != nil {
			return err, ""
		}

		return fmt.Errorf("FAILED TO PARSE EVENTS"), ""
	}

	eventsFromDB, err := GetTimedEventsFromDB(db)
	if err != nil {
		return err, ""
	}

	newEventsAdd := CheckExistEvents(eventsFromDB, allParsedEvents)

	err = InsertEventsInDB(db, newEventsAdd)
	if err != nil {
		return err, ""
	}

	details := "PARSED: " + strconv.Itoa(len(allParsedEvents)) + ", ADDED: " + strconv.Itoa(len(newEventsAdd))
	err = storage.InsertLog(db, "ADD EVENTS", "SUCCESS", details)
	if err != nil {
		return err, ""
	}
	return nil, "new events: " + details

}

func UpdateActiveEvents(db *sql.DB, url, token string) (error error, result string) {

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
			return fmt.Errorf("failed to insert events into DB: %w", err), ""
		}
		result = "updated events" + strconv.Itoa(len(updatedEvents))
	}

	return nil, result

}

func timeToAddEvents(db *sql.DB) (bool, error) {
	action := "ADD EVENTS"

	status, minutes, err := storage.GetLastActionResult(db, action)
	if err != nil {
		if err.Error() == fmt.Sprintf("no log entries found for action: %s", action) {
			return true, nil
		}

		return false, fmt.Errorf("error checking last action: %v", err)
	}

	if status == "SUCCESS" {

		if minutes > 1440 {
			return true, nil
		}
	}

	if status == "FAILED" {

		if minutes > 30 {
			return true, nil
		}
	}

	return false, nil
}

func eventsProcessing(db *sql.DB, url, token string) (error error, result string) {

	timeToAddEvents, err := timeToAddEvents(db)
	if err != nil {
		return fmt.Errorf("get time to add events: %v", err), ""
	}

	if timeToAddEvents {
		err, result = AddNewsEvents(db, url, token)
		if err != nil {
			return fmt.Errorf("add new events: %v", err), ""
		} else if result != "" {
			return nil, result
		}
	} else {

		err, result = UpdateActiveEvents(db, url, token)
		if err != nil {
			return fmt.Errorf("update active events: %v", err), ""
		} else if result != "" {
			return nil, result
		}

	}
	return nil, ""

}

func Start(db *sql.DB) {

	err := godotenv.Load(".env")
	if err != nil {
		log.Printf("ERROR loading .env file: %v", err)
	} else {
		log.Printf("env file is loaded")
	}

	token, err := storage.LoadEnvVariable("FOOTBALL_DATA_TOKEN")
	if err != nil {
		log.Printf("ERROR loading FOOTBALL_DATA_TOKEN from .env file: %v", err)
	} else {
		log.Printf("variable FOOTBALL_DATA_TOKEN is loaded")
	}

	url, err := storage.LoadEnvVariable("FOOTBALL_DATA_URL")
	if err != nil {
		log.Printf("error loading FOOTBALL_DATA_URL from .env file: %v", err)
	} else {
		log.Printf("variable FOOTBALL_DATA_URL is loaded")
	}

	err, result := eventsProcessing(db, url, token)
	if err != nil {
		log.Printf("ERROR: %v", err)
	} else if result != "" {
		log.Printf(result)
	}

	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			err, result := eventsProcessing(db, url, token)
			if err != nil {
				log.Printf("ERROR: %v", err)
			} else if result != "" {
				log.Printf(result)
			}
		}
	}
}
