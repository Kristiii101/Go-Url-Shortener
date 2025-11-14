package http

import (
	"log"
	stdhttp "net/http"
	"path/filepath"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Kristiii101/GO_URL_Shortener_ATAD/internal/config"
	"github.com/Kristiii101/GO_URL_Shortener_ATAD/internal/http/handlers"
	"github.com/Kristiii101/GO_URL_Shortener_ATAD/internal/http/middleware"
	"github.com/Kristiii101/GO_URL_Shortener_ATAD/internal/storage/postgres"
)

type Deps struct {
	Config config.Config
	Logger *log.Logger
	DB     *pgxpool.Pool
}

type Middleware func(stdhttp.Handler) stdhttp.Handler

func chain(h stdhttp.Handler, m ...Middleware) stdhttp.Handler {
	for i := len(m) - 1; i >= 0; i-- {
		h = m[i](h)
	}
	return h
}

func NewRouter(d Deps) stdhttp.Handler {
	mux := stdhttp.NewServeMux()

	global := []Middleware{
		middleware.Recover(d.Logger),
		middleware.RequestID(),
		middleware.Logging(d.Logger),
	}

	// Repos
	linksRepo := postgres.NewLinksRepo(d.DB, d.Config.KeyMinLen, d.Config.KeyMaxLen)
	clicksRepo := postgres.NewClicksRepo(d.DB)
	statsRepo := postgres.NewStatsRepo(d.DB)

	// Health
	mux.Handle("/healthz", chain(handlers.Healthz(d.DB), global...))

	// API
	createDeps := handlers.LinkDeps{Config: d.Config, Logger: d.Logger, LinksRepo: linksRepo}
	mux.Handle("/v1/links", chain(
		handlers.CreateLink(createDeps),
		append(global, middleware.RateLimitPerIP(d.Config.RateLimitCreate, d.Config.RateLimitWindow))...,
	))

	statsDeps := handlers.StatsDeps{Config: d.Config, Logger: d.Logger, LinksRepo: linksRepo, StatsRepo: statsRepo}
	mux.Handle("/v1/links/", chain(handlers.Stats(statsDeps), global...)) // handles /v1/links/{key}/stats

	// Static assets (optional)
	// mux.Handle("/static/", chain(handlers.StaticDir("/static/", filepath.Join(d.Config.WebDir, "static")), global...))

	// Root: serve UI at "/" and redirect for "/{key}"
	redirDeps := handlers.RedirectDeps{Config: d.Config, Logger: d.Logger, LinksRepo: linksRepo, ClicksRepo: clicksRepo}
	mux.Handle("/", chain(handlers.Root(d.Config.WebDir, redirDeps), global...))

	_ = filepath.Separator // avoid unused import if StaticDir is commented

	return mux
}
