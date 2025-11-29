// Package evaluator provides method implementations for primitive types
// This file implements the method-call API for String, Array, Integer, Float types
package evaluator

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/sambeau/parsley/pkg/locale"
)

// ============================================================================
// String Methods
// ============================================================================

// evalStringMethod evaluates a method call on a String
func evalStringMethod(str *String, method string, args []Object) Object {
	switch method {
	case "upper":
		if len(args) != 0 {
			return newError("wrong number of arguments for 'upper'. got=%d, want=0", len(args))
		}
		return &String{Value: strings.ToUpper(str.Value)}

	case "lower":
		if len(args) != 0 {
			return newError("wrong number of arguments for 'lower'. got=%d, want=0", len(args))
		}
		return &String{Value: strings.ToLower(str.Value)}

	case "trim":
		if len(args) != 0 {
			return newError("wrong number of arguments for 'trim'. got=%d, want=0", len(args))
		}
		return &String{Value: strings.TrimSpace(str.Value)}

	case "split":
		if len(args) != 1 {
			return newError("wrong number of arguments for 'split'. got=%d, want=1", len(args))
		}
		delim, ok := args[0].(*String)
		if !ok {
			return newError("argument to 'split' must be STRING, got %s", args[0].Type())
		}
		parts := strings.Split(str.Value, delim.Value)
		elements := make([]Object, len(parts))
		for i, part := range parts {
			elements[i] = &String{Value: part}
		}
		return &Array{Elements: elements}

	case "replace":
		if len(args) != 2 {
			return newError("wrong number of arguments for 'replace'. got=%d, want=2", len(args))
		}
		old, ok := args[0].(*String)
		if !ok {
			return newError("first argument to 'replace' must be STRING, got %s", args[0].Type())
		}
		new, ok := args[1].(*String)
		if !ok {
			return newError("second argument to 'replace' must be STRING, got %s", args[1].Type())
		}
		return &String{Value: strings.ReplaceAll(str.Value, old.Value, new.Value)}

	case "length":
		if len(args) != 0 {
			return newError("wrong number of arguments for 'length'. got=%d, want=0", len(args))
		}
		// Return rune count for proper Unicode support
		return &Integer{Value: int64(len([]rune(str.Value)))}

	default:
		return newError("unknown method '%s' for STRING", method)
	}
}

// ============================================================================
// Array Methods
// ============================================================================

