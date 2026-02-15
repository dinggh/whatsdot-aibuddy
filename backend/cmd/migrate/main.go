package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"

	"github.com/jackc/pgx/v5/pgxpool"

	"whatsdot-aibuddy/backend/internal/config"
)

func main() {
	cfg := config.Load()

	paths, err := discoverSQLPaths()
	if err != nil {
		log.Fatalf("discover sql files: %v", err)
	}

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("connect db: %v", err)
	}
	defer pool.Close()

	for _, sqlPath := range paths {
		sql, err := os.ReadFile(sqlPath)
		if err != nil {
			log.Fatalf("read sql file %s: %v", sqlPath, err)
		}
		if _, err := pool.Exec(ctx, string(sql)); err != nil {
			log.Fatalf("execute sql %s: %v", sqlPath, err)
		}
		log.Printf("migration ok: %s", filepath.Base(sqlPath))
	}
}

func discoverSQLPaths() ([]string, error) {
	if len(os.Args) > 1 {
		return []string{os.Args[1]}, nil
	}
	entries, err := filepath.Glob("migrations/*.sql")
	if err != nil {
		return nil, err
	}
	if len(entries) == 0 {
		return nil, fmt.Errorf("no migration files found in migrations/*.sql")
	}
	sort.Strings(entries)
	return entries, nil
}
