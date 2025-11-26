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
	IDENT     // add, foobar, x, y, ...
	INT       // 1343456
	FLOAT     // 3.14159
	STRING           // "foobar"
	TEMPLATE         // `template ${expr}`
	REGEX            // /pattern/flags
	DATETIME_LITERAL // @2024-12-25T14:30:00Z
	TAG              // <tag prop="value" />
	TAG_START // <tag> or <tag attr="value">
	TAG_END   // </tag>
	TAG_TEXT  // raw text content within tags

	// Operators
	ASSIGN    // =
	PLUS      // +
	MINUS     // -
	BANG      // !
	ASTERISK  // *
	SLASH     // /
	PERCENT   // %
	LT        // <
	GT        // >
	LTE       // <=
	GTE       // >=
	EQ        // ==
	NOT_EQ    // !=
	AND       // & or and
	OR        // | or or
	MATCH     // ~
	NOT_MATCH // !~

	// Delimiters
	COMMA     // ,
	SEMICOLON // ;
	COLON     // :
	DOT       // .
	DOTDOTDOT // ... (spread/rest operator)
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
	AS       // "as"
	TRUE     // "true"
	FALSE    // "false"
	IF       // "if"
	ELSE     // "else"
	RETURN   // "return"
	DELETE   // "delete"
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
	case REGEX:
		return "REGEX"
	case DATETIME_LITERAL:
		return "DATETIME_LITERAL"
	case TAG:
		return "TAG"
	case TAG_START:
		return "TAG_START"
	case TAG_END:
		return "TAG_END"
	case TAG_TEXT:
		return "TAG_TEXT"
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
	case PERCENT:
		return "PERCENT"
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
	case MATCH:
		return "MATCH"
	case NOT_MATCH:
		return "NOT_MATCH"
	case COMMA:
		return "COMMA"
	case SEMICOLON:
		return "SEMICOLON"
	case COLON:
		return "COLON"
	case DOT:
		return "DOT"
	case DOTDOTDOT:
		return "DOTDOTDOT"
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
	case AS:
		return "AS"
	case TRUE:
		return "TRUE"
	case FALSE:
		return "FALSE"
	case IF:
		return "IF"
	case ELSE:
		return "ELSE"
	case RETURN:
		return "RETURN"
	case DELETE:
		return "DELETE"
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
	"as":     AS,
	"true":   TRUE,
	"false":  FALSE,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
	"delete": DELETE,
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
	filename      string
	input         string
	position      int       // current position in input (points to current char)
	readPosition  int       // current reading position in input (after current char)
	ch            byte      // current char under examination
	line          int       // current line number
	column        int       // current column number
	inTagContent  bool      // whether we're currently lexing tag content
	tagDepth      int       // nesting depth of tags (for proper TAG_END matching)
	lastTokenType TokenType // last token type for regex context detection
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

	// Special handling when inside tag content
	if l.inTagContent {
		return l.nextTagContentToken()
	}

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
		} else if l.peekChar() == '~' {
			ch := l.ch
			l.readChar()
			tok = Token{Type: NOT_MATCH, Literal: string(ch) + string(l.ch), Line: l.line, Column: l.column - 1}
		} else {
			tok = newToken(BANG, l.ch, l.line, l.column)
		}
	case '~':
		tok = newToken(MATCH, l.ch, l.line, l.column)
	case '/':
		if l.peekChar() == '/' {
			l.skipComment()
			return l.NextToken()
		} else if l.shouldTreatAsRegex(l.lastTokenType) {
			// This is a regex literal
			line := l.line
			column := l.column
			pattern, flags := l.readRegex()
			tok.Type = REGEX
			tok.Literal = "/" + pattern + "/" + flags
			tok.Line = line
			tok.Column = column
			l.lastTokenType = tok.Type
			return tok
		}
		tok = newToken(SLASH, l.ch, l.line, l.column)
	case '%':
		tok = newToken(PERCENT, l.ch, l.line, l.column)
	case '*':
		tok = newToken(ASTERISK, l.ch, l.line, l.column)
	case '<':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = Token{Type: LTE, Literal: string(ch) + string(l.ch), Line: l.line, Column: l.column - 1}
		} else if l.peekChar() == '/' {
			// This is a closing tag </tag>
			line := l.line
			column := l.column
			tok.Type = TAG_END
			tok.Literal = l.readTagEnd()
			tok.Line = line
			tok.Column = column
			l.lastTokenType = tok.Type
			return tok
		} else if isLetter(l.peekChar()) || l.peekChar() == '>' {
			// Could be a tag start <tag> or singleton <tag />
			line := l.line
			column := l.column
			tagContent, isSingleton := l.readTagStartOrSingleton()
			if isSingleton {
				tok.Type = TAG
			} else {
				tok.Type = TAG_START
				// Enter tag content mode
				l.inTagContent = true
				l.tagDepth = 1
			}
			tok.Literal = tagContent
			tok.Line = line
			tok.Column = column
			l.lastTokenType = tok.Type
			return tok
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
	case '.':
		// Check for "..." (spread/rest operator)
		if l.peekChar() == '.' && l.readPosition+1 < len(l.input) && l.input[l.readPosition+1] == '.' {
			line := l.line
			col := l.column
			l.readChar() // consume second '.'
			l.readChar() // consume third '.'
			tok = Token{Type: DOTDOTDOT, Literal: "...", Line: line, Column: col}
		} else {
			tok = newToken(DOT, l.ch, l.line, l.column)
		}
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
	case '@':
		line := l.line
		column := l.column
		tok.Type = DATETIME_LITERAL
		tok.Literal = l.readDatetimeLiteral()
		tok.Line = line
		tok.Column = column
		l.lastTokenType = tok.Type
		return tok
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
			l.lastTokenType = tok.Type
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
			l.lastTokenType = tok.Type
			return tok // early return to avoid readChar()
		} else {
			tok = newToken(ILLEGAL, l.ch, l.line, l.column)
		}
	}

	l.readChar()
	l.lastTokenType = tok.Type
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
			case '{':
				// Use a special marker that evaluator won't interpret
				result = append(result, '\\', '0', '{') // \0{ as escape marker
			case '}':
				// Use a special marker that evaluator won't interpret
				result = append(result, '\\', '0', '}') // \0} as escape marker
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

