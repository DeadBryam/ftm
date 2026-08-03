package process

import (
	"strings"
	"sync"
)

type LogBuffer struct {
	mu        sync.RWMutex
	lines     []string
	maxLen    int
	OnNewLine func(string)
}

func NewLogBuffer() *LogBuffer {
	return &LogBuffer{
		lines:  make([]string, 0, 100),
		maxLen: 500,
	}
}

func (lb *LogBuffer) Write(p []byte) (n int, err error) {
	var added []string
	for _, line := range strings.Split(string(p), "\n") {
		if line = strings.TrimSpace(line); line != "" {
			added = append(added, line)
		}
	}

	lb.mu.Lock()
	lb.lines = append(lb.lines, added...)
	if len(lb.lines) > lb.maxLen {
		lb.lines = lb.lines[len(lb.lines)-lb.maxLen:]
	}
	notify := lb.OnNewLine
	lb.mu.Unlock()

	// Published outside the lock and in sequence. This used to be
	// `go lb.OnNewLine(line)` per line, which handed the lines to subscribers
	// in whatever order the goroutines happened to get scheduled, so live logs
	// in the web UI arrived shuffled.
	if notify != nil {
		for _, line := range added {
			notify(line)
		}
	}

	return len(p), nil
}

func (lb *LogBuffer) GetLines() []string {
	lb.mu.RLock()
	defer lb.mu.RUnlock()

	result := make([]string, len(lb.lines))
	copy(result, lb.lines)
	return result
}
