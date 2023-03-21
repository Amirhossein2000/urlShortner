package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"urlShortner/internal/help"
)

type ShortenRequest struct {
	Url       string `json:"url"`
	CustomUrl string `json:"custom_Url,omitempty"`
}

type ShortenResponse struct {
	ShortUrl string `json:"short_url"`
}

type StatsResponse struct {
	Original  string         `json:"original"`
	Shortened string         `json:"shortened"`
	Visits    map[string]int `json:"visits"`
}

func (h *handlers) Shorten(w http.ResponseWriter, r *http.Request) {
	req := &ShortenRequest{}
	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		help.SendJsonResponseWithCustomMessage(w, http.StatusBadRequest, "Invalid json struct")
		return
	}

	longURL := req.Url
	if !help.IsValidURL(longURL) {
		help.SendJsonResponseWithCustomMessage(w, http.StatusBadRequest, "Invalid URL")
		return
	}

	_, err = h.storage.GetShortURL(longURL)
	if err == nil {
		help.SendJsonResponseWithCustomMessage(w, http.StatusBadRequest, "URL already exist")
		return
	}
	shortURL := req.CustomUrl
	if shortURL == "" {
		repetitious := "repetitious"
		for repetitious != "" {
			shortURL = h.encoder.GenerateShortURL()
			repetitious, _ = h.storage.GetShortURL(longURL)
		}
	}

	err = h.storage.SetURLMapping(shortURL, longURL)
	if err != nil {
		help.SendJsonResponse(w, http.StatusInternalServerError)
		return
	}

	help.SendJsonResponseWithBody(w, http.StatusCreated, ShortenResponse{ShortUrl: shortURL})
}

func (h *handlers) Redirect(w http.ResponseWriter, r *http.Request) {
	split := strings.Split(r.URL.Path, "/")
	shortURL := split[1]

	longURL, err := h.storage.GetLongURL(shortURL)

	if err != nil {
		http.NotFound(w, r)
		return
	}
	h.storage.IncrementStats(longURL)

	http.Redirect(w, r, longURL, http.StatusPermanentRedirect)
}

func (h *handlers) Stats(w http.ResponseWriter, r *http.Request) {
	split := strings.Split(r.URL.Path, "/")

	if split[len(split)-1] != "stats" {
		help.SendJsonResponseWithCustomMessage(w, http.StatusBadRequest, "url should be like {url}/stats")
		return
	}

	shortURL := split[1]
	longURL, err := h.storage.GetLongURL(shortURL)
	if err != nil {
		help.SendJsonResponse(w, http.StatusNotFound)
		return
	}

	resp := StatsResponse{
		Original:  longURL,
		Shortened: shortURL,
		Visits:    h.storage.GetLastThreeDaysStats(longURL),
	}
	_ = json.NewEncoder(w).Encode(resp)
}

func (h *handlers) Main(w http.ResponseWriter, r *http.Request) {
	split := strings.Split(r.URL.Path, "/")
	if len(split) == 2 {
		h.Redirect(w, r)
		return
	} else if len(split) == 3 {
		h.MiddleWare(http.MethodGet, h.Stats)(w, r)
		return
	}

	http.NotFound(w, r)
}
