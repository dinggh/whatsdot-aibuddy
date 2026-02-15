package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"whatsdot-aibuddy/backend/internal/auth"
	"whatsdot-aibuddy/backend/internal/config"
	"whatsdot-aibuddy/backend/internal/httpapi"
	"whatsdot-aibuddy/backend/internal/store"
	"whatsdot-aibuddy/backend/internal/wechat"
)

func main() {
	cfg := config.Load()
	ctx := context.Background()

	db, err := store.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("connect db failed: %v", err)
	}
	defer db.Close()

	svc := &httpapi.Server{
		Store:  &store.Store{DB: db},
		JWT:    auth.New(cfg.JWTSecret, cfg.JWTExpireAfter),
		WeChat: wechat.New(cfg.WeChatAppID, cfg.WeChatSecret, cfg.WeChatAPIBase),
	}

	h := httpapi.Recoverer(svc.Routes())
	httpSrv := &http.Server{
		Addr:         cfg.ServerAddr,
		Handler:      h,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 15 * time.Second,
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
