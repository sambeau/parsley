package main

import (
	"testing"

	"github.com/sambeau/parsley/pkg/evaluator"
)

// SFTP tests are skipped until we have a working SFTP server for integration testing.
// To run these tests, set up an SFTP server and remove the t.Skip() calls.
// See docs/TODO.md for integration testing setup requirements.

// TestSFTPConnectionCreation tests SFTP() builtin connection creation
func TestSFTPConnectionCreation(t *testing.T) {
	t.Skip("SFTP tests suspended - requires SFTP server for integration testing")

	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "SFTP with password",
			input:   `SFTP("sftp://user:pass@example.com/")`,
			wantErr: true, // Will fail without real server, but should parse correctly
		},
		{
			name:    "SFTP with SSH key",
			input:   `SFTP("sftp://user@example.com/", {key: @~/.ssh/id_rsa})`,
			wantErr: true, // Will fail without real server
		},
		{
			name:    "SFTP with timeout",
			input:   `SFTP("sftp://user:pass@example.com/", {timeout: @5s})`,
			wantErr: true,
		},
		{
			name:    "SFTP with port",
			input:   `SFTP("sftp://user:pass@example.com:2222/")`,
			wantErr: true,
		},
		{
			name:    "Invalid URL",
			input:   `SFTP("not-a-url")`,
			wantErr: true,
		},
		{
			name:    "Missing credentials",
			input:   `SFTP("sftp://example.com/")`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := testEvalHelper(tt.input)

			if tt.wantErr {
				_, isErr := result.(*evaluator.Error)
				if !isErr {
					// Check for SFTP_CONNECTION_OBJ type (connection might be created but fail to connect)
					t.Logf("Result type: %s", result.Type())
				}
			}
		})
	}
}

// TestSFTPCallableSyntax tests conn(@/path) syntax
func TestSFTPCallableSyntax(t *testing.T) {
	t.Skip("SFTP tests suspended - requires SFTP server for integration testing")

	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "Create file handle",
			input:   `let conn = SFTP("sftp://user:pass@example.com/"); conn(@/data.json)`,
			wantErr: true, // Will fail without real server
		},
		{
			name:    "Multiple file handles",
			input:   `let conn = SFTP("sftp://user:pass@example.com/"); let f1 = conn(@/file1.txt); let f2 = conn(@/file2.txt); f2`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := testEvalHelper(tt.input)

			if tt.wantErr {
				// Check for error or SFTP_FILE_HANDLE_OBJ type
				t.Logf("Result type: %s", result.Type())
			}
		})
	}
}

// TestSFTPFormatAccessors tests format accessor properties
func TestSFTPFormatAccessors(t *testing.T) {
	t.Skip("SFTP tests suspended - requires SFTP server for integration testing")

	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "JSON accessor",
			input:   `let conn = SFTP("sftp://user:pass@example.com/"); conn(@/data.json).json`,
			wantErr: true, // Will fail without real server
		},
		{
			name:    "Text accessor",
			input:   `let conn = SFTP("sftp://user:pass@example.com/"); conn(@/file.txt).text`,
			wantErr: true,
		},
		{
			name:    "CSV accessor",
			input:   `let conn = SFTP("sftp://user:pass@example.com/"); conn(@/data.csv).csv`,
			wantErr: true,
		},
		{
			name:    "Lines accessor",
			input:   `let conn = SFTP("sftp://user:pass@example.com/"); conn(@/log.txt).lines`,
			wantErr: true,
		},
		{
			name:    "Bytes accessor",
			input:   `let conn = SFTP("sftp://user:pass@example.com/"); conn(@/binary.dat).bytes`,
			wantErr: true,
		},
		{
			name:    "File accessor (auto-detect)",
			input:   `let conn = SFTP("sftp://user:pass@example.com/"); conn(@/data.json).file`,
			wantErr: true,
		},
		{
			name:    "Dir accessor",
			input:   `let conn = SFTP("sftp://user:pass@example.com/"); conn(@/uploads).dir`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := testEvalHelper(tt.input)
			t.Logf("%s result type: %s", tt.name, result.Type())
		})
	}
}

