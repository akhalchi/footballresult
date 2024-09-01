package types

import (
	"time"
)

type Event struct {
	EventID         int64     `db:"event_id"`         // PRIMARY KEY
	EventDate       time.Time `db:"event_date"`       // TIMESTAMP WITH TIME ZONE
	Tournament      string    `db:"event_tournament"` // TEXT
	TeamHome        string    `db:"team_home"`        // TEXT
	TeamAway        string    `db:"team_away"`        // TEXT
	GoalsHome       int       `db:"goals_home"`       // INTEGER
	GoalsAway       int       `db:"goals_away"`       // INTEGER
	PenHome         int       `db:"pen_home"`         // INTEGER
	PenAway         int       `db:"pen_away"`         // INTEGER
	RcHome          int       `db:"rc_home"`          // INTEGER
	RcAway          int       `db:"rc_away"`          // INTEGER
	Importance      bool      `db:"importance"`       // BOOLEAN
	EventStatus     string    `db:"event_status"`     // TEXT
	PublishedStatus string    `db:"published_status"` // TEXT
}

type FootbalDataResponse struct {
	Matches []struct {
		ID         int64     `json:"id"`
		UTCDate    time.Time `json:"utcDate"`
		Status     string    `json:"status"`
		Tournament struct {
			Name string `json:"name"`
		} `json:"competition"`
		HomeTeam struct {
			ShortName string `json:"shortName"`
		} `json:"homeTeam"`
		AwayTeam struct {
			ShortName string `json:"shortName"`
		} `json:"awayTeam"`
		Score struct {
			FullTime struct {
				Home int `json:"home"`
				Away int `json:"away"`
			} `json:"fullTime"`
		} `json:"score"`
	} `json:"matches"`
}
