package get

import (
	"fmt"
	"footballresult/get/footballdata"
)

func Matches() {
	token := footballdata.LoadFootballDataToken()
	url := footballdata.GetMatchesURL(65)
	getJson, _ := footballdata.GetTeamMatchesJSON(url, token)
	games, _ := footballdata.FilterTimedEvents(getJson)

	for _, event := range games {
		fmt.Printf("Event ID: %d\n", event.EventID)
		fmt.Printf("Date: %s\n", event.EventDate.Format("2006-01-02 15:04"))
		fmt.Printf("Teams: %s vs %s\n", event.TeamHome, event.TeamAway)
		fmt.Printf("Tournament: %s\n", event.Tournament)
		fmt.Printf("Score: %d - %d\n", event.GoalsHome, event.GoalsAway)
		fmt.Printf("Status: %s\n", event.EventStatus)
		fmt.Println("----")
	}

}
