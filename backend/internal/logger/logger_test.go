package logger

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNewWritesJSONLogToFile(t *testing.T) {
	filePath := filepath.Join(t.TempDir(), "nested", "app.log")
	log, closer, err := New(filePath, "debug")
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	log.Debug("test message", "user_id", 1)
	if err := closer.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}

	text := string(content)
	if !strings.Contains(text, `"level":"DEBUG"`) ||
		!strings.Contains(text, `"msg":"test message"`) ||
		!strings.Contains(text, `"user_id":1`) {
		t.Fatalf("log content = %s", text)
	}
}

func TestNewRejectsUnknownLevel(t *testing.T) {
	_, _, err := New(filepath.Join(t.TempDir(), "app.log"), "trace")
	if err == nil {
		t.Fatal("New() error = nil, want unsupported level error")
	}
}