// evalArrayMethod evaluates a method call on an Array
func evalArrayMethod(arr *Array, method string, args []Object, env *Environment) Object {
	switch method {
	case "length":
		if len(args) != 0 {
			return newError("wrong number of arguments for 'length'. got=%d, want=0", len(args))
		}
		return &Integer{Value: int64(len(arr.Elements))}

	case "reverse":
		if len(args) != 0 {
			return newError("wrong number of arguments for 'reverse'. got=%d, want=0", len(args))
		}
		// Create a reversed copy
		length := len(arr.Elements)
		newElements := make([]Object, length)
		for i, elem := range arr.Elements {
			newElements[length-1-i] = elem
		}
		return &Array{Elements: newElements}

	case "sort":
		if len(args) != 0 {
			return newError("wrong number of arguments for 'sort'. got=%d, want=0", len(args))
		}
		return naturalSortArray(arr)

	case "sortBy":
		if len(args) != 1 {
			return newError("wrong number of arguments for 'sortBy'. got=%d, want=1", len(args))
		}
		fn, ok := args[0].(*Function)
		if !ok {
			return newError("argument to 'sortBy' must be a function, got %s", args[0].Type())
		}
		return sortArrayByFunction(arr, fn, env)

	case "map":
		if len(args) != 1 {
			return newError("wrong number of arguments for 'map'. got=%d, want=1", len(args))
		}
		fn, ok := args[0].(*Function)
		if !ok {
			return newError("argument to 'map' must be a function, got %s", args[0].Type())
		}
		return mapArrayWithFunction(arr, fn, env)

	case "filter":
		if len(args) != 1 {
			return newError("wrong number of arguments for 'filter'. got=%d, want=1", len(args))
		}
		fn, ok := args[0].(*Function)
		if !ok {
			return newError("argument to 'filter' must be a function, got %s", args[0].Type())
		}
		return filterArrayWithFunction(arr, fn, env)

	case "format":
		// format(style?, locale?)
		if len(args) > 2 {
			return newError("wrong number of arguments for 'format'. got=%d, want=0-2", len(args))
		}

		// Convert array elements to strings
		items := make([]string, len(arr.Elements))
		for i, elem := range arr.Elements {
			items[i] = elem.Inspect()
		}

		// Get style (default to "and")
		style := locale.ListStyleAnd
		localeStr := "en-US"

		if len(args) >= 1 {
			styleStr, ok := args[0].(*String)
			if !ok {
				return newError("first argument to 'format' must be STRING (style), got %s", args[0].Type())
			}
			switch styleStr.Value {
			case "and":
				style = locale.ListStyleAnd
			case "or":
				style = locale.ListStyleOr
			case "unit":
				style = locale.ListStyleUnit
			default:
				return newError("invalid style %q for 'format', use 'and', 'or', or 'unit'", styleStr.Value)
			}
		}

		if len(args) == 2 {
			locStr, ok := args[1].(*String)
			if !ok {
				return newError("second argument to 'format' must be STRING (locale), got %s", args[1].Type())
			}
			localeStr = locStr.Value
		}

		result := locale.FormatList(items, style, localeStr)
		return &String{Value: result}

	case "join":
		// join(separator?) - joins array elements into a string
		if len(args) > 1 {
			return newError("wrong number of arguments for 'join'. got=%d, want=0-1", len(args))
		}

		separator := ""
		if len(args) == 1 {
			sepStr, ok := args[0].(*String)
			if !ok {
				return newError("argument to 'join' must be a STRING, got %s", args[0].Type())
			}
			separator = sepStr.Value
		}

		// Convert array elements to strings
		items := make([]string, len(arr.Elements))
		for i, elem := range arr.Elements {
			items[i] = objectToTemplateString(elem)
		}

		return &String{Value: strings.Join(items, separator)}

	default:
		return newError("unknown method '%s' for ARRAY", method)
	}
}

// naturalSortArray performs a natural sort on an array
func naturalSortArray(arr *Array) *Array {
	// Make a copy of elements
	elements := make([]Object, len(arr.Elements))
	copy(elements, arr.Elements)

	// Sort using natural comparison
	sort.SliceStable(elements, func(i, j int) bool {
		return compareObjects(elements[i], elements[j]) < 0
	})

	return &Array{Elements: elements}
}

// sortArrayByFunction sorts an array using a key function
func sortArrayByFunction(arr *Array, fn *Function, env *Environment) Object {
	// Make a copy of elements
	elements := make([]Object, len(arr.Elements))
	copy(elements, arr.Elements)

	// Compute keys for all elements
	keys := make([]Object, len(elements))
	for i, elem := range elements {
		extendedEnv := extendFunctionEnv(fn, []Object{elem})
		result := Eval(fn.Body, extendedEnv)
		if isError(result) {
			return result
		}
		if returnValue, ok := result.(*ReturnValue); ok {
			result = returnValue.Value
		}
		keys[i] = result
	}

	// Sort by keys
	sort.SliceStable(elements, func(i, j int) bool {
		return compareObjects(keys[i], keys[j]) < 0
	})

	return &Array{Elements: elements}
}

