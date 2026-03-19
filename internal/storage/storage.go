package storage

import "errors"

var (
	ErrURLNotFound = errors.New("url not found :(")
	ErrURLExists   = errors.New("that url already exists")
	ErrCollision   = errors.New("collision appeared")
)
