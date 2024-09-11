package get

import (
	"database/sql"
	"footballresult/get/footballdata"
	"footballresult/types"
	"log"
)

func Events(db *sql.DB) error {
	token := footballdata.LoadFootballDataToken()

	teamIDs, err := footballdata.GetActiveTeamIDsFromDB(db)
	if err != nil {
		log.Fatal(err)
	}

	var allParsedEvents []types.Event
	for _, teamID := range teamIDs {

		url := footballdata.GetMatchesURL(teamID)
		getJson, _ := footballdata.GetJSON(url, token)
		teamParsedEvents, _ := footballdata.FilterTimedEvents(getJson)
		newParsedEvents := footballdata.CompareEvents(allParsedEvents, teamParsedEvents)
		allParsedEvents = append(allParsedEvents, newParsedEvents...)

	}

	eventsFromDB, _ := footballdata.GetTimedEventsFromDB(db)
	newEventsAdd := footballdata.CompareEvents(eventsFromDB, allParsedEvents)

	footballdata.InsertEventsInDB(db, newEventsAdd)

	return nil

}
