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
			"((-3) div 2)",
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
			"for $i in (10, 20), $j in (1, 2), $k in (4, 5) return ($i + $j + $k)",
			"for $i in (10, 20), $j in (1, 2), $k in (4, 5) return (($i + $j) + $k)",
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
		// {
		// 	"let $x := 1, $y := 2\nreturn $x + $y",
		// 	"let $x := 1, $y := 2 return ($x + $y)",
		// },
		// {
		// 	"let $x := doc('a.xml')/*, $y := $x//* return $y[@value gt $x/@min]",
		// 	"let $x := (doc('a.xml') / *), $y := ($x // *) return $y[(@value gt ($x / @min))]",
		// },
		{
			`
				let $tax_rate :=
          function($rate as xs:integer, $amount as xs:decimal) as xs:decimal
          {
            ($rate div 100) * $amount
          },

          $income_tax :=
          function($amount as xs:decimal) as xs:decimal
          {
           $tax_rate(15, ?)($amount)
          },

         $luxury_tax :=
          function($amount as xs:integer) as xs:decimal
          {
            $tax_rate(50, ?)($amount)
          }

        return 
            ($income_tax(300), $luxury_tax(50))
			`,
			"let $tax_rate := function($rate as xs:integer, $amount as xs:decimal) as xs:decimal {(($rate div 100) * $amount)}, $income_tax := function($amount as xs:decimal) as xs:decimal {tax_rate(15, ?)($amount)}, $luxury_tax := function($amount as xs:integer) as xs:decimal {tax_rate(50, ?)($amount)} return (income_tax(300), luxury_tax(50))",
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
			"function($a, $b){$b || $a}('World', 'Hello')",
			"function($a, $b) {($b || $a)}('World', 'Hello')",
		},
		{
			"array { (), (27, 17, 0) }(1)",
			"array{(), (27, 17, 0)}(1)",
		},
		{
			"array { 'licorice', 'ginger' }(20)",
			"array{'licorice', 'ginger'}(20)",
		},
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
			"$f(2, 3)",
			"f(2, 3)",
		},
		{
			"$f[2]('Hi there')",
			"$f[2]('Hi there')",
		},
		{
			"$f()[2]",
			"f()[2]",
		},
		{
			"function() as xs:integer+ { 2, 3, 5, 7, 11, 13 }()",
			"function() as xs:integer+ {2, 3, 5, 7, 11, 13}()",
		},
		{
			"function () {'hello world'}()",
			"function() {'hello world'}()",
		},
		{
			"(1,2,function () {'hello world'})[3]()",
			"(1, 2, function() {'hello world'})[3]()",
		},
		{
			`fn:for-each-pair(("a", "b", "c"), ("x", "y", "z"), concat#2)`,
			"fn:for-each-pair(('a', 'b', 'c'), ('x', 'y', 'z'), concat#2)",
		},
		{
			"for-each-pair( 1 to 5, ( 'London', 'New York', 'Vienna', 'Paris', 'Tokyo' ), concat( ?, ' ',  ? ) )",
			"for-each-pair((1 to 5), ('London', 'New York', 'Vienna', 'Paris', 'Tokyo'), concat(?, ' ', ?))",
		},
		{
			"concat#3('a', 'b', 'c')",
			"concat#3('a', 'b', 'c')",
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
			"map{}?*",
			"map{}?*",
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
		{
			`(map{"name": "Jack", "age": 1}, map{"name": "Mike", "age": 2})[?name='Jack']`,
			`(map{'name': 'Jack', 'age': 1}, map{'name': 'Mike', 'age': 2})[((?)name = 'Jack')]`,
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

func TestStringConcatExpr(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			`"con" || "cat" || "enate"`,
			`(('con' || 'cat') || 'enate')`,
		},
		{
			`1|| "B"`,
			`(1 || 'B')`,
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

func TestSequence(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			`(10, (1, 2), (), (3, 4))`,
			`(10, (1, 2), (), (3, 4))`,
		},
		{
			`10, (1, 2), (), (3, 4)`,
			`10, (1, 2), (), (3, 4)`,
		},
		{
			"(salary, bonus)",
			"(salary, bonus)",
		},
		{
			"($price, $price)",
			"($price, $price)",
		},
		{
			"(10, 1 to 4)",
			"(10, (1 to 4))",
		},
		{
			"15 to 10",
			"(15 to 10)",
		},
		{
			"fn:reverse(10 to 15)",
			"fn:reverse((10 to 15))",
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

func TestSequenceType(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			`5 instance of xs:integer`,
			`5 instance of xs:integer`,
		},
		{
			"(5, 6) instance of xs:integer+",
			"(5, 6) instance of xs:integer+",
		},
		{
			". instance of element()",
			". instance of element()",
		},
		{
			"if ($x castable as hatsize) then $x cast as hatsize else if ($x castable as IQ) then $x cast as IQ else $x cast as xs:string",
			"if($x castable as hatsize) then $x cast as hatsize else if($x castable as IQ) then $x cast as IQ else $x cast as xs:string",
		},
		{
			"$myaddress treat as element(*, USAddress)",
			"$myaddress treat as element(*, USAddress)",
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

func TestPostfixExpr(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"$products[price gt 100]",
			"$products[(price gt 100)]",
		},
		{
			"(1 to 100)[. mod 5 eq 0]",
			"(1 to 100)[((. mod 5) eq 0)]",
		},
		{
			"(21 to 29)[5]",
			"(21 to 29)[5]",
		},
		{
			"$orders[fn:position() = (5 to 9)]",
			"$orders[(fn:position() = (5 to 9))]",
		},
		{
			"$book/(chapter | appendix)[fn:last()]",
			"($book / (chapter | appendix)[fn:last()])",
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

func TestNodeTest(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"node()",
			"node()",
		},
		{
			"text()",
			"text()",
		},
		{
			"comment()",
			"comment()",
		},
		{
			"namespace-node()",
			"namespace-node()",
		},
		{
			"element()",
			"element()",
		},
		{
			"schema-element(person)",
			"schema-element(person)",
		},
		{
			"element(person)",
			"element(person)",
		},
		{
			"element(person, surgeon)",
			"element(person, surgeon)",
		},
		{
			"element(*, surgeon)",
			"element(*, surgeon)",
		},
		{
			"attribute()",
			"attribute()",
		},
		{
			"attribute(price)",
			"attribute(price)",
		},
		{
			"attribute(*, xs:decimal)",
			"attribute(*, xs:decimal)",
		},
		{
			"document-node()",
			"document-node()",
		},
		{
			"document-node(element(book))",
			"document-node(element(book))",
		},
		{
			"child::chapter[2]",
			"child::chapter[2]",
		},
		{
			"descendant::toy[attribute::color = 'red']",
			"descendant::toy[(attribute::color = 'red')]",
		},
		{
			"child::employee[secretary][assistant]",
			"child::employee[secretary][assistant]",
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

func TestUnabbreviatedSyntax(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"child::para",
			"child::para",
		},
		{
			"child::*",
			"child::*",
		},
		{
			"child::text()",
			"child::text()",
		},
		{
			"child::node()",
			"child::node()",
		},
		{
			"attribute::name",
			"attribute::name",
		},
		{
			"attribute::*",
			"attribute::*",
		},
		{
			"parent::node()",
			"parent::node()",
		},
		{
			"descendant::para",
			"descendant::para",
		},
		{
			"ancestor::div",
			"ancestor::div",
		},
		{
			"ancestor-or-self::div",
			"ancestor-or-self::div",
		},
		{
			"descendant-or-self::para",
			"descendant-or-self::para",
		},
		{
			"self::para",
			"self::para",
		},
		{
			"child::chapter/descendant::para",
			"(child::chapter / descendant::para)",
		},
		{
			"child::*/child::para",
			"(child::* / child::para)",
		},
		{
			"/",
			"/",
		},
		{
			"/descendant::para",
			"/descendant::para",
		},
		{
			"/descendant::list/child::member",
			"(/descendant::list / child::member)",
		},
		{
			"child::para[fn:position() = 1]",
			"child::para[(fn:position() = 1)]",
		},
		{
			"child::para[fn:position() = fn:last()]",
			"child::para[(fn:position() = fn:last())]",
		},
		{
			"child::para[fn:position() = fn:last()-1]",
			"child::para[(fn:position() = (fn:last() - 1))]",
		},
		{
			"child::para[fn:position() > 1]",
			"child::para[(fn:position() > 1)]",
		},
		{
			"following-sibling::chapter[fn:position() = 1]",
			"following-sibling::chapter[(fn:position() = 1)]",
		},
		{
			"preceding-sibling::chapter[fn:position() = 1]",
			"preceding-sibling::chapter[(fn:position() = 1)]",
		},
		{
			"/descendant::figure[fn:position() = 42]",
			"/descendant::figure[(fn:position() = 42)]",
		},
		{
			"/child::book/child::chapter[fn:position() = 5]/child::section[fn:position() = 2]",
			"((/child::book / child::chapter[(fn:position() = 5)]) / child::section[(fn:position() = 2)])",
		},
		{
			"child::para[attribute::type eq 'warning']",
			"child::para[(attribute::type eq 'warning')]",
		},
		{
			"child::para[attribute::type eq 'warning'][fn:position() = 5]",
			"child::para[(attribute::type eq 'warning')][(fn:position() = 5)]",
		},
		{
			"child::para[fn:position() = 5][attribute::type eq 'warning']",
			"child::para[(fn:position() = 5)][(attribute::type eq 'warning')]",
		},
		{
			"child::chapter[child::title = 'Introduction']",
			"child::chapter[(child::title = 'Introduction')]",
		},
		{
			"child::chapter[child::title]",
			"child::chapter[child::title]",
		},
		{
			"child::*[self::chapter or self::appendix]",
			"child::*[(self::chapter or self::appendix)]",
		},
		{
			"child::*[self::chapter or self::appendix][fn:position() = fn:last()]",
			"child::*[(self::chapter or self::appendix)][(fn:position() = fn:last())]",
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

func TestAbbreviatedSyntax(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"para",
			"para",
		},
		{
			"*",
			"*",
		},
		{
			"text()",
			"text()",
		},
		{
			"@name",
			"@name",
		},
		{
			"@*",
			"@*",
		},
		{
			"para[1]",
			"para[1]",
		},
		{
			"para[fn:last()]",
			"para[fn:last()]",
		},
		{
			"*/para",
			"(* / para)",
		},
		{
			"/book/chapter[5]/section[2]",
			"((/book / chapter[5]) / section[2])",
		},
		{
			"chapter//para",
			"(chapter // para)",
		},
		{
			"//para",
			"//para",
		},
		{
			"//@version",
			"//@version",
		},
		{
			"//list/member",
			"(//list / member)",
		},
		{
			".//para",
			"(. // para)",
		},
		{
			"..",
			"..",
		},
		{
			"../@lang",
			"(.. / @lang)",
		},
		{
			"para[@type='warning']",
			"para[(@type = 'warning')]",
		},
		{
			"para[@type='warning'][5]",
			"para[(@type = 'warning')][5]",
		},
		{
			"para[5][@type='warning']",
			"para[5][(@type = 'warning')]",
		},
		{
			"chapter[title='Introduction']",
			"chapter[(title = 'Introduction')]",
		},
		{
			"chapter[title]",
			"chapter[title]",
		},
		{
			"employee[@secretary and @assistant]",
			"employee[(@secretary and @assistant)]",
		},
		{
			"book/(chapter|appendix)/section",
			"((book / (chapter | appendix)) / section)",
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

func TestEQName(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"abc",
			"abc",
		},
		{
			"fn:abc",
			"fn:abc",
		},
		{
			"Q{http://www.w3.org/2005/xpath-functions/math}pi",
			"Q{http://www.w3.org/2005/xpath-functions/math}pi",
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

func TestPathExpr(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		// {
		// 	"/company",
		// 	"/company",
		// },
		// {
		// 	"//company",
		// 	"//company",
		// },
		// {
		// 	"//company/office",
		// 	"(//company / office)",
		// },
		{
			"//company/office/department",
			"((//company / office) / department)",
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
