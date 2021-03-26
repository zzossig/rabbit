package eval

import (
	"fmt"
	"strings"
	"testing"

	"github.com/zzossig/xpath/bif"
	"github.com/zzossig/xpath/lexer"
	"github.com/zzossig/xpath/object"
	"github.com/zzossig/xpath/parser"
)

func TestEvalArithmetic(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"5", 5},
		{"10", 10},
		{"-5", -5},
		{"-10", -10},
		{"5 + 5 + 5 idiv 5 - 10", 1},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"-50 + 100 + -50", 0},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"20 + 2 * -10", 0},
		{"50 div 2 * 2 + 10", 60.0},
		{"2 * (5 + 10)", 30},
		{"3 * 3 * 3 + 10", 37},
		{"3 * (3 * 3) + 10", 37},
		{"(5 + 10 * 2 + 15 div 3) * 2 + -10", 50.0},
		{"7 mod 4", 3},
		{"7.7 mod 2.3", 0.8},
		{"7.724 mod 2.7", 2.324},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		sequence := evaluated.(*object.Sequence)
		for _, item := range sequence.Items {
			testNumberObject(t, item, tt.expected)
		}
	}
}

func TestEvalArray(t *testing.T) {
	tests := []struct {
		input    string
		expected []interface{}
	}{
		{"array{1,2,3}", []interface{}{1, 2, 3}},
		{"array{1*2,2+3,3-4,5 idiv 5, 5 div 5}", []interface{}{2, 5, -1, 1, 1.0}},
		{"[3 mod 2, 'a', 'b', 1.1]", []interface{}{1, "a", "b", 1.1}},
	}

	for _, tt := range tests {
		seq := testEval(tt.input)
		sequence := seq.(*object.Sequence)

		for _, item := range sequence.Items {
			switch item := item.(type) {
			case *object.Array:
				if len(item.Items) != len(tt.expected) {
					t.Errorf("Array lenth not match with the expected one. got=%d, want=%d", len(item.Items), len(tt.expected))
				}
				for i, v := range item.Items {
					switch v.(type) {
					case *object.Integer:
						testNumberObject(t, v, tt.expected[i])
					case *object.Decimal:
						testNumberObject(t, v, tt.expected[i])
					case *object.Double:
						testNumberObject(t, v, tt.expected[i])
					case *object.String:
						testStringObject(t, v, tt.expected[i])
					default:
						t.Errorf("Unkown item type. got=%s", item.Type())
					}

				}
			default:
				t.Errorf("Unkown item type. got=%s", item.Type())
			}
		}
	}
}

func TestFunctionCall(t *testing.T) {
	tests := []struct {
		input    string
		expected []interface{}
	}{
		{`abs(-2.5)`, []interface{}{2.5}},
		{`abs(-2)`, []interface{}{2}},
		{
			`for-each-pair( 1 to 5, ( 'A', 'B', 'C', 'D', 'E' ), concat( ?,  ?, '--' ) )`,
			[]interface{}{"1A--", "2B--", "3C--", "4D--", "5E--"},
		},
		{
			`for-each-pair( 1, ( 'A', 'B', 'C', 'D', 'E' ), concat( ?,  ?, '--' ) )`,
			[]interface{}{"1A--"},
		},
		{
			`for-each-pair( (1 to 2), ( 'A', 'B', 'C' ), concat( ?,  ?, '--' ) )`,
			[]interface{}{"1A--", "2B--"},
		},
	}

	for _, tt := range tests {
		seq := testEval(tt.input)
		sequence := seq.(*object.Sequence)
		testSequenceObject(t, sequence, tt.expected)
	}
}

func TestStringConcat(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`"a"||"B"`, "aB"},
		{`1|| "B"`, "1B"},
		{`1 ||1.5`, "11.5"},
		{`1.2 || 1.5`, "1.21.5"},
		{`1.2 || "A" || 1.5`, "1.2A1.5"},
	}

	for _, tt := range tests {
		seq := testEval(tt.input)
		sequence := seq.(*object.Sequence)

		for _, item := range sequence.Items {
			switch item := item.(type) {
			case *object.String:
				testStringObject(t, item, tt.expected)
			default:
				t.Errorf("Unkown item type. got=%s", item.Type())
			}
		}
	}
}

func TestSimpleMapExpr(t *testing.T) {
	tests := []struct {
		input    string
		expected []interface{}
	}{
		{`(1,2,3)!concat("id-",.)`, []interface{}{"id-1", "id-2", "id-3"}},
	}

	for _, tt := range tests {
		seq := testEval(tt.input)
		sequence := seq.(*object.Sequence)
		testSequenceObject(t, sequence, tt.expected)
	}
}

func TestArrowExpr(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"'a' => upper-case()", "A"},
		{"'a' => upper-case() => lower-case()", "a"},
		{"let $var := 'a' return $var => fn:upper-case()", "A"},
		{`let $pi := 3.14, $area := function ($arg)
				{
					'area = ' ||	$pi * $arg * $arg
				},
			$r := 5
			return $r => $area()`,
			"area = 78.5",
		},
		{`let $area := function ($arg, $pi)
      {
         'area = ' ||	$pi * $arg * $arg
      },
			$r := 5
			return $r => $area(3.14)`,
			"area = 78.5",
		},
	}

	for _, tt := range tests {
		seq := testEval(tt.input)
		sequence := seq.(*object.Sequence)

		for _, item := range sequence.Items {
			switch item := item.(type) {
			case *object.Integer:
				testNumberObject(t, item, tt.expected)
			case *object.Decimal:
				testNumberObject(t, item, tt.expected)
			case *object.Double:
				testNumberObject(t, item, tt.expected)
			case *object.String:
				testStringObject(t, item, tt.expected)
			default:
				t.Errorf("Unkown item type. got=%s", item.Type())
			}
		}
	}
}

func TestPredicate(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`(1,2,3,4)[1]`, 1},
		{`(1,2,3,4)[1+1]`, 2},
		{`(2,1,3,4)[.=2]`, 2},
	}

	for _, tt := range tests {
		seq := testEval(tt.input)
		sequence := seq.(*object.Sequence)

		for _, item := range sequence.Items {
			switch item := item.(type) {
			case *object.Integer:
				if item.Value() != tt.expected {
					t.Errorf("item has wrong value. got=%d, want=%d", item.Value(), tt.expected)
				}
			}
		}
	}
}

func TestIfExpr(t *testing.T) {
	seq := testEval("if ('a') then 2 else 3")
	sequence := seq.(*object.Sequence)
	if len(sequence.Items) != 1 {
		t.Errorf("wrong number of items. got=%d, expected=1", len(sequence.Items))
	}
	item, ok := sequence.Items[0].(*object.Integer)
	if !ok {
		t.Errorf("item type should be integer")
	}
	if item.Value() != 2 {
		t.Errorf("item value should be 2. got=%d", item.Value())
	}

	seq2 := testEvalXML2("if (/company/office[@location='Seoul']/employee[1]/age = 35) then 'is 35' else 'is not 35'")
	sequence2 := seq2.(*object.Sequence)
	if len(sequence2.Items) != 1 {
		t.Errorf("wrong number of items. got=%d, expected=1", len(sequence2.Items))
	}
	item2, ok := sequence2.Items[0].(*object.String)
	if !ok {
		t.Errorf("item type should be string")
	}
	if item2.Value() != "is not 35" {
		t.Errorf("item value should be 'is not 35'")
	}
}

func TestForExpr(t *testing.T) {
	tests := []struct {
		input    string
		expected []interface{}
	}{
		{
			`for $i in (10, 20), $j in (1, 2), $k in (4, 5), $l in (9, 11) return ($i + $j + $k*$l)`,
			[]interface{}{47, 55, 56, 66, 48, 56, 57, 67, 57, 65, 66, 76, 58, 66, 67, 77},
		},
	}

	for _, tt := range tests {
		seq := testEval(tt.input)
		sequence := seq.(*object.Sequence)
		testSequenceObject(t, sequence, tt.expected)
	}

	seq2 := testEvalXML2("for $i in //company/office/employee return if ($i/age >= 30) then upper-case($i/last_name) else lower-case($i/last_name)")
	sequence2 := seq2.(*object.Sequence)
	if len(sequence2.Items) != 5 {
		t.Errorf("wrong number of items. got=%d, expected=5", len(sequence2.Items))
	}
	expects2 := []string{"jack", "HWA", "BROWN", "CHI", "JI"}
	for i, item := range sequence2.Items {
		item, ok := item.(*object.String)
		if !ok {
			t.Errorf("item type should be string")
		}
		if item.Value() != expects2[i] {
			t.Errorf("got=%s, expected: %s", item.Value(), expects2[i])
		}
	}

	seq3 := testEvalXML2("for $i in //company/office/employee return if ($i/age >= 32) then upper-case($i/last_name) else lower-case($i/last_name)")
	sequence3 := seq3.(*object.Sequence)
	if len(sequence3.Items) != 5 {
		t.Errorf("wrong number of items. got=%d, expected=5", len(sequence3.Items))
	}
	expects3 := []string{"jack", "hwa", "brown", "CHI", "JI"}
	for i, item := range sequence3.Items {
		item, ok := item.(*object.String)
		if !ok {
			t.Errorf("item type should be string")
		}
		if item.Value() != expects3[i] {
			t.Errorf("got=%s, expected: %s", item.Value(), expects3[i])
		}
	}
}

