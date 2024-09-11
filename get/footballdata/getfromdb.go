package footballdata

import (
	"database/sql"
	"fmt"
	"footballresult/types"
)

func GetActiveTeamIDsFromDB(db *sql.DB) ([]int, error) {

	query := `SELECT team_id FROM teams WHERE team_status = true`

	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %v", err)
	}
	defer rows.Close()

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

	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %v", err)
	}
	defer rows.Close()

	var events []types.Event

	for rows.Next() {
		var event types.Event

		if err := rows.Scan(
			&event.EventID,
			&event.EventDate,
			&event.Tournament,
			&event.TeamHome,
			&event.TeamAway,
			&event.GoalsHome,
			&event.GoalsAway,
			&event.PenHome,
			&event.PenAway,
			&event.RcHome,
			&event.RcAway,
			&event.Importance,
			&event.EventStatus,
			&event.PublishedStatus,
		); err != nil {
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}
		events = append(events, event)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %v", err)
	}

	return events, nil
}

func GetActiveEventsFromDB(db *sql.DB) ([]types.Event, error) {
	query := `
		SELECT event_id, event_date, event_tournament, team_home, team_away, goals_home, goals_away, pen_home, pen_away, rc_home, rc_away, importance, event_status, published_status
		FROM events
		WHERE event_status NOT IN ('FINISHED', 'PASTRONED', 'CANCELED')
		AND event_date < NOW();
	`

	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %v", err)
	}
	defer rows.Close()

	var events []types.Event
	for rows.Next() {
		var event types.Event
		err := rows.Scan(
			&event.EventID,
			&event.EventDate,
			&event.Tournament,
			&event.TeamHome,
			&event.TeamAway,
			&event.GoalsHome,
			&event.GoalsAway,
			&event.PenHome,
			&event.PenAway,
			&event.RcHome,
			&event.RcAway,
			&event.Importance,
			&event.EventStatus,
			&event.PublishedStatus,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}
		events = append(events, event)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error after iterating rows: %v", err)
	}

	return events, nil
}
