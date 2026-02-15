package main

import (
	"context"
	"log"
	"os"
	"path/filepath"

	"github.com/jackc/pgx/v5/pgxpool"

	"whatsdot-aibuddy/backend/internal/config"
)

func main() {
	cfg := config.Load()
	sqlPath := "migrations/001_init.sql"
	if len(os.Args) > 1 {
		sqlPath = os.Args[1]
	}

	sql, err := os.ReadFile(sqlPath)
	if err != nil {
		log.Fatalf("read sql file: %v", err)
	}

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("connect db: %v", err)
	}
	defer pool.Close()

	if _, err := pool.Exec(ctx, string(sql)); err != nil {
		log.Fatalf("execute sql: %v", err)
	}

	log.Printf("migration ok: %s", filepath.Base(sqlPath))
}
