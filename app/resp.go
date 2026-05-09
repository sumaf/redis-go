package main

import (
	"strconv"
	"sync"
)

type RESP struct {
	Type  byte
	Data  []string
	Raw   []byte
	Count int
}

type Store struct {
	mu sync.RWMutex
	data map[string]string
}

func AppendPrefix(s []byte, p byte, n int64) []byte {
	s = append(s, p)
	s = strconv.AppendInt(s, n, 10)
	return append(s, '\r', '\n')
}

func AppendString(s []byte, simple string) []byte {
	s = append(s, '+')
	s = append(s, simple...)
	return append(s, '\r', '\n')
}

func AppendBulkString(s []byte, bulk string) []byte {
	s = AppendPrefix(s, '$', int64(len(bulk)))
	s = append(s, '\r', '\n')
	s = append(s, bulk...)
	return append(s, '\r', '\n')
}


func NewStore() *Store {
	return &Store {
		data: make(map[string]string),
	}
}

func (s *Store) Set(key, value string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data[key] = value
}

func (s *Store) Get(key string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	value, found := s.data[key]
	return value, found
}
