//go:build darwin

package i18n

import (
	"context"
	"os/exec"
	"time"
)

//

func detectPlatformLang() string {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	out, err := exec.CommandContext(ctx, "defaults", "read", "-g", "AppleLanguages").Output()
	if err != nil {
		return ""
	}

	return parseAppleLanguages(string(out))
}
