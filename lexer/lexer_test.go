package lexer

import (
	"testing"

	"github.com/zzossig/xpath/token"
)

func TestNextToken(t *testing.T) {
	input := `(
		100
		, "100"
		, '100'
		, "string"""
		, 'string'''
		, 100.0
		, 1.0e2
		, 1.0E2
		, .2e-2
		, 1+1
		, 1e2  +     1.5 - 1  *   .215
		, 1 -- ++ +- -+ -+- +-+ - -- 1
		, 1 + ++ +- -+ -+- +-+ - - 1
		, xs:true()
		, xs:date('2021-01-13')
		, /xs
		, //*
		, //element(*, xs:date)
		, /company/office[@location = 'Boston']
		, /company/office/../following-sibling::office/@location
		, let $pi := 3.14,
			$area := function ($arg)
				{
					'area = ' ||	$pi * $arg * $arg
				},
			$r := 5
			return $area($r)
		, '6' cast           as   	xs:integer
		, //return[@return="return"]
	)`

	tokens := []struct {
		expectedType    token.Type
		expectedLiteral string
	}{
		{token.LPAREN, "("},
		{token.INT, "100"},
		{token.COMMA, ","},
		{token.STRING, "100"},
		{token.COMMA, ","},
		{token.STRING, "100"},
		{token.COMMA, ","},
		{token.STRING, "string\""},
		{token.COMMA, ","},
		{token.STRING, "string'"},
		{token.COMMA, ","},
		{token.DECIMAL, "100.0"},
		{token.COMMA, ","},
		{token.DOUBLE, "1.0e2"},
		{token.COMMA, ","},
		{token.DOUBLE, "1.0E2"},
		{token.COMMA, ","},
		{token.DOUBLE, ".2e-2"},
		{token.COMMA, ","},
		{token.INT, "1"},
		{token.PLUS, "+"},
		{token.INT, "1"},
		{token.COMMA, ","},
		{token.DOUBLE, "1e2"},
		{token.PLUS, "+"},
		{token.DECIMAL, "1.5"},
		{token.MINUS, "-"},
		{token.INT, "1"},
		{token.ASTERISK, "*"},
		{token.DECIMAL, ".215"},
		{token.COMMA, ","},
		{token.INT, "1"},
		{token.PLUS, "+"},
		{token.INT, "1"},
		{token.COMMA, ","},
		{token.INT, "1"},
		{token.MINUS, "-"},
		{token.INT, "1"},
		{token.COMMA, ","},
		{token.NS, "xs"},
		{token.COLON, ":"},
		{token.BIF, "true"},
		{token.LPAREN, "("},
		{token.RPAREN, ")"},
		{token.COMMA, ","},
		{token.NS, "xs"},
		{token.COLON, ":"},
		{token.XTYPEF, "date"},
		{token.LPAREN, "("},
		{token.STRING, "2021-01-13"},
		{token.RPAREN, ")"},
		{token.COMMA, ","},
		{token.SLASH, "/"},
		{token.IDENT, "xs"},
		{token.COMMA, ","},
		{token.DSLASH, "//"},
		{token.ASTERISK, "*"},
		{token.COMMA, ","},
		{token.DSLASH, "//"},
		{token.BIF, "element"},
		{token.LPAREN, "("},
		{token.ASTERISK, "*"},
		{token.COMMA, ","},
		{token.NS, "xs"},
		{token.COLON, ":"},
		{token.XTYPE, "date"},
		{token.RPAREN, ")"},
		{token.COMMA, ","},
		{token.SLASH, "/"},
		{token.IDENT, "company"},
		{token.SLASH, "/"},
		{token.IDENT, "office"},
		{token.LBRACKET, "["},
		{token.AT, "@"},
		{token.IDENT, "location"},
		{token.EQ, "="},
		{token.STRING, "Boston"},
		{token.RBRACKET, "]"},
		{token.COMMA, ","},
		{token.SLASH, "/"},
		{token.IDENT, "company"},
		{token.SLASH, "/"},
		{token.IDENT, "office"},
		{token.SLASH, "/"},
		{token.DDOT, ".."},
		{token.SLASH, "/"},
		{token.AXIS, "following-sibling"},
		{token.DCOLON, "::"},
		{token.IDENT, "office"},
		{token.SLASH, "/"},
		{token.AT, "@"},
		{token.IDENT, "location"},
		{token.COMMA, ","},
		{token.LET, "let"},
		{token.VAR, "$pi"},
		{token.ASSIGN, ":="},
		{token.DECIMAL, "3.14"},
		{token.COMMA, ","},
		{token.VAR, "$area"},
		{token.ASSIGN, ":="},
		{token.FUNCTION, "function"},
		{token.LPAREN, "("},
		{token.VAR, "$arg"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.STRING, "area = "},
		{token.DVBAR, "||"},
		{token.VAR, "$pi"},
		{token.ASTERISK, "*"},
		{token.VAR, "$arg"},
		{token.ASTERISK, "*"},
		{token.VAR, "$arg"},
		{token.RBRACE, "}"},
		{token.COMMA, ","},
		{token.VAR, "$r"},
		{token.ASSIGN, ":="},
		{token.INT, "5"},
		{token.RETURN, "return"},
		{token.VAR, "$area"},
		{token.LPAREN, "("},
		{token.VAR, "$r"},
		{token.RPAREN, ")"},
		{token.COMMA, ","},
		{token.STRING, "6"},
		{token.CAST, "cast"},
		{token.AS, "as"},
		{token.NS, "xs"},
		{token.COLON, ":"},
		{token.XTYPE, "integer"},
		{token.COMMA, ","},
		{token.DSLASH, "//"},
		{token.IDENT, "return"},
		{token.LBRACKET, "["},
		{token.AT, "@"},
		{token.IDENT, "return"},
		{token.EQ, "="},
		{token.STRING, "return"},
		{token.RBRACKET, "]"},
		{token.RPAREN, ")"},
	}

	lexer := New(input)

	for i, tt := range tokens {
		tok := lexer.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("TestNextToken:type[%d] - expected=%q, got=%q", i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("TestNextToken:literal[%d] - expected=%q, got=%q", i, tt.expectedLiteral, tok.Literal)
		}
	}
}