func TestLetExpr(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			`let $r := 5, $pi := 3.14 return  'area = ' || $pi * ($r * $r)`,
			`area = 78.5`,
		},
		{
			`let $pi := 3.14,
				$area := function ($arg)
				{
					'area = ' ||	$pi * $arg * $arg
				},
				$r := 5
				return $area($r)`,
			`area = 78.5`,
		},
	}

	for _, tt := range tests {
		seq := testEval(tt.input)
		sequence := seq.(*object.Sequence)
		item := sequence.Items[0].(*object.String)

		if item.Value() != tt.expected {
			t.Errorf("got=%s, expected=%s", item.Value(), tt.expected)
		}
	}
}

func TestMapExpr(t *testing.T) {
	tests := []struct {
		input    string
		expected []interface{}
	}{
		{`map{"a":1}?a`, []interface{}{1}},
		{`map{"a":1,"b":2,"c":3}?("a","b")`, []interface{}{1, 2}},
	}

	for _, tt := range tests {
		seq := testEval(tt.input)
		sequence := seq.(*object.Sequence)
		testSequenceObject(t, sequence, tt.expected)
	}
}

func TestQuantifiedExpr(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{`some $i in (1,2,3), $j in (5,6,7,3) satisfies $i = $j`, true},
		{`every $i in (1,2,3), $j in (5,6,7,3) satisfies $i = $j`, false},
		{`some $x in (1, 2, 3), $y in (2, 3, 4) satisfies $x + $y = 4`, true},
		{`every $x in (1, 2, 3), $y in (2, 3, 4) satisfies $x + $y = 4`, false},
	}

	for _, tt := range tests {
		seq := testEval(tt.input)
		sequence := seq.(*object.Sequence)
		for _, item := range sequence.Items {
			bl := item.(*object.Boolean)
			if bl.Value() != tt.expected {
				t.Errorf("got: %v, expected: %v", bl.Value(), tt.expected)
			}
		}
	}
}

func TestDocNode(t *testing.T) {
	tests := []string{
		"/",
		"/document-node()",
		"//document-node()",
	}

	for _, tt := range tests {
		seq := testEval(tt)
		sequence := seq.(*object.Sequence)

		if len(sequence.Items) != 1 {
			t.Errorf("wrong number of seq items. got=%d, expected=1", len(sequence.Items))
		}

		if sequence.Items[0].Type() != object.DocumentNodeType {
			t.Errorf("node is not a document type. got=%s", sequence.Items[0].Type())
		}
	}

	seq := testEval("//")
	sequence := seq.(*object.Sequence)
	if !bif.IsError(sequence.Items[0]) {
		t.Errorf("// is not a valid xpath expression")
	}

	seq2 := testEval("/div")
	sequence2 := seq2.(*object.Sequence)
	if sequence2.Items != nil {
		t.Errorf("the result should be nil")
	}
}

