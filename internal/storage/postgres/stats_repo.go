package postgres

import (
	"context"
	"time"

	"github.com/Kristiii101/GO_URL_Shortener_ATAD/internal/storage"
	"github.com/jackc/pgx/v5/pgxpool"
)

type StatsRepo struct {
	DB *pgxpool.Pool
}

func NewStatsRepo(db *pgxpool.Pool) *StatsRepo {
	return &StatsRepo{DB: db}
}

// Totals returns the total clicks and the timestamp of the last click
func (r *StatsRepo) Totals(ctx context.Context, linkID int64) (int64, *time.Time, error) {
	// Query the 'clicks' table directly
	query := `
		SELECT COUNT(*), MAX(created_at)
		FROM clicks
		WHERE link_id = $1
	`
	var total int64
	var lastClick *time.Time

	err := r.DB.QueryRow(ctx, query, linkID).Scan(&total, &lastClick)
	if err != nil {
		return 0, nil, err
	}
	return total, lastClick, nil
}

// Daily returns a list of clicks grouped by day (YYYY-MM-DD)
func (r *StatsRepo) Daily(ctx context.Context, linkID int64, fromUTC, toUTC time.Time) ([]storage.DayCount, error) {
	// We use date_trunc to group by day.
	// Note: We use 'created_at' because that's what we defined in the clicks table.
	query := `
		SELECT date_trunc('day', created_at) as day, COUNT(*) as count
		FROM clicks
		WHERE link_id = $1 AND created_at >= $2 AND created_at <= $3
		GROUP BY day
		ORDER BY day ASC
	`

	rows, err := r.DB.Query(ctx, query, linkID, fromUTC, toUTC)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []storage.DayCount
	for rows.Next() {
		var d storage.DayCount
		if err := rows.Scan(&d.Day, &d.Clicks); err != nil {
			return nil, err
		}
		results = append(results, d)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}
