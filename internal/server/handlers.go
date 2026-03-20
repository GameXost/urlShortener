package server

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"net/url"
	"url_shortener/internal/storage"
)

type Request struct {
	URL string `json:"url_to_save"`
}
type Response struct {
	Alias string `json:"alias"`
}

type URLService interface {
	SaveURL(ctx context.Context, url string) (string, error)
	GetURL(ctx context.Context, alias string) (string, error)
}

type Handler struct {
	Service URLService
}

func NewHandler(srv URLService) *Handler {
	return &Handler{Service: srv}
}

func (h *Handler) SaveURL(w http.ResponseWriter, r *http.Request) {
	var req Request

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	if req.URL == "" {
		http.Error(w, "need url to save", http.StatusBadRequest)
		return
	}

	_, err = url.ParseRequestURI(req.URL)
	if err != nil {
		http.Error(w, "invalid url", http.StatusBadRequest)
		return
	}

	alias, err := h.Service.SaveURL(r.Context(), req.URL)
	if err != nil {
		http.Error(w, "failed to save url", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(Response{Alias: alias})
	if err != nil {
		log.Printf("failed to response: %v", err)
	}

}

func (h *Handler) GetURL(w http.ResponseWriter, r *http.Request) {
	alias := chi.URLParam(r, "alias")
	if alias == "" {
		http.Error(w, "alias is required", http.StatusBadRequest)
		return
	}

	resURL, err := h.Service.GetURL(r.Context(), alias)
	if err != nil {
		if errors.Is(err, storage.ErrURLNotFound) {
			http.Error(w, " url not found in storage", http.StatusNotFound)
			return
		}
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, resURL, http.StatusFound)
}
