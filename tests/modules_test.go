package main

import (
	"strings"
	"testing"

	"github.com/sambeau/parsley/pkg/evaluator"
	"github.com/sambeau/parsley/pkg/lexer"
	"github.com/sambeau/parsley/pkg/parser"
)

func evalModule(input string, filename string) evaluator.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	env := evaluator.NewEnvironment()
	env.Filename = filename
	return evaluator.Eval(program, env)
}

func TestBasicModuleImport(t *testing.T) {
	input := `
		let mod = import(@./test_fixtures/modules/simple.pars)
		mod.value
	`

	result := evalModule(input, "/Users/samphillips/Dev/parsley/test.pars")

	if result.Type() == evaluator.ERROR_OBJ {
		t.Fatalf("evaluation error: %s", result.Inspect())
	}

	integer, ok := result.(*evaluator.Integer)
	if !ok {
		t.Fatalf("expected Integer, got %T", result)
	}

	if integer.Value != 42 {
		t.Errorf("expected 42, got %d", integer.Value)
	}
}

func TestModuleDestructuring(t *testing.T) {
	input := `
		let {add, PI} = import(@./test_fixtures/modules/math.pars)
		add(10, 32)
	`

	result := evalModule(input, "/Users/samphillips/Dev/parsley/test.pars")

	if result.Type() == evaluator.ERROR_OBJ {
		t.Fatalf("evaluation error: %s", result.Inspect())
	}

	integer, ok := result.(*evaluator.Integer)
	if !ok {
		t.Fatalf("expected Integer, got %T", result)
	}

	if integer.Value != 42 {
		t.Errorf("expected 42, got %d", integer.Value)
	}
}

func TestModuleAlias(t *testing.T) {
	input := `
		let {square as sq} = import(@./test_fixtures/modules/math.pars)
		sq(8)
	`

	result := evalModule(input, "/Users/samphillips/Dev/parsley/test.pars")

	if result.Type() == evaluator.ERROR_OBJ {
		t.Fatalf("evaluation error: %s", result.Inspect())
	}

	integer, ok := result.(*evaluator.Integer)
	if !ok {
		t.Fatalf("expected Integer, got %T", result)
	}

	if integer.Value != 64 {
		t.Errorf("expected 64, got %d", integer.Value)
	}
}

func TestModuleCaching(t *testing.T) {
	input := `
		let mod1 = import(@./test_fixtures/modules/simple.pars)
		let mod2 = import(@./test_fixtures/modules/simple.pars)
		mod1 == mod2
	`

	result := evalModule(input, "/Users/samphillips/Dev/parsley/test.pars")

	if result.Type() == evaluator.ERROR_OBJ {
		t.Fatalf("evaluation error: %s", result.Inspect())
	}

	boolean, ok := result.(*evaluator.Boolean)
	if !ok {
		t.Fatalf("expected Boolean, got %T", result)
	}

	if !boolean.Value {
		t.Errorf("expected modules to be cached and equal")
	}
}

func TestModuleStringPath(t *testing.T) {
	input := `
		let mod = import("./test_fixtures/modules/simple.pars")
		mod.text
	`

	result := evalModule(input, "/Users/samphillips/Dev/parsley/test.pars")

	if result.Type() == evaluator.ERROR_OBJ {
		t.Fatalf("evaluation error: %s", result.Inspect())
	}

	str, ok := result.(*evaluator.String)
	if !ok {
		t.Fatalf("expected String, got %T", result)
	}

	if str.Value != "hello" {
		t.Errorf("expected 'hello', got %s", str.Value)
	}
}

func TestModuleNotFound(t *testing.T) {
	input := `
		let mod = import(@./nonexistent.pars)
	`

	result := evalModule(input, "/Users/samphillips/Dev/parsley/test.pars")

	if result.Type() != evaluator.ERROR_OBJ {
		t.Fatalf("expected error for missing module, got %T", result)
	}

	errStr := result.Inspect()
	if !strings.Contains(errStr, "failed to read module file") {
		t.Errorf("expected file not found error, got: %s", errStr)
	}
}

func TestModuleFunction(t *testing.T) {
	input := `
		let math = import(@./test_fixtures/modules/math.pars)
		math.add(math.multiply(2, 3), 4)
	`

	result := evalModule(input, "/Users/samphillips/Dev/parsley/test.pars")

	if result.Type() == evaluator.ERROR_OBJ {
		t.Fatalf("evaluation error: %s", result.Inspect())
	}

	integer, ok := result.(*evaluator.Integer)
	if !ok {
		t.Fatalf("expected Integer, got %T", result)
	}

	// (2 * 3) + 4 = 10
	if integer.Value != 10 {
		t.Errorf("expected 10, got %d", integer.Value)
	}
}

func TestModuleClosures(t *testing.T) {
	input := `
		let {double} = import(@./test_fixtures/modules/simple.pars)
		double(21)
	`

	result := evalModule(input, "/Users/samphillips/Dev/parsley/test.pars")

	if result.Type() == evaluator.ERROR_OBJ {
		t.Fatalf("evaluation error: %s", result.Inspect())
	}

	integer, ok := result.(*evaluator.Integer)
	if !ok {
		t.Fatalf("expected Integer, got %T", result)
	}

	if integer.Value != 42 {
		t.Errorf("expected 42, got %d", integer.Value)
	}
}