// readTag reads a singleton tag like <input type="text" />
func (l *Lexer) readTag() string {
	var result []byte
	l.readChar() // skip opening <

	// Read until we find />
	for {
		if l.ch == 0 {
			// Unexpected EOF
			break
		}

		// Check for closing />
		if l.ch == '/' && l.peekChar() == '>' {
			l.readChar() // consume /
			l.readChar() // consume >
			break
		}

		// Handle string literals within the tag
		if l.ch == '"' {
			result = append(result, l.ch)
			l.readChar()
			// Read until closing quote
			for l.ch != '"' && l.ch != 0 {
				if l.ch == '\\' {
					result = append(result, l.ch)
					l.readChar()
					if l.ch != 0 {
						result = append(result, l.ch)
						l.readChar()
					}
				} else {
					result = append(result, l.ch)
					l.readChar()
				}
			}
			if l.ch == '"' {
				result = append(result, l.ch)
				l.readChar()
			}
			continue
		}

		// Handle interpolation braces {}
		if l.ch == '{' {
			result = append(result, l.ch)
			l.readChar()
			braceDepth := 1
			// Read until matching closing brace
			for braceDepth > 0 && l.ch != 0 {
				if l.ch == '{' {
					braceDepth++
				} else if l.ch == '}' {
					braceDepth--
				} else if l.ch == '"' {
					// Handle string inside interpolation
					result = append(result, l.ch)
					l.readChar()
					for l.ch != '"' && l.ch != 0 {
						if l.ch == '\\' {
							result = append(result, l.ch)
							l.readChar()
							if l.ch != 0 {
								result = append(result, l.ch)
								l.readChar()
							}
							continue
						}
						result = append(result, l.ch)
						l.readChar()
					}
					if l.ch == '"' {
						result = append(result, l.ch)
						l.readChar()
					}
					continue
				}
				result = append(result, l.ch)
				l.readChar()
			}
			continue
		}

		result = append(result, l.ch)
		l.readChar()
	}

	return string(result)
}

// readTagEnd reads a closing tag like </div>
func (l *Lexer) readTagEnd() string {
	var result []byte
	l.readChar() // skip <
	l.readChar() // skip /

	// Read tag name
	for isLetter(l.ch) || isDigit(l.ch) {
		result = append(result, l.ch)
		l.readChar()
	}

	// Skip whitespace before >
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}

	// Expect >
	if l.ch != '>' {
		// Error: expected >
		return string(result)
	}
	l.readChar() // consume >

	return string(result)
}

