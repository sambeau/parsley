package lexer

import (
	"testing"
)

func TestDurationLiterals(t *testing.T) {
	tests := []struct {
		input    string
		expected []Token
	}{
		{
			input: "@2h30m",
			expected: []Token{
				{Type: DURATION_LITERAL, Literal: "2h30m"},
				{Type: EOF, Literal: ""},
			},
		},
		{
			input: "@7d",
			expected: []Token{
				{Type: DURATION_LITERAL, Literal: "7d"},
				{Type: EOF, Literal: ""},
			},
		},
		{
			input: "@1y6mo",
			expected: []Token{
				{Type: DURATION_LITERAL, Literal: "1y6mo"},
				{Type: EOF, Literal: ""},
			},
		},
		{
			input: "@30s",
			expected: []Token{
				{Type: DURATION_LITERAL, Literal: "30s"},
				{Type: EOF, Literal: ""},
			},
		},
		{
			input: "@2024-12-25",
			expected: []Token{
				{Type: DATETIME_LITERAL, Literal: "2024-12-25"},
				{Type: EOF, Literal: ""},
			},
		},
		{
			input: "@1h @2024-12-25",
			expected: []Token{
				{Type: DURATION_LITERAL, Literal: "1h"},
				{Type: DATETIME_LITERAL, Literal: "2024-12-25"},
				{Type: EOF, Literal: ""},
			},
		},
		{
			input: "@3w2d5h",
			expected: []Token{
				{Type: DURATION_LITERAL, Literal: "3w2d5h"},
				{Type: EOF, Literal: ""},
			},
		},
		{
			input: "@1y2mo3w4d5h6m7s",
			expected: []Token{
				{Type: DURATION_LITERAL, Literal: "1y2mo3w4d5h6m7s"},
				{Type: EOF, Literal: ""},
			},
		},
	}

	for _, tt := range tests {
		l := New(tt.input)

		for i, expectedToken := range tt.expected {
			tok := l.NextToken()

			if tok.Type != expectedToken.Type {
				t.Fatalf("tests[%q] - token %d - tokentype wrong. expected=%q, got=%q",
					tt.input, i, expectedToken.Type, tok.Type)
			}

			if tok.Literal != expectedToken.Literal {
				t.Fatalf("tests[%q] - token %d - literal wrong. expected=%q, got=%q",
					tt.input, i, expectedToken.Literal, tok.Literal)
			}
		}
	}
}
