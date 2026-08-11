package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSlugify(t *testing.T) {
	tests := map[string]string{
		"Curse of Strahd":  "curse-of-strahd",
		"Curse of Strahd!": "curse-of-strahd",
		"  spaced  out  ":  "spaced-out",
		"Sesión uno":       "sesión-uno",
		"a/b?c#d":          "a-b-c-d",
		"Port 30000":       "port-30000",
		"---":              "",
		"":                 "",
	}

	for input, want := range tests {
		if got := slugify(input); got != want {
			t.Errorf("slugify(%q) = %q, want %q", input, got, want)
		}
	}
}

func TestNewTunnelIDIsUniquePerName(t *testing.T) {
	cfg := DefaultConfig()

	first := cfg.NewTunnelID("Curse of Strahd")
	cfg.AddTunnel(TunnelConfig{ID: first, Name: "Curse of Strahd"})

	second := cfg.NewTunnelID("Curse of Strahd")
	cfg.AddTunnel(TunnelConfig{ID: second, Name: "Curse of Strahd"})

	third := cfg.NewTunnelID("Curse of Strahd")

	if first != "curse-of-strahd" {
		t.Errorf("first id = %q, want the slug of the name", first)
	}
	if second == first || third == first || third == second {
		t.Fatalf("ids collide: %q, %q, %q", first, second, third)
	}
}

func TestNewTunnelIDFallsBackWhenTheNameHasNoSlug(t *testing.T) {
	cfg := DefaultConfig()

	first := cfg.NewTunnelID("...")
	cfg.AddTunnel(TunnelConfig{ID: first, Name: "..."})
	second := cfg.NewTunnelID("!!!")

	if first != fallbackTunnelID {
		t.Errorf("id = %q, want %q for a name with nothing to slug", first, fallbackTunnelID)
	}
	if second == first {
		t.Errorf("second unslugabble name reused %q", first)
	}
}

func TestBackToBackCreationsGetDistinctIDs(t *testing.T) {
	cfg := DefaultConfig()

	seen := make(map[string]bool)
	for i := 0; i < 50; i++ {
		id := cfg.NewTunnelID("Mesa")
		if seen[id] {
			t.Fatalf("id %q handed out twice after %d creations", id, i)
		}
		seen[id] = true
		cfg.AddTunnel(TunnelConfig{ID: id, Name: "Mesa"})
	}
}

func TestLoadRepairsDuplicateIDs(t *testing.T) {
	home := isolateHome(t)

	dir := filepath.Join(home, ".config", AppName)
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatal(err)
	}

	broken := "version: 4\nlanguage: auto\nonboarded: true\ntunnels:\n" +
		"    - id: tunnel-1786488003\n      name: Curse of Strahd\n      provider: cloudflared\n      local_port: 30000\n" +
		"    - id: tunnel-1786488003\n      name: Mesa dos\n      provider: pinggy\n      local_port: 30001\n"
	if err := os.WriteFile(filepath.Join(dir, ConfigFile), []byte(broken), 0644); err != nil {
		t.Fatal(err)
	}

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}

	if len(cfg.Tunnels) != 2 {
		t.Fatalf("Tunnels length = %d, want both kept", len(cfg.Tunnels))
	}
	if cfg.Tunnels[0].ID == cfg.Tunnels[1].ID {
		t.Fatalf("ids still collide: %q", cfg.Tunnels[0].ID)
	}
	if cfg.Tunnels[0].ID != "tunnel-1786488003" {
		t.Errorf("first id = %q, want the original left alone", cfg.Tunnels[0].ID)
	}
	if cfg.Tunnels[1].ID != "mesa-dos" {
		t.Errorf("second id = %q, want it renamed after its own name", cfg.Tunnels[1].ID)
	}

	reloaded, err := Load()
	if err != nil {
		t.Fatalf("second Load() failed: %v", err)
	}
	if reloaded.Tunnels[1].ID != "mesa-dos" {
		t.Errorf("id = %q on reload, want the repair persisted", reloaded.Tunnels[1].ID)
	}
}

func TestLoadLeavesDistinctIDsAlone(t *testing.T) {
	home := isolateHome(t)

	dir := filepath.Join(home, ".config", AppName)
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatal(err)
	}

	fine := "version: 4\nlanguage: auto\nonboarded: true\ntunnels:\n" +
		"    - id: foundry-default\n      name: Foundry VTT\n      provider: cloudflared\n      local_port: 30000\n" +
		"    - id: mesa\n      name: Mesa\n      provider: pinggy\n      local_port: 30001\n"
	if err := os.WriteFile(filepath.Join(dir, ConfigFile), []byte(fine), 0644); err != nil {
		t.Fatal(err)
	}

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}

	if cfg.Tunnels[0].ID != "foundry-default" || cfg.Tunnels[1].ID != "mesa" {
		t.Errorf("ids = %q, %q, want them untouched", cfg.Tunnels[0].ID, cfg.Tunnels[1].ID)
	}
}
