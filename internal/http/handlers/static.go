package handlers

import (
	"net/http"
	"os"
	"path/filepath"
)

func Home(webDir string) http.Handler {
	indexPath := filepath.Join(webDir, "index.html")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Serve only the homepage for now
		if r.URL.Path != "/" && r.URL.Path != "/index.html" {
			http.NotFound(w, r)
			return
		}
		if _, err := os.Stat(indexPath); err != nil {
			http.Error(w, "index.html not found", http.StatusNotFound)
			return
		}
		http.ServeFile(w, r, indexPath)
	})
}

// Optional: static files (CSS/JS/images) if you add them later under web/static
func StaticDir(prefix, dir string) http.Handler {
	fs := http.FileServer(http.Dir(dir))
	return http.StripPrefix(prefix, fs)
}
