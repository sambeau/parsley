package evaluator

import (
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
	"unicode"

	"github.com/sambeau/parsley/pkg/ast"
	"github.com/sambeau/parsley/pkg/lexer"
	"github.com/sambeau/parsley/pkg/parser"
)

// ObjectType represents the type of objects in our language
type ObjectType string

const (
	INTEGER_OBJ    = "INTEGER"
	FLOAT_OBJ      = "FLOAT"
	BOOLEAN_OBJ    = "BOOLEAN"
	STRING_OBJ     = "STRING"
	NULL_OBJ       = "NULL"
	RETURN_OBJ     = "RETURN_VALUE"
	ERROR_OBJ      = "ERROR"
	FUNCTION_OBJ   = "FUNCTION"
	BUILTIN_OBJ    = "BUILTIN"
	ARRAY_OBJ      = "ARRAY"
	DICTIONARY_OBJ = "DICTIONARY"
)

// Object represents all values in our language
type Object interface {
	Type() ObjectType
	Inspect() string
}

// Integer represents integer objects
type Integer struct {
	Value int64
}

func (i *Integer) Inspect() string  { return strconv.FormatInt(i.Value, 10) }
func (i *Integer) Type() ObjectType { return INTEGER_OBJ }

// Float represents floating-point objects
type Float struct {
	Value float64
}

func (f *Float) Inspect() string  { return fmt.Sprintf("%g", f.Value) }
func (f *Float) Type() ObjectType { return FLOAT_OBJ }

// Boolean represents boolean objects
type Boolean struct {
	Value bool
}

func (b *Boolean) Inspect() string  { return strconv.FormatBool(b.Value) }
func (b *Boolean) Type() ObjectType { return BOOLEAN_OBJ }

// String represents string objects
type String struct {
	Value string
}

func (s *String) Inspect() string  { return s.Value }
func (s *String) Type() ObjectType { return STRING_OBJ }

// Null represents null/nil objects
type Null struct{}

func (n *Null) Inspect() string  { return "null" }
func (n *Null) Type() ObjectType { return NULL_OBJ }

// ReturnValue wraps other objects when returned
type ReturnValue struct {
	Value Object
}

func (rv *ReturnValue) Type() ObjectType { return RETURN_OBJ }
func (rv *ReturnValue) Inspect() string  { return rv.Value.Inspect() }

// Error represents error objects
type Error struct {
	Message string
}

func (e *Error) Type() ObjectType { return ERROR_OBJ }
func (e *Error) Inspect() string  { return "ERROR: " + e.Message }

// Function represents function objects
type Function struct {
	Parameters []*ast.Identifier        // deprecated - kept for compatibility
	Params     []*ast.FunctionParameter // new parameter list with destructuring support
	Body       *ast.BlockStatement
	Env        *Environment
}

func (f *Function) Type() ObjectType { return FUNCTION_OBJ }
func (f *Function) Inspect() string {
	if len(f.Params) > 0 {
		return fmt.Sprintf("fn(%v) {\n%s\n}", f.Params, f.Body.String())
	}
	return fmt.Sprintf("fn(%v) {\n%s\n}", f.Parameters, f.Body.String())
}

// ParamCount returns the number of parameters for this function
func (f *Function) ParamCount() int {
	if len(f.Params) > 0 {
		return len(f.Params)
	}
	return len(f.Parameters)
}

// BuiltinFunction represents a built-in function
type BuiltinFunction func(args ...Object) Object

// Builtin represents built-in function objects
type Builtin struct {
	Fn BuiltinFunction
}

func (b *Builtin) Type() ObjectType { return BUILTIN_OBJ }
func (b *Builtin) Inspect() string  { return "builtin function" }

// Array represents array objects
type Array struct {
	Elements []Object
}

func (a *Array) Type() ObjectType { return ARRAY_OBJ }
func (a *Array) Inspect() string {
	var out strings.Builder
	elements := []string{}
	for _, e := range a.Elements {
		if e != nil {
			elements = append(elements, e.Inspect())
		} else {
			elements = append(elements, "nil")
		}
	}
	out.WriteString(strings.Join(elements, ", "))
	return out.String()
}

// Dictionary represents dictionary objects with lazy evaluation
type Dictionary struct {
	Pairs map[string]ast.Expression // Store expressions for lazy evaluation
	Env   *Environment              // Environment for evaluation (for 'this' binding)
}

func (d *Dictionary) Type() ObjectType { return DICTIONARY_OBJ }
func (d *Dictionary) Inspect() string {
	var out strings.Builder
	pairs := []string{}
	for key, expr := range d.Pairs {
		// For inspection, we show the expression, not the evaluated value
		pairs = append(pairs, fmt.Sprintf("%s: %s", key, expr.String()))
	}
	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")
	return out.String()
}

// Environment represents the environment for variable bindings
type Environment struct {
	store     map[string]Object
	outer     *Environment
	Filename  string
	LastToken *lexer.Token
}

// NewEnvironment creates a new environment
func NewEnvironment() *Environment {
	s := make(map[string]Object)
	return &Environment{store: s, outer: nil}
}

// NewEnclosedEnvironment creates a new environment with outer reference
func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	// Preserve filename and token from outer environment
	if outer != nil {
		env.Filename = outer.Filename
		env.LastToken = outer.LastToken
	}
	return env
}

// Get retrieves a value from the environment
func (e *Environment) Get(name string) (Object, bool) {
	value, ok := e.store[name]
	if !ok && e.outer != nil {
		value, ok = e.outer.Get(name)
	}
	return value, ok
}

// Set stores a value in the environment
func (e *Environment) Set(name string, val Object) Object {
	e.store[name] = val
	return val
}

// Update stores a value in the environment where it's defined (current or outer)
// If the variable doesn't exist anywhere, it creates it in the current scope
func (e *Environment) Update(name string, val Object) Object {
	// Check if variable exists in current scope
	if _, ok := e.store[name]; ok {
		e.store[name] = val
		return val
	}

	// Check if it exists in outer scope
	if e.outer != nil {
		if _, ok := e.outer.Get(name); ok {
			return e.outer.Update(name, val)
		}
	}

	// Variable doesn't exist anywhere, create it in current scope
	e.store[name] = val
	return val
}

// Global constants
var (
	NULL  = &Null{}
	TRUE  = &Boolean{Value: true}
	FALSE = &Boolean{Value: false}
)

// naturalCompare compares two objects using natural sort order
// Returns true if a < b in natural sort order
func naturalCompare(a, b Object) bool {
	// Type-based ordering: numbers < strings
	aType := getTypeOrder(a)
	bType := getTypeOrder(b)

	if aType != bType {
		return aType < bType
	}

	// Both are numbers
	if aType == 0 {
		return compareNumbers(a, b)
	}

	// Both are strings - use natural string comparison
	if aType == 1 {
		aStr := a.(*String).Value
		bStr := b.(*String).Value
		return naturalStringCompare(aStr, bStr)
	}

	// Other types (shouldn't happen with current implementation)
	return false
}

// getTypeOrder returns a sort order for types
// 0 = numbers (Integer, Float)
// 1 = strings
// 2 = other
func getTypeOrder(obj Object) int {
	switch obj.Type() {
	case INTEGER_OBJ, FLOAT_OBJ:
		return 0
	case STRING_OBJ:
		return 1
	default:
		return 2
	}
}

// compareNumbers compares two numeric objects
func compareNumbers(a, b Object) bool {
	aVal := getNumericValue(a)
	bVal := getNumericValue(b)
	return aVal < bVal
}

// getNumericValue extracts numeric value as float64
func getNumericValue(obj Object) float64 {
	switch obj := obj.(type) {
	case *Integer:
		return float64(obj.Value)
	case *Float:
		return obj.Value
	default:
		return 0
	}
}

// naturalStringCompare compares strings using natural sort order
// It treats consecutive digits as numbers and compares them numerically
func naturalStringCompare(a, b string) bool {
	aRunes := []rune(a)
	bRunes := []rune(b)

	i, j := 0, 0

	for i < len(aRunes) && j < len(bRunes) {
		aChar := aRunes[i]
		bChar := bRunes[j]

		// Both are digits - compare numerically
		if unicode.IsDigit(aChar) && unicode.IsDigit(bChar) {
			// Extract the full number from both strings
			aNum, aEnd := extractNumber(aRunes, i)
			bNum, bEnd := extractNumber(bRunes, j)

			if aNum != bNum {
				return aNum < bNum
			}

			i = aEnd
			j = bEnd
			continue
		}

		// Character comparison
		if aChar != bChar {
			return aChar < bChar
		}

		i++
		j++
	}

	// If we've exhausted one string, the shorter one comes first
	return len(aRunes) < len(bRunes)
}

// extractNumber extracts a number from a rune slice starting at the given position
// Returns the number and the position after the last digit
func extractNumber(runes []rune, start int) (int64, int) {
	end := start
	for end < len(runes) && unicode.IsDigit(runes[end]) {
		end++
	}

	numStr := string(runes[start:end])
	num, err := strconv.ParseInt(numStr, 10, 64)
	if err != nil {
		return 0, end
	}

	return num, end
}

// objectsEqual compares two objects for equality
func objectsEqual(a, b Object) bool {
	if a.Type() != b.Type() {
		return false
	}

	switch a := a.(type) {
	case *Integer:
		return a.Value == b.(*Integer).Value
	case *Float:
		return a.Value == b.(*Float).Value
	case *String:
		return a.Value == b.(*String).Value
	case *Boolean:
		return a.Value == b.(*Boolean).Value
	case *Null:
		return true
	default:
		return false
	}
}

