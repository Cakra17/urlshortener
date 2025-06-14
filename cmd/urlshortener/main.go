package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/cakra17/urlshortener/internal/handler"
	"github.com/cakra17/urlshortener/internal/storage"
	_ "github.com/mattn/go-sqlite3"
)

func main() { 
  
  port := os.Getenv("URLSHORT_PORT")
	if port == "" {
		port = "6969"
	}

	dbPath := os.Getenv("URLSHORT_DB_PATH")
	if dbPath == "" {
		dbPath = "./urls.db"
	}

  db := storage.InitDB(dbPath)
  defer db.Close()
  
  logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

  server := &handler.Server{
    DB: db,
    Logger: logger,
  }

  mux := http.NewServeMux()

  mux.HandleFunc("POST /shorten", server.ShortenHandler)
  mux.HandleFunc("GET /{code}", server.RedirectHandler)

  logger.Info("[INFO] Server is running on port 6969", "port", port)
  log.Fatal(http.ListenAndServe(":6969", mux))
}
