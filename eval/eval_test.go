package eval

import (
	"fmt"
	"testing"

	"github.com/zzossig/xpath/lexer"
	"github.com/zzossig/xpath/object"
	"github.com/zzossig/xpath/parser"
)

func testEval(input string) object.Item {
	l := lexer.New(input)
	p := parser.New(l)
	xpath := p.ParseXPath()
	env := object.NewEnv()

	return Eval(xpath, env)
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

func TestEvalArithmeticExpr(t *testing.T) {
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

func TestEvalArrayExpr(t *testing.T) {
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

func TestBuiltinFunctions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`abs(-2.5)`, 2.5},
	}

	for _, tt := range tests {
		seq := testEval(tt.input)
		sequence := seq.(*object.Sequence)

		for _, item := range sequence.Items {
			switch item := item.(type) {
			case *object.Builtin:
				testNumberObject(t, item.Func(item.Args...), tt.expected)
			default:
				t.Errorf("Unkown item type. got=%s", item.Type())
			}
		}
	}
}
