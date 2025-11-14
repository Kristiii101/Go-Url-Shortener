package postgres

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ClicksRepo struct {
	pool *pgxpool.Pool
}

func NewClicksRepo(pool *pgxpool.Pool) *ClicksRepo {
	return &ClicksRepo{pool: pool}
}

func (r *ClicksRepo) Insert(ctx context.Context, linkID int64, occurredAt time.Time, visitorHash *string, countryCode *string, userAgent *string) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO clicks (link_id, occurred_at, visitor_hash, country_code, user_agent)
		VALUES ($1, $2, $3, $4, $5)`,
		linkID, occurredAt, visitorHash, countryCode, userAgent,
	)
	return err
}
