package http

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Kristiii101/GO_URL_Shortener_ATAD/internal/config"
)

func NewServer(cfg config.Config, logger *log.Logger, handler http.Handler) *http.Server {
	addr := fmt.Sprintf(":%d", cfg.Port)
	return &http.Server{
		Addr:         addr,
		Handler:      handler,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}
}
