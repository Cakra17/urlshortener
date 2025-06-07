package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"
)

type URL struct {
  LongURL string `json:"url"`
}

type SuccessRes struct {
  Status    string  `json:"status"`
  ShortURL  string  `json:"shorturl"`
}

var urlStore = make(map[string]string)

func generateShortURL() string {
  r := rand.New(rand.NewSource(time.Now().UnixNano()))
  res := ""
  alpha := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTU"
  for {
    if len(res) == 5 && urlStore[res] == "" {
      return res 
    } else if len(res) == 5 && urlStore[res] != "" {
      res = ""
    }
    random := r.Intn(len(alpha) - 1)
    res += string(alpha[random])
  }
}

func shortenHandler(w http.ResponseWriter, r *http.Request) {
  if r.Method != "POST"{
    http.Error(w, "", http.StatusBadRequest) 
    return
  }

  var body URL
  
  if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
    http.Error(w, err.Error(), http.StatusBadRequest)
    return
  }

  shortUrl := generateShortURL()
  urlStore[shortUrl] = body.LongURL

  response := &SuccessRes{ Status: "success", ShortURL: shortUrl}
  jsonByte, err := json.Marshal(response)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }

  w.Header().Set("Content-Type", "application/json")
  w.Write(jsonByte)
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
  if r.Method != "GET" {
    http.Error(w, "", http.StatusBadRequest)
    return
  }

  p := r.PathValue("code")
  longURL, found := urlStore[p]
  if !found {
    http.NotFound(w, r)
    return
  }

  http.Redirect(w, r, longURL, http.StatusFound)
}

func main() {
  http.HandleFunc("/shorten", shortenHandler)
  http.HandleFunc("/{code}", redirectHandler)

  fmt.Println("[INFO] Server is running on port 6969")
  log.Fatal(http.ListenAndServe(":6969", nil))
}
