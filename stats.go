package main

import (
	"encoding/json"
	"os"
	"sync"
	"time"
)

var startedAt = time.Now()

// stats are shared by every visitor and persisted to disk.
type stats struct {
	mu     sync.Mutex
	path   string
	Visits int `json:"visits"`
	High   int `json:"high"`
	HighBy int `json:"high_by"`
}

func loadStats(path string) *stats {
	s := &stats{path: path}
	if b, err := os.ReadFile(path); err == nil {
		_ = json.Unmarshal(b, s)
	}
	return s
}

func (s *stats) visit() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Visits++
	s.save()
	return s.Visits
}

func (s *stats) high() (int, int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.High, s.HighBy
}

// submit records a score and reports whether it is a new all-time high.
func (s *stats) submit(score, visitor int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if score <= s.High {
		return false
	}
	s.High, s.HighBy = score, visitor
	s.save()
	return true
}

func (s *stats) visits() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.Visits
}

// save must be called with mu held.
func (s *stats) save() {
	b, _ := json.Marshal(s)
	tmp := s.path + ".tmp"
	if os.WriteFile(tmp, b, 0o600) == nil {
		_ = os.Rename(tmp, s.path)
	}
}
