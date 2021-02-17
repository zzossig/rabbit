package eval

import (
	"fmt"
	"testing"

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
				if item.Value != tt.expected {
					t.Errorf("item has wrong value. got=%d, want=%d", item.Value, tt.expected)
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
				if item.Value != tt.expected {
					t.Errorf("item has wrong value. got=%d, want=%d", item.Value, tt.expected)
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

		if item.Value != tt.expected {
			t.Errorf("got=%s, expected=%s", item.Value, tt.expected)
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
			if bl.Value != tt.expected {
				t.Errorf("got: %v, expected: %v", bl.Value, tt.expected)
			}
		}
	}
}

func testEval(input string) object.Item {
	l := lexer.New(input)
	p := parser.New(l)
	xpath := p.ParseXPath()
	ctx := object.NewContext()
	ctx.NewReaderFile("text.txt", true)

	return Eval(xpath, ctx)
}

func testNumberObject(t *testing.T, item object.Item, expected interface{}) {
	switch item := item.(type) {
	case *object.Integer:
		if item.Value != expected {
			t.Errorf("object.Integer has wrong value. got=%d, want=%d", item.Value, expected)
		}
	case *object.Decimal:
		e := fmt.Sprintf("%f", expected)
		v := fmt.Sprintf("%f", item.Value)
		if v != e {
			t.Errorf("object.Decimal has wrong value. got=%f, want=%f", item.Value, expected)
		}
	case *object.Double:
		e := fmt.Sprintf("%f", expected)
		v := fmt.Sprintf("%f", item.Value)
		if v != e {
			t.Errorf("object.Double has wrong value. got=%f, want=%f", item.Value, expected)
		}
	default:
		t.Errorf("Unkown item type. got=%s", item.Type())
	}
}

func testStringObject(t *testing.T, item object.Item, expected interface{}) {
	switch item := item.(type) {
	case *object.String:
		if item.Value != expected {
			t.Errorf("object.String has wrong value. got=%s, want=%s", item.Value, expected)
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
