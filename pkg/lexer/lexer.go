package lexer

import (
	"fmt"
	"unicode"
	"unicode/utf8"
)

// TokenType represents different types of tokens
type TokenType int

const (
	// Special tokens
	ILLEGAL TokenType = iota
	EOF

	// Identifiers and literals
	IDENT    // add, foobar, x, y, ...
	INT      // 1343456
	FLOAT    // 3.14159
	STRING   // "foobar"
	TEMPLATE // `template ${expr}`

	// Operators
	ASSIGN   // =
	PLUS     // +
	MINUS    // -
	BANG     // !
	ASTERISK // *
	SLASH    // /
	LT       // <
	GT       // >
	LTE      // <=
	GTE      // >=
	EQ       // ==
	NOT_EQ   // !=
	AND      // & or and
	OR       // | or or

	// Delimiters
	COMMA     // ,
	SEMICOLON // ;
	COLON     // :
	LPAREN    // (
	RPAREN    // )
	LBRACE    // {
	RBRACE    // }
	LBRACKET  // [
	RBRACKET  // ]
	PLUSPLUS  // ++

	// Keywords
	FUNCTION // "fn"
	LET      // "let"
	FOR      // "for"
	IN       // "in"
	TRUE     // "true"
	FALSE    // "false"
	IF       // "if"
	THEN     // "then"
	ELSE     // "else"
	RETURN   // "return"
)

// Token represents a single token
type Token struct {
	Type    TokenType
	Literal string
	Line    int
	Column  int
}

// String returns a string representation of the token
func (t Token) String() string {
	return fmt.Sprintf("{Type: %s, Literal: %s, Line: %d, Column: %d}",
		t.Type.String(), t.Literal, t.Line, t.Column)
}

// String returns a string representation of the token type
func (tt TokenType) String() string {
	switch tt {
	case ILLEGAL:
		return "ILLEGAL"
	case EOF:
		return "EOF"
	case IDENT:
		return "IDENT"
	case INT:
		return "INT"
	case FLOAT:
		return "FLOAT"
	case STRING:
		return "STRING"
	case TEMPLATE:
		return "TEMPLATE"
	case ASSIGN:
		return "ASSIGN"
	case PLUS:
		return "PLUS"
	case MINUS:
		return "MINUS"
	case BANG:
		return "BANG"
	case ASTERISK:
		return "ASTERISK"
	case SLASH:
		return "SLASH"
	case LT:
		return "LT"
	case GT:
		return "GT"
	case LTE:
		return "LTE"
	case GTE:
		return "GTE"
	case EQ:
		return "EQ"
	case NOT_EQ:
		return "NOT_EQ"
	case AND:
		return "AND"
	case OR:
		return "OR"
	case COMMA:
		return "COMMA"
	case SEMICOLON:
		return "SEMICOLON"
	case COLON:
		return "COLON"
	case LPAREN:
		return "LPAREN"
	case RPAREN:
		return "RPAREN"
	case LBRACE:
		return "LBRACE"
	case RBRACE:
		return "RBRACE"
	case LBRACKET:
		return "LBRACKET"
	case RBRACKET:
		return "RBRACKET"
	case PLUSPLUS:
		return "PLUSPLUS"
	case FUNCTION:
		return "FUNCTION"
	case LET:
		return "LET"
	case FOR:
		return "FOR"
	case IN:
		return "IN"
	case TRUE:
		return "TRUE"
	case FALSE:
		return "FALSE"
	case IF:
		return "IF"
	case THEN:
		return "THEN"
	case ELSE:
		return "ELSE"
	case RETURN:
		return "RETURN"
	default:
		return "UNKNOWN"
	}
}

// Keywords map for identifying language keywords
var keywords = map[string]TokenType{
	"fn":     FUNCTION,
	"let":    LET,
	"for":    FOR,
	"in":     IN,
	"true":   TRUE,
	"false":  FALSE,
	"if":     IF,
	"then":   THEN,
	"else":   ELSE,
	"return": RETURN,
	"and":    AND,
	"or":     OR,
	"not":    BANG,
}

