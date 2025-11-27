package parser

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/sambeau/parsley/pkg/ast"
	"github.com/sambeau/parsley/pkg/lexer"
)

// Precedence levels for operators
const (
	_ int = iota
	LOWEST
	COMMA_PREC  // ,
	LOGIC_OR    // or, |
	LOGIC_AND   // and, &
	EQUALS      // ==
	LESSGREATER // > or <
	SUM         // +
	CONCAT      // ++
	PRODUCT     // *
	PREFIX      // -X or !X
	INDEX       // array[index]
	CALL        // myFunction(X)
)

// precedences maps tokens to their precedence
var precedences = map[lexer.TokenType]int{
	lexer.COMMA:     COMMA_PREC,
	lexer.OR:        LOGIC_OR,
	lexer.NULLISH:   LOGIC_OR,
	lexer.AND:       LOGIC_AND,
	lexer.EQ:        EQUALS,
	lexer.NOT_EQ:    EQUALS,
	lexer.MATCH:     EQUALS,
	lexer.NOT_MATCH: EQUALS,
	lexer.LT:        LESSGREATER,
	lexer.GT:        LESSGREATER,
	lexer.LTE:       LESSGREATER,
	lexer.GTE:       LESSGREATER,
	lexer.PLUS:      SUM,
	lexer.MINUS:     SUM,
	lexer.PLUSPLUS:  CONCAT,
	lexer.SLASH:     PRODUCT,
	lexer.ASTERISK:  PRODUCT,
	lexer.PERCENT:   PRODUCT,
	lexer.LBRACKET:  INDEX,
	lexer.DOT:       INDEX,
	lexer.LPAREN:    CALL,
}

// Parser represents the parser
type Parser struct {
	l *lexer.Lexer

	errors []string

	prevToken lexer.Token
	curToken  lexer.Token
	peekToken lexer.Token

	prefixParseFns map[lexer.TokenType]prefixParseFn
	infixParseFns  map[lexer.TokenType]infixParseFn
}

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

// New creates a new parser instance
func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	// Initialize prefix parse functions
	p.prefixParseFns = make(map[lexer.TokenType]prefixParseFn)
	p.registerPrefix(lexer.IDENT, p.parseIdentifier)
	p.registerPrefix(lexer.INT, p.parseIntegerLiteral)
	p.registerPrefix(lexer.FLOAT, p.parseFloatLiteral)
	p.registerPrefix(lexer.STRING, p.parseStringLiteral)
	p.registerPrefix(lexer.TEMPLATE, p.parseTemplateLiteral)
	p.registerPrefix(lexer.REGEX, p.parseRegexLiteral)
	p.registerPrefix(lexer.DATETIME_LITERAL, p.parseDatetimeLiteral)
	p.registerPrefix(lexer.DURATION_LITERAL, p.parseDurationLiteral)
	p.registerPrefix(lexer.PATH_LITERAL, p.parsePathLiteral)
	p.registerPrefix(lexer.URL_LITERAL, p.parseUrlLiteral)
	p.registerPrefix(lexer.TAG, p.parseTagLiteral)
	p.registerPrefix(lexer.TAG_START, p.parseTagPair)
	p.registerPrefix(lexer.BANG, p.parsePrefixExpression)
	p.registerPrefix(lexer.MINUS, p.parsePrefixExpression)
	p.registerPrefix(lexer.TRUE, p.parseBoolean)
	p.registerPrefix(lexer.FALSE, p.parseBoolean)
	p.registerPrefix(lexer.LPAREN, p.parseGroupedExpression)
	p.registerPrefix(lexer.LBRACKET, p.parseSquareBracketArrayLiteral)
	p.registerPrefix(lexer.IF, p.parseIfExpression)
	p.registerPrefix(lexer.FUNCTION, p.parseFunctionLiteral)
	p.registerPrefix(lexer.FOR, p.parseForExpression)
	p.registerPrefix(lexer.LBRACE, p.parseDictionaryLiteral)

	// Initialize infix parse functions
	p.infixParseFns = make(map[lexer.TokenType]infixParseFn)
	p.registerInfix(lexer.PLUS, p.parseInfixExpression)
	p.registerInfix(lexer.MINUS, p.parseInfixExpression)
	p.registerInfix(lexer.SLASH, p.parseInfixExpression)
	p.registerInfix(lexer.ASTERISK, p.parseInfixExpression)
	p.registerInfix(lexer.PERCENT, p.parseInfixExpression)
	p.registerInfix(lexer.EQ, p.parseInfixExpression)
	p.registerInfix(lexer.NOT_EQ, p.parseInfixExpression)
	p.registerInfix(lexer.LT, p.parseInfixExpression)
	p.registerInfix(lexer.GT, p.parseInfixExpression)
	p.registerInfix(lexer.LTE, p.parseInfixExpression)
	p.registerInfix(lexer.GTE, p.parseInfixExpression)
	p.registerInfix(lexer.AND, p.parseInfixExpression)
	p.registerInfix(lexer.OR, p.parseInfixExpression)
	p.registerInfix(lexer.NULLISH, p.parseInfixExpression)
	p.registerInfix(lexer.MATCH, p.parseInfixExpression)
	p.registerInfix(lexer.NOT_MATCH, p.parseInfixExpression)
	p.registerInfix(lexer.PLUSPLUS, p.parseInfixExpression)
	p.registerInfix(lexer.COMMA, p.parseArrayLiteral)
	p.registerInfix(lexer.LPAREN, p.parseCallExpression)
	p.registerInfix(lexer.LBRACKET, p.parseIndexOrSliceExpression)
	p.registerInfix(lexer.DOT, p.parseDotExpression)

	// Read two tokens, so curToken and peekToken are both set
	p.nextToken()
	p.nextToken()

	return p
}

