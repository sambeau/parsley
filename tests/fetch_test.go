package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/sambeau/parsley/pkg/evaluator"
)

// ============================================================================
// URL Builtin Tests
// ============================================================================

func TestUrlBuiltin(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		check    func(evaluator.Object) bool
		errorMsg string
	}{
		{
			name:  "basic URL parsing",
			input: `url("https://example.com")`,
			check: func(obj evaluator.Object) bool {
				dict, ok := obj.(*evaluator.Dictionary)
				if !ok {
					return false
				}
				// Check it has URL fields
				_, hasScheme := dict.Pairs["scheme"]
				_, hasHost := dict.Pairs["host"]
				return hasScheme && hasHost
			},
		},
		{
			name:  "URL with path",
			input: `url("https://example.com/api/users")`,
			check: func(obj evaluator.Object) bool {
				dict, ok := obj.(*evaluator.Dictionary)
				return ok && dict != nil
			},
		},
		{
			name:  "URL with query params",
			input: `url("https://example.com/search?q=test&page=1")`,
			check: func(obj evaluator.Object) bool {
				dict, ok := obj.(*evaluator.Dictionary)
				return ok && dict != nil
			},
		},
		{
			name:  "URL with port",
			input: `url("https://example.com:8080/api")`,
			check: func(obj evaluator.Object) bool {
				dict, ok := obj.(*evaluator.Dictionary)
				return ok && dict != nil
			},
		},
		{
			name:     "URL with missing argument",
			input:    `url()`,
			errorMsg: "wrong number of arguments",
		},
		{
			name:     "URL with wrong type",
			input:    `url(123)`,
			errorMsg: "must be a string",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := testEvalHelper(tt.input)

			if tt.errorMsg != "" {
				errObj, ok := result.(*evaluator.Error)
				if !ok {
					t.Fatalf("expected error, got %T (%+v)", result, result)
				}
				if !strings.Contains(errObj.Message, tt.errorMsg) {
					t.Errorf("error message should contain %q, got %q", tt.errorMsg, errObj.Message)
				}
				return
			}

			if tt.check != nil && !tt.check(result) {
				t.Errorf("check failed for input %q, got %T (%+v)", tt.input, result, result)
			}
		})
	}
}

// ============================================================================
// Request Handle Tests (JSON, YAML format factories)
// ============================================================================

func TestRequestHandleCreation(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		check    func(evaluator.Object) bool
		errorMsg string
	}{
		{
			name:  "JSON request handle from URL",
			input: `JSON(url("https://api.example.com/data"))`,
			check: func(obj evaluator.Object) bool {
				dict, ok := obj.(*evaluator.Dictionary)
				if !ok {
					return false
				}
				// Should be a request dictionary with __type field
				_, hasType := dict.Pairs["__type"]
				_, hasFormat := dict.Pairs["format"]
				return hasType && hasFormat
			},
		},
		{
			name:  "YAML request handle from URL",
			input: `YAML(url("https://api.example.com/config.yaml"))`,
			check: func(obj evaluator.Object) bool {
				dict, ok := obj.(*evaluator.Dictionary)
				return ok && dict != nil
			},
		},
		{
			name:  "text request handle from URL",
			input: `text(url("https://example.com/page.html"))`,
			check: func(obj evaluator.Object) bool {
				dict, ok := obj.(*evaluator.Dictionary)
				return ok && dict != nil
			},
		},
		{
			name:  "lines request handle from URL",
			input: `lines(url("https://example.com/data.txt"))`,
			check: func(obj evaluator.Object) bool {
				dict, ok := obj.(*evaluator.Dictionary)
				return ok && dict != nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := testEvalHelper(tt.input)

			if tt.errorMsg != "" {
				errObj, ok := result.(*evaluator.Error)
				if !ok {
					t.Fatalf("expected error, got %T (%+v)", result, result)
				}
				if !strings.Contains(errObj.Message, tt.errorMsg) {
					t.Errorf("error message should contain %q, got %q", tt.errorMsg, errObj.Message)
				}
				return
			}

			if tt.check != nil && !tt.check(result) {
				t.Errorf("check failed for input %q, got %T (%+v)", tt.input, result, result)
			}
		})
	}
}

// ============================================================================
// Live HTTP Fetch Tests (using httptest)
// ============================================================================

