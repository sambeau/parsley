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
	IDENT            // add, foobar, x, y, ...
	INT              // 1343456
	FLOAT            // 3.14159
	STRING           // "foobar"
	TEMPLATE         // `template ${expr}`
	REGEX            // /pattern/flags
	DATETIME_LITERAL // @2024-12-25T14:30:00Z
	DURATION_LITERAL // @2h30m, @7d, @1y6mo
	PATH_LITERAL     // @/usr/local, @./config
	URL_LITERAL      // @https://example.com
	TAG              // <tag prop="value" />
	TAG_START        // <tag> or <tag attr="value">
	TAG_END          // </tag>
	TAG_TEXT         // raw text content within tags

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
	NULLISH   // ??
	MATCH     // ~
	NOT_MATCH // !~

	// File I/O operators
	READ_FROM // <==
	WRITE_TO  // ==>
	APPEND_TO // ==>>

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
	case DURATION_LITERAL:
		return "DURATION_LITERAL"
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
	case NULLISH:
		return "NULLISH"
	case MATCH:
		return "MATCH"
	case NOT_MATCH:
		return "NOT_MATCH"
	case READ_FROM:
		return "READ_FROM"
	case WRITE_TO:
		return "WRITE_TO"
	case APPEND_TO:
		return "APPEND_TO"
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
	inRawTextTag  string    // non-empty when inside <style> or <script> - stores tag name (for @{} mode)
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