func TestPathExpr(t *testing.T) {
	seq := testEval("//title")
	sequence := seq.(*object.Sequence)
	if len(sequence.Items) != 1 {
		t.Errorf("wrong number of items. got=%d, expected=1", len(sequence.Items))
	}

	seq2 := testEval("//title/text()")
	sequence2 := seq2.(*object.Sequence)
	node2 := sequence2.Items[0].(*object.BaseNode)
	if node2.Type() != object.TextNodeType {
		t.Errorf("wrong node type. got=%s, expected=TextNodeType", node2.Type())
	}
	if node2.Tree().Data != "Quotes to Scrape" {
		t.Errorf("wrong text value. got=%s, expected='Quotes to Scrape'", node2.Tree().Data)
	}

	seq3 := testEval("//div")
	sequence3 := seq3.(*object.Sequence)
	if len(sequence3.Items) != 28 {
		t.Errorf("wrong number of items. got=%d, expected=28", len(sequence3.Items))
	}

	seq4 := testEval("//span/small")
	sequence4 := seq4.(*object.Sequence)
	if len(sequence4.Items) != 10 {
		t.Errorf("wrong number of items. got=%d, expected=10", len(sequence4.Items))
	}

	seq5 := testEvalXML("//book")
	sequence5 := seq5.(*object.Sequence)
	if len(sequence5.Items) != 5 {
		t.Errorf("wrong number of items. got=%d, expected=5", len(sequence5.Items))
	}

	seq6 := testEvalXML("//book/month")
	sequence6 := seq6.(*object.Sequence)
	if sequence6.Items != nil {
		t.Errorf("//book/month expr doesn't have items")
	}

	seq7 := testEvalXML("//book//month")
	sequence7 := seq7.(*object.Sequence)
	if len(sequence7.Items) != 1 {
		t.Errorf("//book//month must select one item")
	}

	seq8 := testEvalXML("//book/year")
	sequence8 := seq8.(*object.Sequence)
	if len(sequence8.Items) != 3 {
		t.Errorf("wrong number of items. got=%d, expected=3", len(sequence8.Items))
	}

	seq9 := testEvalXML("//book/haha")
	sequence9 := seq9.(*object.Sequence)
	if len(sequence9.Items) != 1 {
		t.Errorf("wrong number of items. got=%d, expected=1", len(sequence9.Items))
	}

	seq10 := testEvalXML("//book/haha/parent::book")
	sequence10 := seq10.(*object.Sequence)
	if len(sequence10.Items) != 1 {
		t.Errorf("wrong number of items. got=%d, expected=1", len(sequence10.Items))
	}

	seq11 := testEvalXML("//book/haha/parent::book/text()")
	sequence11 := seq11.(*object.Sequence)
	node11 := sequence11.Items[0].(*object.BaseNode)
	if node11.Type() != object.TextNodeType {
		t.Errorf("wrong node type. got=%s, expected=TextNodeType", node11.Type())
	}
	if strings.TrimSpace(node11.Tree().Data) != "hterkj" {
		t.Errorf("wrong text value. got=%s, expected='hterkj'", node11.Tree().Data)
	}

	seq12 := testEvalXML("//book/haha/ancestor::book")
	sequence12 := seq12.(*object.Sequence)
	if len(sequence12.Items) != 3 {
		t.Errorf("wrong number of items. got=%d, expected=3", len(sequence12.Items))
	}

	seq13 := testEvalXML("//book/ancestor::book")
	sequence13 := seq13.(*object.Sequence)
	if len(sequence13.Items) != 2 {
		t.Errorf("wrong number of items. got=%d, expected=2", len(sequence13.Items))
	}

	seq14 := testEvalXML("//book/ancestor-or-self::book")
	sequence14 := seq14.(*object.Sequence)
	if len(sequence14.Items) != 5 {
		t.Errorf("wrong number of items. got=%d, expected=5", len(sequence14.Items))
	}

	seq15 := testEvalXML("//book/preceding-sibling::book")
	sequence15 := seq15.(*object.Sequence)
	if len(sequence15.Items) != 2 {
		t.Errorf("wrong number of items. got=%d, expected=2", len(sequence15.Items))
	}

	seq16 := testEvalXML("//book/following-sibling::book")
	sequence16 := seq16.(*object.Sequence)
	if len(sequence16.Items) != 2 {
		t.Errorf("wrong number of items. got=%d, expected=2", len(sequence16.Items))
	}

	seq17 := testEvalXML("//month/following::year")
	sequence17 := seq17.(*object.Sequence)
	if len(sequence17.Items) != 2 {
		t.Errorf("wrong number of items. got=%d, expected=2", len(sequence17.Items))
	}

	seq18 := testEvalXML("//month/following-sibling::year")
	sequence18 := seq18.(*object.Sequence)
	if sequence18.Items != nil {
		t.Errorf("the result should be nil")
	}

	seq19 := testEvalXML("//author/following-sibling::year")
	sequence19 := seq19.(*object.Sequence)
	if len(sequence19.Items) != 2 {
		t.Errorf("wrong number of items. got=%d, expected=2", len(sequence19.Items))
	}

	seq20 := testEvalXML("//go/following::year")
	sequence20 := seq20.(*object.Sequence)
	if len(sequence20.Items) != 4 {
		t.Errorf("wrong number of items. got=%d, expected=4", len(sequence20.Items))
	}

	seq21 := testEvalXML("//go/following-sibling::year")
	sequence21 := seq21.(*object.Sequence)
	if sequence21.Items != nil {
		t.Errorf("the result should be nil")
	}

	seq22 := testEvalXML("//month/preceding::year")
	sequence22 := seq22.(*object.Sequence)
	if len(sequence22.Items) != 1 {
		t.Errorf("wrong number of items. got=%d, expected=1", len(sequence22.Items))
	}

	seq23 := testEvalXML("//month/preceding-sibling::year")
	sequence23 := seq23.(*object.Sequence)
	if sequence23.Items != nil {
		t.Errorf("the result should be nil")
	}

	seq24 := testEvalXML("//year/preceding-sibling::book")
	sequence24 := seq24.(*object.Sequence)
	if len(sequence24.Items) != 1 {
		t.Errorf("wrong number of items. got=%d, expected=1", len(sequence24.Items))
	}

	seq25 := testEvalXML("//year/preceding-sibling::book/text()")
	sequence25 := seq25.(*object.Sequence)
	if len(sequence25.Items) != 2 {
		t.Errorf("wrong number of items. got=%d, expected=2", len(sequence25.Items))
	}
	node25 := sequence25.Items[0].(*object.BaseNode)
	if strings.TrimSpace(node25.Tree().Data) != "godc" {
		t.Errorf("wrong text value. got=%s, expected='godc'", node25.Tree().Data)
	}

	seq26 := testEvalXML("//year/ancestor-or-self::book")
	sequence26 := seq26.(*object.Sequence)
	if len(sequence26.Items) != 5 {
		t.Errorf("wrong number of items. got=%d, expected=5", len(sequence26.Items))
	}

	seq27 := testEvalXML("//month/ancestor-or-self::book")
	sequence27 := seq27.(*object.Sequence)
	if len(sequence27.Items) != 1 {
		t.Errorf("wrong number of items. got=%d, expected=1", len(sequence27.Items))
	}

	seq28 := testEvalXML("//book/descendant::book")
	sequence28 := seq28.(*object.Sequence)
	if len(sequence28.Items) != 2 {
		t.Errorf("wrong number of items. got=%d, expected=2", len(sequence28.Items))
	}

	seq29 := testEvalXML("//book/descendant-or-self::book")
	sequence29 := seq29.(*object.Sequence)
	if len(sequence29.Items) != 5 {
		t.Errorf("wrong number of items. got=%d, expected=5", len(sequence29.Items))
	}

	seq30 := testEvalXML("//book/*")
	sequence30 := seq30.(*object.Sequence)
	if len(sequence30.Items) != 18 {
		t.Errorf("wrong number of items. got=%d, expected=18", len(sequence30.Items))
	}

	seq31 := testEvalXML("//book/haha/year/text()")
	sequence31 := seq31.(*object.Sequence)
	node31 := sequence31.Items[0].(*object.BaseNode)
	if node31.Tree().Data != "001" {
		t.Errorf("wrong result. got=%s, expected='001'", node31.Tree().Data)
	}

	seq32 := testEvalXML("//tt:book")
	sequence32 := seq32.(*object.Sequence)
	if len(sequence32.Items) != 1 {
		t.Errorf("wrong number of items. got=%d, expected=1", len(sequence32.Items))
	}

	seq33 := testEvalXML("/html/body/tt:bookstore")
	sequence33 := seq33.(*object.Sequence)
	if len(sequence33.Items) != 1 {
		t.Errorf("wrong number of items. got=%d, expected=1", len(sequence33.Items))
	}

	seq34 := testEvalXML("//@category")
	sequence34 := seq34.(*object.Sequence)
	if len(sequence34.Items) != 7 {
		t.Errorf("wrong number of items. got=%d, expected=7", len(sequence34.Items))
	}

	seq35 := testEvalXML("//book/@category")
	sequence35 := seq35.(*object.Sequence)
	if len(sequence35.Items) != 5 {
		t.Errorf("wrong number of items. got=%d, expected=5", len(sequence35.Items))
	}

	seq36 := testEvalXML("//attribute::lang")
	sequence36 := seq36.(*object.Sequence)
	if len(sequence36.Items) != 4 {
		t.Errorf("wrong number of items. got=%d, expected=4", len(sequence36.Items))
	}

	seq37 := testEvalXML("//book/title/attribute::lang")
	sequence37 := seq37.(*object.Sequence)
	if len(sequence37.Items) != 3 {
		t.Errorf("wrong number of items. got=%d, expected=3", len(sequence37.Items))
	}

	seq38 := testEvalXML("//book/attribute::cover")
	sequence38 := seq38.(*object.Sequence)
	if len(sequence38.Items) != 1 {
		t.Errorf("wrong number of items. got=%d, expected=1", len(sequence38.Items))
	}
	node38 := sequence38.Items[0].(*object.AttrNode)
	if node38.Type() != object.AttributeNodeType {
		t.Errorf("wrong node type. got=%s", node38.Type())
	}
	if node38.Inspect() != "paperback" {
		t.Errorf("wrong attr value. got=%s, expected='paperback'", node38.Inspect())
	}

	seq39 := testEvalXML("//book/self::book")
	sequence39 := seq39.(*object.Sequence)
	if len(sequence39.Items) != 5 {
		t.Errorf("wrong number of items. got=%d, expected=5", len(sequence39.Items))
	}

	seq40 := testEvalXML("//book/self::year")
	sequence40 := seq40.(*object.Sequence)
	if sequence40.Items != nil {
		t.Errorf("the result should be nil")
	}

	seq41 := testEvalXML("/*")
	sequence41 := seq41.(*object.Sequence)
	if len(sequence41.Items) != 1 {
		t.Errorf("wrong number of items. got=%d, expected=1", len(sequence41.Items))
	}

	seq42 := testEvalXML("//*")
	sequence42 := seq42.(*object.Sequence)
	if len(sequence42.Items) != 37 {
		t.Errorf("wrong number of items. got=%d, expected=37", len(sequence42.Items))
	}

	seq43 := testEvalXML("/child::*")
	sequence43 := seq43.(*object.Sequence)
	if len(sequence43.Items) != 1 {
		t.Errorf("wrong number of items. got=%d, expected=1", len(sequence43.Items))
	}

	seq44 := testEvalXML("//haha/descendant-or-self::*")
	sequence44 := seq44.(*object.Sequence)
	if len(sequence44.Items) != 3 {
		t.Errorf("wrong number of items. got=%d, expected=3", len(sequence44.Items))
	}

	seq45 := testEvalXML("//book/descendant::*")
	sequence45 := seq45.(*object.Sequence)
	if len(sequence45.Items) != 23 {
		t.Errorf("wrong number of items. got=%d, expected=23", len(sequence45.Items))
	}

	seq46 := testEvalXML("//book/attribute::*")
	sequence46 := seq46.(*object.Sequence)
	if len(sequence46.Items) != 7 {
		t.Errorf("wrong number of items. got=%d, expected=7", len(sequence46.Items))
	}

	seq47 := testEvalXML("//book/self::*")
	sequence47 := seq47.(*object.Sequence)
	if len(sequence47.Items) != 5 {
		t.Errorf("wrong number of items. got=%d, expected=5", len(sequence47.Items))
	}

	seq48 := testEvalXML("//book/descendant-or-self::*")
	sequence48 := seq48.(*object.Sequence)
	if len(sequence48.Items) != 26 {
		t.Errorf("wrong number of items. got=%d, expected=26", len(sequence48.Items))
	}

	seq49 := testEvalXML("//book/following-sibling::*")
	sequence49 := seq49.(*object.Sequence)
	if len(sequence49.Items) != 4 {
		t.Errorf("wrong number of items. got=%d, expected=4", len(sequence49.Items))
	}

	seq50 := testEvalXML("//book/following::*")
	sequence50 := seq50.(*object.Sequence)
	if len(sequence50.Items) != 20 {
		t.Errorf("wrong number of items. got=%d, expected=20", len(sequence50.Items))
	}

	seq51 := testEvalXML("//book/parent::*")
	sequence51 := seq51.(*object.Sequence)
	if len(sequence51.Items) != 3 {
		t.Errorf("wrong number of items. got=%d, expected=3", len(sequence51.Items))
	}

	seq52 := testEvalXML("//no/ancestor::*")
	sequence52 := seq52.(*object.Sequence)
	if len(sequence52.Items) != 6 {
		t.Errorf("wrong number of items. got=%d, expected=6", len(sequence52.Items))
	}

	seq53 := testEvalXML("//no/preceding-sibling::*")
	sequence53 := seq53.(*object.Sequence)
	if len(sequence53.Items) != 2 {
		t.Errorf("wrong number of items. got=%d, expected=2", len(sequence53.Items))
	}

	seq54 := testEvalXML("//no/preceding::*")
	sequence54 := seq54.(*object.Sequence)
	if len(sequence54.Items) != 15 {
		t.Errorf("wrong number of items. got=%d, expected=15", len(sequence54.Items))
	}

	seq55 := testEvalXML("//no/ancestor-or-self::*")
	sequence55 := seq55.(*object.Sequence)
	if len(sequence55.Items) != 7 {
		t.Errorf("wrong number of items. got=%d, expected=7", len(sequence55.Items))
	}

	seq56 := testEvalXML("//book/author")
	sequence56 := seq56.(*object.Sequence)
	if len(sequence56.Items) != 6 {
		t.Errorf("wrong number of items. got=%d, expected=6", len(sequence56.Items))
	}

	seq57 := testEvalXML("//book/(title|author)")
	sequence57 := seq57.(*object.Sequence)
	if len(sequence57.Items) != 9 {
		t.Errorf("wrong number of items. got=%d, expected=9", len(sequence57.Items))
	}

	seq58 := testEvalXML("//book/(title|author)/title")
	sequence58 := seq58.(*object.Sequence)
	if sequence58.Items != nil {
		t.Errorf("the result should be nil")
	}

	seq59 := testEvalXML("//book/(title|author)/self::author")
	sequence59 := seq59.(*object.Sequence)
	if len(sequence59.Items) != 6 {
		t.Errorf("wrong number of items. got=%d, expected=6", len(sequence59.Items))
	}

	seq60 := testEvalXML("//book/(title|author)//*")
	sequence60 := seq60.(*object.Sequence)
	if sequence60.Items != nil {
		t.Errorf("the result should be nil")
	}

	seq61 := testEvalXML("book")
	sequence61 := seq61.(*object.Sequence)
	if sequence61.Items != nil {
		t.Errorf("the result should be nil")
	}

	seq62 := testEvalXML("html")
	sequence62 := seq62.(*object.Sequence)
	if len(sequence62.Items) != 1 {
		t.Errorf("wrong number of items. got=%d, expected=1", len(sequence62.Items))
	}
}

