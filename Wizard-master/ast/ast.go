package ast //该包提供了ast的一些方法

import (
	"bytes"
	"strings"

	"my.com/myfile/token"
)

type Node interface { //定义了AST(语法树)中所有节点必须实现的两个方法
	TokenLiteral() string
	String() string
}

// All statement nodes implement this
// 在Go语言中，一个类型被视为实现了某个接口，如果它实现了该接口中声明的所有方法。
// 所以Statement接口和Expression接口中没有任何实际功能的方法是为了方便扩展
type Statement interface { //扩展了Node接口，添加了statement()方法
	Node
	statementNode()
}

type Expression interface { //Ecpression接口，扩展了Node接口，添加了一个expressionNode()方法
	Node
	expressionNode()
}

type Program struct { //每个程序可以有许多语句
	Statements []Statement
}

func (p *Program) TokenLiteral() string { //接受一个程序，然后返回这个程序第一个语句
	if len(p.Statements) > 0 { //如果程序中存在语句，调用Node接口的TokenLiteral，反之，则返回空字符串
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

func (p *Program) String() string { //
	var out bytes.Buffer

	for _, s := range p.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

// Statements
type LetStatement struct { //
	Token token.Token // the token.LET token
	Name  *Identifier
	Value Expression
}

func (ls *LetStatement) statementNode()       {}
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }
func (ls *LetStatement) String() string {
	var out bytes.Buffer

	out.WriteString(ls.TokenLiteral() + " ")
	out.WriteString(ls.Name.String())
	out.WriteString(" = ")

	if ls.Value != nil {
		out.WriteString(ls.Value.String())
	}

	out.WriteString(";")

	return out.String()
}

type ReturnStatement struct {
	Token       token.Token // the 'return' token
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

type BreakStatement struct {
	Token token.Token // the 'break' token
}

func (bs *BreakStatement) statementNode()       {}
func (bs *BreakStatement) TokenLiteral() string { return bs.Token.Literal }
func (bs *BreakStatement) String() string       { return bs.Token.Literal }

type ContinueStatement struct {
	Token token.Token // the 'continue' token
}

func (cs *ContinueStatement) statementNode()       {}
func (cs *ContinueStatement) TokenLiteral() string { return cs.Token.Literal }
func (cs *ContinueStatement) String() string       { return cs.Token.Literal }

type ExpressionStatement struct { //表达式结构体
	Token      token.Token // the first token of the expression
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

type BlockStatement struct {
	Token      token.Token // the { token
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

// Expressions
type Identifier struct {
	Token token.Token // the token.IDENT token
	Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) String() string       { return i.Value }

type Boolean struct {
	Token token.Token
	Value bool
}

func (b *Boolean) expressionNode()      {}
func (b *Boolean) TokenLiteral() string { return b.Token.Literal }
func (b *Boolean) String() string       { return b.Token.Literal }

type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (il *IntegerLiteral) expressionNode()      {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }
func (il *IntegerLiteral) String() string       { return il.Token.Literal }

type PrefixExpression struct {
	Token    token.Token // The prefix token, e.g. !
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

type InfixExpression struct {
	Token    token.Token // The operator token, e.g. +
	Left     Expression
	Operator string
	Right    Expression
}

func (ie *InfixExpression) expressionNode()      {}
func (ie *InfixExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *InfixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString(" " + ie.Operator + " ")
	out.WriteString(ie.Right.String())
	out.WriteString(")")

	return out.String()
}

type IfExpression struct {
	Token       token.Token // The 'if' token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (ie *IfExpression) statementNode()       {}
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

type FunctionLiteral struct {
	Token      token.Token // The 'fn' token
	Parameters []*Identifier
	Body       *BlockStatement
}

func (fl *FunctionLiteral) expressionNode()      {}
func (fl *FunctionLiteral) TokenLiteral() string { return fl.Token.Literal }
func (fl *FunctionLiteral) String() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range fl.Parameters {
		params = append(params, p.String())
	}

	out.WriteString(fl.TokenLiteral())
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") ")
	out.WriteString(fl.Body.String())

	return out.String()
}

type CallExpression struct {
	Token     token.Token // The '(' token
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

type ForExpression struct { //For的抽象语法树
	Token      token.Token
	Initialize Expression //可以为空
	Condition  Expression
	Cycleop    Expression
	Body       *BlockStatement
}

/*
定义For语句的结构:

	for int i = 0:i < 10:i++ {
		//Block
	}
*/
func (fs *ForExpression) expressionNode()      {}
func (fs *ForExpression) TokenLiteral() string { return fs.Token.Literal }
func (fs *ForExpression) String() string {
	var out bytes.Buffer
	out.WriteString("For")
	out.WriteString(fs.Initialize.String())
	out.WriteString(fs.Condition.String())
	out.WriteString(fs.Cycleop.String())
	//out.WriteString(") ")
	out.WriteString(fs.Body.String())
	return out.String()
}

type WhileExpression struct {
	Token     token.Token
	Condition Expression
	Body      *BlockStatement
}

func (fs *WhileExpression) expressionNode()      {}
func (fs *WhileExpression) TokenLiteral() string { return fs.Token.Literal }
func (fs *WhileExpression) String() string {
	var out bytes.Buffer
	out.WriteString("while")
	out.WriteString(fs.Condition.String())
	//out.WriteString(") ")
	out.WriteString(fs.Body.String())
	return out.String()
}

type StringLiteral struct {
	Token token.Token
	Value string
}

func (sl *StringLiteral) expressionNode()      {}
func (sl *StringLiteral) TokenLiteral() string { return sl.Token.Literal }
func (sl *StringLiteral) String() string       { return sl.Token.Literal }

//type AssignExpression struct {
//	Token token.Token // The = token
//	Name  Expression  // Name of the variable being assigned
//	Value Expression  // Value to be assigned
//}
//
//func (ae *AssignExpression) statementNode()       {}
//func (ae *AssignExpression) expressionNode()      {}
//func (ae *AssignExpression) TokenLiteral() string { return ae.Token.Literal }
//func (ae *AssignExpression) String() string {
//	var out bytes.Buffer
//
//	out.WriteString(ae.Name.String())
//	out.WriteString(" = ")
//	if ae.Value != nil {
//		out.WriteString(ae.Value.String())
//	}
//
//	return out.String()
//}