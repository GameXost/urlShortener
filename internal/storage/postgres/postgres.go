package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"url_shortener/internal/storage"
)

type Storage struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *Storage {
	return &Storage{db: db}
}

func (s *Storage) SaveURL(ctx context.Context, urlToSave, alias string) (string, error) {
	query := `INSERT INTO url (alias, url) VALUES ($1, $2) ON CONFLICT (url) DO UPDATE SET alias = url.alias RETURNING alias`
	var res string
	err := s.db.QueryRow(ctx, query, alias, urlToSave).Scan(&res)
	if err != nil {
		var postgresErr *pgconn.PgError
		if errors.As(err, &postgresErr) && postgresErr.Code == "23505" {
			return "", storage.ErrCollision
		}
		return "", fmt.Errorf("failed to Save url: %w", err)
	}
	return res, nil
}

func (s *Storage) GetURL(ctx context.Context, alias string) (string, error) {
	query := `SELECT u.url FROM url AS u WHERE u.alias = $1`
	var res string
	err := s.db.QueryRow(ctx, query, alias).Scan(&res)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", storage.ErrURLNotFound
		}
		return "", fmt.Errorf("failed to Get url: %w", err)
	}
	return res, nil
}