// Errors returns parser errors
func (p *Parser) Errors() []string {
	return p.errors
}

// registerPrefix registers a prefix parse function
func (p *Parser) registerPrefix(tokenType lexer.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

// registerInfix registers an infix parse function
func (p *Parser) registerInfix(tokenType lexer.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

// nextToken advances prevToken, curToken, and peekToken
func (p *Parser) nextToken() {
	p.prevToken = p.curToken
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

// ParseProgram parses the program and returns the AST
func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for !p.curTokenIs(lexer.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}

	return program
}

// parseStatement parses statements
func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case lexer.LET:
		return p.parseLetStatement()
	case lexer.RETURN:
		return p.parseReturnStatement()
	case lexer.DELETE:
		return p.parseDeleteStatement()
	case lexer.LBRACE:
		// Check if this is a dictionary destructuring assignment
		// We need to look ahead to see if this is {a, b} = ... or just a dict literal
		// For now, try parsing as destructuring assignment
		savedCur := p.curToken
		savedPeek := p.peekToken
		savedPrev := p.prevToken
		savedErrors := len(p.errors)

		stmt := p.parseDictDestructuringAssignment()

		// If parsing failed (no = found), restore and parse as expression
		if stmt == nil || len(p.errors) > savedErrors {
			p.curToken = savedCur
			p.peekToken = savedPeek
			p.prevToken = savedPrev
			p.errors = p.errors[:savedErrors]
			return p.parseExpressionStatement()
		}
		return stmt
	case lexer.IDENT:
		// Check if this is an assignment statement
		if p.peekTokenIs(lexer.ASSIGN) {
			return p.parseAssignmentStatement()
		}
		// Check for potential destructuring: IDENT followed by COMMA
		// We need to peek further to determine if this is `x,y = ...` or just `x,y` expression
		// For now, try parsing as assignment if comma follows
		if p.peekTokenIs(lexer.COMMA) {
			// Tentatively parse as destructuring assignment
			savedCur := p.curToken
			savedPeek := p.peekToken
			savedPrev := p.prevToken
			savedErrors := len(p.errors)

			stmt := p.parseAssignmentStatement()

			// If parsing failed (no = found), restore and parse as expression
			if stmt == nil || len(p.errors) > savedErrors {
				p.curToken = savedCur
				p.peekToken = savedPeek
				p.prevToken = savedPrev
				p.errors = p.errors[:savedErrors]
				return p.parseExpressionStatement()
			}
			return stmt
		}
		// Otherwise, treat as expression statement
		return p.parseExpressionStatement()
	default:
		return p.parseExpressionStatement()
	}
}

