package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Kristiii101/GO_URL_Shortener_ATAD/internal/domain"
	"github.com/Kristiii101/GO_URL_Shortener_ATAD/internal/id"
)

type LinksRepo struct {
	pool   *pgxpool.Pool
	minLen int
	maxLen int
}

func NewLinksRepo(pool *pgxpool.Pool, minLen, maxLen int) *LinksRepo {
	return &LinksRepo{pool: pool, minLen: minLen, maxLen: maxLen}
}

func (r *LinksRepo) GetByKey(ctx context.Context, key string) (*domain.Link, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT id, key, long_url, is_custom, created_at, expires_at, is_disabled
		FROM links WHERE key = $1`, key)
	var l domain.Link
	if err := row.Scan(&l.ID, &l.Key, &l.LongURL, &l.IsCustom, &l.CreatedAt, &l.ExpiresAt, &l.IsDisabled); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return &l, nil
}

func (r *LinksRepo) GetSystemByCanonicalURL(ctx context.Context, canonicalURL string) (*domain.Link, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT id, key, long_url, is_custom, created_at, expires_at, is_disabled
		FROM links WHERE long_url = $1 AND is_custom = FALSE`, canonicalURL)
	var l domain.Link
	if err := row.Scan(&l.ID, &l.Key, &l.LongURL, &l.IsCustom, &l.CreatedAt, &l.ExpiresAt, &l.IsDisabled); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return &l, nil
}

func (r *LinksRepo) CreateAlias(ctx context.Context, alias string, canonicalURL string, expiresAt *time.Time) (*domain.Link, error) {
	var l domain.Link
	err := r.pool.QueryRow(ctx, `
		INSERT INTO links (key, long_url, is_custom, expires_at)
		VALUES ($1, $2, TRUE, $3)
		RETURNING id, key, long_url, is_custom, created_at, expires_at, is_disabled
	`, alias, canonicalURL, expiresAt).Scan(&l.ID, &l.Key, &l.LongURL, &l.IsCustom, &l.CreatedAt, &l.ExpiresAt, &l.IsDisabled)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			return nil, domain.ErrAliasInUse
		}
		return nil, err
	}
	return &l, nil
}

func (r *LinksRepo) CreateSystem(ctx context.Context, canonicalURL string, expiresAt *time.Time) (*domain.Link, error) {
	// fast path: try to find existing
	if l, err := r.GetSystemByCanonicalURL(ctx, canonicalURL); err == nil {
		return l, nil
	}

	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	var idNum int64
	var createdAt time.Time
	if err := tx.QueryRow(ctx,
		`INSERT INTO links (long_url, is_custom, expires_at) VALUES ($1, FALSE, $2)
		 RETURNING id, created_at`, canonicalURL, expiresAt).Scan(&idNum, &createdAt); err != nil {
		// Unique on (long_url where is_custom=false) might fire under race; re-select
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			if l, selErr := r.GetSystemByCanonicalURL(ctx, canonicalURL); selErr == nil {
				return l, nil
			}
		}
		return nil, err
	}

	primary := id.Pad(id.Encode(int64(idNum)), r.minLen)
	candidate := primary

	// try primary, then primary + one-char suffixes if key is taken (e.g., by a custom alias)
	const suffixAlphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	tryUpdate := func(k string) error {
		_, err := tx.Exec(ctx, `UPDATE links SET key = $1 WHERE id = $2`, k, idNum)
		if err != nil {
			if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
				return domain.ErrAliasInUse // reuse error symbol for "key taken"
			}
			return err
		}
		return nil
	}

	err = tryUpdate(candidate)
	if errors.Is(err, domain.ErrAliasInUse) {
		// suffix attempts (one char)
	updated:
		for i := 0; i < len(suffixAlphabet); i++ {
			if r.maxLen > 0 && len(primary)+1 > r.maxLen {
				break // cannot append without exceeding maxLen
			}
			candidate = primary + string(suffixAlphabet[i])
			if err2 := tryUpdate(candidate); err2 == nil {
				err = nil
				break updated
			} else if !errors.Is(err2, domain.ErrAliasInUse) {
				err = err2
				break updated
			}
		}
	}
	if err != nil {
		return nil, fmt.Errorf("set key collision: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	// fetch full row
	return r.GetByKey(ctx, candidate)
}

func (r *LinksRepo) Disable(ctx context.Context, key string) error {
	ct, err := r.pool.Exec(ctx, `UPDATE links SET is_disabled = TRUE WHERE key = $1`, key)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}
