package parser

import (
	"github.com/zzossig/xpath/ast"
	"github.com/zzossig/xpath/token"
)

func (p *Parser) parseItem() []ast.Item {
	switch p.curToken.Type {
	case token.LPAREN:
		return p.parseGroupedItem()
	default:
		return p.parseExpressionItem()
	}
}

func (p *Parser) parseGroupedItem() []ast.Item {
	p.nextToken()
	items := []ast.Item{}

	for !p.curTokenIs(token.RPAREN) {
		item := p.parseItem()
		if items != nil {
			items = append(items, item...)
		}
		p.nextToken()
	}

	return items
}

func (p *Parser) parseExpressionItem() []ast.Item {
	item := &ast.ExprItem{Token: p.curToken}
	item.Expression = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.COMMA) {
		p.nextToken()
	}

	return []ast.Item{item}
}