// parseLetStatement parses let statements
func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.curToken}

	// Check for dictionary destructuring pattern
	if p.peekTokenIs(lexer.LBRACE) {
		p.nextToken() // move to '{'
		stmt.DictPattern = p.parseDictDestructuringPattern()
		if stmt.DictPattern == nil {
			return nil
		}

		if !p.expectPeek(lexer.ASSIGN) {
			return nil
		}

		p.nextToken()
		stmt.Value = p.parseExpression(LOWEST)

		if p.peekTokenIs(lexer.SEMICOLON) {
			p.nextToken()
		}

		return stmt
	}

	if !p.expectPeek(lexer.IDENT) {
		return nil
	}

	// Collect identifiers (for array destructuring)
	names := []*ast.Identifier{
		{Token: p.curToken, Value: p.curToken.Literal},
	}

	// Check for comma-separated identifiers (array destructuring pattern)
	for p.peekTokenIs(lexer.COMMA) {
		p.nextToken() // consume comma
		if !p.expectPeek(lexer.IDENT) {
			return nil
		}
		names = append(names, &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal})
	}

	// Set either single name or multiple names
	if len(names) == 1 {
		stmt.Name = names[0]
	} else {
		stmt.Names = names
	}

	if !p.expectPeek(lexer.ASSIGN) {
		return nil
	}

	p.nextToken()

	stmt.Value = p.parseExpression(LOWEST)

	if p.peekTokenIs(lexer.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

// parseAssignmentStatement parses assignment statements like 'x = 5;' or 'x,y,z = 1,2,3;'
func (p *Parser) parseAssignmentStatement() *ast.AssignmentStatement {
	stmt := &ast.AssignmentStatement{Token: p.curToken}

	// Collect identifiers (for destructuring)
	names := []*ast.Identifier{
		{Token: p.curToken, Value: p.curToken.Literal},
	}

	// Check for comma-separated identifiers (destructuring pattern)
	for p.peekTokenIs(lexer.COMMA) {
		p.nextToken() // consume comma
		if !p.expectPeek(lexer.IDENT) {
			return nil
		}
		names = append(names, &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal})
	}

	// Set either single name or multiple names
	if len(names) == 1 {
		stmt.Name = names[0]
	} else {
		stmt.Names = names
	}

	if !p.expectPeek(lexer.ASSIGN) {
		return nil
	}

	p.nextToken()

	stmt.Value = p.parseExpression(LOWEST)

	if p.peekTokenIs(lexer.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

// parseDictDestructuringAssignment parses dictionary destructuring assignments like '{a, b} = dict;'
func (p *Parser) parseDictDestructuringAssignment() *ast.AssignmentStatement {
	stmt := &ast.AssignmentStatement{Token: p.curToken} // the '{' token

	stmt.DictPattern = p.parseDictDestructuringPattern()
	if stmt.DictPattern == nil {
		return nil
	}

	if !p.expectPeek(lexer.ASSIGN) {
		return nil
	}

	p.nextToken()

	stmt.Value = p.parseExpression(LOWEST)

	if p.peekTokenIs(lexer.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

// parseReturnStatement parses return statements
func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken}

	p.nextToken()

	stmt.ReturnValue = p.parseExpression(LOWEST)

	if p.peekTokenIs(lexer.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

// parseExpressionStatement parses expression statements
func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}

	stmt.Expression = p.parseExpression(LOWEST)

	if p.peekTokenIs(lexer.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

// parseExpression parses expressions using Pratt parsing
func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}

	leftExp := prefix()

	for !p.peekTokenIs(lexer.SEMICOLON) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}

		p.nextToken()

		leftExp = infix(leftExp)
	}

	return leftExp
}

// Parse functions for different expression types
func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.curToken}

	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	lit.Value = value
	return lit
}

