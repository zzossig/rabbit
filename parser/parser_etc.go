package parser

import (
	"github.com/zzossig/xpath/ast"
	"github.com/zzossig/xpath/token"
)

func (p *Parser) parseArgumentList() ast.ArgumentList {
	li := ast.ArgumentList{}

	if !p.expectPeek(token.LPAREN) {
		return li
	}
	p.nextToken()

	for !p.curTokenIs(token.RPAREN) {
		arg := ast.Argument{}

		switch p.curToken.Type {
		case token.QUESTION:
			t := token.Token{Type: token.QUESTION, Literal: "?"}
			arg.ArgumentPlaceholder.Token = t
		default:
			arg.ExprSingle = p.parseExpression(LOWEST)
		}

		li.Args = append(li.Args, arg)
		p.nextToken()
		if p.curTokenIs(token.COMMA) {
			p.nextToken()
		}
	}

	return li
}
