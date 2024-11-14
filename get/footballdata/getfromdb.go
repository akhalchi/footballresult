package footballdata

import (
	"database/sql"
	"fmt"
	"footballresult/storage"
	"footballresult/types"
)

func GetActiveTeamIDsFromDB(db *sql.DB) ([]int, error) {

	query := `SELECT team_id FROM teams WHERE team_status = true`

	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %v", err)
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)

	var teamIDs []int

	for rows.Next() {
		var teamID int

		if err := rows.Scan(&teamID); err != nil {
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}
		teamIDs = append(teamIDs, teamID)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %v", err)
	}

	return teamIDs, nil
}

func GetTimedEventsFromDB(db *sql.DB) ([]types.Event, error) {
	query := `SELECT event_id, event_date, event_tournament, team_home, team_away, goals_home, goals_away, pen_home, pen_away, rc_home, rc_away, importance, event_status, published_status 
              FROM events 
              WHERE event_status = 'TIMED'`
	return storage.GetEventsFromDB(db, query)
}

func GetActiveEventsFromDB(db *sql.DB) ([]types.Event, error) {
	query := `
		SELECT event_id, event_date, event_tournament, team_home, team_away, goals_home, goals_away, pen_home, pen_away, rc_home, rc_away, importance, event_status, published_status
		FROM events
		WHERE event_status NOT IN ('FINISHED', 'POSTPONED', 'CANCELED')
		AND event_date < NOW();`

	return storage.GetEventsFromDB(db, query)

}
