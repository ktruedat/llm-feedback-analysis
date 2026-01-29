package log

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/rs/zerolog"
)

func TestNewLogger_Production(t *testing.T) {
	logger := NewLogger("production")
	if logger == nil {
		t.Fatal("Expected logger to be non-nil")
	}

	// Verify it implements the Logger interface by calling a method
	logger.Info("test")
}

func TestNewLogger_Development(t *testing.T) {
	logger := NewLogger("development")
	if logger == nil {
		t.Fatal("Expected logger to be non-nil")
	}

	// Verify it implements the Logger interface by calling a method
	logger.Info("test")
}

func TestNewLogger_DefaultToDevelopment(t *testing.T) {
	logger := NewLogger("unknown")
	if logger == nil {
		t.Fatal("Expected logger to be non-nil")
	}

	// Should default to development logger
	logger.Info("test")
}

func TestLogger_Error(t *testing.T) {
	var buf bytes.Buffer
	zlog := zerolog.New(&buf).Level(zerolog.DebugLevel)
	logger := &logger{zlog: &zlog}

	testErr := errors.New("test error")
	logger.Error("error message", testErr, "key1", "value1", "key2", 42)

	output := buf.String()
	if !strings.Contains(output, "error message") {
		t.Errorf("Expected log output to contain 'error message', got: %s", output)
	}
	if !strings.Contains(output, "test error") {
		t.Errorf("Expected log output to contain 'test error', got: %s", output)
	}
	if !strings.Contains(output, "key1") {
		t.Errorf("Expected log output to contain 'key1', got: %s", output)
	}
	if !strings.Contains(output, "value1") {
		t.Errorf("Expected log output to contain 'value1', got: %s", output)
	}
}

func TestLogger_Info(t *testing.T) {
	var buf bytes.Buffer
	zlog := zerolog.New(&buf).Level(zerolog.DebugLevel)
	logger := &logger{zlog: &zlog}

	logger.Info("info message", "key1", "value1", "key2", 42)

	output := buf.String()
	if !strings.Contains(output, "info message") {
		t.Errorf("Expected log output to contain 'info message', got: %s", output)
	}
	if !strings.Contains(output, "key1") {
		t.Errorf("Expected log output to contain 'key1', got: %s", output)
	}
	if !strings.Contains(output, "value1") {
		t.Errorf("Expected log output to contain 'value1', got: %s", output)
	}
}

func TestLogger_Warning(t *testing.T) {
	var buf bytes.Buffer
	zlog := zerolog.New(&buf).Level(zerolog.DebugLevel)
	logger := &logger{zlog: &zlog}

	logger.Warning("warning message", "key1", "value1")

	output := buf.String()
	if !strings.Contains(output, "warning message") {
		t.Errorf("Expected log output to contain 'warning message', got: %s", output)
	}
	if !strings.Contains(output, "key1") {
		t.Errorf("Expected log output to contain 'key1', got: %s", output)
	}
}

func TestLogger_Debug(t *testing.T) {
	var buf bytes.Buffer
	zlog := zerolog.New(&buf).Level(zerolog.DebugLevel)
	logger := &logger{zlog: &zlog}

	logger.Debug("debug message", "key1", "value1")

	output := buf.String()
	if !strings.Contains(output, "debug message") {
		t.Errorf("Expected log output to contain 'debug message', got: %s", output)
	}
	if !strings.Contains(output, "key1") {
		t.Errorf("Expected log output to contain 'key1', got: %s", output)
	}
}

func TestLogger_handleArgs_ValidPairs(t *testing.T) {
	var buf bytes.Buffer
	zlog := zerolog.New(&buf).Level(zerolog.DebugLevel)
	logger := &logger{zlog: &zlog}

	event := zlog.Info()
	logger.handleArgs(event, "key1", "value1", "key2", 42, "key3", true)

	event.Msg("test")
	output := buf.String()

	if !strings.Contains(output, "key1") {
		t.Errorf("Expected log output to contain 'key1', got: %s", output)
	}
	if !strings.Contains(output, "value1") {
		t.Errorf("Expected log output to contain 'value1', got: %s", output)
	}
	if !strings.Contains(output, "key2") {
		t.Errorf("Expected log output to contain 'key2', got: %s", output)
	}
}

func TestLogger_handleArgs_InvalidKeyType(t *testing.T) {
	var buf bytes.Buffer
	zlog := zerolog.New(&buf).Level(zerolog.DebugLevel)
	logger := &logger{zlog: &zlog}

	event := zlog.Info()
	logger.handleArgs(event, 123, "value1")

	event.Msg("test")
	output := buf.String()

	if !strings.Contains(output, "invalid_key_type") {
		t.Errorf("Expected log output to contain 'invalid_key_type', got: %s", output)
	}
}

