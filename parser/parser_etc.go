package parser

import (
	"strconv"
	"strings"

	"github.com/zzossig/xpath/ast"
	"github.com/zzossig/xpath/token"
	"github.com/zzossig/xpath/util"
)

func (p *Parser) parseParam() ast.Param {
	pr := ast.Param{}

	p.nextToken()

	pr.EQName = p.parseEQName()

	if p.peekTokenIs(token.AS) {
		p.nextToken()
		pr.TypeDeclaration = p.parseTypeDeclaration()
	}

	return pr
}

func (p *Parser) parseParamList() ast.ParamList {
	pl := ast.ParamList{}

	for {
		pr := p.parseParam()
		pl.Params = append(pl.Params, pr)

		if !p.peekTokenIs(token.COMMA) {
			break
		}
		p.nextToken()
		p.nextToken()
	}

	return pl
}

func (p *Parser) parsePredicate() ast.Predicate {
	pc := ast.Predicate{}
	p.nextToken()

	e := p.parseExpr()
	ex, ok := e.(*ast.Expr)
	if !ok {
		p.newError("cannot parse Predicate")
		return pc
	}
	pc.Exprs = ex.Exprs

	if !p.expectPeek(token.RBRACKET) {
		return pc
	}

	return pc
}

func (p *Parser) parsePredicateList() ast.PredicateList {
	pl := ast.PredicateList{}

	for {
		pc := p.parsePredicate()
		pl.PL = append(pl.PL, pc)

		if !p.peekTokenIs(token.LBRACKET) {
			break
		}
		p.nextToken()
	}

	return pl
}

func (p *Parser) parseArgument() ast.Argument {
	a := ast.Argument{}

	switch p.curToken.Type {
	case token.QUESTION:
		a.ArgumentPlaceholder.Token = p.curToken
		a.TypeID = 2
	default:
		a.ExprSingle = p.parseExprSingle(LOWEST)
		a.TypeID = 1
	}

	return a
}

func (p *Parser) parseArgumentList() ast.ArgumentList {
	al := ast.ArgumentList{}

	if !p.curTokenIs(token.LPAREN) {
		p.newError("error while parsing ArgumentList: expectCur: (, got=%s", p.curToken.Literal)
		return al
	}
	p.nextToken()

	for !p.curTokenIs(token.RPAREN) {
		arg := p.parseArgument()

		al.Args = append(al.Args, arg)
		p.nextToken()
		if p.curTokenIs(token.COMMA) {
			p.nextToken()
		}
	}

	return al
}

func (p *Parser) parseArrowFunctionSpecifier() ast.ArrowFunctionSpecifier {
	afs := ast.ArrowFunctionSpecifier{}

	switch p.curToken.Type {
	case token.DOLLAR:
		p.nextToken()

		vr := &ast.VarRef{}
		vr.VarName = p.parseEQName()

		afs.TypeID = 2
		afs.VarRef = *vr
	case token.LPAREN:
		e := p.parseParenthesizedExpr()
		pe, ok := e.(*ast.ParenthesizedExpr)
		if !ok {
			p.newError("cannot parse ParenthesizedExpr")
			return afs
		}

		afs.TypeID = 3
		afs.ParenthesizedExpr = *pe
	default:
		afs.TypeID = 1
		afs.EQName = p.parseEQName()
	}

	return afs
}

func (p *Parser) parseNodeTest() ast.NodeTest {
	if util.CheckKindTest(p.curToken.Literal) != 0 && p.peekTokenIs(token.LPAREN) {
		return p.parseKindTest()
	}
	return p.parseNameTest()
}

