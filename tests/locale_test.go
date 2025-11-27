package main

import (
	"strings"
	"testing"

	"github.com/sambeau/parsley/pkg/evaluator"
	"github.com/sambeau/parsley/pkg/lexer"
	"github.com/sambeau/parsley/pkg/parser"
)

func TestFormatNumber(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		// Default locale (English)
		{`formatNumber(1234567.89)`, "1,234,567.89"},
		{`formatNumber(1234567)`, "1,234,567"},

		// US English
		{`formatNumber(1234567.89, "en-US")`, "1,234,567.89"},

		// German - uses period for thousands, comma for decimal
		{`formatNumber(1234567.89, "de-DE")`, "1.234.567,89"},

		// French - uses space for thousands, comma for decimal
		{`formatNumber(1234567.89, "fr-FR")`, "1 234 567,89"},

		// Indian English - lakh/crore grouping
		{`formatNumber(1234567.89, "en-IN")`, "12,34,567.89"},

		// Spanish
		{`formatNumber(1234567.89, "es-ES")`, "1.234.567,89"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := parser.New(l)
			program := p.ParseProgram()
			env := evaluator.NewEnvironment()
			result := evaluator.Eval(program, env)

			str, ok := result.(*evaluator.String)
			if !ok {
				t.Fatalf("expected String, got %T (%+v)", result, result)
			}
			// Normalize spaces (some locales use non-breaking space)
			actual := strings.ReplaceAll(str.Value, "\u00a0", " ")
			expected := strings.ReplaceAll(tt.expected, "\u00a0", " ")
			if actual != expected {
				t.Errorf("expected '%s', got '%s'", expected, actual)
			}
		})
	}
}

func TestFormatCurrency(t *testing.T) {
	tests := []struct {
		input    string
		contains string // Use contains since exact formatting varies
	}{
		// USD
		{`formatCurrency(1234.56, "USD", "en-US")`, "$"},
		{`formatCurrency(1234.56, "USD", "en-US")`, "1,234.56"},

		// EUR in different locales
		{`formatCurrency(1234.56, "EUR", "de-DE")`, "€"},
		{`formatCurrency(1234.56, "EUR", "de-DE")`, "1.234,56"},

		// GBP
		{`formatCurrency(1234.56, "GBP", "en-GB")`, "£"},

		// JPY (no decimal places) - uses fullwidth yen sign ￥
		{`formatCurrency(1234, "JPY", "ja-JP")`, "￥"},

		// CHF
		{`formatCurrency(1234.56, "CHF", "de-CH")`, "CHF"},
	}

	for _, tt := range tests {
		t.Run(tt.input+"_contains_"+tt.contains, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := parser.New(l)
			program := p.ParseProgram()
			env := evaluator.NewEnvironment()
			result := evaluator.Eval(program, env)

			str, ok := result.(*evaluator.String)
			if !ok {
				t.Fatalf("expected String, got %T (%+v)", result, result)
			}
			if !strings.Contains(str.Value, tt.contains) {
				t.Errorf("expected to contain '%s', got '%s'", tt.contains, str.Value)
			}
		})
	}
}

func TestFormatPercent(t *testing.T) {
	tests := []struct {
		input    string
		contains string
	}{
		// Basic percentage
		{`formatPercent(0.12)`, "12"},
		{`formatPercent(0.12)`, "%"},

		// US format
		{`formatPercent(0.1234, "en-US")`, "12"},

		// German format (space before %)
		{`formatPercent(0.1234, "de-DE")`, "%"},

		// Turkish (% before number)
		{`formatPercent(0.1234, "tr-TR")`, "%"},
	}

	for _, tt := range tests {
		t.Run(tt.input+"_contains_"+tt.contains, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := parser.New(l)
			program := p.ParseProgram()
			env := evaluator.NewEnvironment()
			result := evaluator.Eval(program, env)

			str, ok := result.(*evaluator.String)
			if !ok {
				t.Fatalf("expected String, got %T (%+v)", result, result)
			}
			if !strings.Contains(str.Value, tt.contains) {
				t.Errorf("expected to contain '%s', got '%s'", tt.contains, str.Value)
			}
		})
	}
}

func TestFormatNumberErrors(t *testing.T) {
	tests := []struct {
		input       string
		errContains string
	}{
		{`formatNumber("not a number")`, "must be INTEGER or FLOAT"},
		{`formatNumber(123, 456)`, "must be STRING"},
		{`formatNumber(123, "invalid-locale-xyz")`, "invalid locale"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := parser.New(l)
			program := p.ParseProgram()
			env := evaluator.NewEnvironment()
			result := evaluator.Eval(program, env)

			err, ok := result.(*evaluator.Error)
			if !ok {
				t.Fatalf("expected Error, got %T (%+v)", result, result)
			}
			if !strings.Contains(err.Message, tt.errContains) {
				t.Errorf("expected error to contain '%s', got '%s'", tt.errContains, err.Message)
			}
		})
	}
}

