package main

import (
	"testing"

	"github.com/sambeau/parsley/pkg/evaluator"
	"github.com/sambeau/parsley/pkg/lexer"
	"github.com/sambeau/parsley/pkg/parser"
)

func TestSQLiteConnection(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected interface{}
	}{
		{
			name:     "Create SQLite connection",
			input:    `let db = SQLITE(":memory:")`,
			expected: "DB_CONNECTION",
		},
		{
			name:     "Check connection type",
			input:    `let db = SQLITE(":memory:"); db`,
			expected: "<DBConnection driver=sqlite>",
		},
		{
			name:     "Ping connection",
			input:    `let db = SQLITE(":memory:"); db.ping()`,
			expected: true,
		},
		{
			name:     "Begin transaction",
			input:    `let db = SQLITE(":memory:"); db.begin()`,
			expected: true,
		},
		{
			name:     "Begin and commit transaction",
			input:    `let db = SQLITE(":memory:"); db.begin(); db.commit()`,
			expected: true,
		},
		{
			name:     "Begin and rollback transaction",
			input:    `let db = SQLITE(":memory:"); db.begin(); db.rollback()`,
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := parser.New(l)
			program := p.ParseProgram()

			if len(p.Errors()) != 0 {
				t.Fatalf("Parser errors: %v", p.Errors())
			}

			env := evaluator.NewEnvironment()
			result := evaluator.Eval(program, env)

			if result == nil {
				t.Fatalf("Eval returned nil")
			}

			switch expected := tt.expected.(type) {
			case string:
				if expected == "DB_CONNECTION" {
					if result.Type() != "DB_CONNECTION" {
						t.Errorf("Expected DB_CONNECTION object, got %s", result.Type())
					}
				} else {
					if result.Inspect() != expected {
						t.Errorf("Expected %q, got %q", expected, result.Inspect())
					}
				}
			case bool:
				boolean, ok := result.(*evaluator.Boolean)
				if !ok {
					t.Errorf("Expected Boolean object, got %T", result)
				} else if boolean.Value != expected {
					t.Errorf("Expected %v, got %v", expected, boolean.Value)
				}
			}
		})
	}
}
