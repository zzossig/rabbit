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
		, upper-case('a')
		, abs(-2)
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
		, 1 eq 2
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
		{token.IDENT, "upper-case"},
		{token.LPAREN, "("},
		{token.STRING, "a"},
		{token.RPAREN, ")"},
		{token.COMMA, ","},
		{token.IDENT, "abs"},
		{token.LPAREN, "("},
		{token.MINUS, "-"},
		{token.INT, "2"},
		{token.RPAREN, ")"},
		{token.COMMA, ","},
		{token.IDENT, "xs"},
		{token.COLON, ":"},
		{token.IDENT, "true"},
		{token.LPAREN, "("},
		{token.RPAREN, ")"},
		{token.COMMA, ","},
		{token.IDENT, "xs"},
		{token.COLON, ":"},
		{token.IDENT, "date"},
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
		{token.ELEMENT, "element"},
		{token.LPAREN, "("},
		{token.ASTERISK, "*"},
		{token.COMMA, ","},
		{token.IDENT, "xs"},
		{token.COLON, ":"},
		{token.IDENT, "date"},
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
		{token.IDENT, "following-sibling"},
		{token.DCOLON, "::"},
		{token.IDENT, "office"},
		{token.SLASH, "/"},
		{token.AT, "@"},
		{token.IDENT, "location"},
		{token.COMMA, ","},
		{token.LET, "let"},
		{token.DOLLAR, "$"},
		{token.IDENT, "pi"},
		{token.ASSIGN, ":="},
		{token.DECIMAL, "3.14"},
		{token.COMMA, ","},
		{token.DOLLAR, "$"},
		{token.IDENT, "area"},
		{token.ASSIGN, ":="},
		{token.FUNCTION, "function"},
		{token.LPAREN, "("},
		{token.DOLLAR, "$"},
		{token.IDENT, "arg"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.STRING, "area = "},
		{token.DVBAR, "||"},
		{token.DOLLAR, "$"},
		{token.IDENT, "pi"},
		{token.ASTERISK, "*"},
		{token.DOLLAR, "$"},
		{token.IDENT, "arg"},
		{token.ASTERISK, "*"},
		{token.DOLLAR, "$"},
		{token.IDENT, "arg"},
		{token.RBRACE, "}"},
		{token.COMMA, ","},
		{token.DOLLAR, "$"},
		{token.IDENT, "r"},
		{token.ASSIGN, ":="},
		{token.INT, "5"},
		{token.RETURN, "return"},
		{token.DOLLAR, "$"},
		{token.IDENT, "area"},
		{token.LPAREN, "("},
		{token.DOLLAR, "$"},
		{token.IDENT, "r"},
		{token.RPAREN, ")"},
		{token.COMMA, ","},
		{token.STRING, "6"},
		{token.CAST, "cast"},
		{token.AS, "as"},
		{token.IDENT, "xs"},
		{token.COLON, ":"},
		{token.IDENT, "integer"},
		{token.COMMA, ","},
		{token.DSLASH, "//"},
		{token.RETURN, "return"},
		{token.LBRACKET, "["},
		{token.AT, "@"},
		{token.RETURN, "return"},
		{token.EQ, "="},
		{token.STRING, "return"},
		{token.RBRACKET, "]"},
		{token.COMMA, ","},
		{token.INT, "1"},
		{token.EQV, "eq"},
		{token.INT, "2"},
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
