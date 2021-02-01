package parser

import (
	"strconv"

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

	if util.IsForwardAxis(name.Value()) {
		as := &ast.AxisStep{}

		as.TypeID = 2
		as.ForwardStep.TypeID = 1
		as.ForwardAxis.SetValue(name.Value())

		p.nextToken()
		as.ForwardStep.NodeTest = p.parseNodeTest()

		p.nextToken()
		as.PredicateList = p.parsePredicateList()
		return as
	}

	if util.IsReverseAxis(name.Value()) {
		as := &ast.AxisStep{}

		as.TypeID = 1
		as.ReverseStep.TypeID = 1
		as.ReverseAxis.SetValue(name.Value())

		p.nextToken()
		as.ReverseStep.NodeTest = p.parseNodeTest()

		p.nextToken()
		as.PredicateList = p.parsePredicateList()
		return as
	}

	i := &ast.Identifier{}
	i.EQName = name

	return i
}

func (p *Parser) parseIntegerLiteral() ast.ExprSingle {
	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	il := &ast.IntegerLiteral{Value: int(value)}

	if err != nil {
		// TODO error
		return nil
	}

	if p.peekTokenIs(token.LBRACKET, token.LPAREN, token.QUESTION) {
		return p.parsePostfixExpr(il)
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

	if p.peekTokenIs(token.LBRACKET, token.LPAREN, token.QUESTION) {
		return p.parsePostfixExpr(dl)
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

	if p.peekTokenIs(token.LBRACKET, token.LPAREN, token.QUESTION) {
		return p.parsePostfixExpr(dl)
	}

	return dl
}

func (p *Parser) parseStringLiteral() ast.ExprSingle {
	sl := &ast.StringLiteral{Value: p.curToken.Literal}
	if p.peekTokenIs(token.LBRACKET, token.LPAREN, token.QUESTION) {
		return p.parsePostfixExpr(sl)
	}
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

		return fc
	}

	return vr
}

func (p *Parser) parseGroupedExpr() ast.ExprSingle {
	if p.hasComma() {
		return p.parseSequenceExpr()
	}

	p.nextToken()

	expr := p.parseExprSingle(LOWEST)
	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	return expr
}

func (p *Parser) parseSequenceExpr() ast.ExprSingle {
	p.nextToken()
	expr := &ast.Expr{}

	for !p.curTokenIs(token.RPAREN) {
		e := p.parseExprSingle(LOWEST)
		expr.Exprs = append(expr.Exprs, e)

		p.nextToken()
		if p.curTokenIs(token.COMMA) {
			p.nextToken()
		}
	}

	if p.peekTokenIs(token.LBRACKET, token.LPAREN, token.QUESTION) {
		return p.parsePostfixExpr(expr)
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
	right := p.parseExprSingle(precedence)

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
	right := p.parseExprSingle(precedence)

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

	expr.ExprSingle = p.parseExprSingle(UNARY)

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

		if !p.expectPeek(token.COMMA) {
			break
		}
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

func (p *Parser) parseBangExpr(left ast.ExprSingle) ast.ExprSingle {
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
	expr := &ast.ComparisonExpr{LeftExpr: left}
	expr.SetToken(p.curToken)

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

	expr.TestExpr = p.parseExprSingle(LOWEST)

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

		binding.ExprSingle = p.parseExprSingle(LOWEST)
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

	expr.ExprSingle = p.parseExprSingle(LOWEST)

	return expr
}

func (p *Parser) parseOrExpr(left ast.ExprSingle) ast.ExprSingle {
	expr := &ast.OrExpr{Token: p.curToken, LeftExpr: left}

	precedence := p.curPrecedence()
	p.nextToken()
	expr.RightExpr = p.parseExprSingle(precedence)

	return expr
}

func (p *Parser) parseAndExpr(left ast.ExprSingle) ast.ExprSingle {
	expr := &ast.AndExpr{Token: p.curToken, LeftExpr: left}

	precedence := p.curPrecedence()
	p.nextToken()
	expr.RightExpr = p.parseExprSingle(precedence)

	return expr
}

func (p *Parser) parseRangeExpr(left ast.ExprSingle) ast.ExprSingle {
	expr := &ast.RangeExpr{Token: p.curToken}

	precedence := p.curPrecedence()
	p.nextToken()
	expr.RightExpr = p.parseExprSingle(precedence)

	if p.peekTokenIs(token.LBRACKET, token.LPAREN, token.QUESTION) {
		return p.parsePostfixExpr(expr)
	}

	return expr
}

func (p *Parser) parseUnionExpr(left ast.ExprSingle) ast.ExprSingle {
	expr := &ast.UnionExpr{Token: p.curToken}

	precedence := p.curPrecedence()
	p.nextToken()
	expr.RightExpr = p.parseExprSingle(precedence)

	return expr
}

func (p *Parser) parseIntersectExceptExpr(left ast.ExprSingle) ast.ExprSingle {
	expr := &ast.IntersectExceptExpr{Token: p.curToken}

	precedence := p.curPrecedence()
	p.nextToken()
	expr.RightExpr = p.parseExprSingle(precedence)

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

func (p *Parser) parseInlineFunctionExpr() ast.ExprSingle {
	expr := &ast.InlineFunctionExpr{}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	if !p.expectPeek(token.RPAREN) {
		expr.ParamList = p.parseParamList()
		p.nextToken()
	}

	if p.expectPeek(token.AS) {
		p.nextToken()
		expr.SequenceType = p.parseSequenceType()
		p.nextToken()
	}

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
	pe := &ast.PostfixExpr{}

	switch left.(type) {
	case *ast.IntegerLiteral:
		e := left.(*ast.IntegerLiteral)
		pe.PrimaryExpr = e
	case *ast.DecimalLiteral:
		e := left.(*ast.DecimalLiteral)
		pe.PrimaryExpr = e
	case *ast.DoubleLiteral:
		e := left.(*ast.DoubleLiteral)
		pe.PrimaryExpr = e
	case *ast.StringLiteral:
		e := left.(*ast.StringLiteral)
		pe.PrimaryExpr = e
	case *ast.VarRef:
		e := left.(*ast.VarRef)
		pe.PrimaryExpr = e
	case *ast.ParenthesizedExpr:
		e := left.(*ast.ParenthesizedExpr)
		pe.PrimaryExpr = e
	case *ast.ContextItemExpr:
		e := left.(*ast.ContextItemExpr)
		pe.PrimaryExpr = e
	case *ast.FunctionCall:
		e := left.(*ast.FunctionCall)
		pe.PrimaryExpr = e
	case *ast.NamedFunctionRef:
		e := left.(*ast.NamedFunctionRef)
		pe.PrimaryExpr = e
	case *ast.InlineFunctionExpr:
		e := left.(*ast.InlineFunctionExpr)
		pe.PrimaryExpr = e
	case *ast.MapConstructor:
		e := left.(*ast.MapConstructor)
		pe.PrimaryExpr = e
	case *ast.SquareArrayConstructor:
		e := left.(*ast.SquareArrayConstructor)
		pe.PrimaryExpr = e
	case *ast.CurlyArrayConstructor:
		e := left.(*ast.CurlyArrayConstructor)
		pe.PrimaryExpr = e
	case *ast.UnaryLookup:
		e := left.(*ast.UnaryLookup)
		pe.PrimaryExpr = e
	default:
		// panic
		return nil
	}

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

	return pe
}

func (p *Parser) parsePathExpr() ast.ExprSingle {
	expr := &ast.PathExpr{Token: p.curToken}

	p.nextToken()
	expr.ExprSingle = p.parseRelativePathExpr()

	return expr
}

func (p *Parser) parseRelativePathExpr() ast.ExprSingle {
	expr := &ast.RelativePathExpr{}

	for {
		t := p.curToken
		e := p.parseExprSingle(LOWEST)

		expr.Exprs = append(expr.Exprs, e)
		expr.Tokens = append(expr.Tokens, t)

		if !p.peekTokenIs(token.SLASH, token.DSLASH) {
			break
		}
		p.nextToken()
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

	p.nextToken()
	expr.PredicateList = p.parsePredicateList()

	return expr
}
