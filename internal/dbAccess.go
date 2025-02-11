package internal

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type DbAccess struct {
	db *sql.DB
}

func NewDbAccess() *DbAccess {
	db, err := sql.Open("sqlite3", "saveLocations.db")
	if err != nil {
		log.Fatalln(err)
	}
	return &DbAccess{
		db: db,
	}
}

func (access *DbAccess) GetGameMatchingPattern(pattern string) []string {
	const SELECTGAMES = `
		SELECT
			Game
		FROM Saves
		WHERE Game LIKE ?;`

	rows, err := access.db.Query(SELECTGAMES, "%"+pattern+"%")
	if err != nil {
		log.Fatalln(err)
	}
	defer rows.Close()

	games := make([]string, 0, 50)
	for rows.Next() {
		game := ""
		err := rows.Scan(&game)
		if err != nil {
			log.Fatalln(err)
		}
		games = append(games, game)
	}
	return games
}

func (access *DbAccess) Close() {
	access.db.Close()
}