func TestPathPredicateExpr(t *testing.T) {
	seq := testEvalXML("//book[0]")
	sequence := seq.(*object.Sequence)
	if sequence.Items != nil {
		t.Errorf("the result should be nil")
	}

	seq2 := testEvalXML("//book[1]")
	sequence2 := seq2.(*object.Sequence)
	if len(sequence2.Items) != 3 {
		t.Errorf("wrong number of items. got=%d, expected=3", len(sequence2.Items))
	}

	seq3 := testEvalXML("//book[2]")
	sequence3 := seq3.(*object.Sequence)
	if len(sequence3.Items) != 1 {
		t.Errorf("wrong number of items. got=%d, expected=1", len(sequence3.Items))
	}

	seq4 := testEvalXML("//book[2]/@category")
	sequence4 := seq4.(*object.Sequence)
	if len(sequence4.Items) != 1 {
		t.Errorf("wrong number of items. got=%d, expected=1", len(sequence4.Items))
	}
	node4, ok := sequence4.Items[0].(*object.AttrNode)
	if !ok {
		t.Errorf("node type must be an attribute node")
	}
	if node4.Inspect() != "web" {
		t.Errorf("wrong attribute value. got=%s, expected='web'", node4.Inspect())
	}

	seq5 := testEvalXML("//book[@category='web']")
	sequence5 := seq5.(*object.Sequence)
	if len(sequence5.Items) != 2 {
		t.Errorf("wrong number of items. got=%d, expected=2", len(sequence5.Items))
	}

	seq6 := testEvalXML("//book['web'=@category]")
	sequence6 := seq6.(*object.Sequence)
	if len(sequence6.Items) != 2 {
		t.Errorf("wrong number of items. got=%d, expected=2", len(sequence6.Items))
	}

	seq7 := testEvalXML("//book['web'=@category][@cover='paperback']")
	sequence7 := seq7.(*object.Sequence)
	if len(sequence7.Items) != 1 {
		t.Errorf("wrong number of items. got=%d, expected=1", len(sequence7.Items))
	}

	seq8 := testEvalXML("//*[title]")
	sequence8 := seq8.(*object.Sequence)
	if len(sequence8.Items) != 3 {
		t.Errorf("wrong number of items. got=%d, expected=3", len(sequence8.Items))
	}

	seq9 := testEvalXML("//child::*[title]")
	sequence9 := seq9.(*object.Sequence)
	if len(sequence9.Items) != 3 {
		t.Errorf("wrong number of items. got=%d, expected=3", len(sequence9.Items))
	}

	seq10 := testEvalXML("//child::book[author]")
	sequence10 := seq10.(*object.Sequence)
	if len(sequence10.Items) != 2 {
		t.Errorf("wrong number of items. got=%d, expected=2", len(sequence10.Items))
	}

	seq11 := testEvalXML("//child::book[title]")
	sequence11 := seq11.(*object.Sequence)
	if len(sequence11.Items) != 3 {
		t.Errorf("wrong number of items. got=%d, expected=3", len(sequence11.Items))
	}

	seq12 := testEvalXML("//descendant::title[.='Harry Potter']")
	sequence12 := seq12.(*object.Sequence)
	if len(sequence12.Items) != 1 {
		t.Errorf("wrong number of items. got=%d, expected=1", len(sequence12.Items))
	}

	seq13 := testEvalXML("//attribute::lang['en'=.]")
	sequence13 := seq13.(*object.Sequence)
	if len(sequence13.Items) != 2 {
		t.Errorf("wrong number of items. got=%d, expected=2", len(sequence13.Items))
	}

	seq14 := testEvalXML("//attribute::*[.='cooking']")
	sequence14 := seq14.(*object.Sequence)
	if len(sequence14.Items) != 1 {
		t.Errorf("wrong number of items. got=%d, expected=1", len(sequence14.Items))
	}

	seq15 := testEvalXML("//book/self::*[@cover]")
	sequence15 := seq15.(*object.Sequence)
	if len(sequence15.Items) != 1 {
		t.Errorf("wrong number of items. got=%d, expected=1", len(sequence15.Items))
	}

	seq16 := testEvalXML("//book/self::book[@cover]")
	sequence16 := seq16.(*object.Sequence)
	if len(sequence16.Items) != 1 {
		t.Errorf("wrong number of items. got=%d, expected=1", len(sequence16.Items))
	}

	seq17 := testEvalXML("//book/descendant-or-self::author[.='James McGovern']")
	sequence17 := seq17.(*object.Sequence)
	if len(sequence17.Items) != 1 {
		t.Errorf("wrong number of items. got=%d, expected=1", len(sequence17.Items))
	}

	seq18 := testEvalXML("//book/descendant-or-self::author[.]")
	sequence18 := seq18.(*object.Sequence)
	if len(sequence18.Items) != 6 {
		t.Errorf("wrong number of items. got=%d, expected=6", len(sequence18.Items))
	}

	seq19 := testEvalXML("//book/descendant-or-self::*[child::price]")
	sequence19 := seq19.(*object.Sequence)
	if len(sequence19.Items) != 3 {
		t.Errorf("wrong number of items. got=%d, expected=3", len(sequence19.Items))
	}

	seq20 := testEvalXML("//book/descendant-or-self::*[price]")
	sequence20 := seq20.(*object.Sequence)
	if len(sequence20.Items) != 3 {
		t.Errorf("wrong number of items. got=%d, expected=3", len(sequence20.Items))
	}

	seq21 := testEvalXML("//book/following-sibling::*[year]")
	sequence21 := seq21.(*object.Sequence)
	if len(sequence21.Items) != 2 {
		t.Errorf("wrong number of items. got=%d, expected=2", len(sequence21.Items))
	}

	seq22 := testEvalXML("//book/following-sibling::year[.]")
	sequence22 := seq22.(*object.Sequence)
	if len(sequence22.Items) != 1 {
		t.Errorf("wrong number of items. got=%d, expected=1", len(sequence22.Items))
	}

	seq23 := testEvalXML("//book/following::*[.]")
	sequence23 := seq23.(*object.Sequence)
	if len(sequence23.Items) != 20 {
		t.Errorf("wrong number of items. got=%d, expected=20", len(sequence23.Items))
	}

	seq24 := testEvalXML("//book/following::*[year]")
	sequence24 := seq24.(*object.Sequence)
	if len(sequence24.Items) != 2 {
		t.Errorf("wrong number of items. got=%d, expected=2", len(sequence24.Items))
	}

	seq25 := testEvalXML("//book/following::year[.]")
	sequence25 := seq25.(*object.Sequence)
	if len(sequence25.Items) != 3 {
		t.Errorf("wrong number of items. got=%d, expected=3", len(sequence25.Items))
	}

	seq26 := testEvalXML("//year/parent::*[book]")
	sequence26 := seq26.(*object.Sequence)
	if len(sequence26.Items) != 1 {
		t.Errorf("wrong number of items. got=%d, expected=1", len(sequence26.Items))
	}

	seq27 := testEvalXML("//year/parent::*[.]")
	sequence27 := seq27.(*object.Sequence)
	if len(sequence27.Items) != 4 {
		t.Errorf("wrong number of items. got=%d, expected=4", len(sequence27.Items))
	}

	seq28 := testEvalXML("//year/parent::book[@category='1']")
	sequence28 := seq28.(*object.Sequence)
	if len(sequence28.Items) != 1 {
		t.Errorf("wrong number of items. got=%d, expected=1", len(sequence28.Items))
	}
	node28 := sequence28.Items[0].(*object.BaseNode)
	if node28.Parent().Tree().Data != "tt:bookstore" {
		t.Errorf("parent node tag name must be tt:bookstore, got=%s", node28.Parent().Tree().Data)
	}

	seq29 := testEvalXML("//year/ancestor::*[book]")
	sequence29 := seq29.(*object.Sequence)
	if len(sequence29.Items) != 3 {
		t.Errorf("wrong number of items. got=%d, expected=3", len(sequence29.Items))
	}

	seq30 := testEvalXML("//year/ancestor::*[.]")
	sequence30 := seq30.(*object.Sequence)
	if len(sequence30.Items) != 9 {
		t.Errorf("wrong number of items. got=%d, expected=9", len(sequence30.Items))
	}

	seq31 := testEvalXML("//year/ancestor::book[1]")
	sequence31 := seq31.(*object.Sequence)
	if len(sequence31.Items) != 4 {
		t.Errorf("wrong number of items. got=%d, expected=4", len(sequence31.Items))
	}

	seq32 := testEvalXML("//year/ancestor::book[2]")
	sequence32 := seq32.(*object.Sequence)
	if len(sequence32.Items) != 1 {
		t.Errorf("wrong number of items. got=%d, expected=1", len(sequence32.Items))
	}

	seq33 := testEvalXML("//no/preceding-sibling::*[.]")
	sequence33 := seq33.(*object.Sequence)
	if len(sequence33.Items) != 2 {
		t.Errorf("wrong number of items. got=%d, expected=2", len(sequence33.Items))
	}

	seq34 := testEvalXML("//no/preceding-sibling::*[haha]")
	sequence34 := seq34.(*object.Sequence)
	if sequence34.Items != nil {
		t.Errorf("the result should be nil")
	}

	seq35 := testEvalXML("//no/preceding-sibling::*[preceding-sibling::*]")
	sequence35 := seq35.(*object.Sequence)
	if len(sequence35.Items) != 1 {
		t.Errorf("wrong number of items. got=%d, expected=1", len(sequence35.Items))
	}

	seq36 := testEvalXML("//no/preceding-sibling::day[preceding-sibling::*]")
	sequence36 := seq36.(*object.Sequence)
	if len(sequence36.Items) != 1 {
		t.Errorf("wrong number of items. got=%d, expected=1", len(sequence36.Items))
	}

	seq37 := testEvalXML("//no/preceding::*[.]")
	sequence37 := seq37.(*object.Sequence)
	if len(sequence37.Items) != 15 {
		t.Errorf("wrong number of items. got=%d, expected=15", len(sequence37.Items))
	}

	seq38 := testEvalXML("//no/preceding::*[preceding::*]")
	sequence38 := seq38.(*object.Sequence)
	if len(sequence38.Items) != 14 {
		t.Errorf("wrong number of items. got=%d, expected=14", len(sequence38.Items))
	}

	seq39 := testEvalXML("//no/preceding::year[.]")
	sequence39 := seq39.(*object.Sequence)
	if len(sequence39.Items) != 1 {
		t.Errorf("wrong number of items. got=%d, expected=1", len(sequence39.Items))
	}

	seq40 := testEvalXML("//no/ancestor-or-self::*[.]")
	sequence40 := seq40.(*object.Sequence)
	if len(sequence40.Items) != 7 {
		t.Errorf("wrong number of items. got=%d, expected=7", len(sequence40.Items))
	}

	seq41 := testEvalXML("//no/ancestor-or-self::*[no]")
	sequence41 := seq41.(*object.Sequence)
	if len(sequence41.Items) != 1 {
		t.Errorf("wrong number of items. got=%d, expected=1", len(sequence41.Items))
	}
	node41 := sequence41.Items[0].(*object.BaseNode)
	if node41.Tree().Data != "month" {
		t.Errorf("selected tag name must be a [month]")
	}

	seq42 := testEvalXML("//no/ancestor-or-self::no[.]")
	sequence42 := seq42.(*object.Sequence)
	if len(sequence42.Items) != 1 {
		t.Errorf("wrong number of items. got=%d, expected=1", len(sequence42.Items))
	}
	node42 := sequence42.Items[0].(*object.BaseNode)
	if node42.Tree().Data != "no" {
		t.Errorf("selected tag name must be a [no]")
	}

	seq43 := testEvalXML("//attribute::category[. = '1']")
	sequence43 := seq43.(*object.Sequence)
	if len(sequence43.Items) != 1 {
		t.Errorf("wrong number of items. got=%d, expected=1", len(sequence43.Items))
	}
	node43, ok := sequence43.Items[0].(*object.AttrNode)
	if !ok {
		t.Errorf("node must be an AttrNode. got=%s", node43.Type())
	}
	if node43.Inspect() != "1" {
		t.Errorf("wrong attribute value. got=%s, expected='1'", node43.Inspect())
	}

	seq44 := testEvalXML("//@category[. = '2']")
	sequence44 := seq44.(*object.Sequence)
	if len(sequence44.Items) != 1 {
		t.Errorf("wrong number of items. got=%d, expected=1", len(sequence44.Items))
	}
	node44, ok := sequence44.Items[0].(*object.AttrNode)
	if !ok {
		t.Errorf("node must be an AttrNode. got=%s", node44.Type())
	}
	if node44.Inspect() != "2" {
		t.Errorf("wrong attribute value. got=%s, expected='2'", node44.Inspect())
	}

	seq45 := testEvalXML("//book[year='2003'][1]")
	sequence45 := seq45.(*object.Sequence)
	if len(sequence45.Items) != 1 {
		t.Errorf("wrong number of items. got=%d, expected=1", len(sequence45.Items))
	}

	seq46 := testEvalXML("//book[year='2003'][1][1]")
	sequence46 := seq46.(*object.Sequence)
	if len(sequence46.Items) != 1 {
		t.Errorf("wrong number of items. got=%d, expected=1", len(sequence46.Items))
	}

	seq47 := testEvalXML("//book[year='2003'][1][1][1]")
	sequence47 := seq47.(*object.Sequence)
	if len(sequence47.Items) != 1 {
		t.Errorf("wrong number of items. got=%d, expected=1", len(sequence47.Items))
	}

	seq48 := testEvalXML("//book[year='2003']/preceding::*[book]")
	sequence48 := seq48.(*object.Sequence)
	if len(sequence48.Items) != 2 {
		t.Errorf("wrong number of items. got=%d, expected=2", len(sequence48.Items))
	}

	seq49 := testEvalXML("//book[year='2003']/preceding::*[book][1]")
	sequence49 := seq49.(*object.Sequence)
	if len(sequence49.Items) != 1 {
		t.Errorf("wrong number of items. got=%d, expected=1", len(sequence49.Items))
	}
	node49 := sequence49.Items[0].(*object.BaseNode)
	if len(node49.Attr()) != 1 {
		t.Errorf("wrong number of attrs. got=%d, expected=1", len(node49.Attr()))
	}
	attr49, ok := node49.Attr()[0].(*object.AttrNode)
	if !ok {
		t.Errorf("node type should be AttrNode")
	}
	if attr49.Key() != "category" || attr49.Inspect() != "2" {
		t.Errorf("expected attr: %s='%s'", attr49.Key(), attr49.Inspect())
	}

	seq50 := testEvalXML("//book[year='2003']/preceding::*[book][2]")
	sequence50 := seq50.(*object.Sequence)
	if len(sequence50.Items) != 1 {
		t.Errorf("wrong number of items. got=%d, expected=1", len(sequence50.Items))
	}
	node50 := sequence50.Items[0].(*object.BaseNode)
	if len(node50.Attr()) != 1 {
		t.Errorf("wrong number of attrs. got=%d, expected=1", len(node50.Attr()))
	}
	attr50, ok := node50.Attr()[0].(*object.AttrNode)
	if !ok {
		t.Errorf("node type should be AttrNode")
	}
	if attr50.Key() != "category" || attr50.Inspect() != "1" {
		t.Errorf("expected attr: %s='%s'", attr50.Key(), attr50.Inspect())
	}
}

