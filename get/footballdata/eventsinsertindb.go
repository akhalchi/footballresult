package footballdata

import (
	"database/sql"
	"fmt"
	"footballresult/types"
)

func InsertEventsInDB(db *sql.DB, events []types.Event) error {
	query := `
    INSERT INTO events (event_id, event_date, event_tournament, team_home, team_away, goals_home, goals_away, pen_home, pen_away, rc_home, rc_away, importance, event_status, published_status)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
    ON CONFLICT (event_id) DO UPDATE
    SET event_date = EXCLUDED.event_date,
        event_tournament = EXCLUDED.event_tournament,
        team_home = EXCLUDED.team_home,
        team_away = EXCLUDED.team_away,
        goals_home = EXCLUDED.goals_home,
        goals_away = EXCLUDED.goals_away,
        pen_home = EXCLUDED.pen_home,
        pen_away = EXCLUDED.pen_away,
        rc_home = EXCLUDED.rc_home,
        rc_away = EXCLUDED.rc_away,
        importance = EXCLUDED.importance,
        event_status = EXCLUDED.event_status,
        published_status = EXCLUDED.published_status;
    `
	for _, event := range events {
		_, err := db.Exec(query,
			event.EventID,
			event.EventDate,
			event.Tournament,
			event.TeamHome,
			event.TeamAway,
			event.GoalsHome,
			event.GoalsAway,
			event.PenHome,
			event.PenAway,
			event.RcHome,
			event.RcAway,
			event.Importance,
			event.EventStatus,
			event.PublishedStatus,
		)
		if err != nil {
			return fmt.Errorf("failed to insert event %d: %v", event.EventID, err)
		}
	}

	return nil
}
