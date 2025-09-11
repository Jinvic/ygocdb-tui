package log

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// LogLevel represents the severity level of a log message
type LogLevel int

const (
	// OffLevel means no logging
	OffLevel LogLevel = iota
	// ErrorLevel logs only errors
	ErrorLevel
	// WarnLevel logs warnings and errors
	WarnLevel
	// InfoLevel logs info, warnings, and errors
	InfoLevel
	// DebugLevel logs everything including debug messages
	DebugLevel
)

// String returns the string representation of a LogLevel
func (l LogLevel) String() string {
	switch l {
	case OffLevel:
		return "OFF"
	case ErrorLevel:
		return "ERROR"
	case WarnLevel:
		return "WARN"
	case InfoLevel:
		return "INFO"
	case DebugLevel:
		return "DEBUG"
	default:
		return "UNKNOWN"
	}
}

// Logger represents a logger instance
type Logger struct {
	level  LogLevel
	file   *os.File
	logger *log.Logger
	mu     sync.Mutex
}

var (
	// Default logger instance
	defaultLogger *Logger
	// Ensure singleton initialization
	once sync.Once
)

// Init initializes the logger with the specified log level
func Init(level LogLevel) error {
	var err error
	once.Do(func() {
		defaultLogger, err = NewLogger(level)
	})
	return err
}

// NewLogger creates a new logger instance
func NewLogger(level LogLevel) (*Logger, error) {
	if level == OffLevel {
		return &Logger{level: level}, nil
	}

	// Create logs directory if it doesn't exist
	logsDir := "logs"
	if err := os.MkdirAll(logsDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create logs directory: %w", err)
	}

	// Create log file with timestamp
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	filename := filepath.Join(logsDir, fmt.Sprintf("ygocdb-tui_%s.log", timestamp))
	
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to create log file: %w", err)
	}

	logger := &Logger{
		level:  level,
		file:   file,
		logger: log.New(file, "", log.LstdFlags|log.Lshortfile),
	}

	// Log initialization
	logger.logger.Printf("Logger initialized with level: %s", level.String())

	return logger, nil
}

// SetLevel sets the logging level
func (l *Logger) SetLevel(level LogLevel) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.level = level
}

// Debug logs a debug message
func (l *Logger) Debug(format string, v ...interface{}) {
	if l.level >= DebugLevel {
		l.mu.Lock()
		defer l.mu.Unlock()
		l.logger.Printf("[DEBUG] "+format, v...)
	}
}

// Info logs an info message
func (l *Logger) Info(format string, v ...interface{}) {
	if l.level >= InfoLevel {
		l.mu.Lock()
		defer l.mu.Unlock()
		l.logger.Printf("[INFO] "+format, v...)
	}
}

// Warn logs a warning message
func (l *Logger) Warn(format string, v ...interface{}) {
	if l.level >= WarnLevel {
		l.mu.Lock()
		defer l.mu.Unlock()
		l.logger.Printf("[WARN] "+format, v...)
	}
}

// Error logs an error message
func (l *Logger) Error(format string, v ...interface{}) {
	if l.level >= ErrorLevel {
		l.mu.Lock()
		defer l.mu.Unlock()
		l.logger.Printf("[ERROR] "+format, v...)
	}
}

// Close closes the logger and the underlying file
func (l *Logger) Close() error {
	if l.file != nil {
		l.logger.Printf("Logger closed")
		return l.file.Close()
	}
	return nil
}

// Global functions that use the default logger

// Debug logs a debug message using the default logger
func Debug(format string, v ...interface{}) {
	if defaultLogger != nil {
		defaultLogger.Debug(format, v...)
	}
}

// Info logs an info message using the default logger
func Info(format string, v ...interface{}) {
	if defaultLogger != nil {
		defaultLogger.Info(format, v...)
	}
}

// Warn logs a warning message using the default logger
func Warn(format string, v ...interface{}) {
	if defaultLogger != nil {
		defaultLogger.Warn(format, v...)
	}
}

// Error logs an error message using the default logger
func Error(format string, v ...interface{}) {
	if defaultLogger != nil {
		defaultLogger.Error(format, v...)
	}
}

// SetLevel sets the logging level for the default logger
func SetLevel(level LogLevel) {
	if defaultLogger != nil {
		defaultLogger.SetLevel(level)
	}
}

// Close closes the default logger
func Close() error {
	if defaultLogger != nil {
		return defaultLogger.Close()
	}
	return nil
}