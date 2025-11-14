package app

import (
	"context"
	"time"

	"github.com/Kristiii101/GO_URL_Shortener_ATAD/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

// App holds application-wide dependencies (optional helper)
type App struct {
	Config config.Config
	DB     *pgxpool.Pool
}

// HealthCheck performs a health check
func (a *App) HealthCheck(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	return a.DB.Ping(ctx)
}
