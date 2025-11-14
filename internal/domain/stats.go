package domain

import (
	"time"
)

// LinkStats represents statistics for a link
type LinkStats struct {
	Link          *Link
	TotalClicks   int64
	LastClickedAt *time.Time
	DailyClicks   []DayStats
}

// DayStats represents daily statistics
type DayStats struct {
	Date   time.Time
	Clicks int64
}

// StatsRange represents a time range for statistics
type StatsRange struct {
	From time.Time
	To   time.Time
}

// NewStatsRange creates a stats range for the last N days
func NewStatsRange(days int) StatsRange {
	now := time.Now().UTC()
	return StatsRange{
		From: now.AddDate(0, 0, -days).Truncate(24 * time.Hour),
		To:   now,
	}
}