// LookupIdent checks if an identifier is a keyword
func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}

// Lexer represents the lexical analyzer
type Lexer struct {
	filename     string
	input        string
	position     int  // current position in input (points to current char)
	readPosition int  // current reading position in input (after current char)
	ch           byte // current char under examination
	line         int  // current line number
	column       int  // current column number
}

// New creates a new lexer instance
func New(input string) *Lexer {
	l := &Lexer{
		filename: "<input>",
		input:    input,
		line:     1,
		column:   0,
	}
	l.readChar()
	return l
}

// NewWithFilename creates a new lexer instance with a specific filename
func NewWithFilename(input string, filename string) *Lexer {
	l := &Lexer{
		filename: filename,
		input:    input,
		line:     1,
		column:   0,
	}
	l.readChar()
	return l
}

// readChar reads the next character and advances position
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0 // ASCII NUL character represents EOF
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++

	if l.ch == '\n' {
		l.line++
		l.column = 0
	} else {
		l.column++
	}
}

// peekChar returns the next character without advancing position
func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

// NextToken scans the input and returns the next token
func (l *Lexer) NextToken() Token {
	var tok Token

	l.skipWhitespace()

	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = Token{Type: EQ, Literal: string(ch) + string(l.ch), Line: l.line, Column: l.column - 1}
		} else {
			tok = newToken(ASSIGN, l.ch, l.line, l.column)
		}
	case '+':
		if l.peekChar() == '+' {
			ch := l.ch
			l.readChar()
			tok = Token{Type: PLUSPLUS, Literal: string(ch) + string(l.ch), Line: l.line, Column: l.column - 1}
		} else {
			tok = newToken(PLUS, l.ch, l.line, l.column)
		}
	case '-':
		tok = newToken(MINUS, l.ch, l.line, l.column)
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = Token{Type: NOT_EQ, Literal: string(ch) + string(l.ch), Line: l.line, Column: l.column - 1}
		} else {
			tok = newToken(BANG, l.ch, l.line, l.column)
		}
	case '/':
		if l.peekChar() == '/' {
			l.skipComment()
			return l.NextToken()
		}
		tok = newToken(SLASH, l.ch, l.line, l.column)
	case '*':
		tok = newToken(ASTERISK, l.ch, l.line, l.column)
	case '<':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = Token{Type: LTE, Literal: string(ch) + string(l.ch), Line: l.line, Column: l.column - 1}
		} else {
			tok = newToken(LT, l.ch, l.line, l.column)
		}
	case '>':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = Token{Type: GTE, Literal: string(ch) + string(l.ch), Line: l.line, Column: l.column - 1}
		} else {
			tok = newToken(GT, l.ch, l.line, l.column)
		}
	case '&':
		tok = newToken(AND, l.ch, l.line, l.column)
	case '|':
		tok = newToken(OR, l.ch, l.line, l.column)
	case ';':
		tok = newToken(SEMICOLON, l.ch, l.line, l.column)
	case ',':
		tok = newToken(COMMA, l.ch, l.line, l.column)
	case ':':
		tok = newToken(COLON, l.ch, l.line, l.column)
	case '[':
		tok = newToken(LBRACKET, l.ch, l.line, l.column)
	case ']':
		tok = newToken(RBRACKET, l.ch, l.line, l.column)
	case '(':
		tok = newToken(LPAREN, l.ch, l.line, l.column)
	case ')':
		tok = newToken(RPAREN, l.ch, l.line, l.column)
	case '{':
		tok = newToken(LBRACE, l.ch, l.line, l.column)
	case '}':
		tok = newToken(RBRACE, l.ch, l.line, l.column)
	case '"':
		line := l.line
		column := l.column
		tok.Type = STRING
		tok.Literal = l.readString()
		tok.Line = line
		tok.Column = column
	case '`':
		line := l.line
		column := l.column
		tok.Type = TEMPLATE
		tok.Literal = l.readTemplate()
		tok.Line = line
		tok.Column = column
	case 0:
		tok.Literal = ""
		tok.Type = EOF
		tok.Line = l.line
		tok.Column = l.column
	default:
		if isLetter(l.ch) {
			// Save position before reading
			line := l.line
			column := l.column
			tok.Literal = l.readIdentifier()
			tok.Type = LookupIdent(tok.Literal)
			tok.Line = line
			tok.Column = column
			return tok // early return to avoid readChar()
		} else if isDigit(l.ch) {
			// Save position before reading
			line := l.line
			column := l.column
			tok.Literal = l.readNumber()
			// Check if it's a float or integer
			if containsDot(tok.Literal) {
				tok.Type = FLOAT
			} else {
				tok.Type = INT
			}
			tok.Line = line
			tok.Column = column
			return tok // early return to avoid readChar()
		} else {
			tok = newToken(ILLEGAL, l.ch, l.line, l.column)
		}
	}

	l.readChar()
	return tok
}

