package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/cakra17/urlshortener/internal/models"
)

type Server struct {
  DB *sql.DB
}

func (s *Server) ShortenHandler(w http.ResponseWriter, r *http.Request) {
	var body models.URL

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	stmt, err := s.DB.Prepare("INSERT INTO tb_urls (short_url, long_url) VALUES(?, ?);")
	if err != nil {
		http.Error(w, "Failed to Prepare statment", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(body.ShortURL, body.LongURL)
	if err != nil {
		http.Error(w, fmt.Sprintf("[ERROR] %s", err.Error()), http.StatusBadRequest)
		return
	}

	response := &models.SuccessRes{Status: "success", ShortURL: body.ShortURL}
	jsonByte, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonByte)
}

func (s *Server) RedirectHandler(w http.ResponseWriter, r *http.Request) {
	p := r.PathValue("code")

  stmt, err := s.DB.Prepare("SELECT long_url FROM tb_urls WHERE short_url = ?;")
	if err != nil {
		http.Error(w, "Failed to Prepare statment", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

  var longUrl string
  if err := stmt.QueryRow(p).Scan(&longUrl); err != nil {
    if err == sql.ErrNoRows {
      http.NotFound(w, r)
      return
    }
    http.Error(w, "Database Error", http.StatusInternalServerError)
  }

	http.Redirect(w, r, longUrl, http.StatusFound)
}
