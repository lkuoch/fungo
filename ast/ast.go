package ast

import (
	"bytes"
	"fmt"
	"fungo/token"
	"strings"
)

type Node interface {
	TokenLiteral() string
	String() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Identifier struct {
	Token token.Token
	Value string
}

type IntergerLiteral struct {
	Token token.Token
	Value int64
}

type LetStatement struct {
	Token token.Token
	Name  *Identifier
	Value Expression
}

type ReturnStatement struct {
	Token       token.Token
	ReturnValue Expression
}

type ExpressionStatement struct {
	Token      token.Token
	Expression Expression
}

type BlockStatement struct {
	Token      token.Token
	Statements []Statement
}

type Boolean struct {
	Token token.Token
	Value bool
}

type PrefixExpression struct {
	Token    token.Token
	Operator string
	Right    Expression
}

type InfixExpression struct {
	Token    token.Token
	Right    Expression
	Operator string
	Left     Expression
}

// if (<condition>) { <IfCondition> } else { <ElseCondition> }
type IfExpression struct {
	Token         token.Token
	Condition     Expression
	IfCondition   *BlockStatement
	ElseCondition *BlockStatement
}

// fn <parameters> <block statement>
type FunctionLiteral struct {
	Token      token.Token
	Parameters []*Identifier
	Body       *BlockStatement
}

type Program struct {
	Statements []Statement
}

func (i *InfixExpression) expressionNode() {}
func (i *InfixExpression) TokenLiteral() string {
	return i.Token.Literal
}
func (i *InfixExpression) String() string {
	return fmt.Sprintf("(%s %s %s)", i.Left.String(), i.Operator, i.Right.String())
}

func (p *PrefixExpression) expressionNode() {}
func (p *PrefixExpression) TokenLiteral() string {
	return p.Token.Literal
}
func (p *PrefixExpression) String() string {
	return fmt.Sprintf("(%s%s)", p.Operator, p.Right.String())
}

func (i *IfExpression) expressionNode() {}
func (i *IfExpression) TokenLiteral() string {
	return i.Token.Literal
}
func (i *IfExpression) String() string {
	var out bytes.Buffer

	if i.ElseCondition != nil {
		out.WriteString("if" + i.Condition.String() + " " + i.IfCondition.String() + "else " + i.ElseCondition.String())
	} else {
		out.WriteString("if" + i.Condition.String() + " " + i.IfCondition.String())
	}

	return out.String()
}

func (i *Identifier) expressionNode() {}
func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}
func (i *Identifier) String() string {
	return i.Value
}

func (i *IntergerLiteral) expressionNode() {}
func (i *IntergerLiteral) TokenLiteral() string {
	return i.Token.Literal
}
func (i *IntergerLiteral) String() string {
	return i.Token.Literal
}

func (l *LetStatement) statementNode() {}
func (l *LetStatement) TokenLiteral() string {
	return l.Token.Literal
}
func (l *LetStatement) String() string {
	var out bytes.Buffer

	out.WriteString(l.TokenLiteral() + " " + l.Name.String() + " = " + l.Value.String() + ";")

	return out.String()
}

func (r *ReturnStatement) statementNode() {}
func (r *ReturnStatement) TokenLiteral() string {
	return r.Token.Literal
}
func (r *ReturnStatement) String() string {
	var out bytes.Buffer

	out.WriteString(r.TokenLiteral() + " " + r.ReturnValue.String() + ";")

	return out.String()
}

func (e *ExpressionStatement) statementNode() {}
func (e *ExpressionStatement) TokenLiteral() string {
	return e.Token.Literal
}
func (e *ExpressionStatement) String() string {
	if e.Expression != nil {
		return e.Expression.String()
	}

	return ""
}

func (b *BlockStatement) statementNode() {}
func (b *BlockStatement) TokenLiteral() string {
	return b.Token.Literal
}
func (b *BlockStatement) String() string {
	var out bytes.Buffer

	for _, s := range b.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

func (f *FunctionLiteral) expressionNode() {}
func (f *FunctionLiteral) TokenLiteral() string {
	return f.Token.Literal
}
func (f *FunctionLiteral) String() string {
	var out bytes.Buffer

	params := []string{}

	out.WriteString(f.TokenLiteral() + "(" + strings.Join(params, ", ") + ")" + f.Body.String())

	return out.String()
}

func (b *Boolean) expressionNode() {}
func (b *Boolean) TokenLiteral() string {
	return b.Token.Literal
}
func (b *Boolean) String() string {
	return b.TokenLiteral()
}

func (p *Program) String() string {
	var out bytes.Buffer

	for _, s := range p.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}