// getBuiltins returns the map of built-in functions
func getBuiltins() map[string]*Builtin {
	return map[string]*Builtin{
		"sin": {
			Fn: func(args ...Object) Object {
				if len(args) != 1 {
					return newError("wrong number of arguments. got=%d, want=1", len(args))
				}

				arg := args[0]
				switch arg := arg.(type) {
				case *Integer:
					return &Float{Value: math.Sin(float64(arg.Value))}
				case *Float:
					return &Float{Value: math.Sin(arg.Value)}
				default:
					return newError("argument to `sin` not supported, got %T", arg)
				}
			},
		},
		"cos": {
			Fn: func(args ...Object) Object {
				if len(args) != 1 {
					return newError("wrong number of arguments. got=%d, want=1", len(args))
				}

				arg := args[0]
				switch arg := arg.(type) {
				case *Integer:
					return &Float{Value: math.Cos(float64(arg.Value))}
				case *Float:
					return &Float{Value: math.Cos(arg.Value)}
				default:
					return newError("argument to `cos` not supported, got %T", arg)
				}
			},
		},
		"tan": {
			Fn: func(args ...Object) Object {
				if len(args) != 1 {
					return newError("wrong number of arguments. got=%d, want=1", len(args))
				}

				arg := args[0]
				switch arg := arg.(type) {
				case *Integer:
					return &Float{Value: math.Tan(float64(arg.Value))}
				case *Float:
					return &Float{Value: math.Tan(arg.Value)}
				default:
					return newError("argument to `tan` not supported, got %T", arg)
				}
			},
		},
		"asin": {
			Fn: func(args ...Object) Object {
				if len(args) != 1 {
					return newError("wrong number of arguments. got=%d, want=1", len(args))
				}

				arg := args[0]
				switch arg := arg.(type) {
				case *Integer:
					return &Float{Value: math.Asin(float64(arg.Value))}
				case *Float:
					return &Float{Value: math.Asin(arg.Value)}
				default:
					return newError("argument to `asin` not supported, got %T", arg)
				}
			},
		},
		"acos": {
			Fn: func(args ...Object) Object {
				if len(args) != 1 {
					return newError("wrong number of arguments. got=%d, want=1", len(args))
				}

				arg := args[0]
				switch arg := arg.(type) {
				case *Integer:
					return &Float{Value: math.Acos(float64(arg.Value))}
				case *Float:
					return &Float{Value: math.Acos(arg.Value)}
				default:
					return newError("argument to `acos` not supported, got %T", arg)
				}
			},
		},
		"atan": {
			Fn: func(args ...Object) Object {
				if len(args) != 1 {
					return newError("wrong number of arguments. got=%d, want=1", len(args))
				}

				arg := args[0]
				switch arg := arg.(type) {
				case *Integer:
					return &Float{Value: math.Atan(float64(arg.Value))}
				case *Float:
					return &Float{Value: math.Atan(arg.Value)}
				default:
					return newError("argument to `atan` not supported, got %T", arg)
				}
			},
		},
		"sqrt": {
			Fn: func(args ...Object) Object {
				if len(args) != 1 {
					return newError("wrong number of arguments. got=%d, want=1", len(args))
				}

				arg := args[0]
				switch arg := arg.(type) {
				case *Integer:
					return &Float{Value: math.Sqrt(float64(arg.Value))}
				case *Float:
					return &Float{Value: math.Sqrt(arg.Value)}
				default:
					return newError("argument to `sqrt` not supported, got %T", arg)
				}
			},
		},
		"round": {
			Fn: func(args ...Object) Object {
				if len(args) != 1 {
					return newError("wrong number of arguments. got=%d, want=1", len(args))
				}

				arg := args[0]
				switch arg := arg.(type) {
				case *Integer:
					return arg // already an integer
				case *Float:
					return &Integer{Value: int64(math.Round(arg.Value))}
				default:
					return newError("argument to `round` not supported, got %T", arg)
				}
			},
		},
		"pow": {
			Fn: func(args ...Object) Object {
				if len(args) != 2 {
					return newError("wrong number of arguments. got=%d, want=2", len(args))
				}

				base := args[0]
				exp := args[1]

				var baseVal, expVal float64

				switch base := base.(type) {
				case *Integer:
					baseVal = float64(base.Value)
				case *Float:
					baseVal = base.Value
				default:
					return newError("first argument to `pow` not supported, got %T", base)
				}

				switch exp := exp.(type) {
				case *Integer:
					expVal = float64(exp.Value)
				case *Float:
					expVal = exp.Value
				default:
					return newError("second argument to `pow` not supported, got %T", exp)
				}

				return &Float{Value: math.Pow(baseVal, expVal)}
			},
		},
		"pi": {
			Fn: func(args ...Object) Object {
				if len(args) != 0 {
					return newError("wrong number of arguments. got=%d, want=0", len(args))
				}
				return &Float{Value: math.Pi}
			},
		},
		"map": {
			Fn: func(args ...Object) Object {
				if len(args) < 2 {
					return newError("wrong number of arguments to `map`. got=%d, want at least 2", len(args))
				}

				fn, ok := args[0].(*Function)
				if !ok {
					return newError("first argument to `map` must be a function, got %s", args[0].Type())
				}

				// If second argument is an array, use it; otherwise create array from remaining args
				var arr *Array
				if a, ok := args[1].(*Array); ok && len(args) == 2 {
					arr = a
				} else {
					// Create array from all arguments after the function
					arr = &Array{Elements: args[1:]}
				}

				// Validate function parameter count
				if fn.ParamCount() != 1 {
					return newError("function passed to `map` must take exactly 1 parameter, got %d", fn.ParamCount())
				}

				result := []Object{}
				for _, elem := range arr.Elements {
					// Apply function to each element
					extendedEnv := extendFunctionEnv(fn, []Object{elem})

					// Evaluate the function body
					var evaluated Object
					for _, stmt := range fn.Body.Statements {
						evaluated = evalStatement(stmt, extendedEnv)
						if returnValue, ok := evaluated.(*ReturnValue); ok {
							evaluated = returnValue.Value
							break
						}
						if isError(evaluated) {
							return evaluated
						}
					}

					// Skip null values (filter behavior)
					if evaluated != NULL {
						result = append(result, evaluated)
					}
				}

				return &Array{Elements: result}
			},
		},
		"toUpper": {
			Fn: func(args ...Object) Object {
				if len(args) != 1 {
					return newError("wrong number of arguments to `toUpper`. got=%d, want=1", len(args))
				}

				str, ok := args[0].(*String)
				if !ok {
					return newError("argument to `toUpper` must be a string, got %s", args[0].Type())
				}

				return &String{Value: strings.ToUpper(str.Value)}
			},
		},
		"toLower": {
			Fn: func(args ...Object) Object {
				if len(args) != 1 {
					return newError("wrong number of arguments to `toLower`. got=%d, want=1", len(args))
				}

				str, ok := args[0].(*String)
				if !ok {
					return newError("argument to `toLower` must be a string, got %s", args[0].Type())
				}

				return &String{Value: strings.ToLower(str.Value)}
			},
		},
		"len": {
			Fn: func(args ...Object) Object {
				if len(args) != 1 {
					return newError("wrong number of arguments to `len`. got=%d, want=1", len(args))
				}

				switch arg := args[0].(type) {
				case *String:
					return &Integer{Value: int64(len(arg.Value))}
				case *Array:
					return &Integer{Value: int64(len(arg.Elements))}
				default:
					return newError("argument to `len` not supported, got %s", args[0].Type())
				}
			},
		},
		"toInt": {
			Fn: func(args ...Object) Object {
				if len(args) != 1 {
					return newError("wrong number of arguments to `toInt`. got=%d, want=1", len(args))
				}

				str, ok := args[0].(*String)
				if !ok {
					return newError("argument to `toInt` must be a string, got %s", args[0].Type())
				}

				var val int64
				_, err := fmt.Sscanf(str.Value, "%d", &val)
				if err != nil {
					return newError("cannot convert '%s' to integer", str.Value)
				}

				return &Integer{Value: val}
			},
		},
		"toFloat": {
			Fn: func(args ...Object) Object {
				if len(args) != 1 {
					return newError("wrong number of arguments to `toFloat`. got=%d, want=1", len(args))
				}

				str, ok := args[0].(*String)
				if !ok {
					return newError("argument to `toFloat` must be a string, got %s", args[0].Type())
				}

				var val float64
				_, err := fmt.Sscanf(str.Value, "%f", &val)
				if err != nil {
					return newError("cannot convert '%s' to float", str.Value)
				}

				return &Float{Value: val}
			},
		},
		"toNumber": {
			Fn: func(args ...Object) Object {
				if len(args) != 1 {
					return newError("wrong number of arguments to `toNumber`. got=%d, want=1", len(args))
				}

				str, ok := args[0].(*String)
				if !ok {
					return newError("argument to `toNumber` must be a string, got %s", args[0].Type())
				}

				// Try to parse as integer first
				var intVal int64
				if _, err := fmt.Sscanf(str.Value, "%d", &intVal); err == nil {
					// Check if the string has a decimal point - if so, it's a float
					if !strings.Contains(str.Value, ".") {
						return &Integer{Value: intVal}
					}
				}

				// Parse as float
				var floatVal float64
				if _, err := fmt.Sscanf(str.Value, "%f", &floatVal); err == nil {
					return &Float{Value: floatVal}
				}

				return newError("cannot convert '%s' to number", str.Value)
			},
		},
		"toString": {
			Fn: func(args ...Object) Object {
				var result strings.Builder

				for _, arg := range args {
					result.WriteString(objectToPrintString(arg))
				}

				return &String{Value: result.String()}
			},
		},
		"toDebug": {
			Fn: func(args ...Object) Object {
				var result strings.Builder

				for i, arg := range args {
					if i > 0 {
						result.WriteString(", ")
					}
					result.WriteString(objectToDebugString(arg))
				}

				return &String{Value: result.String()}
			},
		},
		"log": {
			Fn: func(args ...Object) Object {
				var result strings.Builder

				for i, arg := range args {
					if i == 0 {
						// First argument: if it's a string, show without quotes
						if str, ok := arg.(*String); ok {
							result.WriteString(str.Value)
						} else {
							result.WriteString(objectToDebugString(arg))
						}
					} else {
						// Subsequent arguments: add separator and debug format
						if i == 1 {
							// After first string, no comma - just space
							if _, firstWasString := args[0].(*String); firstWasString {
								result.WriteString(" ")
							} else {
								result.WriteString(", ")
							}
						} else {
							result.WriteString(", ")
						}
						result.WriteString(objectToDebugString(arg))
					}
				}

				// Write immediately to stdout
				fmt.Fprintln(os.Stdout, result.String())

				// Return null
				return NULL
			},
		},
		"logLine": {
			Fn: func(args ...Object) Object {
				// This is a placeholder - will be replaced with actual implementation
				// that has access to environment
				return NULL
			},
		},
		"sort": {
			Fn: func(args ...Object) Object {
				if len(args) != 1 {
					return newError("wrong number of arguments to `sort`. got=%d, want=1", len(args))
				}

				arr, ok := args[0].(*Array)
				if !ok {
					return newError("argument to `sort` must be an array, got %s", args[0].Type())
				}

				// Create a copy to avoid modifying the original
				sortedElements := make([]Object, len(arr.Elements))
				copy(sortedElements, arr.Elements)

				// Sort using natural sort comparison
				sort.Slice(sortedElements, func(i, j int) bool {
					return naturalCompare(sortedElements[i], sortedElements[j])
				})

				return &Array{Elements: sortedElements}
			},
		},
		"reverse": {
			Fn: func(args ...Object) Object {
				if len(args) != 1 {
					return newError("wrong number of arguments to `reverse`. got=%d, want=1", len(args))
				}

				arr, ok := args[0].(*Array)
				if !ok {
					return newError("argument to `reverse` must be an array, got %s", args[0].Type())
				}

				// Create a reversed copy
				reversed := make([]Object, len(arr.Elements))
				for i, elem := range arr.Elements {
					reversed[len(arr.Elements)-1-i] = elem
				}

				return &Array{Elements: reversed}
			},
		},
		"sortBy": {
			Fn: func(args ...Object) Object {
				if len(args) != 2 {
					return newError("wrong number of arguments to `sortBy`. got=%d, want=2", len(args))
				}

				arr, ok := args[0].(*Array)
				if !ok {
					return newError("first argument to `sortBy` must be an array, got %s", args[0].Type())
				}

				compareFn := args[1]

				// Verify it's a function
				fn, ok := compareFn.(*Function)
				if !ok {
					return newError("second argument to `sortBy` must be a function, got %s", compareFn.Type())
				}

				// Verify the function takes exactly 2 parameters
				if fn.ParamCount() != 2 {
					return newError("comparison function must take exactly 2 parameters, got %d", fn.ParamCount())
				}

				// Create a copy to avoid modifying the original
				sortedElements := make([]Object, len(arr.Elements))
				copy(sortedElements, arr.Elements)

				// Sort using the custom comparison function
				sort.Slice(sortedElements, func(i, j int) bool {
					// Call the comparison function with the two elements
					result := applyFunction(fn, []Object{sortedElements[i], sortedElements[j]})

					// The function should return a 2-element array
					resultArr, ok := result.(*Array)
					if !ok || len(resultArr.Elements) != 2 {
						return false
					}

					// Check if the first element equals sortedElements[i]
					// If so, it means i comes before j (ascending order)
					return objectsEqual(resultArr.Elements[0], sortedElements[i])
				})

				return &Array{Elements: sortedElements}
			},
		},
		"keys": {
			Fn: func(args ...Object) Object {
				if len(args) != 1 {
					return newError("wrong number of arguments to `keys`. got=%d, want=1", len(args))
				}

				dict, ok := args[0].(*Dictionary)
				if !ok {
					return newError("argument to `keys` must be a dictionary, got %s", args[0].Type())
				}

				keys := make([]Object, 0, len(dict.Pairs))
				for key := range dict.Pairs {
					keys = append(keys, &String{Value: key})
				}
				return &Array{Elements: keys}
			},
		},
		"values": {
			Fn: func(args ...Object) Object {
				if len(args) != 1 {
					return newError("wrong number of arguments to `values`. got=%d, want=1", len(args))
				}

				dict, ok := args[0].(*Dictionary)
				if !ok {
					return newError("argument to `values` must be a dictionary, got %s", args[0].Type())
				}

				// Create environment for evaluation with 'this'
				dictEnv := NewEnclosedEnvironment(dict.Env)
				dictEnv.Set("this", dict)

				values := make([]Object, 0, len(dict.Pairs))
				for _, expr := range dict.Pairs {
					val := Eval(expr, dictEnv)
					values = append(values, val)
				}
				return &Array{Elements: values}
			},
		},
		"has": {
			Fn: func(args ...Object) Object {
				if len(args) != 2 {
					return newError("wrong number of arguments to `has`. got=%d, want=2", len(args))
				}

				dict, ok := args[0].(*Dictionary)
				if !ok {
					return newError("first argument to `has` must be a dictionary, got %s", args[0].Type())
				}

				key, ok := args[1].(*String)
				if !ok {
					return newError("second argument to `has` must be a string, got %s", args[1].Type())
				}

				_, exists := dict.Pairs[key.Value]
				return nativeBoolToParsBoolean(exists)
			},
		},
		"toArray": {
			Fn: func(args ...Object) Object {
				if len(args) != 1 {
					return newError("wrong number of arguments to `toArray`. got=%d, want=1", len(args))
				}

				dict, ok := args[0].(*Dictionary)
				if !ok {
					return newError("argument to `toArray` must be a dictionary, got %s", args[0].Type())
				}

				// Create environment for evaluation with 'this'
				dictEnv := NewEnclosedEnvironment(dict.Env)
				dictEnv.Set("this", dict)

				pairs := make([]Object, 0, len(dict.Pairs))
				for key, expr := range dict.Pairs {
					val := Eval(expr, dictEnv)

					// Skip functions with parameters (they can't be called without args)
					if fn, ok := val.(*Function); ok && fn.ParamCount() > 0 {
						continue
					}

					// If it's a function with no parameters, call it
					if fn, ok := val.(*Function); ok && fn.ParamCount() == 0 {
						val = applyFunction(fn, []Object{})
					}

					// Create [key, value] pair
					pair := &Array{Elements: []Object{
						&String{Value: key},
						val,
					}}
					pairs = append(pairs, pair)
				}
				return &Array{Elements: pairs}
			},
		},
		"toDict": {
			Fn: func(args ...Object) Object {
				if len(args) != 1 {
					return newError("wrong number of arguments to `toDict`. got=%d, want=1", len(args))
				}

				arr, ok := args[0].(*Array)
				if !ok {
					return newError("argument to `toDict` must be an array, got %s", args[0].Type())
				}

				dict := &Dictionary{
					Pairs: make(map[string]ast.Expression),
					Env:   NewEnvironment(),
				}

				for _, elem := range arr.Elements {
					pair, ok := elem.(*Array)
					if !ok || len(pair.Elements) != 2 {
						return newError("toDict requires array of [key, value] pairs")
					}

					keyObj, ok := pair.Elements[0].(*String)
					if !ok {
						return newError("dictionary keys must be strings, got %s", pair.Elements[0].Type())
					}

					// Create a literal expression from the value
					valueObj := pair.Elements[1]
					var expr ast.Expression

					switch v := valueObj.(type) {
					case *Integer:
						expr = &ast.IntegerLiteral{Value: v.Value}
					case *Float:
						expr = &ast.FloatLiteral{Value: v.Value}
					case *String:
						expr = &ast.StringLiteral{Value: v.Value}
					case *Boolean:
						expr = &ast.Boolean{Value: v.Value}
					case *Array:
						// For arrays, we'll store a reference that evaluates to the array
						// This is a workaround - store in environment and reference it
						tempKey := "__toDict_temp_" + keyObj.Value
						dict.Env.Set(tempKey, v)
						expr = &ast.Identifier{Value: tempKey}
					default:
						return newError("toDict: unsupported value type %s", valueObj.Type())
					}

					dict.Pairs[keyObj.Value] = expr
				}

				return dict
			},
		},
	}
} // Helper function to evaluate a statement
func evalStatement(stmt ast.Statement, env *Environment) Object {
	switch stmt := stmt.(type) {
	case *ast.ExpressionStatement:
		return Eval(stmt.Expression, env)
	case *ast.ReturnStatement:
		val := Eval(stmt.ReturnValue, env)
		if isError(val) {
			return val
		}
		return &ReturnValue{Value: val}
	default:
		return Eval(stmt, env)
	}
}

