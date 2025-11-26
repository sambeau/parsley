package ast

import (
	"bytes"
	"strings"

	"github.com/sambeau/parsley/pkg/lexer"
)

// Node represents any node in the AST
type Node interface {
	TokenLiteral() string
	String() string
}

// Statement represents statement nodes
type Statement interface {
	Node
	statementNode()
}

// Expression represents expression nodes
type Expression interface {
	Node
	expressionNode()
}

// Program represents the root node of every AST
type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}
	return ""
}

func (p *Program) String() string {
	var out bytes.Buffer

	for _, s := range p.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

// LetStatement represents let statements like 'let x = 5;' or 'let x,y,z = 1,2,3;'
type LetStatement struct {
	Token       lexer.Token               // the lexer.LET token
	Name        *Identifier               // single name (for backwards compatibility)
	Names       []*Identifier             // multiple names for array destructuring
	DictPattern *DictDestructuringPattern // pattern for dictionary destructuring
	Value       Expression
}

func (ls *LetStatement) statementNode()       {}
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }
func (ls *LetStatement) String() string {
	var out bytes.Buffer

	out.WriteString(ls.TokenLiteral() + " ")
	if ls.DictPattern != nil {
		out.WriteString(ls.DictPattern.String())
	} else if len(ls.Names) > 0 {
		for i, name := range ls.Names {
			if i > 0 {
				out.WriteString(", ")
			}
			out.WriteString(name.String())
		}
	} else {
		out.WriteString(ls.Name.String())
	}
	out.WriteString(" = ")

	if ls.Value != nil {
		out.WriteString(ls.Value.String())
	}

	out.WriteString(";")
	return out.String()
}

// AssignmentStatement represents assignment statements like 'x = 5;' or 'x,y,z = 1,2,3;'
type AssignmentStatement struct {
	Token       lexer.Token               // the identifier token
	Name        *Identifier               // single name (for backwards compatibility)
	Names       []*Identifier             // multiple names for array destructuring
	DictPattern *DictDestructuringPattern // pattern for dictionary destructuring
	Value       Expression
}

func (as *AssignmentStatement) statementNode()       {}
func (as *AssignmentStatement) TokenLiteral() string { return as.Token.Literal }
func (as *AssignmentStatement) String() string {
	var out bytes.Buffer

	if as.DictPattern != nil {
		out.WriteString(as.DictPattern.String())
	} else if len(as.Names) > 0 {
		for i, name := range as.Names {
			if i > 0 {
				out.WriteString(", ")
			}
			out.WriteString(name.String())
		}
	} else {
		out.WriteString(as.Name.String())
	}
	out.WriteString(" = ")

	if as.Value != nil {
		out.WriteString(as.Value.String())
	}

	out.WriteString(";")
	return out.String()
}

// ReturnStatement represents return statements like 'return 5;'
type ReturnStatement struct {
	Token       lexer.Token // the 'return' token
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode()       {}
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }
func (rs *ReturnStatement) String() string {
	var out bytes.Buffer

	out.WriteString(rs.TokenLiteral() + " ")

	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String())
	}

	out.WriteString(";")
	return out.String()
}

// ExpressionStatement represents expression statements
type ExpressionStatement struct {
	Token      lexer.Token // the first token of the expression
	Expression Expression
}

