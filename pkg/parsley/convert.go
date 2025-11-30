package parsley

import (
	"fmt"
	"reflect"
	"time"

	"github.com/sambeau/parsley/pkg/evaluator"
)

// ToParsley converts a Go value to a Parsley Object.
//
// Supported conversions:
//   - int, int8, int16, int32, int64  → Integer
//   - uint, uint8, uint16, uint32, uint64 → Integer
//   - float32, float64 → Float
//   - string → String
//   - bool → Boolean
//   - []byte → String (as raw bytes)
//   - []interface{} or []T → Array
//   - map[string]interface{} or map[string]T → Dictionary
//   - time.Time → DateTime dictionary
//   - time.Duration → Duration dictionary
//   - nil → Null
func ToParsley(v interface{}) (evaluator.Object, error) {
	if v == nil {
		return evaluator.NULL, nil
	}

	switch val := v.(type) {
	case int:
		return &evaluator.Integer{Value: int64(val)}, nil
	case int8:
		return &evaluator.Integer{Value: int64(val)}, nil
	case int16:
		return &evaluator.Integer{Value: int64(val)}, nil
	case int32:
		return &evaluator.Integer{Value: int64(val)}, nil
	case int64:
		return &evaluator.Integer{Value: val}, nil
	case uint:
		return &evaluator.Integer{Value: int64(val)}, nil
	case uint8:
		return &evaluator.Integer{Value: int64(val)}, nil
	case uint16:
		return &evaluator.Integer{Value: int64(val)}, nil
	case uint32:
		return &evaluator.Integer{Value: int64(val)}, nil
	case uint64:
		return &evaluator.Integer{Value: int64(val)}, nil
	case float32:
		return &evaluator.Float{Value: float64(val)}, nil
	case float64:
		return &evaluator.Float{Value: val}, nil
	case string:
		return &evaluator.String{Value: val}, nil
	case bool:
		if val {
			return evaluator.TRUE, nil
		}
		return evaluator.FALSE, nil
	case []byte:
		return &evaluator.String{Value: string(val)}, nil
	case time.Time:
		return timeToParsley(val), nil
	case time.Duration:
		return durationToParsley(val), nil
	case []interface{}:
		return sliceToParsley(val)
	case map[string]interface{}:
		return mapToParsley(val)
	default:
		// Use reflection for other slice and map types
		rv := reflect.ValueOf(v)
		switch rv.Kind() {
		case reflect.Slice, reflect.Array:
			return reflectSliceToParsley(rv)
		case reflect.Map:
			return reflectMapToParsley(rv)
		case reflect.Ptr:
			if rv.IsNil() {
				return evaluator.NULL, nil
			}
			return ToParsley(rv.Elem().Interface())
		}
		return nil, fmt.Errorf("unsupported Go type: %T", v)
	}
}

// FromParsley converts a Parsley Object to a Go value.
//
// Returns:
//   - Integer → int64
//   - Float → float64
//   - String → string
//   - Boolean → bool
//   - Array → []interface{}
//   - Dictionary → map[string]interface{}
//   - Null → nil
func FromParsley(obj evaluator.Object) interface{} {
	if obj == nil {
		return nil
	}

	switch o := obj.(type) {
	case *evaluator.Integer:
		return o.Value
	case *evaluator.Float:
		return o.Value
	case *evaluator.String:
		return o.Value
	case *evaluator.Boolean:
		return o.Value
	case *evaluator.Null:
		return nil
	case *evaluator.Array:
		return arrayToGo(o)
	case *evaluator.Dictionary:
		return dictToGo(o)
	case *evaluator.Error:
		return fmt.Errorf("%s", o.Message)
	default:
		// For other types, return the string representation
		return obj.Inspect()
	}
}

// Helper functions