func TestLogger_handleArgs_DanglingArg(t *testing.T) {
	var buf bytes.Buffer
	zlog := zerolog.New(&buf).Level(zerolog.DebugLevel)
	logger := &logger{zlog: &zlog}

	event := zlog.Info()
	logger.handleArgs(event, "key1", "value1", "dangling")

	event.Msg("test")
	output := buf.String()

	if !strings.Contains(output, "dangling_arg") {
		t.Errorf("Expected log output to contain 'dangling_arg', got: %s", output)
	}
}

func TestLogger_handleArgs_EmptyArgs(t *testing.T) {
	var buf bytes.Buffer
	zlog := zerolog.New(&buf).Level(zerolog.DebugLevel)
	logger := &logger{zlog: &zlog}

	event := zlog.Info()
	logger.handleArgs(event)

	event.Msg("test")
	output := buf.String()

	if !strings.Contains(output, "test") {
		t.Errorf("Expected log output to contain 'test', got: %s", output)
	}
}

func TestLogger_NewGroup(t *testing.T) {
	var buf bytes.Buffer
	zlog := zerolog.New(&buf).Level(zerolog.DebugLevel)
	logger := &logger{zlog: &zlog}

	groupLogger := logger.NewGroup("test-group")
	if groupLogger == nil {
		t.Fatal("Expected NewGroup to return a non-nil logger")
	}

	groupLogger.Info("test message")
	output := buf.String()

	if !strings.Contains(output, "group") {
		t.Errorf("Expected log output to contain 'group', got: %s", output)
	}
	if !strings.Contains(output, "test-group") {
		t.Errorf("Expected log output to contain 'test-group', got: %s", output)
	}
}

func TestLogger_With(t *testing.T) {
	var buf bytes.Buffer
	zlog := zerolog.New(&buf).Level(zerolog.DebugLevel)
	logger := &logger{zlog: &zlog}

	withLogger := logger.With("key1", "value1", "key2", 42)
	if withLogger == nil {
		t.Fatal("Expected With to return a non-nil logger")
	}

	withLogger.Info("test message")
	output := buf.String()

	if !strings.Contains(output, "additional_info") {
		t.Errorf("Expected log output to contain 'additionalInfo', got: %s", output)
	}
	if !strings.Contains(output, "key1") {
		t.Errorf("Expected log output to contain 'key1', got: %s", output)
	}
	if !strings.Contains(output, "value1") {
		t.Errorf("Expected log output to contain 'value1', got: %s", output)
	}
}

func TestLogger_With_EmptyArgs(t *testing.T) {
	var buf bytes.Buffer
	zlog := zerolog.New(&buf).Level(zerolog.DebugLevel)
	logger := &logger{zlog: &zlog}

	withLogger := logger.With()
	if withLogger == nil {
		t.Fatal("Expected With to return a non-nil logger")
	}

	withLogger.Info("test message")
	output := buf.String()

	if !strings.Contains(output, "test message") {
		t.Errorf("Expected log output to contain 'test message', got: %s", output)
	}
}

func TestLogger_Chaining(t *testing.T) {
	var buf bytes.Buffer
	zlog := zerolog.New(&buf).Level(zerolog.DebugLevel)
	logger := &logger{zlog: &zlog}

	// Chain NewGroup and With
	chainedLogger := logger.NewGroup("api").With("request_id", "12345")
	if chainedLogger == nil {
		t.Fatal("Expected chained logger to be non-nil")
	}

	chainedLogger.Info("chained message")
	output := buf.String()

	if !strings.Contains(output, "group") {
		t.Errorf("Expected log output to contain 'group', got: %s", output)
	}
	if !strings.Contains(output, "api") {
		t.Errorf("Expected log output to contain 'api', got: %s", output)
	}
	if !strings.Contains(output, "additional_info") {
		t.Errorf("Expected log output to contain 'additionalInfo', got: %s", output)
	}
}

func TestLogger_Error_NilError(t *testing.T) {
	var buf bytes.Buffer
	zlog := zerolog.New(&buf).Level(zerolog.DebugLevel)
	logger := &logger{zlog: &zlog}

	// Should not panic with nil error
	logger.Error("error message", nil, "key1", "value1")

	output := buf.String()
	if !strings.Contains(output, "error message") {
		t.Errorf("Expected log output to contain 'error message', got: %s", output)
	}
}

func TestLogger_AllLevels(t *testing.T) {
	var buf bytes.Buffer
	zlog := zerolog.New(&buf).Level(zerolog.TraceLevel)
	logger := &logger{zlog: &zlog}

	logger.Debug("debug", "level", "debug")
	logger.Info("info", "level", "info")
	logger.Warning("warning", "level", "warning")
	logger.Error("error", errors.New("test"), "level", "error")

	output := buf.String()
	if !strings.Contains(output, "debug") {
		t.Errorf("Expected log output to contain 'debug', got: %s", output)
	}
	if !strings.Contains(output, "info") {
		t.Errorf("Expected log output to contain 'info', got: %s", output)
	}
	if !strings.Contains(output, "warning") {
		t.Errorf("Expected log output to contain 'warning', got: %s", output)
	}
	if !strings.Contains(output, "error") {
		t.Errorf("Expected log output to contain 'error', got: %s", output)
	}
}