// TestSFTPReadOperatorSyntax tests <=/= operator syntax
func TestSFTPReadOperatorSyntax(t *testing.T) {
	t.Skip("SFTP tests suspended - requires SFTP server for integration testing")

	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "Read JSON file",
			input:   `let conn = SFTP("sftp://user:pass@example.com/"); data <=/= conn(@/data.json).json`,
			wantErr: true, // Will fail without real server
		},
		{
			name:    "Read text file",
			input:   `let conn = SFTP("sftp://user:pass@example.com/"); content <=/= conn(@/file.txt).text`,
			wantErr: true,
		},
		{
			name:    "Read with error capture",
			input:   `let conn = SFTP("sftp://user:pass@example.com/"); {data, error} <=/= conn(@/data.json).json; error`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := testEvalHelper(tt.input)
			t.Logf("%s result type: %s", tt.name, result.Type())
		})
	}
}

// TestSFTPWriteOperatorSyntax tests =/=> operator syntax
func TestSFTPWriteOperatorSyntax(t *testing.T) {
	t.Skip("SFTP tests suspended - requires SFTP server for integration testing")

	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "Write JSON file",
			input:   `let conn = SFTP("sftp://user:pass@example.com/"); {name: "test"} =/=> conn(@/data.json).json`,
			wantErr: true,
		},
		{
			name:    "Write text file",
			input:   `let conn = SFTP("sftp://user:pass@example.com/"); "Hello World" =/=> conn(@/file.txt).text`,
			wantErr: true,
		},
		{
			name:    "Write with error capture",
			input:   `let conn = SFTP("sftp://user:pass@example.com/"); error = {test: 123} =/=> conn(@/data.json).json; error`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := testEvalHelper(tt.input)
			t.Logf("%s result type: %s", tt.name, result.Type())
		})
	}
}

// TestSFTPAppendOperatorSyntax tests =/=>> operator syntax
func TestSFTPAppendOperatorSyntax(t *testing.T) {
	t.Skip("SFTP tests suspended - requires SFTP server for integration testing")

	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "Append to text file",
			input:   `let conn = SFTP("sftp://user:pass@example.com/"); "New line\n" =/=>> conn(@/log.txt).text`,
			wantErr: true,
		},
		{
			name:    "Append line",
			input:   `let conn = SFTP("sftp://user:pass@example.com/"); "Log entry" =/=>> conn(@/app.log).lines`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := testEvalHelper(tt.input)
			t.Logf("%s result type: %s", tt.name, result.Type())
		})
	}
}

// TestSFTPDirectoryOperations tests directory methods
func TestSFTPDirectoryOperations(t *testing.T) {
	t.Skip("SFTP tests suspended - requires SFTP server for integration testing")

	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "List directory",
			input:   `let conn = SFTP("sftp://user:pass@example.com/"); files <=/= conn(@/uploads).dir`,
			wantErr: true,
		},
		{
			name:    "Create directory",
			input:   `let conn = SFTP("sftp://user:pass@example.com/"); conn(@/newdir).mkdir()`,
			wantErr: true,
		},
		{
			name:    "Create directory with mode",
			input:   `let conn = SFTP("sftp://user:pass@example.com/"); conn(@/newdir).mkdir({mode: 0755})`,
			wantErr: true,
		},
		{
			name:    "Remove directory",
			input:   `let conn = SFTP("sftp://user:pass@example.com/"); conn(@/olddir).rmdir()`,
			wantErr: true,
		},
		{
			name:    "Remove file",
			input:   `let conn = SFTP("sftp://user:pass@example.com/"); conn(@/old.txt).remove()`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := testEvalHelper(tt.input)
			t.Logf("%s result type: %s", tt.name, result.Type())
		})
	}
}

// TestSFTPConnectionMethods tests connection lifecycle methods
func TestSFTPConnectionMethods(t *testing.T) {
	t.Skip("SFTP tests suspended - requires SFTP server for integration testing")

	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "Close connection",
			input:   `let conn = SFTP("sftp://user:pass@example.com/"); conn.close()`,
			wantErr: true, // Will fail to connect but should parse
		},
		{
			name:    "Use after close",
			input:   `let conn = SFTP("sftp://user:pass@example.com/"); conn.close(); conn(@/file.txt)`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := testEvalHelper(tt.input)
			t.Logf("%s result type: %s", tt.name, result.Type())
		})
	}
}

