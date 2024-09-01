package footballdata

import (
	"database/sql"
	"fmt"
)

// GetActiveTeamIDs возвращает массив идентификаторов команд, у которых статус true
func GetActiveTeamIDs(db *sql.DB) ([]int, error) {
	// SQL-запрос для выборки team_id команд со статусом true
	query := `SELECT team_id FROM teams WHERE team_status = true`

	// Выполнение запроса
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %v", err)
	}
	defer rows.Close()

	var teamIDs []int

	// Проход по всем строкам результата
	for rows.Next() {
		var teamID int
		// Сканирование значения team_id из строки
		if err := rows.Scan(&teamID); err != nil {
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}
		teamIDs = append(teamIDs, teamID)
	}

	// Проверка на наличие ошибок при итерации по строкам
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %v", err)
	}

	return teamIDs, nil
}
