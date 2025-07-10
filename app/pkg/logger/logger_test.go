package logger

import (
	"bytes"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
)

func TestNew(t *testing.T) {
	tempDir := t.TempDir()
	
	tests := []struct {
		name    string
		config  Config
		wantErr bool
	}{
		{
			name: "File and console",
			config: Config{
				Level:      "debug",
				File:       filepath.Join(tempDir, "test.log"),
				MaxSize:    10,
				MaxBackups: 3,
				MaxAge:     7,
				Compress:   false,
				Console:    true,
			},
			wantErr: false,
		},
		{
			name: "Console only",
			config: Config{
				Level:   "info",
				Console: true,
			},
			wantErr: false,
		},
		{
			name: "File only",
			config: Config{
				Level:      "warn",
				File:       filepath.Join(tempDir, "test2.log"),
				MaxSize:    5,
				MaxBackups: 1,
				MaxAge:     1,
			},
			wantErr: false,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger, err := New(tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			
			if logger == nil && !tt.wantErr {
				t.Error("Expected logger to be created")
			}
		})
	}
}

func TestLogger_Levels(t *testing.T) {
	// Create logger with custom writer to capture output
	var buf bytes.Buffer
	logger := &Logger{
		level:  InfoLevel,
		logger: log.New(&buf, "", 0), // No flags for easier testing
	}
	
	// Test different levels
	logger.Debug("debug message")
	logger.Info("info message")
	logger.Warn("warn message")
	logger.Error("error message")
	
	output := buf.String()
	
	// Debug should not appear (level is Info)
	if strings.Contains(output, "debug message") {
		t.Error("Debug message should not appear when level is Info")
	}
	
	// Info, Warn, Error should appear
	if !strings.Contains(output, "info message") {
		t.Error("Info message should appear")
	}
	if !strings.Contains(output, "warn message") {
		t.Error("Warn message should appear")
	}
	if !strings.Contains(output, "error message") {
		t.Error("Error message should appear")
	}
}

func TestLogger_SetLevel(t *testing.T) {
	var buf bytes.Buffer
	logger := &Logger{
		level:  ErrorLevel,
		logger: log.New(&buf, "", 0),
	}
	
	// Only error should log
	logger.Info("info 1")
	logger.Error("error 1")
	
	// Change level to Info
	logger.SetLevel(InfoLevel)
	
	// Now info should log too
	logger.Info("info 2")
	logger.Error("error 2")
	
	output := buf.String()
	
	if strings.Contains(output, "info 1") {
		t.Error("First info message should not appear")
	}
	if !strings.Contains(output, "error 1") {
		t.Error("First error message should appear")
	}
	if !strings.Contains(output, "info 2") {
		t.Error("Second info message should appear")
	}
	if !strings.Contains(output, "error 2") {
		t.Error("Second error message should appear")
	}
}

func TestLogger_SetPrefix(t *testing.T) {
	var buf bytes.Buffer
	logger := &Logger{
		level:  InfoLevel,
		logger: log.New(&buf, "", 0),
	}
	
	logger.Info("without prefix")
	
	logger.SetPrefix("TEST")
	logger.Info("with prefix")
	
	output := buf.String()
	
	if strings.Contains(output, "[TEST] without prefix") {
		t.Error("First message should not have prefix")
	}
	if !strings.Contains(output, "[TEST] with prefix") {
		t.Error("Second message should have prefix")
	}
}

func TestParseLevel(t *testing.T) {
	tests := []struct {
		input    string
		expected Level
	}{
		{"debug", DebugLevel},
		{"info", InfoLevel},
		{"warn", WarnLevel},
		{"warning", WarnLevel},
		{"error", ErrorLevel},
		{"fatal", FatalLevel},
		{"unknown", InfoLevel}, // default
		{"", InfoLevel},        // default
	}
	
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if got := parseLevel(tt.input); got != tt.expected {
				t.Errorf("parseLevel(%q) = %v, want %v", tt.input, got, tt.expected)
			}
		})
	}
}

func TestInitialize(t *testing.T) {
	// Reset defaultLogger for testing
	defaultLogger = nil
	once = sync.Once{}
	
	tempDir := t.TempDir()
	cfg := Config{
		Level:   "debug",
		File:    filepath.Join(tempDir, "test.log"),
		Console: false,
	}
	
	err := Initialize(cfg)
	if err != nil {
		t.Fatalf("Initialize() error = %v", err)
	}
	
	// Test that default logger methods work
	Debug("debug test")
	Info("info test")
	Warn("warn test")
	Error("error test")
	
	// Check log file was created
	if _, err := os.Stat(cfg.File); os.IsNotExist(err) {
		t.Error("Log file was not created")
	}
	
	// Read log file
	content, err := os.ReadFile(cfg.File)
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}
	
	// Check content
	logContent := string(content)
	if !strings.Contains(logContent, "debug test") {
		t.Error("Debug message not found in log file")
	}
	if !strings.Contains(logContent, "info test") {
		t.Error("Info message not found in log file")
	}
	if !strings.Contains(logContent, "warn test") {
		t.Error("Warn message not found in log file")
	}
	if !strings.Contains(logContent, "error test") {
		t.Error("Error message not found in log file")
	}
}