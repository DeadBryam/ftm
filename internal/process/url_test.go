package process

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/sthbryan/ftm/internal/config"
)

func TestRepeatedURLInLogsNotifiesOnce(t *testing.T) {
	script := `echo "https://example.trycloudflare.com"
echo "https://example.trycloudflare.com"
echo "https://example.trycloudflare.com"
sleep 1`

	m, updates := newScriptManager(t, script)
	t.Cleanup(m.StopAll)

	var online int32
	m.SetNotificationHandler(func(status config.TunnelStatus) {
		if status.State == config.TunnelStateOnline {
			atomic.AddInt32(&online, 1)
		}
	})

	if err := m.Start(tunnelFixture(), nil); err != nil {
		t.Fatalf("Start: %v", err)
	}

	waitForState(t, updates, config.TunnelStateOnline)
	time.Sleep(500 * time.Millisecond)

	if got := atomic.LoadInt32(&online); got != 1 {
		t.Errorf("every request line notified: got %d online notifications, want 1", got)
	}
}

func TestURLChangeStillNotifies(t *testing.T) {
	script := `echo "https://first.trycloudflare.com"
sleep 0.3
echo "https://second.trycloudflare.com"
sleep 1`

	m, _ := newScriptManager(t, script)
	t.Cleanup(m.StopAll)

	seen := make(chan string, 8)
	m.SetNotificationHandler(func(status config.TunnelStatus) {
		if status.State == config.TunnelStateOnline {
			seen <- status.PublicURL
		}
	})

	if err := m.Start(tunnelFixture(), nil); err != nil {
		t.Fatalf("Start: %v", err)
	}

	want := []string{"https://first.trycloudflare.com", "https://second.trycloudflare.com"}
	for _, url := range want {
		select {
		case got := <-seen:
			if got != url {
				t.Fatalf("got notification for %q, want %q", got, url)
			}
		case <-time.After(5 * time.Second):
			t.Fatalf("never notified for %q", url)
		}
	}
}
