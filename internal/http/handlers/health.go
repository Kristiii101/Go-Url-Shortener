package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type healthResponse struct {
	Status string `json:"status"`
	DB     string `json:"db"`
}

func Healthz(db *pgxpool.Pool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()
		dbStatus := "ok"
		if err := db.Ping(ctx); err != nil {
			dbStatus = "error"
			w.WriteHeader(http.StatusServiceUnavailable)
		}
		resp := healthResponse{
			Status: "ok",
			DB:     dbStatus,
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	})
}
