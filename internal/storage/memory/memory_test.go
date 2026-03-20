package memory

import (
	"context"
	"errors"
	"sync"
	"testing"
	"url_shortener/internal/storage"
)

func TestLRUCache_SetnGet(t *testing.T) {
	cache := NewCache(10)
	ctx := context.Background()
	url := "https://habr.com"
	alias := "habr_alias"

	savedAlias, err := cache.SaveURL(ctx, url, alias)
	if err != nil {
		t.Errorf("unexpected error while saving: %v", err)
	}

	if savedAlias != alias {
		t.Errorf("want %s, got %s", alias, savedAlias)
	}

	gotURL, err := cache.GetURL(ctx, alias)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if gotURL != url {
		t.Errorf("want: %s, got: %s", url, gotURL)
	}
}

func TestLRUCacheSameAliases(t *testing.T) {
	cache := NewCache(10)
	ctx := context.Background()
	url := "https://habr.com"
	savedAliase1, err := cache.SaveURL(ctx, url, "first")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	savedAliase2, err := cache.SaveURL(ctx, url, "second")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if savedAliase1 != savedAliase2 {
		t.Errorf("expected aliases to be the same, got1: %s, got2 %s", savedAliase1, savedAliase2)
	}
}

func TestLRUCacheCapacity(t *testing.T) {
	cache := NewCache(2)
	ctx := context.Background()
	_, err := cache.SaveURL(ctx, "habr.com", "hb11111111")
	if err != nil {
		t.Errorf("got unexpected error: %v", err)
	}
	_, err = cache.SaveURL(ctx, "youtube.com", "yt11111111")
	if err != nil {
		t.Errorf("got unexpected error: %v", err)
	}
	_, err = cache.SaveURL(ctx, "instagram.com", "instagramm")
	if err != nil {
		t.Errorf("got unexpected error: %v", err)
	}

	_, err = cache.GetURL(ctx, "hb11111111")
	if err == nil {
		t.Errorf("expected an error, got nothing") // habr уже не в кэше
	}

	if !errors.Is(err, storage.ErrURLNotFound) {
		t.Errorf("expected: %v, got %v", storage.ErrURLNotFound, err)
	}

	if _, err = cache.GetURL(ctx, "yt11111111"); err != nil {
		t.Errorf("youtube must be in cache")
	}
	if _, err = cache.GetURL(ctx, "instagramm"); err != nil {
		t.Errorf("instagram must be in cache")
	}
}

func TestLRUCacheConcurrency(t *testing.T) {
	wantSize := 10
	cache := NewCache(wantSize)
	ctx := context.Background()
	var wg sync.WaitGroup

	for i := range 100 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, _ = cache.SaveURL(ctx, "alias"+string(rune(i)), "url"+string(rune(i)))
		}()
	}
	wg.Wait()
	if cache.size != 10 {
		t.Errorf("want cache size: %d, got %d", wantSize, cache.size)
	}
}
