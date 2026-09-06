package http

import (
	"encoding/json"
	"net/http"
	"github.com/rikut0904/starter/create/go/internal/usecase"
)

func NewHandler(health usecase.Health) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) { w.Header().Set("Content-Type", "application/json"); _ = json.NewEncoder(w).Encode(health.Check()) })
	return mux
}
