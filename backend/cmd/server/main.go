package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"trello-clone/internal/database"
	"trello-clone/internal/server"
)

func env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func main() {
	ctx := context.Background()

	dsn := env("DATABASE_URL", "postgres://trello:trello@127.0.0.1:5432/trello?sslmode=disable")
	port := env("PORT", "8080")

	db, err := database.Open(ctx, dsn)
	if err != nil {
		log.Fatalf("database: %v", err)
	}
	defer db.Close()

	if err := database.Migrate(ctx, db); err != nil {
		log.Fatalf("migrate: %v", err)
	}

	srv := server.New(server.Config{
		DB:              db,
		CookieDomain:    env("COOKIE_DOMAIN", "localhost"),
		AllowOrigin:     env("BASE_URL", "http://localhost:5173"),
		BaseURL:         "http://localhost:" + port,
		GoogleID:        os.Getenv("GOOGLE_CLIENT_ID"),
		GoogleSecret:    os.Getenv("GOOGLE_CLIENT_SECRET"),
		MicrosoftID:     os.Getenv("MICROSOFT_CLIENT_ID"),
		MicrosoftSecret: os.Getenv("MICROSOFT_CLIENT_SECRET"),
	})

	srv.Addr = ":" + port
	ln, err := net.Listen("tcp", srv.Addr)
	if err != nil {
		log.Fatalf("listen: %v", err)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Printf("server listening on %s", srv.Addr)
		if err := srv.Serve(ln); err != nil && err != http.ErrServerClosed {
			log.Fatalf("serve: %v", err)
		}
	}()

	<-quit
	log.Println("shutting down...")

	shutdownCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("shutdown: %v", err)
	}
	log.Println("server stopped")
}