// readTagStartOrSingleton reads a tag start <tag> or singleton <tag />
// Returns the tag content and a boolean indicating if it's a singleton
func (l *Lexer) readTagStartOrSingleton() (string, bool) {
	var result []byte
	l.readChar() // skip opening <

	// Check for empty grouping tag <>
	if l.ch == '>' {
		l.readChar()     // consume >
		return "", false // empty grouping tag is a tag start
	}

	// Read until we find > or />
	isSingleton := false
	for {
		if l.ch == 0 {
			// Unexpected EOF
			break
		}

		// Check for closing />
		if l.ch == '/' && l.peekChar() == '>' {
			l.readChar() // consume /
			l.readChar() // consume >
			isSingleton = true
			break
		}

		// Check for just >
		if l.ch == '>' {
			l.readChar() // consume >
			break
		}

		// Handle string literals within the tag
		if l.ch == '"' {
			result = append(result, l.ch)
			l.readChar()
			// Read until closing quote
			for l.ch != '"' && l.ch != 0 {
				if l.ch == '\\' {
					result = append(result, l.ch)
					l.readChar()
					if l.ch != 0 {
						result = append(result, l.ch)
						l.readChar()
					}
				} else {
					result = append(result, l.ch)
					l.readChar()
				}
			}
			if l.ch == '"' {
				result = append(result, l.ch)
				l.readChar()
			}
			continue
		}

		// Handle interpolation braces {}
		if l.ch == '{' {
			result = append(result, l.ch)
			l.readChar()
			braceDepth := 1
			// Read until matching closing brace
			for braceDepth > 0 && l.ch != 0 {
				if l.ch == '{' {
					braceDepth++
				} else if l.ch == '}' {
					braceDepth--
				} else if l.ch == '"' {
					// Handle string inside interpolation
					result = append(result, l.ch)
					l.readChar()
					for l.ch != '"' && l.ch != 0 {
						if l.ch == '\\' {
							result = append(result, l.ch)
							l.readChar()
							if l.ch != 0 {
								result = append(result, l.ch)
								l.readChar()
							}
							continue
						}
						result = append(result, l.ch)
						l.readChar()
					}
					if l.ch == '"' {
						result = append(result, l.ch)
						l.readChar()
					}
					continue
				}
				result = append(result, l.ch)
				l.readChar()
			}
			continue
		}

		result = append(result, l.ch)
		l.readChar()
	}

	return string(result), isSingleton
}

// nextTagContentToken returns the next token while in tag content mode
func (l *Lexer) nextTagContentToken() Token {
	var tok Token

	// Consolidate multiple newlines into single spaces
	if l.ch == '\n' || l.ch == '\r' {
		for l.ch == '\n' || l.ch == '\r' || l.ch == ' ' || l.ch == '\t' {
			l.readChar()
		}
		// Don't return whitespace token, just continue to next content
		if l.ch == 0 || l.ch == '<' || l.ch == '{' {
			// Fall through to handle these special cases
		} else {
			// Start with a space before the next text
			line := l.line
			column := l.column
			text := l.readTagText()
			tok = Token{Type: TAG_TEXT, Literal: " " + text, Line: line, Column: column}
			return tok
		}
	}

	line := l.line
	column := l.column

	switch l.ch {
	case 0:
		tok = Token{Type: EOF, Literal: "", Line: l.line, Column: l.column}
		l.inTagContent = false

	case '<':
		if l.peekChar() == '/' {
			// Closing tag
			tok.Type = TAG_END
			tok.Literal = l.readTagEnd()
			tok.Line = line
			tok.Column = column
			l.tagDepth--
			if l.tagDepth == 0 {
				l.inTagContent = false
			}
			return tok
		} else if l.peekChar() == '>' {
			// Empty grouping tag start
			l.readChar() // skip <
			l.readChar() // skip >
			tok = Token{Type: TAG_START, Literal: "", Line: line, Column: column}
			l.tagDepth++
			return tok
		} else if isLetter(l.peekChar()) {
			// Nested tag (start or singleton)
			tagContent, isSingleton := l.readTagStartOrSingleton()
			if isSingleton {
				tok.Type = TAG
			} else {
				tok.Type = TAG_START
				l.tagDepth++
			}
			tok.Literal = tagContent
			tok.Line = line
			tok.Column = column
			return tok
		} else {
			// Literal < character in content
			tok.Type = TAG_TEXT
			tok.Literal = string(l.ch)
			tok.Line = line
			tok.Column = column
			l.readChar()
			return tok
		}

	case '{':
		// Interpolation - temporarily exit tag content mode
		tok = newToken(LBRACE, l.ch, l.line, l.column)
		l.readChar()
		l.inTagContent = false
		return tok

	default:
		// Regular text content
		tok.Type = TAG_TEXT
		tok.Literal = l.readTagText()
		tok.Line = line
		tok.Column = column
	}

	return tok
}

