package entity

import "errors"

var (
	ErrProgressNotFound           = errors.New("progress is not found in the cache")
	ErrProgressCachedTypeMismatch = errors.New("type mismatch")
)