func (es *ExpressionStatement) statementNode()       {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

// BlockStatement represents block statements like '{...}'
type BlockStatement struct {
	Token      lexer.Token // the '{' token
	Statements []Statement
}

func (bs *BlockStatement) statementNode()       {}
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Literal }
func (bs *BlockStatement) String() string {
	var out bytes.Buffer

	for _, s := range bs.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

// Identifier represents identifier expressions
type Identifier struct {
	Token lexer.Token // the lexer.IDENT token
	Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) String() string       { return i.Value }

// IntegerLiteral represents integer literals
type IntegerLiteral struct {
	Token lexer.Token // the lexer.INT token
	Value int64
}

func (il *IntegerLiteral) expressionNode()      {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }
func (il *IntegerLiteral) String() string       { return il.Token.Literal }

// FloatLiteral represents floating-point literals
type FloatLiteral struct {
	Token lexer.Token // the lexer.FLOAT token
	Value float64
}

func (fl *FloatLiteral) expressionNode()      {}
func (fl *FloatLiteral) TokenLiteral() string { return fl.Token.Literal }
func (fl *FloatLiteral) String() string       { return fl.Token.Literal }

// StringLiteral represents string literals
type StringLiteral struct {
	Token lexer.Token // the lexer.STRING token
	Value string
}

func (sl *StringLiteral) expressionNode()      {}
func (sl *StringLiteral) TokenLiteral() string { return sl.Token.Literal }
func (sl *StringLiteral) String() string       { return sl.Token.Literal }

// TemplateLiteral represents template literals with interpolation
type TemplateLiteral struct {
	Token lexer.Token // the lexer.TEMPLATE token
	Value string      // the raw template string
}

func (tl *TemplateLiteral) expressionNode()      {}
func (tl *TemplateLiteral) TokenLiteral() string { return tl.Token.Literal }
func (tl *TemplateLiteral) String() string       { return "`" + tl.Value + "`" }

// RegexLiteral represents regular expression literals like /pattern/flags
type RegexLiteral struct {
	Token   lexer.Token // the lexer.REGEX token
	Pattern string      // the regex pattern
	Flags   string      // the regex flags
}

func (rl *RegexLiteral) expressionNode()      {}
func (rl *RegexLiteral) TokenLiteral() string { return rl.Token.Literal }
func (rl *RegexLiteral) String() string {
	return "/" + rl.Pattern + "/" + rl.Flags
}

// DatetimeLiteral represents datetime literals like @2024-12-25T14:30:00Z
type DatetimeLiteral struct {
	Token lexer.Token // the lexer.DATETIME_LITERAL token
	Value string      // the ISO-8601 datetime string
}

func (dl *DatetimeLiteral) expressionNode()      {}
func (dl *DatetimeLiteral) TokenLiteral() string { return dl.Token.Literal }
func (dl *DatetimeLiteral) String() string       { return "@" + dl.Value }

// DurationLiteral represents duration literals like @2h30m, @7d, @1y6mo
type DurationLiteral struct {
	Token lexer.Token // the lexer.DURATION_LITERAL token
	Value string      // the duration string (e.g., "2h30m", "7d")
}

func (dr *DurationLiteral) expressionNode()      {}
func (dr *DurationLiteral) TokenLiteral() string { return dr.Token.Literal }
func (dr *DurationLiteral) String() string       { return "@" + dr.Value }

// TagLiteral represents singleton tags like <input type="text" />
type TagLiteral struct {
	Token lexer.Token // the lexer.TAG token
	Raw   string      // the raw tag content (everything between < and />)
}

func (tg *TagLiteral) expressionNode()      {}
func (tg *TagLiteral) TokenLiteral() string { return tg.Token.Literal }
func (tg *TagLiteral) String() string       { return "<" + tg.Raw + " />" }

// TagPairExpression represents paired tags like <div>content</div>
type TagPairExpression struct {
	Token    lexer.Token // the lexer.TAG_START token
	Name     string      // tag name (empty string for grouping tags <>)
	Props    string      // raw props content
	Contents []Node      // mixed content: text nodes, expressions, nested tags
}

func (tp *TagPairExpression) expressionNode()      {}
func (tp *TagPairExpression) TokenLiteral() string { return tp.Token.Literal }
func (tp *TagPairExpression) String() string {
	var out bytes.Buffer
	if tp.Name == "" {
		out.WriteString("<>")
	} else {
		out.WriteString("<" + tp.Name)
		if tp.Props != "" {
			out.WriteString(" " + tp.Props)
		}
		out.WriteString(">")
	}
	for _, content := range tp.Contents {
		out.WriteString(content.String())
	}
	if tp.Name == "" {
		out.WriteString("</>")
	} else {
		out.WriteString("</" + tp.Name + ">")
	}
	return out.String()
}

// TextNode represents raw text content within tags
type TextNode struct {
	Token lexer.Token // the lexer.TAG_TEXT token
	Value string      // the text content
}

func (tn *TextNode) expressionNode()      {}
func (tn *TextNode) TokenLiteral() string { return tn.Token.Literal }
func (tn *TextNode) String() string       { return tn.Value }

// Boolean represents boolean literals
type Boolean struct {
	Token lexer.Token // the lexer.TRUE or lexer.FALSE token
	Value bool
}

func (b *Boolean) expressionNode()      {}
func (b *Boolean) TokenLiteral() string { return b.Token.Literal }
func (b *Boolean) String() string       { return b.Token.Literal }

// PrefixExpression represents prefix expressions like '!x' or '-x'
type PrefixExpression struct {
	Token    lexer.Token // the prefix token, e.g. !
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) expressionNode()      {}
func (pe *PrefixExpression) TokenLiteral() string { return pe.Token.Literal }
func (pe *PrefixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")

	return out.String()
}

// InfixExpression represents infix expressions like 'x + y'
type InfixExpression struct {
	Token    lexer.Token // the operator token, e.g. +
	Left     Expression
	Operator string
	Right    Expression
}

func (oe *InfixExpression) expressionNode()      {}
func (oe *InfixExpression) TokenLiteral() string { return oe.Token.Literal }
func (oe *InfixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(oe.Left.String())
	out.WriteString(" " + oe.Operator + " ")
	out.WriteString(oe.Right.String())
	out.WriteString(")")

	return out.String()
}

// IfExpression represents if expressions
type IfExpression struct {
	Token       lexer.Token // the 'if' token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (ie *IfExpression) expressionNode()      {}
func (ie *IfExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IfExpression) String() string {
	var out bytes.Buffer

	out.WriteString("if")
	out.WriteString(ie.Condition.String())
	out.WriteString(" ")
	out.WriteString(ie.Consequence.String())

	if ie.Alternative != nil {
		out.WriteString("else ")
		out.WriteString(ie.Alternative.String())
	}

	return out.String()
}

// FunctionLiteral represents function literals
// FunctionParameter represents a function parameter (identifier, array pattern, or dict pattern)
type FunctionParameter struct {
	Ident        *Identifier               // simple identifier parameter
	ArrayPattern []*Identifier             // array destructuring pattern
	DictPattern  *DictDestructuringPattern // dict destructuring pattern
}

func (fp *FunctionParameter) String() string {
	if fp.DictPattern != nil {
		return fp.DictPattern.String()
	}
	if len(fp.ArrayPattern) > 0 {
		var out bytes.Buffer
		out.WriteString("[")
		for i, p := range fp.ArrayPattern {
			if i > 0 {
				out.WriteString(", ")
			}
			out.WriteString(p.String())
		}
		out.WriteString("]")
		return out.String()
	}
	return fp.Ident.String()
}

type FunctionLiteral struct {
	Token      lexer.Token          // the 'fn' token
	Parameters []*Identifier        // deprecated - kept for backwards compatibility
	Params     []*FunctionParameter // new parameter list supporting destructuring
	Body       *BlockStatement
}

func (fl *FunctionLiteral) expressionNode()      {}
func (fl *FunctionLiteral) TokenLiteral() string { return fl.Token.Literal }
func (fl *FunctionLiteral) String() string {
	var out bytes.Buffer

	params := []string{}
	if len(fl.Params) > 0 {
		for _, p := range fl.Params {
			params = append(params, p.String())
		}
	} else {
		// Fallback for old-style parameters
		for _, p := range fl.Parameters {
			params = append(params, p.String())
		}
	}

	out.WriteString(fl.TokenLiteral())
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") ")
	out.WriteString(fl.Body.String())

	return out.String()
}

// CallExpression represents function calls
type CallExpression struct {
	Token     lexer.Token // the '(' token
	Function  Expression  // Identifier or FunctionLiteral
	Arguments []Expression
}

func (ce *CallExpression) expressionNode()      {}
func (ce *CallExpression) TokenLiteral() string { return ce.Token.Literal }
func (ce *CallExpression) String() string {
	var out bytes.Buffer

	args := []string{}
	for _, a := range ce.Arguments {
		args = append(args, a.String())
	}

	out.WriteString(ce.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")

	return out.String()
}

// ArrayLiteral represents array literals like [1, 2, 3]
type ArrayLiteral struct {
	Token    lexer.Token // the first element token
	Elements []Expression
}

func (al *ArrayLiteral) expressionNode()      {}
func (al *ArrayLiteral) TokenLiteral() string { return al.Token.Literal }
func (al *ArrayLiteral) String() string {
	var out bytes.Buffer

	elements := []string{}
	for _, el := range al.Elements {
		elements = append(elements, el.String())
	}

	out.WriteString(strings.Join(elements, ", "))

	return out.String()
}

// ForExpression represents for expressions
// Two forms: for(array) func  OR  for(var in array) body
type ForExpression struct {
	Token         lexer.Token // the 'for' token
	Array         Expression  // the array to iterate over
	Function      Expression  // the function to apply (for simple form)
	Variable      *Identifier // the loop variable (for 'in' form) or key variable (for dict)
	ValueVariable *Identifier // the value variable (for dict 'key, value in dict' form)
	Body          Expression  // the body expression (for 'in' form)
}

func (fe *ForExpression) expressionNode()      {}
func (fe *ForExpression) TokenLiteral() string { return fe.Token.Literal }
func (fe *ForExpression) String() string {
	var out bytes.Buffer

	out.WriteString("for(")
	if fe.Variable != nil {
		out.WriteString(fe.Variable.String())
		out.WriteString(" in ")
	}
	out.WriteString(fe.Array.String())
	out.WriteString(")")

	if fe.Function != nil {
		out.WriteString(" ")
		out.WriteString(fe.Function.String())
	} else if fe.Body != nil {
		out.WriteString(" ")
		out.WriteString(fe.Body.String())
	}

	return out.String()
}

// IndexExpression represents array/string indexing like arr[0] or str[1]
type IndexExpression struct {
	Token lexer.Token // the '[' token
	Left  Expression  // the array or string being indexed
	Index Expression  // the index expression
}

func (ie *IndexExpression) expressionNode()      {}
func (ie *IndexExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IndexExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString("[")
	out.WriteString(ie.Index.String())
	out.WriteString("])")

	return out.String()
}

// SliceExpression represents array/string slicing like arr[1:4]
type SliceExpression struct {
	Token lexer.Token // the '[' token
	Left  Expression  // the array or string being sliced
	Start Expression  // the start index (can be nil)
	End   Expression  // the end index (can be nil)
}

func (se *SliceExpression) expressionNode()      {}
func (se *SliceExpression) TokenLiteral() string { return se.Token.Literal }
func (se *SliceExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(se.Left.String())
	out.WriteString("[")
	if se.Start != nil {
		out.WriteString(se.Start.String())
	}
	out.WriteString(":")
	if se.End != nil {
		out.WriteString(se.End.String())
	}
	out.WriteString("])")

	return out.String()
}

// DictionaryLiteral represents dictionary literals like { key: value, ... }
type DictionaryLiteral struct {
	Token lexer.Token // the '{' token
	Pairs map[string]Expression
}

func (dl *DictionaryLiteral) expressionNode()      {}
func (dl *DictionaryLiteral) TokenLiteral() string { return dl.Token.Literal }
func (dl *DictionaryLiteral) String() string {
	var out bytes.Buffer

	pairs := []string{}
	for key, value := range dl.Pairs {
		pairs = append(pairs, key+": "+value.String())
	}

	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")

	return out.String()
}

// DotExpression represents dot notation access like dict.key
type DotExpression struct {
	Token lexer.Token // the '.' token
	Left  Expression  // the object being accessed
	Key   string      // the property name
}

func (de *DotExpression) expressionNode()      {}
func (de *DotExpression) TokenLiteral() string { return de.Token.Literal }
func (de *DotExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(de.Left.String())
	out.WriteString(".")
	out.WriteString(de.Key)
	out.WriteString(")")

	return out.String()
}

// DeleteStatement represents delete dict.key or delete dict["key"]
type DeleteStatement struct {
	Token  lexer.Token // the 'delete' token
	Target Expression  // the property access expression to delete
}

func (ds *DeleteStatement) statementNode()       {}
func (ds *DeleteStatement) TokenLiteral() string { return ds.Token.Literal }
func (ds *DeleteStatement) String() string {
	var out bytes.Buffer

	out.WriteString("delete ")
	out.WriteString(ds.Target.String())
	out.WriteString(";")

	return out.String()
}

// DictDestructuringPattern represents a dictionary destructuring pattern like {a, b as c, ...rest}
type DictDestructuringPattern struct {
	Token lexer.Token             // the '{' token
	Keys  []*DictDestructuringKey // the keys to extract
	Rest  *Identifier             // optional rest identifier (for ...rest)
}

func (ddp *DictDestructuringPattern) expressionNode()      {}
func (ddp *DictDestructuringPattern) TokenLiteral() string { return ddp.Token.Literal }
func (ddp *DictDestructuringPattern) String() string {
	var out bytes.Buffer

	out.WriteString("{")
	for i, key := range ddp.Keys {
		if i > 0 {
			out.WriteString(", ")
		}
		out.WriteString(key.String())
	}
	if ddp.Rest != nil {
		if len(ddp.Keys) > 0 {
			out.WriteString(", ")
		}
		out.WriteString("...")
		out.WriteString(ddp.Rest.String())
	}
	out.WriteString("}")

	return out.String()
}

// DictDestructuringKey represents a single key in a dictionary destructuring pattern
// Can be: 'a' or 'a as b' or nested pattern
type DictDestructuringKey struct {
	Token  lexer.Token // the identifier token
	Key    *Identifier // the key name from the dictionary
	Alias  *Identifier // optional alias (for 'as' syntax)
	Nested Expression  // optional nested pattern for nested destructuring
}

func (ddk *DictDestructuringKey) expressionNode()      {}
func (ddk *DictDestructuringKey) TokenLiteral() string { return ddk.Token.Literal }
func (ddk *DictDestructuringKey) String() string {
	var out bytes.Buffer

	if ddk.Nested != nil {
		out.WriteString(ddk.Key.String())
		out.WriteString(": ")
		out.WriteString(ddk.Nested.String())
	} else if ddk.Alias != nil {
		out.WriteString(ddk.Key.String())
		out.WriteString(" as ")
		out.WriteString(ddk.Alias.String())
	} else {
		out.WriteString(ddk.Key.String())
	}

	return out.String()
}
