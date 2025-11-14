package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Kristiii101/GO_URL_Shortener_ATAD/internal/config"
	apphttp "github.com/Kristiii101/GO_URL_Shortener_ATAD/internal/http"
	"github.com/Kristiii101/GO_URL_Shortener_ATAD/internal/observability"
	"github.com/Kristiii101/GO_URL_Shortener_ATAD/internal/storage/postgres"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	logger := observability.NewLogger()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	pool, err := postgres.Open(ctx, cfg.DatabaseURL)
	if err != nil {
		logger.Fatalf("db connect: %v", err)
	}
	defer pool.Close()

	router := apphttp.NewRouter(apphttp.Deps{
		Config: cfg,
		Logger: logger,
		DB:     pool,
	})

	srv := apphttp.NewServer(cfg, logger, router)

	go func() {
		logger.Printf("http listening on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("http server: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	shutdownCtx, cancelShutdown := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelShutdown()
	logger.Println("shutting down...")
	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Printf("graceful shutdown error: %v", err)
	}
}
