package parsley_test

import (
	"database/sql"
	"testing"

	"github.com/sambeau/parsley/pkg/parsley"
	_ "modernc.org/sqlite"
)

func TestEval(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected interface{}
	}{
		{"integer", "42", int64(42)},
		{"float", "3.14", 3.14},
		{"string", `"hello"`, "hello"},
		{"boolean true", "true", true},
		{"boolean false", "false", false},
		{"addition", "1 + 2", int64(3)},
		{"string concat", `"a" ++ "b"`, []interface{}{"a", "b"}}, // ++ creates array
		{"array", "[1, 2, 3]", []interface{}{int64(1), int64(2), int64(3)}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parsley.Eval(tt.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			got := result.GoValue()

			// Compare arrays specially
			if expectedArr, ok := tt.expected.([]interface{}); ok {
				gotArr, ok := got.([]interface{})
				if !ok {
					t.Fatalf("expected array, got %T", got)
				}
				if len(gotArr) != len(expectedArr) {
					t.Fatalf("expected array length %d, got %d", len(expectedArr), len(gotArr))
				}
				for i := range expectedArr {
					if gotArr[i] != expectedArr[i] {
						t.Errorf("array[%d]: expected %v, got %v", i, expectedArr[i], gotArr[i])
					}
				}
			} else if got != tt.expected {
				t.Errorf("expected %v (%T), got %v (%T)", tt.expected, tt.expected, got, got)
			}
		})
	}
}

func TestEvalWithVar(t *testing.T) {
	result, err := parsley.Eval(`name ++ "!"`,
		parsley.WithVar("name", "Hello"),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.String() != "Hello!" {
		t.Errorf("expected 'Hello!', got '%s'", result.String())
	}
}

func TestEvalWithMultipleVars(t *testing.T) {
	result, err := parsley.Eval(`a + b`,
		parsley.WithVar("a", 10),
		parsley.WithVar("b", 20),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.GoValue().(int64) != 30 {
		t.Errorf("expected 30, got %v", result.GoValue())
	}
}

func TestEvalWithLogger(t *testing.T) {
	logger := parsley.NewBufferedLogger()

	result, err := parsley.Eval(`log("hello"); 42`,
		parsley.WithLogger(logger),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.GoValue().(int64) != 42 {
		t.Errorf("expected 42, got %v", result.GoValue())
	}

	logged := logger.String()
	if logged != "hello\n" {
		t.Errorf("expected 'hello\\n', got '%s'", logged)
	}
}

func TestParseError(t *testing.T) {
	_, err := parsley.Eval(`let x = `)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	_, ok := err.(*parsley.ParseError)
	if !ok {
		t.Errorf("expected ParseError, got %T", err)
	}
}

func TestRuntimeError(t *testing.T) {
	// Test an actual runtime error - undefined variable
	_, err := parsley.Eval(`undefinedVariable`)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	_, ok := err.(*parsley.RuntimeError)
	if !ok {
		t.Errorf("expected RuntimeError, got %T: %v", err, err)
	}
}

func TestToParsley(t *testing.T) {
	tests := []struct {
		name  string
		input interface{}
		check func(parsley.Object) bool
	}{
		{"int", 42, func(o parsley.Object) bool {
			i, ok := o.(*parsley.Integer)
			return ok && i.Value == 42
		}},
		{"float", 3.14, func(o parsley.Object) bool {
			f, ok := o.(*parsley.Float)
			return ok && f.Value == 3.14
		}},
		{"string", "hello", func(o parsley.Object) bool {
			s, ok := o.(*parsley.String)
			return ok && s.Value == "hello"
		}},
		{"bool true", true, func(o parsley.Object) bool {
			b, ok := o.(*parsley.Boolean)
			return ok && b.Value == true
		}},
		{"nil", nil, func(o parsley.Object) bool {
			_, ok := o.(*parsley.Null)
			return ok
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obj, err := parsley.ToParsley(tt.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !tt.check(obj) {
				t.Errorf("conversion failed for %v, got %v", tt.input, obj)
			}
		})
	}
}

func TestFromParsley(t *testing.T) {
	tests := []struct {
		name     string
		obj      parsley.Object
		expected interface{}
	}{
		{"integer", &parsley.Integer{Value: 42}, int64(42)},
		{"float", &parsley.Float{Value: 3.14}, 3.14},
		{"string", &parsley.String{Value: "hello"}, "hello"},
		{"boolean", parsley.TRUE, true},
		{"null", parsley.NULL, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parsley.FromParsley(tt.obj)
			if got != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, got)
			}
		})
	}
}

func TestWithDB(t *testing.T) {
	// Open an in-memory SQLite database
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}
	defer db.Close()

	// Create a test table
	_, err = db.Exec("CREATE TABLE users (id INTEGER PRIMARY KEY, name TEXT)")
	if err != nil {
		t.Fatalf("failed to create table: %v", err)
	}

	// Insert test data
	_, err = db.Exec("INSERT INTO users (name) VALUES ('Alice')")
	if err != nil {
		t.Fatalf("failed to insert data: %v", err)
	}

	// Use WithDB to inject the connection
	result, err := parsley.Eval(`
		let user = db <=?=> "SELECT * FROM users WHERE id = 1"
		user.name
	`, parsley.WithDB("db", db, "sqlite"))

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.String() != "Alice" {
		t.Errorf("expected 'Alice', got '%s'", result.String())
	}
}

func TestWithDBManagedConnectionCannotBeClosed(t *testing.T) {
	// Open an in-memory SQLite database
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}
	defer db.Close()

	// Try to close the managed connection from Parsley
	result, err := parsley.Eval(`db.close()`, parsley.WithDB("db", db, "sqlite"))

	// Should return an error about managed connections
	if err == nil && !result.IsError() {
		t.Error("expected error when closing managed connection")
	}
}