func TestKindTest(t *testing.T) {
	seq := testEvalXML("//book/haha/year/text()")
	sequence := seq.(*object.Sequence)
	if len(sequence.Items) != 1 {
		t.Errorf("wrong number of items. got=%d, expected=1", len(sequence.Items))
	}
	node, ok := sequence.Items[0].(*object.BaseNode)
	if !ok {
		t.Errorf("sequence item must be a node")
	}
	if node.Text() != "001" {
		t.Errorf("wrong node value. got=%s, expected='001'", node.Text())
	}

	seq2 := testEvalXML("//book/haha/year/text()[0]")
	sequence2 := seq2.(*object.Sequence)
	if sequence2.Items != nil {
		t.Errorf("the result should be nil")
	}

	seq3 := testEvalXML("//book/haha/year/text()[1]")
	sequence3 := seq3.(*object.Sequence)
	if len(sequence3.Items) != 1 {
		t.Errorf("wrong number of items. got=%d, expected=1", len(sequence3.Items))
	}

	seq4 := testEvalXML("//book/haha/year/child::text()[1]")
	sequence4 := seq4.(*object.Sequence)
	if len(sequence4.Items) != 1 {
		t.Errorf("wrong number of items. got=%d, expected=1", len(sequence4.Items))
	}

	seq5 := testEvalXML("//book/document-node()")
	sequence5 := seq5.(*object.Sequence)
	if sequence5.Items != nil {
		t.Errorf("the result should be nil")
	}

	seq6 := testEvalXML("//document-node()")
	sequence6 := seq6.(*object.Sequence)
	if len(sequence6.Items) != 1 {
		t.Errorf("wrong number of items. got=%d, expected=1", len(sequence6.Items))
	}
	node6, ok := sequence6.Items[0].(*object.BaseNode)
	if !ok {
		t.Errorf("node should be BaseNode")
	}
	if node6.Type() != object.DocumentNodeType {
		t.Errorf("node type should be DocumentNodeType")
	}

	seq7 := testEvalXML("/document-node()")
	sequence7 := seq7.(*object.Sequence)
	if len(sequence7.Items) != 1 {
		t.Errorf("wrong number of items. got=%d, expected=1", len(sequence7.Items))
	}
	node7, ok := sequence7.Items[0].(*object.BaseNode)
	if !ok {
		t.Errorf("node should be BaseNode")
	}
	if node7.Type() != object.DocumentNodeType {
		t.Errorf("node type should be DocumentNodeType")
	}

	seq8 := testEvalXML("/document-node()[.]")
	sequence8 := seq8.(*object.Sequence)
	if len(sequence8.Items) != 1 {
		t.Errorf("wrong number of items. got=%d, expected=1", len(sequence8.Items))
	}
	node8, ok := sequence8.Items[0].(*object.BaseNode)
	if !ok {
		t.Errorf("node should be BaseNode")
	}
	if node8.Type() != object.DocumentNodeType {
		t.Errorf("node type should be DocumentNodeType")
	}

	seq9 := testEvalXML("//document-node()[1]")
	sequence9 := seq9.(*object.Sequence)
	if len(sequence9.Items) != 1 {
		t.Errorf("wrong number of items. got=%d, expected=1", len(sequence9.Items))
	}
	node9, ok := sequence9.Items[0].(*object.BaseNode)
	if !ok {
		t.Errorf("node should be BaseNode")
	}
	if node9.Type() != object.DocumentNodeType {
		t.Errorf("node type should be DocumentNodeType")
	}

	seq10 := testEvalXML("/document-node()[.][1]")
	sequence10 := seq10.(*object.Sequence)
	if len(sequence10.Items) != 1 {
		t.Errorf("wrong number of items. got=%d, expected=1", len(sequence10.Items))
	}
	node10, ok := sequence10.Items[0].(*object.BaseNode)
	if !ok {
		t.Errorf("node should be BaseNode")
	}
	if node10.Type() != object.DocumentNodeType {
		t.Errorf("node type should be DocumentNodeType")
	}

	seq11 := testEvalXML("/element()")
	sequence11 := seq11.(*object.Sequence)
	if len(sequence11.Items) != 1 {
		t.Errorf("wrong number of items. got=%d, expected=1", len(sequence11.Items))
	}
	node11, ok := sequence11.Items[0].(*object.BaseNode)
	if !ok {
		t.Errorf("node should be BaseNode")
	}
	if node11.Tree().Data != "html" {
		t.Errorf("root node should be [html]")
	}

	seq12 := testEvalXML("//element()")
	sequence12 := seq12.(*object.Sequence)
	if len(sequence12.Items) != 37 {
		t.Errorf("wrong number of items. got=%d, expected=37", len(sequence12.Items))
	}

	seq13 := testEvalXML("//element()[.]")
	sequence13 := seq13.(*object.Sequence)
	if len(sequence13.Items) != 37 {
		t.Errorf("wrong number of items. got=%d, expected=37", len(sequence13.Items))
	}

	seq14 := testEvalXML("//element()[1]")
	sequence14 := seq14.(*object.Sequence)
	if len(sequence14.Items) != 14 {
		t.Errorf("wrong number of items. got=%d, expected=14", len(sequence14.Items))
	}

	seq15 := testEvalXML("//attribute()")
	sequence15 := seq15.(*object.Sequence)
	if len(sequence15.Items) != 16 {
		t.Errorf("wrong number of items. got=%d, expected=16", len(sequence15.Items))
	}

	seq16 := testEvalXML("//attribute()[1]")
	sequence16 := seq16.(*object.Sequence)
	if len(sequence16.Items) != 12 {
		t.Errorf("wrong number of items. got=%d, expected=12", len(sequence16.Items))
	}

	seq17 := testEvalXML("//attribute()[.]")
	sequence17 := seq17.(*object.Sequence)
	if len(sequence17.Items) != 16 {
		t.Errorf("wrong number of items. got=%d, expected=16", len(sequence17.Items))
	}

	seq18 := testEvalXML("//book[@category='1']/attribute()")
	sequence18 := seq18.(*object.Sequence)
	if len(sequence18.Items) != 1 {
		t.Errorf("wrong number of items. got=%d, expected=1", len(sequence18.Items))
	}

	seq19 := testEvalXML("//attribute::*/attribute::attribute()")
	sequence19 := seq19.(*object.Sequence)
	if sequence19.Items != nil {
		t.Errorf("the result should be nil")
	}

	seq20 := testEvalXML("//tt:book/attribute()")
	sequence20 := seq20.(*object.Sequence)
	if len(sequence20.Items) != 2 {
		t.Errorf("wrong number of items. got=%d, expected=2", len(sequence20.Items))
	}

	seq21 := testEvalXML("//book[@category='web']/descendant::attribute()")
	sequence21 := seq21.(*object.Sequence)
	if sequence21.Items != nil {
		t.Errorf("the result should be nil")
	}

	seq22 := testEvalXML("//book[@category='web']/child::attribute()")
	sequence22 := seq22.(*object.Sequence)
	if len(sequence22.Items) != 4 {
		t.Errorf("wrong number of items. got=%d, expected=4", len(sequence22.Items))
	}

	seq23 := testEvalXML("//book[@category='web']/following::node()")
	sequence23 := seq23.(*object.Sequence)
	if len(sequence23.Items) != 21 {
		t.Errorf("wrong number of items. got=%d, expected=21", len(sequence23.Items))
	}

	seq24 := testEvalXML("//book[@category='web']/following::element()")
	sequence24 := seq24.(*object.Sequence)
	if len(sequence24.Items) != 5 {
		t.Errorf("wrong number of items. got=%d, expected=5", len(sequence24.Items))
	}

	seq25 := testEvalXML("//book[@category='web']/following::comment()")
	sequence25 := seq25.(*object.Sequence)
	if len(sequence25.Items) != 2 {
		t.Errorf("wrong number of items. got=%d, expected=2", len(sequence25.Items))
	}

	seq26 := testEvalXML("//book[@category='web']/following::text()")
	sequence26 := seq26.(*object.Sequence)
	if len(sequence26.Items) != 14 {
		t.Errorf("wrong number of items. got=%d, expected=14", len(sequence26.Items))
	}

	seq27 := testEvalXML("//title/following-sibling::element()")
	sequence27 := seq27.(*object.Sequence)
	if len(sequence27.Items) != 13 {
		t.Errorf("wrong number of items. got=%d, expected=13", len(sequence27.Items))
	}

	seq28 := testEvalXML("//title/following-sibling::attribute()")
	sequence28 := seq28.(*object.Sequence)
	if sequence28.Items != nil {
		t.Errorf("the result should be nil")
	}

	seq29 := testEvalXML("//book[@category='web']/child::text()")
	sequence29 := seq29.(*object.Sequence)
	if len(sequence29.Items) != 14 {
		t.Errorf("wrong number of items. got=%d, expected=14", len(sequence29.Items))
	}

	seq30 := testEvalXML("//book[@category='web']/child::element()")
	sequence30 := seq30.(*object.Sequence)
	if len(sequence30.Items) != 12 {
		t.Errorf("wrong number of items. got=%d, expected=12", len(sequence30.Items))
	}

	seq31 := testEvalXML("//book[@category='web']/self::element()")
	sequence31 := seq31.(*object.Sequence)
	if len(sequence31.Items) != 2 {
		t.Errorf("wrong number of items. got=%d, expected=2", len(sequence31.Items))
	}

	seq32 := testEvalXML("//book[@category='web']/self::text()")
	sequence32 := seq32.(*object.Sequence)
	if sequence32.Items != nil {
		t.Errorf("the result should be nil")
	}

	seq33 := testEvalXML("//book[@category='web']/descendant::text()")
	sequence33 := seq33.(*object.Sequence)
	if len(sequence33.Items) != 26 {
		t.Errorf("wrong number of items. got=%d, expected=26", len(sequence33.Items))
	}

	seq34 := testEvalXML("//book[@category='web']/descendant-or-self::element()")
	sequence34 := seq34.(*object.Sequence)
	if len(sequence34.Items) != 14 {
		t.Errorf("wrong number of items. got=%d, expected=14", len(sequence34.Items))
	}

	seq35 := testEvalXML("//book[@category='web']/descendant-or-self::text()")
	sequence35 := seq35.(*object.Sequence)
	if len(sequence35.Items) != 26 {
		t.Errorf("wrong number of items. got=%d, expected=26", len(sequence35.Items))
	}

	seq36 := testEvalXML("//book[@category='web']/parent::text()")
	sequence36 := seq36.(*object.Sequence)
	if sequence36.Items != nil {
		t.Errorf("the result should be nil")
	}

	seq37 := testEvalXML("//book[@category='web']/ancestor::text()")
	sequence37 := seq37.(*object.Sequence)
	if sequence37.Items != nil {
		t.Errorf("the result should be nil")
	}

	seq38 := testEvalXML("//book[@category='web']/parent::element()")
	sequence38 := seq38.(*object.Sequence)
	if len(sequence38.Items) != 1 {
		t.Errorf("wrong number of items. got=%d, expected=1", len(sequence38.Items))
	}

	seq39 := testEvalXML("//book[@category='web']/ancestor::element()")
	sequence39 := seq39.(*object.Sequence)
	if len(sequence39.Items) != 3 {
		t.Errorf("wrong number of items. got=%d, expected=3", len(sequence39.Items))
	}

	seq40 := testEvalXML("//book[@category='web']/preceding-sibling::element()")
	sequence40 := seq40.(*object.Sequence)
	if len(sequence40.Items) != 4 {
		t.Errorf("wrong number of items. got=%d, expected=4", len(sequence40.Items))
	}

	seq41 := testEvalXML("//book[@category='web']/preceding::element()")
	sequence41 := seq41.(*object.Sequence)
	if len(sequence41.Items) != 29 {
		t.Errorf("wrong number of items. got=%d, expected=29", len(sequence41.Items))
	}

	seq42 := testEvalXML("//book[@category='web']/preceding::text()")
	sequence42 := seq42.(*object.Sequence)
	if len(sequence42.Items) != 60 {
		t.Errorf("wrong number of items. got=%d, expected=60", len(sequence42.Items))
	}

	seq43 := testEvalXML("//book[@category='web']/ancestor-or-self::text()")
	sequence43 := seq43.(*object.Sequence)
	if sequence43.Items != nil {
		t.Errorf("the result should be nil")
	}

	seq44 := testEvalXML("//book[@category='web']/ancestor-or-self::element()")
	sequence44 := seq44.(*object.Sequence)
	if len(sequence44.Items) != 5 {
		t.Errorf("wrong number of items. got=%d, expected=5", len(sequence44.Items))
	}

	seq45 := testEvalXML("//book[@category='web']/ancestor-or-self::element()[1]")
	sequence45 := seq45.(*object.Sequence)
	if len(sequence45.Items) != 2 {
		t.Errorf("wrong number of items. got=%d, expected=2", len(sequence45.Items))
	}

	seq46 := testEvalXML("//book[@category='web']/ancestor-or-self::element()[2]")
	sequence46 := seq46.(*object.Sequence)
	if len(sequence46.Items) != 1 {
		t.Errorf("wrong number of items. got=%d, expected=1", len(sequence46.Items))
	}
	node46 := sequence46.Items[0].(*object.BaseNode)
	if node46.Tree().Data != "tt:bookstore" {
		t.Errorf("selected node should be [tt:bookstore]")
	}

	seq47 := testEvalXML("//book[@category='web']/preceding::element()[2]")
	sequence47 := seq47.(*object.Sequence)
	if len(sequence47.Items) != 2 {
		t.Errorf("wrong number of items. got=%d, expected=2", len(sequence47.Items))
	}

	seq48 := testEvalXML("//book[@category='web']/preceding::element()[1]")
	sequence48 := seq48.(*object.Sequence)
	if len(sequence48.Items) != 2 {
		t.Errorf("wrong number of items. got=%d, expected=2", len(sequence48.Items))
	}

	seq49 := testEvalXML("//book[@category='web']/preceding::node()")
	sequence49 := seq49.(*object.Sequence)
	if len(sequence49.Items) != 90 {
		t.Errorf("wrong number of items. got=%d, expected=90", len(sequence49.Items))
	}

	seq50 := testEvalXML("//book[@category='web']/preceding::node()[1]")
	sequence50 := seq50.(*object.Sequence)
	if len(sequence50.Items) != 2 {
		t.Errorf("wrong number of items. got=%d, expected=2", len(sequence50.Items))
	}

	seq51 := testEvalXML("//book[@category='web']/preceding::node()[2]")
	sequence51 := seq51.(*object.Sequence)
	if len(sequence51.Items) != 2 {
		t.Errorf("wrong number of items. got=%d, expected=2", len(sequence51.Items))
	}

	seq52 := testEvalXML("//book[@category='web']/preceding-sibling::node()")
	sequence52 := seq52.(*object.Sequence)
	if len(sequence52.Items) != 11 {
		t.Errorf("wrong number of items. got=%d, expected=11", len(sequence52.Items))
	}

	seq53 := testEvalXML("//book[@category='web']/preceding-sibling::element()[1]")
	sequence53 := seq53.(*object.Sequence)
	if len(sequence53.Items) != 2 {
		t.Errorf("wrong number of items. got=%d, expected=2", len(sequence53.Items))
	}

	seq54 := testEvalXML("//book[@category='web']/preceding-sibling::element()[2]")
	sequence54 := seq54.(*object.Sequence)
	if len(sequence54.Items) != 2 {
		t.Errorf("wrong number of items. got=%d, expected=2", len(sequence54.Items))
	}

	seq55 := testEvalXML("//book[@category='web']/(ancestor::*)[1]")
	sequence55 := seq55.(*object.Sequence)
	if len(sequence55.Items) != 1 {
		t.Errorf("wrong number of items. got=%d, expected=1", len(sequence55.Items))
	}
	node55 := sequence55.Items[0].(*object.BaseNode)
	if node55.Tree().Data != "html" {
		t.Errorf("selected node should be [html]")
	}

	seq56 := testEvalXML("element()")
	sequence56 := seq56.(*object.Sequence)
	if len(sequence56.Items) != 1 {
		t.Errorf("wrong number of items. got=%d, expected=1", len(sequence56.Items))
	}
	node56 := sequence56.Items[0].(*object.BaseNode)
	if node56.Tree().Data != "html" {
		t.Errorf("selected node should be [html]")
	}

	seq57 := testEvalXML("node()")
	sequence57 := seq57.(*object.Sequence)
	if len(sequence57.Items) != 1 {
		t.Errorf("wrong number of items. got=%d, expected=1", len(sequence57.Items))
	}
	node57 := sequence57.Items[0].(*object.BaseNode)
	if node57.Tree().Data != "html" {
		t.Errorf("selected node should be [html]")
	}

	seq58 := testEvalXML("document-node()")
	sequence58 := seq58.(*object.Sequence)
	if len(sequence58.Items) != 1 {
		t.Errorf("wrong number of items. got=%d, expected=1", len(sequence58.Items))
	}
	node58 := sequence58.Items[0].(*object.BaseNode)
	if node58.Type() != object.DocumentNodeType {
		t.Errorf("node type should be document-node")
	}
}

