package parser

import (
	"testing"

	"github.com/zzossig/xpath/lexer"
)

func TestGrouped(t *testing.T) {
	input := "(5 + 5) * 2"
	l := lexer.New(input)
	p := New(l)
	xpath := p.ParseXPath()

	// if len(xpath.Items) != 2 {
	// 	t.Errorf("wrong num, got:%d, %q", len(xpath.Items), xpath.Items)
	// }

	if xpath.String() != "((5 + 5) * 2)" {
		t.Errorf("wrong value, got:%q", xpath.String())
	}
}

func TestArithmeticExpr(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"2 + 2",
			"(2 + 2)",
		},
		{
			"(2 * 2)",
			"(2 * 2)",
		},
		{
			"2-2",
			"(2 - 2)",
		},
		{
			"(2*2)",
			"(2 * 2)",
		},
		{
			"(2 idiv 2)",
			"(2 idiv 2)",
		},
		{
			"(2 div 2)",
			"(2 div 2)",
		},
		{
			"(2 mod 2)",
			"(2 mod 2)",
		},
		{
			"2 -- ++ +- -+ -+- +-+ - -- 2",
			"(2 + 2)",
		},
		{
			"2 + ++ +- -+ -+- +-+ - - 2",
			"(2 - 2)",
		},
		{
			"2 + 2 - 2",
			"((2 + 2) - 2)",
		},
		{
			"2 + 2 - 2 + 2",
			"(((2 + 2) - 2) + 2)",
		},
		{
			"(2-2, 2 + 2)",
			"((2 - 2), (2 + 2))",
		},
		{
			"(2-2, 2 + 2,(1,3,4+ 1))",
			"((2 - 2), (2 + 2), 1, 3, (4 + 1))",
		},
		{
			"(2-2, 2 + 2,(1,3,4+ 1,( 8 -1)))",
			"((2 - 2), (2 + 2), 1, 3, (4 + 1), (8 - 1))",
		},
		{
			"1 + (2 + 3) + 4",
			"((1 + (2 + 3)) + 4)",
		},
		{
			"5 - 5 * 2",
			"(5 - (5 * 2))",
		},
		// {
		// 	"(5 + 5) * 2",
		// 	"((5 + 5) * 2)",
		// },
		{
			"2 div 2 * 2 idiv 2 mod 2",
			"((((2 div 2) * 2) idiv 2) mod 2)",
		},
		{
			"-1",
			"(-1)",
		},
		{
			"-(5 + 5)",
			"(-(5 + 5))",
		},
		{
			"-(5 + 5) * 4",
			"((-(5 + 5)) * 4)",
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		xpath := p.ParseXPath()

		actual := xpath.String()
		if actual != tt.expected {
			t.Errorf("expected=%q, got=%q", tt.expected, actual)
		}
	}
}
