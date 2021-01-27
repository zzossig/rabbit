package parser

import (
	"strconv"

	"github.com/zzossig/xpath/ast"
	"github.com/zzossig/xpath/token"
)

func (p *Parser) parseExpression(precedence int) ast.ExprSingle {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		// TODO error
		return nil
	}
	leftExp := prefix()

	for !p.peekTokenIs(token.COMMA) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}

		p.nextToken()

		leftExp = infix(leftExp)
	}

	return leftExp
}

func (p *Parser) parseIntegerLiteral() ast.ExprSingle {
	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		// TODO error
		return nil
	}
	return &ast.IntegerLiteral{Value: int(value)}
}

func (p *Parser) parseGroupedExpr() ast.ExprSingle {
	p.nextToken()

	exp := p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return exp
}

func (p *Parser) parseAdditiveExpr(left ast.ExprSingle) ast.ExprSingle {
	expression := &ast.AdditiveExpr{
		LeftExpr: left,
		Token:    p.curToken,
	}

	precedence := p.curPrecedence()
	p.nextToken()
	right := p.parseExpression(precedence)

	if right != nil {
		expression.RightExpr = right
	}

	return expression
}

func (p *Parser) parseMultiplicativeExpr(left ast.ExprSingle) ast.ExprSingle {
	expression := &ast.MultiplicativeExpr{
		LeftExpr: left,
		Token:    p.curToken,
	}

	precedence := p.curPrecedence()
	p.nextToken()
	right := p.parseExpression(precedence)

	if right != nil {
		expression.RightExpr = right
	}

	return expression
}

func (p *Parser) parsePrefixExpr() ast.ExprSingle {
	expression := &ast.UnaryExpr{Token: p.curToken}

	p.nextToken()

	expression.ExprSingle = p.parseExpression(UNARY)

	return expression
}
