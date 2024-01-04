package ast

import (
	"fungo/token"
	"testing"

	"github.com/stretchr/testify/suite"
)

type AstTestSuite struct {
	suite.Suite
}

func TestAstTestSuite(t *testing.T) {
	suite.Run(t, &AstTestSuite{})
}

func (t *AstTestSuite) TestString() {
	program := &Program{
		Statements: []Statement{
			&LetStatement{
				Token: token.Token{Type: token.LET, Literal: "let"},
				Name: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "myVar"},
					Value: "myVar",
				},
				Value: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "anotherVar"},
					Value: "anotherVar",
				},
			},
		},
	}

	t.Equal("let myVar = anotherVar;", program.String())
}
