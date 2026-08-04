package proxy

import (
	"sync"
	"time"
)

type counter struct {
	mu       sync.Mutex
	window   time.Duration
	lastSeen map[string]time.Time
	sessions int
	requests int64
	now      func() time.Time
}

func newCounter(window time.Duration) *counter {
	if window <= 0 {
		window = 5 * time.Minute
	}

	return &counter{
		window:   window,
		lastSeen: make(map[string]time.Time),
		now:      time.Now,
	}
}

func (c *counter) seen(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.requests++
	c.lastSeen[key] = c.now()
	c.pruneLocked()
}

func (c *counter) openSession() {
	c.mu.Lock()
	c.sessions++
	c.mu.Unlock()
}

func (c *counter) closeSession() {
	c.mu.Lock()
	if c.sessions > 0 {
		c.sessions--
	}
	c.mu.Unlock()
}

func (c *counter) stats() Stats {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.pruneLocked()

	return Stats{
		ActiveSessions: c.sessions,
		Visitors:       len(c.lastSeen),
		Requests:       c.requests,
	}
}

func (c *counter) pruneLocked() {
	cutoff := c.now().Add(-c.window)

	for key, at := range c.lastSeen {
		if at.Before(cutoff) {
			delete(c.lastSeen, key)
		}
	}
}
