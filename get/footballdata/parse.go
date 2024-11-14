package footballdata

import (
	"encoding/json"
	"fmt"
	"footballresult/types"
	"io"
	"io/ioutil"
	"net/http"
)

func GetEventsURL(url string, teamID int) string {
	baseURL := url + "/teams/"
	return fmt.Sprintf("%s%d/matches", baseURL, teamID)
}

func GetOneEventURL(url string, eventID int) string {
	baseURL := url + "/matches/"
	return fmt.Sprintf("%s%d", baseURL, eventID)
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
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

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
		if match.Status == "TIMED" ||
			match.Status == "IN_PLAY" ||
			match.Status == "PAUSED" ||
			match.Status == "SUSPENDED" {
			event := types.Event{
				EventID:         match.ID,
				EventDate:       match.UTCDate,
				Tournament:      match.Tournament.Name,
				TeamHome:        match.HomeTeam.ShortName,
				TeamAway:        match.AwayTeam.ShortName,
				GoalsHome:       0,
				GoalsAway:       0,
				PenHome:         0,
				PenAway:         0,
				RcHome:          0,
				RcAway:          0,
				EventStatus:     match.Status,
				PublishedStatus: "PLANNED",
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
		EventID:         response.ID,
		EventDate:       response.UTCDate,
		Tournament:      response.Tournament.Name,
		TeamHome:        response.HomeTeam.ShortName,
		TeamAway:        response.AwayTeam.ShortName,
		EventStatus:     response.Status,
		PublishedStatus: "PLANNED",
	}

	if response.Score.Duration == "EXTRA_TIME" || response.Score.Duration == "PENALTY_SHOOTOUT" {
		event.GoalsHome = response.Score.RegularTime.HomeTeam + response.Score.ExtraTime.HomeTeam
		event.GoalsAway = response.Score.RegularTime.AwayTeam + response.Score.ExtraTime.AwayTeam
	} else {
		event.GoalsHome = response.Score.FullTime.Home
		event.GoalsAway = response.Score.FullTime.Away
	}

	if response.Score.Duration == "PENALTY_SHOOTOUT" {
		event.PenHome = response.Score.Penalties.HomeTeam
		event.PenAway = response.Score.Penalties.AwayTeam
	}

	return event, nil
}
