package footballdata

import (
	"database/sql"
	"fmt"
	"footballresult/types"
	"log"
)

func GetEvents(db *sql.DB, url, token string) error {

	teamIDs, err := GetActiveTeamIDsFromDB(db)
	if err != nil {
		return err
	}

	var allParsedEvents []types.Event
	for _, teamID := range teamIDs {

		url := GetEventsURL(url, teamID)
		getJson, _ := GetJSON(url, token)
		teamParsedEvents, _ := FilterTimedEvents(getJson)
		newParsedEvents := CheckExistEvents(allParsedEvents, teamParsedEvents)
		allParsedEvents = append(allParsedEvents, newParsedEvents...)

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