func (p *Parser) parseItemType() ast.NodeTest {
	test := &ast.ItemType{}

	switch util.CheckItemType(p.curToken.Literal) {
	case 1:
		test.NodeTest = p.parseKindTest()
		test.TypeID = 1
	case 2:
		test.NodeTest = p.parseItemTest()
		test.TypeID = 2
	case 3:
		test.NodeTest = p.parseFunctionTest()
		test.TypeID = 3
	case 4:
		test.NodeTest = p.parseMapTest()
		test.TypeID = 4
	case 5:
		test.NodeTest = p.parseArrayTest()
		test.TypeID = 5
	case 6:
		test.NodeTest = p.parseAtomicOrUnionType()
		test.TypeID = 6
	case 7:
		test.NodeTest = p.parseParenthesizedItemType()
		test.TypeID = 7
	default:
		test.TypeID = 0
	}

	return test
}

func (p *Parser) parseKindTest() ast.NodeTest {
	test := &ast.KindTest{}

	switch util.CheckKindTest(p.curToken.Literal) {
	case 1:
		test.NodeTest = p.parseDocumentTest()
		test.TypeID = 1
	case 2:
		test.NodeTest = p.parseElementTest()
		test.TypeID = 2
	case 3:
		test.NodeTest = p.parseAttributeTest()
		test.TypeID = 3
	case 4:
		test.NodeTest = p.parseSchemaElementTest()
		test.TypeID = 4
	case 5:
		test.NodeTest = p.parseSchemaAttributeTest()
		test.TypeID = 5
	case 6:
		test.NodeTest = p.parsePITest()
		test.TypeID = 6
	case 7:
		test.NodeTest = p.parseCommentTest()
		test.TypeID = 7
	case 8:
		test.NodeTest = p.parseTextTest()
		test.TypeID = 8
	case 9:
		test.NodeTest = p.parseNamespaceNodeTest()
		test.TypeID = 9
	case 10:
		test.NodeTest = p.parseAnyKindTest()
		test.TypeID = 10
	default:
		test.TypeID = 0
	}

	return test
}

func (p *Parser) parseNameTest() ast.NodeTest {
	test := &ast.NameTest{}

	switch p.curToken.Type {
	case token.ASTERISK:
		if p.peekTokenIs(token.COLON) {
			p.nextToken()
			test.Wildcard.NCName.SetValue(p.readNCName())
			test.Wildcard.TypeID = 3
			test.TypeID = 2
		} else {
			test.Wildcard.TypeID = 1
			test.TypeID = 2
		}
	default:
		if p.curToken.Literal == "Q" && p.peekTokenIs(token.LBRACE) {
			bracedURI := p.readBracedURI()
			if !p.expectPeek(token.ASTERISK) {
				p.newError("error while parsing Wildcard: expectPeek: *, got=%s", p.peekToken.Literal)
				return nil
			}
			test.Wildcard.BracedURILiteral.SetValue(bracedURI)
			test.Wildcard.TypeID = 4
			test.TypeID = 2
		} else {
			name := p.readNCName()
			if !p.peekTokenIs(token.COLON) {
				test.TypeID = 1
				test.EQName.SetValue(name)
			} else {
				p.nextToken()
				if p.peekTokenIs(token.ASTERISK) {
					p.nextToken()
					test.Wildcard.TypeID = 2
					test.Wildcard.NCName.SetValue(name)
					test.TypeID = 2
				} else {
					p.nextToken()
					localPart := p.readNCName()

					var sb strings.Builder
					sb.WriteString(name)
					sb.WriteString(":")
					sb.WriteString(localPart)

					test.TypeID = 1
					test.EQName.SetValue(sb.String())
				}
			}
		}
	}

	return test
}

