package parser

import (
	"github.com/zzossig/xpath/ast"
	"github.com/zzossig/xpath/token"
)

func (p *Parser) parseItem() []ast.Item {
	switch p.curToken.Type {
	case token.COMMA, token.RPAREN:
		return p.pass()
	default:
		return p.parseExpressionItem()
	}
}

func (p *Parser) pass() []ast.Item {
	return []ast.Item{}
}

func (p *Parser) parseExpressionItem() []ast.Item {
	item := &ast.ExprItem{Token: p.curToken}
	item.Expression = p.parseExprSingle(LOWEST)

	if p.peekTokenIs(token.COMMA) {
		p.nextToken()
	}

	return []ast.Item{item}
}
