package handlers

import (
	"context"
	"log"
	"net"
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
	// Prepare the file server for the frontend
	index := Home(webDir)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		// 1. Routing Logic: Serve Frontend or API
		// If root or explicit index.html, serve the React/HTML app
		if path == "/" || path == "/index.html" {
			index.ServeHTTP(w, r)
			return
		}
		// Block direct access to API paths or health checks via this handler
		if strings.HasPrefix(path, "/v1/") || strings.HasPrefix(path, "/healthz") {
			http.NotFound(w, r)
			return
		}

		// 2. Extract Short Code
		// Remove the leading slash (e.g., "/AbCd" -> "AbCd")
		key := strings.TrimPrefix(path, "/")
		if key == "" {
			http.NotFound(w, r)
			return
		}

		// 3. Lookup Link in Database
		// We use r.Context() so we don't waste resources if the user disconnects
		link, err := d.LinksRepo.GetByKey(r.Context(), key)
		if err != nil {
			// If link doesn't exist, return 404
			http.NotFound(w, r)
			return
		}

		// 4. Check Active/Expiration Status (USER STORY #5)
		// We check if the link is manually disabled OR if the expiration date has passed
		if link.IsDisabled {
			http.Error(w, "Link has been disabled", http.StatusGone) // 410 Gone
			return
		}
		if link.ExpiresAt != nil && link.ExpiresAt.Before(time.Now().UTC()) {
			http.Error(w, "Link has expired", http.StatusGone) // 410 Gone
			return
		}

		// 5. Async Analytics (Fire and Forget)
		// We capture request data BEFORE starting the goroutine to avoid race conditions
		userAgent := r.UserAgent()
		ip := getRealIP(r)
		referer := r.Referer()

		go func() {
			// CRITICAL: Create a new context.
			// r.Context() dies when the ServeHTTP function returns.
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			// Prepare the time
			occurredAt := time.Now().UTC()

			// Prepare country code (nil for now, since we don't have GeoIP yet)
			var countryCode *string = nil

			// Call Insert with exactly 7 arguments as defined in your repo
			// We pass addresses (&ip, &userAgent, etc) because the repo expects *string
			err := d.ClicksRepo.Insert(ctx, link.ID, occurredAt, &ip, countryCode, &userAgent, &referer)

			if err != nil {
				d.Logger.Printf("Analytics error (key=%s): %v", key, err)
			}
		}()

		// 6. Perform Redirect
		// 307 Temporary Redirect is preferred over 302 for API-like redirects
		// to preserve the HTTP method, though 302 is also acceptable.
		http.Redirect(w, r, link.LongURL, http.StatusTemporaryRedirect)
	})
}

// Helper: Extract the correct user IP (handles Proxies/Cloudflare)
func getRealIP(r *http.Request) string {
	// 1. Check X-Forwarded-For (Standard for proxies)
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		parts := strings.Split(xff, ",")
		return strings.TrimSpace(parts[0])
	}
	// 2. Check X-Real-IP
	if xrip := r.Header.Get("X-Real-IP"); xrip != "" {
		return xrip
	}
	// 3. Fallback to RemoteAddr (contains port, e.g. "127.0.0.1:8989")
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}