func (p *Parser) parseFloatLiteral() ast.Expression {
	lit := &ast.FloatLiteral{Token: p.curToken}

	value, err := strconv.ParseFloat(p.curToken.Literal, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as float", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	lit.Value = value
	return lit
}

func (p *Parser) parseStringLiteral() ast.Expression {
	return &ast.StringLiteral{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseTemplateLiteral() ast.Expression {
	return &ast.TemplateLiteral{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseRegexLiteral() ast.Expression {
	// Token.Literal is in the form "/pattern/flags"
	literal := p.curToken.Literal
	if len(literal) < 2 || literal[0] != '/' {
		p.errors = append(p.errors, fmt.Sprintf("invalid regex literal: %s", literal))
		return nil
	}

	// Find the closing / by looking from the end backwards
	// This handles /pattern/ and /pattern/flags
	lastSlash := strings.LastIndex(literal[1:], "/")
	if lastSlash == -1 {
		p.errors = append(p.errors, fmt.Sprintf("unterminated regex literal: %s", literal))
		return nil
	}
	lastSlash++ // adjust for the slice offset

	pattern := literal[1:lastSlash]
	flags := ""
	if lastSlash+1 < len(literal) {
		flags = literal[lastSlash+1:]
	}

	return &ast.RegexLiteral{
		Token:   p.curToken,
		Pattern: pattern,
		Flags:   flags,
	}
}

func (p *Parser) parseDatetimeLiteral() ast.Expression {
	// Token.Literal contains the ISO-8601 datetime string (without the @ prefix)
	return &ast.DatetimeLiteral{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}
}

func (p *Parser) parseDurationLiteral() ast.Expression {
	// Token.Literal contains the duration string (without the @ prefix)
	return &ast.DurationLiteral{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}
}

func (p *Parser) parsePathLiteral() ast.Expression {
	// Token.Literal contains the path string (without the @ prefix)
	return &ast.PathLiteral{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}
}

func (p *Parser) parseUrlLiteral() ast.Expression {
	// Token.Literal contains the URL string (without the @ prefix)
	return &ast.UrlLiteral{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}
}

func (p *Parser) parseTagLiteral() ast.Expression {
	return &ast.TagLiteral{Token: p.curToken, Raw: p.curToken.Literal}
}

func (p *Parser) parseTagPair() ast.Expression {
	tagExpr := &ast.TagPairExpression{
		Token:    p.curToken,
		Contents: []ast.Node{},
	}

	// Parse tag name and props from the TAG_START token literal
	// Format: "tagname attr1="value" attr2={expr}" or empty string for <>
	raw := p.curToken.Literal
	tagExpr.Name, tagExpr.Props = parseTagNameAndProps(raw)

	// Parse tag contents
	p.nextToken()
	tagExpr.Contents = p.parseTagContents(tagExpr.Name)

	// Current token should be TAG_END
	if !p.curTokenIs(lexer.TAG_END) {
		p.errors = append(p.errors, fmt.Sprintf("expected closing tag, got %s at line %d, column %d",
			p.curToken.Type, p.curToken.Line, p.curToken.Column))
		return nil
	}

	// Validate closing tag matches opening tag
	closingName := p.curToken.Literal
	if closingName != tagExpr.Name {
		p.errors = append(p.errors, fmt.Sprintf("mismatched tags: opening <%s> but closing </%s> at line %d, column %d",
			tagExpr.Name, closingName, p.curToken.Line, p.curToken.Column))
		return nil
	}

	return tagExpr
}

// parseTagContents parses the contents between opening and closing tags
func (p *Parser) parseTagContents(tagName string) []ast.Node {
	var contents []ast.Node

	for !p.curTokenIs(lexer.TAG_END) && !p.curTokenIs(lexer.EOF) {
		switch p.curToken.Type {
		case lexer.TAG_TEXT:
			// Raw text content
			textNode := &ast.TextNode{
				Token: p.curToken,
				Value: p.curToken.Literal,
			}
			contents = append(contents, textNode)
			p.nextToken()

		case lexer.TAG_START:
			// Nested tag pair
			nestedTag := p.parseTagPair()
			if nestedTag != nil {
				contents = append(contents, nestedTag)
			}
			p.nextToken()

		case lexer.TAG:
			// Singleton tag
			singletonTag := p.parseTagLiteral()
			if singletonTag != nil {
				contents = append(contents, singletonTag)
			}
			p.nextToken()

		case lexer.LBRACE:
			// Interpolation expression
			p.nextToken() // skip {
			expr := p.parseExpression(LOWEST)
			if expr != nil {
				contents = append(contents, expr)
			}
			// Re-enter tag content mode BEFORE checking for }
			p.l.EnterTagContentMode()
			if !p.expectPeek(lexer.RBRACE) {
				return contents
			}
			p.nextToken() // move past }

		default:
			// Unexpected token
			p.errors = append(p.errors, fmt.Sprintf("unexpected token in tag contents: %s at line %d, column %d",
				p.curToken.Type, p.curToken.Line, p.curToken.Column))
			p.nextToken()
		}
	}

	return contents
}

// parseTagNameAndProps splits raw tag content into name and props
// Examples: "div class=\"foo\"" -> ("div", "class=\"foo\"")
//
//	"" -> ("", "")
func parseTagNameAndProps(raw string) (string, string) {
	if raw == "" {
		return "", "" // empty grouping tag
	}

	// Find first space to separate name from props
	for i, ch := range raw {
		if ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r' {
			return raw[:i], strings.TrimSpace(raw[i:])
		}
	}

	// No spaces, all tag name
	return raw, ""
}

func (p *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{Token: p.curToken, Value: p.curTokenIs(lexer.TRUE)}
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}

	p.nextToken()

	expression.Right = p.parseExpression(PREFIX)

	return expression
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    p.curToken,
		Left:     left,
		Operator: p.curToken.Literal,
	}

	precedence := p.curPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)

	return expression
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()

	exp := p.parseExpression(LOWEST)

	if !p.expectPeek(lexer.RPAREN) {
		return nil
	}

	return exp
}

func (p *Parser) parseSquareBracketArrayLiteral() ast.Expression {
	array := &ast.ArrayLiteral{Token: p.curToken}
	array.Elements = []ast.Expression{}

	// Check for empty array []
	if p.peekTokenIs(lexer.RBRACKET) {
		p.nextToken()
		return array
	}

	// Parse first element - use COMMA_PREC to prevent comma from being treated as infix
	p.nextToken()
	array.Elements = append(array.Elements, p.parseExpression(COMMA_PREC))

	// Parse remaining elements separated by commas
	for p.peekTokenIs(lexer.COMMA) {
		p.nextToken() // consume current element
		p.nextToken() // consume comma
		array.Elements = append(array.Elements, p.parseExpression(COMMA_PREC))
	}

	if !p.expectPeek(lexer.RBRACKET) {
		return nil
	}

	return array
}

func (p *Parser) parseIfExpression() ast.Expression {
	expression := &ast.IfExpression{Token: p.curToken}

	if !p.expectPeek(lexer.LPAREN) {
		return nil
	}

	p.nextToken()

	// Check for assignment in condition
	if p.curTokenIs(lexer.IDENT) && p.peekTokenIs(lexer.ASSIGN) {
		varName := p.curToken.Literal
		msg := fmt.Sprintf("assignment is not allowed inside if condition. Use a separate statement:\n  let %s = ...\n  if (%s) { ... }",
			varName, varName)
		p.errors = append(p.errors, msg)
		return nil
	}

	expression.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(lexer.RPAREN) {
		return nil
	}

	if p.peekTokenIs(lexer.LBRACE) {
		// Block form: if (...) { ... }
		p.nextToken()
		expression.Consequence = p.parseBlockStatement()
	} else {
		// Single statement/expression form: if (...) expr or if (...) return expr
		p.nextToken()

		// Check if it's a return statement
		if p.curTokenIs(lexer.RETURN) {
			stmt := p.parseReturnStatement()
			expression.Consequence = &ast.BlockStatement{
				Token:      p.curToken,
				Statements: []ast.Statement{stmt},
			}
		} else {
			// Regular expression
			stmt := &ast.ExpressionStatement{Token: p.curToken}
			stmt.Expression = p.parseExpression(LOWEST)
			expression.Consequence = &ast.BlockStatement{
				Token:      p.curToken,
				Statements: []ast.Statement{stmt},
			}
		}
	}

	// Optional else clause
	if p.peekTokenIs(lexer.ELSE) {
		p.nextToken()

		// Check if alternative is a block statement or single statement/expression
		if p.peekTokenIs(lexer.LBRACE) {
			p.nextToken()
			expression.Alternative = p.parseBlockStatement()
		} else {
			// Single statement/expression form
			p.nextToken()

			// Check if it's a return statement
			if p.curTokenIs(lexer.RETURN) {
				stmt := p.parseReturnStatement()
				expression.Alternative = &ast.BlockStatement{
					Token:      p.curToken,
					Statements: []ast.Statement{stmt},
				}
			} else {
				// Parse single expression as alternative
				stmt := &ast.ExpressionStatement{Token: p.curToken}
				stmt.Expression = p.parseExpression(LOWEST)
				expression.Alternative = &ast.BlockStatement{
					Token:      p.curToken,
					Statements: []ast.Statement{stmt},
				}
			}
		}
	}

	return expression
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: p.curToken}
	block.Statements = []ast.Statement{}

	p.nextToken()

	for !p.curTokenIs(lexer.RBRACE) && !p.curTokenIs(lexer.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}

	return block
}

