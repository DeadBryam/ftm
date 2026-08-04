package updater

import "testing"

func TestMethodFor(t *testing.T) {
	tests := []struct {
		name     string
		execPath string
		goos     string
		desktop  bool
		want     Method
	}{
		{
			name:     "plain binary on linux",
			execPath: "/usr/local/bin/ftm",
			goos:     "linux",
			want:     MethodSelf,
		},
		{
			name:     "msix install is owned by the store",
			execPath: `C:\Program Files\WindowsApps\sthbryan.ftm_0.11.0.0_x64__abc\ftm-desktop.exe`,
			goos:     "windows",
			want:     MethodStore,
		},
		{
			name:     "windowsapps path outside windows is not a store install",
			execPath: "/home/me/WindowsApps/ftm",
			goos:     "linux",
			want:     MethodSelf,
		},
		{
			name:     "homebrew cellar",
			execPath: "/opt/homebrew/Cellar/ftm/0.11.0/bin/ftm",
			goos:     "darwin",
			want:     MethodHomebrew,
		},
		{
			name:     "macos app bundle",
			execPath: "/Applications/Foundry Tunnel Manager.app/Contents/MacOS/ftm-desktop",
			goos:     "darwin",
			want:     MethodDownload,
		},
		{
			name:     "loose desktop binary identifies itself",
			execPath: `C:\Users\me\Downloads\ftm-desktop-windows.exe`,
			goos:     "windows",
			desktop:  true,
			want:     MethodDownload,
		},
		{
			name:     "store install wins over the desktop marker",
			execPath: `C:\Program Files\WindowsApps\sthbryan.ftm\ftm-desktop.exe`,
			goos:     "windows",
			desktop:  true,
			want:     MethodStore,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := methodFor(tt.execPath, tt.goos, tt.desktop); got != tt.want {
				t.Errorf("methodFor(%q, %q, %v) = %q, want %q", tt.execPath, tt.goos, tt.desktop, got, tt.want)
			}
		})
	}
}

func TestApplyRefusesNonSelfInstalls(t *testing.T) {
	for _, method := range []Method{MethodStore, MethodHomebrew, MethodDownload} {
		if err := New("").Apply(&Info{Method: method}); err != ErrNotSelfUpdatable {
			t.Errorf("Apply with method %q = %v, want ErrNotSelfUpdatable", method, err)
		}
	}
}
