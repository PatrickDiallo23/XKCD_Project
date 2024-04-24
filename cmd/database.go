package cmd

import (
	"database/sql"
	"fmt"
	"xkcdcomics/model"

	_ "github.com/mattn/go-sqlite3"
)

var Db *sql.DB

func SetupDatabase() (*sql.DB, error) {
	Db, err := sql.Open("sqlite3", "./xkcd_comics.db")
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}

	// Create table if not exists
	_, err = Db.Exec(`CREATE TABLE IF NOT EXISTS comics (
		num INTEGER PRIMARY KEY,
		safeTitle TEXT,
		img TEXT,
		title TEXT,
		transcript TEXT
	)`)
	if err != nil {
		return nil, fmt.Errorf("failed to create table: %v", err)
	}

	return Db, nil
}

func InsertComic(db *sql.DB, c model.Comic) error {
	_, err := db.Exec("INSERT INTO comics (num, safeTitle, img, title, transcript) VALUES (?, ?, ?, ?, ?)",
		c.Num, c.SafeTitle, c.Img, c.Title, c.Transcript)
	if err != nil {
		return fmt.Errorf("failed to insert comic %d into database: %v", c.Num, err)
	}
	return nil
}

func GetComicsFromDatabase(offset, limit int) ([]model.Comic, error) {
	var comics []model.Comic
	rows, err := Db.Query("SELECT num, safeTitle, img, title, transcript FROM comics ORDER BY num ASC LIMIT ?, ?", offset, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve comics from database: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var c model.Comic
		if err := rows.Scan(&c.Num, &c.SafeTitle, &c.Img, &c.Title, &c.Transcript); err != nil {
			return nil, fmt.Errorf("failed to scan comic row: %v", err)
		}
		comics = append(comics, c)
	}

	if len(comics) == 0 {
		// Fetch from XKCD API
		//sequential option
		// for i := offset + 1; i <= offset+limit; i++ {
		// 	c, err := GetComic(i)
		// 	if err != nil {
		// 		return nil, fmt.Errorf("failed to retrieve comic %d: %v", i, err)
		// 	}
		// 	err = InsertComic(Db, c)
		// 	if err != nil {
		// 		log.Printf("Failed to insert comic %d into database: %v", i, err)
		// 	}
		// 	comics = append(comics, c)
		// }

		//concurrently option
		comics, _ = RetrieveComicsConcurrently2(offset+1, offset+limit)
		for _, c := range comics {
			err = InsertComic(Db, c)
			if err != nil {
				return nil, fmt.Errorf("failed to insert comic %d into database: %v", c.Num, err)
			}
		}
	}

	return comics, nil
}
