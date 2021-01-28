package parser

import (
	"github.com/zzossig/xpath/ast"
	"github.com/zzossig/xpath/token"
)

func (p *Parser) parseItem() []ast.Item {
	switch p.curToken.Type {
	case token.ARRAY:
		return p.parseCurlyArrayItem()
	case token.LBRACKET:
		return p.parseSquareArrayItem()
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
	item.Expression = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.COMMA) {
		p.nextToken()
	}

	return []ast.Item{item}
}

func (p *Parser) parseCurlyArrayItem() []ast.Item {
	item := &ast.ExprItem{Token: token.Token{Type: token.ARRAY, Literal: "array"}}
	item.Expression = p.parseCurlyArrayExpr()

	return []ast.Item{item}
}

func (p *Parser) parseSquareArrayItem() []ast.Item {
	item := &ast.ExprItem{Token: token.Token{Type: token.LBRACKET, Literal: "["}}
	item.Expression = p.parseSquareArrayExpr()

	return []ast.Item{item}
}