func (p *Parser) parseFunctionLiteral() ast.Expression {
	lit := &ast.FunctionLiteral{Token: p.curToken}

	if !p.expectPeek(lexer.LPAREN) {
		return nil
	}

	// Use new parameter parsing that supports destructuring
	lit.Params = p.parseFunctionParametersNew()

	if !p.expectPeek(lexer.LBRACE) {
		return nil
	}

	lit.Body = p.parseBlockStatement()

	return lit
}

func (p *Parser) parseFunctionParameters() []*ast.Identifier {
	identifiers := []*ast.Identifier{}

	if p.peekTokenIs(lexer.RPAREN) {
		p.nextToken()
		return identifiers
	}

	p.nextToken()

	ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	identifiers = append(identifiers, ident)

	for p.peekTokenIs(lexer.COMMA) {
		p.nextToken()
		p.nextToken()
		ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
		identifiers = append(identifiers, ident)
	}

	if !p.expectPeek(lexer.RPAREN) {
		return nil
	}

	return identifiers
}

// parseFunctionParametersNew parses function parameters with destructuring support
func (p *Parser) parseFunctionParametersNew() []*ast.FunctionParameter {
	params := []*ast.FunctionParameter{}

	if p.peekTokenIs(lexer.RPAREN) {
		p.nextToken()
		return params
	}

	p.nextToken()

	// Parse first parameter
	param := p.parseFunctionParameter()
	if param != nil {
		params = append(params, param)
	}

	for p.peekTokenIs(lexer.COMMA) {
		p.nextToken() // consume comma
		p.nextToken() // move to next parameter
		param := p.parseFunctionParameter()
		if param != nil {
			params = append(params, param)
		}
	}

	if !p.expectPeek(lexer.RPAREN) {
		return nil
	}

	return params
}

// parseFunctionParameter parses a single function parameter (can be identifier, array, or dict pattern)
func (p *Parser) parseFunctionParameter() *ast.FunctionParameter {
	param := &ast.FunctionParameter{}

	switch p.curToken.Type {
	case lexer.LBRACE:
		// Dictionary destructuring pattern
		param.DictPattern = p.parseDictDestructuringPattern()
		return param

	case lexer.LBRACKET:
		// Array destructuring pattern
		if !p.expectPeek(lexer.IDENT) {
			return nil
		}
		idents := []*ast.Identifier{
			{Token: p.curToken, Value: p.curToken.Literal},
		}
		for p.peekTokenIs(lexer.COMMA) {
			p.nextToken() // consume comma
			if !p.expectPeek(lexer.IDENT) {
				return nil
			}
			idents = append(idents, &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal})
		}
		if !p.expectPeek(lexer.RBRACKET) {
			return nil
		}
		param.ArrayPattern = idents
		return param

	case lexer.IDENT:
		// Simple identifier
		param.Ident = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
		return param

	default:
		return nil
	}
}

