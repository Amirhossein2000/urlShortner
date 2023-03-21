package handlers

import (
	"log"
	"net/http"
	"urlShortner/internal/help"
	"urlShortner/internal/ports"
)

type handlers struct {
	storage ports.Storage
	encoder ports.Encoder
}

func NewHandlers(s ports.Storage, e ports.Encoder) *handlers {
	return &handlers{
		storage: s,
		encoder: e,
	}
}

func (h *handlers) MiddleWare(method string, handlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handlerFunc = h.JsonResponseMiddleWare(handlerFunc)
		if r.Method != method {
			log.Println("MethodNotAllowed request, details:", help.HttpRequestDetails(r))
			help.SendJsonResponse(w, http.StatusMethodNotAllowed)
			return
		}
		handlerFunc(w, r)
	}
}

func (h *handlers) JsonResponseMiddleWare(handlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		handlerFunc(w, r)
	}
}
