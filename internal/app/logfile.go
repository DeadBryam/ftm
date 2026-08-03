package app

import (
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/sthbryan/ftm/internal/config"
)

const logFileName = "ftm.log"

func LogPath() string {
	return filepath.Join(config.ConfigDir(), logFileName)
}

func redirectLog() io.Closer {
	if err := os.MkdirAll(config.ConfigDir(), 0o755); err != nil {
		log.SetOutput(io.Discard)
		return nil
	}

	file, err := os.OpenFile(LogPath(), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		log.SetOutput(io.Discard)
		return nil
	}

	log.SetOutput(file)

	return file
}
