package footballdata

import (
	"database/sql"
	"fmt"
)

func GetActiveTeamIDs(db *sql.DB) ([]int, error) {

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
