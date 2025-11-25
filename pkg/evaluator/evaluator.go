package evaluator

import (
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
	"unicode"

	"pars/pkg/ast"
	"pars/pkg/lexer"
	"pars/pkg/parser"
)

// ObjectType represents the type of objects in our language
type ObjectType string

const (
	INTEGER_OBJ  = "INTEGER"
	FLOAT_OBJ    = "FLOAT"
	BOOLEAN_OBJ  = "BOOLEAN"
	STRING_OBJ   = "STRING"
	NULL_OBJ     = "NULL"
	RETURN_OBJ   = "RETURN_VALUE"
	ERROR_OBJ    = "ERROR"
	FUNCTION_OBJ = "FUNCTION"
	BUILTIN_OBJ  = "BUILTIN"
	ARRAY_OBJ    = "ARRAY"
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

func (i *Integer) Inspect() string  { return fmt.Sprintf("%d", i.Value) }
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

func (b *Boolean) Inspect() string  { return fmt.Sprintf("%t", b.Value) }
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
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

func (f *Function) Type() ObjectType { return FUNCTION_OBJ }
func (f *Function) Inspect() string {
	return fmt.Sprintf("fn(%v) {\n%s\n}", f.Parameters, f.Body.String())
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

// Environment represents the environment for variable bindings
type Environment struct {
	store map[string]Object
	outer *Environment
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

				if len(fn.Parameters) != 1 {
					return newError("function passed to `map` must take exactly 1 parameter, got %d", len(fn.Parameters))
				}

				result := []Object{}
				for _, elem := range arr.Elements {
					// Apply function to each element by creating new environment
					extendedEnv := NewEnclosedEnvironment(fn.Env)
					extendedEnv.Set(fn.Parameters[0].Value, elem)

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
					if i > 0 {
						result.WriteString(", ")
					}
					result.WriteString(objectToDebugString(arg))
				}

				// Write immediately to stdout
				fmt.Fprintln(os.Stdout, result.String())

				// Return null
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

		// Handle destructuring assignment
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

		// Handle destructuring assignment
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
		params := node.Parameters
		body := node.Body
		return &Function{Parameters: params, Env: env, Body: body}

	case *ast.ArrayLiteral:
		elements := evalExpressions(node.Elements, env)
		if len(elements) == 1 && isError(elements[0]) {
			return elements[0]
		}
		return &Array{Elements: elements}

	case *ast.CallExpression:
		function := Eval(node.Function, env)
		if isError(function) {
			return function
		}
		args := evalExpressions(node.Arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}
		return applyFunction(function, args)

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

func extendFunctionEnv(fn *Function, args []Object) *Environment {
	env := NewEnclosedEnvironment(fn.Env)

	for paramIdx, param := range fn.Parameters {
		env.Set(param.Value, args[paramIdx])
	}

	return env
}

func unwrapReturnValue(obj Object) Object {
	if returnValue, ok := obj.(*ReturnValue); ok {
		return returnValue.Value
	}
	return obj
}

// evalForExpression evaluates for expressions
func evalForExpression(node *ast.ForExpression, env *Environment) Object {
	// Evaluate the array expression
	arrayObj := Eval(node.Array, env)
	if isError(arrayObj) {
		return arrayObj
	}

	// Convert to array (handle strings as rune arrays)
	var elements []Object
	switch arr := arrayObj.(type) {
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
		return newError("for expects an array or string, got %s", arrayObj.Type())
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
			if len(f.Parameters) != 1 {
				return newError("function passed to for must take exactly 1 parameter, got %d", len(f.Parameters))
			}

			// Create a new environment that extends the function's environment
			// This allows the loop body to access and modify parent scope variables
			extendedEnv := NewEnclosedEnvironment(f.Env)
			extendedEnv.Set(f.Parameters[0].Value, elem)

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

// evalTemplateLiteral evaluates a template literal with interpolation
func evalTemplateLiteral(node *ast.TemplateLiteral, env *Environment) Object {
	template := node.Value
	var result strings.Builder

	i := 0
	for i < len(template) {
		// Check for escaped dollar sign marker \0$
		if i < len(template)-2 && template[i] == '\\' && template[i+1] == '0' && template[i+2] == '$' {
			result.WriteByte('$')
			i += 3
			continue
		}

		// Look for ${
		if i < len(template)-1 && template[i] == '$' && template[i+1] == '{' {
			// Find the closing }
			i += 2 // skip ${
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
				return newError("unclosed ${ in template literal")
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

// objectToTemplateString converts an object to its string representation for template interpolation
func objectToTemplateString(obj Object) string {
	switch obj := obj.(type) {
	case *Integer:
		return fmt.Sprintf("%d", obj.Value)
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
	switch obj := obj.(type) {
	case *Integer:
		return fmt.Sprintf("%d", obj.Value)
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

// objectToDebugString converts an object to its debug string representation
func objectToDebugString(obj Object) string {
	switch obj := obj.(type) {
	case *Integer:
		return fmt.Sprintf("%d", obj.Value)
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
