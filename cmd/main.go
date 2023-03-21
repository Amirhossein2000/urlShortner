package main

import (
	"net/http"
	"urlShortner/internal/handlers"
	"urlShortner/internal/repository/inMemory"
	"urlShortner/internal/repository/urlGen"
)

func main() {
	s := inMemory.NewInMemory()
	e := urlGen.NewURLGenerator()
	h := handlers.NewHandlers(s, e)
	c := GetConfig()

	http.HandleFunc("/shorten", h.MiddleWare(http.MethodPost, h.Shorten))
	http.HandleFunc("/", h.Main)
	http.ListenAndServe(c.Server.Address, nil)
}
