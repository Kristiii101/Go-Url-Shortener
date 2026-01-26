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
	index := Home(webDir)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		// 1. Routing Logic: Serve Frontend or API
		if path == "/" || path == "/index.html" {
			index.ServeHTTP(w, r)
			return
		}
		if strings.HasPrefix(path, "/v1/") || strings.HasPrefix(path, "/healthz") {
			http.NotFound(w, r)
			return
		}

		// 2. Extract Short Code
		key := strings.TrimPrefix(path, "/")
		if key == "" {
			http.NotFound(w, r)
			return
		}

		// 3. Lookup Link
		// We use r.Context() here because we WANT to cancel if the user disconnects
		link, err := d.LinksRepo.GetByKey(r.Context(), key)
		if err != nil {
			// Using 404 is better than 500 for "Link not found"
			http.NotFound(w, r)
			return
		}

		// 4. Check Active/Expiration Status
		now := time.Now().UTC()
		if link.IsDisabled || (link.ExpiresAt != nil && !link.ExpiresAt.After(now)) {
			http.Error(w, "Link Expired", http.StatusGone)
			return
		}

		// 5. Async Analytics (Fire and Forget)
		// We capture the data we need BEFORE starting the goroutine
		userAgent := r.UserAgent()
		ip := getRealIP(r)
		referer := r.Referer()

		go func() {
			// CRITICAL: Create a new context.
			// r.Context() dies when the ServeHTTP function returns.
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			// Pass the IP and Referer to your repo
			// (Make sure your ClicksRepo.Insert signature matches this!)
			err := d.ClicksRepo.Insert(ctx, link.ID, time.Now().UTC(), &ip, nil, &userAgent, &referer)
			if err != nil {
				d.Logger.Printf("Analytics error (key=%s): %v", key, err)
			}
		}()

		// 6. Perform Redirect
		// 307 preserves the method (POST/GET), 302 is standard "Found".
		// 307 is preferred for modern apps.
		w.Header().Set("Location", link.LongURL)
		w.WriteHeader(http.StatusTemporaryRedirect)
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