// TestSFTPErrorCapturePattern tests {data, error} pattern
func TestSFTPErrorCapturePattern(t *testing.T) {
	t.Skip("SFTP tests suspended - requires SFTP server for integration testing")

	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "Capture read error",
			input:   `let conn = SFTP("sftp://user:pass@example.com/"); {data, error} <=/= conn(@/missing.txt).text; error != null`,
			wantErr: false, // Should return true (error exists)
		},
		{
			name:    "Capture write error",
			input:   `let conn = SFTP("sftp://user:pass@example.com/"); {data, error} = "test" =/=> conn(@/readonly.txt).text; error != null`,
			wantErr: false, // Should return true (error exists)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := testEvalHelper(tt.input)

			if tt.wantErr {
				if _, isErr := result.(*evaluator.Error); !isErr {
					t.Errorf("Expected error, got %s", result.Type())
				}
			} else {
				// Check if result is boolean true (error exists)
				if boolean, ok := result.(*evaluator.Boolean); ok {
					t.Logf("Error exists: %v", boolean.Value)
				}
			}
		})
	}
}

// TestSFTPConnectionCaching tests connection pooling behavior
func TestSFTPConnectionCaching(t *testing.T) {
	t.Skip("SFTP tests suspended - requires SFTP server for integration testing")

	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "Reuse connection",
			input:   `let conn1 = SFTP("sftp://user:pass@example.com/"); let conn2 = SFTP("sftp://user:pass@example.com/"); conn1 == conn2`,
			wantErr: true, // Will fail to connect, but tests caching logic
		},
		{
			name:    "Different hosts create different connections",
			input:   `let conn1 = SFTP("sftp://user:pass@host1.com/"); let conn2 = SFTP("sftp://user:pass@host2.com/"); conn1 == conn2`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := testEvalHelper(tt.input)
			t.Logf("%s result type: %s", tt.name, result.Type())
		})
	}
}

// TestSFTPURLParsing tests SFTP URL parsing edge cases
func TestSFTPURLParsing(t *testing.T) {
	t.Skip("SFTP tests suspended - requires SFTP server for integration testing")

	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "URL with path",
			input:   `SFTP("sftp://user:pass@example.com/home/user/")`,
			wantErr: true,
		},
		{
			name:    "URL with special characters in password",
			input:   `SFTP("sftp://user:p@ss%20word@example.com/")`,
			wantErr: true,
		},
		{
			name:    "URL without trailing slash",
			input:   `SFTP("sftp://user:pass@example.com")`,
			wantErr: true,
		},
		{
			name:    "Invalid scheme",
			input:   `SFTP("http://user:pass@example.com/")`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := testEvalHelper(tt.input)

			if tt.wantErr {
				t.Logf("%s result type: %s", tt.name, result.Type())
			}
		})
	}
}

// TestSFTPFormatEncoding tests format-specific encoding
func TestSFTPFormatEncoding(t *testing.T) {
	t.Skip("SFTP tests suspended - requires SFTP server for integration testing")

	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "JSON encoding",
			input:   `let conn = SFTP("sftp://user:pass@example.com/"); {a: 1, b: [2, 3]} =/=> conn(@/test.json).json`,
			wantErr: true,
		},
		{
			name:    "Text encoding",
			input:   `let conn = SFTP("sftp://user:pass@example.com/"); "Hello\nWorld" =/=> conn(@/test.txt).text`,
			wantErr: true,
		},
		{
			name:    "Lines encoding",
			input:   `let conn = SFTP("sftp://user:pass@example.com/"); ["line1", "line2", "line3"] =/=> conn(@/test.txt).lines`,
			wantErr: true,
		},
		{
			name:    "Bytes encoding",
			input:   `let conn = SFTP("sftp://user:pass@example.com/"); [72, 101, 108, 108, 111] =/=> conn(@/test.bin).bytes`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := testEvalHelper(tt.input)
			t.Logf("%s result type: %s", tt.name, result.Type())
		})
	}
}