// parseForExpression parses for expressions
// Two forms: for(array) func  OR  for(var in array) body
func (p *Parser) parseForExpression() ast.Expression {
	expression := &ast.ForExpression{Token: p.curToken}

	if !p.expectPeek(lexer.LPAREN) {
		return nil
	}

	p.nextToken()

	// Check if this is the "for(var in array)" form
	// We need to peek ahead to see if there's an IN token or COMMA (for key,value syntax)
	// But only if current token is an identifier
	if p.curToken.Type == lexer.IDENT && (p.peekTokenIs(lexer.IN) || p.peekTokenIs(lexer.COMMA)) {
		// Parse: for(var in array) body OR for(key, value in dict) body
		expression.Variable = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

		// Check for comma (key, value in dict form)
		if p.peekTokenIs(lexer.COMMA) {
			p.nextToken() // move to COMMA
			if !p.expectPeek(lexer.IDENT) {
				return nil
			}
			expression.ValueVariable = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
		}

		if !p.expectPeek(lexer.IN) {
			return nil
		}
		p.nextToken() // move past IN to array/dict expression

		expression.Array = p.parseExpression(LOWEST)

		if !p.expectPeek(lexer.RPAREN) {
			return nil
		}

		// Parse body - must be a block expression
		if !p.expectPeek(lexer.LBRACE) {
			return nil
		}

		// Create a function literal for the body
		bodyFn := &ast.FunctionLiteral{
			Token: p.curToken,
		}

		// Set parameters based on whether we have one or two variables
		if expression.ValueVariable != nil {
			bodyFn.Parameters = []*ast.Identifier{expression.Variable, expression.ValueVariable}
		} else {
			bodyFn.Parameters = []*ast.Identifier{expression.Variable}
		}

		bodyFn.Body = p.parseBlockStatement()
		expression.Body = bodyFn
	} else {
		// Parse: for(array) func
		expression.Array = p.parseExpression(LOWEST)

		if !p.expectPeek(lexer.RPAREN) {
			return nil
		}

		p.nextToken() // move past RPAREN to function

		expression.Function = p.parseExpression(LOWEST)
	}

	return expression
}

func (p *Parser) parseCallExpression(fn ast.Expression) ast.Expression {
	exp := &ast.CallExpression{Token: p.curToken, Function: fn}
	exp.Arguments = p.parseExpressionList(lexer.RPAREN)
	return exp
}

func (p *Parser) parseExpressionList(end lexer.TokenType) []ast.Expression {
	args := []ast.Expression{}

	if p.peekTokenIs(end) {
		p.nextToken()
		return args
	}

	p.nextToken()
	args = append(args, p.parseExpression(COMMA_PREC+1))

	for p.peekTokenIs(lexer.COMMA) {
		p.nextToken()
		p.nextToken()
		args = append(args, p.parseExpression(COMMA_PREC+1))
	}

	if !p.expectPeek(end) {
		return nil
	}

	return args
}

func (p *Parser) parseArrayLiteral(left ast.Expression) ast.Expression {
	// The left expression may already be an array if we're chaining commas
	// curToken is COMMA at this point
	var array *ast.ArrayLiteral

	if leftArray, ok := left.(*ast.ArrayLiteral); ok {
		// Left is already an array, extend it
		array = leftArray
	} else {
		// Create new array with left as first element
		array = &ast.ArrayLiteral{Token: p.curToken}
		array.Elements = []ast.Expression{left}
	}

	// Parse the right side of the current comma
	p.nextToken()
	array.Elements = append(array.Elements, p.parseExpression(COMMA_PREC))

	return array
}

func (p *Parser) parseIndexOrSliceExpression(left ast.Expression) ast.Expression {
	exp := &ast.IndexExpression{Token: p.curToken, Left: left}

	p.nextToken()

	// Check for slice (colon before any expression, or expression followed by colon)
	if p.curTokenIs(lexer.COLON) {
		// Slice with no start: [:end]
		return p.parseSliceExpression(left, nil)
	}

	// Parse the first expression (could be index or slice start)
	firstExp := p.parseExpression(LOWEST)

	// Check if this is a slice
	if p.peekTokenIs(lexer.COLON) {
		p.nextToken() // consume colon
		return p.parseSliceExpression(left, firstExp)
	}

	// It's an index expression
	exp.Index = firstExp

	if !p.expectPeek(lexer.RBRACKET) {
		return nil
	}

	return exp
}

func (p *Parser) parseSliceExpression(left ast.Expression, start ast.Expression) ast.Expression {
	exp := &ast.SliceExpression{
		Token: p.curToken,
		Left:  left,
		Start: start,
	}

	// We're at the colon, move to next token
	p.nextToken()

	// Check if there's an end expression
	if !p.curTokenIs(lexer.RBRACKET) {
		// Parse the end expression
		exp.End = p.parseExpression(LOWEST)
		// After parsing, expect the closing bracket
		if !p.expectPeek(lexer.RBRACKET) {
			return nil
		}
	}
	// If we're already at RBRACKET, this handles open-ended slices like arr[1:] or arr[:]

	return exp
}