func timeToParsley(t time.Time) *evaluator.Dictionary {
	pairs := make(map[string]evaluator.Object)
	pairs["year"] = &evaluator.Integer{Value: int64(t.Year())}
	pairs["month"] = &evaluator.Integer{Value: int64(t.Month())}
	pairs["day"] = &evaluator.Integer{Value: int64(t.Day())}
	pairs["hour"] = &evaluator.Integer{Value: int64(t.Hour())}
	pairs["minute"] = &evaluator.Integer{Value: int64(t.Minute())}
	pairs["second"] = &evaluator.Integer{Value: int64(t.Second())}
	pairs["weekday"] = &evaluator.Integer{Value: int64(t.Weekday())}
	pairs["yearday"] = &evaluator.Integer{Value: int64(t.YearDay())}
	pairs["kind"] = &evaluator.String{Value: "datetime"}

	// Store Unix timestamp for operations
	pairs["unix"] = &evaluator.Integer{Value: t.Unix()}

	return evaluator.NewDictionaryFromObjects(pairs)
}

func durationToParsley(d time.Duration) *evaluator.Dictionary {
	pairs := make(map[string]evaluator.Object)

	// Total values
	totalSeconds := int64(d.Seconds())
	totalMinutes := int64(d.Minutes())
	totalHours := int64(d.Hours())
	totalDays := totalHours / 24

	// Component values
	hours := totalHours % 24
	minutes := totalMinutes % 60
	seconds := totalSeconds % 60

	pairs["days"] = &evaluator.Integer{Value: totalDays}
	pairs["hours"] = &evaluator.Integer{Value: hours}
	pairs["minutes"] = &evaluator.Integer{Value: minutes}
	pairs["seconds"] = &evaluator.Integer{Value: seconds}
	pairs["totalHours"] = &evaluator.Integer{Value: totalHours}
	pairs["totalMinutes"] = &evaluator.Integer{Value: totalMinutes}
	pairs["totalSeconds"] = &evaluator.Integer{Value: totalSeconds}
	pairs["kind"] = &evaluator.String{Value: "duration"}

	return evaluator.NewDictionaryFromObjects(pairs)
}

func sliceToParsley(slice []interface{}) (*evaluator.Array, error) {
	elements := make([]evaluator.Object, len(slice))
	for i, v := range slice {
		obj, err := ToParsley(v)
		if err != nil {
			return nil, fmt.Errorf("element %d: %w", i, err)
		}
		elements[i] = obj
	}
	return &evaluator.Array{Elements: elements}, nil
}

func reflectSliceToParsley(rv reflect.Value) (*evaluator.Array, error) {
	elements := make([]evaluator.Object, rv.Len())
	for i := 0; i < rv.Len(); i++ {
		obj, err := ToParsley(rv.Index(i).Interface())
		if err != nil {
			return nil, fmt.Errorf("element %d: %w", i, err)
		}
		elements[i] = obj
	}
	return &evaluator.Array{Elements: elements}, nil
}

func mapToParsley(m map[string]interface{}) (*evaluator.Dictionary, error) {
	pairs := make(map[string]evaluator.Object)
	for k, v := range m {
		obj, err := ToParsley(v)
		if err != nil {
			return nil, fmt.Errorf("key %q: %w", k, err)
		}
		pairs[k] = obj
	}
	return evaluator.NewDictionaryFromObjects(pairs), nil
}

func reflectMapToParsley(rv reflect.Value) (*evaluator.Dictionary, error) {
	if rv.Type().Key().Kind() != reflect.String {
		return nil, fmt.Errorf("map keys must be strings, got %s", rv.Type().Key())
	}

	pairs := make(map[string]evaluator.Object)
	iter := rv.MapRange()
	for iter.Next() {
		k := iter.Key().String()
		v := iter.Value().Interface()
		obj, err := ToParsley(v)
		if err != nil {
			return nil, fmt.Errorf("key %q: %w", k, err)
		}
		pairs[k] = obj
	}
	return evaluator.NewDictionaryFromObjects(pairs), nil
}

func arrayToGo(arr *evaluator.Array) []interface{} {
	result := make([]interface{}, len(arr.Elements))
	for i, elem := range arr.Elements {
		result[i] = FromParsley(elem)
	}
	return result
}

func dictToGo(dict *evaluator.Dictionary) map[string]interface{} {
	result := make(map[string]interface{})

	// Dictionaries store expressions that need to be evaluated
	// We create a temporary environment to evaluate them
	env := evaluator.NewEnvironment()

	if dict.Pairs != nil {
		for k, expr := range dict.Pairs {
			// Evaluate the expression to get the actual value
			obj := evaluator.Eval(expr, env)
			result[k] = FromParsley(obj)
		}
	}

	return result
}
