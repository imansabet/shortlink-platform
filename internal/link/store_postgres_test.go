package link

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"
	"time"
)

func TestPostgresStore(t *testing.T) {
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("TEST_DATABASE_URL not set")
	}
	s, err := NewPostgresStore(context.Background(), dsn)
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()

	code := fmt.Sprintf("t%d", time.Now().UnixNano())

	if err := s.Save(code, "https://example.com"); err != nil {
		t.Fatalf("save: %v", err)
	}
	if err := s.Save(code, "https://example.org"); !errors.Is(err, ErrExists) {
		t.Fatalf("want ErrExists, got %v", err)
	}
	got, err := s.Get(code)
	if err != nil || got != "https://example.com" {
		t.Fatalf("get: %q, %v", got, err)
	}
	if _, err := s.Get(code + "x"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("want ErrNotFound, got %v", err)
	}
}
