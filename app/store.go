package main

import (
	"sync"
	"time"
)

type Store struct {
	mu   sync.RWMutex
	data map[string]entry
}

type entry struct {
	value     string
	expiresAt time.Time
}

func (e entry) expired(now time.Time) bool {
	return !e.expiresAt.IsZero() && !now.Before(e.expiresAt)
}

func NewStore() *Store {
	return &Store{
		data: make(map[string]entry),
	}
}

func (s *Store) Set(key, value string, ttl time.Duration) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	var exp time.Time
	if ttl > 0 {
		exp = time.Now().Add(ttl)
	}
	s.data[key] = entry{value, exp}
	return nil
}

func (s *Store) Get(key string) (string, bool) {
	now := time.Now()

	s.mu.RLock()
	e, found := s.data[key]
	s.mu.RUnlock()

	if !found {
		return "", false
	}

	if !e.expired(now) {
		return e.value, found
	}

	s.mu.Lock()
	e, found = s.data[key]
	if found && e.expired(now) {
		delete(s.data, key)
	}
	s.mu.Unlock()

	return "", false
}