// Eval evaluates AST nodes and returns objects
func Eval(node ast.Node, env *Environment) Object {
	switch node := node.(type) {

	// Statements
	case *ast.Program:
		return evalProgram(node.Statements, env)

	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)

	case *ast.BlockStatement:
		return evalBlockStatement(node, env)

	case *ast.LetStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}

		// Handle dictionary destructuring
		if node.DictPattern != nil {
			return evalDictDestructuringAssignment(node.DictPattern, val, env, true)
		}

		// Handle array destructuring assignment
		if len(node.Names) > 0 {
			return evalDestructuringAssignment(node.Names, val, env)
		}

		// Single assignment
		// Special handling for '_' - don't store it
		if node.Name.Value != "_" {
			env.Set(node.Name.Value, val)
		}
		return val

	case *ast.AssignmentStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}

		// Handle dictionary destructuring
		if node.DictPattern != nil {
			return evalDictDestructuringAssignment(node.DictPattern, val, env, false)
		}

		// Handle array destructuring assignment
		if len(node.Names) > 0 {
			return evalDestructuringAssignment(node.Names, val, env)
		}

		// Single assignment
		// Special handling for '_' - don't store it
		if node.Name.Value != "_" {
			env.Update(node.Name.Value, val)
		}
		return val

	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue, env)
		if isError(val) {
			return val
		}
		return &ReturnValue{Value: val}

	// Expressions
	case *ast.IntegerLiteral:
		return &Integer{Value: node.Value}

	case *ast.FloatLiteral:
		return &Float{Value: node.Value}

	case *ast.StringLiteral:
		return &String{Value: node.Value}

	case *ast.TemplateLiteral:
		return evalTemplateLiteral(node, env)

	case *ast.TagLiteral:
		return evalTagLiteral(node, env)

	case *ast.TagPairExpression:
		return evalTagPair(node, env)

	case *ast.TextNode:
		return &String{Value: node.Value}

	case *ast.Boolean:
		return nativeBoolToParsBoolean(node.Value)

	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right)

	case *ast.InfixExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalInfixExpression(node.Operator, left, right)

	case *ast.IfExpression:
		return evalIfExpression(node, env)

	case *ast.Identifier:
		return evalIdentifier(node, env)

	case *ast.FunctionLiteral:
		body := node.Body
		// Use new-style params if available, otherwise fall back to old parameters
		if len(node.Params) > 0 {
			return &Function{Params: node.Params, Env: env, Body: body}
		}
		return &Function{Parameters: node.Parameters, Env: env, Body: body}

	case *ast.ArrayLiteral:
		elements := evalExpressions(node.Elements, env)
		if len(elements) == 1 && isError(elements[0]) {
			return elements[0]
		}
		return &Array{Elements: elements}

	case *ast.DictionaryLiteral:
		return evalDictionaryLiteral(node, env)

	case *ast.DotExpression:
		return evalDotExpression(node, env)

	case *ast.DeleteStatement:
		return evalDeleteStatement(node, env)

	case *ast.CallExpression:
		// Store current token in environment for logLine
		env.LastToken = &node.Token

		// Check if this is a call to logLine
		if ident, ok := node.Function.(*ast.Identifier); ok && ident.Value == "logLine" {
			args := evalExpressions(node.Arguments, env)
			if len(args) == 1 && isError(args[0]) {
				return args[0]
			}
			return evalLogLine(args, env)
		}

		function := Eval(node.Function, env)
		if isError(function) {
			return function
		}
		args := evalExpressions(node.Arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}
		return applyFunctionWithEnv(function, args, env)

	case *ast.ForExpression:
		return evalForExpression(node, env)

	case *ast.IndexExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		index := Eval(node.Index, env)
		if isError(index) {
			return index
		}
		return evalIndexExpression(left, index)

	case *ast.SliceExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}

		var start, end Object
		if node.Start != nil {
			start = Eval(node.Start, env)
			if isError(start) {
				return start
			}
		}
		if node.End != nil {
			end = Eval(node.End, env)
			if isError(end) {
				return end
			}
		}
		return evalSliceExpression(left, start, end)
	}

	return newError("unknown node type: %T", node)
}