func (p *Parser) parseDocumentTest() ast.NodeTest {
	test := &ast.DocumentTest{}

	if !p.expectPeek(token.LPAREN) {
		p.newError("error while parsing DocumentTest: expectPeek: }, got=%s", p.peekToken.Literal)
		return nil
	}
	p.nextToken()

	if p.curTokenIs(token.RPAREN) {
		return test
	}

	test.NodeTest = p.parseKindTest()

	t, ok := test.NodeTest.(*ast.KindTest)
	if !ok {
		p.newError("cannot parse DocumentTest")
		return nil
	}
	if t.TypeID != 2 && t.TypeID != 4 {
		p.newError("error while parsing DocumentTest: expected ElementTest, SchemaElementTest")
		return nil
	}

	if !p.expectPeek(token.RPAREN) {
		p.newError("error while parsing DocumentTest: expectPeek: ), got=%s", p.peekToken.Literal)
		return nil
	}

	return test
}

func (p *Parser) parseElementTest() ast.NodeTest {
	test := &ast.ElementTest{}

	if !p.expectPeek(token.LPAREN) {
		p.newError("error while parsing ElementTest: expectPeek: (, got=%s", p.peekToken.Literal)
		return nil
	}

	p.nextToken()
	if p.curTokenIs(token.RPAREN) {
		return test
	}

	test.ElementNameOrWildcard = p.parseElementNameOrWildcard()

	p.nextToken()
	if !p.curTokenIs(token.COMMA) {
		return test
	}
	p.nextToken()

	test.TypeName = p.parseEQName()

	p.nextToken()
	if !p.curTokenIs(token.QUESTION) {
		return test
	}

	test.Token = p.curToken

	if !p.expectPeek(token.RPAREN) {
		p.newError("error while parsing ElementTest: expectPeek: ), got=%s", p.peekToken.Literal)
		return nil
	}

	return test
}

func (p *Parser) parseAttributeTest() ast.NodeTest {
	at := &ast.AttributeTest{}

	if !p.expectPeek(token.LPAREN) {
		p.newError("error while parsing AttributeTest: expectPeek: (, got=%s", p.peekToken.Literal)
		return nil
	}

	p.nextToken()
	if p.curTokenIs(token.RPAREN) {
		return at
	}

	at.AttribNameOrWildcard = p.parseAttribNameOrWildcard()

	p.nextToken()
	if !p.curTokenIs(token.COMMA) {
		return at
	}
	p.nextToken()

	at.TypeName = p.parseEQName()

	if !p.expectPeek(token.RPAREN) {
		p.newError("error while parsing AttributeTest: expectPeek: ), got=%s", p.peekToken.Literal)
		return nil
	}

	return at
}

func (p *Parser) parseSchemaElementTest() ast.NodeTest {
	set := &ast.SchemaElementTest{}

	if !p.expectPeek(token.LPAREN) {
		p.newError("error while parsing SchemaElementTest: expectPeek: (, got=%s", p.peekToken.Literal)
		return nil
	}
	p.nextToken()

	set.ElementDeclaration = p.parseEQName()

	if !p.expectPeek(token.RPAREN) {
		p.newError("error while parsing SchemaElementTest: expectPeek: ), got=%s", p.peekToken.Literal)
		return nil
	}

	return set
}

func (p *Parser) parseSchemaAttributeTest() ast.NodeTest {
	sat := &ast.SchemaAttributeTest{}

	if !p.expectPeek(token.LPAREN) {
		p.newError("error while parsing SchemaAttributeTest: expectPeek: (, got=%s", p.peekToken.Literal)
		return nil
	}
	p.nextToken()

	sat.AttributeDeclaration = p.parseEQName()

	if !p.expectPeek(token.RPAREN) {
		p.newError("error while parsing SchemaAttributeTest: expectPeek: ), got=%s", p.peekToken.Literal)
		return nil
	}

	return sat
}

