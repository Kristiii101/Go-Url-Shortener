package domain

import "errors"

var (
	ErrInvalidURL   = errors.New("invalid_url")
	ErrInvalidAlias = errors.New("invalid_alias")
	ErrReservedKey  = errors.New("reserved_key")
	ErrAliasInUse   = errors.New("alias_in_use")
	ErrNotFound     = errors.New("not_found")
	ErrExpired      = errors.New("expired")
	ErrDisabled     = errors.New("disabled")
)
