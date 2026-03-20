package server

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"net/http"
	"net/http/httptest"
	"testing"
	"url_shortener/internal/storage"
)

type mockService struct {
	saveURLFunc func(ctx context.Context, url string) (string, error)
	getURLFunc  func(ctx context.Context, alias string) (string, error)
}

func (m *mockService) SaveURL(ctx context.Context, url string) (string, error) {
	return m.saveURLFunc(ctx, url)
}

func (m *mockService) GetURL(ctx context.Context, alias string) (string, error) {
	return m.getURLFunc(ctx, alias)
}

func aliasHelp(r *http.Request, alias string) *http.Request {
	routCtx := chi.NewRouteContext()
	routCtx.URLParams.Add("alias", alias)
	return r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, routCtx))
}

func TestSaveURL(t *testing.T) {
	cases := []struct {
		name        string
		body        string
		saveURLFunc func(ctx context.Context, url string) (string, error)
		wantStatus  int
		wantAlias   string
	}{
		{
			name:       "invalid json",
			body:       "invalid bruh",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "empty url",
			body:       `{"url_to_save": ""}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "invalid url",
			body:       `{"url_to_save": "bruh no url here"}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "service error",
			body: `{"url_to_save": "https://youtube.com"}`,
			saveURLFunc: func(_ context.Context, _ string) (string, error) {
				return "", errors.New("some error unexpected")
			},
			wantStatus: http.StatusInternalServerError,
		},
		{
			name: "success",
			body: `{"url_to_save": "https://youtube.com"}`,
			saveURLFunc: func(_ context.Context, _ string) (string, error) {
				return "1234567890", nil
			},
			wantStatus: http.StatusCreated,
			wantAlias:  "1234567890",
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			srv := &mockService{saveURLFunc: tt.saveURLFunc}
			handl := NewHandler(srv)

			req := httptest.NewRequest(http.MethodPost, "/api/shorten", bytes.NewBufferString(tt.body))
			w := httptest.NewRecorder()
			handl.SaveURL(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("expected: %d, got %d", tt.wantStatus, w.Code)
			}

			if tt.wantAlias != "" {
				var resp Response
				if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
					t.Errorf("failed to decode resp: %v", err)
				}
				if resp.Alias != tt.wantAlias {
					t.Errorf("wanted alias: %s, got %s", tt.wantAlias, resp.Alias)
				}
			}

		})
	}
}

func TestGetURL(t *testing.T) {
	cases := []struct {
		name         string
		alias        string
		getURLFunc   func(ctx context.Context, alias string) (string, error)
		wantStatus   int
		wantURLState string
	}{
		{
			name:       "empty alias",
			alias:      "",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:  "not found",
			alias: "1234567890",
			getURLFunc: func(_ context.Context, _ string) (string, error) {
				return "", storage.ErrURLNotFound
			},
			wantStatus: http.StatusNotFound,
		},
		{
			name:  "service error",
			alias: "1234567890",
			getURLFunc: func(_ context.Context, _ string) (string, error) {
				return "", errors.New("unexpected error")
			},
			wantStatus: http.StatusInternalServerError,
		},
		{
			name:  "succcess",
			alias: "1234567890",
			getURLFunc: func(_ context.Context, _ string) (string, error) {
				return "https://youtube.com", nil
			},
			wantStatus:   http.StatusFound,
			wantURLState: "https://youtube.com",
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			srv := &mockService{getURLFunc: tt.getURLFunc}
			handl := NewHandler(srv)
			req := httptest.NewRequest(http.MethodGet, "/"+tt.alias, nil)
			req = aliasHelp(req, tt.alias)
			w := httptest.NewRecorder()
			handl.GetURL(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("wanted status: %d, got %d", tt.wantStatus, w.Code)
			}
			if tt.wantURLState != "" {
				if urlGot := w.Header().Get("Location"); urlGot != tt.wantURLState {
					t.Errorf("wanted location: %s, got :%s", tt.wantURLState, urlGot)
				}
			}
		})
	}
}