func TestNodeComp(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{`//attribute::category[. = '1'] << //@category[. = '2']`, true},
		{`//attribute::category[. = '1'] >> //@category[. = '2']`, false},
		{`//attribute::category[. = '1'] is //@category[. = '2']`, false},
		{`//attribute::category[. = '1'] is //@category[. = '1']`, true},
		{`//attribute::category[. = '1'] is //book/haha`, false},
		{`//book/haha is //@category[. = '2']`, false},
		{`//month/haha is //year//haha`, true},
		{`//attribute::category[.='cooking'] << //tt:title`, true},
		{`//attribute::category[.='cooking'] >> //tt:title`, false},
		{`//attribute::category[.='cooking'] << //tt:book`, false},
		{`//attribute::category[.='cooking'] >> //tt:book`, true},
		{`//tt:title << //attribute::category[.='cooking']`, false},
		{`//tt:title >> //attribute::category[.='cooking']`, true},
		{`//tt:book << //attribute::category[.='cooking']`, true},
		{`//tt:book >> //attribute::category[.='cooking']`, false},
	}

	for _, tt := range tests {
		seq := testEvalXML(tt.input)
		sequence := seq.(*object.Sequence)
		item := sequence.Items[0].(*object.Boolean)

		if item.Value() != tt.expected {
			t.Errorf("got=%t, expected=%t", item.Value(), tt.expected)
		}
	}
}