// peekCharN returns the character n positions ahead without advancing position
func (l *Lexer) peekCharN(n int) byte {
	pos := l.readPosition + n - 1
	if pos >= len(l.input) {
		return 0
	}
	return l.input[pos]
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
			line := l.line
			col := l.column
			l.readChar() // consume second '='
			if l.peekChar() == '>' {
				l.readChar() // consume '>'
				if l.peekChar() == '>' {
					l.readChar() // consume second '>'
					tok = Token{Type: APPEND_TO, Literal: "==>>", Line: line, Column: col}
				} else {
					tok = Token{Type: WRITE_TO, Literal: "==>", Line: line, Column: col}
				}
			} else {
				tok = Token{Type: EQ, Literal: string(ch) + string(l.ch), Line: line, Column: col}
			}
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
		if l.peekChar() == '=' && l.peekCharN(2) == '=' {
			// <== (read from file)
			line := l.line
			col := l.column
			l.readChar() // consume first '='
			l.readChar() // consume second '='
			tok = Token{Type: READ_FROM, Literal: "<==", Line: line, Column: col}
		} else if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = Token{Type: LTE, Literal: string(ch) + string(l.ch), Line: l.line, Column: l.column - 1}
		} else if l.peekChar() == '?' {
			// XML processing instruction <?xml ... ?> - pass through as string
			line := l.line
			column := l.column
			content := l.readProcessingInstruction()
			tok.Type = STRING
			tok.Literal = content
			tok.Line = line
			tok.Column = column
			l.lastTokenType = tok.Type
			return tok
		} else if l.peekChar() == '!' {
			// Could be XML comment <!-- -->, CDATA <![CDATA[]]>, or DOCTYPE <!DOCTYPE>
			if l.peekCharN(2) == '-' && l.peekCharN(3) == '-' {
				// XML comment - skip it and get next token
				l.skipXMLComment()
				return l.NextToken()
			} else if l.peekCharN(2) == '[' && l.peekCharN(3) == 'C' {
				// CDATA section - return as string
				line := l.line
				column := l.column
				content, ok := l.readCDATA()
				if ok {
					tok.Type = STRING
					tok.Literal = content
					tok.Line = line
					tok.Column = column
					l.lastTokenType = tok.Type
					return tok
				}
			} else if l.peekCharN(2) == 'D' || l.peekCharN(2) == 'd' {
				// DOCTYPE declaration - pass through as string
				line := l.line
				column := l.column
				content := l.readDoctype()
				tok.Type = STRING
				tok.Literal = content
				tok.Line = line
				tok.Column = column
				l.lastTokenType = tok.Type
				return tok
			}
			// Not a comment, CDATA, or DOCTYPE - treat < as less-than
			tok = newToken(LT, l.ch, l.line, l.column)
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
				// Check if this is a raw text tag (style or script)
				tagName := extractTagName(tagContent)
				if tagName == "style" || tagName == "script" {
					l.inRawTextTag = tagName
				}
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
	case '?':
		if l.peekChar() == '?' {
			ch := l.ch
			l.readChar()
			tok = Token{Type: NULLISH, Literal: string(ch) + string(l.ch), Line: l.line, Column: l.column - 1}
		} else {
			tok = newToken(ILLEGAL, l.ch, l.line, l.column)
		}
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
		// Peek ahead to determine the literal type
		literalType := l.detectAtLiteralType()
		switch literalType {
		case DATETIME_LITERAL:
			tok.Type = DATETIME_LITERAL
			tok.Literal = l.readDatetimeLiteral()
		case DURATION_LITERAL:
			tok.Type = DURATION_LITERAL
			tok.Literal = l.readDurationLiteral()
		case PATH_LITERAL:
			tok.Type = PATH_LITERAL
			tok.Literal = l.readPathLiteral()
		case URL_LITERAL:
			tok.Type = URL_LITERAL
			tok.Literal = l.readUrlLiteral()
		default:
			tok.Type = ILLEGAL
			tok.Literal = string(l.ch)
		}
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

// skipXMLComment skips an XML comment <!-- ... -->
// Returns true if a comment was successfully skipped
func (l *Lexer) skipXMLComment() bool {
	// We're at '<', peek for '!--'
	if l.peekChar() != '!' || l.peekCharN(2) != '-' || l.peekCharN(3) != '-' {
		return false
	}

	// Skip the <!--
	l.readChar() // skip <
	l.readChar() // skip !
	l.readChar() // skip -
	l.readChar() // skip -

	// Read until we find -->
	for {
		if l.ch == 0 {
			// Unexpected EOF
			break
		}
		if l.ch == '-' && l.peekChar() == '-' && l.peekCharN(2) == '>' {
			l.readChar() // skip -
			l.readChar() // skip -
			l.readChar() // skip >
			break
		}
		l.readChar()
	}

	return true
}

// readCDATA reads a CDATA section <![CDATA[ ... ]]> and returns its content
func (l *Lexer) readCDATA() (string, bool) {
	// We're at '<', check for '![CDATA['
	if l.peekChar() != '!' || l.peekCharN(2) != '[' || l.peekCharN(3) != 'C' ||
		l.peekCharN(4) != 'D' || l.peekCharN(5) != 'A' || l.peekCharN(6) != 'T' ||
		l.peekCharN(7) != 'A' || l.peekCharN(8) != '[' {
		return "", false
	}

	// Skip the <![CDATA[
	for i := 0; i < 9; i++ {
		l.readChar()
	}

	var content []byte
	// Read until we find ]]>
	for {
		if l.ch == 0 {
			// Unexpected EOF
			break
		}
		if l.ch == ']' && l.peekChar() == ']' && l.peekCharN(2) == '>' {
			l.readChar() // skip ]
			l.readChar() // skip ]
			l.readChar() // skip >
			break
		}
		content = append(content, l.ch)
		l.readChar()
	}

	return string(content), true
}

// readProcessingInstruction reads a processing instruction <?...?>
// Returns the full content including delimiters
func (l *Lexer) readProcessingInstruction() string {
	var result []byte
	result = append(result, '<', '?')
	l.readChar() // skip <
	l.readChar() // skip ?

	// Read until we find ?>
	for {
		if l.ch == 0 {
			// Unexpected EOF
			break
		}
		if l.ch == '?' && l.peekChar() == '>' {
			result = append(result, '?', '>')
			l.readChar() // skip ?
			l.readChar() // skip >
			break
		}
		result = append(result, l.ch)
		l.readChar()
	}

	return string(result)
}

// readDoctype reads a DOCTYPE declaration <!DOCTYPE...>
// Returns the full content including delimiters
func (l *Lexer) readDoctype() string {
	var result []byte
	result = append(result, '<', '!')
	l.readChar() // skip <
	l.readChar() // skip !

	// Read until we find >
	for {
		if l.ch == 0 {
			// Unexpected EOF
			break
		}
		if l.ch == '>' {
			result = append(result, '>')
			l.readChar() // skip >
			break
		}
		result = append(result, l.ch)
		l.readChar()
	}

	return string(result)
}

// readTagEnd reads a closing tag like </div> or </my-component>
func (l *Lexer) readTagEnd() string {
	var result []byte
	l.readChar() // skip <
	l.readChar() // skip /

	// Read tag name (allow hyphens for web components like my-component)
	for isLetter(l.ch) || isDigit(l.ch) || l.ch == '-' {
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

	// In raw text mode (style/script), check for @{ which triggers interpolation
	inRawMode := l.inRawTextTag != ""

	// Consolidate multiple newlines into single spaces (but not in raw mode)
	if !inRawMode && (l.ch == '\n' || l.ch == '\r') {
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
			closingTagName := l.readTagEnd()
			tok.Literal = closingTagName
			tok.Line = line
			tok.Column = column
			l.tagDepth--
			if l.tagDepth == 0 {
				l.inTagContent = false
			}
			// If we're closing a raw text tag, exit raw text mode
			if l.inRawTextTag != "" && closingTagName == l.inRawTextTag {
				l.inRawTextTag = ""
			}
			return tok
		} else if l.peekChar() == '!' {
			// Could be XML comment <!-- --> or CDATA <![CDATA[]]>
			if l.peekCharN(2) == '-' && l.peekCharN(3) == '-' {
				// XML comment - skip it and get next token
				l.skipXMLComment()
				return l.nextTagContentToken()
			} else if l.peekCharN(2) == '[' && l.peekCharN(3) == 'C' {
				// CDATA section - return as TAG_TEXT
				content, ok := l.readCDATA()
				if ok {
					tok.Type = TAG_TEXT
					tok.Literal = content
					tok.Line = line
					tok.Column = column
					return tok
				}
			}
			// Not a comment or CDATA, treat as literal text
			tok.Type = TAG_TEXT
			tok.Literal = string(l.ch)
			tok.Line = line
			tok.Column = column
			l.readChar()
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
				// Check if this is a raw text tag (style or script)
				tagName := extractTagName(tagContent)
				if tagName == "style" || tagName == "script" {
					l.inRawTextTag = tagName
				}
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

	case '/':
		// In normal mode, skip Parsley comments //
		// In raw mode (style/script), keep // as literal text (valid JS comments)
		if !inRawMode && l.peekChar() == '/' {
			l.skipComment()
			return l.nextTagContentToken()
		}
		// Not a comment (or in raw mode), treat as regular text
		tok.Type = TAG_TEXT
		if inRawMode {
			tok.Literal = l.readRawTagText()
		} else {
			tok.Literal = l.readTagText()
		}
		tok.Line = line
		tok.Column = column
		return tok

	case '@':
		// In raw text mode, @{ triggers interpolation
		if inRawMode && l.peekChar() == '{' {
			l.readChar() // skip @
			tok = newToken(LBRACE, l.ch, l.line, l.column)
			l.readChar() // skip {
			l.inTagContent = false
			return tok
		}
		// Not @{, treat @ as regular text
		tok.Type = TAG_TEXT
		if inRawMode {
			tok.Literal = l.readRawTagText()
		} else {
			tok.Literal = l.readTagText()
		}
		tok.Line = line
		tok.Column = column
		return tok

	case '{':
		if inRawMode {
			// In raw text mode, { is literal - read as text
			tok.Type = TAG_TEXT
			tok.Literal = l.readRawTagText()
			tok.Line = line
			tok.Column = column
			return tok
		}
		// Normal mode: interpolation - temporarily exit tag content mode
		tok = newToken(LBRACE, l.ch, l.line, l.column)
		l.readChar()
		l.inTagContent = false
		return tok

	default:
		// Regular text content
		tok.Type = TAG_TEXT
		if inRawMode {
			tok.Literal = l.readRawTagText()
		} else {
			tok.Literal = l.readTagText()
		}
		tok.Line = line
		tok.Column = column
	}

	return tok
}

// readTagText reads text content until we hit <, {, or EOF
func (l *Lexer) readTagText() string {
	var result []byte

	for l.ch != 0 && l.ch != '<' && l.ch != '{' {
		// Skip Parsley comments (//)
		if l.ch == '/' && l.peekChar() == '/' {
			l.skipComment()
			continue
		}
		result = append(result, l.ch)
		l.readChar()
	}

	return string(result)
}

// readRawTagText reads text content in raw text mode (style/script)
// In raw text mode, {} is literal and @{} is used for interpolation
// // comments are preserved (valid in JavaScript, harmless in CSS)
// Stops at <, @{, or EOF
func (l *Lexer) readRawTagText() string {
	var result []byte

	for l.ch != 0 && l.ch != '<' {
		// Check for @{ which starts interpolation in raw text mode
		// This works even inside // comments, allowing datestamps etc.
		if l.ch == '@' && l.peekChar() == '{' {
			break
		}
		result = append(result, l.ch)
		l.readChar()
	}

	return string(result)
}

// extractTagName extracts the tag name from tag content (e.g., "div class=\"foo\"" -> "div")
func extractTagName(tagContent string) string {
	var name []byte
	for i := 0; i < len(tagContent); i++ {
		ch := tagContent[i]
		if ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r' {
			break
		}
		name = append(name, ch)
	}
	return string(name)
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

// isDatetimeLiteral checks if the @ symbol is followed by a datetime pattern (YYYY-MM-DD)
// rather than a duration pattern (like 2h30m or 7d)
func (l *Lexer) isDatetimeLiteral() bool {
	// Look ahead to check the pattern
	// Datetime starts with 4 digits (year)
	// Duration starts with digits followed by a unit letter
	pos := l.readPosition

	// Skip @ (already at current position)
	if pos >= len(l.input) {
		return false
	}

	// Count consecutive digits
	digitCount := 0
	for pos < len(l.input) && isDigit(l.input[pos]) {
		digitCount++
		pos++
	}

	// If we have 4 digits followed by '-', it's a datetime (YYYY-MM-DD)
	if digitCount == 4 && pos < len(l.input) && l.input[pos] == '-' {
		return true
	}

	// Otherwise, it's a duration
	return false
}

// readDurationLiteral reads a duration literal after @
// Supports formats: @2h30m, @7d, @1y6mo, @30s, @-1d, @-2w (negative durations)
// Units: y (years), mo (months), w (weeks), d (days), h (hours), m (minutes), s (seconds)
func (l *Lexer) readDurationLiteral() string {
	var duration []byte
	l.readChar() // skip @

	// Check for negative duration
	if l.ch == '-' {
		duration = append(duration, l.ch)
		l.readChar()
	}

	// Read pairs of number + unit
	for {
		// Read number
		if !isDigit(l.ch) {
			break
		}

		for isDigit(l.ch) {
			duration = append(duration, l.ch)
			l.readChar()
		}

		// Read unit (could be single letter or "mo" for months)
		if !isLetter(l.ch) {
			break
		}

		// Check for "mo" (months)
		if l.ch == 'm' && l.peekChar() == 'o' {
			duration = append(duration, l.ch)
			l.readChar()
			duration = append(duration, l.ch)
			l.readChar()
		} else {
			// Single letter unit
			duration = append(duration, l.ch)
			l.readChar()
		}
	}

	// Back up one char since we consumed one too many
	l.position = l.readPosition - 1
	if l.position < len(l.input) {
		l.ch = l.input[l.position]
	} else {
		l.ch = 0
	}

	return string(duration)
}

// isWhitespace checks if the given byte is whitespace
func isWhitespace(ch byte) bool {
	return ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r'
}

// detectAtLiteralType determines what type of @ literal this is
// Returns the appropriate TokenType for the literal
func (l *Lexer) detectAtLiteralType() TokenType {
	pos := l.readPosition

	if pos >= len(l.input) {
		return ILLEGAL
	}

	// Check for URL: @scheme://
	// Look for characters followed by ://
	colonPos := pos
	for colonPos < len(l.input) && colonPos < pos+20 {
		if l.input[colonPos] == ':' {
			if colonPos+2 < len(l.input) && l.input[colonPos+1] == '/' && l.input[colonPos+2] == '/' {
				return URL_LITERAL
			}
			break
		}
		if !isLetter(l.input[colonPos]) && l.input[colonPos] != '+' && l.input[colonPos] != '-' {
			break
		}
		colonPos++
	}

	// Check for path: @/ or @./ or @~/ or @../
	firstChar := l.input[pos]
	if firstChar == '/' {
		return PATH_LITERAL
	}
	if firstChar == '.' && pos+1 < len(l.input) && (l.input[pos+1] == '/' || l.input[pos+1] == '.') {
		return PATH_LITERAL
	}
	if firstChar == '~' && pos+1 < len(l.input) && l.input[pos+1] == '/' {
		return PATH_LITERAL
	}

	// Check for negative duration: @-1d, @-2w, etc.
	// A minus followed by a digit indicates negative duration
	if firstChar == '-' && pos+1 < len(l.input) && isDigit(l.input[pos+1]) {
		return DURATION_LITERAL
	}

	// Check for datetime: 4 digits followed by '-'
	digitCount := 0
	checkPos := pos
	for checkPos < len(l.input) && isDigit(l.input[checkPos]) {
		digitCount++
		checkPos++
	}

	if digitCount == 4 && checkPos < len(l.input) && l.input[checkPos] == '-' {
		return DATETIME_LITERAL
	}

	// Default to duration
	return DURATION_LITERAL
}

// readPathLiteral reads a path literal after @
// Supports: @/absolute/path, @./relative/path, @~/home/path
func (l *Lexer) readPathLiteral() string {
	l.readChar() // skip @

	var path []byte
	// Read until whitespace or delimiter
	for l.ch != 0 && !isWhitespace(l.ch) {
		// Stop at delimiters that can't be in a path literal
		if l.ch == ')' || l.ch == ']' || l.ch == '}' || l.ch == ',' || l.ch == ';' {
			break
		}
		// Stop at dot if it's NOT followed by / and the previous char WAS / (i.e., /file.ext stops at . before property name)
		// This allows ./path and ../path and file.txt but stops at .basename
		if l.ch == '.' && len(path) > 0 {
			// If previous char is '/' and next char is a letter, this is .property access
			if path[len(path)-1] == '/' && isLetter(l.peekChar()) {
				break
			}
			// If next char is neither '/' nor '.', and we have text before, it might be property access
			nextCh := l.peekChar()
			if nextCh != '/' && nextCh != '.' && !isPathChar(nextCh) && isLetter(nextCh) {
				break
			}
		}
		path = append(path, l.ch)
		l.readChar()
	}

	return string(path)
}

// isPathChar checks if a character is valid in a path (but not at the start of a property)
func isPathChar(ch byte) bool {
	return ch == '/' || ch == '-' || ch == '_' || ch == '~' || isLetter(ch) || isDigit(ch)
}

// readUrlLiteral reads a URL literal after @
// Supports: @scheme://host/path?query#fragment
func (l *Lexer) readUrlLiteral() string {
	l.readChar() // skip @

	var url []byte
	hasScheme := false // Track if we've seen ://

	// Read until whitespace or delimiter
	for l.ch != 0 && !isWhitespace(l.ch) {
		// Track if we've seen the :// pattern
		if !hasScheme && l.ch == ':' && len(url) > 0 {
			if l.peekChar() == '/' {
				hasScheme = true
			}
		}

		// Stop at delimiters that can't be in a URL literal
		if l.ch == ')' || l.ch == ']' || l.ch == '}' || l.ch == ',' || l.ch == ';' {
			break
		}

		// For dots: if we've seen ://, ALL dots are part of the URL until we hit a delimiter or whitespace
		// This handles .com, .org, file.html, etc.
		// We only stop at . for property access if there's no scheme (edge case)
		if l.ch == '.' && isLetter(l.peekChar()) && !hasScheme {
			// No :// seen yet, so this might be property access
			break
		}

		url = append(url, l.ch)
		l.readChar()
	}

	return string(url)
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
