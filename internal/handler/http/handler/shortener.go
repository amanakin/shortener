package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/amanakin/shortener/internal/service"
	"github.com/go-chi/chi"
	"golang.org/x/exp/slog"
)

const (
	setLink = "/setlink"
	getLink = "/getlink/{shortened}"
)

type ShortenerHandler struct {
	readLimit int64
	shortener service.Shortener
	logger    *slog.Logger
}

func NewShortener(logger *slog.Logger, shortener service.Shortener, readLimit int64) *ShortenerHandler {
	return &ShortenerHandler{
		shortener: shortener,
		readLimit: readLimit,
		logger:    logger,
	}
}

type errorHandleFunc func(w http.ResponseWriter, r *http.Request) error

func (h *ShortenerHandler) errorLogger(f errorHandleFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := f(w, r)
		if err != nil {
			h.logger.Error("shortener handler: %s", err)
		}
	}
}

func (h *ShortenerHandler) Register(r chi.Router) {
	r.Post(setLink, h.errorLogger(h.SetLink))
	r.Get(getLink, h.errorLogger(h.GetLink))
}

// SetLinkRequest is a request for setting link.
type SetLinkRequest struct {
	Original string `json:"url"`
}

// SetLinkResponse is a response for setting link.
// Original could differ from request, because it could be fixed (added schema, etc.).
type SetLinkResponse struct {
	Original  string `json:"original"`
	Shortened string `json:"shortened"`
	Created   bool   `json:"created"`
}

func (h *ShortenerHandler) SetLink(w http.ResponseWriter, r *http.Request) error {
	r.Body = http.MaxBytesReader(w, r.Body, h.readLimit)
	dec := json.NewDecoder(r.Body)

	var req SetLinkRequest
	err := dec.Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return fmt.Errorf("decode request: %w", err)
	}

	link, created, err := h.shortener.Shorten(r.Context(), req.Original)
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return fmt.Errorf("shorten: %w", err)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	resp := SetLinkResponse{
		Original:  link.OriginalURL,
		Shortened: link.ShortenedURL,
		Created:   created,
	}

	if created {
		w.WriteHeader(http.StatusCreated)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		// No make sense to write error header, because we already have done it :(
		return fmt.Errorf("encode response: %w", err)
	}

	return nil
}

type GetLinkResponse struct {
	Original string `json:"original"`
}

func (h *ShortenerHandler) GetLink(w http.ResponseWriter, r *http.Request) error {
	shortened := chi.URLParam(r, "shortened")

	original, err := h.shortener.Resolve(r.Context(), shortened)
	if errors.Is(err, service.ErrNotFound) {
		http.Error(w, "Not found", http.StatusNotFound)
		return fmt.Errorf("resolve: %w", err)
	}
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return fmt.Errorf("resolve: %w", err)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	resp := GetLinkResponse{
		Original: original,
	}

	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return fmt.Errorf("encode response: %w", err)
	}

	return nil
}
