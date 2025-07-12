package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sync"

	"gopkg.in/natefinch/lumberjack.v2"
)

// Level represents log level
type Level int

const (
	// DebugLevel is for debug messages
	DebugLevel Level = iota
	// InfoLevel is for informational messages
	InfoLevel
	// WarnLevel is for warning messages
	WarnLevel
	// ErrorLevel is for error messages
	ErrorLevel
	// FatalLevel is for fatal messages
	FatalLevel
)

// Logger wraps standard logger with level support
type Logger struct {
	mu         sync.RWMutex
	level      Level
	logger     *log.Logger
	prefix     string
	outputs    []io.Writer
	fileOutput *lumberjack.Logger // Keep reference for closing
}

// Config represents logger configuration
type Config struct {
	Level      string // debug, info, warn, error
	File       string // log file path
	MaxSize    int    // megabytes
	MaxBackups int    // number of old files to keep
	MaxAge     int    // days
	Compress   bool   // compress old files
	Console    bool   // also log to console
}

var (
	defaultLogger *Logger
	once          sync.Once
)

// Initialize sets up the default logger
func Initialize(cfg Config) error {
	var err error
	once.Do(func() {
		defaultLogger, err = New(cfg)
	})
	return err
}

// New creates a new logger
func New(cfg Config) (*Logger, error) {
	outputs := make([]io.Writer, 0)

	// Set up file output with rotation
	if cfg.File != "" {
		// Ensure log directory exists
		logDir := filepath.Dir(cfg.File)
		if err := os.MkdirAll(logDir, 0o755); err != nil {
			return nil, fmt.Errorf("failed to create log directory: %w", err)
		}

		// Configure lumberjack
		fileOutput := &lumberjack.Logger{
			Filename:   cfg.File,
			MaxSize:    cfg.MaxSize, // megabytes
			MaxBackups: cfg.MaxBackups,
			MaxAge:     cfg.MaxAge, // days
			Compress:   cfg.Compress,
		}
		outputs = append(outputs, fileOutput)
	}

	// Add console output if requested
	if cfg.Console || cfg.File == "" {
		outputs = append(outputs, os.Stderr)
	}

	// Create multi-writer
	multiWriter := io.MultiWriter(outputs...)

	// Create logger
	l := &Logger{
		level:   parseLevel(cfg.Level),
		logger:  log.New(multiWriter, "", log.LstdFlags|log.Lmicroseconds),
		outputs: outputs,
	}

	// Keep reference to file output for closing
	if cfg.File != "" && len(outputs) > 0 {
		if fileLogger, ok := outputs[0].(*lumberjack.Logger); ok {
			l.fileOutput = fileLogger
		}
	}

	return l, nil
}

// SetLevel sets the log level
func (l *Logger) SetLevel(level Level) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.level = level
}

// SetPrefix sets the logger prefix
func (l *Logger) SetPrefix(prefix string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.prefix = prefix
}

// Debug logs a debug message
func (l *Logger) Debug(format string, v ...any) {
	l.log(DebugLevel, "DEBUG", format, v...)
}

// Info logs an info message
func (l *Logger) Info(format string, v ...any) {
	l.log(InfoLevel, "INFO", format, v...)
}

// Warn logs a warning message
func (l *Logger) Warn(format string, v ...any) {
	l.log(WarnLevel, "WARN", format, v...)
}

// Error logs an error message
func (l *Logger) Error(format string, v ...any) {
	l.log(ErrorLevel, "ERROR", format, v...)
}

// Fatal logs a fatal message and exits
func (l *Logger) Fatal(format string, v ...any) {
	l.log(FatalLevel, "FATAL", format, v...)
	os.Exit(1)
}

// log is the internal logging method
func (l *Logger) log(level Level, levelStr, format string, v ...any) {
	l.mu.RLock()
	currentLevel := l.level
	prefix := l.prefix
	l.mu.RUnlock()

	if level < currentLevel {
		return
	}

	// Get caller information
	_, file, line, ok := runtime.Caller(2)
	if ok {
		file = filepath.Base(file)
	} else {
		file = "???"
		line = 0
	}

	// Format message
	msg := fmt.Sprintf(format, v...)
	if prefix != "" {
		msg = fmt.Sprintf("[%s] %s", prefix, msg)
	}

	// Log with level, file, and line
	l.logger.Printf("[%s] %s:%d %s", levelStr, file, line, msg)
}

// Close closes the logger and releases any resources
func (l *Logger) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.fileOutput != nil {
		return l.fileOutput.Close()
	}
	return nil
}

// parseLevel parses string level to Level type
func parseLevel(level string) Level {
	switch level {
	case "debug":
		return DebugLevel
	case "info":
		return InfoLevel
	case "warn", "warning":
		return WarnLevel
	case "error":
		return ErrorLevel
	case "fatal":
		return FatalLevel
	default:
		return InfoLevel
	}
}

// Default logger methods

// Debug logs a debug message using the default logger
func Debug(format string, v ...any) {
	if defaultLogger != nil {
		defaultLogger.Debug(format, v...)
	}
}

// Info logs an info message using the default logger
func Info(format string, v ...any) {
	if defaultLogger != nil {
		defaultLogger.Info(format, v...)
	}
}

// Warn logs a warning message using the default logger
func Warn(format string, v ...any) {
	if defaultLogger != nil {
		defaultLogger.Warn(format, v...)
	}
}

// Error logs an error message using the default logger
func Error(format string, v ...any) {
	if defaultLogger != nil {
		defaultLogger.Error(format, v...)
	}
}

// Fatal logs a fatal message using the default logger and exits
func Fatal(format string, v ...any) {
	if defaultLogger != nil {
		defaultLogger.Fatal(format, v...)
	} else {
		log.Fatalf(format, v...)
	}
}

// Close closes the default logger
func Close() error {
	if defaultLogger != nil {
		err := defaultLogger.Close()
		defaultLogger = nil
		return err
	}
	return nil
}
