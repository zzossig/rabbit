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

func (p *Parser) parseStringLiteral() ast.ExprSingle {
	return &ast.StringLiteral{Value: p.curToken.Literal}
}

func (p *Parser) parseGroupedExpr() ast.ExprSingle {
	p.nextToken()

	tokens := p.collectTokenUntil(token.RPAREN)
	isCommaExist := false
	for _, t := range tokens {
		if t.Type == token.COMMA {
			isCommaExist = true
			break
		}
	}

	if isCommaExist {
		return p.parseSequenceExpr()
	}

	expr := p.parseExpression(LOWEST)
	if p.peekTokenIs(token.COMMA) || p.peekTokenIs(token.RPAREN) {
		p.nextToken()
	}
	return expr
}

func (p *Parser) parseSequenceExpr() ast.ExprSingle {
	expr := &ast.Expr{}

	for !p.curTokenIs(token.RPAREN) {
		e := p.parseExpression(LOWEST)
		expr.Exprs = append(expr.Exprs, e)

		p.nextToken()
		if p.curTokenIs(token.COMMA) {
			p.nextToken()
		}
	}

	return expr
}

func (p *Parser) parseAdditiveExpr(left ast.ExprSingle) ast.ExprSingle {
	expr := &ast.AdditiveExpr{
		LeftExpr: left,
		Token:    p.curToken,
	}

	precedence := p.curPrecedence()
	p.nextToken()
	right := p.parseExpression(precedence)

	if right != nil {
		expr.RightExpr = right
	}

	return expr
}

func (p *Parser) parseMultiplicativeExpr(left ast.ExprSingle) ast.ExprSingle {
	expr := &ast.MultiplicativeExpr{
		LeftExpr: left,
		Token:    p.curToken,
	}

	precedence := p.curPrecedence()
	p.nextToken()
	right := p.parseExpression(precedence)

	if right != nil {
		expr.RightExpr = right
	}

	return expr
}

func (p *Parser) parsePrefixExpr() ast.ExprSingle {
	expr := &ast.UnaryExpr{Token: p.curToken}

	p.nextToken()

	expr.ExprSingle = p.parseExpression(UNARY)

	return expr
}

func (p *Parser) parseCurlyArrayExpr() ast.ExprSingle {
	expr := &ast.CurlyArrayConstructor{}
	enclosedExpr := &ast.EnclosedExpr{}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	for !p.curTokenIs(token.RBRACE) {
		e := p.parseExpression(LOWEST)
		if e != nil {
			enclosedExpr.Exprs = append(enclosedExpr.Exprs, e)
		}
		p.nextToken()
	}

	expr.ExprSingle = enclosedExpr
	return expr
}

func (p *Parser) parseSquareArrayExpr() ast.ExprSingle {
	expr := &ast.SquareArrayConstructor{}

	p.nextToken()

	for !p.curTokenIs(token.RBRACKET) {
		e := p.parseExpression(LOWEST)
		if e != nil {
			expr.Exprs = append(expr.Exprs, e)
		}
		p.nextToken()
	}

	return expr
}

func (p *Parser) parseArrowExpr(left ast.ExprSingle) ast.ExprSingle {
	expr := &ast.ArrowExpr{ExprSingle: left}

	for p.curTokenIs(token.ARROW) || p.expectPeek(token.ARROW) {
		b := ast.ArrowBinding{}
		precedence := p.curPrecedence()
		p.nextToken()

		switch p.curToken.Type {
		case token.VAR:
			b.VarName.SetValue(p.curToken.Literal)
			b.TypeID = 2
		case token.LPAREN:
			e := p.parseExpression(precedence)
			b.ParenthesizedExpr.Exprs = append(b.ParenthesizedExpr.Exprs, e)
			b.TypeID = 3
		default:
			b.EQName.SetValue(p.curToken.Literal)
			b.TypeID = 1
		}

		b.ArgumentList = p.parseArgumentList()
		expr.Bindings = append(expr.Bindings, b)
	}

	return expr
}

func (p *Parser) parseMapExpr(left ast.ExprSingle) ast.ExprSingle {
	expr := &ast.SimpleMapExpr{LeftExpr: left, Token: p.curToken}

	precedence := p.curPrecedence()
	p.nextToken()
	right := p.parseExpression(precedence)

	if right != nil {
		expr.RightExpr = right
	}

	return expr
}

func (p *Parser) parseComparisonExpr(left ast.ExprSingle) ast.ExprSingle {
	expr := &ast.ComparisonExpr{LeftExpr: left}
	expr.SetToken(p.curToken)

	precedence := p.curPrecedence()
	p.nextToken()
	right := p.parseExpression(precedence)

	if right != nil {
		expr.RightExpr = right
	}

	return expr
}