func TestNodeExpr(t *testing.T) {
	seq := testEvalXML("//book union //author")
	sequence := seq.(*object.Sequence)
	if len(sequence.Items) != 11 {
		t.Errorf("wrong number of items. got=%d, expected=11", len(sequence.Items))
	}

	seq2 := testEvalXML("//book | //author")
	sequence2 := seq2.(*object.Sequence)
	if len(sequence2.Items) != 11 {
		t.Errorf("wrong number of items. got=%d, expected=11", len(sequence2.Items))
	}

	seq3 := testEvalXML("//book[@category='1'] union //author")
	sequence3 := seq3.(*object.Sequence)
	if len(sequence3.Items) != 7 {
		t.Errorf("wrong number of items. got=%d, expected=7", len(sequence3.Items))
	}

	seq4 := testEvalXML("//book union //author[.='Per Bothner']")
	sequence4 := seq4.(*object.Sequence)
	if len(sequence4.Items) != 6 {
		t.Errorf("wrong number of items. got=%d, expected=6", len(sequence4.Items))
	}

	seq5 := testEvalXML("//book[@category='1'] | //author[.='Per Bothner']")
	sequence5 := seq5.(*object.Sequence)
	if len(sequence5.Items) != 2 {
		t.Errorf("wrong number of items. got=%d, expected=2", len(sequence5.Items))
	}

	seq6 := testEvalXML("//book except //book[@category='1']")
	sequence6 := seq6.(*object.Sequence)
	if len(sequence6.Items) != 4 {
		t.Errorf("wrong number of items. got=%d, expected=4", len(sequence6.Items))
	}

	seq7 := testEvalXML("//book intersect //book[@category='1']")
	sequence7 := seq7.(*object.Sequence)
	if len(sequence7.Items) != 1 {
		t.Errorf("wrong number of items. got=%d, expected=1", len(sequence7.Items))
	}
}

