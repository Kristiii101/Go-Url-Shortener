package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"

	"github.com/Kristiii101/GO_URL_Shortener_ATAD/internal/config"
	apphttp "github.com/Kristiii101/GO_URL_Shortener_ATAD/internal/http"
	"github.com/Kristiii101/GO_URL_Shortener_ATAD/internal/observability"
	"github.com/Kristiii101/GO_URL_Shortener_ATAD/internal/storage/postgres"
)

func main() {
	// 1. Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on system environment variables")
	}

	// 2. Load Config
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config error: %v", err)
	}

	// 3. Setup Logger
	logger := observability.NewLogger()

	// 4. Connect to Database (Using your robust internal package)
	// We increase the timeout to 30s in case Supabase is "waking up"
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	logger.Println("Connecting to Database...")
	pool, err := postgres.Open(ctx, cfg.DatabaseURL)
	if err != nil {
		logger.Fatalf("db connect error: %v", err)
	}
	defer pool.Close()
	logger.Println("Database connection established successfully.")

	// 5. Setup Router & Server
	router := apphttp.NewRouter(apphttp.Deps{
		Config: cfg,
		Logger: logger,
		DB:     pool,
	})

	srv := apphttp.NewServer(cfg, logger, router)

	// 6. Start Server in Background
	go func() {
		logger.Printf("http listening on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("http server error: %v", err)
		}
	}()

	// 7. Graceful Shutdown
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