func (p *Parser) parsePITest() ast.NodeTest {
	pit := &ast.PITest{}

	if !p.expectPeek(token.LPAREN) {
		p.newError("error while parsing ProcessingInstructionTest: expectPeek: (, got=%s", p.peekToken.Literal)
		return nil
	}

	p.nextToken()
	if p.curTokenIs(token.RPAREN) {
		return pit
	}

	if p.curTokenIs(token.STRING) {
		pit.StringLiteral = ast.StringLiteral{Value: p.curToken.Literal}
	} else {
		name := ast.NCName{}
		name.SetValue(p.curToken.Literal)
		pit.NCName = name
	}

	if !p.expectPeek(token.RPAREN) {
		p.newError("error while parsing ProcessingInstructionTest: expectPeek: (, got=%s", p.peekToken.Literal)
		return nil
	}

	return pit
}

func (p *Parser) parseCommentTest() ast.NodeTest {
	ct := &ast.CommentTest{}

	if !p.expectPeek(token.LPAREN) {
		p.newError("error while parsing CommentTest: expectPeek: (, got=%s", p.peekToken.Literal)
		return nil
	}

	if !p.expectPeek(token.RPAREN) {
		p.newError("error while parsing CommentTest: expectPeek: ), got=%s", p.peekToken.Literal)
		return nil
	}

	return ct
}

func (p *Parser) parseTextTest() ast.NodeTest {
	tt := &ast.TextTest{}

	if !p.expectPeek(token.LPAREN) {
		p.newError("error while parsing TextTest: expectPeek: (, got=%s", p.peekToken.Literal)
		return nil
	}

	if !p.expectPeek(token.RPAREN) {
		p.newError("error while parsing TextTest: expectPeek: ), got=%s", p.peekToken.Literal)
		return nil
	}

	return tt
}

func (p *Parser) parseNamespaceNodeTest() ast.NodeTest {
	nnt := &ast.NamespaceNodeTest{}

	if !p.expectPeek(token.LPAREN) {
		p.newError("error while parsing NamespaceNodeTest: expectPeek: (, got=%s", p.peekToken.Literal)
		return nil
	}

	if !p.expectPeek(token.RPAREN) {
		p.newError("error while parsing NamespaceNodeTest: expectPeek: ), got=%s", p.peekToken.Literal)
		return nil
	}

	return nnt
}

func (p *Parser) parseAnyKindTest() ast.NodeTest {
	akt := &ast.AnyKindTest{}

	if !p.expectPeek(token.LPAREN) {
		p.newError("error while parsing AnyKindTest: expectPeek: (, got=%s", p.peekToken.Literal)
		return nil
	}

	if !p.expectPeek(token.RPAREN) {
		p.newError("error while parsing AnyKindTest: expectPeek: ), got=%s", p.peekToken.Literal)
		return nil
	}

	return akt
}

func (p *Parser) parseItemTest() ast.NodeTest {
	it := &ast.ItemTest{}

	if !p.expectPeek(token.LPAREN) {
		p.newError("error while parsing ItemTest: expectPeek: (, got=%s", p.peekToken.Literal)
		return nil
	}

	if !p.expectPeek(token.RPAREN) {
		p.newError("error while parsing ItemTest: expectPeek: ), got=%s", p.peekToken.Literal)
		return nil
	}

	return it
}

func (p *Parser) parseFunctionTest() ast.NodeTest {
	ft := &ast.FunctionTest{}

	if !p.expectPeek(token.LPAREN) {
		p.newError("error while parsing FunctionTest: expectPeek: (, got=%s", p.peekToken.Literal)
		return nil
	}

	if p.peekTokenIs(token.ASTERISK) {
		p.nextToken()
		ft.NodeTest = &ast.AnyFunctionTest{}

		if !p.expectPeek(token.RPAREN) {
			p.newError("error while parsing AnyFunctionTest: expectPeek: ), got=%s", p.peekToken.Literal)
			return nil
		}
	} else {
		tft := &ast.TypedFunctionTest{}

		if !p.peekTokenIs(token.RPAREN) {
			p.nextToken()

			for {
				st := p.parseSequenceType()
				tft.ParamSTypes = append(tft.ParamSTypes, st)

				if !p.peekTokenIs(token.COMMA) {
					break
				}
				p.nextToken()
				p.nextToken()
			}

			if !p.expectPeek(token.RPAREN) {
				p.newError("error while parsing TypedFunctionTest: expectPeek: ), got=%s", p.peekToken.Literal)
				return nil
			}
		}

		if !p.expectPeek(token.AS) {
			p.newError("error while parsing TypedFunctionTest: expectPeek: as, got=%s", p.peekToken.Literal)
			return nil
		}
		p.nextToken()

		tft.AsSType = p.parseSequenceType()
		ft.NodeTest = tft
	}

	return ft
}