func TestPathWithTypes(t *testing.T) {
	seq := testEvalXML("//1")
	sequence := seq.(*object.Sequence)
	if len(sequence.Items) != 112 {
		t.Errorf("wrong number of items. got=%d, expected=112", len(sequence.Items))
	}

	seq2 := testEvalXML("//(1+2)")
	sequence2 := seq2.(*object.Sequence)
	if len(sequence2.Items) != 112 {
		t.Errorf("wrong number of items. got=%d, expected=112", len(sequence2.Items))
	}

	seq3 := testEvalXML("/1.1")
	sequence3 := seq3.(*object.Sequence)
	if len(sequence3.Items) != 1 {
		t.Errorf("wrong number of items. got=%d, expected=1", len(sequence3.Items))
	}
	item3, ok := sequence3.Items[0].(*object.Decimal)
	if !ok {
		t.Errorf("item should be decimal type")
	}
	if item3.Value() != 1.1 {
		t.Errorf("item value should be 1.1. got=%f", item3.Value())
	}

	seq4 := testEvalXML("//book/1.1e6")
	sequence4 := seq4.(*object.Sequence)
	if len(sequence4.Items) != 5 {
		t.Errorf("wrong number of items. got=%d, expected=5", len(sequence4.Items))
	}
	item4, ok := sequence4.Items[0].(*object.Double)
	if !ok {
		t.Errorf("item should be double type")
	}
	if item4.Value() != 1.1e6 {
		t.Errorf("item value should be 1.1e6. got=%f", item4.Value())
	}

	seq5 := testEvalXML("//year/'abc'")
	sequence5 := seq5.(*object.Sequence)
	if len(sequence5.Items) != 4 {
		t.Errorf("wrong number of items. got=%d, expected=4", len(sequence5.Items))
	}
	item5, ok := sequence5.Items[0].(*object.String)
	if !ok {
		t.Errorf("item should be string type")
	}
	if item5.Value() != "abc" {
		t.Errorf("item value should be 'abc'. got=%s", item5.Value())
	}

	seq6 := testEvalXML("//'abc'")
	sequence6 := seq6.(*object.Sequence)
	if len(sequence6.Items) != 112 {
		t.Errorf("wrong number of items. got=%d, expected=112", len(sequence6.Items))
	}
}

func testEval(input string) object.Item {
	l := lexer.New(input)
	p := parser.New(l)
	xpath := p.ParseXPath()
	ctx := object.NewContext()

	docFunc := bif.F["fn:doc"]
	str := &object.String{}
	str.SetValue("testdata/quotes-1.html")
	err := docFunc(ctx, str)
	if err != nil {
		return err
	}

	return Eval(xpath, ctx)
}

func testEvalXML(input string) object.Item {
	l := lexer.New(input)
	p := parser.New(l)
	xpath := p.ParseXPath()
	ctx := object.NewContext()

	docFunc := bif.F["fn:doc"]
	str := &object.String{}
	str.SetValue("testdata/company.xml")
	err := docFunc(ctx, str)
	if err != nil {
		return err
	}

	return Eval(xpath, ctx)
}

func testEvalXML2(input string) object.Item {
	l := lexer.New(input)
	p := parser.New(l)
	xpath := p.ParseXPath()
	ctx := object.NewContext()

	docFunc := bif.F["fn:doc"]
	str := &object.String{}
	str.SetValue("testdata/company_2.xml")
	err := docFunc(ctx, str)
	if err != nil {
		return err
	}

	return Eval(xpath, ctx)
}

func testNumberObject(t *testing.T, item object.Item, expected interface{}) {
	switch item := item.(type) {
	case *object.Integer:
		if item.Value() != expected {
			t.Errorf("object.Integer has wrong value. got=%d, want=%d", item.Value(), expected)
		}
	case *object.Decimal:
		e := fmt.Sprintf("%f", expected)
		v := fmt.Sprintf("%f", item.Value())
		if v != e {
			t.Errorf("object.Decimal has wrong value. got=%f, want=%f", item.Value(), expected)
		}
	case *object.Double:
		e := fmt.Sprintf("%f", expected)
		v := fmt.Sprintf("%f", item.Value())
		if v != e {
			t.Errorf("object.Double has wrong value. got=%f, want=%f", item.Value(), expected)
		}
	default:
		t.Errorf("Unkown item type. got=%s", item.Type())
	}
}

func testStringObject(t *testing.T, item object.Item, expected interface{}) {
	switch item := item.(type) {
	case *object.String:
		if item.Value() != expected {
			t.Errorf("object.String has wrong value. got=%s, want=%s", item.Value(), expected)
		}
	default:
		t.Errorf("item type must object.String. got=%s", item.Type())
	}
}

func testSequenceObject(t *testing.T, item object.Item, expected []interface{}) {
	switch item := item.(type) {
	case *object.Sequence:
		if len(item.Items) != len(expected) {
			t.Errorf("length of the item must be the same. got=%d, want=%d", len(item.Items), len(expected))
		}
		for i := 0; i < len(item.Items); i++ {
			switch item.Items[i].(type) {
			case *object.Integer:
				testNumberObject(t, item.Items[i], expected[i])
			case *object.Decimal:
				testNumberObject(t, item.Items[i], expected[i])
			case *object.Double:
				testNumberObject(t, item.Items[i], expected[i])
			case *object.String:
				testStringObject(t, item.Items[i], expected[i])
			default:
				t.Errorf("Unkown item type. got=%s", item.Items[i].Type())
			}
		}
	default:
		t.Errorf("item type must object.String. got=%s", item.Type())
	}
}
