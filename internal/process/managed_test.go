package process

import (
	"context"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/sthbryan/ftm/internal/config"
	"github.com/sthbryan/ftm/internal/providers"
)

// fakeProvider reports a URL for any line containing one, mimicking the real
// providers closely enough to exercise the writer's line splitting.
type fakeProvider struct{}

func (fakeProvider) Name() string       { return "Fake" }
func (fakeProvider) BinaryName() string { return "fake" }

func (fakeProvider) Start(context.Context, config.TunnelConfig, io.Writer) (*providers.Process, error) {
	return &providers.Process{}, nil
}

func (fakeProvider) ParseURL(line string) string {
	idx := strings.Index(line, "https://")
	if idx == -1 {
		return ""
	}
	return strings.Fields(line[idx:])[0]
}

func collectURLs(t *testing.T, writes []string) []string {
	t.Helper()

	var got []string
	w := newURLCapture(fakeProvider{}, func(url string) { got = append(got, url) })

	for _, chunk := range writes {
		n, err := w.Write([]byte(chunk))
		if err != nil {
			t.Fatalf("Write(%q) failed: %v", chunk, err)
		}
		if n != len(chunk) {
			t.Fatalf("Write(%q) = %d bytes, want %d", chunk, n, len(chunk))
		}
	}

	return got
}

func TestURLCaptureFindsURLInCompleteLine(t *testing.T) {
	got := collectURLs(t, []string{"tunnel ready at https://abc.trycloudflare.com\n"})

	if len(got) != 1 || got[0] != "https://abc.trycloudflare.com" {
		t.Fatalf("captured %v, want one trycloudflare URL", got)
	}
}

// Provider output arrives in arbitrary chunks, so a URL can be split across two
// writes. The tail has to be buffered until its newline shows up.
func TestURLCaptureJoinsSplitWrites(t *testing.T) {
	got := collectURLs(t, []string{"tunnel ready at https://abc.tryclo", "udflare.com\n"})

	if len(got) != 1 || got[0] != "https://abc.trycloudflare.com" {
		t.Fatalf("captured %v, want the URL reassembled across writes", got)
	}
}

func TestURLCaptureDoesNotEmitIncompleteLine(t *testing.T) {
	got := collectURLs(t, []string{"tunnel ready at https://abc.trycloudflare.com"})

	if len(got) != 0 {
		t.Fatalf("captured %v from a line with no newline yet, want nothing", got)
	}
}

func TestURLCaptureHandlesMultipleLinesInOneWrite(t *testing.T) {
	got := collectURLs(t, []string{
		"starting\nurl https://one.example.com\nnoise\nurl https://two.example.com\n",
	})

	want := []string{"https://one.example.com", "https://two.example.com"}
	if len(got) != len(want) {
		t.Fatalf("captured %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("URL %d = %q, want %q", i, got[i], want[i])
		}
	}
}

func TestURLCaptureIgnoresLinesWithoutURL(t *testing.T) {
	got := collectURLs(t, []string{"INF connecting\nINF retrying\n"})

	if len(got) != 0 {
		t.Fatalf("captured %v, want nothing", got)
	}
}

func TestLogBufferTrimsToMaxLength(t *testing.T) {
	lb := NewLogBuffer()

	for i := 0; i < lb.maxLen+250; i++ {
		if _, err := lb.Write([]byte("line\n")); err != nil {
			t.Fatalf("Write failed: %v", err)
		}
	}

	if got := len(lb.GetLines()); got > lb.maxLen {
		t.Fatalf("buffer holds %d lines, want at most %d", got, lb.maxLen)
	}
}

func TestLogBufferSkipsBlankLines(t *testing.T) {
	lb := NewLogBuffer()

	if _, err := lb.Write([]byte("first\n\n   \nsecond\n")); err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	got := lb.GetLines()
	want := []string{"first", "second"}

	if len(got) != len(want) {
		t.Fatalf("lines = %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("line %d = %q, want %q", i, got[i], want[i])
		}
	}
}

// Live log subscribers must see lines in the order the process wrote them.
// Publishing each line from its own goroutine shuffled them.
func TestLogBufferPublishesLinesInOrder(t *testing.T) {
	lb := NewLogBuffer()

	var got []string
	lb.OnNewLine = func(line string) { got = append(got, line) }

	if _, err := lb.Write([]byte("one\ntwo\nthree\n")); err != nil {
		t.Fatalf("Write failed: %v", err)
	}
	if _, err := lb.Write([]byte("four\nfive\n")); err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	want := []string{"one", "two", "three", "four", "five"}
	if len(got) != len(want) {
		t.Fatalf("published %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("published %v, want %v", got, want)
		}
	}
}

// A subscriber that blocks must not deadlock the writer against lb.mu.
func TestLogBufferWriteDoesNotHoldLockDuringPublish(t *testing.T) {
	lb := NewLogBuffer()
	lb.OnNewLine = func(string) {
		// Reading the buffer from inside the callback deadlocks if Write still
		// holds the lock while publishing.
		lb.GetLines()
	}

	done := make(chan struct{})
	go func() {
		defer close(done)
		lb.Write([]byte("line\n")) //nolint:errcheck // the write cannot fail
	}()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("Write deadlocked while publishing to a subscriber")
	}
}

// GetLines must hand back a copy: the TUI renders the result while the process
// keeps writing, and a shared slice would race.
func TestLogBufferGetLinesReturnsCopy(t *testing.T) {
	lb := NewLogBuffer()
	if _, err := lb.Write([]byte("original\n")); err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	lines := lb.GetLines()
	lines[0] = "mutated"

	if got := lb.GetLines()[0]; got != "original" {
		t.Fatalf("buffer line = %q, want the caller's mutation not to leak in", got)
	}
}
