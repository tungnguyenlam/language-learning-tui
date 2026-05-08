package app

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

// LogLevel represents the severity level for logging
type LogLevel int

const (
	// Log levels
	LogLevelDebug LogLevel = iota
	LogLevelInfo
	LogLevelWarn
	LogLevelError
)

// ParseLogLevel converts a string to a LogLevel
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

// String returns the string representation of a LogLevel
func (l LogLevel) String() string {
	switch l {
	case LogLevelDebug:
		return "DEBUG"
	case LogLevelWarn:
		return "WARN"
	case LogLevelError:
		return "ERROR"
	default:
		return "INFO"
	}
}

// LeveledLogger wraps the standard logger with level-based filtering
type LeveledLogger struct {
	logger   *log.Logger
	minLevel LogLevel
}

// NewLeveledLogger creates a new leveled logger
func NewLeveledLogger(logger *log.Logger, minLevel LogLevel) *LeveledLogger {
	return &LeveledLogger{
		logger:   logger,
		minLevel: minLevel,
	}
}

// Debug logs a debug message
func (l *LeveledLogger) Debug(format string, v ...interface{}) {
	if l.minLevel <= LogLevelDebug {
		l.logger.Printf("[DEBUG] "+format, v...)
	}
}

// Info logs an info message
func (l *LeveledLogger) Info(format string, v ...interface{}) {
	if l.minLevel <= LogLevelInfo {
		l.logger.Printf("[INFO] "+format, v...)
	}
}

// Warn logs a warning message
func (l *LeveledLogger) Warn(format string, v ...interface{}) {
	if l.minLevel <= LogLevelWarn {
		l.logger.Printf("[WARN] "+format, v...)
	}
}

// Error logs an error message
func (l *LeveledLogger) Error(format string, v ...interface{}) {
	if l.minLevel <= LogLevelError {
		l.logger.Printf("[ERROR] "+format, v...)
	}
}

// Fatal logs a fatal error and exits
func (l *LeveledLogger) Fatal(format string, v ...interface{}) {
	l.logger.Fatalf("[FATAL] "+format, v...)
}

func OpenLog(dataDir string) (*os.File, *log.Logger, error) {
	file, err := os.OpenFile(filepath.Join(dataDir, LogFileName), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return nil, nil, err
	}
	logger := log.New(file, "deutsch-tui ", log.LstdFlags|log.LUTC|log.Lshortfile)
	return file, logger, nil
}
