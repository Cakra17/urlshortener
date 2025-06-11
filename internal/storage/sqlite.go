package storage

import (
	"database/sql"
	"log"
)

func InitDB(filepath string) *sql.DB {
  db, err := sql.Open("sqlite3", filepath)
	if err != nil {
		log.Fatal(err)
	}

	createTableQuery :=
		`CREATE TABLE IF NOT EXISTS tb_urls (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    short_url TEXT NOT NULL UNIQUE, 
    long_url TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
  );`

	_, err = db.Exec(createTableQuery)
	if err != nil {
		log.Fatal(err)
	}

  return db
}
