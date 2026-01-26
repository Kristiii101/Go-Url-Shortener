package postgres

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ClicksRepo struct {
	DB *pgxpool.Pool
}

func NewClicksRepo(db *pgxpool.Pool) *ClicksRepo {
	return &ClicksRepo{DB: db}
}

func (r *ClicksRepo) Insert(ctx context.Context, linkID int64, occurredAt time.Time, visitorHash, countryCode, userAgent, referer *string) error {
	query := `
        INSERT INTO clicks (link_id, created_at, visitor_hash, country_code, user_agent, referer)
        VALUES ($1, $2, $3, $4, $5, $6)
    `
	_, err := r.DB.Exec(ctx, query, linkID, occurredAt, visitorHash, countryCode, userAgent, referer)
	return err
}