// Helper functions
func evalProgram(stmts []ast.Statement, env *Environment) Object {
	var result Object

	for _, statement := range stmts {
		result = Eval(statement, env)

		switch result := result.(type) {
		case *ReturnValue:
			return result.Value
		case *Error:
			return result
		}
	}

	return result
}

func evalBlockStatement(block *ast.BlockStatement, env *Environment) Object {
	var result Object

	for _, statement := range block.Statements {
		result = Eval(statement, env)

		if result != nil {
			rt := result.Type()
			if rt == RETURN_OBJ || rt == ERROR_OBJ {
				return result
			}
		}
	}

	return result
}

func nativeBoolToParsBoolean(input bool) *Boolean {
	if input {
		return TRUE
	}
	return FALSE
}

func evalPrefixExpression(operator string, right Object) Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "not":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	default:
		return newError("unknown operator: %s%s", operator, right.Type())
	}
}

func evalBangOperatorExpression(right Object) Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		return FALSE
	}
}

func evalMinusPrefixOperatorExpression(right Object) Object {
	if right.Type() != INTEGER_OBJ {
		return newError("unknown operator: -%s", right.Type())
	}

	value := right.(*Integer).Value
	return &Integer{Value: -value}
}

func evalInfixExpression(operator string, left, right Object) Object {
	switch {
	case operator == "&" || operator == "and":
		return nativeBoolToParsBoolean(isTruthy(left) && isTruthy(right))
	case operator == "|" || operator == "or":
		return nativeBoolToParsBoolean(isTruthy(left) || isTruthy(right))
	case operator == "++":
		return evalConcatExpression(left, right)
	case operator == "+" && (left.Type() == STRING_OBJ || right.Type() == STRING_OBJ):
		// String concatenation with automatic type conversion
		return evalStringConcatExpression(left, right)
	case left.Type() == INTEGER_OBJ && right.Type() == INTEGER_OBJ:
		return evalIntegerInfixExpression(operator, left, right)
	case left.Type() == FLOAT_OBJ && right.Type() == FLOAT_OBJ:
		return evalFloatInfixExpression(operator, left, right)
	case left.Type() == INTEGER_OBJ && right.Type() == FLOAT_OBJ:
		return evalMixedInfixExpression(operator, left, right)
	case left.Type() == FLOAT_OBJ && right.Type() == INTEGER_OBJ:
		return evalMixedInfixExpression(operator, left, right)
	case left.Type() == STRING_OBJ && right.Type() == STRING_OBJ:
		return evalStringInfixExpression(operator, left, right)
	case operator == "==":
		return nativeBoolToParsBoolean(left == right)
	case operator == "!=":
		return nativeBoolToParsBoolean(left != right)
	case left.Type() != right.Type():
		return newError("type mismatch: %s %s %s", left.Type(), operator, right.Type())
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalIntegerInfixExpression(operator string, left, right Object) Object {
	leftVal := left.(*Integer).Value
	rightVal := right.(*Integer).Value

	switch operator {
	case "+":
		return &Integer{Value: leftVal + rightVal}
	case "-":
		return &Integer{Value: leftVal - rightVal}
	case "*":
		return &Integer{Value: leftVal * rightVal}
	case "/":
		if rightVal == 0 {
			return newError("division by zero")
		}
		return &Integer{Value: leftVal / rightVal}
	case "%":
		if rightVal == 0 {
			return newError("modulo by zero")
		}
		return &Integer{Value: leftVal % rightVal}
	case "<":
		return nativeBoolToParsBoolean(leftVal < rightVal)
	case ">":
		return nativeBoolToParsBoolean(leftVal > rightVal)
	case "<=":
		return nativeBoolToParsBoolean(leftVal <= rightVal)
	case ">=":
		return nativeBoolToParsBoolean(leftVal >= rightVal)
	case "==":
		return nativeBoolToParsBoolean(leftVal == rightVal)
	case "!=":
		return nativeBoolToParsBoolean(leftVal != rightVal)
	default:
		return newError("unknown operator: %s", operator)
	}
}

func evalFloatInfixExpression(operator string, left, right Object) Object {
	leftVal := left.(*Float).Value
	rightVal := right.(*Float).Value

	switch operator {
	case "+":
		return &Float{Value: leftVal + rightVal}
	case "-":
		return &Float{Value: leftVal - rightVal}
	case "*":
		return &Float{Value: leftVal * rightVal}
	case "/":
		if rightVal == 0 {
			return newError("division by zero")
		}
		return &Float{Value: leftVal / rightVal}
	case "<":
		return nativeBoolToParsBoolean(leftVal < rightVal)
	case ">":
		return nativeBoolToParsBoolean(leftVal > rightVal)
	case "<=":
		return nativeBoolToParsBoolean(leftVal <= rightVal)
	case ">=":
		return nativeBoolToParsBoolean(leftVal >= rightVal)
	case "==":
		return nativeBoolToParsBoolean(leftVal == rightVal)
	case "!=":
		return nativeBoolToParsBoolean(leftVal != rightVal)
	default:
		return newError("unknown operator: %s", operator)
	}
}

func evalMixedInfixExpression(operator string, left, right Object) Object {
	var leftVal, rightVal float64

	// Convert both operands to float64
	switch left := left.(type) {
	case *Integer:
		leftVal = float64(left.Value)
	case *Float:
		leftVal = left.Value
	default:
		return newError("unsupported type for mixed arithmetic: %T", left)
	}

	switch right := right.(type) {
	case *Integer:
		rightVal = float64(right.Value)
	case *Float:
		rightVal = right.Value
	default:
		return newError("unsupported type for mixed arithmetic: %T", right)
	}

	switch operator {
	case "+":
		return &Float{Value: leftVal + rightVal}
	case "-":
		return &Float{Value: leftVal - rightVal}
	case "*":
		return &Float{Value: leftVal * rightVal}
	case "/":
		if rightVal == 0 {
			return newError("division by zero")
		}
		return &Float{Value: leftVal / rightVal}
	case "<":
		return nativeBoolToParsBoolean(leftVal < rightVal)
	case ">":
		return nativeBoolToParsBoolean(leftVal > rightVal)
	case "<=":
		return nativeBoolToParsBoolean(leftVal <= rightVal)
	case ">=":
		return nativeBoolToParsBoolean(leftVal >= rightVal)
	case "==":
		return nativeBoolToParsBoolean(leftVal == rightVal)
	case "!=":
		return nativeBoolToParsBoolean(leftVal != rightVal)
	default:
		return newError("unknown operator: %s", operator)
	}
}

func evalStringInfixExpression(operator string, left, right Object) Object {
	leftVal := left.(*String).Value
	rightVal := right.(*String).Value

	switch operator {
	case "+":
		return &String{Value: leftVal + rightVal}
	case "==":
		return nativeBoolToParsBoolean(leftVal == rightVal)
	case "!=":
		return nativeBoolToParsBoolean(leftVal != rightVal)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

// evalStringConcatExpression handles string concatenation with automatic type conversion
func evalStringConcatExpression(left, right Object) Object {
	leftStr := objectToTemplateString(left)
	rightStr := objectToTemplateString(right)
	return &String{Value: leftStr + rightStr}
}

func evalIfExpression(ie *ast.IfExpression, env *Environment) Object {
	condition := Eval(ie.Condition, env)
	if isError(condition) {
		return condition
	}

	if isTruthy(condition) {
		return Eval(ie.Consequence, env)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative, env)
	} else {
		return NULL
	}
}

func isTruthy(obj Object) bool {
	switch obj {
	case NULL:
		return false
	case TRUE:
		return true
	case FALSE:
		return false
	default:
		return true
	}
}

func evalIdentifier(node *ast.Identifier, env *Environment) Object {
	// Special handling for '_' - always returns null
	if node.Value == "_" {
		return NULL
	}

	val, ok := env.Get(node.Value)
	if !ok {
		if builtin, ok := getBuiltins()[node.Value]; ok {
			return builtin
		}
		return newError("identifier not found: " + node.Value)
	}

	return val
}

func evalExpressions(exps []ast.Expression, env *Environment) []Object {
	var result []Object

	for _, e := range exps {
		evaluated := Eval(e, env)
		if isError(evaluated) {
			return []Object{evaluated}
		}
		result = append(result, evaluated)
	}

	return result
}

func applyFunction(fn Object, args []Object) Object {
	switch fn := fn.(type) {
	case *Function:
		extendedEnv := extendFunctionEnv(fn, args)
		evaluated := Eval(fn.Body, extendedEnv)
		return unwrapReturnValue(evaluated)
	case *Builtin:
		return fn.Fn(args...)
	default:
		return newError("not a function: %T", fn)
	}
}

func applyFunctionWithEnv(fn Object, args []Object, env *Environment) Object {
	switch fn := fn.(type) {
	case *Function:
		extendedEnv := extendFunctionEnv(fn, args)
		evaluated := Eval(fn.Body, extendedEnv)
		return unwrapReturnValue(evaluated)
	case *Builtin:
		return fn.Fn(args...)
	default:
		return newError("not a function: %T", fn)
	}
}

// evalLogLine implements logLine with filename and line number
func evalLogLine(args []Object, env *Environment) Object {
	var result strings.Builder

	// Add filename and line number prefix
	filename := env.Filename
	if filename == "" {
		filename = "<unknown>"
	}
	line := 1
	if env.LastToken != nil {
		line = env.LastToken.Line
	}
	result.WriteString(fmt.Sprintf("%s:%d: ", filename, line))

	// Process arguments like log()
	for i, arg := range args {
		if i == 0 {
			// First argument: if it's a string, show without quotes
			if str, ok := arg.(*String); ok {
				result.WriteString(str.Value)
			} else {
				result.WriteString(objectToDebugString(arg))
			}
		} else {
			// Subsequent arguments: add separator and debug format
			if i == 1 {
				// After first string, no comma - just space
				if _, firstWasString := args[0].(*String); firstWasString {
					result.WriteString(" ")
				} else {
					result.WriteString(", ")
				}
			} else {
				result.WriteString(", ")
			}
			result.WriteString(objectToDebugString(arg))
		}
	}

	// Write immediately to stdout
	fmt.Fprintln(os.Stdout, result.String())

	// Return null
	return NULL
}

func extendFunctionEnv(fn *Function, args []Object) *Environment {
	env := NewEnclosedEnvironment(fn.Env)

	// Use new-style parameters if available
	if len(fn.Params) > 0 {
		for paramIdx, param := range fn.Params {
			if paramIdx >= len(args) {
				break
			}
			arg := args[paramIdx]

			// Handle different parameter types
			if param.DictPattern != nil {
				// Dictionary destructuring
				evalDictDestructuringAssignment(param.DictPattern, arg, env, true)
			} else if len(param.ArrayPattern) > 0 {
				// Array destructuring
				evalArrayDestructuringForParam(param.ArrayPattern, arg, env)
			} else if param.Ident != nil {
				// Simple identifier
				env.Set(param.Ident.Value, arg)
			}
		}
	} else {
		// Fallback to old-style parameters
		for paramIdx, param := range fn.Parameters {
			if paramIdx >= len(args) {
				break
			}
			env.Set(param.Value, args[paramIdx])
		}
	}

	return env
}

// evalArrayDestructuringForParam handles array destructuring in function parameters
func evalArrayDestructuringForParam(pattern []*ast.Identifier, val Object, env *Environment) {
	// Convert value to array if it isn't already
	var elements []Object

	switch v := val.(type) {
	case *Array:
		elements = v.Elements
	default:
		// Single value becomes single-element array
		elements = []Object{v}
	}

	// Assign each element to corresponding variable
	for i, name := range pattern {
		if i < len(elements) {
			if name.Value != "_" {
				env.Set(name.Value, elements[i])
			}
		} else {
			// No more elements, assign null
			if name.Value != "_" {
				env.Set(name.Value, NULL)
			}
		}
	}

	// If there are more elements than names, assign remaining as array to last variable
	if len(elements) > len(pattern) && len(pattern) > 0 {
		lastIdx := len(pattern) - 1
		lastName := pattern[lastIdx]
		if lastName.Value != "_" {
			// Replace the last assignment with an array of remaining elements
			remaining := &Array{Elements: elements[lastIdx:]}
			env.Set(lastName.Value, remaining)
		}
	}
}

func unwrapReturnValue(obj Object) Object {
	if returnValue, ok := obj.(*ReturnValue); ok {
		return returnValue.Value
	}
	return obj
}

// evalForExpression evaluates for expressions
func evalForExpression(node *ast.ForExpression, env *Environment) Object {
	// Evaluate the array/dict expression
	iterableObj := Eval(node.Array, env)
	if isError(iterableObj) {
		return iterableObj
	}

	// Handle dictionary iteration
	if dict, ok := iterableObj.(*Dictionary); ok {
		return evalForDictExpression(node, dict, env)
	}

	// Convert to array (handle strings as rune arrays)
	var elements []Object
	switch arr := iterableObj.(type) {
	case *Array:
		elements = arr.Elements
	case *String:
		// Convert string to array of single-character strings
		runes := []rune(arr.Value)
		elements = make([]Object, len(runes))
		for i, r := range runes {
			elements[i] = &String{Value: string(r)}
		}
	default:
		return newError("for expects an array, string, or dictionary, got %s", iterableObj.Type())
	}

	// Determine which function to use
	var fn Object
	if node.Function != nil {
		// Simple form: for(array) func
		fn = Eval(node.Function, env)
		if isError(fn) {
			return fn
		}
		// Accept both functions and builtins
		switch fn.(type) {
		case *Function, *Builtin:
			// OK
		default:
			return newError("for expects a function or builtin, got %s", fn.Type())
		}
	} else if node.Body != nil {
		// 'in' form: for(var in array) body
		// node.Body is already a FunctionLiteral with the variable as parameter
		fn = &Function{
			Parameters: node.Body.(*ast.FunctionLiteral).Parameters,
			Body:       node.Body.(*ast.FunctionLiteral).Body,
			Env:        env,
		}
	} else {
		return newError("for expression missing function or body")
	}

	// Map function over array elements
	result := []Object{}
	for _, elem := range elements {
		var evaluated Object

		switch f := fn.(type) {
		case *Builtin:
			// Call builtin with single element
			evaluated = f.Fn(elem)
		case *Function:
			// Call user function
			if f.ParamCount() != 1 {
				return newError("function passed to for must take exactly 1 parameter, got %d", f.ParamCount())
			}

			// Create a new environment and bind the parameter
			extendedEnv := extendFunctionEnv(f, []Object{elem})

			// Evaluate all statements in the body
			for _, stmt := range f.Body.Statements {
				evaluated = evalStatement(stmt, extendedEnv)
				if returnValue, ok := evaluated.(*ReturnValue); ok {
					evaluated = returnValue.Value
					break
				}
				if isError(evaluated) {
					return evaluated
				}
			}
		}

		// Skip null values (filter behavior)
		if evaluated != NULL {
			result = append(result, evaluated)
		}
	}

	return &Array{Elements: result}
}

// evalForDictExpression handles for loops over dictionaries
func evalForDictExpression(node *ast.ForExpression, dict *Dictionary, env *Environment) Object {
	// Create environment for evaluation with 'this'
	dictEnv := NewEnclosedEnvironment(dict.Env)
	dictEnv.Set("this", dict)

	// Determine which function to use
	var fn *Function
	if node.Body != nil {
		bodyFn := node.Body.(*ast.FunctionLiteral)
		if len(bodyFn.Params) > 0 {
			fn = &Function{
				Params: bodyFn.Params,
				Body:   bodyFn.Body,
				Env:    env,
			}
		} else {
			fn = &Function{
				Parameters: bodyFn.Parameters,
				Body:       bodyFn.Body,
				Env:        env,
			}
		}
	} else {
		return newError("for loop over dictionary requires body with key, value parameters")
	}

	// Check parameter count
	if fn.ParamCount() != 2 {
		return newError("for loop over dictionary requires exactly 2 parameters (key, value), got %d", fn.ParamCount())
	}

	// Iterate over dictionary key-value pairs
	result := []Object{}
	for key, expr := range dict.Pairs {
		// Evaluate the value
		value := Eval(expr, dictEnv)
		if isError(value) {
			return value
		}

		// Create environment for loop body with both key and value
		extendedEnv := extendFunctionEnv(fn, []Object{&String{Value: key}, value})

		// Evaluate all statements in the body
		var evaluated Object
		for _, stmt := range fn.Body.Statements {
			evaluated = evalStatement(stmt, extendedEnv)
			if returnValue, ok := evaluated.(*ReturnValue); ok {
				evaluated = returnValue.Value
				break
			}
			if isError(evaluated) {
				return evaluated
			}
		}

		// Skip null values (filter behavior)
		if evaluated != NULL {
			result = append(result, evaluated)
		}
	}

	return &Array{Elements: result}
}

func newError(format string, a ...interface{}) *Error {
	return &Error{Message: fmt.Sprintf(format, a...)}
}

func isError(obj Object) bool {
	if obj != nil {
		return obj.Type() == ERROR_OBJ
	}
	return false
}

// evalDestructuringAssignment handles array destructuring assignment
func evalDestructuringAssignment(names []*ast.Identifier, val Object, env *Environment) Object {
	// Convert value to array if it isn't already
	var elements []Object

	switch v := val.(type) {
	case *Array:
		elements = v.Elements
	default:
		// Single value becomes single-element array
		elements = []Object{v}
	}

	// Assign each element to corresponding variable
	for i, name := range names {
		if i < len(elements) {
			// Direct assignment for elements within bounds
			if name.Value != "_" {
				env.Update(name.Value, elements[i])
			}
		} else {
			// No more elements, assign null
			if name.Value != "_" {
				env.Update(name.Value, NULL)
			}
		}
	}

	// If there are more elements than names, assign remaining as array to last variable
	if len(elements) > len(names) && len(names) > 0 {
		lastIdx := len(names) - 1
		lastName := names[lastIdx]
		if lastName.Value != "_" {
			// Replace the last assignment with an array of remaining elements
			remaining := &Array{Elements: elements[lastIdx:]}
			env.Update(lastName.Value, remaining)
		}
	}

	// Return the original value
	return val
}

// evalDictDestructuringAssignment evaluates dictionary destructuring patterns
func evalDictDestructuringAssignment(pattern *ast.DictDestructuringPattern, val Object, env *Environment, isLet bool) Object {
	// Type check: value must be a dictionary
	dict, ok := val.(*Dictionary)
	if !ok {
		return newError("dictionary destructuring requires a dictionary value, got %s", val.Type())
	}

	// Track which keys we've extracted (for rest operator)
	extractedKeys := make(map[string]bool)

	// Process each key in the pattern
	for _, keyPattern := range pattern.Keys {
		keyName := keyPattern.Key.Value
		extractedKeys[keyName] = true

		// Get expression from dictionary and evaluate it
		var value Object
		if expr, exists := dict.Pairs[keyName]; exists {
			// Evaluate the expression in the dictionary's environment
			value = Eval(expr, dict.Env)
			if isError(value) {
				return value
			}
		} else {
			// If key not found, assign null
			value = NULL
		}

		// Handle nested destructuring
		if keyPattern.Nested != nil {
			if nestedPattern, ok := keyPattern.Nested.(*ast.DictDestructuringPattern); ok {
				result := evalDictDestructuringAssignment(nestedPattern, value, env, isLet)
				if isError(result) {
					return result
				}
			} else {
				return newError("unsupported nested destructuring pattern")
			}
		} else {
			// Determine the target variable name (alias or original key)
			targetName := keyName
			if keyPattern.Alias != nil {
				targetName = keyPattern.Alias.Value
			}

			// Assign to environment
			if targetName != "_" {
				if isLet {
					env.Set(targetName, value)
				} else {
					env.Update(targetName, value)
				}
			}
		}
	}

	// Handle rest operator
	if pattern.Rest != nil {
		restPairs := make(map[string]ast.Expression)
		for key, expr := range dict.Pairs {
			if !extractedKeys[key] {
				restPairs[key] = expr
			}
		}

		restDict := &Dictionary{Pairs: restPairs, Env: dict.Env}
		if pattern.Rest.Value != "_" {
			if isLet {
				env.Set(pattern.Rest.Value, restDict)
			} else {
				env.Update(pattern.Rest.Value, restDict)
			}
		}
	}

	// Return the original value
	return val
}

// evalTemplateLiteral evaluates a template literal with interpolation
func evalTemplateLiteral(node *ast.TemplateLiteral, env *Environment) Object {
	template := node.Value
	var result strings.Builder

	i := 0
	for i < len(template) {
		// Check for escaped brace markers \0{ and \0}
		if i < len(template)-2 && template[i] == '\\' && template[i+1] == '0' {
			if template[i+2] == '{' {
				result.WriteByte('{')
				i += 3
				continue
			} else if template[i+2] == '}' {
				result.WriteByte('}')
				i += 3
				continue
			}
		}

		// Look for {
		if template[i] == '{' {
			// Find the closing }
			i++ // skip {
			braceCount := 1
			exprStart := i

			for i < len(template) && braceCount > 0 {
				if template[i] == '{' {
					braceCount++
				} else if template[i] == '}' {
					braceCount--
				}
				if braceCount > 0 {
					i++
				}
			}

			if braceCount != 0 {
				return newError("unclosed { in template literal")
			}

			// Extract and evaluate the expression
			exprStr := template[exprStart:i]
			i++ // skip closing }

			// Parse and evaluate the expression
			l := lexer.New(exprStr)
			p := parser.New(l)
			program := p.ParseProgram()

			if len(p.Errors()) > 0 {
				return newError("error parsing template expression: %s", p.Errors()[0])
			}

			// Evaluate the expression
			var evaluated Object
			for _, stmt := range program.Statements {
				evaluated = Eval(stmt, env)
				if isError(evaluated) {
					return evaluated
				}
			}

			// Convert result to string
			if evaluated != nil {
				result.WriteString(objectToTemplateString(evaluated))
			}
		} else {
			// Regular character
			result.WriteByte(template[i])
			i++
		}
	}

	return &String{Value: result.String()}
}

// evalTagLiteral evaluates a singleton tag
func evalTagLiteral(node *ast.TagLiteral, env *Environment) Object {
	raw := node.Raw

	// Parse tag name (first word)
	i := 0
	for i < len(raw) && !unicode.IsSpace(rune(raw[i])) {
		i++
	}
	tagName := raw[:i]
	rest := raw[i:]

	// Check if it's a custom tag (starts with uppercase)
	isCustom := len(tagName) > 0 && unicode.IsUpper(rune(tagName[0]))

	if isCustom {
		// Custom tag - call function with props dictionary
		return evalCustomTag(tagName, rest, env)
	} else {
		// Standard tag - return as interpolated string
		return evalStandardTag(tagName, rest, env)
	}
}

// evalTagPair evaluates a paired tag like <div>content</div> or <Component>content</Component>
func evalTagPair(node *ast.TagPairExpression, env *Environment) Object {
	// Empty grouping tag <> just returns its contents
	if node.Name == "" {
		return evalTagContents(node.Contents, env)
	}

	// Check if it's a custom component (starts with uppercase)
	isCustom := len(node.Name) > 0 && unicode.IsUpper(rune(node.Name[0]))

	if isCustom {
		// Custom component - call function with props dictionary including contents
		return evalCustomTagPair(node, env)
	} else {
		// Standard tag - return as HTML string
		return evalStandardTagPair(node, env)
	}
}

// evalStandardTagPair evaluates a standard (lowercase) tag pair as HTML string
func evalStandardTagPair(node *ast.TagPairExpression, env *Environment) Object {
	var result strings.Builder

	result.WriteByte('<')
	result.WriteString(node.Name)

	// Process props with interpolation (similar to singleton tags)
	if node.Props != "" {
		result.WriteByte(' ')
		propsResult := evalTagProps(node.Props, env)
		if isError(propsResult) {
			return propsResult
		}
		result.WriteString(propsResult.(*String).Value)
	}

	result.WriteByte('>')

	// Evaluate and append contents
	contentsObj := evalTagContents(node.Contents, env)
	if isError(contentsObj) {
		return contentsObj
	}
	result.WriteString(contentsObj.(*String).Value)

	result.WriteString("</")
	result.WriteString(node.Name)
	result.WriteByte('>')

	return &String{Value: result.String()}
}

// evalCustomTagPair evaluates a custom (uppercase) tag pair as a function call
func evalCustomTagPair(node *ast.TagPairExpression, env *Environment) Object {
	// Look up the component function
	fn, ok := env.Get(node.Name)
	if !ok {
		return newError("undefined component: %s", node.Name)
	}

	// Parse props into a dictionary and add contents
	propsDict := parseTagProps(node.Props, env)
	if isError(propsDict) {
		return propsDict
	}

	dict := propsDict.(*Dictionary)

	// Evaluate contents and add to props as "contents"
	contentsObj := evalTagContentsAsArray(node.Contents, env)
	if isError(contentsObj) {
		return contentsObj
	}

	// Create a literal expression for the contents array
	// We need to wrap the evaluated contents in an expression
	dict.Pairs["contents"] = &ast.ArrayLiteral{Elements: []ast.Expression{}}

	// Store the evaluated contents directly in the environment temporarily
	contentsEnv := NewEnclosedEnvironment(env)
	contentsEnv.Set("__tag_contents__", contentsObj)

	// Actually, let's simplify - evaluate contents as a single value
	if contentsArray, ok := contentsObj.(*Array); ok && len(contentsArray.Elements) == 1 {
		// Single item - pass directly
		dict.Pairs["contents"] = createLiteralExpression(contentsArray.Elements[0])
	} else {
		// Multiple items or empty - pass as array
		dict.Pairs["contents"] = createLiteralExpression(contentsObj)
	}

	// Call the function with the props dictionary
	return applyFunction(fn, []Object{dict})
}

// evalTagContents evaluates tag contents and returns as a concatenated string
func evalTagContents(contents []ast.Node, env *Environment) Object {
	var result strings.Builder

	for _, node := range contents {
		obj := Eval(node, env)
		if isError(obj) {
			return obj
		}
		result.WriteString(objectToTemplateString(obj))
	}

	return &String{Value: result.String()}
}

// evalTagContentsAsArray evaluates tag contents and returns as an array
func evalTagContentsAsArray(contents []ast.Node, env *Environment) Object {
	if len(contents) == 0 {
		return NULL
	}

	elements := make([]Object, 0, len(contents))
	for _, node := range contents {
		obj := Eval(node, env)
		if isError(obj) {
			return obj
		}
		// Convert to string for consistency
		elements = append(elements, &String{Value: objectToTemplateString(obj)})
	}

	return &Array{Elements: elements}
}

// evalTagProps evaluates tag props string with interpolations
func evalTagProps(propsStr string, env *Environment) Object {
	var result strings.Builder

	i := 0
	for i < len(propsStr) {
		// Look for {expr}
		if propsStr[i] == '{' {
			// Find the closing }
			i++ // skip {
			braceCount := 1
			exprStart := i

			for i < len(propsStr) && braceCount > 0 {
				if propsStr[i] == '"' {
					// Skip quoted strings
					i++
					for i < len(propsStr) && propsStr[i] != '"' {
						if propsStr[i] == '\\' {
							i += 2
						} else {
							i++
						}
					}
					if i < len(propsStr) {
						i++
					}
					continue
				}
				if propsStr[i] == '{' {
					braceCount++
				} else if propsStr[i] == '}' {
					braceCount--
				}
				if braceCount > 0 {
					i++
				}
			}

			if braceCount != 0 {
				return newError("unclosed { in tag props")
			}

			// Extract and evaluate the expression
			exprStr := propsStr[exprStart:i]
			i++ // skip closing }

			// Parse and evaluate the expression
			l := lexer.New(exprStr)
			p := parser.New(l)
			program := p.ParseProgram()

			if len(p.Errors()) > 0 {
				return newError("error parsing tag prop expression: %s", p.Errors()[0])
			}

			// Evaluate the expression
			var evaluated Object
			for _, stmt := range program.Statements {
				evaluated = Eval(stmt, env)
				if isError(evaluated) {
					return evaluated
				}
			}

			// Convert result to string
			if evaluated != nil {
				result.WriteString(objectToTemplateString(evaluated))
			}
		} else {
			// Regular character
			result.WriteByte(propsStr[i])
			i++
		}
	}

	return &String{Value: result.String()}
}

// createLiteralExpression creates an AST expression from an evaluated object
// This is a helper for passing evaluated values back through the AST
func createLiteralExpression(obj Object) ast.Expression {
	switch obj := obj.(type) {
	case *Integer:
		return &ast.IntegerLiteral{Value: obj.Value}
	case *Float:
		return &ast.FloatLiteral{Value: obj.Value}
	case *String:
		return &ast.StringLiteral{Value: obj.Value}
	case *Boolean:
		return &ast.Boolean{Value: obj.Value}
	case *Array:
		// For arrays, create array literal with elements
		elements := make([]ast.Expression, len(obj.Elements))
		for i, elem := range obj.Elements {
			elements[i] = createLiteralExpression(elem)
		}
		return &ast.ArrayLiteral{Elements: elements}
	default:
		// For other types, return a string literal
		return &ast.StringLiteral{Value: obj.Inspect()}
	}
}

// evalStandardTag evaluates a standard (lowercase) tag as an interpolated string
func evalStandardTag(tagName string, propsStr string, env *Environment) Object {
	var result strings.Builder
	result.WriteByte('<')
	result.WriteString(tagName)

	// Process props with interpolation
	i := 0
	for i < len(propsStr) {
		// Look for {expr}
		if propsStr[i] == '{' {
			// Find the closing }
			i++ // skip {
			braceCount := 1
			exprStart := i

			for i < len(propsStr) && braceCount > 0 {
				if propsStr[i] == '"' {
					// Skip quoted strings
					i++
					for i < len(propsStr) && propsStr[i] != '"' {
						if propsStr[i] == '\\' {
							i += 2
						} else {
							i++
						}
					}
					if i < len(propsStr) {
						i++
					}
					continue
				}
				if propsStr[i] == '{' {
					braceCount++
				} else if propsStr[i] == '}' {
					braceCount--
				}
				if braceCount > 0 {
					i++
				}
			}

			if braceCount != 0 {
				return newError("unclosed { in tag")
			}

			// Extract and evaluate the expression
			exprStr := propsStr[exprStart:i]
			i++ // skip closing }

			// Parse and evaluate the expression
			l := lexer.New(exprStr)
			p := parser.New(l)
			program := p.ParseProgram()

			if len(p.Errors()) > 0 {
				return newError("error parsing tag expression: %s", p.Errors()[0])
			}

			// Evaluate the expression
			var evaluated Object
			for _, stmt := range program.Statements {
				evaluated = Eval(stmt, env)
				if isError(evaluated) {
					return evaluated
				}
			}

			// Convert result to string (don't add quotes - they should be in the tag already)
			if evaluated != nil {
				result.WriteString(objectToTemplateString(evaluated))
			}
		} else {
			// Regular character
			result.WriteByte(propsStr[i])
			i++
		}
	}

	result.WriteString(" />")
	return &String{Value: result.String()}
}

// evalCustomTag evaluates a custom (uppercase) tag as a function call
func evalCustomTag(tagName string, propsStr string, env *Environment) Object {
	// Look up the function
	fn, ok := env.Get(tagName)
	if !ok {
		if builtin, ok := getBuiltins()[tagName]; ok {
			fn = builtin
		} else {
			return newError("function not found: " + tagName)
		}
	}

	// Parse props into a dictionary
	props := parseTagProps(propsStr, env)
	if isError(props) {
		return props
	}

	// Call the function with the props dictionary
	return applyFunction(fn, []Object{props})
}

// parseTagProps parses tag properties into a dictionary
func parseTagProps(propsStr string, env *Environment) Object {
	pairs := make(map[string]ast.Expression)

	i := 0
	for i < len(propsStr) {
		// Skip whitespace
		for i < len(propsStr) && unicode.IsSpace(rune(propsStr[i])) {
			i++
		}
		if i >= len(propsStr) {
			break
		}

		// Read prop name
		nameStart := i
		for i < len(propsStr) && !unicode.IsSpace(rune(propsStr[i])) && propsStr[i] != '=' {
			i++
		}
		if nameStart == i {
			break
		}
		propName := propsStr[nameStart:i]

		// Skip whitespace
		for i < len(propsStr) && unicode.IsSpace(rune(propsStr[i])) {
			i++
		}

		// Check for = or standalone prop
		if i >= len(propsStr) || propsStr[i] != '=' {
			// Standalone prop (boolean)
			pairs[propName] = &ast.Boolean{Value: true}
			continue
		}

		i++ // skip =

		// Skip whitespace
		for i < len(propsStr) && unicode.IsSpace(rune(propsStr[i])) {
			i++
		}

		if i >= len(propsStr) {
			break
		}

		// Read prop value
		var valueStr string
		if propsStr[i] == '"' {
			// Quoted string - check if it contains interpolation
			i++ // skip opening quote
			valueStart := i
			hasInterpolation := false
			tempI := i
			for tempI < len(propsStr) && propsStr[tempI] != '"' {
				if propsStr[tempI] == '{' {
					hasInterpolation = true
					break
				}
				if propsStr[tempI] == '\\' {
					tempI += 2
				} else {
					tempI++
				}
			}

			if hasInterpolation {
				// The string contains {expr}, treat it as an interpolation
				// Extract content between quotes
				for i < len(propsStr) && propsStr[i] != '"' {
					if propsStr[i] == '\\' {
						i += 2
					} else {
						i++
					}
				}
				valueStr = propsStr[valueStart:i]
				if i < len(propsStr) {
					i++ // skip closing quote
				}

				// Now parse the interpolation - find the {expr}
				j := 0
				for j < len(valueStr) {
					if valueStr[j] == '{' {
						j++ // skip {
						exprStart := j
						braceCount := 1
						for j < len(valueStr) && braceCount > 0 {
							if valueStr[j] == '{' {
								braceCount++
							} else if valueStr[j] == '}' {
								braceCount--
							}
							if braceCount > 0 {
								j++
							}
						}
						exprStr := valueStr[exprStart:j]
						// Parse the expression
						l := lexer.New(exprStr)
						p := parser.New(l)
						program := p.ParseProgram()

						if len(p.Errors()) > 0 {
							return newError("error parsing tag prop expression: %s", p.Errors()[0])
						}

						// Store as expression statement
						if len(program.Statements) > 0 {
							if exprStmt, ok := program.Statements[0].(*ast.ExpressionStatement); ok {
								pairs[propName] = exprStmt.Expression
							}
						}
						break
					}
					j++
				}
			} else {
				// Plain string with no interpolation
				for i < len(propsStr) && propsStr[i] != '"' {
					if propsStr[i] == '\\' {
						i += 2
					} else {
						i++
					}
				}
				valueStr = propsStr[valueStart:i]
				if i < len(propsStr) {
					i++ // skip closing quote
				}
				pairs[propName] = &ast.StringLiteral{Value: valueStr}
			}
		} else if propsStr[i] == '{' {
			// Expression in braces
			i++ // skip {
			braceCount := 1
			exprStart := i

			for i < len(propsStr) && braceCount > 0 {
				if propsStr[i] == '"' {
					// Skip quoted strings
					i++
					for i < len(propsStr) && propsStr[i] != '"' {
						if propsStr[i] == '\\' {
							i += 2
						} else {
							i++
						}
					}
					if i < len(propsStr) {
						i++
					}
					continue
				}
				if propsStr[i] == '{' {
					braceCount++
				} else if propsStr[i] == '}' {
					braceCount--
				}
				if braceCount > 0 {
					i++
				}
			}

			if braceCount != 0 {
				return newError("unclosed { in tag prop")
			}

			exprStr := propsStr[exprStart:i]
			i++ // skip }

			// Parse the expression
			l := lexer.New(exprStr)
			p := parser.New(l)
			program := p.ParseProgram()

			if len(p.Errors()) > 0 {
				return newError("error parsing tag prop expression: %s", p.Errors()[0])
			}

			// Store as expression statement
			if len(program.Statements) > 0 {
				if exprStmt, ok := program.Statements[0].(*ast.ExpressionStatement); ok {
					pairs[propName] = exprStmt.Expression
				}
			}
		}
	}

	return &Dictionary{Pairs: pairs, Env: env}
}

// objectToTemplateString converts an object to its string representation for template interpolation
func objectToTemplateString(obj Object) string {
	switch obj := obj.(type) {
	case *Integer:
		return strconv.FormatInt(obj.Value, 10)
	case *Float:
		return fmt.Sprintf("%g", obj.Value)
	case *Boolean:
		if obj.Value {
			return "true"
		}
		return "false"
	case *String:
		return obj.Value
	case *Array:
		// Arrays are printed without commas in templates
		var result strings.Builder
		for _, elem := range obj.Elements {
			result.WriteString(objectToTemplateString(elem))
		}
		return result.String()
	case *Null:
		return ""
	default:
		return obj.Inspect()
	}
}

// objectToPrintString converts an object to its string representation for print function
func objectToPrintString(obj Object) string {
	if obj == nil {
		return ""
	}

	switch obj := obj.(type) {
	case *Integer:
		return strconv.FormatInt(obj.Value, 10)
	case *Float:
		return fmt.Sprintf("%g", obj.Value)
	case *Boolean:
		if obj.Value {
			return "true"
		}
		return "false"
	case *String:
		return obj.Value
	case *Array:
		// Arrays: recursively print each element without any separators
		var result strings.Builder
		for _, elem := range obj.Elements {
			result.WriteString(objectToPrintString(elem))
		}
		return result.String()
	case *Null:
		return ""
	default:
		return obj.Inspect()
	}
}

// ObjectToPrintString is the exported version for use outside the package
func ObjectToPrintString(obj Object) string {
	return objectToPrintString(obj)
}

// objectToDebugString converts an object to its debug string representation
func objectToDebugString(obj Object) string {
	switch obj := obj.(type) {
	case *Integer:
		return strconv.FormatInt(obj.Value, 10)
	case *Float:
		return fmt.Sprintf("%g", obj.Value)
	case *Boolean:
		if obj.Value {
			return "true"
		}
		return "false"
	case *String:
		// Strings are wrapped in quotes for debug output
		return fmt.Sprintf("\"%s\"", obj.Value)
	case *Array:
		// Arrays: recursively debug print each element with separators, wrapped in brackets
		var result strings.Builder
		result.WriteString("[")
		for i, elem := range obj.Elements {
			if i > 0 {
				result.WriteString(", ")
			}
			result.WriteString(objectToDebugString(elem))
		}
		result.WriteString("]")
		return result.String()
	case *Null:
		return "null"
	default:
		return obj.Inspect()
	}
}

// evalConcatExpression handles the ++ operator for array concatenation
func evalConcatExpression(left, right Object) Object {
	// Handle dictionary concatenation
	if left.Type() == DICTIONARY_OBJ && right.Type() == DICTIONARY_OBJ {
		leftDict := left.(*Dictionary)
		rightDict := right.(*Dictionary)

		// Create new dictionary with merged pairs
		merged := &Dictionary{
			Pairs: make(map[string]ast.Expression),
			Env:   leftDict.Env, // Use left dict's environment
		}

		// Copy left dictionary pairs
		for k, v := range leftDict.Pairs {
			merged.Pairs[k] = v
		}

		// Copy right dictionary pairs (overwrites left if keys match)
		for k, v := range rightDict.Pairs {
			merged.Pairs[k] = v
		}

		return merged
	}

	// Convert single values to arrays
	var leftElements, rightElements []Object

	switch l := left.(type) {
	case *Array:
		leftElements = l.Elements
	default:
		leftElements = []Object{left}
	}

	switch r := right.(type) {
	case *Array:
		rightElements = r.Elements
	default:
		rightElements = []Object{right}
	}

	// Concatenate the arrays
	result := make([]Object, 0, len(leftElements)+len(rightElements))
	result = append(result, leftElements...)
	result = append(result, rightElements...)

	return &Array{Elements: result}
}

// evalIndexExpression handles array and string indexing
func evalIndexExpression(left, index Object) Object {
	switch {
	case left.Type() == ARRAY_OBJ && index.Type() == INTEGER_OBJ:
		return evalArrayIndexExpression(left, index)
	case left.Type() == STRING_OBJ && index.Type() == INTEGER_OBJ:
		return evalStringIndexExpression(left, index)
	case left.Type() == DICTIONARY_OBJ && index.Type() == STRING_OBJ:
		return evalDictionaryIndexExpression(left, index)
	default:
		return newError("index operator not supported: %s[%s]", left.Type(), index.Type())
	}
}

// evalArrayIndexExpression handles array indexing with support for negative indices
func evalArrayIndexExpression(array, index Object) Object {
	arrayObject := array.(*Array)
	idx := index.(*Integer).Value
	max := int64(len(arrayObject.Elements))

	// Handle negative indices
	if idx < 0 {
		idx = max + idx
	}

	if idx < 0 || idx >= max {
		return newError("index out of range: %d", index.(*Integer).Value)
	}

	return arrayObject.Elements[idx]
}

// evalStringIndexExpression handles string indexing with support for negative indices
func evalStringIndexExpression(str, index Object) Object {
	stringObject := str.(*String)
	idx := index.(*Integer).Value
	max := int64(len(stringObject.Value))

	// Handle negative indices
	if idx < 0 {
		idx = max + idx
	}

	if idx < 0 || idx >= max {
		return newError("index out of range: %d", index.(*Integer).Value)
	}

	return &String{Value: string(stringObject.Value[idx])}
}

// evalSliceExpression handles array and string slicing
func evalSliceExpression(left, start, end Object) Object {
	switch left.Type() {
	case ARRAY_OBJ:
		return evalArraySliceExpression(left, start, end)
	case STRING_OBJ:
		return evalStringSliceExpression(left, start, end)
	default:
		return newError("slice operator not supported: %s", left.Type())
	}
}

// evalArraySliceExpression handles array slicing
func evalArraySliceExpression(array, start, end Object) Object {
	arrayObject := array.(*Array)
	max := int64(len(arrayObject.Elements))

	var startIdx, endIdx int64

	// Determine start index
	if start == nil {
		startIdx = 0
	} else if start.Type() == INTEGER_OBJ {
		startIdx = start.(*Integer).Value
		if startIdx < 0 {
			startIdx = max + startIdx
		}
	} else {
		return newError("slice start index must be an integer, got %s", start.Type())
	}

	// Determine end index
	if end == nil {
		endIdx = max
	} else if end.Type() == INTEGER_OBJ {
		endIdx = end.(*Integer).Value
		if endIdx < 0 {
			endIdx = max + endIdx
		}
	} else {
		return newError("slice end index must be an integer, got %s", end.Type())
	}

	// Validate indices
	if startIdx < 0 || startIdx > max {
		return newError("slice start index out of range: %d", startIdx)
	}
	if endIdx < 0 || endIdx > max {
		return newError("slice end index out of range: %d", endIdx)
	}
	if startIdx > endIdx {
		return newError("slice start index %d is greater than end index %d", startIdx, endIdx)
	}

	// Create the slice
	return &Array{Elements: arrayObject.Elements[startIdx:endIdx]}
}

// evalStringSliceExpression handles string slicing
func evalStringSliceExpression(str, start, end Object) Object {
	stringObject := str.(*String)
	max := int64(len(stringObject.Value))

	var startIdx, endIdx int64

	// Determine start index
	if start == nil {
		startIdx = 0
	} else if start.Type() == INTEGER_OBJ {
		startIdx = start.(*Integer).Value
		if startIdx < 0 {
			startIdx = max + startIdx
		}
	} else {
		return newError("slice start index must be an integer, got %s", start.Type())
	}

	// Determine end index
	if end == nil {
		endIdx = max
	} else if end.Type() == INTEGER_OBJ {
		endIdx = end.(*Integer).Value
		if endIdx < 0 {
			endIdx = max + endIdx
		}
	} else {
		return newError("slice end index must be an integer, got %s", end.Type())
	}

	// Validate indices
	if startIdx < 0 || startIdx > max {
		return newError("slice start index out of range: %d", startIdx)
	}
	if endIdx < 0 || endIdx > max {
		return newError("slice end index out of range: %d", endIdx)
	}
	if startIdx > endIdx {
		return newError("slice start index %d is greater than end index %d", startIdx, endIdx)
	}

	// Create the slice
	return &String{Value: stringObject.Value[startIdx:endIdx]}
}

// evalDictionaryLiteral evaluates dictionary literals
func evalDictionaryLiteral(node *ast.DictionaryLiteral, env *Environment) Object {
	dict := &Dictionary{
		Pairs: node.Pairs,
		Env:   env,
	}
	return dict
}

// evalDotExpression evaluates dot notation access (dict.key)
func evalDotExpression(node *ast.DotExpression, env *Environment) Object {
	left := Eval(node.Left, env)
	if isError(left) {
		return left
	}

	dict, ok := left.(*Dictionary)
	if !ok {
		return newError("dot notation can only be used on dictionaries, got %s", left.Type())
	}

	// Get the expression from the dictionary
	expr, ok := dict.Pairs[node.Key]
	if !ok {
		return NULL
	}

	// Create a new environment with 'this' bound to the dictionary
	dictEnv := NewEnclosedEnvironment(dict.Env)
	dictEnv.Set("this", dict)

	// Evaluate the expression in the dictionary's environment
	return Eval(expr, dictEnv)
}

// evalDeleteStatement evaluates delete statements
func evalDeleteStatement(node *ast.DeleteStatement, env *Environment) Object {
	// The target must be a dot expression or index expression
	switch target := node.Target.(type) {
	case *ast.DotExpression:
		// Get the dictionary
		left := Eval(target.Left, env)
		if isError(left) {
			return left
		}

		dict, ok := left.(*Dictionary)
		if !ok {
			return newError("can only delete from dictionaries, got %s", left.Type())
		}

		// Delete the key
		delete(dict.Pairs, target.Key)
		return NULL

	case *ast.IndexExpression:
		// Get the dictionary
		left := Eval(target.Left, env)
		if isError(left) {
			return left
		}

		dict, ok := left.(*Dictionary)
		if !ok {
			return newError("can only delete from dictionaries, got %s", left.Type())
		}

		// Get the key
		index := Eval(target.Index, env)
		if isError(index) {
			return index
		}

		keyStr, ok := index.(*String)
		if !ok {
			return newError("dictionary key must be a string, got %s", index.Type())
		}

		// Delete the key
		delete(dict.Pairs, keyStr.Value)
		return NULL

	default:
		return newError("invalid delete target")
	}
}

// evalDictionaryIndexExpression handles dictionary access via dict["key"]
func evalDictionaryIndexExpression(dict, index Object) Object {
	dictObject := dict.(*Dictionary)
	key := index.(*String).Value

	// Get the expression from the dictionary
	expr, ok := dictObject.Pairs[key]
	if !ok {
		return NULL
	}

	// Create a new environment with 'this' bound to the dictionary
	dictEnv := NewEnclosedEnvironment(dictObject.Env)
	dictEnv.Set("this", dictObject)

	// Evaluate the expression in the dictionary's environment
	return Eval(expr, dictEnv)
}
