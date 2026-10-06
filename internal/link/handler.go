package link

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

type Handler struct {
	store Store
}

func NewHandler(store Store) *Handler {
	return &Handler{store: store}
}

func (h *Handler) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", h.health)
	mux.Handle("GET /metrics", promhttp.Handler())
	mux.HandleFunc("POST /api/links", h.create)
	mux.HandleFunc("GET /{code}", h.redirect)
	return instrument(mux)
}

func (h *Handler) health(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}

type createRequest struct {
	URL string `json:"url"`
}

type createResponse struct {
	Code string `json:"code"`
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	var req createRequest
	body := http.MaxBytesReader(w, r.Body, 1<<20)
	if err := json.NewDecoder(body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	u, err := url.ParseRequestURI(req.URL)
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") || u.Host == "" {
		http.Error(w, "url must be an absolute http or https url", http.StatusBadRequest)
		return
	}

	for i := 0; i < 5; i++ {
		code, err := newCode(7)
		if err != nil {
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
		err = h.store.Save(code, req.URL)
		if errors.Is(err, ErrExists) {
			continue
		}
		if err != nil {
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(createResponse{Code: code})
		return
	}
	http.Error(w, "could not allocate a code", http.StatusInternalServerError)
}

func (h *Handler) redirect(w http.ResponseWriter, r *http.Request) {
	target, err := h.store.Get(r.PathValue("code"))
	if errors.Is(err, ErrNotFound) {
		http.NotFound(w, r)
		return
	}
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, target, http.StatusFound)
}

func newCode(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	for i := range b {
		b[i] = alphabet[int(b[i])%len(alphabet)]
	}
	return string(b), nil
}
