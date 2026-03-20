package storage

import "errors"

var (
	ErrURLNotFound = errors.New("url not found :(")
	ErrCollision   = errors.New("collision appeared")
)
