package main

import (
	"context"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"url_shortener/internal/config"
	"url_shortener/internal/server"
	"url_shortener/internal/service"
	"url_shortener/internal/storage/memory"
	"url_shortener/internal/storage/postgres"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	var repo service.URLStorage
	var pool *pgxpool.Pool

	switch cfg.StorageType {
	case "postgres":
		if cfg.DataBaseUrl == "" {
			log.Fatal("DATABASE_URL is not set")
		}

		pgConfig, err := pgxpool.ParseConfig(cfg.DataBaseUrl)
		if err != nil {
			log.Fatalf("failed to parse conn: %v", err)
		}

		pool, err = pgxpool.NewWithConfig(context.Background(), pgConfig)
		if err != nil {
			log.Fatalf("failed to connect to DB: %v", err)
		}
		if err = pool.Ping(context.Background()); err != nil {
			log.Fatalf("failed to ping DB: %v", err)
		}

		repo = postgres.New(pool)
		log.Println("DB started successfully")
	case "memory":
		repo = memory.NewCache(cfg.MemoryCapacity)
		log.Println("in memory storage successfully started")
	default:
		log.Fatal("choose store type")
	}

	srv := service.New(repo)
	handler := server.NewHandler(srv)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "web/index.html")
	})
	r.Post("/api/shorten", handler.SaveURL)
	r.Get("/{alias}", handler.GetURL)

	httpServer := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}

	go func() {
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("failed to listen to 8080: %v", err)
		}
	}()

	sign := make(chan os.Signal, 1)
	signal.Notify(sign, syscall.SIGINT, syscall.SIGTERM)
	<-sign

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := httpServer.Shutdown(ctx); err != nil {
		log.Printf("server forced to shutdownЖ %v", err)
	}
	if pool != nil {
		pool.Close()
	}
}
