package link

import (
	"errors"
	"sync"
)

var (
	ErrNotFound = errors.New("link not found")
	ErrExists   = errors.New("code already exists")
)

// Store persists short codes and the URLs they point to.
type Store interface {
	Save(code, url string) error
	Get(code string) (string, error)
}

// MemoryStore is an in-memory Store, used until a database is wired in.
type MemoryStore struct {
	mu    sync.RWMutex
	links map[string]string
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{links: make(map[string]string)}
}

func (s *MemoryStore) Save(code, url string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.links[code]; ok {
		return ErrExists
	}
	s.links[code] = url
	return nil
}

func (s *MemoryStore) Get(code string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	url, ok := s.links[code]
	if !ok {
		return "", ErrNotFound
	}
	return url, nil
}
