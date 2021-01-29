package parser

import (
	"github.com/zzossig/xpath/ast"
	"github.com/zzossig/xpath/token"
)

func (p *Parser) parseItem() []ast.Item {
	switch p.curToken.Type {
	case token.ARRAY:
		return p.parseCurlyArrayItem()
	case token.LPAREN:
		return p.parseGroupedItem()
	case token.LBRACKET:
		return p.parseSquareArrayItem()
	case token.EQ, token.NE, token.LT, token.LE, token.GT, token.GE,
		token.EQV, token.NEV, token.LTV, token.LEV, token.GTV, token.GEV,
		token.IS, token.DLT, token.DGT:
		return p.parseComparisonItem()
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
	item := &ast.ExprItem{Token: p.curToken}
	item.Expression = p.parseCurlyArrayExpr()

	return []ast.Item{item}
}

func (p *Parser) parseSquareArrayItem() []ast.Item {
	item := &ast.ExprItem{Token: p.curToken}
	item.Expression = p.parseSquareArrayExpr()

	return []ast.Item{item}
}

func (p *Parser) parseGroupedItem() []ast.Item {
	item := &ast.ExprItem{Token: p.curToken}
	item.Expression = p.parseGroupedExpr()

	return []ast.Item{item}
}

func (p *Parser) parseComparisonItem() []ast.Item {
	if len(p.xpath.Items) == 0 {
		// TODO error
		return []ast.Item{}
	}

	item := &ast.ExprItem{Token: p.curToken}
	it, ok := p.xpath.Items[len(p.xpath.Items)-1].(*ast.ExprItem)
	if !ok {
		// TODO error
		return []ast.Item{}
	}

	p.xpath.Items = p.xpath.Items[:len(p.xpath.Items)-1]
	expr := p.parseComparisonExpr(it.Expression)
	item.Expression = expr

	return []ast.Item{item}
}
