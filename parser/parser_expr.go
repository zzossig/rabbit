package parser

import (
	"strconv"
	"strings"

	"github.com/zzossig/xpath/ast"
	"github.com/zzossig/xpath/token"
	"github.com/zzossig/xpath/util"
)

func (p *Parser) parseExpr() ast.ExprSingle {
	expr := &ast.Expr{}

	for {
		e := p.parseExprSingle(LOWEST)
		expr.Exprs = append(expr.Exprs, e)

		if !p.peekTokenIs(token.COMMA) {
			break
		}
		p.nextToken()
		p.nextToken()
	}

	return expr
}

func (p *Parser) parseExprSingle(precedence int) ast.ExprSingle {
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

func (p *Parser) parseIdentifier() ast.ExprSingle {
	name := p.parseEQName()

	if p.peekTokenIs(token.LPAREN) {
		p.nextToken()

		fc := &ast.FunctionCall{}
		fc.ArgumentList = p.parseArgumentList()
		fc.EQName = name

		return fc
	}

	if p.peekTokenIs(token.HASH) {
		p.nextToken()

		i := &ast.Identifier{}
		i.EQName = name

		// TODO should check name is reserved-function0name

		return p.parseNamedFunctionRef(i)
	}

	if p.peekTokenIs(token.DCOLON) {
		p.nextToken()

		var sb strings.Builder
		sb.WriteString(name.Value())
		sb.WriteString(p.curToken.Literal)
		axis := sb.String()

		if util.IsForwardAxis(axis) {
			as := &ast.AxisStep{}
			as.TypeID = 2
			as.ForwardStep.TypeID = 1
			as.ForwardAxis.SetValue(axis)

			p.nextToken()
			as.ForwardStep.NodeTest = p.parseNodeTest()

			if p.peekTokenIs(token.LBRACKET) {
				p.nextToken()
				as.PredicateList = p.parsePredicateList()
			}

			return as
		}

		if util.IsReverseAxis(axis) {
			as := &ast.AxisStep{}

			as.TypeID = 1
			as.ReverseStep.TypeID = 1
			as.ReverseAxis.SetValue(axis)

			p.nextToken()
			as.ReverseStep.NodeTest = p.parseNodeTest()

			if p.peekTokenIs(token.LBRACKET) {
				p.nextToken()
				as.PredicateList = p.parsePredicateList()
			}

			return as
		}
	}

	i := &ast.Identifier{}
	i.EQName = name

	if p.peekTokenIs(token.LBRACKET, token.LPAREN, token.QUESTION) {
		return p.parsePostfixExpr(i)
	}

	return i
}

func (p *Parser) parseIntegerLiteral() ast.ExprSingle {
	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	il := &ast.IntegerLiteral{Value: int(value)}

	if err != nil {
		// TODO error
		return nil
	}

	return il
}

func (p *Parser) parseDecimalLiteral() ast.ExprSingle {
	value, err := strconv.ParseFloat(p.curToken.Literal, 64)
	dl := &ast.DecimalLiteral{Value: value}

	if err != nil {
		// TODO error
		return nil
	}

	return dl
}

func (p *Parser) parseDoubleLiteral() ast.ExprSingle {
	value, err := strconv.ParseFloat(p.curToken.Literal, 64)
	dl := &ast.DoubleLiteral{Value: value}

	if err != nil {
		// TODO error
		return nil
	}

	return dl
}

func (p *Parser) parseStringLiteral() ast.ExprSingle {
	sl := &ast.StringLiteral{Value: p.curToken.Literal}
	return sl
}

func (p *Parser) parseVariable() ast.ExprSingle {
	vr := &ast.VarRef{}

	p.nextToken()

	name := p.parseEQName()
	vr.VarName = name

	if p.peekTokenIs(token.LBRACKET, token.QUESTION) {
		return p.parsePostfixExpr(vr)
	}

	if p.peekTokenIs(token.LPAREN) {
		p.nextToken()

		fc := &ast.FunctionCall{}
		fc.ArgumentList = p.parseArgumentList()
		fc.EQName = name

		if p.peekTokenIs(token.LBRACKET, token.LPAREN, token.QUESTION) {
			return p.parsePostfixExpr(fc)
		}

		return fc
	}

	return vr
}

func (p *Parser) parseGroupedExpr() ast.ExprSingle {
	if p.hasComma() {
		return p.parseSequenceExpr()
	}

	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return &ast.ParenthesizedExpr{}
	}

	p.nextToken()
	expr := p.parseExprSingle(LOWEST)
	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	if p.peekTokenIs(token.LBRACKET, token.LPAREN, token.QUESTION) {
		return p.parsePostfixExpr(expr)
	}

	return expr
}

