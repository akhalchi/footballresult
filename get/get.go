package get

import (
	"database/sql"
	"footballresult/get/footballdata"
	"footballresult/types"
	"log"
)

func Matches(db *sql.DB) error {
	token := footballdata.LoadFootballDataToken()

	teamIDs, err := footballdata.GetActiveTeamIDs(db)
	if err != nil {
		log.Fatal(err)
	}

	var allParsedEvents []types.Event
	for _, teamID := range teamIDs {

		url := footballdata.GetMatchesURL(teamID)
		getJson, _ := footballdata.GetTeamMatchesJSON(url, token)
		teamParsedEvents, _ := footballdata.FilterTimedEvents(getJson)
		newParsedEvents := footballdata.CompareEvents(allParsedEvents, teamParsedEvents)
		allParsedEvents = append(allParsedEvents, newParsedEvents...)

	}

	eventsFromDB, _ := footballdata.GetEventsFromDB(db)
	newEventsAdd := footballdata.CompareEvents(eventsFromDB, allParsedEvents)

	footballdata.InsertEventsInDB(db, newEventsAdd)

	/* url := footballdata.GetMatchesURL(65)
	getJson, _ := footballdata.GetTeamMatchesJSON(url, token)
	games, _ := footballdata.FilterTimedEvents(getJson)

	err := footballdata.InsertEventsInDB(db,games)

	if err != nil {
		return fmt.Errorf("failed to insert event: %v", err)
	}

	 for _, event := range allParsedEvents {
		fmt.Printf("Event ID: %d\n", event.EventID)
		fmt.Printf("Date: %s\n", event.EventDate.Format("2006-01-02 15:04"))
		fmt.Printf("Teams: %s vs %s\n", event.TeamHome, event.TeamAway)
		fmt.Printf("Tournament: %s\n", event.Tournament)
		fmt.Printf("Score: %d - %d\n", event.GoalsHome, event.GoalsAway)
		fmt.Printf("Status: %s\n", event.EventStatus)
		fmt.Println("----")
	} */

	return nil

}
