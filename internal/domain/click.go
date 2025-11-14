package domain

import (
	"time"
)

// Click represents a click/visit event
type Click struct {
	ID          int64
	LinkID      int64
	OccurredAt  time.Time
	VisitorHash *string
	CountryCode *string
	UserAgent   *string
}

// ClickStats represents aggregated click statistics
type ClickStats struct {
	Total         int64
	LastClickedAt *time.Time
	ByCountry     map[string]int64
	ByDate        map[string]int64 // date string -> count
}
