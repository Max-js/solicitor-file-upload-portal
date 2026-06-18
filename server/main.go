package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	databaseURL := getenv("DATABASE_URL", "postgres://maxjs:maxjs@localhost:5432/maxjs?sslmode=disable")
	uploadDir := getenv("UPLOAD_DIR", "./uploads")
	port := getenv("PORT", "3000")
	const maxUploadBytes = 10 << 20

	ctx := context.Background()

	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		log.Fatalf("connect db: %v", err)
	}
	defer pool.Close()

	if err := waitForDB(ctx, pool); err != nil {
		log.Fatalf("database not ready: %v", err)
	}

	store, err := NewStore(pool, uploadDir)
	if err != nil {
		log.Fatalf("init store: %v", err)
	}

	api := &API{store: store, maxUploadBytes: maxUploadBytes}

	srv := &http.Server{
		Addr:              ":" + port,
		Handler:           api.routes(),
		ReadHeaderTimeout: 10 * time.Second,
	}

	log.Printf("listening on :%s (uploads -> %s)", port, uploadDir)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("server: %v", err)
	}
}

func waitForDB(ctx context.Context, pool *pgxpool.Pool) error {
	var err error
	for i := 0; i < 10; i++ {
		if err = pool.Ping(ctx); err == nil {
			return nil
		}
		time.Sleep(500 * time.Millisecond)
	}
	return err
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}