func TestFetchOperatorWithMockServer(t *testing.T) {
	// Create a test server that returns JSON
	jsonServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"name": "test", "value": 42}`))
	}))
	defer jsonServer.Close()

	// Create a test server that returns plain text
	textServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello, World!"))
	}))
	defer textServer.Close()

	// Create a test server that returns an error
	errorServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
	}))
	defer errorServer.Close()

	// Create a test server that echoes the request method
	methodServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"method": "` + r.Method + `"}`))
	}))
	defer methodServer.Close()

	t.Run("fetch JSON with error capture pattern", func(t *testing.T) {
		input := `{data, error} <=/= JSON(url("` + jsonServer.URL + `")); data`
		result := testEvalHelper(input)

		// Should return a dictionary with name and value
		dict, ok := result.(*evaluator.Dictionary)
		if !ok {
			t.Fatalf("expected Dictionary, got %T (%+v)", result, result)
		}

		// Check for expected fields
		if _, hasName := dict.Pairs["name"]; !hasName {
			t.Error("expected 'name' field in response")
		}
	})

	t.Run("fetch plain text", func(t *testing.T) {
		input := `{data, error} <=/= text(url("` + textServer.URL + `")); data`
		result := testEvalHelper(input)

		str, ok := result.(*evaluator.String)
		if !ok {
			t.Fatalf("expected String, got %T (%+v)", result, result)
		}
		if str.Value != "Hello, World!" {
			t.Errorf("expected 'Hello, World!', got %q", str.Value)
		}
	})

	t.Run("fetch with error capture - server error", func(t *testing.T) {
		input := `{data, error, status} <=/= text(url("` + errorServer.URL + `")); status`
		result := testEvalHelper(input)

		// Should return status code
		num, ok := result.(*evaluator.Integer)
		if !ok {
			t.Fatalf("expected Integer (status code), got %T (%+v)", result, result)
		}
		if num.Value != 500 {
			t.Errorf("expected status 500, got %d", num.Value)
		}
	})

	t.Run("fetch error returns null for data", func(t *testing.T) {
		// When fetching from a server that returns 500, data should still be the content
		input := `{data, error, status} <=/= text(url("` + errorServer.URL + `")); data`
		result := testEvalHelper(input)

		// Data should be the error text content
		str, ok := result.(*evaluator.String)
		if !ok {
			t.Fatalf("expected String for data, got %T (%+v)", result, result)
		}
		if str.Value != "Internal Server Error" {
			t.Errorf("expected 'Internal Server Error', got %q", str.Value)
		}
	})
}

// ============================================================================
// Fetch Error Handling Tests
// ============================================================================

func TestFetchErrorHandling(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		errorMsg string
	}{
		{
			name:     "fetch with invalid source type",
			input:    `data <=/= 123`,
			errorMsg: "fetch operator",
		},
		{
			name:     "fetch with string instead of URL",
			input:    `data <=/= "https://example.com"`,
			errorMsg: "fetch operator",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := testEvalHelper(tt.input)

			errObj, ok := result.(*evaluator.Error)
			if !ok {
				t.Fatalf("expected error, got %T (%+v)", result, result)
			}
			if !strings.Contains(errObj.Message, tt.errorMsg) {
				t.Errorf("error message should contain %q, got %q", tt.errorMsg, errObj.Message)
			}
		})
	}
}

// ============================================================================
// URL Computed Properties Tests
// ============================================================================

func TestURLComputedProperties(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "URL scheme property",
			input:    `u = url("https://example.com"); u.scheme`,
			expected: "https",
		},
		{
			name:     "URL host property",
			input:    `u = url("https://example.com"); u.host`,
			expected: "example.com",
		},
		{
			name:     "URL href property",
			input:    `u = url("https://example.com/path"); u.href`,
			expected: "https://example.com/path",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := testEvalHelper(tt.input)

			str, ok := result.(*evaluator.String)
			if !ok {
				t.Fatalf("expected String, got %T (%+v)", result, result)
			}
			if str.Value != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, str.Value)
			}
		})
	}
}

// ============================================================================
// Request Options Tests
// ============================================================================

func TestRequestWithOptions(t *testing.T) {
	// Server that echoes request details
	echoServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		// Echo back relevant request info
		w.Write([]byte(`{
			"method": "` + r.Method + `",
			"contentType": "` + r.Header.Get("Content-Type") + `"
		}`))
	}))
	defer echoServer.Close()

	t.Run("basic GET request", func(t *testing.T) {
		input := `{data, error} <=/= JSON(url("` + echoServer.URL + `")); data.method`
		result := testEvalHelper(input)

		str, ok := result.(*evaluator.String)
		if !ok {
			t.Fatalf("expected String, got %T (%+v)", result, result)
		}
		if str.Value != "GET" {
			t.Errorf("expected 'GET', got %q", str.Value)
		}
	})
}

// ============================================================================
// Lines Format Tests
// ============================================================================

func TestFetchLinesFormat(t *testing.T) {
	// Server that returns multiple lines
	linesServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("line1\nline2\nline3"))
	}))
	defer linesServer.Close()

	t.Run("fetch as lines array", func(t *testing.T) {
		input := `{data, error} <=/= lines(url("` + linesServer.URL + `")); data.length()`
		result := testEvalHelper(input)

		num, ok := result.(*evaluator.Integer)
		if !ok {
			t.Fatalf("expected Integer, got %T (%+v)", result, result)
		}
		if num.Value != 3 {
			t.Errorf("expected 3 lines, got %d", num.Value)
		}
	})
}