func (p *Parser) parseMapTest() ast.NodeTest {
	mt := &ast.MapTest{}

	if !p.expectPeek(token.LPAREN) {
		p.newError("error while parsing MapTest: expectPeek: (, got=%s", p.peekToken.Literal)
		return nil
	}

	if p.peekTokenIs(token.ASTERISK) {
		p.nextToken()
		mt.NodeTest = &ast.AnyMapTest{}
	} else {
		p.nextToken()
		tmt := &ast.TypedMapTest{}
		tmt.AtomicOrUnionType.EQName = p.parseEQName()

		if !p.expectPeek(token.COMMA) {
			p.newError("error while parsing TypedMapTest: expectPeek: [,], got=%s", p.peekToken.Literal)
			return nil
		}
		p.nextToken()

		tmt.SequenceType = p.parseSequenceType()
		mt.NodeTest = tmt
	}

	if !p.expectPeek(token.RPAREN) {
		p.newError("error while parsing MapTest: expectPeek: ), got=%s", p.peekToken.Literal)
		return nil
	}

	return mt
}

func (p *Parser) parseArrayTest() ast.NodeTest {
	at := &ast.ArrayTest{}

	if !p.expectPeek(token.LPAREN) {
		p.newError("error while parsing ArrayTest: expectPeek: (, got=%s", p.peekToken.Literal)
		return nil
	}

	if p.peekTokenIs(token.ASTERISK) {
		p.nextToken()
		at.NodeTest = &ast.AnyArrayTest{}
	} else {
		p.nextToken()
		tat := &ast.TypedArrayTest{}
		tat.SequenceType = p.parseSequenceType()

		at.NodeTest = tat
	}

	if !p.expectPeek(token.RPAREN) {
		p.newError("error while parsing ArrayTest: expectPeek: ), got=%s", p.peekToken.Literal)
		return nil
	}

	return at
}

func (p *Parser) parseAtomicOrUnionType() ast.NodeTest {
	aout := &ast.AtomicOrUnionType{}
	name := p.parseEQName()
	aout.EQName.SetValue(name.Value())

	return aout
}

func (p *Parser) parseParenthesizedItemType() ast.NodeTest {
	pit := &ast.ParenthesizedItemType{}

	p.nextToken()
	pit.NodeTest = p.parseItemType()

	if !p.expectPeek(token.RPAREN) {
		p.newError("error while parsing ParenthesizedItemType: expectPeek: ), got=%s", p.peekToken.Literal)
		return nil
	}

	return pit
}

func (p *Parser) parseElementNameOrWildcard() ast.ElementNameOrWildcard {
	enow := ast.ElementNameOrWildcard{}

	if p.curTokenIs(token.ASTERISK) {
		enow.WC = "*"
	} else {
		enow.ElementName = p.parseEQName()
	}

	return enow
}

func (p *Parser) parseAttribNameOrWildcard() ast.AttribNameOrWildcard {
	anow := ast.AttribNameOrWildcard{}

	if p.curTokenIs(token.ASTERISK) {
		anow.WC = "*"
	} else {
		anow.AttributeName = p.parseEQName()
	}

	return anow
}

// parseEQName can returns BracedURILiteral which is not an EQName
func (p *Parser) parseEQName() ast.EQName {
	eqn := ast.EQName{}
	eqn.SetValue(p.readEQName())

	return eqn
}

