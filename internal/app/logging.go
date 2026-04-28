package app

import (
	"log"
	"os"
	"path/filepath"
)

func OpenLog(dataDir string) (*os.File, *log.Logger, error) {
	file, err := os.OpenFile(filepath.Join(dataDir, LogFileName), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return nil, nil, err
	}
	logger := log.New(file, "deutsch-tui ", log.LstdFlags|log.LUTC|log.Lshortfile)
	return file, logger, nil
}
