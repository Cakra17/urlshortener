package models

type URL struct {
  ShortURL  string `json:"short_url"`
  LongURL   string `json:"long_url"`
}

type SuccessRes struct {
  Status    string  `json:"status"`
  ShortURL  string  `json:"shorturl"`
}
