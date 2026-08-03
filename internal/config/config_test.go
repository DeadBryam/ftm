package config

import (
	"os"
	"path/filepath"
	"testing"
)

// isolateHome points the config directory at a temp dir so tests never touch
// the developer's real ~/.config/foundry-tunnel.
func isolateHome(t *testing.T) string {
	t.Helper()

	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)

	return home
}

func TestLoadCreatesDefaultConfigWhenMissing(t *testing.T) {
	home := isolateHome(t)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}

	if cfg.NotificationsStatus != NotificationPending {
		t.Errorf("NotificationsStatus = %q, want %q", cfg.NotificationsStatus, NotificationPending)
	}
	if len(cfg.Tunnels) != 0 {
		t.Errorf("Tunnels = %v, want empty", cfg.Tunnels)
	}

	path := filepath.Join(home, ".config", AppName, ConfigFile)
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("Load() did not persist a default config at %s: %v", path, err)
	}
}

func TestSaveLoadRoundTrip(t *testing.T) {
	isolateHome(t)

	want := DefaultConfig()
	want.WebPort = 40507
	want.NotificationSound = false
	want.NotificationsStatus = NotificationGranted
	want.Tunnels = []TunnelConfig{
		{ID: "a", Name: "Foundry VTT", Provider: ProviderCloudflared, LocalPort: 30000},
		{ID: "b", Name: "Partida de prueba", Provider: ProviderPinggy, LocalPort: 30001},
	}

	if err := want.Save(); err != nil {
		t.Fatalf("Save() failed: %v", err)
	}

	got, err := Load()
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}

	if got.WebPort != want.WebPort {
		t.Errorf("WebPort = %d, want %d", got.WebPort, want.WebPort)
	}
	if got.NotificationSound != want.NotificationSound {
		t.Errorf("NotificationSound = %v, want %v", got.NotificationSound, want.NotificationSound)
	}
	if got.NotificationsStatus != want.NotificationsStatus {
		t.Errorf("NotificationsStatus = %q, want %q", got.NotificationsStatus, want.NotificationsStatus)
	}
	if len(got.Tunnels) != len(want.Tunnels) {
		t.Fatalf("Tunnels length = %d, want %d", len(got.Tunnels), len(want.Tunnels))
	}
	for i, tunnel := range want.Tunnels {
		if got.Tunnels[i] != tunnel {
			t.Errorf("Tunnels[%d] = %+v, want %+v", i, got.Tunnels[i], tunnel)
		}
	}
}

func TestLoadRejectsMalformedYAML(t *testing.T) {
	home := isolateHome(t)

	dir := filepath.Join(home, ".config", AppName)
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, ConfigFile), []byte("tunnels: [oops"), 0644); err != nil {
		t.Fatal(err)
	}

	if _, err := Load(); err == nil {
		t.Fatal("Load() accepted malformed YAML, want an error")
	}
}

func TestNormalizeNotificationsStatus(t *testing.T) {
	tests := map[string]string{
		NotificationGranted:  NotificationGranted,
		NotificationRejected: NotificationRejected,
		NotificationPending:  NotificationPending,
		"":                   NotificationPending,
		"nonsense":           NotificationPending,
	}

	for input, want := range tests {
		cfg := &Config{NotificationsStatus: input}
		cfg.NormalizeNotificationsStatus()

		if cfg.NotificationsStatus != want {
			t.Errorf("normalizing %q gave %q, want %q", input, cfg.NotificationsStatus, want)
		}
	}
}

func TestTunnelCRUD(t *testing.T) {
	cfg := DefaultConfig()

	cfg.AddTunnel(TunnelConfig{ID: "a", Name: "First", LocalPort: 30000})
	cfg.AddTunnel(TunnelConfig{ID: "b", Name: "Second", LocalPort: 30001})

	if got := cfg.GetTunnel("b"); got == nil || got.Name != "Second" {
		t.Fatalf("GetTunnel(b) = %+v, want the second tunnel", got)
	}
	if got := cfg.GetTunnel("missing"); got != nil {
		t.Errorf("GetTunnel(missing) = %+v, want nil", got)
	}

	if !cfg.RemoveTunnel("a") {
		t.Error("RemoveTunnel(a) = false, want true")
	}
	if cfg.RemoveTunnel("a") {
		t.Error("RemoveTunnel(a) twice = true, want false")
	}
	if len(cfg.Tunnels) != 1 || cfg.Tunnels[0].ID != "b" {
		t.Fatalf("Tunnels = %+v, want only b", cfg.Tunnels)
	}
}

// GetTunnel hands out a pointer into the slice, so callers can mutate config
// in place. Removing an entry must not leave that aliasing behind.
func TestGetTunnelReturnsLivePointer(t *testing.T) {
	cfg := DefaultConfig()
	cfg.AddTunnel(TunnelConfig{ID: "a", Name: "First", LocalPort: 30000})

	cfg.GetTunnel("a").LocalPort = 30500

	if got := cfg.Tunnels[0].LocalPort; got != 30500 {
		t.Fatalf("LocalPort = %d, want the mutation to stick", got)
	}
}
