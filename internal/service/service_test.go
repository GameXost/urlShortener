package service

import (
	"context"
	"testing"
	"url_shortener/internal/storage"
)

type MockStorage struct {
	SaveURLFunc func(ctx context.Context, url, alias string) (string, error)
	GetURLFunc  func(ctx context.Context, alias string) (string, error)
}

func (m *MockStorage) SaveURL(ctx context.Context, url, alias string) (string, error) {
	return m.SaveURLFunc(ctx, url, alias)
}

func (m *MockStorage) GetURL(ctx context.Context, alias string) (string, error) {
	return m.GetURLFunc(ctx, alias)
}

func TestServiceSaveURL(t *testing.T) {
	cases := []struct {
		name      string
		inputURL  string
		mockFuncs func(ctx context.Context, url, alias string) (string, error)
		wantError bool
		wantTries int
	}{
		{
			name:     "success from 1 try",
			inputURL: "https://habr.com",
			mockFuncs: func(ctx context.Context, url, alias string) (string, error) {
				return alias, nil
			},
			wantError: false,
			wantTries: 1,
		},
		{
			name:     "success after 9 collisions",
			inputURL: "https://habr.com",
			mockFuncs: func() func(context.Context, string, string) (string, error) {
				tries := 0
				return func(ctx context.Context, url, alias string) (string, error) {
					tries++
					if tries < 9 {
						return "", storage.ErrCollision
					}
					return alias, nil
				}
			}(),
			wantError: false,
			wantTries: 9,
		},
		{
			name:     "failure after 10 tries",
			inputURL: "https://habr.com",
			mockFuncs: func(ctx context.Context, url, alias string) (string, error) {
				return "", storage.ErrCollision
			},
			wantError: true,
			wantTries: 10,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			curTries := 0
			mockRepo := &MockStorage{
				SaveURLFunc: func(ctx context.Context, url, alias string) (string, error) {
					curTries++
					return tt.mockFuncs(ctx, url, alias)
				},
			}
			srvs := New(mockRepo)
			resAlias, err := srvs.SaveURL(context.Background(), tt.inputURL)
			
			if tt.wantError {
				if err == nil {
					t.Errorf("wanted error, but got success")
				}
			} else if !tt.wantError {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if resAlias == "" {
					t.Errorf("wnted alias, got nothing")
				}
			}
			if curTries != tt.wantTries {
				t.Errorf("wanted %d tries, got %d ", tt.wantTries, curTries)
			}

		})
	}
}