// Helper functions
func (p *Parser) curTokenIs(t lexer.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t lexer.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t lexer.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}

func (p *Parser) peekError(t lexer.TokenType) {
	tokenName := tokenTypeToReadableName(t)
	gotName := tokenTypeToReadableName(p.peekToken.Type)
	gotLiteral := p.peekToken.Literal
	if gotLiteral == "" {
		gotLiteral = gotName
	}

	// Report error at the position after the last successfully parsed token (curToken)
	line := p.curToken.Line
	column := p.curToken.Column + len(p.curToken.Literal)

	msg := fmt.Sprintf("line %d, column %d: expected %s, got '%s'",
		line, column, tokenName, gotLiteral)
	p.errors = append(p.errors, msg)
}

func (p *Parser) noPrefixParseFnError(t lexer.TokenType) {
	literal := p.curToken.Literal
	if literal == "" {
		literal = tokenTypeToReadableName(t)
	}

	// If curToken is on a new line compared to prevToken,
	// report the error at the previous token (where the expression should have been)
	line := p.curToken.Line
	column := p.curToken.Column + len(p.curToken.Literal)

	if p.prevToken.Type != lexer.ILLEGAL && p.curToken.Line > p.prevToken.Line {
		// Current token is on a new line, point to after the previous token
		line = p.prevToken.Line
		column = p.prevToken.Column + len(p.prevToken.Literal)
	} else if p.prevToken.Type != lexer.ILLEGAL {
		// Same line, point to after the previous token
		column = p.prevToken.Column + len(p.prevToken.Literal)
	}

	msg := fmt.Sprintf("line %d, column %d: unexpected '%s'",
		line, column, literal)
	p.errors = append(p.errors, msg)
}

// tokenTypeToReadableName converts token types to human-readable names
func tokenTypeToReadableName(t lexer.TokenType) string {
	switch t {
	case lexer.IDENT:
		return "identifier"
	case lexer.INT:
		return "integer"
	case lexer.FLOAT:
		return "float"
	case lexer.STRING:
		return "string"
	case lexer.TEMPLATE:
		return "template literal"
	case lexer.ASSIGN:
		return "'='"
	case lexer.PLUS:
		return "'+'"
	case lexer.MINUS:
		return "'-'"
	case lexer.BANG:
		return "'!'"
	case lexer.ASTERISK:
		return "'*'"
	case lexer.SLASH:
		return "'/'"
	case lexer.PERCENT:
		return "'%'"
	case lexer.LT:
		return "'<'"
	case lexer.GT:
		return "'>'"
	case lexer.EQ:
		return "'=='"
	case lexer.NOT_EQ:
		return "'!='"
	case lexer.COMMA:
		return "','"
	case lexer.SEMICOLON:
		return "';'"
	case lexer.COLON:
		return "':'"
	case lexer.LPAREN:
		return "'('"
	case lexer.RPAREN:
		return "')'"
	case lexer.LBRACE:
		return "'{'"
	case lexer.RBRACE:
		return "'}'"
	case lexer.LBRACKET:
		return "'['"
	case lexer.RBRACKET:
		return "']'"
	case lexer.FUNCTION:
		return "'fn'"
	case lexer.LET:
		return "'let'"
	case lexer.TRUE:
		return "'true'"
	case lexer.FALSE:
		return "'false'"
	case lexer.IF:
		return "'if'"
	case lexer.ELSE:
		return "'else'"
	case lexer.RETURN:
		return "'return'"
	case lexer.FOR:
		return "'for'"
	case lexer.IN:
		return "'in'"
	case lexer.EOF:
		return "end of file"
	case lexer.ILLEGAL:
		return "illegal character"
	default:
		return string(t.String())
	}
}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}
	return LOWEST
}

// parseDictionaryLiteral parses dictionary literals like { key: value, ... }
func (p *Parser) parseDictionaryLiteral() ast.Expression {
	dict := &ast.DictionaryLiteral{Token: p.curToken}
	dict.Pairs = make(map[string]ast.Expression)

	// Empty dictionary
	if p.peekTokenIs(lexer.RBRACE) {
		p.nextToken()
		return dict
	}

	// Parse key-value pairs
	for !p.curTokenIs(lexer.RBRACE) {
		p.nextToken()

		// Key must be an identifier
		if !p.curTokenIs(lexer.IDENT) {
			p.errors = append(p.errors, fmt.Sprintf("expected identifier as dictionary key, got %s at line %d, column %d",
				p.curToken.Type, p.curToken.Line, p.curToken.Column))
			return nil
		}
		key := p.curToken.Literal

		// Expect colon
		if !p.expectPeek(lexer.COLON) {
			return nil
		}

		// Parse value expression with COMMA_PREC+1 to avoid consuming commas
		p.nextToken()
		value := p.parseExpression(COMMA_PREC + 1)
		if value == nil {
			return nil
		}

		dict.Pairs[key] = value

		// Check for comma, semicolon, or closing brace
		if p.peekTokenIs(lexer.RBRACE) {
			p.nextToken()
			break
		}
		if p.peekTokenIs(lexer.COMMA) || p.peekTokenIs(lexer.SEMICOLON) {
			p.nextToken()
			// Skip any extra commas/semicolons
			for p.peekTokenIs(lexer.COMMA) || p.peekTokenIs(lexer.SEMICOLON) {
				p.nextToken()
			}
		}
	}

	return dict
}

