package app

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

type LogLevel int

const (
	LogLevelDebug LogLevel = iota
	LogLevelInfo
	LogLevelWarn
	LogLevelError
)

func ParseLogLevel(level string) LogLevel {
	switch strings.ToLower(strings.TrimSpace(level)) {
	case "debug":
		return LogLevelDebug
	case "warn", "warning":
		return LogLevelWarn
	case "error":
		return LogLevelError
	default:
		return LogLevelInfo
	}
}

type LeveledLogger struct {
	logger   *log.Logger
	minLevel LogLevel
}

func NewLeveledLogger(logger *log.Logger, minLevel LogLevel) *LeveledLogger {
	return &LeveledLogger{logger: logger, minLevel: minLevel}
}

func (l *LeveledLogger) Debug(format string, v ...interface{}) {
	if l.minLevel <= LogLevelDebug {
		l.logger.Printf("[DEBUG] "+format, v...)
	}
}

func (l *LeveledLogger) Info(format string, v ...interface{}) {
	if l.minLevel <= LogLevelInfo {
		l.logger.Printf("[INFO] "+format, v...)
	}
}

func (l *LeveledLogger) Error(format string, v ...interface{}) {
	if l.minLevel <= LogLevelError {
		l.logger.Printf("[ERROR] "+format, v...)
	}
}

func OpenLog(dataDir string) (*os.File, *log.Logger, error) {
	file, err := os.OpenFile(filepath.Join(dataDir, LogFileName), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return nil, nil, err
	}
	logger := log.New(file, "deutsch-tui ", log.LstdFlags|log.LUTC|log.Lshortfile)
	return file, logger, nil
}
