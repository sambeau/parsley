package parser

import (
	"fmt"
	"strconv"

	"pars/pkg/ast"
	"pars/pkg/lexer"
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
	lexer.COMMA:    COMMA_PREC,
	lexer.OR:       LOGIC_OR,
	lexer.AND:      LOGIC_AND,
	lexer.EQ:       EQUALS,
	lexer.NOT_EQ:   EQUALS,
	lexer.LT:       LESSGREATER,
	lexer.GT:       LESSGREATER,
	lexer.LTE:      LESSGREATER,
	lexer.GTE:      LESSGREATER,
	lexer.PLUS:     SUM,
	lexer.MINUS:    SUM,
	lexer.PLUSPLUS: CONCAT,
	lexer.SLASH:    PRODUCT,
	lexer.ASTERISK: PRODUCT,
	lexer.PERCENT:  PRODUCT,
	lexer.LBRACKET: INDEX,
	lexer.LPAREN:   CALL,
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
	p.registerPrefix(lexer.BANG, p.parsePrefixExpression)
	p.registerPrefix(lexer.MINUS, p.parsePrefixExpression)
	p.registerPrefix(lexer.TRUE, p.parseBoolean)
	p.registerPrefix(lexer.FALSE, p.parseBoolean)
	p.registerPrefix(lexer.LPAREN, p.parseGroupedExpression)
	p.registerPrefix(lexer.LBRACKET, p.parseSquareBracketArrayLiteral)
	p.registerPrefix(lexer.IF, p.parseIfExpression)
	p.registerPrefix(lexer.FUNCTION, p.parseFunctionLiteral)
	p.registerPrefix(lexer.FOR, p.parseForExpression)

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
	p.registerInfix(lexer.PLUSPLUS, p.parseInfixExpression)
	p.registerInfix(lexer.COMMA, p.parseArrayLiteral)
	p.registerInfix(lexer.LPAREN, p.parseCallExpression)
	p.registerInfix(lexer.LBRACKET, p.parseIndexOrSliceExpression)

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

	if !p.expectPeek(lexer.IDENT) {
		return nil
	}

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

	lit.Parameters = p.parseFunctionParameters()

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

// parseForExpression parses for expressions
// Two forms: for(array) func  OR  for(var in array) body
func (p *Parser) parseForExpression() ast.Expression {
	expression := &ast.ForExpression{Token: p.curToken}

	if !p.expectPeek(lexer.LPAREN) {
		return nil
	}

	p.nextToken()

	// Check if this is the "for(var in array)" form
	// We need to peek ahead to see if there's an IN token
	if p.peekTokenIs(lexer.IN) {
		// Parse: for(var in array) body
		if p.curToken.Type != lexer.IDENT {
			p.peekError(lexer.IDENT)
			return nil
		}
		expression.Variable = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

		p.nextToken() // move to IN
		p.nextToken() // move past IN to array expression

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
			Token:      p.curToken,
			Parameters: []*ast.Identifier{expression.Variable},
			Body:       p.parseBlockStatement(),
		}
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
		exp.End = p.parseExpression(LOWEST)
	}

	if !p.expectPeek(lexer.RBRACKET) {
		return nil
	}

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
