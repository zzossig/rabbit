package parser

import (
	"testing"

	"github.com/zzossig/xpath/lexer"
)

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
			"((2 - 2), (2 + 2), (1, 3, (4 + 1)))",
		},
		{
			"(2-2, 2 + 2,((1+2, (3+4, 5+6)),3,4+ 1,( 8 -1)))",
			"((2 - 2), (2 + 2), (((1 + 2), ((3 + 4), (5 + 6))), 3, (4 + 1), (8 - 1)))",
		},
		{
			"1 + (2 + 3) + 4",
			"((1 + (2 + 3)) + 4)",
		},
		{
			"5 - 5 * 2",
			"(5 - (5 * 2))",
		},
		{
			"(5 + 5) * 2",
			"((5 + 5) * 2)",
		},
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
		{
			"1,2,3",
			"(1, 2, 3)",
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

func TestArrayExpr(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"array{}",
			"array{()}",
		},
		{
			"array	{}",
			"array{()}",
		},
		{
			"(array{}, 1+1)",
			"(array{()}, (1 + 1))",
		},
		{
			"(array{}, array {2 idiv 2 * 3})",
			"(array{()}, array{(((2 idiv 2) * 3))})",
		},
		{
			"array{1,2+3,4+5,6*(7+8)}",
			"array{(1, (2 + 3), (4 + 5), (6 * (7 + 8)))}",
		},
		{
			"array{2 div 2 * 2 idiv 2 mod 2}",
			"array{(((((2 div 2) * 2) idiv 2) mod 2))}",
		},
		{
			"[1,2,3]",
			"[1, 2, 3]",
		},
		{
			`[ "Obama", "Nixon", "Kennedy" ]`,
			`['Obama', 'Nixon', 'Kennedy']`,
		},
		{
			"[1-1,2,2-3*5]",
			"[(1 - 1), 2, (2 - (3 * 5))]",
		},
		{
			"[(1-1),2,2-3*5]",
			"[(1 - 1), 2, (2 - (3 * 5))]",
		},
		{
			"(1,[(2,3+3,4),5],6)",
			"(1, [(2, (3 + 3), 4), 5], 6)",
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

func TestArrowExpr(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"'a' => upper-case()",
			"'a' => upper-case()",
		},
		{
			"('a' => upper-case(),1,2)",
			"('a' => upper-case(), 1, 2)",
		},
		{
			"'a' => upper-case() => normalize-unicode()",
			"'a' => upper-case() => normalize-unicode()",
		},
		{
			"'a' => tokenize('\\s+', 'abc')",
			"'a' => tokenize('\\s+', 'abc')",
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

func TestMapExpr(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"1!1",
			"(1 ! 1)",
		},
		{
			"(1 +1 )!1-1",
			"(((1 + 1) ! 1) - 1)",
		},
		{
			"1 + 1 ! 1 * 1",
			"(1 + ((1 ! 1) * 1))",
		},
		{
			"(1!1,1,2,3)",
			"((1 ! 1), 1, 2, 3)",
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

func TestComparisonExpr(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"1 <= 1",
			"(1 <= 1)",
		},
		{
			"(1, 2) != (2, 3)",
			"((1, 2) != (2, 3))",
		},
		{
			"(1, 2) != (2, 3), 5, 6",
			"(((1, 2) != (2, 3)), 5, 6)",
		},
		{
			"((6,5),9,7,(1, 2) != (2, 3))",
			"((6, 5), 9, 7, ((1, 2) != (2, 3)))",
		},
		{
			`[ "Obama", "Nixon", "Kennedy" ] = "Kennedy"`,
			`(['Obama', 'Nixon', 'Kennedy'] = 'Kennedy')`,
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

func TestIfExpr(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"if   (0) then 2 else 4",
			"if(0) then 2 else 4",
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
