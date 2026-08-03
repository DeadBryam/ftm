package process

import (
	"testing"
	"time"

	"github.com/sthbryan/ftm/internal/config"
)

// A status consumer that never reads must not be able to wedge the Manager.
// callStatusUpdate runs with m.mu held, so a blocking send would deadlock every
// other operation -- including the TUI, which shares the lock.
func TestCallStatusUpdateDoesNotBlockOnFullChannel(t *testing.T) {
	m := NewManager()
	m.SetStatusChannel(make(chan config.TunnelStatus, 2))

	done := make(chan struct{})
	go func() {
		defer close(done)
		for i := 0; i < 100; i++ {
			m.callStatusUpdate(config.TunnelStatus{ID: "t1", State: config.TunnelStateOnline})
		}
	}()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("callStatusUpdate blocked on a full channel, want the update dropped")
	}
}

func TestCallStatusUpdateWithoutChannelIsNoop(t *testing.T) {
	m := NewManager()

	done := make(chan struct{})
	go func() {
		defer close(done)
		m.callStatusUpdate(config.TunnelStatus{ID: "t1", State: config.TunnelStateStopped})
	}()

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("callStatusUpdate blocked with no channel set")
	}
}

func TestCallStatusUpdateDeliversWhenDrained(t *testing.T) {
	m := NewManager()
	ch := make(chan config.TunnelStatus, 4)
	m.SetStatusChannel(ch)

	want := config.TunnelStatus{ID: "t1", Name: "Foundry", State: config.TunnelStateOnline}
	m.callStatusUpdate(want)

	select {
	case got := <-ch:
		if got.ID != want.ID || got.State != want.State {
			t.Fatalf("got %+v, want %+v", got, want)
		}
	case <-time.After(time.Second):
		t.Fatal("status update was not delivered")
	}
}