func (p *Parser) parseSequenceExpr() ast.ExprSingle {
	expr := &ast.ParenthesizedExpr{}

	for {
		p.nextToken()
		e := p.parseExprSingle(LOWEST)
		expr.Exprs = append(expr.Exprs, e)

		if !p.peekTokenIs(token.COMMA) {
			break
		}
		p.nextToken()
	}

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	if p.peekTokenIs(token.LBRACKET, token.LPAREN, token.QUESTION) {
		return p.parsePostfixExpr(expr)
	}

	return expr
}

func (p *Parser) parseAdditiveExpr(left ast.ExprSingle) ast.ExprSingle {
	expr := &ast.AdditiveExpr{LeftExpr: left, Token: p.curToken}

	precedence := p.curPrecedence()
	p.nextToken()
	right := p.parseExprSingle(precedence)

	if right != nil {
		expr.RightExpr = right
	}

	return expr
}

func (p *Parser) parseMultiplicativeExpr(left ast.ExprSingle) ast.ExprSingle {
	expr := &ast.MultiplicativeExpr{LeftExpr: left, Token: p.curToken}

	precedence := p.curPrecedence()
	p.nextToken()
	right := p.parseExprSingle(precedence)

	if right != nil {
		expr.RightExpr = right
	}

	return expr
}

func (p *Parser) parseUnaryExpr() ast.ExprSingle {
	expr := &ast.UnaryExpr{Token: p.curToken}

	precedence := p.precedence(token.UPLUS)
	p.nextToken()
	expr.ExprSingle = p.parseExprSingle(precedence)

	return expr
}

func (p *Parser) parseUnaryLookupExpr() ast.ExprSingle {
	expr := &ast.UnaryLookup{}

	if !p.curTokenIs(token.QUESTION) {
		return nil
	}

	expr.Token = token.Token{Type: token.UQUESTION, Literal: "(?)"}
	p.nextToken()
	expr.KeySpecifier = p.parseKeySpecifier()

	if p.peekTokenIs(token.LBRACKET, token.LPAREN, token.QUESTION) {
		return p.parsePostfixExpr(expr)
	}

	return expr
}

func (p *Parser) parseCurlyArrayExpr() ast.ExprSingle {
	expr := &ast.CurlyArrayConstructor{}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	expr.EnclosedExpr = p.parseEnclosedExpr()
	if p.peekTokenIs(token.LBRACKET, token.LPAREN, token.QUESTION) {
		return p.parsePostfixExpr(expr)
	}

	return expr
}

func (p *Parser) parseSquareArrayExpr() ast.ExprSingle {
	expr := &ast.SquareArrayConstructor{}

	p.nextToken()

	for !p.curTokenIs(token.RBRACKET) {
		e := p.parseExprSingle(LOWEST)
		if e != nil {
			expr.Exprs = append(expr.Exprs, e)
		}
		p.nextToken()
	}

	if p.peekTokenIs(token.LBRACKET, token.LPAREN, token.QUESTION) {
		return p.parsePostfixExpr(expr)
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

		if p.peekTokenIs(token.LBRACKET, token.LPAREN, token.QUESTION) {
			return p.parsePostfixExpr(expr)
		}

		return expr
	}

	for {
		p.nextToken()

		entry := ast.MapConstructorEntry{}
		entry.MapKeyExpr.ExprSingle = p.parseExprSingle(LOWEST)
		if !p.expectPeek(token.COLON) {
			return nil
		}

		p.nextToken()
		entry.MapValueExpr.ExprSingle = p.parseExprSingle(LOWEST)

		expr.Entries = append(expr.Entries, entry)

		if !p.peekTokenIs(token.COMMA) {
			break
		}
		p.nextToken()
	}

	if !p.expectPeek(token.RBRACE) {
		return nil
	}

	if p.peekTokenIs(token.LBRACKET, token.LPAREN, token.QUESTION) {
		return p.parsePostfixExpr(expr)
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

func (p *Parser) parseSimpleMapExpr(left ast.ExprSingle) ast.ExprSingle {
	expr := &ast.SimpleMapExpr{LeftExpr: left, Token: p.curToken}

	precedence := p.curPrecedence()
	p.nextToken()
	right := p.parseExprSingle(precedence)

	if right != nil {
		expr.RightExpr = right
	}

	return expr
}

func (p *Parser) parseComparisonExpr(left ast.ExprSingle) ast.ExprSingle {
	expr := &ast.ComparisonExpr{LeftExpr: left, Token: p.curToken}

	precedence := p.curPrecedence()
	p.nextToken()
	right := p.parseExprSingle(precedence)

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
	p.nextToken()

	expr.TestExpr = p.parseExpr()

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	if !p.expectPeek(token.THEN) {
		return nil
	}
	p.nextToken()

	expr.ThenExpr = p.parseExprSingle(LOWEST)

	if !p.expectPeek(token.ELSE) {
		return nil
	}
	p.nextToken()

	expr.ElseExpr = p.parseExprSingle(LOWEST)

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

		binding.ExprSingle = p.parseExprSingle(LOWEST)
		expr.Bindings = append(expr.Bindings, binding)

		if !p.peekTokenIs(token.COMMA) {
			break
		}
		p.nextToken()
	}

	if !p.expectPeek(token.RETURN) {
		return nil
	}
	p.nextToken()

	expr.ExprSingle = p.parseExprSingle(LOWEST)

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

		binding.ExprSingle = p.parseExprSingle(LOWEST)
		expr.Bindings = append(expr.Bindings, binding)

		if !p.expectPeek(token.COMMA) {
			break
		}
	}

	if !p.expectPeek(token.RETURN) {
		return nil
	}
	p.nextToken()

	expr.ExprSingle = p.parseExprSingle(LOWEST)

	return expr
}

