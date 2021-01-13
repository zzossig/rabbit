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
		, 100.0
		, 1.0e2
		, 1.0E2
		, xs:true()
		, /xs
		, /company
		, //company
		, /company/office[@location = 'Boston']
		, /company/office/../following-sibling::office/@location
	)`
	// , xs:date('2021-01-13')
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
		{token.DECIMAL, "100.0"},
		{token.COMMA, ","},
		{token.DOUBLE, "1.0e2"},
		{token.COMMA, ","},
		{token.DOUBLE, "1.0E2"},
		{token.COMMA, ","},
		{token.XSCHEMA, "xs"},
		{token.COLON, ":"},
		{token.IDENT, "true"},
		{token.LPAREN, "("},
		{token.RPAREN, ")"},
		{token.COMMA, ","},
		{token.SLASH, "/"},
		{token.IDENT, "xs"},
		{token.COMMA, ","},
		{token.SLASH, "/"},
		{token.IDENT, "company"},
		{token.COMMA, ","},
		{token.DSLASH, "//"},
		{token.IDENT, "company"},
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
