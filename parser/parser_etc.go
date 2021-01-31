package parser

import (
	"strings"

	"github.com/zzossig/xpath/ast"
	"github.com/zzossig/xpath/token"
	"github.com/zzossig/xpath/util"
)

// Argument
// TypeDeclaration
// Param
// ParamList

func (p *Parser) parseArgumentList() ast.ArgumentList {
	al := ast.ArgumentList{}

	if !p.expectPeek(token.LPAREN) {
		return al
	}
	p.nextToken()

	for !p.curTokenIs(token.RPAREN) {
		arg := ast.Argument{}

		switch p.curToken.Type {
		case token.QUESTION:
			t := token.Token{Type: token.QUESTION, Literal: "?"}
			arg.ArgumentPlaceholder.Token = t
		default:
			arg.ExprSingle = p.parseExpression(LOWEST)
		}

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
		v := p.parseVariable()
		vr, ok := v.(*ast.VarRef)
		if !ok {
			// TODO error
			return afs
		}

		afs.TypeID = 2
		afs.VarRef = *vr
	case token.LPAREN:
		e := p.parseExpression(LOWEST)
		pe, ok := e.(*ast.ParenthesizedExpr)
		if !ok {
			// TODO error
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

func (p *Parser) parseDocumentTest() ast.NodeTest {
	test := &ast.DocumentTest{}

	if !p.expectPeek(token.LPAREN) {
		// TODO error
		return nil
	}
	p.nextToken()

	if p.curTokenIs(token.RPAREN) {
		return test
	}

	test.NodeTest = p.parseKindTest()

	t, ok := test.NodeTest.(*ast.KindTest)
	if !ok {
		// TODO error
		return nil
	}
	if t.TypeID != 2 && t.TypeID != 4 {
		// TODO error
		return nil
	}

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return test
}

func (p *Parser) parseElementTest() ast.NodeTest {
	test := &ast.ElementTest{}

	if !p.expectPeek(token.LPAREN) {
		// TODO error
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
		return nil
	}

	return test
}

func (p *Parser) parseAttributeTest() ast.NodeTest {
	at := &ast.AttributeTest{}

	if !p.expectPeek(token.LPAREN) {
		// TODO error
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
		return nil
	}

	return at
}

func (p *Parser) parseSchemaElementTest() ast.NodeTest {
	set := &ast.SchemaElementTest{}

	if !p.expectPeek(token.LPAREN) {
		// TODO error
		return nil
	}

	set.ElementDeclaration = p.parseEQName()

	if !p.expectPeek(token.RPAREN) {
		// TODO error
		return nil
	}

	return set
}

func (p *Parser) parseSchemaAttributeTest() ast.NodeTest {
	sat := &ast.SchemaAttributeTest{}

	if !p.expectPeek(token.LPAREN) {
		// TODO error
		return nil
	}

	sat.AttributeDeclaration = p.parseEQName()

	if !p.expectPeek(token.RPAREN) {
		// TODO error
		return nil
	}

	return sat
}

func (p *Parser) parsePITest() ast.NodeTest {
	pit := &ast.PITest{}

	if !p.expectPeek(token.LPAREN) {
		// TODO error
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
		return nil
	}

	return pit
}

func (p *Parser) parseCommentTest() ast.NodeTest {
	ct := &ast.CommentTest{}

	if !p.expectPeek(token.LPAREN) {
		// TODO error
		return nil
	}

	if !p.expectPeek(token.RPAREN) {
		// TODO error
		return nil
	}

	return ct
}

func (p *Parser) parseTextTest() ast.NodeTest {
	tt := &ast.TextTest{}

	if !p.expectPeek(token.LPAREN) {
		// TODO error
		return nil
	}

	if !p.expectPeek(token.RPAREN) {
		// TODO error
		return nil
	}

	return tt
}

func (p *Parser) parseNamespaceNodeTest() ast.NodeTest {
	nnt := &ast.NamespaceNodeTest{}

	if !p.expectPeek(token.LPAREN) {
		// TODO error
		return nil
	}

	if !p.expectPeek(token.RPAREN) {
		// TODO error
		return nil
	}

	return nnt
}

func (p *Parser) parseAnyKindTest() ast.NodeTest {
	akt := &ast.AnyKindTest{}

	if !p.expectPeek(token.LPAREN) {
		// TODO error
		return nil
	}

	if !p.expectPeek(token.RPAREN) {
		// TODO error
		return nil
	}

	return akt
}

func (p *Parser) parseItemTest() ast.NodeTest {
	it := &ast.ItemTest{}

	if !p.expectPeek(token.LPAREN) {
		// TODO error
		return nil
	}

	if !p.expectPeek(token.RPAREN) {
		// TODO error
		return nil
	}

	return it
}

func (p *Parser) parseFunctionTest() ast.NodeTest {
	ft := &ast.FunctionTest{}

	if !p.expectPeek(token.LPAREN) {
		// TODO error
		return nil
	}

	if p.expectPeek(token.ASTERISK) {
		ft.NodeTest = &ast.AnyFunctionTest{}

		if !p.expectPeek(token.RPAREN) {
			// TODO error
			return nil
		}
	} else {
		tft := &ast.TypedFunctionTest{}

		if !p.curTokenIs(token.RPAREN) {
			for {
				st := p.parseSequenceType()
				tft.ParamSTypes = append(tft.ParamSTypes, st)

				if !p.expectPeek(token.COMMA) {
					break
				}
				p.nextToken()
			}
		}

		if !p.expectPeek(token.AS) {
			// TODO error
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
		// TODO error
		return nil
	}

	if p.expectPeek(token.ASTERISK) {
		mt.NodeTest = &ast.AnyMapTest{}
	} else {
		tmt := &ast.TypedMapTest{}
		tmt.AtomicOrUnionType.EQName = p.parseEQName()

		if !p.expectPeek(token.COMMA) {
			// TODO error
			return nil
		}
		p.nextToken()

		tmt.SequenceType = p.parseSequenceType()
		mt.NodeTest = tmt
	}

	if !p.expectPeek(token.RPAREN) {
		// TODO error
		return nil
	}

	return mt
}

func (p *Parser) parseArrayTest() ast.NodeTest {
	at := &ast.ArrayTest{}

	if !p.expectPeek(token.LPAREN) {
		// TODO error
		return nil
	}

	if p.expectPeek(token.ASTERISK) {
		at.NodeTest = &ast.AnyArrayTest{}
	} else {
		tat := &ast.TypedArrayTest{}
		tat.SequenceType = p.parseSequenceType()

		at.NodeTest = tat
	}

	if !p.expectPeek(token.RPAREN) {
		// TODO error
		return nil
	}

	return at
}

func (p *Parser) parseAtomicOrUnionType() ast.NodeTest {
	aout := &ast.AtomicOrUnionType{}
	aout.EQName.SetValue(p.curToken.Literal)

	return aout
}

func (p *Parser) parseParenthesizedItemType() ast.NodeTest {
	pit := &ast.ParenthesizedItemType{}

	p.nextToken()
	pit.NodeTest = p.parseItemType()

	if !p.expectPeek(token.RPAREN) {
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

func (p *Parser) parseEQName() ast.EQName {
	name := ast.EQName{}

	if p.peekTokenIs(token.COLON) {
		var sb strings.Builder
		sb.WriteString(p.curToken.Literal)
		p.nextToken()
		sb.WriteString(p.curToken.Literal)
		p.nextToken()
		sb.WriteString(p.curToken.Literal)

		name.SetValue(sb.String())
	} else {
		name.SetValue(p.curToken.Literal)
	}

	return name
}

func (p *Parser) parseSequenceType() ast.SequenceType {
	st := ast.SequenceType{}

	if p.curToken.Literal == "empty-sequence" {
		if !p.expectPeek(token.LPAREN) {
			// TODO error
		}

		if !p.expectPeek(token.RPAREN) {
			// TODO error
		}
		st.TypeID = 1
	} else {
		st.NodeTest = p.parseItemType()

		if !p.expectPeek(token.RPAREN) {
			st.OccurrenceIndicator = ast.OccurrenceIndicator{Token: p.curToken}
			p.nextToken()
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
