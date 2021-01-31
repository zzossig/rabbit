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

func (p *Parser) parseDecimalLiteral() ast.ExprSingle {
	value, err := strconv.ParseFloat(p.curToken.Literal, 64)
	if err != nil {
		// TODO error
		return nil
	}
	return &ast.DecimalLiteral{Value: value}
}

func (p *Parser) parseDoubleLiteral() ast.ExprSingle {
	value, err := strconv.ParseFloat(p.curToken.Literal, 64)
	if err != nil {
		// TODO error
		return nil
	}
	return &ast.DoubleLiteral{Value: value}
}

func (p *Parser) parseStringLiteral() ast.ExprSingle {
	return &ast.StringLiteral{Value: p.curToken.Literal}
}

func (p *Parser) parseVariable() ast.ExprSingle {
	vr := &ast.VarRef{}

	p.nextToken()
	vr.VarName = p.parseEQName()

	return vr
}

func (p *Parser) parseGroupedExpr() ast.ExprSingle {
	if p.hasComma() {
		return p.parseSequenceExpr()
	}

	p.nextToken()

	expr := p.parseExpression(LOWEST)
	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	return expr
}

func (p *Parser) parseSequenceExpr() ast.ExprSingle {
	p.nextToken()
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
	expr := &ast.UnaryExpr{}

	if p.curTokenIs(token.PLUS) {
		expr.Token = token.Token{Type: token.UPLUS, Literal: "(+)"}
	} else {
		expr.Token = token.Token{Type: token.UMINUS, Literal: "(-)"}
	}

	p.nextToken()

	expr.ExprSingle = p.parseExpression(UNARY)

	return expr
}

func (p *Parser) parseUnaryLookupExpr() ast.ExprSingle {
	expr := &ast.UnaryLookup{}

	if p.curTokenIs(token.QUESTION) {
		expr.Token = token.Token{Type: token.UQUESTION, Literal: "(?)"}
	}

	p.nextToken()

	if p.curTokenIs(token.ASTERISK) {
		expr.TypeID = 4
	} else if p.curTokenIs(token.INT) {
		expr.TypeID = 2

		i, _ := strconv.Atoi(p.curToken.Literal)
		expr.IntegerLiteral.Value = i
	} else if p.curTokenIs(token.LPAREN) {
		expr.TypeID = 3

		e := p.parseExpression(LOWEST)
		expr.ParenthesizedExpr.Exprs = append(expr.ParenthesizedExpr.Exprs, e)
	} else {
		expr.TypeID = 1
		expr.NCName.SetValue(p.curToken.Literal)
	}

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

func (p *Parser) parseMapExpr() ast.ExprSingle {
	expr := &ast.MapConstructor{}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}
	if p.peekTokenIs(token.RBRACE) {
		p.nextToken()
		return expr
	}

	for {
		p.nextToken()

		entry := ast.MapConstructorEntry{}
		entry.MapKeyExpr.ExprSingle = p.parseExpression(LOWEST)
		if !p.expectPeek(token.COLON) {
			return nil
		}

		p.nextToken()
		entry.MapValueExpr.ExprSingle = p.parseExpression(LOWEST)

		expr.Entries = append(expr.Entries, entry)

		if !p.expectPeek(token.COMMA) {
			break
		}
	}

	if !p.expectPeek(token.RBRACE) {
		return nil
	}

	return expr
}

func (p *Parser) parseArrowExpr(left ast.ExprSingle) ast.ExprSingle {
	expr := &ast.ArrowExpr{ExprSingle: left}

	for {
		p.nextToken()

		b := ast.ArrowBinding{}
		b.ArrowFunctionSpecifier = p.parseArrowFunctionSpecifier()

		p.nextToken()

		b.ArgumentList = p.parseArgumentList()
		expr.Bindings = append(expr.Bindings, b)

		if !p.expectPeek(token.ARROW) {
			break
		}
	}

	return expr
}

