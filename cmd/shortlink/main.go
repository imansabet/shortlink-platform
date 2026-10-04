package main

import (
	"errors"
	"log/slog"
	"net/http"
	"os"
	"time"
	"context"
	"github.com/imansabet/shortlink-platform/internal/link"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	addr := ":8080"
	if v := os.Getenv("LISTEN_ADDR"); v != "" {
		addr = v
	}


var store link.Store = link.NewMemoryStore()
if dsn := os.Getenv("DATABASE_URL"); dsn != "" {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	pg, err := link.NewPostgresStore(ctx, dsn)
	cancel()
	if err != nil {
		logger.Error("database connection failed", "error", err)
		os.Exit(1)
	}
	store = pg
	logger.Info("using postgres store")
} else {
	logger.Info("using in-memory store")
}

handler := link.NewHandler(store)
	srv := &http.Server{
		Addr:              addr,
		Handler:           handler.Routes(),
		ReadHeaderTimeout: 5 * time.Second,
	}

	logger.Info("starting server", "addr", addr)
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Error("server failed", "error", err)
		os.Exit(1)
	}
}
