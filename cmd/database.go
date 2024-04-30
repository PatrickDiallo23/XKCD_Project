package cmd

import (
	"database/sql"
	"fmt"
	"xkcdcomics/model"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func SetupDatabase() error {
	var err error
	db, err = sql.Open("sqlite3", "./xkcd_comics.db")
	if err != nil {
		return fmt.Errorf("failed to open database: %v", err)
	}

	// Create table if not exists
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS comics (
		num INTEGER PRIMARY KEY,
		safeTitle TEXT,
		img TEXT,
		title TEXT,
		transcript TEXT
	)`)
	if err != nil {
		return fmt.Errorf("failed to create table: %v", err)
	}

	return nil
}

func CloseDBConn() {
	db.Close()
}

func InsertComic(db *sql.DB, c model.Comic) error {
	_, err := db.Exec("INSERT INTO comics (num, safeTitle, img, title, transcript) VALUES (?, ?, ?, ?, ?)",
		c.Num, c.SafeTitle, c.Img, c.Title, c.Transcript)
	if err != nil {
		return fmt.Errorf("failed to insert comic %d into database: %v", c.Num, err)
	}
	return nil
}

func ComicExists(db *sql.DB, num int) (bool, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM comics WHERE num = ?", num).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check if comic %d exists in database: %v", num, err)
	}
	return count > 0, nil
}

func InsertComics(db *sql.DB, comics []model.Comic) error {
	// Deduplicate comics based on Num
	comicMap := make(map[int]model.Comic)
	for _, comic := range comics {
		comicMap[comic.Num] = comic
	}

	// Validate and deduplicate comics before insertion
	var validComics []model.Comic
	for _, comic := range comicMap {
		// Check if the comic already exists in the database
		exists, err := ComicExists(db, comic.Num)
		if err != nil {
			fmt.Printf("Error checking if comic %d exists in database: %v", comic.Num, err)
			continue
		}
		if !exists {
			validComics = append(validComics, comic)
		}
	}

	// Prepare the SQL statement
	sqlStr := "INSERT INTO comics (num, safeTitle, img, title, transcript) VALUES "
	var vals []interface{}

	// Iterate over the comics and build the SQL statement and values
	for _, comic := range validComics {
		sqlStr += "(?, ?, ?, ?, ?),"
		vals = append(vals, comic.Num, comic.SafeTitle, comic.Img, comic.Title, comic.Transcript)
	}

	// Trim the last comma from the SQL statement
	sqlStr = sqlStr[:len(sqlStr)-1]

	// Prepare the statement
	stmt, err := db.Prepare(sqlStr)
	if err != nil {
		return fmt.Errorf("failed to prepare SQL statement: %v", err)
	}
	defer stmt.Close()

	// Execute the statement with the values
	_, err = stmt.Exec(vals...)
	if err != nil {
		return fmt.Errorf("failed to insert comics into database: %v", err)
	}

	return nil
}

func GetComicsFromDatabase(offset, limit int) ([]model.Comic, error) {
	var comics []model.Comic
	rows, err := db.Query("SELECT num, safeTitle, img, title, transcript FROM comics ORDER BY num ASC LIMIT ?, ?", offset, limit)
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
		go InsertComics(db, comics)
	}

	return comics, nil
}