func (p *Parser) parseQuantifiedExpr() ast.ExprSingle {
	expr := &ast.QuantifiedExpr{Token: p.curToken}

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

		binding.ExprSingle = p.parseExprSingle(LOWEST)
		expr.Bindings = append(expr.Bindings, binding)

		if !p.peekTokenIs(token.COMMA) {
			break
		}
		p.nextToken()
	}

	if !p.expectPeek(token.SATISFIES) {
		// TODO error
		return nil
	}
	p.nextToken()

	expr.ExprSingle = p.parseExprSingle(LOWEST)

	return expr
}

func (p *Parser) parseOrExpr(left ast.ExprSingle) ast.ExprSingle {
	expr := &ast.OrExpr{LeftExpr: left, Token: p.curToken}

	precedence := p.curPrecedence()
	p.nextToken()
	expr.RightExpr = p.parseExprSingle(precedence)

	return expr
}

func (p *Parser) parseAndExpr(left ast.ExprSingle) ast.ExprSingle {
	expr := &ast.AndExpr{LeftExpr: left, Token: p.curToken}

	precedence := p.curPrecedence()
	p.nextToken()
	expr.RightExpr = p.parseExprSingle(precedence)

	return expr
}

func (p *Parser) parseRangeExpr(left ast.ExprSingle) ast.ExprSingle {
	expr := &ast.RangeExpr{LeftExpr: left, Token: p.curToken}

	precedence := p.curPrecedence()
	p.nextToken()
	expr.RightExpr = p.parseExprSingle(precedence)

	if p.peekTokenIs(token.LBRACKET, token.LPAREN, token.QUESTION) {
		return p.parsePostfixExpr(expr)
	}

	return expr
}

func (p *Parser) parseUnionExpr(left ast.ExprSingle) ast.ExprSingle {
	expr := &ast.UnionExpr{LeftExpr: left, Token: p.curToken}

	precedence := p.curPrecedence()
	p.nextToken()
	expr.RightExpr = p.parseExprSingle(precedence)

	return expr
}

func (p *Parser) parseStringConcatExpr(left ast.ExprSingle) ast.ExprSingle {
	expr := &ast.StringConcatExpr{LeftExpr: left, Token: p.curToken}

	precedence := p.curPrecedence()
	p.nextToken()
	expr.RightExpr = p.parseExprSingle(precedence)

	return expr
}

func (p *Parser) parseIntersectExceptExpr(left ast.ExprSingle) ast.ExprSingle {
	expr := &ast.IntersectExceptExpr{LeftExpr: left, Token: p.curToken}

	precedence := p.curPrecedence()
	p.nextToken()
	expr.RightExpr = p.parseExprSingle(precedence)

	return expr
}

func (p *Parser) parseInstanceofExpr(left ast.ExprSingle) ast.ExprSingle {
	expr := &ast.InstanceofExpr{ExprSingle: left}

	if !p.expectPeek(token.OF) {
		return nil
	}
	p.nextToken()

	expr.SequenceType = p.parseSequenceType()

	return expr
}

func (p *Parser) parseCastExpr(left ast.ExprSingle) ast.ExprSingle {
	expr := &ast.CastExpr{ExprSingle: left}

	if !p.expectPeek(token.AS) {
		return nil
	}
	p.nextToken()

	expr.SingleType = p.parseSingleType()

	return expr
}

func (p *Parser) parseCastableExpr(left ast.ExprSingle) ast.ExprSingle {
	expr := &ast.CastableExpr{ExprSingle: left}

	if !p.expectPeek(token.AS) {
		return nil
	}
	p.nextToken()

	expr.SingleType = p.parseSingleType()

	return expr
}

