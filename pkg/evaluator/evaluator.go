package evaluator

import (
	"fmt"
	"math"

	"pars/pkg/ast"
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

// Global constants
var (
	NULL  = &Null{}
	TRUE  = &Boolean{Value: true}
	FALSE = &Boolean{Value: false}
)

// Built-in functions
var builtins = map[string]*Builtin{
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
		env.Set(node.Name.Value, val)
		return val

	case *ast.AssignmentStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		env.Set(node.Name.Value, val)
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
		return &Integer{Value: leftVal / rightVal}
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
	val, ok := env.Get(node.Value)
	if !ok {
		if builtin, ok := builtins[node.Value]; ok {
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

func newError(format string, a ...interface{}) *Error {
	return &Error{Message: fmt.Sprintf(format, a...)}
}

func isError(obj Object) bool {
	if obj != nil {
		return obj.Type() == ERROR_OBJ
	}
	return false
}
