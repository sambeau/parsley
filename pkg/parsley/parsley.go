// Package parsley provides a public API for embedding the Parsley language interpreter.
//
// Basic usage:
//
//	result, err := parsley.Eval(`1 + 2`)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(result.String()) // "3"
//
// With variables:
//
//	result, err := parsley.Eval(`name ++ "!"`,
//	    parsley.WithVar("name", "Hello"),
//	)
//
// With a file:
//
//	result, err := parsley.EvalFile("script.pars",
//	    parsley.WithSecurity(policy),
//	)
package parsley

import (
	"fmt"
	"os"

	"github.com/sambeau/parsley/pkg/evaluator"
	"github.com/sambeau/parsley/pkg/lexer"
	"github.com/sambeau/parsley/pkg/parser"
)

// Result wraps evaluation output
type Result struct {
	// Value is the Parsley object result
	Value evaluator.Object
	// Err contains any runtime error (also returned by Eval functions)
	Err error
}

// GoValue converts the result to a Go value using FromParsley
func (r *Result) GoValue() interface{} {
	if r.Value == nil {
		return nil
	}
	return FromParsley(r.Value)
}

// String returns the string representation of the result
func (r *Result) String() string {
	if r.Value == nil {
		return ""
	}
	return evaluator.ObjectToPrintString(r.Value)
}

// IsNull returns true if the result is null
func (r *Result) IsNull() bool {
	if r.Value == nil {
		return true
	}
	_, ok := r.Value.(*evaluator.Null)
	return ok
}

// IsError returns true if the result is an error
func (r *Result) IsError() bool {
	if r.Value == nil {
		return false
	}
	return r.Value.Type() == evaluator.ERROR_OBJ
}

// Eval evaluates Parsley source code and returns the result.
func Eval(source string, opts ...Option) (*Result, error) {
	config := newConfig(opts...)

	// Create or use provided environment
	env := config.Env
	if env == nil {
		env = evaluator.NewEnvironment()
	}

	// Apply configuration
	if err := applyConfig(env, config); err != nil {
		return nil, fmt.Errorf("configuration error: %w", err)
	}

	// Create lexer and parser
	var l *lexer.Lexer
	if config.Filename != "" {
		l = lexer.NewWithFilename(source, config.Filename)
		env.Filename = config.Filename
	} else {
		l = lexer.New(source)
	}

	p := parser.New(l)
	program := p.ParseProgram()

	// Check for parse errors
	if errors := p.Errors(); len(errors) != 0 {
		return nil, &ParseError{Errors: errors}
	}

	// Evaluate the program
	result := evaluator.Eval(program, env)

	// Check for runtime errors
	if result != nil && result.Type() == evaluator.ERROR_OBJ {
		errObj := result.(*evaluator.Error)
		return &Result{Value: result, Err: &RuntimeError{Message: errObj.Message, Line: errObj.Line, Column: errObj.Column}},
			&RuntimeError{Message: errObj.Message, Line: errObj.Line, Column: errObj.Column}
	}

	return &Result{Value: result}, nil
}

// EvalFile evaluates a Parsley file and returns the result.
func EvalFile(filename string, opts ...Option) (*Result, error) {
	// Read the file
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("reading file: %w", err)
	}

	// Add filename to options
	opts = append([]Option{WithFilename(filename)}, opts...)

	return Eval(string(content), opts...)
}

// ParseError represents one or more parse errors
type ParseError struct {
	Errors []string
}

func (e *ParseError) Error() string {
	if len(e.Errors) == 1 {
		return e.Errors[0]
	}
	return fmt.Sprintf("%d parse errors: %s", len(e.Errors), e.Errors[0])
}

// RuntimeError represents a runtime error during evaluation
type RuntimeError struct {
	Message string
	Line    int
	Column  int
}

func (e *RuntimeError) Error() string {
	if e.Line > 0 {
		return fmt.Sprintf("line %d, column %d: %s", e.Line, e.Column, e.Message)
	}
	return e.Message
}

// Re-export types from evaluator for convenience
type (
	// Object is the interface for all Parsley values
	Object = evaluator.Object

	// Environment holds variable bindings
	Environment = evaluator.Environment

	// SecurityPolicy controls file system access
	SecurityPolicy = evaluator.SecurityPolicy

	// Integer represents integer values
	Integer = evaluator.Integer

	// Float represents floating-point values
	Float = evaluator.Float

	// String represents string values
	String = evaluator.String

	// Boolean represents boolean values
	Boolean = evaluator.Boolean

	// Array represents array values
	Array = evaluator.Array

	// Dictionary represents dictionary values
	Dictionary = evaluator.Dictionary

	// Null represents the null value
	Null = evaluator.Null
)

// Re-export constants
const (
	INTEGER_OBJ    = evaluator.INTEGER_OBJ
	FLOAT_OBJ      = evaluator.FLOAT_OBJ
	BOOLEAN_OBJ    = evaluator.BOOLEAN_OBJ
	STRING_OBJ     = evaluator.STRING_OBJ
	NULL_OBJ       = evaluator.NULL_OBJ
	ARRAY_OBJ      = evaluator.ARRAY_OBJ
	DICTIONARY_OBJ = evaluator.DICTIONARY_OBJ
	ERROR_OBJ      = evaluator.ERROR_OBJ
)

// Re-export functions
var (
	// NewEnvironment creates a fresh evaluation environment
	NewEnvironment = evaluator.NewEnvironment

	// NewEnclosedEnvironment creates a new environment with outer reference
	NewEnclosedEnvironment = evaluator.NewEnclosedEnvironment

	// NewDictionaryFromObjects creates a Dictionary from a map of Objects
	NewDictionaryFromObjects = evaluator.NewDictionaryFromObjects

	// ObjectToPrintString converts an object to its print representation
	ObjectToPrintString = evaluator.ObjectToPrintString
)

// Re-export singleton values
var (
	// NULL is the null value
	NULL = evaluator.NULL

	// TRUE is the boolean true value
	TRUE = evaluator.TRUE

	// FALSE is the boolean false value
	FALSE = evaluator.FALSE
)