// mapArrayWithFunction applies a function to each element
func mapArrayWithFunction(arr *Array, fn *Function, env *Environment) Object {
	result := make([]Object, 0, len(arr.Elements))

	for _, elem := range arr.Elements {
		extendedEnv := extendFunctionEnv(fn, []Object{elem})

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

// filterArrayWithFunction filters array elements based on a predicate function
func filterArrayWithFunction(arr *Array, fn *Function, env *Environment) Object {
	result := make([]Object, 0, len(arr.Elements))

	for _, elem := range arr.Elements {
		extendedEnv := extendFunctionEnv(fn, []Object{elem})

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

		// Include element if predicate returns truthy value
		if isTruthy(evaluated) {
			result = append(result, elem)
		}
	}

	return &Array{Elements: result}
}

// compareObjects compares two objects for sorting
func compareObjects(a, b Object) int {
	// Handle nil/NULL
	if a == nil || a == NULL {
		if b == nil || b == NULL {
			return 0
		}
		return -1
	}
	if b == nil || b == NULL {
		return 1
	}

	// Compare by type
	switch av := a.(type) {
	case *Integer:
		if bv, ok := b.(*Integer); ok {
			if av.Value < bv.Value {
				return -1
			} else if av.Value > bv.Value {
				return 1
			}
			return 0
		}
		if bv, ok := b.(*Float); ok {
			af := float64(av.Value)
			if af < bv.Value {
				return -1
			} else if af > bv.Value {
				return 1
			}
			return 0
		}
	case *Float:
		if bv, ok := b.(*Float); ok {
			if av.Value < bv.Value {
				return -1
			} else if av.Value > bv.Value {
				return 1
			}
			return 0
		}
		if bv, ok := b.(*Integer); ok {
			bf := float64(bv.Value)
			if av.Value < bf {
				return -1
			} else if av.Value > bf {
				return 1
			}
			return 0
		}
	case *String:
		if bv, ok := b.(*String); ok {
			return strings.Compare(av.Value, bv.Value)
		}
	case *Boolean:
		if bv, ok := b.(*Boolean); ok {
			if !av.Value && bv.Value {
				return -1
			} else if av.Value && !bv.Value {
				return 1
			}
			return 0
		}
	}

	// Fall back to string comparison
	return strings.Compare(a.Inspect(), b.Inspect())
}

// ============================================================================
// Dictionary Methods
// ============================================================================

// evalDictionaryMethod evaluates a method call on a Dictionary
func evalDictionaryMethod(dict *Dictionary, method string, args []Object, env *Environment) Object {
	switch method {
	case "keys":
		if len(args) != 0 {
			return newError("wrong number of arguments for 'keys'. got=%d, want=0", len(args))
		}
		keys := make([]Object, 0, len(dict.Pairs))
		for k := range dict.Pairs {
			// Skip internal fields
			if !strings.HasPrefix(k, "__") {
				keys = append(keys, &String{Value: k})
			}
		}
		return &Array{Elements: keys}

	case "values":
		if len(args) != 0 {
			return newError("wrong number of arguments for 'values'. got=%d, want=0", len(args))
		}
		values := make([]Object, 0, len(dict.Pairs))
		for k, expr := range dict.Pairs {
			// Skip internal fields
			if !strings.HasPrefix(k, "__") {
				val := Eval(expr, dict.Env)
				values = append(values, val)
			}
		}
		return &Array{Elements: values}

	case "has":
		if len(args) != 1 {
			return newError("wrong number of arguments for 'has'. got=%d, want=1", len(args))
		}
		key, ok := args[0].(*String)
		if !ok {
			return newError("argument to 'has' must be STRING, got %s", args[0].Type())
		}
		_, exists := dict.Pairs[key.Value]
		return nativeBoolToParsBoolean(exists)

	default:
		// Return nil for unknown methods to allow user-defined methods to be checked
		return nil
	}
}

// ============================================================================
// Number Methods (Integer and Float)
// ============================================================================

// evalIntegerMethod evaluates a method call on an Integer
func evalIntegerMethod(num *Integer, method string, args []Object) Object {
	switch method {
	case "format":
		// format(locale?)
		if len(args) > 1 {
			return newError("wrong number of arguments for 'format'. got=%d, want=0-1", len(args))
		}
		localeStr := "en-US"
		if len(args) == 1 {
			loc, ok := args[0].(*String)
			if !ok {
				return newError("argument to 'format' must be STRING, got %s", args[0].Type())
			}
			localeStr = loc.Value
		}
		return formatNumberWithLocale(float64(num.Value), localeStr)

	case "currency":
		// currency(code, locale?)
		if len(args) < 1 || len(args) > 2 {
			return newError("wrong number of arguments for 'currency'. got=%d, want=1-2", len(args))
		}
		code, ok := args[0].(*String)
		if !ok {
			return newError("first argument to 'currency' must be STRING, got %s", args[0].Type())
		}
		localeStr := "en-US"
		if len(args) == 2 {
			loc, ok := args[1].(*String)
			if !ok {
				return newError("second argument to 'currency' must be STRING, got %s", args[1].Type())
			}
			localeStr = loc.Value
		}
		return formatCurrencyWithLocale(float64(num.Value), code.Value, localeStr)

	case "percent":
		// percent(locale?)
		if len(args) > 1 {
			return newError("wrong number of arguments for 'percent'. got=%d, want=0-1", len(args))
		}
		localeStr := "en-US"
		if len(args) == 1 {
			loc, ok := args[0].(*String)
			if !ok {
				return newError("argument to 'percent' must be STRING, got %s", args[0].Type())
			}
			localeStr = loc.Value
		}
		return formatPercentWithLocale(float64(num.Value), localeStr)

	default:
		return newError("unknown method '%s' for INTEGER", method)
	}
}

// evalFloatMethod evaluates a method call on a Float
func evalFloatMethod(num *Float, method string, args []Object) Object {
	switch method {
	case "format":
		// format(locale?)
		if len(args) > 1 {
			return newError("wrong number of arguments for 'format'. got=%d, want=0-1", len(args))
		}
		localeStr := "en-US"
		if len(args) == 1 {
			loc, ok := args[0].(*String)
			if !ok {
				return newError("argument to 'format' must be STRING, got %s", args[0].Type())
			}
			localeStr = loc.Value
		}
		return formatNumberWithLocale(num.Value, localeStr)

	case "currency":
		// currency(code, locale?)
		if len(args) < 1 || len(args) > 2 {
			return newError("wrong number of arguments for 'currency'. got=%d, want=1-2", len(args))
		}
		code, ok := args[0].(*String)
		if !ok {
			return newError("first argument to 'currency' must be STRING, got %s", args[0].Type())
		}
		localeStr := "en-US"
		if len(args) == 2 {
			loc, ok := args[1].(*String)
			if !ok {
				return newError("second argument to 'currency' must be STRING, got %s", args[1].Type())
			}
			localeStr = loc.Value
		}
		return formatCurrencyWithLocale(num.Value, code.Value, localeStr)

	case "percent":
		// percent(locale?)
		if len(args) > 1 {
			return newError("wrong number of arguments for 'percent'. got=%d, want=0-1", len(args))
		}
		localeStr := "en-US"
		if len(args) == 1 {
			loc, ok := args[0].(*String)
			if !ok {
				return newError("argument to 'percent' must be STRING, got %s", args[0].Type())
			}
			localeStr = loc.Value
		}
		return formatPercentWithLocale(num.Value, localeStr)

	default:
		return newError("unknown method '%s' for FLOAT", method)
	}
}

// ============================================================================
// Datetime Methods
// ============================================================================

// evalDatetimeMethod evaluates a method call on a datetime dictionary
func evalDatetimeMethod(dict *Dictionary, method string, args []Object, env *Environment) Object {
	switch method {
	case "toDict":
		// toDict() - returns the raw dictionary representation for debugging
		if len(args) != 0 {
			return newError("wrong number of arguments for 'toDict'. got=%d, want=0", len(args))
		}
		return dict

	case "format":
		// format(style?, locale?)
		if len(args) > 2 {
			return newError("wrong number of arguments for 'format'. got=%d, want=0-2", len(args))
		}

		style := "long"
		localeStr := "en-US"

		if len(args) >= 1 {
			styleArg, ok := args[0].(*String)
			if !ok {
				return newError("first argument to 'format' must be STRING (style), got %s", args[0].Type())
			}
			style = styleArg.Value
		}

		if len(args) == 2 {
			locArg, ok := args[1].(*String)
			if !ok {
				return newError("second argument to 'format' must be STRING (locale), got %s", args[1].Type())
			}
			localeStr = locArg.Value
		}

		// Delegate to the formatDate builtin logic
		return formatDateWithStyleAndLocale(dict, style, localeStr, env)

	case "dayOfYear":
		if len(args) != 0 {
			return newError("wrong number of arguments for 'dayOfYear'. got=%d, want=0", len(args))
		}
		return evalDatetimeComputedProperty(dict, "dayOfYear", env)

	case "week":
		if len(args) != 0 {
			return newError("wrong number of arguments for 'week'. got=%d, want=0", len(args))
		}
		return evalDatetimeComputedProperty(dict, "week", env)

	case "timestamp":
		if len(args) != 0 {
			return newError("wrong number of arguments for 'timestamp'. got=%d, want=0", len(args))
		}
		return evalDatetimeComputedProperty(dict, "timestamp", env)

	default:
		return newError("unknown method '%s' for datetime", method)
	}
}

// ============================================================================
// Duration Methods
// ============================================================================

// evalDurationMethod evaluates a method call on a duration dictionary
func evalDurationMethod(dict *Dictionary, method string, args []Object, env *Environment) Object {
	switch method {
	case "toDict":
		// toDict() - returns the raw dictionary representation for debugging
		if len(args) != 0 {
			return newError("wrong number of arguments for 'toDict'. got=%d, want=0", len(args))
		}
		return dict

	case "format":
		// format(locale?)
		if len(args) > 1 {
			return newError("wrong number of arguments for 'format'. got=%d, want=0-1", len(args))
		}

		// Extract months and seconds from duration
		months, seconds, err := getDurationComponents(dict, env)
		if err != nil {
			return newError("invalid duration: %s", err.Error())
		}

		// Get locale (default to en-US)
		localeStr := "en-US"
		if len(args) == 1 {
			locStr, ok := args[0].(*String)
			if !ok {
				return newError("argument to 'format' must be STRING, got %s", args[0].Type())
			}
			localeStr = locStr.Value
		}

		// Format the duration as relative time
		result := locale.DurationToRelativeTime(months, seconds, localeStr)
		return &String{Value: result}

	default:
		return newError("unknown method '%s' for duration", method)
	}
}

// ============================================================================
// Path Methods
// ============================================================================

// evalPathMethod evaluates a method call on a path dictionary
func evalPathMethod(dict *Dictionary, method string, args []Object, env *Environment) Object {
	switch method {
	case "toDict":
		// toDict() - returns the raw dictionary representation for debugging
		if len(args) != 0 {
			return newError("wrong number of arguments for 'toDict'. got=%d, want=0", len(args))
		}
		return dict

	case "isAbsolute":
		if len(args) != 0 {
			return newError("wrong number of arguments for 'isAbsolute'. got=%d, want=0", len(args))
		}
		// Get the absolute property
		if absExpr, ok := dict.Pairs["absolute"]; ok {
			result := Eval(absExpr, env)
			if b, ok := result.(*Boolean); ok {
				return b
			}
		}
		return FALSE

	case "isRelative":
		if len(args) != 0 {
			return newError("wrong number of arguments for 'isRelative'. got=%d, want=0", len(args))
		}
		// Get the absolute property and negate it
		if absExpr, ok := dict.Pairs["absolute"]; ok {
			result := Eval(absExpr, env)
			if b, ok := result.(*Boolean); ok {
				return nativeBoolToParsBoolean(!b.Value)
			}
		}
		return TRUE

	default:
		return newError("unknown method '%s' for path", method)
	}
}

// ============================================================================
// URL Methods
// ============================================================================

// evalUrlMethod evaluates a method call on a URL dictionary
func evalUrlMethod(dict *Dictionary, method string, args []Object, env *Environment) Object {
	switch method {
	case "toDict":
		// toDict() - returns the raw dictionary representation for debugging
		if len(args) != 0 {
			return newError("wrong number of arguments for 'toDict'. got=%d, want=0", len(args))
		}
		return dict

	case "origin":
		if len(args) != 0 {
			return newError("wrong number of arguments for 'origin'. got=%d, want=0", len(args))
		}
		// origin = scheme + "://" + host + (port ? ":" + port : "")
		scheme := ""
		host := ""
		port := ""

		if schemeExpr, ok := dict.Pairs["scheme"]; ok {
			if s := Eval(schemeExpr, env); s != nil {
				if str, ok := s.(*String); ok {
					scheme = str.Value
				}
			}
		}
		if hostExpr, ok := dict.Pairs["host"]; ok {
			if h := Eval(hostExpr, env); h != nil {
				if str, ok := h.(*String); ok {
					host = str.Value
				}
			}
		}
		if portExpr, ok := dict.Pairs["port"]; ok {
			if p := Eval(portExpr, env); p != nil {
				switch pv := p.(type) {
				case *Integer:
					if pv.Value > 0 {
						port = fmt.Sprintf(":%d", pv.Value)
					}
				case *String:
					if pv.Value != "" {
						port = ":" + pv.Value
					}
				}
			}
		}
		return &String{Value: scheme + "://" + host + port}

	case "pathname":
		if len(args) != 0 {
			return newError("wrong number of arguments for 'pathname'. got=%d, want=0", len(args))
		}
		// pathname = "/" + path components joined by "/"
		if pathExpr, ok := dict.Pairs["path"]; ok {
			if p := Eval(pathExpr, env); p != nil {
				if arr, ok := p.(*Array); ok {
					parts := make([]string, 0, len(arr.Elements))
					for _, elem := range arr.Elements {
						if s, ok := elem.(*String); ok && s.Value != "" {
							parts = append(parts, s.Value)
						}
					}
					return &String{Value: "/" + strings.Join(parts, "/")}
				}
			}
		}
		return &String{Value: "/"}

	case "search":
		if len(args) != 0 {
			return newError("wrong number of arguments for 'search'. got=%d, want=0", len(args))
		}
		// search = query string representation
		if queryExpr, ok := dict.Pairs["query"]; ok {
			if q := Eval(queryExpr, env); q != nil {
				if queryDict, ok := q.(*Dictionary); ok {
					if len(queryDict.Pairs) == 0 {
						return &String{Value: ""}
					}
					parts := make([]string, 0, len(queryDict.Pairs))
					for k, v := range queryDict.Pairs {
						if strings.HasPrefix(k, "__") {
							continue
						}
						val := Eval(v, env)
						parts = append(parts, k+"="+val.Inspect())
					}
					return &String{Value: "?" + strings.Join(parts, "&")}
				}
			}
		}
		return &String{Value: ""}

	case "href":
		if len(args) != 0 {
			return newError("wrong number of arguments for 'href'. got=%d, want=0", len(args))
		}
		// href = full URL string representation
		return &String{Value: urlDictToString(dict)}

	default:
		return newError("unknown method '%s' for url", method)
	}
}

// ============================================================================
// Regex Methods
// ============================================================================

// evalRegexMethod evaluates a method call on a regex dictionary
func evalRegexMethod(dict *Dictionary, method string, args []Object, env *Environment) Object {
	switch method {
	case "toDict":
		// toDict() - returns the raw dictionary representation for debugging
		if len(args) != 0 {
			return newError("wrong number of arguments for 'toDict'. got=%d, want=0", len(args))
		}
		return dict

	case "format":
		// format(style?)
		// Styles: "pattern" (just pattern), "literal" (with slashes/flags), "verbose" (pattern and flags separated)
		if len(args) > 1 {
			return newError("wrong number of arguments for 'format'. got=%d, want=0-1", len(args))
		}

		// Get pattern and flags
		var pattern, flags string
		if patternExpr, ok := dict.Pairs["pattern"]; ok {
			if p := Eval(patternExpr, env); p != nil {
				if str, ok := p.(*String); ok {
					pattern = str.Value
				}
			}
		}
		if flagsExpr, ok := dict.Pairs["flags"]; ok {
			if f := Eval(flagsExpr, env); f != nil {
				if str, ok := f.(*String); ok {
					flags = str.Value
				}
			}
		}

		// Get style (default to "literal")
		style := "literal"
		if len(args) == 1 {
			styleArg, ok := args[0].(*String)
			if !ok {
				return newError("argument to 'format' must be STRING (style), got %s", args[0].Type())
			}
			style = styleArg.Value
		}

		switch style {
		case "pattern":
			return &String{Value: pattern}
		case "literal":
			return &String{Value: "/" + pattern + "/" + flags}
		case "verbose":
			if flags == "" {
				return &String{Value: "pattern: " + pattern}
			}
			return &String{Value: "pattern: " + pattern + ", flags: " + flags}
		default:
			return newError("invalid style %q for 'format', use 'pattern', 'literal', or 'verbose'", style)
		}

	case "test":
		// test(string) - returns boolean if the regex matches the string
		if len(args) != 1 {
			return newError("wrong number of arguments for 'test'. got=%d, want=1", len(args))
		}
		str, ok := args[0].(*String)
		if !ok {
			return newError("argument to 'test' must be STRING, got %s", args[0].Type())
		}

		// Get pattern and flags
		var pattern, flags string
		if patternExpr, ok := dict.Pairs["pattern"]; ok {
			if p := Eval(patternExpr, env); p != nil {
				if s, ok := p.(*String); ok {
					pattern = s.Value
				}
			}
		}
		if flagsExpr, ok := dict.Pairs["flags"]; ok {
			if f := Eval(flagsExpr, env); f != nil {
				if s, ok := f.(*String); ok {
					flags = s.Value
				}
			}
		}

		// Compile regex with flags
		re, err := compileRegex(pattern, flags)
		if err != nil {
			return newError("invalid regex pattern: %s", err.Error())
		}

		return nativeBoolToParsBoolean(re.MatchString(str.Value))

	default:
		return newError("unknown method '%s' for regex", method)
	}
}

// ============================================================================
// File Methods
// ============================================================================

// evalFileMethod evaluates a method call on a file dictionary
func evalFileMethod(dict *Dictionary, method string, args []Object, env *Environment) Object {
	switch method {
	case "toDict":
		// toDict() - returns the raw dictionary representation for debugging
		if len(args) != 0 {
			return newError("wrong number of arguments for 'toDict'. got=%d, want=0", len(args))
		}
		return dict

	case "remove":
		// remove() - removes/deletes the file from filesystem
		if len(args) != 0 {
			return newError("wrong number of arguments for 'remove'. got=%d, want=0", len(args))
		}
		return evalFileRemove(dict, env)

	case "mkdir":
		// Create directory
		pathStr := getFilePathString(dict, env)
		if pathStr == "" {
			return newError("file handle has no valid path")
		}

		absPath, pathErr := resolveModulePath(pathStr, env.Filename)
		if pathErr != nil {
			return newError("failed to resolve path '%s': %s", pathStr, pathErr.Error())
		}

		var recursive bool
		if len(args) > 0 {
			if optDict, ok := args[0].(*Dictionary); ok {
				if parentsExpr, ok := optDict.Pairs["parents"]; ok {
					if parentsVal := Eval(parentsExpr, optDict.Env); parentsVal != nil {
						if boolVal, ok := parentsVal.(*Boolean); ok {
							recursive = boolVal.Value
						}
					}
				}
			}
		}

		// Security check (treat as write operation)
		if err := env.checkPathAccess(absPath, "write"); err != nil {
			return newError("security: %s", err.Error())
		}

		var err error
		if recursive {
			err = os.MkdirAll(absPath, 0755)
		} else {
			err = os.Mkdir(absPath, 0755)
		}

		if err != nil {
			return newError("failed to create directory: %s", err.Error())
		}
		return NULL

	case "rmdir":
		// Remove directory
		pathStr := getFilePathString(dict, env)
		if pathStr == "" {
			return newError("file handle has no valid path")
		}

		absPath, pathErr := resolveModulePath(pathStr, env.Filename)
		if pathErr != nil {
			return newError("failed to resolve path '%s': %s", pathStr, pathErr.Error())
		}

		var recursive bool
		if len(args) > 0 {
			if optDict, ok := args[0].(*Dictionary); ok {
				if recursiveExpr, ok := optDict.Pairs["recursive"]; ok {
					if recursiveVal := Eval(recursiveExpr, optDict.Env); recursiveVal != nil {
						if boolVal, ok := recursiveVal.(*Boolean); ok {
							recursive = boolVal.Value
						}
					}
				}
			}
		}

		// Security check (treat as write operation)
		if err := env.checkPathAccess(absPath, "write"); err != nil {
			return newError("security: %s", err.Error())
		}

		var err error
		if recursive {
			err = os.RemoveAll(absPath)
		} else {
			err = os.Remove(absPath)
		}

		if err != nil {
			return newError("failed to remove directory: %s", err.Error())
		}
		return NULL

	default:
		return newError("unknown method '%s' for file", method)
	}
}

// ============================================================================
// Dir Methods
// ============================================================================

// evalDirMethod evaluates a method call on a directory dictionary
func evalDirMethod(dict *Dictionary, method string, args []Object, env *Environment) Object {
	switch method {
	case "toDict":
		// toDict() - returns the raw dictionary representation for debugging
		if len(args) != 0 {
			return newError("wrong number of arguments for 'toDict'. got=%d, want=0", len(args))
		}
		return dict

	case "mkdir":
		// Create directory
		pathStr := getFilePathString(dict, env)
		if pathStr == "" {
			return newError("directory handle has no valid path")
		}

		absPath, pathErr := resolveModulePath(pathStr, env.Filename)
		if pathErr != nil {
			return newError("failed to resolve path '%s': %s", pathStr, pathErr.Error())
		}

		var recursive bool
		if len(args) > 0 {
			if optDict, ok := args[0].(*Dictionary); ok {
				if parentsExpr, ok := optDict.Pairs["parents"]; ok {
					if parentsVal := Eval(parentsExpr, optDict.Env); parentsVal != nil {
						if boolVal, ok := parentsVal.(*Boolean); ok {
							recursive = boolVal.Value
						}
					}
				}
			}
		}

		// Security check (treat as write operation)
		if err := env.checkPathAccess(absPath, "write"); err != nil {
			return newError("security: %s", err.Error())
		}

		var err error
		if recursive {
			err = os.MkdirAll(absPath, 0755)
		} else {
			err = os.Mkdir(absPath, 0755)
		}

		if err != nil {
			return newError("failed to create directory: %s", err.Error())
		}
		return NULL

	case "rmdir":
		// Remove directory
		pathStr := getFilePathString(dict, env)
		if pathStr == "" {
			return newError("directory handle has no valid path")
		}

		absPath, pathErr := resolveModulePath(pathStr, env.Filename)
		if pathErr != nil {
			return newError("failed to resolve path '%s': %s", pathStr, pathErr.Error())
		}

		var recursive bool
		if len(args) > 0 {
			if optDict, ok := args[0].(*Dictionary); ok {
				if recursiveExpr, ok := optDict.Pairs["recursive"]; ok {
					if recursiveVal := Eval(recursiveExpr, optDict.Env); recursiveVal != nil {
						if boolVal, ok := recursiveVal.(*Boolean); ok {
							recursive = boolVal.Value
						}
					}
				}
			}
		}

		// Security check (treat as write operation)
		if err := env.checkPathAccess(absPath, "write"); err != nil {
			return newError("security: %s", err.Error())
		}

		var err error
		if recursive {
			err = os.RemoveAll(absPath)
		} else {
			err = os.Remove(absPath)
		}

		if err != nil {
			return newError("failed to remove directory: %s", err.Error())
		}
		return NULL

	default:
		return newError("unknown method '%s' for dir", method)
	}
}

// ============================================================================
// Request Methods
// ============================================================================

// evalRequestMethod evaluates a method call on a request dictionary
func evalRequestMethod(dict *Dictionary, method string, args []Object, env *Environment) Object {
	switch method {
	case "toDict":
		// toDict() - returns the raw dictionary representation for debugging
		if len(args) != 0 {
			return newError("wrong number of arguments for 'toDict'. got=%d, want=0", len(args))
		}
		return dict

	default:
		return newError("unknown method '%s' for request", method)
	}
}
