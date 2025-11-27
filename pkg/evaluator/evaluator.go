package evaluator

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
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
	Line    int
	Column  int
}

func (e *Error) Type() ObjectType { return ERROR_OBJ }
func (e *Error) Inspect() string {
	if e.Line > 0 {
		return fmt.Sprintf("line %d, column %d: %s", e.Line, e.Column, e.Message)
	}
	return "ERROR: " + e.Message
}

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
	store       map[string]Object
	outer       *Environment
	Filename    string
	LastToken   *lexer.Token
	letBindings map[string]bool // tracks which variables were declared with 'let'
}

// NewEnvironment creates a new environment
func NewEnvironment() *Environment {
	s := make(map[string]Object)
	l := make(map[string]bool)
	return &Environment{store: s, outer: nil, letBindings: l}
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

// SetLet stores a value in the environment and marks it as a let binding
func (e *Environment) SetLet(name string, val Object) Object {
	e.store[name] = val
	e.letBindings[name] = true
	return val
}

// IsLetBinding checks if a variable was declared with let
func (e *Environment) IsLetBinding(name string) bool {
	// Check current environment
	if e.letBindings[name] {
		return true
	}
	// Don't check outer environments - each module has its own scope
	return false
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

// ModuleCache caches imported modules
type ModuleCache struct {
	modules map[string]*Dictionary // absolute path -> module dictionary
	loading map[string]bool        // tracks currently loading modules for cycle detection
}

// Global module cache
var moduleCache = &ModuleCache{
	modules: make(map[string]*Dictionary),
	loading: make(map[string]bool),
}

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

// timeToDict converts a Go time.Time to a Parsley Dictionary
func timeToDict(t time.Time, env *Environment) *Dictionary {
	pairs := make(map[string]ast.Expression)

	// Mark this as a datetime dictionary for special operator handling
	pairs["__type"] = &ast.StringLiteral{
		Token: lexer.Token{Type: lexer.STRING, Literal: "datetime"},
		Value: "datetime",
	}

	// Create integer literals for numeric values with proper tokens
	pairs["year"] = &ast.IntegerLiteral{
		Token: lexer.Token{Type: lexer.INT, Literal: fmt.Sprintf("%d", t.Year())},
		Value: int64(t.Year()),
	}
	pairs["month"] = &ast.IntegerLiteral{
		Token: lexer.Token{Type: lexer.INT, Literal: fmt.Sprintf("%d", t.Month())},
		Value: int64(t.Month()),
	}
	pairs["day"] = &ast.IntegerLiteral{
		Token: lexer.Token{Type: lexer.INT, Literal: fmt.Sprintf("%d", t.Day())},
		Value: int64(t.Day()),
	}
	pairs["hour"] = &ast.IntegerLiteral{
		Token: lexer.Token{Type: lexer.INT, Literal: fmt.Sprintf("%d", t.Hour())},
		Value: int64(t.Hour()),
	}
	pairs["minute"] = &ast.IntegerLiteral{
		Token: lexer.Token{Type: lexer.INT, Literal: fmt.Sprintf("%d", t.Minute())},
		Value: int64(t.Minute()),
	}
	pairs["second"] = &ast.IntegerLiteral{
		Token: lexer.Token{Type: lexer.INT, Literal: fmt.Sprintf("%d", t.Second())},
		Value: int64(t.Second()),
	}
	pairs["unix"] = &ast.IntegerLiteral{
		Token: lexer.Token{Type: lexer.INT, Literal: fmt.Sprintf("%d", t.Unix())},
		Value: t.Unix(),
	}

	// Create string literals for string values with proper tokens
	weekday := t.Weekday().String()
	pairs["weekday"] = &ast.StringLiteral{
		Token: lexer.Token{Type: lexer.STRING, Literal: weekday},
		Value: weekday,
	}
	iso := t.Format(time.RFC3339)
	pairs["iso"] = &ast.StringLiteral{
		Token: lexer.Token{Type: lexer.STRING, Literal: iso},
		Value: iso,
	}

	return &Dictionary{Pairs: pairs, Env: env}
}

// dictToTime converts a Parsley Dictionary to a Go time.Time
func dictToTime(dict *Dictionary, env *Environment) (time.Time, error) {
	// Evaluate the year field
	yearExpr, ok := dict.Pairs["year"]
	if !ok {
		return time.Time{}, fmt.Errorf("missing 'year' field")
	}
	yearObj := Eval(yearExpr, env)
	yearInt, ok := yearObj.(*Integer)
	if !ok {
		return time.Time{}, fmt.Errorf("'year' must be an integer")
	}

	// Evaluate the month field
	monthExpr, ok := dict.Pairs["month"]
	if !ok {
		return time.Time{}, fmt.Errorf("missing 'month' field")
	}
	monthObj := Eval(monthExpr, env)
	monthInt, ok := monthObj.(*Integer)
	if !ok {
		return time.Time{}, fmt.Errorf("'month' must be an integer")
	}

	// Evaluate the day field
	dayExpr, ok := dict.Pairs["day"]
	if !ok {
		return time.Time{}, fmt.Errorf("missing 'day' field")
	}
	dayObj := Eval(dayExpr, env)
	dayInt, ok := dayObj.(*Integer)
	if !ok {
		return time.Time{}, fmt.Errorf("'day' must be an integer")
	}

	// Hour, minute, second are optional (default to 0)
	var hour, minute, second int64

	if hExpr, ok := dict.Pairs["hour"]; ok {
		hObj := Eval(hExpr, env)
		if hInt, ok := hObj.(*Integer); ok {
			hour = hInt.Value
		}
	}

	if mExpr, ok := dict.Pairs["minute"]; ok {
		mObj := Eval(mExpr, env)
		if mInt, ok := mObj.(*Integer); ok {
			minute = mInt.Value
		}
	}

	if sExpr, ok := dict.Pairs["second"]; ok {
		sObj := Eval(sExpr, env)
		if sInt, ok := sObj.(*Integer); ok {
			second = sInt.Value
		}
	}

	return time.Date(
		int(yearInt.Value),
		time.Month(monthInt.Value),
		int(dayInt.Value),
		int(hour),
		int(minute),
		int(second),
		0,
		time.UTC,
	), nil
}

// isDatetimeDict checks if a dictionary is a datetime by looking for __type field
func isDatetimeDict(dict *Dictionary) bool {
	if typeExpr, ok := dict.Pairs["__type"]; ok {
		if strLit, ok := typeExpr.(*ast.StringLiteral); ok {
			return strLit.Value == "datetime"
		}
	}
	return false
}

// isDurationDict checks if a dictionary is a duration by looking for __type field
func isDurationDict(dict *Dictionary) bool {
	if typeExpr, ok := dict.Pairs["__type"]; ok {
		if strLit, ok := typeExpr.(*ast.StringLiteral); ok {
			return strLit.Value == "duration"
		}
	}
	return false
}

// getDurationComponents extracts months and seconds from a duration dictionary
func getDurationComponents(dict *Dictionary, env *Environment) (int64, int64, error) {
	monthsExpr, ok := dict.Pairs["months"]
	if !ok {
		return 0, 0, fmt.Errorf("duration dictionary missing months field")
	}
	monthsObj := Eval(monthsExpr, env)
	monthsInt, ok := monthsObj.(*Integer)
	if !ok {
		return 0, 0, fmt.Errorf("months must be an integer")
	}

	secondsExpr, ok := dict.Pairs["seconds"]
	if !ok {
		return 0, 0, fmt.Errorf("duration dictionary missing seconds field")
	}
	secondsObj := Eval(secondsExpr, env)
	secondsInt, ok := secondsObj.(*Integer)
	if !ok {
		return 0, 0, fmt.Errorf("seconds must be an integer")
	}

	return monthsInt.Value, secondsInt.Value, nil
}

// getDatetimeUnix extracts the unix timestamp from a datetime dictionary
func getDatetimeUnix(dict *Dictionary, env *Environment) (int64, error) {
	unixExpr, ok := dict.Pairs["unix"]
	if !ok {
		return 0, fmt.Errorf("datetime dictionary missing unix field")
	}
	unixObj := Eval(unixExpr, env)
	unixInt, ok := unixObj.(*Integer)
	if !ok {
		return 0, fmt.Errorf("unix field is not an integer")
	}
	return unixInt.Value, nil
}

// applyDelta applies time deltas to a time.Time
func applyDelta(t time.Time, delta *Dictionary, env *Environment) time.Time {
	// Apply date-based deltas first (years, months, days)
	if yearsExpr, ok := delta.Pairs["years"]; ok {
		yearsObj := Eval(yearsExpr, env)
		if yearsInt, ok := yearsObj.(*Integer); ok {
			t = t.AddDate(int(yearsInt.Value), 0, 0)
		}
	}

	if monthsExpr, ok := delta.Pairs["months"]; ok {
		monthsObj := Eval(monthsExpr, env)
		if monthsInt, ok := monthsObj.(*Integer); ok {
			t = t.AddDate(0, int(monthsInt.Value), 0)
		}
	}

	if daysExpr, ok := delta.Pairs["days"]; ok {
		daysObj := Eval(daysExpr, env)
		if daysInt, ok := daysObj.(*Integer); ok {
			t = t.AddDate(0, 0, int(daysInt.Value))
		}
	}

	// Apply time-based deltas (hours, minutes, seconds)
	if hoursExpr, ok := delta.Pairs["hours"]; ok {
		hoursObj := Eval(hoursExpr, env)
		if hoursInt, ok := hoursObj.(*Integer); ok {
			t = t.Add(time.Duration(hoursInt.Value) * time.Hour)
		}
	}

	if minutesExpr, ok := delta.Pairs["minutes"]; ok {
		minutesObj := Eval(minutesExpr, env)
		if minutesInt, ok := minutesObj.(*Integer); ok {
			t = t.Add(time.Duration(minutesInt.Value) * time.Minute)
		}
	}

	if secondsExpr, ok := delta.Pairs["seconds"]; ok {
		secondsObj := Eval(secondsExpr, env)
		if secondsInt, ok := secondsObj.(*Integer); ok {
			t = t.Add(time.Duration(secondsInt.Value) * time.Second)
		}
	}

	return t
}

// evalRegexLiteral evaluates a regex literal and returns a Dictionary with __type: "regex"
func evalRegexLiteral(node *ast.RegexLiteral, env *Environment) Object {
	pairs := make(map[string]ast.Expression)

	// Mark this as a regex dictionary
	pairs["__type"] = &ast.StringLiteral{Value: "regex"}
	pairs["pattern"] = &ast.StringLiteral{Value: node.Pattern}
	pairs["flags"] = &ast.StringLiteral{Value: node.Flags}

	// Try to compile the regex to validate it
	_, err := compileRegex(node.Pattern, node.Flags)
	if err != nil {
		return newError("invalid regex pattern: %s", err.Error())
	}

	return &Dictionary{Pairs: pairs, Env: env}
}

// evalDatetimeLiteral evaluates a datetime literal like @2024-12-25T14:30:00Z
func evalDatetimeLiteral(node *ast.DatetimeLiteral, env *Environment) Object {
	// Parse the ISO-8601 datetime string
	var t time.Time
	var err error

	// Try parsing as RFC3339 first (most complete format with timezone)
	t, err = time.Parse(time.RFC3339, node.Value)
	if err != nil {
		// Try date-only format (2024-12-25) - interpret as UTC
		t, err = time.ParseInLocation("2006-01-02", node.Value, time.UTC)
		if err != nil {
			// Try datetime without timezone (2024-12-25T14:30:05) - interpret as UTC
			t, err = time.ParseInLocation("2006-01-02T15:04:05", node.Value, time.UTC)
			if err != nil {
				return newError("invalid datetime literal: %s", node.Value)
			}
		}
	}

	// Convert to dictionary using the same function as the time() builtin
	return timeToDict(t, env)
}

// evalDurationLiteral parses a duration literal like @2h30m, @7d, @1y6mo
func evalDurationLiteral(node *ast.DurationLiteral, env *Environment) Object {
	// Parse the duration string into months and seconds
	months, seconds, err := parseDurationString(node.Value)
	if err != nil {
		return newError("invalid duration literal: %s", err.Error())
	}

	return durationToDict(months, seconds, env)
}

// evalPathLiteral parses a path literal like @/usr/local/bin or @./config.json
func evalPathLiteral(node *ast.PathLiteral, env *Environment) Object {
	// Parse the path string into components
	components, isAbsolute := parsePathString(node.Value)

	// Create path dictionary
	return pathToDict(components, isAbsolute, env)
}

// evalUrlLiteral parses a URL literal like @https://example.com/api
func evalUrlLiteral(node *ast.UrlLiteral, env *Environment) Object {
	// Parse the URL string
	urlDict, err := parseUrlString(node.Value, env)
	if err != nil {
		return newError("invalid URL literal: %s", err.Error())
	}

	return urlDict
}

// parseDurationString parses a duration string like "2h30m" or "1y6mo" into months and seconds
// Returns (months, seconds, error)
func parseDurationString(s string) (int64, int64, error) {
	var months int64
	var seconds int64

	i := 0
	for i < len(s) {
		// Read number
		if !isDigit(rune(s[i])) {
			return 0, 0, fmt.Errorf("expected digit at position %d", i)
		}

		numStart := i
		for i < len(s) && isDigit(rune(s[i])) {
			i++
		}

		num, err := strconv.ParseInt(s[numStart:i], 10, 64)
		if err != nil {
			return 0, 0, err
		}

		// Read unit
		if i >= len(s) {
			return 0, 0, fmt.Errorf("missing unit after number at position %d", i)
		}

		var unit string
		// Check for "mo" (months)
		if i+1 < len(s) && s[i:i+2] == "mo" {
			unit = "mo"
			i += 2
		} else {
			// Single letter unit
			unit = string(s[i])
			i++
		}

		// Convert to months or seconds
		switch unit {
		case "y": // years = 12 months
			months += num * 12
		case "mo": // months
			months += num
		case "w": // weeks = 7 days = 7 * 24 * 60 * 60 seconds
			seconds += num * 7 * 24 * 60 * 60
		case "d": // days = 24 * 60 * 60 seconds
			seconds += num * 24 * 60 * 60
		case "h": // hours = 60 * 60 seconds
			seconds += num * 60 * 60
		case "m": // minutes = 60 seconds
			seconds += num * 60
		case "s": // seconds
			seconds += num
		default:
			return 0, 0, fmt.Errorf("unknown unit: %s", unit)
		}
	}

	return months, seconds, nil
}

// durationToDict converts months and seconds into a duration dictionary
func durationToDict(months, seconds int64, env *Environment) *Dictionary {
	dict := &Dictionary{Pairs: make(map[string]ast.Expression)}

	// Add __type field
	dict.Pairs["__type"] = &ast.StringLiteral{
		Token: lexer.Token{Type: lexer.STRING, Literal: "duration"},
		Value: "duration",
	}

	// Add months field
	dict.Pairs["months"] = &ast.IntegerLiteral{
		Token: lexer.Token{Type: lexer.INT, Literal: fmt.Sprintf("%d", months)},
		Value: months,
	}

	// Add seconds field
	dict.Pairs["seconds"] = &ast.IntegerLiteral{
		Token: lexer.Token{Type: lexer.INT, Literal: fmt.Sprintf("%d", seconds)},
		Value: seconds,
	}

	// Add totalSeconds field (only present if no months)
	if months == 0 {
		dict.Pairs["totalSeconds"] = &ast.IntegerLiteral{
			Token: lexer.Token{Type: lexer.INT, Literal: fmt.Sprintf("%d", seconds)},
			Value: seconds,
		}
	}

	return dict
}

// isDigit checks if a rune is a digit
func isDigit(ch rune) bool {
	return ch >= '0' && ch <= '9'
}

// isRegexDict checks if a dictionary is a regex by looking for __type field
func isRegexDict(dict *Dictionary) bool {
	if typeExpr, ok := dict.Pairs["__type"]; ok {
		if strLit, ok := typeExpr.(*ast.StringLiteral); ok {
			return strLit.Value == "regex"
		}
	}
	return false
}

// isPathDict checks if a dictionary is a path by looking for __type field
func isPathDict(dict *Dictionary) bool {
	if typeExpr, ok := dict.Pairs["__type"]; ok {
		if strLit, ok := typeExpr.(*ast.StringLiteral); ok {
			return strLit.Value == "path"
		}
	}
	return false
}

// isUrlDict checks if a dictionary is a URL by looking for __type field
func isUrlDict(dict *Dictionary) bool {
	if typeExpr, ok := dict.Pairs["__type"]; ok {
		if strLit, ok := typeExpr.(*ast.StringLiteral); ok {
			return strLit.Value == "url"
		}
	}
	return false
}

// compileRegex compiles a regex pattern with optional flags
// Go's regexp doesn't support all Perl flags, so we map what we can
func compileRegex(pattern, flags string) (*regexp.Regexp, error) {
	// Process flags - Go regexp supports (?flags) syntax
	prefix := ""
	for _, flag := range flags {
		switch flag {
		case 'i': // case-insensitive
			prefix += "(?i)"
		case 'm': // multi-line (^ and $ match line boundaries)
			prefix += "(?m)"
		case 's': // dot matches newline
			prefix += "(?s)"
			// 'g' (global) is handled by match operator, not compilation
			// Other flags like 'x' (verbose) could be added
		}
	}

	fullPattern := prefix + pattern
	return regexp.Compile(fullPattern)
}

// evalMatchExpression handles string ~ regex matching
// Returns an array of matches (with captures) or null if no match
func evalMatchExpression(tok lexer.Token, text string, regexDict *Dictionary, env *Environment) Object {
	// Extract pattern and flags from regex dictionary
	patternExpr, ok := regexDict.Pairs["pattern"]
	if !ok {
		return newErrorWithPos(tok, "regex dictionary missing pattern field")
	}
	patternObj := Eval(patternExpr, env)
	patternStr, ok := patternObj.(*String)
	if !ok {
		return newErrorWithPos(tok, "regex pattern must be a string")
	}

	flagsExpr, ok := regexDict.Pairs["flags"]
	var flags string
	if ok {
		flagsObj := Eval(flagsExpr, env)
		if flagsStr, ok := flagsObj.(*String); ok {
			flags = flagsStr.Value
		}
	}

	// Compile the regex
	re, err := compileRegex(patternStr.Value, flags)
	if err != nil {
		return newErrorWithPos(tok, "invalid regex: %s", err.Error())
	}

	// Find matches
	matches := re.FindStringSubmatch(text)
	if matches == nil {
		return NULL // No match - returns null (falsy)
	}

	// Convert matches to array of strings
	elements := make([]Object, len(matches))
	for i, match := range matches {
		elements[i] = &String{Value: match}
	}

	return &Array{Elements: elements}
}

// parsePathString parses a file path string into components
// Returns components array and whether path is absolute
func parsePathString(pathStr string) ([]string, bool) {
	if pathStr == "" {
		return []string{}, false
	}

	// Detect absolute vs relative
	isAbsolute := false
	if pathStr[0] == '/' {
		isAbsolute = true
	} else if len(pathStr) >= 2 && pathStr[1] == ':' {
		// Windows drive letter (C:, D:, etc.)
		isAbsolute = true
	}

	// Split on forward slashes (handle both Unix and Windows)
	pathStr = strings.ReplaceAll(pathStr, "\\", "/")
	parts := strings.Split(pathStr, "/")

	// Filter empty strings except for leading slash
	components := []string{}
	for i, part := range parts {
		if part == "" && i == 0 && isAbsolute {
			// Keep leading empty string to indicate absolute path
			components = append(components, "")
		} else if part != "" {
			components = append(components, part)
		}
	}

	return components, isAbsolute
}

// pathToDict creates a path dictionary from components
func pathToDict(components []string, isAbsolute bool, env *Environment) *Dictionary {
	pairs := make(map[string]ast.Expression)

	// Add __type field
	pairs["__type"] = &ast.StringLiteral{
		Token: lexer.Token{Type: lexer.STRING, Literal: "path"},
		Value: "path",
	}

	// Add components as array literal
	componentExprs := make([]ast.Expression, len(components))
	for i, comp := range components {
		componentExprs[i] = &ast.StringLiteral{
			Token: lexer.Token{Type: lexer.STRING, Literal: comp},
			Value: comp,
		}
	}
	pairs["components"] = &ast.ArrayLiteral{
		Token:    lexer.Token{Type: lexer.LBRACKET, Literal: "["},
		Elements: componentExprs,
	}

	// Add absolute flag
	tokenType := lexer.FALSE
	tokenLiteral := "false"
	if isAbsolute {
		tokenType = lexer.TRUE
		tokenLiteral = "true"
	}
	pairs["absolute"] = &ast.Boolean{
		Token: lexer.Token{Type: tokenType, Literal: tokenLiteral},
		Value: isAbsolute,
	}

	return &Dictionary{Pairs: pairs, Env: env}
}

// parseUrlString parses a URL string into components
// Supports: scheme://[user:pass@]host[:port]/path?query#fragment
func parseUrlString(urlStr string, env *Environment) (*Dictionary, error) {
	// Simple URL parsing (not using net/url to keep it simple)
	pairs := make(map[string]ast.Expression)

	// Add __type field
	pairs["__type"] = &ast.StringLiteral{
		Token: lexer.Token{Type: lexer.STRING, Literal: "url"},
		Value: "url",
	}

	// Parse scheme
	schemeEnd := strings.Index(urlStr, "://")
	if schemeEnd == -1 {
		return nil, fmt.Errorf("invalid URL: missing scheme (expected scheme://...)")
	}
	scheme := urlStr[:schemeEnd]
	rest := urlStr[schemeEnd+3:]

	pairs["scheme"] = &ast.StringLiteral{
		Token: lexer.Token{Type: lexer.STRING, Literal: scheme},
		Value: scheme,
	}

	// Parse fragment (if present)
	var fragment string
	if fragIdx := strings.Index(rest, "#"); fragIdx != -1 {
		fragment = rest[fragIdx+1:]
		rest = rest[:fragIdx]
	}

	// Parse query (if present)
	queryPairs := make(map[string]ast.Expression)
	if queryIdx := strings.Index(rest, "?"); queryIdx != -1 {
		queryStr := rest[queryIdx+1:]
		rest = rest[:queryIdx]

		// Parse query parameters
		for _, param := range strings.Split(queryStr, "&") {
			if param == "" {
				continue
			}
			parts := strings.SplitN(param, "=", 2)
			key := parts[0]
			value := ""
			if len(parts) > 1 {
				value = parts[1]
			}
			queryPairs[key] = &ast.StringLiteral{
				Token: lexer.Token{Type: lexer.STRING, Literal: value},
				Value: value,
			}
		}
	}
	pairs["query"] = &ast.DictionaryLiteral{
		Token: lexer.Token{Type: lexer.LBRACE, Literal: "{"},
		Pairs: queryPairs,
	}

	// Parse path (if present)
	pathComponents := []string{}
	var pathStr string
	if pathIdx := strings.Index(rest, "/"); pathIdx != -1 {
		pathStr = rest[pathIdx:]
		rest = rest[:pathIdx]
		pathComponents, _ = parsePathString(pathStr)
	}

	pathExprs := make([]ast.Expression, len(pathComponents))
	for i, comp := range pathComponents {
		pathExprs[i] = &ast.StringLiteral{
			Token: lexer.Token{Type: lexer.STRING, Literal: comp},
			Value: comp,
		}
	}
	pairs["path"] = &ast.ArrayLiteral{
		Token:    lexer.Token{Type: lexer.LBRACKET, Literal: "["},
		Elements: pathExprs,
	}

	// Parse authority (user:pass@host:port)
	var username, password, host string
	var port int64 = 0

	// Check for userinfo (user:pass@)
	if atIdx := strings.Index(rest, "@"); atIdx != -1 {
		userinfo := rest[:atIdx]
		rest = rest[atIdx+1:]

		if colonIdx := strings.Index(userinfo, ":"); colonIdx != -1 {
			username = userinfo[:colonIdx]
			password = userinfo[colonIdx+1:]
		} else {
			username = userinfo
		}
	}

	// Parse host:port
	if colonIdx := strings.Index(rest, ":"); colonIdx != -1 {
		host = rest[:colonIdx]
		portStr := rest[colonIdx+1:]
		if p, err := strconv.ParseInt(portStr, 10, 64); err == nil {
			port = p
		}
	} else {
		host = rest
	}

	pairs["host"] = &ast.StringLiteral{
		Token: lexer.Token{Type: lexer.STRING, Literal: host},
		Value: host,
	}

	pairs["port"] = &ast.IntegerLiteral{
		Token: lexer.Token{Type: lexer.INT, Literal: fmt.Sprintf("%d", port)},
		Value: port,
	}

	if username != "" {
		pairs["username"] = &ast.StringLiteral{
			Token: lexer.Token{Type: lexer.STRING, Literal: username},
			Value: username,
		}
	} else {
		pairs["username"] = &ast.Identifier{
			Token: lexer.Token{Type: lexer.IDENT, Literal: "null"},
			Value: "null",
		}
	}

	if password != "" {
		pairs["password"] = &ast.StringLiteral{
			Token: lexer.Token{Type: lexer.STRING, Literal: password},
			Value: password,
		}
	} else {
		pairs["password"] = &ast.Identifier{
			Token: lexer.Token{Type: lexer.IDENT, Literal: "null"},
			Value: "null",
		}
	}

	if fragment != "" {
		pairs["fragment"] = &ast.StringLiteral{
			Token: lexer.Token{Type: lexer.STRING, Literal: fragment},
			Value: fragment,
		}
	} else {
		pairs["fragment"] = &ast.Identifier{
			Token: lexer.Token{Type: lexer.IDENT, Literal: "null"},
			Value: "null",
		}
	}

	return &Dictionary{Pairs: pairs, Env: env}, nil
}

// evalPathComputedProperty returns computed properties for path dictionaries
// Returns nil if the property doesn't exist
func evalPathComputedProperty(dict *Dictionary, key string, env *Environment) Object {
	switch key {
	case "basename":
		// Get last component
		componentsExpr, ok := dict.Pairs["components"]
		if !ok {
			return NULL
		}
		componentsObj := Eval(componentsExpr, env)
		arr, ok := componentsObj.(*Array)
		if !ok || len(arr.Elements) == 0 {
			return NULL
		}
		return arr.Elements[len(arr.Elements)-1]

	case "dirname", "parent":
		// Get all but last component, return as path dict
		componentsExpr, ok := dict.Pairs["components"]
		if !ok {
			return NULL
		}
		componentsObj := Eval(componentsExpr, env)
		arr, ok := componentsObj.(*Array)
		if !ok || len(arr.Elements) == 0 {
			return NULL
		}

		// Get absolute flag
		absoluteExpr, ok := dict.Pairs["absolute"]
		isAbsolute := false
		if ok {
			absoluteObj := Eval(absoluteExpr, env)
			if b, ok := absoluteObj.(*Boolean); ok {
				isAbsolute = b.Value
			}
		}

		// Create new components array (all but last)
		parentComponents := []string{}
		for i := 0; i < len(arr.Elements)-1; i++ {
			if str, ok := arr.Elements[i].(*String); ok {
				parentComponents = append(parentComponents, str.Value)
			}
		}

		return pathToDict(parentComponents, isAbsolute, env)

	case "extension", "ext":
		// Get extension from basename
		componentsExpr, ok := dict.Pairs["components"]
		if !ok {
			return NULL
		}
		componentsObj := Eval(componentsExpr, env)
		arr, ok := componentsObj.(*Array)
		if !ok || len(arr.Elements) == 0 {
			return NULL
		}
		basename, ok := arr.Elements[len(arr.Elements)-1].(*String)
		if !ok {
			return NULL
		}

		// Find last dot
		lastDot := strings.LastIndex(basename.Value, ".")
		if lastDot == -1 || lastDot == 0 {
			return &String{Value: ""}
		}
		return &String{Value: basename.Value[lastDot+1:]}

	case "stem":
		// Get filename without extension
		componentsExpr, ok := dict.Pairs["components"]
		if !ok {
			return NULL
		}
		componentsObj := Eval(componentsExpr, env)
		arr, ok := componentsObj.(*Array)
		if !ok || len(arr.Elements) == 0 {
			return NULL
		}
		basename, ok := arr.Elements[len(arr.Elements)-1].(*String)
		if !ok {
			return NULL
		}

		// Find last dot
		lastDot := strings.LastIndex(basename.Value, ".")
		if lastDot == -1 || lastDot == 0 {
			return basename
		}
		return &String{Value: basename.Value[:lastDot]}

	case "name":
		// Alias for basename
		return evalPathComputedProperty(dict, "basename", env)

	case "suffix":
		// Alias for extension
		return evalPathComputedProperty(dict, "extension", env)

	case "suffixes":
		// Get all extensions as array (e.g., ["tar", "gz"] from file.tar.gz)
		componentsExpr, ok := dict.Pairs["components"]
		if !ok {
			return NULL
		}
		componentsObj := Eval(componentsExpr, env)
		arr, ok := componentsObj.(*Array)
		if !ok || len(arr.Elements) == 0 {
			return &Array{Elements: []Object{}}
		}
		basename, ok := arr.Elements[len(arr.Elements)-1].(*String)
		if !ok {
			return &Array{Elements: []Object{}}
		}

		// Find all dots and extract suffixes
		var suffixes []Object
		parts := strings.Split(basename.Value, ".")
		if len(parts) > 1 {
			// Skip the first part (filename), collect rest as suffixes
			for i := 1; i < len(parts); i++ {
				if parts[i] != "" {
					suffixes = append(suffixes, &String{Value: parts[i]})
				}
			}
		}
		return &Array{Elements: suffixes}

	case "parts":
		// Alias for components
		componentsExpr, ok := dict.Pairs["components"]
		if !ok {
			return NULL
		}
		return Eval(componentsExpr, env)

	case "isAbsolute":
		// Boolean indicating if path is absolute
		absoluteExpr, ok := dict.Pairs["absolute"]
		if !ok {
			return FALSE
		}
		return Eval(absoluteExpr, env)

	case "isRelative":
		// Boolean indicating if path is relative (opposite of absolute)
		absoluteExpr, ok := dict.Pairs["absolute"]
		if !ok {
			return TRUE
		}
		absoluteObj := Eval(absoluteExpr, env)
		if b, ok := absoluteObj.(*Boolean); ok {
			return nativeBoolToParsBoolean(!b.Value)
		}
		return TRUE
	}

	return nil // Property doesn't exist
}

// evalUrlComputedProperty returns computed properties for URL dictionaries
// Returns nil if the property doesn't exist
func evalUrlComputedProperty(dict *Dictionary, key string, env *Environment) Object {
	switch key {
	case "origin":
		// scheme://host[:port]
		var result strings.Builder

		if schemeExpr, ok := dict.Pairs["scheme"]; ok {
			schemeObj := Eval(schemeExpr, env)
			if str, ok := schemeObj.(*String); ok {
				result.WriteString(str.Value)
				result.WriteString("://")
			}
		}

		if hostExpr, ok := dict.Pairs["host"]; ok {
			hostObj := Eval(hostExpr, env)
			if str, ok := hostObj.(*String); ok {
				result.WriteString(str.Value)
			}
		}

		if portExpr, ok := dict.Pairs["port"]; ok {
			portObj := Eval(portExpr, env)
			if i, ok := portObj.(*Integer); ok && i.Value != 0 {
				result.WriteString(":")
				result.WriteString(strconv.FormatInt(i.Value, 10))
			}
		}

		return &String{Value: result.String()}

	case "pathname":
		// Just the path part as a string
		if pathExpr, ok := dict.Pairs["path"]; ok {
			pathObj := Eval(pathExpr, env)
			if arr, ok := pathObj.(*Array); ok {
				var result strings.Builder
				hasLeadingSlash := false
				for i, elem := range arr.Elements {
					if str, ok := elem.(*String); ok {
						if i == 0 && str.Value == "" {
							// Leading empty string means absolute path
							result.WriteString("/")
							hasLeadingSlash = true
						} else if str.Value != "" {
							// Add slash before element (but not if we just added leading slash)
							if i > 0 && !(i == 1 && hasLeadingSlash) {
								result.WriteString("/")
							}
							result.WriteString(str.Value)
						}
					}
				}
				return &String{Value: result.String()}
			}
		}
		return &String{Value: ""}

	case "hostname":
		// Alias for host
		if hostExpr, ok := dict.Pairs["host"]; ok {
			return Eval(hostExpr, env)
		}
		return &String{Value: ""}

	case "protocol":
		// Scheme with colon suffix (e.g., "https:")
		if schemeExpr, ok := dict.Pairs["scheme"]; ok {
			schemeObj := Eval(schemeExpr, env)
			if str, ok := schemeObj.(*String); ok {
				return &String{Value: str.Value + ":"}
			}
		}
		return &String{Value: ""}

	case "search":
		// Query string with ? prefix (e.g., "?key=value&foo=bar")
		if queryExpr, ok := dict.Pairs["query"]; ok {
			queryObj := Eval(queryExpr, env)
			if queryDict, ok := queryObj.(*Dictionary); ok {
				if len(queryDict.Pairs) == 0 {
					return &String{Value: ""}
				}
				var result strings.Builder
				result.WriteString("?")
				first := true
				for key, expr := range queryDict.Pairs {
					val := Eval(expr, env)
					if str, ok := val.(*String); ok {
						if !first {
							result.WriteString("&")
						}
						result.WriteString(key)
						result.WriteString("=")
						result.WriteString(str.Value)
						first = false
					}
				}
				return &String{Value: result.String()}
			}
		}
		return &String{Value: ""}

	case "href":
		// Full URL as string (alias for toString)
		return &String{Value: urlDictToString(dict)}
	}

	return nil // Property doesn't exist
}

// pathDictToString converts a path dictionary back to a string
func pathDictToString(dict *Dictionary) string {
	// Get components array
	componentsExpr, ok := dict.Pairs["components"]
	if !ok {
		return ""
	}

	// Evaluate the array expression
	componentsObj := Eval(componentsExpr, dict.Env)
	arr, ok := componentsObj.(*Array)
	if !ok {
		return ""
	}

	// Check if absolute
	absoluteExpr, ok := dict.Pairs["absolute"]
	isAbsolute := false
	if ok {
		absoluteObj := Eval(absoluteExpr, dict.Env)
		if b, ok := absoluteObj.(*Boolean); ok {
			isAbsolute = b.Value
		}
	}

	// Build path string
	var result strings.Builder
	for i, elem := range arr.Elements {
		if str, ok := elem.(*String); ok {
			if str.Value == "" && i == 0 && isAbsolute {
				// Leading empty string means absolute path
				result.WriteString("/")
			} else {
				if i > 0 && (i > 1 || !isAbsolute) {
					result.WriteString("/")
				}
				result.WriteString(str.Value)
			}
		}
	}

	return result.String()
}

// urlDictToString converts a URL dictionary back to a string
func urlDictToString(dict *Dictionary) string {
	var result strings.Builder

	// Scheme
	if schemeExpr, ok := dict.Pairs["scheme"]; ok {
		schemeObj := Eval(schemeExpr, dict.Env)
		if str, ok := schemeObj.(*String); ok {
			result.WriteString(str.Value)
			result.WriteString("://")
		}
	}

	// Username and password
	if usernameExpr, ok := dict.Pairs["username"]; ok {
		usernameObj := Eval(usernameExpr, dict.Env)
		if str, ok := usernameObj.(*String); ok && str.Value != "" {
			result.WriteString(str.Value)

			if passwordExpr, ok := dict.Pairs["password"]; ok {
				passwordObj := Eval(passwordExpr, dict.Env)
				if pstr, ok := passwordObj.(*String); ok && pstr.Value != "" {
					result.WriteString(":")
					result.WriteString(pstr.Value)
				}
			}
			result.WriteString("@")
		}
	}

	// Host
	if hostExpr, ok := dict.Pairs["host"]; ok {
		hostObj := Eval(hostExpr, dict.Env)
		if str, ok := hostObj.(*String); ok {
			result.WriteString(str.Value)
		}
	}

	// Port (if non-zero)
	if portExpr, ok := dict.Pairs["port"]; ok {
		portObj := Eval(portExpr, dict.Env)
		if i, ok := portObj.(*Integer); ok && i.Value != 0 {
			result.WriteString(":")
			result.WriteString(strconv.FormatInt(i.Value, 10))
		}
	}

	// Path
	if pathExpr, ok := dict.Pairs["path"]; ok {
		pathObj := Eval(pathExpr, dict.Env)
		if arr, ok := pathObj.(*Array); ok && len(arr.Elements) > 0 {
			// Check if first element is empty string (indicates leading slash)
			startIdx := 0
			if str, ok := arr.Elements[0].(*String); ok && str.Value == "" {
				// Leading empty string means path starts with /
				result.WriteString("/")
				startIdx = 1
			} else if len(arr.Elements) > 0 {
				// No leading empty, but we have segments, so add /
				result.WriteString("/")
			}

			// Add remaining path segments
			for i := startIdx; i < len(arr.Elements); i++ {
				if str, ok := arr.Elements[i].(*String); ok && str.Value != "" {
					if i > startIdx {
						result.WriteString("/")
					}
					result.WriteString(str.Value)
				}
			}
		}
	}

	// Query
	if queryExpr, ok := dict.Pairs["query"]; ok {
		queryObj := Eval(queryExpr, dict.Env)
		if queryDict, ok := queryObj.(*Dictionary); ok && len(queryDict.Pairs) > 0 {
			result.WriteString("?")
			first := true
			for key, expr := range queryDict.Pairs {
				if !first {
					result.WriteString("&")
				}
				first = false
				result.WriteString(key)
				result.WriteString("=")
				valObj := Eval(expr, dict.Env)
				if str, ok := valObj.(*String); ok {
					result.WriteString(str.Value)
				}
			}
		}
	}

	// Fragment
	if fragmentExpr, ok := dict.Pairs["fragment"]; ok {
		fragmentObj := Eval(fragmentExpr, dict.Env)
		if str, ok := fragmentObj.(*String); ok && str.Value != "" {
			result.WriteString("#")
			result.WriteString(str.Value)
		}
	}

	return result.String()
}

// getBuiltins returns the map of built-in functions
func getBuiltins() map[string]*Builtin {
	return map[string]*Builtin{
		"import": {
			Fn: func(args ...Object) Object {
				// This is a placeholder - actual implementation happens in CallExpression
				// where we have access to the environment for path resolution
				return newError("import() requires environment context")
			},
		},
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
		"now": {
			Fn: func(args ...Object) Object {
				if len(args) != 0 {
					return newError("wrong number of arguments. got=%d, want=0", len(args))
				}
				// Get current environment from context (we'll pass it through the Builtin)
				// For now, create a new environment for the dictionary
				env := NewEnvironment()
				return timeToDict(time.Now(), env)
			},
		},
		"time": {
			Fn: func(args ...Object) Object {
				if len(args) < 1 || len(args) > 2 {
					return newError("wrong number of arguments. got=%d, want=1 or 2", len(args))
				}

				env := NewEnvironment()
				var t time.Time
				var err error

				switch arg := args[0].(type) {
				case *String:
					// Try parsing as ISO 8601 first, then fall back to date-only format
					t, err = time.Parse(time.RFC3339, arg.Value)
					if err != nil {
						t, err = time.Parse("2006-01-02", arg.Value)
					}
					if err != nil {
						t, err = time.Parse("2006-01-02T15:04:05", arg.Value)
					}
					if err != nil {
						return newError("invalid datetime string: %s", arg.Value)
					}
				case *Integer:
					// Unix timestamp
					t = time.Unix(arg.Value, 0).UTC()
				case *Dictionary:
					// From dictionary
					t, err = dictToTime(arg, env)
					if err != nil {
						return newError("invalid datetime dictionary: %s", err)
					}
				default:
					return newError("argument to `time` must be STRING, INTEGER, or DICTIONARY, got %s", args[0].Type())
				}

				// Apply delta if provided
				if len(args) == 2 {
					delta, ok := args[1].(*Dictionary)
					if !ok {
						return newError("second argument to `time` must be DICTIONARY, got %s", args[1].Type())
					}
					t = applyDelta(t, delta, env)
				}

				return timeToDict(t, env)
			},
		},
		"path": {
			Fn: func(args ...Object) Object {
				if len(args) != 1 {
					return newError("wrong number of arguments to `path`. got=%d, want=1", len(args))
				}

				str, ok := args[0].(*String)
				if !ok {
					return newError("argument to `path` must be a string, got %s", args[0].Type())
				}

				components, isAbsolute := parsePathString(str.Value)
				env := NewEnvironment()
				return pathToDict(components, isAbsolute, env)
			},
		},
		"url": {
			Fn: func(args ...Object) Object {
				if len(args) != 1 {
					return newError("wrong number of arguments to `url`. got=%d, want=1", len(args))
				}

				str, ok := args[0].(*String)
				if !ok {
					return newError("argument to `url` must be a string, got %s", args[0].Type())
				}

				env := NewEnvironment()
				urlDict, err := parseUrlString(str.Value, env)
				if err != nil {
					return newError("invalid URL: %s", err.Error())
				}

				return urlDict
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
		"regex": {
			Fn: func(args ...Object) Object {
				if len(args) < 1 || len(args) > 2 {
					return newError("wrong number of arguments to `regex`. got=%d, want=1 or 2", len(args))
				}

				pattern, ok := args[0].(*String)
				if !ok {
					return newError("first argument to `regex` must be a string, got %s", args[0].Type())
				}

				flags := ""
				if len(args) == 2 {
					flagsStr, ok := args[1].(*String)
					if !ok {
						return newError("second argument to `regex` must be a string, got %s", args[1].Type())
					}
					flags = flagsStr.Value
				}

				// Validate the regex
				_, err := compileRegex(pattern.Value, flags)
				if err != nil {
					return newError("invalid regex pattern: %s", err.Error())
				}

				// Create regex dictionary
				pairs := make(map[string]ast.Expression)
				pairs["__type"] = &ast.StringLiteral{Value: "regex"}
				pairs["pattern"] = &ast.StringLiteral{Value: pattern.Value}
				pairs["flags"] = &ast.StringLiteral{Value: flags}

				return &Dictionary{Pairs: pairs, Env: NewEnvironment()}
			},
		},
		"replace": {
			Fn: func(args ...Object) Object {
				if len(args) != 3 {
					return newError("wrong number of arguments to `replace`. got=%d, want=3", len(args))
				}

				text, ok := args[0].(*String)
				if !ok {
					return newError("first argument to `replace` must be a string, got %s", args[0].Type())
				}

				// Second arg can be string or regex
				var pattern string
				var flags string
				if str, ok := args[1].(*String); ok {
					// String pattern - use literal replacement
					replacement, ok := args[2].(*String)
					if !ok {
						return newError("third argument to `replace` must be a string, got %s", args[2].Type())
					}
					return &String{Value: strings.Replace(text.Value, str.Value, replacement.Value, -1)}
				} else if dict, ok := args[1].(*Dictionary); ok && isRegexDict(dict) {
					// Regex pattern
					patternExpr, _ := dict.Pairs["pattern"]
					patternObj := Eval(patternExpr, NewEnvironment())
					patternStr := patternObj.(*String)
					pattern = patternStr.Value

					flagsExpr, ok := dict.Pairs["flags"]
					if ok {
						flagsObj := Eval(flagsExpr, NewEnvironment())
						if flagsStr, ok := flagsObj.(*String); ok {
							flags = flagsStr.Value
						}
					}
				} else {
					return newError("second argument to `replace` must be a string or regex, got %s", args[1].Type())
				}

				replacement, ok := args[2].(*String)
				if !ok {
					return newError("third argument to `replace` must be a string, got %s", args[2].Type())
				}

				re, err := compileRegex(pattern, flags)
				if err != nil {
					return newError("invalid regex: %s", err.Error())
				}

				result := re.ReplaceAllString(text.Value, replacement.Value)
				return &String{Value: result}
			},
		},
		"split": {
			Fn: func(args ...Object) Object {
				if len(args) != 2 {
					return newError("wrong number of arguments to `split`. got=%d, want=2", len(args))
				}

				text, ok := args[0].(*String)
				if !ok {
					return newError("first argument to `split` must be a string, got %s", args[0].Type())
				}

				// Second arg can be string or regex
				var parts []string
				if str, ok := args[1].(*String); ok {
					// String delimiter
					parts = strings.Split(text.Value, str.Value)
				} else if dict, ok := args[1].(*Dictionary); ok && isRegexDict(dict) {
					// Regex pattern
					patternExpr, _ := dict.Pairs["pattern"]
					patternObj := Eval(patternExpr, NewEnvironment())
					patternStr := patternObj.(*String)
					pattern := patternStr.Value

					flags := ""
					flagsExpr, ok := dict.Pairs["flags"]
					if ok {
						flagsObj := Eval(flagsExpr, NewEnvironment())
						if flagsStr, ok := flagsObj.(*String); ok {
							flags = flagsStr.Value
						}
					}

					re, err := compileRegex(pattern, flags)
					if err != nil {
						return newError("invalid regex: %s", err.Error())
					}

					parts = re.Split(text.Value, -1)
				} else {
					return newError("second argument to `split` must be a string or regex, got %s", args[1].Type())
				}

				elements := make([]Object, len(parts))
				for i, part := range parts {
					elements[i] = &String{Value: part}
				}

				return &Array{Elements: elements}
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
			return evalDestructuringAssignment(node.Names, val, env, true)
		}

		// Single assignment
		// Special handling for '_' - don't store it
		if node.Name.Value != "_" {
			env.SetLet(node.Name.Value, val)
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
			return evalDestructuringAssignment(node.Names, val, env, false)
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

	case *ast.RegexLiteral:
		return evalRegexLiteral(node, env)

	case *ast.DatetimeLiteral:
		return evalDatetimeLiteral(node, env)

	case *ast.DurationLiteral:
		return evalDurationLiteral(node, env)

	case *ast.PathLiteral:
		return evalPathLiteral(node, env)

	case *ast.UrlLiteral:
		return evalUrlLiteral(node, env)

	case *ast.TagLiteral:
		return evalTagLiteral(node, env)

	case *ast.TagPairExpression:
		return evalTagPair(node, env)

	case *ast.TextNode:
		return &String{Value: node.Value}

	case *ast.Boolean:
		return nativeBoolToParsBoolean(node.Value)

	case *ast.ObjectLiteralExpression:
		return node.Obj.(Object)

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
		return evalInfixExpression(node.Token, node.Operator, left, right)

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

		// Check if this is a call to import
		if ident, ok := node.Function.(*ast.Identifier); ok && ident.Value == "import" {
			args := evalExpressions(node.Arguments, env)
			if len(args) == 1 && isError(args[0]) {
				return args[0]
			}
			return evalImport(args, env)
		}

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
		return evalIndexExpression(node.Token, left, index)

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

func evalInfixExpression(tok lexer.Token, operator string, left, right Object) Object {
	switch {
	case operator == "&" || operator == "and":
		return nativeBoolToParsBoolean(isTruthy(left) && isTruthy(right))
	case operator == "|" || operator == "or":
		return nativeBoolToParsBoolean(isTruthy(left) || isTruthy(right))
	case operator == "++":
		return evalConcatExpression(left, right)
	// Path and URL operators with strings (must come before general string concatenation)
	case left.Type() == DICTIONARY_OBJ && right.Type() == STRING_OBJ:
		if dict := left.(*Dictionary); isPathDict(dict) {
			return evalPathStringInfixExpression(tok, operator, dict, right.(*String))
		}
		if dict := left.(*Dictionary); isUrlDict(dict) {
			return evalUrlStringInfixExpression(tok, operator, dict, right.(*String))
		}
		// Fall through to string concatenation if not path/url
		if operator == "+" {
			return evalStringConcatExpression(left, right)
		}
		return newErrorWithPos(tok, "unknown operator: %s %s %s", left.Type(), operator, right.Type())
	case operator == "+" && (left.Type() == STRING_OBJ || right.Type() == STRING_OBJ):
		// String concatenation with automatic type conversion
		return evalStringConcatExpression(left, right)
	// Regex match operators
	case operator == "~" || operator == "!~":
		if left.Type() != STRING_OBJ {
			return newErrorWithPos(tok, "left operand of %s must be a string, got %s", operator, left.Type())
		}
		if right.Type() != DICTIONARY_OBJ {
			return newErrorWithPos(tok, "right operand of %s must be a regex, got %s", operator, right.Type())
		}
		rightDict := right.(*Dictionary)
		if !isRegexDict(rightDict) {
			return newErrorWithPos(tok, "right operand of %s must be a regex dictionary", operator)
		}
		result := evalMatchExpression(tok, left.(*String).Value, rightDict, NewEnvironment())
		if operator == "!~" {
			// !~ returns boolean: true if no match, false if match
			return nativeBoolToParsBoolean(result == NULL)
		}
		return result // ~ returns array or null
	// Datetime dictionary operations
	case left.Type() == DICTIONARY_OBJ && right.Type() == DICTIONARY_OBJ:
		leftDict := left.(*Dictionary)
		rightDict := right.(*Dictionary)
		if isDatetimeDict(leftDict) && isDatetimeDict(rightDict) {
			return evalDatetimeInfixExpression(tok, operator, leftDict, rightDict)
		}
		if isDurationDict(leftDict) && isDurationDict(rightDict) {
			return evalDurationInfixExpression(tok, operator, leftDict, rightDict)
		}
		if isDatetimeDict(leftDict) && isDurationDict(rightDict) {
			return evalDatetimeDurationInfixExpression(tok, operator, leftDict, rightDict)
		}
		if isDurationDict(leftDict) && isDatetimeDict(rightDict) {
			// duration + datetime not allowed, only datetime + duration
			return newErrorWithPos(tok, "cannot add datetime to duration (use datetime + duration instead)")
		}
		// Path dictionary operations
		if isPathDict(leftDict) && isPathDict(rightDict) {
			return evalPathInfixExpression(tok, operator, leftDict, rightDict)
		}
		// URL dictionary operations
		if isUrlDict(leftDict) && isUrlDict(rightDict) {
			return evalUrlInfixExpression(tok, operator, leftDict, rightDict)
		}
		// Fall through to default comparison for non-datetime dicts
		if operator == "==" {
			return nativeBoolToParsBoolean(left == right)
		} else if operator == "!=" {
			return nativeBoolToParsBoolean(left != right)
		}
		return newErrorWithPos(tok, "unknown operator: %s %s %s", left.Type(), operator, right.Type())
	case left.Type() == DICTIONARY_OBJ && right.Type() == INTEGER_OBJ:
		if dict := left.(*Dictionary); isDatetimeDict(dict) {
			return evalDatetimeIntegerInfixExpression(tok, operator, dict, right.(*Integer))
		}
		if dict := left.(*Dictionary); isDurationDict(dict) {
			return evalDurationIntegerInfixExpression(tok, operator, dict, right.(*Integer))
		}
		return newErrorWithPos(tok, "unknown operator: %s %s %s", left.Type(), operator, right.Type())
	case left.Type() == INTEGER_OBJ && right.Type() == DICTIONARY_OBJ:
		if dict := right.(*Dictionary); isDatetimeDict(dict) {
			return evalIntegerDatetimeInfixExpression(tok, operator, left.(*Integer), dict)
		}
		return newErrorWithPos(tok, "unknown operator: %s %s %s", left.Type(), operator, right.Type())
	case left.Type() == INTEGER_OBJ && right.Type() == INTEGER_OBJ:
		return evalIntegerInfixExpression(tok, operator, left, right)
	case left.Type() == FLOAT_OBJ && right.Type() == FLOAT_OBJ:
		return evalFloatInfixExpression(tok, operator, left, right)
	case left.Type() == INTEGER_OBJ && right.Type() == FLOAT_OBJ:
		return evalMixedInfixExpression(tok, operator, left, right)
	case left.Type() == FLOAT_OBJ && right.Type() == INTEGER_OBJ:
		return evalMixedInfixExpression(tok, operator, left, right)
	case left.Type() == STRING_OBJ && right.Type() == STRING_OBJ:
		return evalStringInfixExpression(operator, left, right)
	case operator == "==":
		return nativeBoolToParsBoolean(left == right)
	case operator == "!=":
		return nativeBoolToParsBoolean(left != right)
	case left.Type() != right.Type():
		return newErrorWithPos(tok, "type mismatch: %s %s %s", left.Type(), operator, right.Type())
	default:
		return newErrorWithPos(tok, "unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalIntegerInfixExpression(tok lexer.Token, operator string, left, right Object) Object {
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
			return newErrorWithPos(tok, "division by zero")
		}
		return &Integer{Value: leftVal / rightVal}
	case "%":
		if rightVal == 0 {
			return newErrorWithPos(tok, "modulo by zero")
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

func evalFloatInfixExpression(tok lexer.Token, operator string, left, right Object) Object {
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
			return newErrorWithPos(tok, "division by zero")
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

func evalMixedInfixExpression(tok lexer.Token, operator string, left, right Object) Object {
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
			return newErrorWithPos(tok, "division by zero")
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

// evalDatetimeInfixExpression handles operations between two datetime dictionaries
func evalDatetimeInfixExpression(tok lexer.Token, operator string, left, right *Dictionary) Object {
	env := NewEnvironment()
	leftUnix, err := getDatetimeUnix(left, env)
	if err != nil {
		return newErrorWithPos(tok, "invalid datetime: %s", err)
	}
	rightUnix, err := getDatetimeUnix(right, env)
	if err != nil {
		return newErrorWithPos(tok, "invalid datetime: %s", err)
	}

	switch operator {
	case "<":
		return nativeBoolToParsBoolean(leftUnix < rightUnix)
	case ">":
		return nativeBoolToParsBoolean(leftUnix > rightUnix)
	case "<=":
		return nativeBoolToParsBoolean(leftUnix <= rightUnix)
	case ">=":
		return nativeBoolToParsBoolean(leftUnix >= rightUnix)
	case "==":
		return nativeBoolToParsBoolean(leftUnix == rightUnix)
	case "!=":
		return nativeBoolToParsBoolean(leftUnix != rightUnix)
	case "-":
		// BREAKING CHANGE: Return Duration instead of Integer
		// Calculate difference in seconds
		diffSeconds := leftUnix - rightUnix
		// Return as duration (0 months, diffSeconds seconds)
		return durationToDict(0, diffSeconds, env)
	default:
		return newErrorWithPos(tok, "unknown operator for datetime: %s", operator)
	}
}

// evalDatetimeIntegerInfixExpression handles datetime + integer or datetime - integer
func evalDatetimeIntegerInfixExpression(tok lexer.Token, operator string, dt *Dictionary, seconds *Integer) Object {
	env := NewEnvironment()
	unixTime, err := getDatetimeUnix(dt, env)
	if err != nil {
		return newErrorWithPos(tok, "invalid datetime: %s", err)
	}

	switch operator {
	case "+":
		// Add seconds to datetime
		newTime := time.Unix(unixTime+seconds.Value, 0).UTC()
		return timeToDict(newTime, env)
	case "-":
		// Subtract seconds from datetime
		newTime := time.Unix(unixTime-seconds.Value, 0).UTC()
		return timeToDict(newTime, env)
	default:
		return newErrorWithPos(tok, "unknown operator for datetime and integer: %s", operator)
	}
}

// evalIntegerDatetimeInfixExpression handles integer + datetime
func evalIntegerDatetimeInfixExpression(tok lexer.Token, operator string, seconds *Integer, dt *Dictionary) Object {
	env := NewEnvironment()
	unixTime, err := getDatetimeUnix(dt, env)
	if err != nil {
		return newErrorWithPos(tok, "invalid datetime: %s", err)
	}

	switch operator {
	case "+":
		// Add seconds to datetime (commutative)
		newTime := time.Unix(unixTime+seconds.Value, 0).UTC()
		return timeToDict(newTime, env)
	default:
		return newErrorWithPos(tok, "unknown operator for integer and datetime: %s", operator)
	}
}

// evalDurationInfixExpression handles duration + duration or duration - duration
func evalDurationInfixExpression(tok lexer.Token, operator string, left, right *Dictionary) Object {
	env := NewEnvironment()

	leftMonths, leftSeconds, err := getDurationComponents(left, env)
	if err != nil {
		return newErrorWithPos(tok, "invalid duration: %s", err)
	}

	rightMonths, rightSeconds, err := getDurationComponents(right, env)
	if err != nil {
		return newErrorWithPos(tok, "invalid duration: %s", err)
	}

	switch operator {
	case "+":
		return durationToDict(leftMonths+rightMonths, leftSeconds+rightSeconds, env)
	case "-":
		return durationToDict(leftMonths-rightMonths, leftSeconds-rightSeconds, env)
	case "<", ">", "<=", ">=", "==", "!=":
		// Comparison only allowed for pure-seconds durations (no months)
		if leftMonths != 0 || rightMonths != 0 {
			return newErrorWithPos(tok, "cannot compare durations with month components (months have variable length)")
		}
		switch operator {
		case "<":
			return nativeBoolToParsBoolean(leftSeconds < rightSeconds)
		case ">":
			return nativeBoolToParsBoolean(leftSeconds > rightSeconds)
		case "<=":
			return nativeBoolToParsBoolean(leftSeconds <= rightSeconds)
		case ">=":
			return nativeBoolToParsBoolean(leftSeconds >= rightSeconds)
		case "==":
			return nativeBoolToParsBoolean(leftSeconds == rightSeconds && leftMonths == rightMonths)
		case "!=":
			return nativeBoolToParsBoolean(leftSeconds != rightSeconds || leftMonths != rightMonths)
		}
	}

	return newErrorWithPos(tok, "unknown operator for duration: %s", operator)
}

// evalDurationIntegerInfixExpression handles duration * integer or duration / integer
func evalDurationIntegerInfixExpression(tok lexer.Token, operator string, dur *Dictionary, num *Integer) Object {
	env := NewEnvironment()

	months, seconds, err := getDurationComponents(dur, env)
	if err != nil {
		return newErrorWithPos(tok, "invalid duration: %s", err)
	}

	switch operator {
	case "*":
		return durationToDict(months*num.Value, seconds*num.Value, env)
	case "/":
		if num.Value == 0 {
			return newErrorWithPos(tok, "division by zero")
		}
		return durationToDict(months/num.Value, seconds/num.Value, env)
	default:
		return newErrorWithPos(tok, "unknown operator for duration and integer: %s", operator)
	}
}

// evalDatetimeDurationInfixExpression handles datetime + duration or datetime - duration
func evalDatetimeDurationInfixExpression(tok lexer.Token, operator string, dt, dur *Dictionary) Object {
	env := NewEnvironment()

	// Get datetime as time.Time
	t, err := dictToTime(dt, env)
	if err != nil {
		return newErrorWithPos(tok, "invalid datetime: %s", err)
	}

	// Get duration components
	months, seconds, err := getDurationComponents(dur, env)
	if err != nil {
		return newErrorWithPos(tok, "invalid duration: %s", err)
	}

	switch operator {
	case "+":
		// Add months first (using AddDate for proper month arithmetic)
		if months != 0 {
			t = t.AddDate(0, int(months), 0)
		}
		// Then add seconds
		if seconds != 0 {
			t = t.Add(time.Duration(seconds) * time.Second)
		}
		return timeToDict(t, env)
	case "-":
		// Subtract months first
		if months != 0 {
			t = t.AddDate(0, -int(months), 0)
		}
		// Then subtract seconds
		if seconds != 0 {
			t = t.Add(-time.Duration(seconds) * time.Second)
		}
		return timeToDict(t, env)
	default:
		return newErrorWithPos(tok, "unknown operator for datetime and duration: %s", operator)
	}
}

// evalPathInfixExpression handles operations between two path dictionaries
func evalPathInfixExpression(tok lexer.Token, operator string, left, right *Dictionary) Object {
	switch operator {
	case "==":
		// Compare paths by their string representation
		leftStr := pathDictToString(left)
		rightStr := pathDictToString(right)
		return nativeBoolToParsBoolean(leftStr == rightStr)
	case "!=":
		leftStr := pathDictToString(left)
		rightStr := pathDictToString(right)
		return nativeBoolToParsBoolean(leftStr != rightStr)
	default:
		return newErrorWithPos(tok, "unknown operator for path: %s (supported: ==, !=)", operator)
	}
}

// evalPathStringInfixExpression handles path + string or path / string
func evalPathStringInfixExpression(tok lexer.Token, operator string, path *Dictionary, str *String) Object {
	env := path.Env
	if env == nil {
		env = NewEnvironment()
	}

	switch operator {
	case "+", "/":
		// Join path with string segment
		// Get current components
		componentsExpr, ok := path.Pairs["components"]
		if !ok {
			return newErrorWithPos(tok, "path dictionary missing components field")
		}
		componentsObj := Eval(componentsExpr, env)
		if componentsObj.Type() != ARRAY_OBJ {
			return newErrorWithPos(tok, "path components is not an array")
		}
		componentsArr := componentsObj.(*Array)

		// Get absolute flag
		absoluteExpr, ok := path.Pairs["absolute"]
		if !ok {
			return newErrorWithPos(tok, "path dictionary missing absolute field")
		}
		absoluteObj := Eval(absoluteExpr, env)
		if absoluteObj.Type() != BOOLEAN_OBJ {
			return newErrorWithPos(tok, "path absolute is not a boolean")
		}
		isAbsolute := absoluteObj.(*Boolean).Value

		// Parse the string to add as new path segments
		newSegments, _ := parsePathString(str.Value)

		// Combine components
		var newComponents []string
		for _, elem := range componentsArr.Elements {
			if strObj, ok := elem.(*String); ok {
				newComponents = append(newComponents, strObj.Value)
			}
		}

		// Append new segments (skip empty leading segment if present)
		for _, seg := range newSegments {
			if seg != "" || len(newComponents) == 0 {
				newComponents = append(newComponents, seg)
			}
		}

		return pathToDict(newComponents, isAbsolute, env)
	default:
		return newErrorWithPos(tok, "unknown operator for path and string: %s (supported: +, /)", operator)
	}
}

// evalUrlInfixExpression handles operations between two URL dictionaries
func evalUrlInfixExpression(tok lexer.Token, operator string, left, right *Dictionary) Object {
	switch operator {
	case "==":
		// Compare URLs by their string representation
		leftStr := urlDictToString(left)
		rightStr := urlDictToString(right)
		return nativeBoolToParsBoolean(leftStr == rightStr)
	case "!=":
		leftStr := urlDictToString(left)
		rightStr := urlDictToString(right)
		return nativeBoolToParsBoolean(leftStr != rightStr)
	default:
		return newErrorWithPos(tok, "unknown operator for url: %s (supported: ==, !=)", operator)
	}
}

// evalUrlStringInfixExpression handles url + string for path joining
func evalUrlStringInfixExpression(tok lexer.Token, operator string, urlDict *Dictionary, str *String) Object {
	env := urlDict.Env
	if env == nil {
		env = NewEnvironment()
	}

	switch operator {
	case "+":
		// Add string to URL path
		// Get current path array
		pathExpr, ok := urlDict.Pairs["path"]
		if !ok {
			return newErrorWithPos(tok, "url dictionary missing path field")
		}
		pathObj := Eval(pathExpr, env)
		if pathObj.Type() != ARRAY_OBJ {
			return newErrorWithPos(tok, "url path is not an array")
		}
		pathArr := pathObj.(*Array)

		// Parse the string as a path to add
		newSegments, _ := parsePathString(str.Value)

		// Combine path segments
		var newPath []string
		for _, elem := range pathArr.Elements {
			if strObj, ok := elem.(*String); ok {
				newPath = append(newPath, strObj.Value)
			}
		}

		// Append new segments (skip empty leading segment)
		for _, seg := range newSegments {
			if seg != "" {
				newPath = append(newPath, seg)
			}
		}

		// Create new URL dict with updated path
		pairs := make(map[string]ast.Expression)
		for k, v := range urlDict.Pairs {
			if k == "path" {
				// Create new path array
				pathElements := make([]ast.Expression, len(newPath))
				for i, seg := range newPath {
					pathElements[i] = &ast.StringLiteral{Value: seg}
				}
				pairs[k] = &ast.ArrayLiteral{Elements: pathElements}
			} else {
				pairs[k] = v
			}
		}

		return &Dictionary{Pairs: pairs, Env: env}
	default:
		return newErrorWithPos(tok, "unknown operator for url and string: %s (supported: +)", operator)
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
		return newErrorWithPos(node.Token, "identifier not found: %s", node.Value)
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

// evalImport implements the import(path) builtin
func evalImport(args []Object, env *Environment) Object {
	if len(args) != 1 {
		return newError("wrong number of arguments to `import`. got=%d, want=1", len(args))
	}

	// Extract path string from argument (handle both path dictionaries and strings)
	var pathStr string
	switch arg := args[0].(type) {
	case *Dictionary:
		// Handle path literal (@/path/to/file.pars)
		if typeExpr, ok := arg.Pairs["__type"]; ok {
			typeVal := Eval(typeExpr, arg.Env)
			if typeStr, ok := typeVal.(*String); ok && typeStr.Value == "path" {
				pathStr = pathDictToString(arg)
			} else {
				return newError("argument to `import` must be a path or string, got dictionary")
			}
		} else {
			return newError("argument to `import` must be a path or string, got dictionary")
		}
	case *String:
		pathStr = arg.Value
	default:
		return newError("argument to `import` must be a path or string, got %s", arg.Type())
	}

	// Resolve path relative to current file
	absPath, err := resolveModulePath(pathStr, env.Filename)
	if err != nil {
		return newError("failed to resolve module path: %s", err.Error())
	}

	// Check if module is currently being loaded (circular dependency)
	if moduleCache.loading[absPath] {
		return newError("circular dependency detected when importing: %s", absPath)
	}

	// Check cache first
	if cached, ok := moduleCache.modules[absPath]; ok {
		return cached
	}

	// Mark as loading
	moduleCache.loading[absPath] = true
	defer delete(moduleCache.loading, absPath)

	// Read the file
	content, err := os.ReadFile(absPath)
	if err != nil {
		return newError("failed to read module file %s: %s", absPath, err.Error())
	}

	// Parse the module
	l := lexer.New(string(content))
	p := parser.New(l)
	program := p.ParseProgram()
	if len(p.Errors()) > 0 {
		var errMsg strings.Builder
		errMsg.WriteString(fmt.Sprintf("parse errors in module %s:\n", absPath))
		for _, msg := range p.Errors() {
			errMsg.WriteString(fmt.Sprintf("  %s\n", msg))
		}
		return newError(errMsg.String())
	}

	// Create isolated environment for the module
	moduleEnv := NewEnvironment()
	moduleEnv.Filename = absPath

	// Evaluate the module
	Eval(program, moduleEnv)

	// Convert environment to dictionary
	moduleDict := environmentToDict(moduleEnv)

	// Cache the result
	moduleCache.modules[absPath] = moduleDict

	return moduleDict
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

// newErrorWithPos creates an error with position information from a token
func newErrorWithPos(tok lexer.Token, format string, a ...interface{}) *Error {
	return &Error{
		Message: fmt.Sprintf(format, a...),
		Line:    tok.Line,
		Column:  tok.Column,
	}
}

func isError(obj Object) bool {
	if obj != nil {
		return obj.Type() == ERROR_OBJ
	}
	return false
}

// evalDestructuringAssignment handles array destructuring assignment
func evalDestructuringAssignment(names []*ast.Identifier, val Object, env *Environment, isLet bool) Object {
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
				if isLet {
					env.SetLet(name.Value, elements[i])
				} else {
					env.Update(name.Value, elements[i])
				}
			}
		} else {
			// No more elements, assign null
			if name.Value != "_" {
				if isLet {
					env.SetLet(name.Value, NULL)
				} else {
					env.Update(name.Value, NULL)
				}
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
			if isLet {
				env.SetLet(lastName.Value, remaining)
			} else {
				env.Update(lastName.Value, remaining)
			}
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
				env.SetLet(pattern.Rest.Value, restDict)
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
			return newError("function not found: %s", tagName)
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
	case *Dictionary:
		// Check for special dictionary types
		if isPathDict(obj) {
			// Convert path dictionary back to string
			return pathDictToString(obj)
		}
		if isUrlDict(obj) {
			// Convert URL dictionary back to string
			return urlDictToString(obj)
		}
		return obj.Inspect()
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
func evalIndexExpression(tok lexer.Token, left, index Object) Object {
	switch {
	case left.Type() == ARRAY_OBJ && index.Type() == INTEGER_OBJ:
		return evalArrayIndexExpression(tok, left, index)
	case left.Type() == STRING_OBJ && index.Type() == INTEGER_OBJ:
		return evalStringIndexExpression(tok, left, index)
	case left.Type() == DICTIONARY_OBJ && index.Type() == STRING_OBJ:
		return evalDictionaryIndexExpression(left, index)
	default:
		return newErrorWithPos(tok, "index operator not supported: %s[%s]", left.Type(), index.Type())
	}
}

// evalArrayIndexExpression handles array indexing with support for negative indices
func evalArrayIndexExpression(tok lexer.Token, array, index Object) Object {
	arrayObject := array.(*Array)
	idx := index.(*Integer).Value
	max := int64(len(arrayObject.Elements))

	// Handle negative indices
	if idx < 0 {
		idx = max + idx
	}

	if idx < 0 || idx >= max {
		return newErrorWithPos(tok, "index out of range: %d", index.(*Integer).Value)
	}

	return arrayObject.Elements[idx]
}

// evalStringIndexExpression handles string indexing with support for negative indices
func evalStringIndexExpression(tok lexer.Token, str, index Object) Object {
	stringObject := str.(*String)
	idx := index.(*Integer).Value
	max := int64(len(stringObject.Value))

	// Handle negative indices
	if idx < 0 {
		idx = max + idx
	}

	if idx < 0 || idx >= max {
		return newErrorWithPos(tok, "index out of range: %d", index.(*Integer).Value)
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
		return newErrorWithPos(node.Token, "dot notation can only be used on dictionaries, got %s", left.Type())
	}

	// Check for computed properties on special dictionary types
	if isPathDict(dict) {
		if computed := evalPathComputedProperty(dict, node.Key, env); computed != nil {
			return computed
		}
	}
	if isUrlDict(dict) {
		if computed := evalUrlComputedProperty(dict, node.Key, env); computed != nil {
			return computed
		}
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

// environmentToDict converts an environment's store to a Dictionary object
// Only includes variables that were declared with 'let'
func environmentToDict(env *Environment) *Dictionary {
	pairs := make(map[string]ast.Expression)

	// Only export variables that were declared with 'let'
	for name, value := range env.store {
		if env.IsLetBinding(name) {
			// Wrap the object as a literal expression
			pairs[name] = objectToExpression(value)
		}
	}

	// Create dictionary with the module's environment for evaluation
	return &Dictionary{Pairs: pairs, Env: env}
}

// objectToExpression wraps an Object as an AST expression
func objectToExpression(obj Object) ast.Expression {
	switch v := obj.(type) {
	case *Integer:
		return &ast.IntegerLiteral{Value: v.Value}
	case *Float:
		return &ast.FloatLiteral{Value: v.Value}
	case *String:
		return &ast.StringLiteral{Value: v.Value}
	case *Boolean:
		return &ast.Boolean{Value: v.Value}
	default:
		// For complex types (functions, arrays, dictionaries, null), we create
		// an expression that returns the object directly when evaluated
		return &ast.ObjectLiteralExpression{Obj: obj}
	}
}

// objectLiteralExpression removed - now using ast.ObjectLiteralExpression

// resolveModulePath resolves a module path relative to the current file
func resolveModulePath(pathStr string, currentFile string) (string, error) {
	var absPath string

	// If path is absolute, use it directly
	if strings.HasPrefix(pathStr, "/") {
		absPath = pathStr
	} else {
		// Resolve relative to the current file's directory
		var baseDir string
		if currentFile != "" {
			baseDir = filepath.Dir(currentFile)
		} else {
			// If no current file, use current working directory
			cwd, err := os.Getwd()
			if err != nil {
				return "", err
			}
			baseDir = cwd
		}

		// Join and clean the path
		absPath = filepath.Join(baseDir, pathStr)
	}

	// Clean the path (resolve . and ..)
	absPath = filepath.Clean(absPath)

	return absPath, nil
}
