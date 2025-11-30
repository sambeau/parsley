package main

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/sambeau/parsley/pkg/evaluator"
	"github.com/sambeau/parsley/pkg/lexer"
	"github.com/sambeau/parsley/pkg/parser"
)

// Helper to run Parsley code with simulated stdin
func runWithStdin(t *testing.T, code string, stdinData string) (string, string) {
	// Save original stdin and stdout
	oldStdin := os.Stdin
	oldStdout := os.Stdout
	oldStderr := os.Stderr
	defer func() {
		os.Stdin = oldStdin
		os.Stdout = oldStdout
		os.Stderr = oldStderr
	}()

	// Create a pipe for stdin
	stdinReader, stdinWriter, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create stdin pipe: %v", err)
	}
	os.Stdin = stdinReader

	// Write stdin data in a goroutine
	go func() {
		defer stdinWriter.Close()
		stdinWriter.WriteString(stdinData)
	}()

	// Create pipes for stdout and stderr
	stdoutReader, stdoutWriter, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create stdout pipe: %v", err)
	}
	os.Stdout = stdoutWriter

	stderrReader, stderrWriter, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create stderr pipe: %v", err)
	}
	os.Stderr = stderrWriter

	// Run the code
	l := lexer.New(code)
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		t.Fatalf("parser errors: %v", p.Errors())
	}

	env := evaluator.NewEnvironment()
	evaluator.Eval(program, env)

	// Close writers and read output
	stdoutWriter.Close()
	stderrWriter.Close()

	var stdoutBuf, stderrBuf bytes.Buffer
	io.Copy(&stdoutBuf, stdoutReader)
	io.Copy(&stderrBuf, stderrReader)

	return stdoutBuf.String(), stderrBuf.String()
}

func TestStdinJSONRead(t *testing.T) {
	code := `let data <== JSON(@-)
data ==> JSON(@-)`

	stdout, _ := runWithStdin(t, code, `{"name":"test","value":42}`)

	if !strings.Contains(stdout, `"name"`) || !strings.Contains(stdout, `"test"`) {
		t.Errorf("Expected JSON output with name field, got: %s", stdout)
	}
	if !strings.Contains(stdout, `"value"`) || !strings.Contains(stdout, `42`) {
		t.Errorf("Expected JSON output with value field, got: %s", stdout)
	}
}

func TestStdinTextRead(t *testing.T) {
	code := `let data <== text(@-)
data ==> text(@-)`

	stdout, _ := runWithStdin(t, code, "Hello World")

	if stdout != "Hello World" {
		t.Errorf("Expected 'Hello World', got: %s", stdout)
	}
}

func TestStdinLinesRead(t *testing.T) {
	code := `let data <== lines(@-)
len(data) ==> text(@-)`

	stdout, _ := runWithStdin(t, code, "line1\nline2\nline3")

	if stdout != "3" {
		t.Errorf("Expected '3' lines, got: %s", stdout)
	}
}

func TestStdinAlias(t *testing.T) {
	code := `let data <== JSON(@stdin)
data ==> JSON(@stdout)`

	stdout, _ := runWithStdin(t, code, `{"test":true}`)

	if !strings.Contains(stdout, `"test"`) || !strings.Contains(stdout, `true`) {
		t.Errorf("Expected JSON output with test field, got: %s", stdout)
	}
}

func TestStderrWrite(t *testing.T) {
	code := `"error message" ==> text(@stderr)`

	_, stderr := runWithStdin(t, code, "")

	if stderr != "error message" {
		t.Errorf("Expected 'error message' on stderr, got: %s", stderr)
	}
}

func TestStdioSeparation(t *testing.T) {
	code := `"stdout message" ==> text(@stdout)
"stderr message" ==> text(@stderr)`

	stdout, stderr := runWithStdin(t, code, "")

	if stdout != "stdout message" {
		t.Errorf("Expected 'stdout message' on stdout, got: %s", stdout)
	}
	if stderr != "stderr message" {
		t.Errorf("Expected 'stderr message' on stderr, got: %s", stderr)
	}
}
