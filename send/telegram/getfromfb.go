package telegram

import (
	"database/sql"
	"footballresult/storage"
	"footballresult/types"
)

func GetFinishedEventsFromDB(db *sql.DB) ([]types.Event, error) {
	query := `
		SELECT event_id, event_date, event_tournament, team_home, team_away, goals_home, goals_away, pen_home, pen_away, rc_home, rc_away, importance, event_status, published_status
		FROM events
		WHERE event_status = 'FINISHED' 
		AND published_status = 'PLANNED';`
	return storage.GetEventsFromDB(db, query)
}
