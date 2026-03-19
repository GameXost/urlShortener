package service

import (
	"context"
	"errors"
	"fmt"
	"url_shortener/internal/generator"
	"url_shortener/internal/storage"
)

const tries = 10

type URLStorage interface {
	SaveURL(ctx context.Context, URLToSave string, alias string) (string, error)
	GetURL(ctx context.Context, alias string) (string, error)
}

type Service struct {
	repo URLStorage
}

func New(repo URLStorage) *Service {
	return &Service{repo: repo}
}

func (s *Service) SaveURL(ctx context.Context, url string) (string, error) {
	for range tries {
		alias := generator.GenerateShortAlias()
		alias, err := s.repo.SaveURL(ctx, url, alias)
		if err != nil {
			if errors.Is(err, storage.ErrCollision) {
				continue
			}
			return "", fmt.Errorf("failed to save url: %w", err)
		}
		return alias, nil
	}
	return "", fmt.Errorf("failed alias generation cnt tries: %d", tries)
}

func (s *Service) GetURL(ctx context.Context, alias string) (string, error) {
	url, err := s.repo.GetURL(ctx, alias)
	if err != nil {
		return "", fmt.Errorf("failed to get link: %w", err)
	}
	return url, nil
}
