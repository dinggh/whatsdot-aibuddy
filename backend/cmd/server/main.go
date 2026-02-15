package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"whatsdot-aibuddy/backend/internal/config"
	"whatsdot-aibuddy/backend/internal/httpapi"
	"whatsdot-aibuddy/backend/internal/logger"
	"whatsdot-aibuddy/backend/internal/openai"
	"whatsdot-aibuddy/backend/internal/store"
)

func main() {
	cfg := config.Load()

	closer, err := logger.Init(cfg.LogDir)
	if err != nil {
		log.Fatalf("init logger: %v", err)
	}
	defer closer.Close()

	ctx := context.Background()
	db, err := store.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("connect db failed: %v", err)
	}
	defer db.Close()

	svc := &httpapi.Server{
		Store:       &store.Store{DB: db},
		OpenAI:      openai.New(cfg.OpenAIBaseURL, cfg.OpenAIAPIKey, cfg.OpenAIModel),
		UploadDir:   cfg.UploadDir,
		AnalyzeMock: cfg.AnalyzeMock,
		Limiter:     httpapi.NewDeviceLimiter(cfg.RateLimitCapacity, cfg.RateLimitRefill),
	}

	httpSrv := &http.Server{
		Addr:         cfg.ServerAddr,
		Handler:      svc.Engine(),
		ReadTimeout:  20 * time.Second,
		WriteTimeout: 60 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("server listening on %s", cfg.ServerAddr)
		if err := httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server failed: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
	log.Println("shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := httpSrv.Shutdown(shutdownCtx); err != nil {
		log.Printf("shutdown error: %v", err)
	}
}
