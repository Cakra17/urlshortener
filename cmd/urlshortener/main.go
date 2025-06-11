package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/cakra17/urlshortener/internal/handler"
	"github.com/cakra17/urlshortener/internal/storage"
	_ "github.com/mattn/go-sqlite3"
)

func main() { 
  db := storage.InitDB("./urls.db")
  defer db.Close()

  server := &handler.Server{
    DB: db,
  }

  mux := http.NewServeMux()

  mux.HandleFunc("POST /shorten", server.ShortenHandler)
  mux.HandleFunc("GET /{code}", server.RedirectHandler)

  fmt.Println("[INFO] Server is running on port 6969")
  log.Fatal(http.ListenAndServe(":6969", mux))
}
