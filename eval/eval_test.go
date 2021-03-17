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
		{"'a' => upper-case() => lower-case()", "a"},
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
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`if ("a") then 2 else 3`, 2},
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
		t.Errorf("doc node not have a div node of child")
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
		t.Errorf("wrong number of items. got=%d, expected=1", len(sequence4.Items))
	}

	seq5 := testEvalXML("//book")
	sequence5 := seq5.(*object.Sequence)
	if len(sequence5.Items) != 5 {
		t.Errorf("wrong number of items. got=%d, expected=1", len(sequence5.Items))
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
		t.Errorf("result must be a nil")
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
		t.Errorf("result must be a nil")
	}

	seq22 := testEvalXML("//month/preceding::year")
	sequence22 := seq22.(*object.Sequence)
	if len(sequence22.Items) != 1 {
		t.Errorf("wrong number of items. got=%d, expected=1", len(sequence22.Items))
	}

	seq23 := testEvalXML("//month/preceding-sibling::year")
	sequence23 := seq23.(*object.Sequence)
	if sequence23.Items != nil {
		t.Errorf("result must be a nil")
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
}

func TestPathPredicateExpr(t *testing.T) {
}

func testEval(input string) object.Item {
	l := lexer.New(input)
	p := parser.New(l)
	xpath := p.ParseXPath()
	ctx := object.NewContext()

	docFunc := bif.Builtins["fn:doc"]
	str := &object.String{}
	str.SetValue("testdata/quotes-1.html")
	docNode := docFunc(str)
	if !bif.IsError(docNode) {
		d := docNode.(*object.BaseNode)
		ctx.Doc = object.Node(d)
	}

	return Eval(xpath, ctx)
}

func testEvalXML(input string) object.Item {
	l := lexer.New(input)
	p := parser.New(l)
	xpath := p.ParseXPath()
	ctx := object.NewContext()

	docFunc := bif.Builtins["fn:doc"]
	str := &object.String{}
	str.SetValue("testdata/company.xml")
	docNode := docFunc(str)
	if !bif.IsError(docNode) {
		d := docNode.(*object.BaseNode)
		ctx.Doc = object.Node(d)
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
