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
		{`formatNumber("not a number")`, "must be an integer or float"},
		{`formatNumber(123, 456)`, "must be a string"},
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
		{`formatCurrency("not a number", "USD")`, "must be an integer or float"},
		{`formatCurrency(123, 456)`, "must be a string"},
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
		{`let d = time({year: 2024, month: 12, day: 25}); formatDate(d, 123)`, "must be a string"},
		{`let d = time({year: 2024, month: 12, day: 25}); formatDate(d, "long", 456)`, "must be a string"},
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

func TestNegativeDurationLiterals(t *testing.T) {
	tests := []struct {
		input           string
		expectedSeconds int64
	}{
		{`@-1d`, -86400},
		{`@-2d`, -172800},
		{`@-1w`, -604800},
		{`@-2w`, -1209600},
		{`@-3h`, -10800},
		{`@-30m`, -1800},
		{`@-3h30m`, -12600},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := parser.New(l)
			program := p.ParseProgram()
			env := evaluator.NewEnvironment()
			result := evaluator.Eval(program, env)

			dict, ok := result.(*evaluator.Dictionary)
			if !ok {
				t.Fatalf("expected Dictionary, got %T (%+v)", result, result)
			}

			// Check __type
			typeExpr, ok := dict.Pairs["__type"]
			if !ok {
				t.Fatal("expected __type field")
			}
			typeObj := evaluator.Eval(typeExpr, env)
			typeStr, ok := typeObj.(*evaluator.String)
			if !ok || typeStr.Value != "duration" {
				t.Fatalf("expected __type='duration', got %v", typeObj)
			}

			// Check seconds
			secondsExpr, ok := dict.Pairs["seconds"]
			if !ok {
				t.Fatal("expected seconds field")
			}
			secondsObj := evaluator.Eval(secondsExpr, env)
			secondsInt, ok := secondsObj.(*evaluator.Integer)
			if !ok {
				t.Fatalf("expected Integer for seconds, got %T", secondsObj)
			}
			if secondsInt.Value != tt.expectedSeconds {
				t.Errorf("expected seconds=%d, got %d", tt.expectedSeconds, secondsInt.Value)
			}
		})
	}
}

func TestFormatDuration(t *testing.T) {
	tests := []struct {
		input    string
		contains string
	}{
		// English (default)
		{`format(@1d)`, "tomorrow"},
		{`format(@-1d)`, "yesterday"},
		{`format(@2d)`, "in 2 days"},
		{`format(@-2d)`, "2 days ago"},
		{`format(@1w)`, "next week"},
		{`format(@-1w)`, "last week"},
		{`format(@3h)`, "in 3 hours"},
		{`format(@-3h)`, "3 hours ago"},

		// German
		{`format(@1d, "de-DE")`, "morgen"},
		{`format(@-1d, "de-DE")`, "gestern"},
		{`format(@-2d, "de-DE")`, "vorgestern"},
		{`format(@2w, "de-DE")`, "in 2 Wochen"},

		// French
		{`format(@1d, "fr-FR")`, "demain"},
		{`format(@-1d, "fr-FR")`, "hier"},
		{`format(@-2d, "fr-FR")`, "avant-hier"},

		// Spanish
		{`format(@1d, "es-ES")`, "mañana"},
		{`format(@-1d, "es-ES")`, "ayer"},
		{`format(@-2d, "es-ES")`, "anteayer"},

		// Japanese
		{`format(@1d, "ja-JP")`, "明日"},
		{`format(@-1d, "ja-JP")`, "昨日"},
		{`format(@-2d, "ja-JP")`, "一昨日"},

		// Russian
		{`format(@1d, "ru-RU")`, "завтра"},
		{`format(@-1d, "ru-RU")`, "вчера"},
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

func TestFormatDurationErrors(t *testing.T) {
	tests := []struct {
		input       string
		errContains string
	}{
		{`format("not a duration")`, "must be a duration or array"},
		{`format({})`, "must be a duration"},
		{`format(@1d, 123)`, "must be a string"},
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

// ============================================================================
// List Formatting Tests (Phase 4)
// ============================================================================

func TestFormatListEnglish(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		// Empty and single-item lists
		{`format([])`, ""},
		{`format(["apple"])`, "apple"},

		// Two-item lists
		{`format(["apple", "banana"])`, "apple and banana"},
		{`format(["apple", "banana"], "or")`, "apple or banana"},

		// Three-item lists (with Oxford comma for en-US)
		{`format(["apple", "banana", "cherry"])`, "apple, banana, and cherry"},
		{`format(["apple", "banana", "cherry"], "or")`, "apple, banana, or cherry"},

		// Four-item lists
		{`format(["apple", "banana", "cherry", "date"])`, "apple, banana, cherry, and date"},

		// Unit style (no conjunction)
		{`format(["5 feet", "6 inches"], "unit")`, "5 feet, 6 inches"},
		{`format(["1 hour", "30 minutes", "15 seconds"], "unit")`, "1 hour, 30 minutes, 15 seconds"},

		// Non-string elements get converted
		{`format([1, 2, 3])`, "1, 2, and 3"},
		{`format([true, false])`, "true and false"},
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
			if str.Value != tt.expected {
				t.Errorf("expected '%s', got '%s'", tt.expected, str.Value)
			}
		})
	}
}

func TestFormatListEnglishGB(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		// No Oxford comma for en-GB
		{`format(["apple", "banana", "cherry"], "and", "en-GB")`, "apple, banana and cherry"},
		{`format(["apple", "banana", "cherry"], "or", "en-GB")`, "apple, banana or cherry"},
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
			if str.Value != tt.expected {
				t.Errorf("expected '%s', got '%s'", tt.expected, str.Value)
			}
		})
	}
}

