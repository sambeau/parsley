// Package parsley provides a public API for embedding the Parsley language interpreter.
package parsley

import (
	"fmt"
	"io"
	"strings"
	"sync"

	"github.com/sambeau/parsley/pkg/evaluator"
)

// Logger is an alias for evaluator.Logger for convenience
type Logger = evaluator.Logger

// stdoutLogger writes to stdout (default behavior)
// This is provided as a convenience - evaluator.DefaultLogger can also be used
type stdoutLogger struct{}

func (l *stdoutLogger) Log(values ...interface{}) {
	fmt.Print(formatLogValues(values...))
}

func (l *stdoutLogger) LogLine(values ...interface{}) {
	fmt.Println(formatLogValues(values...))
}

// StdoutLogger returns a logger that writes to stdout (default for CLI/REPL)
func StdoutLogger() Logger {
	return evaluator.DefaultLogger
}

// writerLogger writes to an io.Writer
type writerLogger struct {
	w io.Writer
}

func (l *writerLogger) Log(values ...interface{}) {
	fmt.Fprint(l.w, formatLogValues(values...))
}

func (l *writerLogger) LogLine(values ...interface{}) {
	fmt.Fprintln(l.w, formatLogValues(values...))
}

// WriterLogger returns a logger that writes to an io.Writer
func WriterLogger(w io.Writer) Logger {
	return &writerLogger{w: w}
}

// BufferedLogger captures log output for later retrieval
type BufferedLogger struct {
	mu    sync.Mutex
	lines []string
	buf   strings.Builder
}

// NewBufferedLogger creates a new buffered logger
func NewBufferedLogger() *BufferedLogger {
	return &BufferedLogger{
		lines: make([]string, 0),
	}
}

func (l *BufferedLogger) Log(values ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.buf.WriteString(formatLogValues(values...))
}

func (l *BufferedLogger) LogLine(values ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	// Flush any pending buffer content as a line
	line := l.buf.String() + formatLogValues(values...)
	l.lines = append(l.lines, line)
	l.buf.Reset()
}

// String returns all captured output as a single string
func (l *BufferedLogger) String() string {
	l.mu.Lock()
	defer l.mu.Unlock()
	result := strings.Join(l.lines, "\n")
	if len(l.lines) > 0 {
		result += "\n"
	}
	// Include any pending buffer content
	if l.buf.Len() > 0 {
		result += l.buf.String()
	}
	return result
}

// Lines returns all captured log lines
func (l *BufferedLogger) Lines() []string {
	l.mu.Lock()
	defer l.mu.Unlock()
	// Make a copy to avoid race conditions
	result := make([]string, len(l.lines))
	copy(result, l.lines)
	return result
}

// Reset clears all captured output
func (l *BufferedLogger) Reset() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.lines = l.lines[:0]
	l.buf.Reset()
}

// nullLogger discards all output
type nullLogger struct{}

func (l *nullLogger) Log(values ...interface{})     {}
func (l *nullLogger) LogLine(values ...interface{}) {}

// NullLogger returns a logger that discards all output
func NullLogger() Logger {
	return &nullLogger{}
}

// formatLogValues formats values for logging, similar to existing behavior
func formatLogValues(values ...interface{}) string {
	if len(values) == 0 {
		return ""
	}
	parts := make([]string, len(values))
	for i, v := range values {
		parts[i] = fmt.Sprint(v)
	}
	return strings.Join(parts, " ")
}

// DefaultLogger is the default logger used when none is specified
var DefaultLogger Logger = evaluator.DefaultLogger