func (p *Parser) parseTreatExpr(left ast.ExprSingle) ast.ExprSingle {
	expr := &ast.TreatExpr{ExprSingle: left}

	if !p.expectPeek(token.AS) {
		return nil
	}
	p.nextToken()

	expr.SequenceType = p.parseSequenceType()

	return expr
}

func (p *Parser) parseInlineFunctionExpr() ast.ExprSingle {
	expr := &ast.InlineFunctionExpr{}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	if !p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		expr.ParamList = p.parseParamList()
	}

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	if p.peekTokenIs(token.AS) {
		p.nextToken()
		p.nextToken()
		expr.SequenceType = p.parseSequenceType()
	}

	p.nextToken()
	expr.FunctionBody = p.parseEnclosedExpr()

	if p.peekTokenIs(token.LBRACKET, token.LPAREN, token.QUESTION) {
		return p.parsePostfixExpr(expr)
	}

	return expr
}

func (p *Parser) parseNamedFunctionRef(left ast.ExprSingle) ast.ExprSingle {
	expr := &ast.NamedFunctionRef{}
	ident, ok := left.(*ast.Identifier)
	if !ok {
		return nil
	}

	expr.EQName = ident.EQName

	p.nextToken()
	il := p.parseIntegerLiteral()
	i, ok := il.(*ast.IntegerLiteral)
	if !ok {
		return nil
	}

	expr.IntegerLiteral.Value = i.Value

	if p.peekTokenIs(token.LBRACKET, token.LPAREN, token.QUESTION) {
		return p.parsePostfixExpr(expr)
	}

	return expr
}

func (p *Parser) parseContextItemExpr() ast.ExprSingle {
	cie := &ast.ContextItemExpr{Token: p.curToken}

	if p.peekTokenIs(token.LBRACKET, token.LPAREN, token.QUESTION) {
		return p.parsePostfixExpr(cie)
	}

	return cie
}

func (p *Parser) parsePostfixExpr(left ast.ExprSingle) ast.ExprSingle {
	pe := &ast.PostfixExpr{ExprSingle: left}

	for p.peekTokenIs(token.LBRACKET, token.LPAREN, token.QUESTION) {
		p.nextToken()

		pal := p.parsePal()
		pe.Pals = append(pe.Pals, pal)
	}

	return pe
}

func (p *Parser) parseParenthesizedExpr() ast.ExprSingle {
	pe := &ast.ParenthesizedExpr{}

	if !p.curTokenIs(token.LPAREN) {
		return nil
	}
	p.nextToken()

	e := p.parseExpr()
	er, ok := e.(*ast.Expr)
	if !ok {
		return nil
	}
	pe.Exprs = er.Exprs

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	if p.peekTokenIs(token.LBRACKET, token.LPAREN, token.QUESTION) {
		return p.parsePostfixExpr(pe)
	}

	return pe
}

func (p *Parser) parsePathExpr() ast.ExprSingle {
	expr := &ast.PathExpr{Token: p.curToken}

	precedence := p.curPrecedence()
	p.nextToken()
	expr.ExprSingle = p.parseExprSingle(precedence)

	if p.peekTokenIs(token.LBRACKET, token.LPAREN, token.QUESTION) {
		return p.parsePostfixExpr(expr)
	}

	return expr
}

func (p *Parser) parseRelativePathExpr(left ast.ExprSingle) ast.ExprSingle {
	expr := &ast.RelativePathExpr{LeftExpr: left, Token: p.curToken}

	precedence := p.curPrecedence()
	p.nextToken()
	right := p.parseExprSingle(precedence)

	if right != nil {
		expr.RightExpr = right
	}

	return expr
}

func (p *Parser) parseAxisStep() ast.ExprSingle {
	expr := &ast.AxisStep{}

	if p.curTokenIs(token.DDOT) {
		expr.TypeID = 1
		expr.ReverseStep.TypeID = 2
		expr.ReverseStep.AbbrevReverseStep.Token = p.curToken
	}
	if p.curTokenIs(token.AT) {
		expr.TypeID = 2
		expr.ForwardStep.TypeID = 2
		expr.ForwardStep.AbbrevForwardStep.Token = p.curToken

		p.nextToken()
		expr.ForwardStep.AbbrevForwardStep.NodeTest = p.parseNodeTest()
	}

	if p.peekTokenIs(token.LBRACKET) {
		p.nextToken()
		expr.PredicateList = p.parsePredicateList()
	}

	return expr
}

func (p *Parser) parseWildcard() ast.ExprSingle {
	expr := &ast.Wildcard{}

	if p.peekTokenIs(token.COLON) {
		p.nextToken()
		expr.NCName.SetValue(p.readNCName())
		expr.TypeID = 3
	} else {
		expr.TypeID = 1
	}

	return expr
}
