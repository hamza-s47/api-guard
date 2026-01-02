package store

import (
	"sync"

	"github.com/hamza-s47/api-guard/internal/limiter"
)

type MemoryStore struct {
	buckets map[string]*limiter.TokenBucket
	mu      sync.Mutex
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		buckets: make(map[string]*limiter.TokenBucket),
	}
}

func (s *MemoryStore) GetBucket(ip string) *limiter.TokenBucket {
	s.mu.Lock()
	defer s.mu.Unlock()

	if bucket, exist := s.buckets[ip]; exist {
		return bucket
	}

	bucket := limiter.NewTokenBucket(5, 1)
	s.buckets[ip] = bucket
	return bucket
}