func TestFormatCurrencyErrors(t *testing.T) {
	tests := []struct {
		input       string
		errContains string
	}{
		{`formatCurrency("not a number", "USD")`, "must be INTEGER or FLOAT"},
		{`formatCurrency(123, 456)`, "must be STRING"},
		{`formatCurrency(123, "INVALID")`, "invalid currency code"},
		{`formatCurrency(123, "USD", "invalid-locale-xyz")`, "invalid locale"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := parser.New(l)
			program := p.ParseProgram()
			env := evaluator.NewEnvironment()
			result := evaluator.Eval(program, env)

			err, ok := result.(*evaluator.Error)
			if !ok {
				t.Fatalf("expected Error, got %T (%+v)", result, result)
			}
			if !strings.Contains(err.Message, tt.errContains) {
				t.Errorf("expected error to contain '%s', got '%s'", tt.errContains, err.Message)
			}
		})
	}
}

func TestFormatDate(t *testing.T) {
	tests := []struct {
		input    string
		contains string
	}{
		// US English formats
		{`let d = time({year: 2024, month: 12, day: 25}); formatDate(d)`, "December 25, 2024"},
		{`let d = time({year: 2024, month: 12, day: 25}); formatDate(d, "short")`, "12/25/24"},
		{`let d = time({year: 2024, month: 12, day: 25}); formatDate(d, "medium")`, "Dec 25, 2024"},
		{`let d = time({year: 2024, month: 12, day: 25}); formatDate(d, "long")`, "December 25, 2024"},
		{`let d = time({year: 2024, month: 12, day: 25}); formatDate(d, "full")`, "December 25, 2024"},

		// French formats
		{`let d = time({year: 2024, month: 12, day: 25}); formatDate(d, "long", "fr-FR")`, "25 décembre 2024"},
		{`let d = time({year: 2024, month: 12, day: 25}); formatDate(d, "medium", "fr-FR")`, "25 déc 2024"},
		{`let d = time({year: 2024, month: 12, day: 25}); formatDate(d, "full", "fr-FR")`, "mercredi"},

		// German formats
		{`let d = time({year: 2024, month: 12, day: 25}); formatDate(d, "long", "de-DE")`, "25. Dezember 2024"},
		{`let d = time({year: 2024, month: 12, day: 25}); formatDate(d, "short", "de-DE")`, "25.12.24"},
		{`let d = time({year: 2024, month: 12, day: 25}); formatDate(d, "full", "de-DE")`, "Mittwoch"},

		// Japanese formats
		{`let d = time({year: 2024, month: 12, day: 25}); formatDate(d, "long", "ja-JP")`, "2024年12月25日"},

		// Spanish formats
		{`let d = time({year: 2024, month: 12, day: 25}); formatDate(d, "long", "es-ES")`, "25 de diciembre de 2024"},
		{`let d = time({year: 2024, month: 12, day: 25}); formatDate(d, "full", "es-ES")`, "miércoles"},
	}

	for _, tt := range tests {
		t.Run(tt.input+"_contains_"+tt.contains, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := parser.New(l)
			program := p.ParseProgram()
			env := evaluator.NewEnvironment()
			result := evaluator.Eval(program, env)

			str, ok := result.(*evaluator.String)
			if !ok {
				t.Fatalf("expected String, got %T (%+v)", result, result)
			}
			if !strings.Contains(str.Value, tt.contains) {
				t.Errorf("expected to contain '%s', got '%s'", tt.contains, str.Value)
			}
		})
	}
}

func TestFormatDateErrors(t *testing.T) {
	tests := []struct {
		input       string
		errContains string
	}{
		{`formatDate("not a date")`, "must be a datetime"},
		{`formatDate({})`, "must be a datetime"},
		{`let d = time({year: 2024, month: 12, day: 25}); formatDate(d, 123)`, "must be STRING"},
		{`let d = time({year: 2024, month: 12, day: 25}); formatDate(d, "long", 456)`, "must be STRING"},
		{`let d = time({year: 2024, month: 12, day: 25}); formatDate(d, "invalid")`, "must be one of: short, medium, long, full"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := parser.New(l)
			program := p.ParseProgram()
			env := evaluator.NewEnvironment()
			result := evaluator.Eval(program, env)

			err, ok := result.(*evaluator.Error)
			if !ok {
				t.Fatalf("expected Error, got %T (%+v)", result, result)
			}
			if !strings.Contains(err.Message, tt.errContains) {
				t.Errorf("expected error to contain '%s', got '%s'", tt.errContains, err.Message)
			}
		})
	}
}
