package storage

import (
	"context"
	"time"

	"github.com/Kristiii101/GO_URL_Shortener_ATAD/internal/domain"
)

type LinksRepo interface {
	GetByKey(ctx context.Context, key string) (*domain.Link, error)
	GetSystemByCanonicalURL(ctx context.Context, canonicalURL string) (*domain.Link, error)
	CreateSystem(ctx context.Context, canonicalURL string, expiresAt *time.Time) (*domain.Link, error)
	CreateAlias(ctx context.Context, alias string, canonicalURL string, expiresAt *time.Time) (*domain.Link, error)
	Disable(ctx context.Context, key string) error
}

type ClicksRepo interface {
	Insert(ctx context.Context, linkID int64, occurredAt time.Time, visitorHash *string, countryCode *string, userAgent *string, referer *string) error
}

type StatsRepo interface {
	Totals(ctx context.Context, linkID int64) (clicksTotal int64, lastClickedAt *time.Time, err error)
	Daily(ctx context.Context, linkID int64, fromUTC, toUTC time.Time) ([]DayCount, error)
}

type DayCount struct {
	Day    time.Time
	Clicks int64
}