// parseDotExpression parses dot notation like dict.key
func (p *Parser) parseDotExpression(left ast.Expression) ast.Expression {
	dotExpr := &ast.DotExpression{
		Token: p.curToken,
		Left:  left,
	}

	if !p.expectPeek(lexer.IDENT) {
		return nil
	}

	dotExpr.Key = p.curToken.Literal
	return dotExpr
}

// parseDeleteStatement parses delete statements
func (p *Parser) parseDeleteStatement() ast.Statement {
	stmt := &ast.DeleteStatement{Token: p.curToken}

	p.nextToken()

	// Parse the target expression (should be a property access)
	stmt.Target = p.parseExpression(LOWEST)

	// Optional semicolon
	if p.peekTokenIs(lexer.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

// parseDictDestructuringPattern parses dictionary destructuring patterns like {a, b as c, ...rest}
func (p *Parser) parseDictDestructuringPattern() *ast.DictDestructuringPattern {
	pattern := &ast.DictDestructuringPattern{Token: p.curToken} // the '{' token

	// Check for empty pattern
	if p.peekTokenIs(lexer.RBRACE) {
		msg := fmt.Sprintf("empty dictionary destructuring pattern at line %d, column %d",
			p.peekToken.Line, p.peekToken.Column)
		p.errors = append(p.errors, msg)
		return nil
	}

	p.nextToken() // move to first identifier or ...

	// Parse keys
	for {
		// Check for rest operator
		if p.curTokenIs(lexer.DOTDOTDOT) {
			if !p.expectPeek(lexer.IDENT) {
				msg := fmt.Sprintf("expected identifier after '...' at line %d, column %d",
					p.peekToken.Line, p.peekToken.Column)
				p.errors = append(p.errors, msg)
				return nil
			}
			pattern.Rest = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

			// Rest must be at the end
			if !p.peekTokenIs(lexer.RBRACE) {
				msg := fmt.Sprintf("rest element must be last in destructuring pattern at line %d, column %d",
					p.peekToken.Line, p.peekToken.Column)
				p.errors = append(p.errors, msg)
				return nil
			}
			break
		}

		// Expect identifier for key
		if !p.curTokenIs(lexer.IDENT) {
			return nil
		}

		// Parse regular key
		key := &ast.DictDestructuringKey{
			Token: p.curToken,
			Key:   &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal},
		}

		// Check for alias (as syntax)
		if p.peekTokenIs(lexer.AS) {
			p.nextToken() // consume 'as'
			if !p.expectPeek(lexer.IDENT) {
				msg := fmt.Sprintf("expected identifier after 'as' at line %d, column %d",
					p.peekToken.Line, p.peekToken.Column)
				p.errors = append(p.errors, msg)
				return nil
			}
			key.Alias = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
		}

		// Check for nested pattern (colon syntax)
		if p.peekTokenIs(lexer.COLON) {
			p.nextToken() // consume ':'
			p.nextToken() // move to pattern start

			// Parse nested pattern
			if p.curTokenIs(lexer.LBRACE) {
				key.Nested = p.parseDictDestructuringPattern()
			} else if p.curTokenIs(lexer.LBRACKET) {
				// For nested array destructuring, we'd need to handle this
				msg := fmt.Sprintf("nested array destructuring not yet supported at line %d, column %d",
					p.curToken.Line, p.curToken.Column)
				p.errors = append(p.errors, msg)
				return nil
			} else {
				msg := fmt.Sprintf("expected destructuring pattern after ':' at line %d, column %d",
					p.curToken.Line, p.curToken.Column)
				p.errors = append(p.errors, msg)
				return nil
			}
		}

		pattern.Keys = append(pattern.Keys, key)

		// Check for more keys
		if !p.peekTokenIs(lexer.COMMA) {
			break
		}
		p.nextToken() // consume comma

		// Check for trailing comma before }
		if p.peekTokenIs(lexer.RBRACE) {
			break
		}

		// Move to next key or rest operator
		p.nextToken()
	}

	if !p.expectPeek(lexer.RBRACE) {
		return nil
	}

	return pattern
}