func (p *Parser) parseBangExpr(left ast.ExprSingle) ast.ExprSingle {
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

func (p *Parser) parseIfExpr() ast.ExprSingle {
	expr := &ast.IfExpr{}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	expr.TestExpr = p.parseExpression(LOWEST)

	if !p.expectPeek(token.THEN) {
		return nil
	}
	p.nextToken()

	expr.ThenExpr = p.parseExpression(LOWEST)

	if !p.expectPeek(token.ELSE) {
		return nil
	}
	p.nextToken()

	expr.ElseExpr = p.parseExpression(LOWEST)

	return expr
}

func (p *Parser) parseForExpr() ast.ExprSingle {
	expr := &ast.ForExpr{}

	for {
		if !p.expectPeek(token.DOLLAR) {
			return nil
		}
		p.nextToken()

		binding := ast.SimpleForBinding{}
		binding.VarName = p.parseEQName()

		if !p.expectPeek(token.IN) {
			return nil
		}
		p.nextToken()

		binding.ExprSingle = p.parseExpression(LOWEST)
		expr.Bindings = append(expr.Bindings, binding)

		if !p.expectPeek(token.COMMA) {
			break
		}
	}

	if !p.expectPeek(token.RETURN) {
		return nil
	}
	p.nextToken()

	expr.ExprSingle = p.parseExpression(LOWEST)

	return expr
}

func (p *Parser) parseLetExpr() ast.ExprSingle {
	expr := &ast.LetExpr{}

	for {
		if !p.expectPeek(token.DOLLAR) {
			return nil
		}
		p.nextToken()

		binding := ast.SimpleLetBinding{}
		binding.VarName = p.parseEQName()

		if !p.expectPeek(token.ASSIGN) {
			return nil
		}
		p.nextToken()

		binding.ExprSingle = p.parseExpression(LOWEST)
		expr.Bindings = append(expr.Bindings, binding)

		if !p.expectPeek(token.COMMA) {
			break
		}
	}

	if !p.expectPeek(token.RETURN) {
		return nil
	}
	p.nextToken()

	expr.ExprSingle = p.parseExpression(LOWEST)

	return expr
}

func (p *Parser) parseQuantifiedExpr() ast.ExprSingle {
	expr := &ast.QuantifiedExpr{}

	for {
		if !p.expectPeek(token.DOLLAR) {
			return nil
		}
		p.nextToken()

		binding := ast.SimpleQBinding{}
		binding.VarName = p.parseEQName()

		if !p.expectPeek(token.IN) {
			return nil
		}
		p.nextToken()

		binding.ExprSingle = p.parseExpression(LOWEST)
		expr.Bindings = append(expr.Bindings, binding)

		if !p.expectPeek(token.COMMA) {
			break
		}
	}

	if !p.curTokenIs(token.SATISFIES) {
		// TODO error
		return nil
	}
	p.nextToken()

	expr.ExprSingle = p.parseExpression(LOWEST)

	return expr
}

func (p *Parser) parseOrExpr(left ast.ExprSingle) ast.ExprSingle {
	expr := &ast.OrExpr{Token: p.curToken, LeftExpr: left}

	precedence := p.curPrecedence()
	p.nextToken()
	expr.RightExpr = p.parseExpression(precedence)

	return expr
}

func (p *Parser) parseAndExpr(left ast.ExprSingle) ast.ExprSingle {
	expr := &ast.AndExpr{Token: p.curToken, LeftExpr: left}

	precedence := p.curPrecedence()
	p.nextToken()
	expr.RightExpr = p.parseExpression(precedence)

	return expr
}

func (p *Parser) parseRangeExpr(left ast.ExprSingle) ast.ExprSingle {
	expr := &ast.RangeExpr{Token: p.curToken}

	precedence := p.curPrecedence()
	p.nextToken()
	expr.RightExpr = p.parseExpression(precedence)

	return expr
}

func (p *Parser) parseUnionExpr(left ast.ExprSingle) ast.ExprSingle {
	expr := &ast.UnionExpr{Token: p.curToken}

	precedence := p.curPrecedence()
	p.nextToken()
	expr.RightExpr = p.parseExpression(precedence)

	return expr
}

func (p *Parser) parseIntersectExceptExpr(left ast.ExprSingle) ast.ExprSingle {
	expr := &ast.IntersectExceptExpr{Token: p.curToken}

	precedence := p.curPrecedence()
	p.nextToken()
	expr.RightExpr = p.parseExpression(precedence)

	return expr
}

func (p *Parser) parseInstanceofExpr(left ast.ExprSingle) ast.ExprSingle {
	expr := &ast.InstanceofExpr{}
	expr.ExprSingle = left

	if !p.expectPeek(token.OF) {
		return nil
	}
	p.nextToken()

	expr.SequenceType = p.parseSequenceType()

	return expr
}

func (p *Parser) parseCastExpr(left ast.ExprSingle) ast.ExprSingle {
	expr := &ast.CastExpr{}
	expr.ExprSingle = left

	if !p.expectPeek(token.AS) {
		return nil
	}
	p.nextToken()

	expr.SingleType = p.parseSingleType()

	return expr
}

func (p *Parser) parseCastableExpr(left ast.ExprSingle) ast.ExprSingle {
	expr := &ast.CastableExpr{}
	expr.ExprSingle = left

	if !p.expectPeek(token.AS) {
		return nil
	}
	p.nextToken()

	expr.SingleType = p.parseSingleType()

	return expr
}

func (p *Parser) parseTreatExpr(left ast.ExprSingle) ast.ExprSingle {
	expr := &ast.TreatExpr{}
	expr.ExprSingle = left

	if !p.expectPeek(token.AS) {
		return nil
	}
	p.nextToken()

	expr.SequenceType = p.parseSequenceType()

	return expr
}
