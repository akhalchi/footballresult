package footballdata

import (
	"database/sql"
	"fmt"
	"footballresult/config"
	"footballresult/logger"
	"footballresult/storage"
	"footballresult/types"
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

	err = storage.InsertEvents(db, newEventsAdd)
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

func UpdateActiveEvents(db *sql.DB, url, token string) (string, error) {

	activeEvents, err := GetActiveEventsFromDB(db)
	if err != nil {
		return "", fmt.Errorf("failed to get active events from DB: %w", err)
	}

	var updatedEvents []types.Event
	result := "updated events: "

	if len(activeEvents) > 0 {
		for _, event := range activeEvents {
			eventURL := GetOneEventURL(url, int(event.EventID))
			getJson, err := GetJSON(eventURL, token)
			if err != nil {
				return "", fmt.Errorf("failed to get JSON for event %d: %w", event.EventID, err)
			}

			parsedEvent, err := ParseFootballEvent(getJson)
			if err != nil {
				return "", fmt.Errorf("failed to parse event %d: %w", event.EventID, err)
			}

			if !CompareEvents(parsedEvent, event) {
				updatedEvents = append(updatedEvents, parsedEvent)
				result = result + parsedEvent.TeamHome + " vs " + parsedEvent.TeamAway + " | "

			}

		}
	}

	if len(updatedEvents) > 0 {
		if err := storage.UpdateEvents(db, updatedEvents); err != nil {
			return "", fmt.Errorf("failed to insert events into DB: %w", err)
		}
	} else {
		result = ""
	}

	return result, nil

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

		result, err = UpdateActiveEvents(db, url, token)
		if err != nil {
			return fmt.Errorf("update active events: %v", err), ""
		} else if result != "" {
			return nil, result
		}

	}
	return nil, ""

}

func Start(db *sql.DB) {

	token := config.Load.FootballDataToken
	url := config.Load.FootballDataURL

	logger.Send("parser is starting...")

	err, result := eventsProcessing(db, url, token)
	if err != nil {
		logger.Send(err.Error())
	} else if result != "" {
		logger.Send(result)
	}

	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()
	i := 1
	for {
		select {
		case <-ticker.C:
			err, result := eventsProcessing(db, url, token)
			if err != nil {
				logger.Send(err.Error())
			} else if result != "" {
				logger.Send(result)
				i = 0
			} else {
				i = i + 1
				if i > 60 {
					logger.Send("no updated events for the last hour")
					i = 0
				}
			}

		}
	}
}
