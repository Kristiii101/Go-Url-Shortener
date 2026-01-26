package postgres

import (
	"context"
	"errors"
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
	// FIXED: 'key' -> 'short_code'
	row := r.pool.QueryRow(ctx, `
        SELECT id, short_code, original_url, is_custom, created_at, expires_at, is_disabled
        FROM links WHERE short_code = $1`, key)

	var l domain.Link
	// FIXED: Scan into correct fields
	if err := row.Scan(&l.ID, &l.Key, &l.LongURL, &l.IsCustom, &l.CreatedAt, &l.ExpiresAt, &l.IsDisabled); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return &l, nil
}

func (r *LinksRepo) GetSystemByCanonicalURL(ctx context.Context, canonicalURL string) (*domain.Link, error) {
	// FIXED: 'key' -> 'short_code'
	row := r.pool.QueryRow(ctx, `
        SELECT id, short_code, original_url, is_custom, created_at, expires_at, is_disabled
        FROM links WHERE original_url = $1 AND is_custom = FALSE`, canonicalURL)

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
	// FIXED: 'key' -> 'short_code'
	err := r.pool.QueryRow(ctx, `
        INSERT INTO links (short_code, original_url, is_custom, expires_at)
        VALUES ($1, $2, TRUE, $3)
        RETURNING id, short_code, original_url, is_custom, created_at, expires_at, is_disabled
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
	// 1. Check if it already exists
	if l, err := r.GetSystemByCanonicalURL(ctx, canonicalURL); err == nil {
		return l, nil
	}

	// 2. Generate a random short code (Simple & Robust)
	// We use your internal 'id' package to generate a random string
	// If id.Generate() is not available, we can use a simple helper.
	// Assuming id.Generate(length) exists based on your architecture.
	code := id.Generate(r.minLen)

	// 3. Insert into DB
	// FIXED: 'key' -> 'short_code'
	query := `
        INSERT INTO links (original_url, short_code, is_custom, expires_at)
        VALUES ($1, $2, FALSE, $3)
        RETURNING id, created_at, is_disabled
    `

	var idVal int64
	var createdAt time.Time
	var isDisabled bool

	err := r.pool.QueryRow(ctx, query, canonicalURL, code, expiresAt).Scan(&idVal, &createdAt, &isDisabled)
	if err != nil {
		// If collision (rare), just return error or retry (for MVP, error is fine)
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			// If we hit a collision, we can check if it was due to the URL already existing
			// (race condition) or the random code colliding.
			if l, selErr := r.GetSystemByCanonicalURL(ctx, canonicalURL); selErr == nil {
				return l, nil
			}
			return nil, domain.ErrAliasInUse
		}
		return nil, err
	}

	return &domain.Link{
		ID:         idVal,
		Key:        code,
		LongURL:    canonicalURL,
		IsCustom:   false,
		CreatedAt:  createdAt,
		ExpiresAt:  expiresAt,
		IsDisabled: isDisabled,
	}, nil
}

func (r *LinksRepo) Disable(ctx context.Context, key string) error {
	// FIXED: 'key' -> 'short_code'
	ct, err := r.pool.Exec(ctx, `UPDATE links SET is_disabled = TRUE WHERE short_code = $1`, key)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}
