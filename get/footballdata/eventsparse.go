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

func GetTeamMatchesJSON(url string, authToken string) ([]byte, error) {
	// Создаем запрос
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Устанавливаем заголовок для авторизации
	req.Header.Set("X-Auth-Token", authToken)

	// Выполняем запрос
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Читаем ответ
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Возвращаем JSON как слайс байтов
	return body, nil
}

func FilterTimedEvents(apiResponse []byte) ([]types.Event, error) {
	var response types.FootbalDataResponse
	if err := json.Unmarshal(apiResponse, &response); err != nil {
		return nil, err
	}

	var events []types.Event
	for _, match := range response.Matches {
		if match.Status == "TIMED" {
			event := types.Event{
				EventID:    match.ID,
				EventDate:  match.UTCDate,
				Tournament: match.Tournament.Name,
				TeamHome:   match.HomeTeam.ShortName,
				TeamAway:   match.AwayTeam.ShortName,
				GoalsHome:  match.Score.FullTime.Home,
				GoalsAway:  match.Score.FullTime.Away,
				// Пропустим PenHome, PenAway, RcHome и RcAway, если эти данные недоступны
				EventStatus: match.Status,
			}
			events = append(events, event)
		}
	}

	return events, nil
}
