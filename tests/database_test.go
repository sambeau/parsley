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

func TestSQLiteQueries(t *testing.T) {
	tests := []struct {
		name  string
		input string
		check func(*testing.T, evaluator.Object)
	}{
		{
			name: "Execute CREATE TABLE",
			input: `
				let db = SQLITE(":memory:")
				let result = db <=!=> "CREATE TABLE users (id INTEGER PRIMARY KEY, name TEXT)"
				result
			`,
			check: func(t *testing.T, result evaluator.Object) {
				dict, ok := result.(*evaluator.Dictionary)
				if !ok {
					t.Fatalf("Expected Dictionary, got %T", result)
				}
				// Check that affected exists
				if _, hasAffected := dict.Pairs["affected"]; !hasAffected {
					t.Error("Result should have 'affected' property")
				}
			},
		},
		{
			name: "Execute INSERT",
			input: `
				let db = SQLITE(":memory:")
				let _ = db <=!=> "DROP TABLE IF EXISTS test_users"
				let _ = db <=!=> "CREATE TABLE test_users (id INTEGER PRIMARY KEY, name TEXT)"
				let result = db <=!=> "INSERT INTO test_users (name) VALUES ('Alice')"
				result
			`,
			check: func(t *testing.T, result evaluator.Object) {
				dict, ok := result.(*evaluator.Dictionary)
				if !ok {
					t.Fatalf("Expected Dictionary, got %T", result)
				}
				// Check for affected rows
				affectedExpr, ok := dict.Pairs["affected"]
				if !ok {
					t.Fatal("Result should have 'affected' property")
				}
				affected := evaluator.Eval(affectedExpr, dict.Env)
				affectedInt, ok := affected.(*evaluator.Integer)
				if !ok || affectedInt.Value != 1 {
					t.Errorf("Expected affected=1, got %v", affected.Inspect())
				}
			},
		},
		{
			name: "Query single row with <=?=>",
			input: `
				let db = SQLITE(":memory:")
				let _ = db <=!=> "DROP TABLE IF EXISTS query_users"
				let _ = db <=!=> "CREATE TABLE query_users (id INTEGER PRIMARY KEY, name TEXT)"
				let _ = db <=!=> "INSERT INTO query_users (name) VALUES ('Alice')"
				let user = db <=?=> "SELECT * FROM query_users WHERE name = 'Alice'"
				user
			`,
			check: func(t *testing.T, result evaluator.Object) {
				dict, ok := result.(*evaluator.Dictionary)
				if !ok {
					t.Fatalf("Expected Dictionary, got %T", result)
				}
				// Check for name field
				if _, hasName := dict.Pairs["name"]; !hasName {
					t.Error("Result should have 'name' field")
				}
			},
		},
		{
			name: "Query multiple rows with <=??=>",
			input: `
				let db = SQLITE(":memory:")
				let _ = db <=!=> "DROP TABLE IF EXISTS many_users"
				let _ = db <=!=> "CREATE TABLE many_users (id INTEGER PRIMARY KEY, name TEXT)"
				let _ = db <=!=> "INSERT INTO many_users (name) VALUES ('Alice')"
				let _ = db <=!=> "INSERT INTO many_users (name) VALUES ('Bob')"
				let users = db <=??=> "SELECT * FROM many_users"
				users
			`,
			check: func(t *testing.T, result evaluator.Object) {
				arr, ok := result.(*evaluator.Array)
				if !ok {
					t.Fatalf("Expected Array, got %T", result)
				}
				if len(arr.Elements) != 2 {
					t.Errorf("Expected 2 users, got %d", len(arr.Elements))
				}
			},
		},
		{
			name: "Query non-existent row returns null",
			input: `
				let db = SQLITE(":memory:")
				let _ = db <=!=> "DROP TABLE IF EXISTS empty_users"
				let _ = db <=!=> "CREATE TABLE empty_users (id INTEGER PRIMARY KEY, name TEXT)"
				let user = db <=?=> "SELECT * FROM empty_users WHERE id = 999"
				user
			`,
			check: func(t *testing.T, result evaluator.Object) {
				if result.Type() != "NULL" {
					t.Errorf("Expected NULL, got %s", result.Type())
				}
			},
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

			if err, ok := result.(*evaluator.Error); ok {
				t.Fatalf("Eval returned error: %s", err.Message)
			}

			tt.check(t, result)
		})
	}
}

func TestSQLTag(t *testing.T) {
	tests := []struct {
		name  string
		input string
		check func(*testing.T, evaluator.Object)
	}{
		{
			name: "SQL tag with params",
			input: `
				let db = SQLITE(":memory:")
				let _ = db <=!=> "DROP TABLE IF EXISTS tag_users"
				let _ = db <=!=> "CREATE TABLE tag_users (id INTEGER PRIMARY KEY, name TEXT)"
				
				let InsertUser = fn(props) {
					<SQL>
						INSERT INTO tag_users (name) VALUES ('Alice')
					</SQL>
				}
				
				let result = db <=!=> <InsertUser />
				result
			`,
			check: func(t *testing.T, result evaluator.Object) {
				dict, ok := result.(*evaluator.Dictionary)
				if !ok {
					t.Fatalf("Expected Dictionary, got %T", result)
				}
				affectedExpr, ok := dict.Pairs["affected"]
				if !ok {
					t.Fatal("Result should have 'affected' property")
				}
				affected := evaluator.Eval(affectedExpr, dict.Env)
				affectedInt, ok := affected.(*evaluator.Integer)
				if !ok || affectedInt.Value != 1 {
					t.Errorf("Expected affected=1, got %v", affected.Inspect())
				}
			},
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

			if err, ok := result.(*evaluator.Error); ok {
				t.Fatalf("Eval returned error: %s", err.Message)
			}

			tt.check(t, result)
		})
	}
}
