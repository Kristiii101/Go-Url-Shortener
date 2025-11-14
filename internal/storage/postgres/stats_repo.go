package postgres

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Kristiii101/GO_URL_Shortener_ATAD/internal/storage"
)

type StatsRepo struct {
	pool *pgxpool.Pool
}

func NewStatsRepo(pool *pgxpool.Pool) *StatsRepo {
	return &StatsRepo{pool: pool}
}

func (r *StatsRepo) Totals(ctx context.Context, linkID int64) (int64, *time.Time, error) {
	var total int64
	var last *time.Time
	err := r.pool.QueryRow(ctx, `
		SELECT COALESCE(COUNT(*),0) AS total, MAX(occurred_at) AS last
		FROM clicks WHERE link_id = $1`, linkID).Scan(&total, &last)
	if err != nil {
		return 0, nil, err
	}
	return total, last, nil
}

func (r *StatsRepo) Daily(ctx context.Context, linkID int64, fromUTC, toUTC time.Time) ([]storage.DayCount, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT date(timezone('UTC', occurred_at)) AS day, COUNT(*) AS clicks
		FROM clicks
		WHERE link_id = $1
		  AND occurred_at >= $2
		  AND occurred_at < $3
		GROUP BY day
		ORDER BY day`,
		linkID, fromUTC, toUTC)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []storage.DayCount
	for rows.Next() {
		var d time.Time
		var c int64
		if err := rows.Scan(&d, &c); err != nil {
			return nil, err
		}
		out = append(out, storage.DayCount{Day: d, Clicks: c})
	}
	return out, rows.Err()
}