func (p *Parser) parseSequenceType() ast.SequenceType {
	st := ast.SequenceType{}

	if p.curToken.Literal == "empty-sequence" {
		if !p.expectPeek(token.LPAREN) {
			p.newError("error while parsing empty-sequence: expectPeek: (, got=%s", p.peekToken.Literal)
		}

		if !p.expectPeek(token.RPAREN) {
			p.newError("error while parsing empty-sequence: expectPeek: ), got=%s", p.peekToken.Literal)
		}
		st.TypeID = 1
	} else {
		st.TypeID = 2
		st.NodeTest = p.parseItemType()

		if p.peekTokenIs(token.QUESTION, token.ASTERISK, token.PLUS) {
			p.nextToken()
			st.OccurrenceIndicator = ast.OccurrenceIndicator{Token: p.curToken}
		}
	}

	return st
}

func (p *Parser) parseSingleType() ast.SingleType {
	st := ast.SingleType{}
	st.SimpleTypeName = p.parseEQName()

	if p.peekTokenIs(token.QUESTION) {
		p.nextToken()
		st.Token = p.curToken
	}

	return st
}

func (p *Parser) parseTypeDeclaration() ast.TypeDeclaration {
	p.nextToken()

	td := ast.TypeDeclaration{}
	td.SequenceType = p.parseSequenceType()

	return td
}

func (p *Parser) parseEnclosedExpr() ast.EnclosedExpr {
	ee := ast.EnclosedExpr{}

	if !p.peekTokenIs(token.RBRACE) {
		p.nextToken()

		e := p.parseExpr()
		er, ok := e.(*ast.Expr)
		if !ok {
			p.newError("cannot parse EnclosedExpr")
			return ee
		}
		ee.Exprs = er.Exprs
	}

	if !p.expectPeek(token.RBRACE) {
		p.newError("error while parsing EnclosedExpr: expectPeek: }, got=%s", p.peekToken.Literal)
		return ee
	}

	return ee
}

func (p *Parser) parsePal() ast.PAL {
	if p.curTokenIs(token.LBRACKET) {
		p.nextToken()
		pal := &ast.Predicate{}

		e := p.parseExpr()
		er, ok := e.(*ast.Expr)
		if !ok {
			return nil
		}

		if !p.expectPeek(token.RBRACKET) {
			p.newError("error while parsing Predicate: expectPeek: ], got=%s", p.peekToken.Literal)
			return nil
		}

		pal.Exprs = er.Exprs
		return pal
	} else if p.curTokenIs(token.LPAREN) {
		pal := &ast.ArgumentList{}
		pal.Args = p.parseArgumentList().Args
		return pal
	} else if p.curTokenIs(token.QUESTION) {
		pal := &ast.Lookup{Token: p.curToken}
		p.nextToken()
		pal.KeySpecifier = p.parseKeySpecifier()
		return pal
	}

	p.newError("cannot parse PAL expression. expectCur: [, {, (. got=%s", p.curToken.Literal)
	return nil
}

func (p *Parser) parseKeySpecifier() ast.KeySpecifier {
	ks := ast.KeySpecifier{}

	switch p.curToken.Type {
	case token.ASTERISK:
		ks.TypeID = 4
	case token.INT:
		ks.TypeID = 2
		i, _ := strconv.Atoi(p.curToken.Literal)
		ks.IntegerLiteral.Value = i
	case token.LPAREN:
		ks.TypeID = 3
		pe := p.parseParenthesizedExpr()
		pep, ok := pe.(*ast.ParenthesizedExpr)
		if !ok {
			p.newError("cannot parse ParenthesizedExpr")
			return ks
		}
		ks.ParenthesizedExpr.Expr = pep.Expr
	default:
		ks.TypeID = 1
		ks.NCName.SetValue(p.curToken.Literal)
	}

	return ks
}