func TestFormatListGerman(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`format(["Apfel", "Banane"], "and", "de-DE")`, "Apfel und Banane"},
		{`format(["Apfel", "Banane", "Kirsche"], "and", "de-DE")`, "Apfel, Banane und Kirsche"},
		{`format(["Apfel", "Banane"], "or", "de-DE")`, "Apfel oder Banane"},
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
			if str.Value != tt.expected {
				t.Errorf("expected '%s', got '%s'", tt.expected, str.Value)
			}
		})
	}
}

func TestFormatListFrench(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`format(["pomme", "banane"], "and", "fr-FR")`, "pomme et banane"},
		{`format(["pomme", "banane", "cerise"], "and", "fr-FR")`, "pomme, banane et cerise"},
		{`format(["pomme", "banane"], "or", "fr-FR")`, "pomme ou banane"},
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
			if str.Value != tt.expected {
				t.Errorf("expected '%s', got '%s'", tt.expected, str.Value)
			}
		})
	}
}

func TestFormatListJapanese(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		// Japanese uses different separators
		{`format(["りんご", "バナナ"], "and", "ja-JP")`, "りんごとバナナ"},
		{`format(["りんご", "バナナ", "さくらんぼ"], "and", "ja-JP")`, "りんご、バナナ、さくらんぼ"},
		{`format(["りんご", "バナナ"], "or", "ja-JP")`, "りんごまたはバナナ"},
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
			if str.Value != tt.expected {
				t.Errorf("expected '%s', got '%s'", tt.expected, str.Value)
			}
		})
	}
}

func TestFormatListChinese(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`format(["苹果", "香蕉"], "and", "zh-CN")`, "苹果和香蕉"},
		{`format(["苹果", "香蕉", "樱桃"], "and", "zh-CN")`, "苹果、香蕉和樱桃"},
		{`format(["苹果", "香蕉"], "or", "zh-CN")`, "苹果或香蕉"},
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
			if str.Value != tt.expected {
				t.Errorf("expected '%s', got '%s'", tt.expected, str.Value)
			}
		})
	}
}

func TestFormatListRussian(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`format(["яблоко", "банан"], "and", "ru-RU")`, "яблоко и банан"},
		{`format(["яблоко", "банан", "вишня"], "and", "ru-RU")`, "яблоко, банан и вишня"},
		{`format(["яблоко", "банан"], "or", "ru-RU")`, "яблоко или банан"},
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
			if str.Value != tt.expected {
				t.Errorf("expected '%s', got '%s'", tt.expected, str.Value)
			}
		})
	}
}

func TestFormatListErrors(t *testing.T) {
	tests := []struct {
		input       string
		errContains string
	}{
		{`format(["a", "b"], 123)`, "must be a string"},
		{`format(["a", "b"], "invalid")`, "invalid style"},
		{`format(["a", "b"], "and", 123)`, "must be a string"},
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
