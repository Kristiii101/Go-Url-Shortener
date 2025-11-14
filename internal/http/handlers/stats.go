package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/Kristiii101/GO_URL_Shortener_ATAD/internal/config"
	"github.com/Kristiii101/GO_URL_Shortener_ATAD/internal/storage"
	"github.com/Kristiii101/GO_URL_Shortener_ATAD/internal/util"
)

type StatsDeps struct {
	Config    config.Config
	Logger    *log.Logger
	LinksRepo storage.LinksRepo
	StatsRepo storage.StatsRepo
}

type statsResponse struct {
	Key           string        `json:"key"`
	ShortURL      string        `json:"short_url"`
	TotalClicks   int64         `json:"total_clicks"`
	LastClickedAt *time.Time    `json:"last_clicked_at,omitempty"`
	Daily         []dailyRecord `json:"daily"`
	From          string        `json:"from"`
	To            string        `json:"to"`
}
type dailyRecord struct {
	Day    string `json:"day"` // YYYY-MM-DD
	Clicks int64  `json:"clicks"`
}

// Handles GET /v1/links/{key}/stats[?from=YYYY-MM-DD&to=YYYY-MM-DD]
func Stats(d StatsDeps) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.Header().Set("Allow", http.MethodGet)
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}
		// Path parsing
		prefix := "/v1/links/"
		if !strings.HasPrefix(r.URL.Path, prefix) {
			http.NotFound(w, r)
			return
		}
		rest := strings.TrimPrefix(r.URL.Path, prefix) // "{key}/stats"
		parts := strings.Split(strings.Trim(rest, "/"), "/")
		if len(parts) != 2 || parts[1] != "stats" {
			http.NotFound(w, r)
			return
		}
		key := parts[0]

		link, err := d.LinksRepo.GetByKey(r.Context(), key)
		if err != nil {
			http.NotFound(w, r)
			return
		}

		// Date range
		q := r.URL.Query()
		var from, to time.Time
		const layout = "2006-01-02"
		today := time.Now().UTC().Truncate(24 * time.Hour)
		if fromStr := q.Get("from"); fromStr != "" {
			f, err := time.ParseInLocation(layout, fromStr, time.UTC)
			if err != nil {
				util.WriteError(w, http.StatusBadRequest, "bad_request", "invalid from date")
				return
			}
			from = f
		} else {
			from = today.AddDate(0, 0, -30)
		}
		if toStr := q.Get("to"); toStr != "" {
			t, err := time.ParseInLocation(layout, toStr, time.UTC)
			if err != nil {
				util.WriteError(w, http.StatusBadRequest, "bad_request", "invalid to date")
				return
			}
			// make 'to' exclusive by adding a day at midnight
			to = t.AddDate(0, 0, 1)
		} else {
			to = today.AddDate(0, 0, 1)
		}

		total, last, err := d.StatsRepo.Totals(r.Context(), link.ID)
		if err != nil {
			d.Logger.Printf("stats totals error: %v", err)
			util.WriteError(w, http.StatusInternalServerError, "server_error", "could not fetch stats")
			return
		}
		days, err := d.StatsRepo.Daily(r.Context(), link.ID, from, to)
		if err != nil {
			d.Logger.Printf("stats daily error: %v", err)
			util.WriteError(w, http.StatusInternalServerError, "server_error", "could not fetch stats")
			return
		}
		out := make([]dailyRecord, 0, len(days))
		for _, dc := range days {
			out = append(out, dailyRecord{
				Day:    dc.Day.Format(layout),
				Clicks: dc.Clicks,
			})
		}

		resp := statsResponse{
			Key:           link.Key,
			ShortURL:      d.Config.BaseURL + "/" + link.Key,
			TotalClicks:   total,
			LastClickedAt: last,
			Daily:         out,
			From:          from.Format(layout),
			To:            to.AddDate(0, 0, -1).Format(layout), // inclusive end date
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	})
}
