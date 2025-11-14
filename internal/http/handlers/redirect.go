package handlers

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/Kristiii101/GO_URL_Shortener_ATAD/internal/config"
	"github.com/Kristiii101/GO_URL_Shortener_ATAD/internal/storage"
)

type RedirectDeps struct {
	Config     config.Config
	Logger     *log.Logger
	LinksRepo  storage.LinksRepo
	ClicksRepo storage.ClicksRepo
}

// Root serves index.html on "/" and treats any other single-segment path as a key to redirect.
func Root(webDir string, d RedirectDeps) http.Handler {
	index := Home(webDir) // reuse existing Home handler
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if path == "/" || path == "/index.html" {
			index.ServeHTTP(w, r)
			return
		}
		// ignore API and health paths (ServeMux should already route them)
		if strings.HasPrefix(path, "/v1/") || strings.HasPrefix(path, "/healthz") {
			http.NotFound(w, r)
			return
		}

		key := strings.TrimPrefix(path, "/")
		if key == "" {
			http.NotFound(w, r)
			return
		}

		link, err := d.LinksRepo.GetByKey(r.Context(), key)
		if err != nil {
			http.NotFound(w, r)
			return
		}
		// active?
		now := time.Now().UTC()
		if link.IsDisabled || (link.ExpiresAt != nil && !link.ExpiresAt.After(now)) {
			http.Error(w, http.StatusText(http.StatusGone), http.StatusGone)
			return
		}

		// record click (synchronously for simplicity)
		ua := r.UserAgent()
		if err := d.ClicksRepo.Insert(r.Context(), link.ID, now, nil, nil, &ua); err != nil {
			d.Logger.Printf("click insert failed (key=%s): %v", key, err)
			// continue anyway
		}

		w.Header().Set("Location", link.LongURL)
		w.WriteHeader(http.StatusTemporaryRedirect) // 307
	})
}
