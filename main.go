package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

type URL struct {
  ShortURL  string `json:"short_url"`
  LongURL   string `json:"long_url"`
}

type SuccessRes struct {
  Status    string  `json:"status"`
  ShortURL  string  `json:"shorturl"`
}

const (
  alphabet  = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
  base      = float64(len(alphabet))
)

var urlStore = make(map[string]string)
var db *sql.DB

func InitDB() {
  var err error
  db, err = sql.Open("sqlite3", "./urls.db")
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
}

func loadData() error {
  clear(urlStore)

  rows, err := db.Query(`SELECT short_url, long_url FROM tb_urls;`) 
  if err != nil {
    return errors.New("Fail to fetch recent data")
  }
  defer rows.Close()

  for rows.Next() { 
    var url URL
    if err := rows.Scan(&url.ShortURL, &url.LongURL); err != nil {
      return errors.New("fail to load recent data")
    }
    urlStore[url.ShortURL] = url.LongURL
  }
  return nil
}

func shortenHandler(w http.ResponseWriter, r *http.Request) {
  var body URL
  
  if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
    http.Error(w, err.Error(), http.StatusBadRequest)
    return
  }

  stmt, err := db.Prepare("INSERT INTO tb_urls (short_url, long_url) VALUES(?, ?);")
  if err != nil {
    http.Error(w, "Failed to Prepare statment", http.StatusInternalServerError)
    return
  }
  defer stmt.Close()

  _ , err = stmt.Exec(body.ShortURL, body.LongURL)
  if err != nil {
    http.Error(w, fmt.Sprintf("[ERROR] %s", err.Error()), http.StatusBadRequest)
    return
  }

  urlStore[body.ShortURL] = body.LongURL

  response := &SuccessRes{ Status: "success", ShortURL: body.ShortURL}
  jsonByte, err := json.Marshal(response)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }

  w.Header().Set("Content-Type", "application/json")
  w.Write(jsonByte)
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
  p := r.PathValue("code")
  url, found := urlStore[p] 
  if !found {
    http.NotFound(w, r)
    return
  }

  http.Redirect(w, r, url, http.StatusFound)
}

func main() { 
  InitDB()
  err := loadData()
  if err != nil {
    fmt.Println(err.Error())
  }

  mux := http.NewServeMux()

  mux.HandleFunc("POST /shorten", shortenHandler)
  mux.HandleFunc("GET /{code}", redirectHandler)

  fmt.Println("[INFO] Server is running on port 6969")
  log.Fatal(http.ListenAndServe(":6969", mux))
}
