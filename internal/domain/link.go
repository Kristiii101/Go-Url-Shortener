package domain

import "time"

type Link struct {
	ID         int64
	Key        string
	LongURL    string
	IsCustom   bool
	CreatedAt  time.Time
	ExpiresAt  *time.Time
	IsDisabled bool
}
