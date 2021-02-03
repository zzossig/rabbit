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
			"-3 div 2",
			"(((-)3) div 2)",
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
			"2*(5 + 5)",
			"(2 * (5 + 5))",
		},
		{
			"2 div 2 * 2 idiv 2 mod 2",
			"((((2 div 2) * 2) idiv 2) mod 2)",
		},
		{
			"-1",
			"((-)1)",
		},
		{
			"-(5 + 5)",
			"((-)(5 + 5))",
		},
		{
			"-(5 + 5) * 4",
			"(((-)(5 + 5)) * 4)",
		},
		{
			"1,2,3",
			"1, 2, 3",
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
			"array{}",
		},
		{
			"array	{}",
			"array{}",
		},
		{
			"(array{}, 1+1)",
			"(array{}, (1 + 1))",
		},
		{
			"(array{}, array {2 idiv 2 * 3})",
			"(array{}, array{((2 idiv 2) * 3)})",
		},
		{
			"array{1,2+3,4+5,6*(7+8)}",
			"array{1, (2 + 3), (4 + 5), (6 * (7 + 8))}",
		},
		{
			"array {2-2, 2 + 2,((1+2, (3+4, 5+6)),3,4+ 1,( 8 -1))}",
			"array{(2 - 2), (2 + 2), (((1 + 2), ((3 + 4), (5 + 6))), 3, (4 + 1), (8 - 1))}",
		},
		{
			"array{(5+5)*2}",
			"array{((5 + 5) * 2)}",
		},
		{
			"array { $x }",
			"array{$x}",
		},
		{
			"array { local:items() }",
			"array{local:items()}",
		},
		{
			"array{2 div 2 * 2 idiv 2 mod 2}",
			"array{((((2 div 2) * 2) idiv 2) mod 2)}",
		},
		{
			"[1,2,3]",
			"[1, 2, 3]",
		},
		{
			"[ (), (27, 17, 0)]",
			"[(), (27, 17, 0)]",
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

func TestSimpleMapExpr(t *testing.T) {
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
		{
			`child::div1 / child::para / string() ! concat("id-", .)`,
			`(((child::div1 / child::para) / string()) ! concat('id-', .))`,
		},
		{
			"$emp ! (@if, @middle, @last)",
			"($emp ! (@if, @middle, @last))",
		},
		{
			"$docs ! ( //employee)",
			"($docs ! //employee)",
		},
		{
			"avg( //employee / salary ! translate(., '$', '') ! number(.))",
			"avg((((//employee / salary) ! translate(., '$', '')) ! number(.)))",
		},
		{
			`fn:string-join((1 to $n)!"*")`,
			`fn:string-join(((1 to $n) ! '*'))`,
		},
		{
			"$values!(.*.) => fn:sum()",
			"($values ! (. * .)) => fn:sum()",
		},
		{
			"string-join(ancestor::*!name(), '/')",
			"string-join((ancestor::* ! name()), '/')",
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
			"((1, 2) != (2, 3)), 5, 6",
		},
		{
			"((6,5),9,7,(1, 2) != (2, 3))",
			"((6, 5), 9, 7, ((1, 2) != (2, 3)))",
		},
		{
			`[ "Obama", "Nixon", "Kennedy" ] = "Kennedy"`,
			`(['Obama', 'Nixon', 'Kennedy'] = 'Kennedy')`,
		},
		{
			`$book1/author eq "Kennedy"`,
			`(($book1 / author) eq 'Kennedy')`,
		},
		{
			`[ "Kennedy" ] eq "Kennedy"`,
			`(['Kennedy'] eq 'Kennedy')`,
		},
		{
			"//product[weight gt 100]",
			"//product[(weight gt 100)]",
		},
		{
			`fn:QName('http://example.com/ns1', 'this:color') eq fn:QName('http://example.com/ns1', 'that:color')`,
			`(fn:QName('http://example.com/ns1', 'this:color') eq fn:QName('http://example.com/ns1', 'that:color'))`,
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
		{
			"if ($widget1/unit-cost < $widget2/unit-cost) then $widget1 else $widget2",
			"if((($widget1 / unit-cost) < ($widget2 / unit-cost))) then $widget1 else $widget2",
		},
		{
			"if ($part/@discounted) then $part/wholesale else $part/retail",
			"if(($part / @discounted)) then ($part / wholesale) else ($part / retail)",
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

func TestForExpr(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"for $i in (10,20),\n$j in (1,2)\nreturn ($i + $j)",
			"for $i in (10, 20), $j in (1, 2) return ($i + $j)",
		},
		{
			"for $a in fn:distinct-values(book/author) return ((book/author[. = $a])[1], book[author = $a]/title)",
			"for $a in fn:distinct-values((book / author)) return ((book / author[(. = $a)])[1], (book[(author = $a)] / title))",
		},
		{
			"for $x in $z, $y in f($x) return g($x, $y)",
			"for $x in $z, $y in f($x) return g($x, $y)",
		},
		{
			"fn:sum(for $i in order-item return $i/@price * $i/@qty)",
			"fn:sum(for $i in order-item return (($i / @price) * ($i / @qty)))",
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

func TestLetExpr(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"let $x := 1, $y := 2\nreturn $x + $y",
			"let $x := 1, $y := 2 return ($x + $y)",
		},
		{
			"let $x := doc('a.xml')/*, $y := $x//* return $y[@value gt $x/@min]",
			"let $x := (doc('a.xml') / *), $y := ($x // *) return $y[(@value gt ($x / @min))]",
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

func TestLogicalExpr(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"1 and 2",
			"(1 and 2)",
		},
		{
			"1 and 1+1 or 2",
			"((1 and (1 + 1)) or 2)",
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
			"map { 'first' : 'Jenna', 'last' : 'Scott' }",
			"map{'first': 'Jenna', 'last': 'Scott'}",
		},
		{
			`map {
				'book': map {
					'title': 'Data on the Web',
					'year': 2000,
					'author': [
						map {
							'last': 'Abiteboul',
							'first': 'Serge'
						},
						map {
							'last': 'Buneman',
							'first': 'Peter'
						},
						map {
							'last': 'Suciu',
							'first': 'Dan'
						}
					],
					'publisher': 'Morgan Kaufmann Publishers',
					'price': 39.95
				}
			}`,
			`map{'book': map{'title': 'Data on the Web', 'year': 2000, 'author': [map{'last': 'Abiteboul', 'first': 'Serge'}, map{'last': 'Buneman', 'first': 'Peter'}, map{'last': 'Suciu', 'first': 'Dan'}], 'publisher': 'Morgan Kaufmann Publishers', 'price': 39.950000}}`,
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

func TestFunctionCall(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"[ 1, 2, 5, 7 ](4)",
			"[1, 2, 5, 7](4)",
		},
		{
			"[ [1, 2, 3], [4, 5, 6]](2)",
			"[[1, 2, 3], [4, 5, 6]](2)",
		},
		{
			"[(), [1, 2, 3],[4, 5, 6]](2)(2)",
			"[(), [1, 2, 3], [4, 5, 6]](2)(2)",
		},
		{
			"array { (), (27, 17, 0) }(1)",
			"array{(), (27, 17, 0)}(1)",
		},
		{
			"array { 'licorice', 'ginger' }(20)",
			"array{'licorice', 'ginger'}(20)",
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

func TestLookup(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"(1)[?2 = 5]",
			"1[((?)2 = 5)]",
		},
		{
			"(1,3)[?2 = 5]",
			"(1, 3)[((?)2 = 5)]",
		},
		{
			"([1,2,3], [1,2,5], [1,2])[?3 = 5]",
			"([1, 2, 3], [1, 2, 5], [1, 2])[((?)3 = 5)]",
		},
		{
			"[1, 2, 5, 7]?*",
			"[1, 2, 5, 7]?*",
		},
		{
			"[[1, 2, 3], [4, 5, 6]]?*",
			"[[1, 2, 3], [4, 5, 6]]?*",
		},
		{
			`map { "first" : "Jenna", "last" : "Scott" }?first`,
			`map{'first': 'Jenna', 'last': 'Scott'}?first`,
		},
		{
			"(map{'first': 'Tom'}, map{'first': 'Dick'}, map{'first': 'Harry'})?first",
			"(map{'first': 'Tom'}, map{'first': 'Dick'}, map{'first': 'Harry'})?first",
		},
		{
			"(map{'first': 'Tom'} ! ?first='Tom')",
			"((map{'first': 'Tom'} ! (?)first) = 'Tom')",
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

func TestQuantifiedExpr(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"every $part in /parts/part satisfies $part/@discounted",
			"every $part in (/parts / part) satisfies ($part / @discounted)",
		},
		{
			"some $emp in /emps/employee satisfies ($emp/bonus > 0.25 * $emp/salary)",
			"some $emp in (/emps / employee) satisfies (($emp / bonus) > (0.250000 * ($emp / salary)))",
		},
		{
			"some $x in (1, 2, 3), $y in (2, 3, 4) satisfies $x + $y = 4",
			"some $x in (1, 2, 3), $y in (2, 3, 4) satisfies (($x + $y) = 4)",
		},
		{
			"every $x in (1, 2, 'cat') satisfies $x * 2 = 4",
			"every $x in (1, 2, 'cat') satisfies (($x * 2) = 4)",
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