// newToken creates a new token with the given parameters
func newToken(tokenType TokenType, ch byte, line, column int) Token {
	return Token{Type: tokenType, Literal: string(ch), Line: line, Column: column}
}

// readIdentifier reads an identifier or keyword
func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) || isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

// readNumber reads a number (integer or float)
func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}

	// Check for decimal point
	if l.ch == '.' && isDigit(l.peekChar()) {
		l.readChar() // consume the '.'
		for isDigit(l.ch) {
			l.readChar()
		}
	}

	return l.input[position:l.position]
}

// readString reads a string literal with escape sequence support
func (l *Lexer) readString() string {
	var result []byte
	l.readChar() // skip opening quote

	for l.ch != '"' && l.ch != 0 {
		if l.ch == '\\' {
			l.readChar() // consume backslash
			switch l.ch {
			case 'n':
				result = append(result, '\n')
			case 't':
				result = append(result, '\t')
			case '\\':
				result = append(result, '\\')
			case '"':
				result = append(result, '"')
			default:
				// Unknown escape, keep as-is
				result = append(result, '\\')
				result = append(result, l.ch)
			}
		} else {
			result = append(result, l.ch)
		}
		l.readChar()
	}

	return string(result)
}

// readTemplate reads a template literal (backtick string)
func (l *Lexer) readTemplate() string {
	var result []byte
	l.readChar() // skip opening backtick

	for l.ch != '`' && l.ch != 0 {
		if l.ch == '\\' {
			l.readChar() // consume backslash
			switch l.ch {
			case '`':
				result = append(result, '`')
			case '$':
				// Use a special marker that evaluator won't interpret
				result = append(result, '\\', '0', '$') // \0$ as escape marker
			default:
				// Unknown escape, keep as-is
				result = append(result, '\\')
				result = append(result, l.ch)
			}
		} else {
			result = append(result, l.ch)
		}
		l.readChar()
	}

	return string(result)
}

// skipWhitespace skips whitespace characters
func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

// skipComment skips single-line comments starting with //
func (l *Lexer) skipComment() {
	// Skip the two slashes
	l.readChar()
	l.readChar()

	// Read until end of line or EOF
	for l.ch != '\n' && l.ch != 0 {
		l.readChar()
	}
}

// isLetter checks if the character is a letter
func isLetter(ch byte) bool {
	r, _ := utf8.DecodeRune([]byte{ch})
	return unicode.IsLetter(r) || ch == '_'
}

// isDigit checks if the character is a digit
func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

// containsDot checks if a string contains a decimal point
func containsDot(s string) bool {
	for _, ch := range s {
		if ch == '.' {
			return true
		}
	}
	return false
}