// readTagText reads text content until we hit <, {, or EOF
func (l *Lexer) readTagText() string {
	var result []byte

	for l.ch != 0 && l.ch != '<' && l.ch != '{' {
		result = append(result, l.ch)
		l.readChar()
	}

	return string(result)
}

// EnterTagContentMode sets the lexer into tag content mode
func (l *Lexer) EnterTagContentMode() {
	if l.tagDepth > 0 {
		l.inTagContent = true
	}
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

// readRegex reads a regex literal like /pattern/flags
func (l *Lexer) readRegex() (string, string) {
	var pattern []byte
	l.readChar() // skip opening /

	// Read pattern until we find unescaped /
	for l.ch != '/' && l.ch != 0 && l.ch != '\n' {
		if l.ch == '\\' {
			pattern = append(pattern, l.ch)
			l.readChar()
			if l.ch != 0 {
				pattern = append(pattern, l.ch)
				l.readChar()
			}
		} else {
			pattern = append(pattern, l.ch)
			l.readChar()
		}
	}

	if l.ch != '/' {
		// Invalid regex (unterminated)
		return string(pattern), ""
	}

	l.readChar() // consume closing /

	// Read flags (letters immediately after closing /)
	var flags []byte
	for isLetter(l.ch) {
		flags = append(flags, l.ch)
		l.readChar()
	}

	// Back up one char since we consumed one too many
	l.position = l.readPosition - 1
	if l.position < len(l.input) {
		l.ch = l.input[l.position]
	} else {
		l.ch = 0
	}

	return string(pattern), string(flags)
}

// readDatetimeLiteral reads a datetime literal after @
// Supports formats: @2024-12-25, @2024-12-25T14:30:00, @2024-12-25T14:30:00Z, @2024-12-25T14:30:00-05:00
func (l *Lexer) readDatetimeLiteral() string {
	var datetime []byte
	l.readChar() // skip @

	// Read date part: YYYY-MM-DD
	for isDigit(l.ch) || l.ch == '-' {
		datetime = append(datetime, l.ch)
		l.readChar()
	}
	
	// Check for time part: T14:30:00
	if l.ch == 'T' {
		datetime = append(datetime, l.ch)
		l.readChar()
		// Read time: HH:MM:SS
		for isDigit(l.ch) || l.ch == ':' {
			datetime = append(datetime, l.ch)
			l.readChar()
		}
	}
	
	// Check for fractional seconds (.123)
	if l.ch == '.' && isDigit(l.peekChar()) {
		datetime = append(datetime, l.ch)
		l.readChar()
		// Read the fractional part
		for isDigit(l.ch) {
			datetime = append(datetime, l.ch)
			l.readChar()
		}
	}
	
	// Check for timezone: Z or +05:00 or -05:00
	if l.ch == 'Z' {
		datetime = append(datetime, l.ch)
		l.readChar()
	} else if l.ch == '+' || l.ch == '-' {
		// Only consume if followed by digit (timezone offset)
		if isDigit(l.peekChar()) {
			datetime = append(datetime, l.ch)
			l.readChar()
			// Read timezone offset
			for isDigit(l.ch) || l.ch == ':' {
				datetime = append(datetime, l.ch)
				l.readChar()
			}
		}
	}

	// Back up one char since we consumed one too many
	l.position = l.readPosition - 1
	if l.position < len(l.input) {
		l.ch = l.input[l.position]
	} else {
		l.ch = 0
	}

	return string(datetime)
}

// shouldTreatAsRegex determines if / should be regex or division
// Regex context: after operators, keywords, commas, open parens/brackets
// But NOT after complete expressions like identifiers, numbers, close parens
func (l *Lexer) shouldTreatAsRegex(lastToken TokenType) bool {
	switch lastToken {
	case ASSIGN, EQ, NOT_EQ, LT, GT, LTE, GTE,
		AND, OR, MATCH, NOT_MATCH,
		LPAREN, LBRACKET, LBRACE,
		COMMA, SEMICOLON, COLON,
		RETURN, LET, IF, ELSE, FOR, IN,
		PLUSPLUS:
		return true
	case 0: // Start of input
		return true
	// Don't treat as regex after arithmetic operators that could be infix
	// These appear in expressions like: x - /y/ which is (x - /) then /y/
	// Instead they're more likely: x-1/2 (division)
	default:
		return false
	}
}
