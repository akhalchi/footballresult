package footballdata

import (
	"encoding/json"
	"fmt"
	"footballresult/types"
	"github.com/joho/godotenv"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func LoadFootballDataToken() string {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	token := os.Getenv("FOOTBALL_DATA_TOKEN")
	if token == "" {
		log.Fatalf("FOOTBALL_DATA_TOKEN not set in .env file")
	}

	return token
}

func GetMatchesURL(teamID int) string {
	baseURL := "https://api.football-data.org/v4/teams/"
	return fmt.Sprintf("%s%d/matches", baseURL, teamID)
}

func GetJSON(url string, authToken string) ([]byte, error) {

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-Auth-Token", authToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func FilterTimedEvents(apiResponse []byte) ([]types.Event, error) {
	var response types.EventsArrayJson
	if err := json.Unmarshal(apiResponse, &response); err != nil {
		return nil, err
	}

	var events []types.Event
	for _, match := range response.Matches {
		if match.Status == "TIMED" {
			event := types.Event{
				EventID:     match.ID,
				EventDate:   match.UTCDate,
				Tournament:  match.Tournament.Name,
				TeamHome:    match.HomeTeam.ShortName,
				TeamAway:    match.AwayTeam.ShortName,
				EventStatus: match.Status,
			}
			events = append(events, event)
		}
	}

	return events, nil
}

func ParseFootballEvent(jsonData []byte) (types.Event, error) {
	var response types.EventJson
	err := json.Unmarshal(jsonData, &response)
	if err != nil {
		return types.Event{}, err
	}

	event := types.Event{
		EventID:     response.ID,
		EventDate:   response.UTCDate,
		Tournament:  response.Tournament.Name,
		TeamHome:    response.HomeTeam.ShortName,
		TeamAway:    response.AwayTeam.ShortName,
		EventStatus: response.Status,
	}

	if response.Score.Duration == "EXTRA_TIME" || response.Score.Duration == "PENALTY_SHOOTOUT" {
		event.GoalsHome = response.Score.RegularTime.HomeTeam + response.Score.ExtraTime.HomeTeam
		event.GoalsAway = response.Score.RegularTime.AwayTeam + response.Score.ExtraTime.AwayTeam
	} else {
		event.GoalsHome = response.Score.RegularTime.HomeTeam
		event.GoalsAway = response.Score.RegularTime.AwayTeam
	}

	if response.Score.Duration == "PENALTY_SHOOTOUT" {
		event.PenHome = response.Score.Penalties.HomeTeam
		event.PenAway = response.Score.Penalties.AwayTeam
	}

	return event, nil
}
