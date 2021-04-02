package eval

import (
	"fmt"
	"math"
	"strings"
	"testing"

	"github.com/zzossig/rabbit/bif"
	"github.com/zzossig/rabbit/lexer"
	"github.com/zzossig/rabbit/object"
	"github.com/zzossig/rabbit/parser"
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
		{`[ [1, 2, 3], [4, 5, 6]](2)`, []interface{}{4, 5, 6}},
		{`array:join((["a", "b"], ["c", "d"], [["e", "f"]]))`, []interface{}{"a", "b", "c", "d", []interface{}{"e", "f"}}},
		{`array:join((["a", "b"], ["c", "d"], [ ]))`, []interface{}{"a", "b", "c", "d"}},
		{`array:join((["a", "b"], ["c", "d"]))`, []interface{}{"a", "b", "c", "d"}},
		{`array:join([1, 2, 3])`, []interface{}{1, 2, 3}},
		{`array:join(())`, []interface{}{}},
		{`array:reverse([])`, []interface{}{}},
		{`array:reverse([(1 to 5)])`, []interface{}{[]interface{}{1, 2, 3, 4, 5}}},
		{`array:reverse([("a", "b"), ("c", "d")])`, []interface{}{[]interface{}{"c", "d"}, []interface{}{"a", "b"}}},
		{`array:reverse(["a", "b", "c", "d"])`, []interface{}{"d", "c", "b", "a"}},
		{`array:tail([5])`, []interface{}{}},
		{`array:tail([5, 6, 7, 8])`, []interface{}{6, 7, 8}},
		{`array:head([["a", "b"], ["c", "d"]])`, []interface{}{"a", "b"}},
		{`array:insert-before(["a", "b", "c", "d"], 3, ["x", "y"])`, []interface{}{"a", "b", []interface{}{"x", "y"}, "c", "d"}},
		{`array:insert-before(["a", "b", "c", "d"], 5, ("x", "y"))`, []interface{}{"a", "b", "c", "d", []interface{}{"x", "y"}}},
		{`array:insert-before(["a", "b", "c", "d"], 3, ("x", "y"))`, []interface{}{"a", "b", []interface{}{"x", "y"}, "c", "d"}},
		{`array:remove(["a", "b", "c", "d"], ())`, []interface{}{"a", "b", "c", "d"}},
		{`array:remove(["a", "b", "c", "d"], 1 to 3)`, []interface{}{"d"}},
		{`array:remove(["a"], 1)`, []interface{}{}},
		{`array:remove(["a", "b", "c", "d"], 2)`, []interface{}{"a", "c", "d"}},
		{`array:remove(["a", "b", "c", "d"], 1)`, []interface{}{"b", "c", "d"}},
		{`array:subarray([ ], 1, 0)`, []interface{}{}},
		{`array:subarray(["a", "b", "c", "d"], 5, 0)`, []interface{}{}},
		{`array:subarray(["a", "b", "c", "d"], 2, 3)`, []interface{}{"b", "c", "d"}},
		{`array:subarray(["a", "b", "c", "d"], 2, 2)`, []interface{}{"b", "c"}},
		{`array:subarray(["a", "b", "c", "d"], 2, 1)`, []interface{}{"b"}},
		{`array:subarray(["a", "b", "c", "d"], 2, 0)`, []interface{}{}},
		{`array:subarray(["a", "b", "c", "d"], 5)`, []interface{}{}},
		{`array:subarray(["a", "b", "c", "d"], 2)`, []interface{}{"b", "c", "d"}},
		{`array:append(["a", "b", "c"], ["d", "e"])`, []interface{}{"a", "b", "c", []interface{}{"d", "e"}}},
		{`array:append(["a", "b", "c"], ("d", "e"))`, []interface{}{"a", "b", "c", []interface{}{"d", "e"}}},
		{`array:append(["a", "b", "c"], "d")`, []interface{}{"a", "b", "c", "d"}},
		{`array:put(["a", "b", "c"], 2, ("d", "e"))`, []interface{}{"a", []interface{}{"d", "e"}, "c"}},
		{`array:put(["a", "b", "c"], 2, "d")`, []interface{}{"a", "d", "c"}},
		{`["a", ["b", "c"]] => array:get(2)`, []interface{}{"b", "c"}},
		{`array:size(["a", "b", "c"])`, []interface{}{3}},
		{"array{1,2,3}", []interface{}{1, 2, 3}},
		{"array{1*2,2+3,3-4,5 idiv 5, 5 div 5}", []interface{}{2, 5, -1, 1, 1.0}},
		{"[3 mod 2, 'a', 'b', 1.1]", []interface{}{1, "a", "b", 1.1}},
	}

	for _, tt := range tests {
		seq := testEval(tt.input)
		sequence := seq.(*object.Sequence)

		for _, item := range sequence.Items {
			switch item := item.(type) {
			case *object.Integer:
			case *object.Decimal:
			case *object.Double:
			case *object.String:
			case *object.Array:
				for i, v := range item.Items {
					switch vv := v.(type) {
					case *object.Integer:
						testNumberObject(t, v, tt.expected[i])
					case *object.Decimal:
						testNumberObject(t, v, tt.expected[i])
					case *object.Double:
						testNumberObject(t, v, tt.expected[i])
					case *object.String:
						testStringObject(t, v, tt.expected[i])
					case *object.Sequence:
						s := tt.expected[i].([]interface{})
						if len(vv.Items) != len(s) {
							t.Errorf("sequence length not match")
						}
						for j, it := range vv.Items {
							switch it := it.(type) {
							case *object.Integer:
								if it.Value() != s[j] {
									t.Errorf("sequence value not match. got=%d, expected=%d", it.Value(), s[j])
								}
							case *object.Decimal:
								if it.Value() != s[j] {
									t.Errorf("sequence value not match. got=%f, expected=%f", it.Value(), s[j])
								}
							case *object.Double:
								if it.Value() != s[j] {
									t.Errorf("sequence value not match. got=%f, expected=%f", it.Value(), s[j])
								}
							case *object.String:
								if it.Value() != s[j] {
									t.Errorf("sequence value not match. got=%s, expected=%s", it.Value(), s[j])
								}
							case *object.Boolean:
								if it.Value() != s[j] {
									t.Errorf("sequence value not match. got=%v, expected=%v", it.Value(), s[j])
								}
							}
						}
					case *object.Array:
						s := tt.expected[i].([]interface{})
						if len(vv.Items) != len(s) {
							t.Errorf("sequence length not match")
						}
						for j, it := range vv.Items {
							switch it := it.(type) {
							case *object.Integer:
								if it.Value() != s[j] {
									t.Errorf("sequence value not match. got=%d, expected=%d", it.Value(), s[j])
								}
							case *object.Decimal:
								if it.Value() != s[j] {
									t.Errorf("sequence value not match. got=%f, expected=%f", it.Value(), s[j])
								}
							case *object.Double:
								if it.Value() != s[j] {
									t.Errorf("sequence value not match. got=%f, expected=%f", it.Value(), s[j])
								}
							case *object.String:
								if it.Value() != s[j] {
									t.Errorf("sequence value not match. got=%s, expected=%s", it.Value(), s[j])
								}
							case *object.Boolean:
								if it.Value() != s[j] {
									t.Errorf("sequence value not match. got=%v, expected=%v", it.Value(), s[j])
								}
							}
						}
					default:
						t.Errorf("Unknown item type. got=%s", item.Type())
					}

				}
			default:
				t.Errorf("Unknown item type. got=%s", item.Type())
			}
		}
	}
}

func TestEvalArray2(t *testing.T) {
	tests := []struct {
		input    string
		expected []interface{}
	}{
		{`([1,2,3], [4,5,6])?2`, []interface{}{2, 5}},
		{`[4, 5, 6]?2`, []interface{}{5}},
		{`([1,2,3], [1,2,5], [1,2])[?2 = 2]`, []interface{}{[]int{1, 2, 3}, []int{1, 2, 5}, []int{1, 2}}},
		{`[[1, 2, 3], [4, 5, 6]]?*`, []interface{}{[]int{1, 2, 3}, []int{4, 5, 6}}},
		{`[1, 2, 5, 7]?*`, []interface{}{1, 2, 5, 7}},
		{`[1, 2, 3, 4]?2`, []interface{}{2}},
		{`array { (), (27, 17, 0) }(2)`, []interface{}{17}},
		{`array { (), (27, 17, 0) }(1)`, []interface{}{27}},
		{`[ [1, 2, 3], [4, 5, 6]](2)(2)`, []interface{}{5}},
		{`[ 1, 2, 5, 7 ](4)`, []interface{}{7}},
		{`array:flatten([(1,0), (1,1), (0,1), (0,0)])`, []interface{}{1, 0, 1, 1, 0, 1, 0, 0}},
		{`array:flatten(([1, 2, 5], [[10, 11], 12], [], 13))`, []interface{}{1, 2, 5, 10, 11, 12, 13}},
		{`array:flatten([1, 4, 6, 5, 3])`, []interface{}{1, 4, 6, 5, 3}},
		{`array:head([("a", "b"), ("c", "d")])`, []interface{}{"a", "b"}},
		{`array:head([5, 6, 7, 8])`, []interface{}{5}},
		{`["a", "b", "c"] => array:get(2)`, []interface{}{"b"}},
		{`array:size([[ ]])`, []interface{}{1}},
		{`array:size([ ])`, []interface{}{0}},
		{`array:size(["a", ["b", "c"]])`, []interface{}{2}},
		{`array:size(["a", "b", "c"])`, []interface{}{3}},
	}

	for _, tt := range tests {
		seq := testEval(tt.input)
		sequence := seq.(*object.Sequence)
		testSequenceObject(t, sequence, tt.expected)
	}
}

func TestFunctionCall(t *testing.T) {
	tests := []struct {
		input    string
		expected []interface{}
	}{
		{`(1 to 3)!(.*.) => fn:sum()`, []interface{}{14}},
		{`fn:string-join((1 to 5)!"*")`, []interface{}{"*****"}},
		{"filter((1, 2), fn:empty#1)", []interface{}{}},
		{"fn:filter((4, 5), fn:exists(?))", []interface{}{4, 5}},
		{"fn:filter((1, 2), fn:exists#1)", []interface{}{1, 2}},
		{"fn:filter(1 to 10, function($a) {$a mod 2 = 0})", []interface{}{2, 4, 6, 8, 10}},
		{"fn:codepoints-to-string(())", []interface{}{""}},
		{"fn:codepoints-to-string((2309, 2358, 2378, 2325))", []interface{}{"अशॊक"}},
		{"fn:codepoints-to-string((66, 65, 67, 72))", []interface{}{"BACH"}},
		{"fn:string-to-codepoints('Thérèse')", []interface{}{84, 104, 233, 114, 232, 115, 101}},
		{"fn:for-each(('john','jane'), fn:string-to-codepoints#1)", []interface{}{106, 111, 104, 110, 106, 97, 110, 101}},
		{"fn:for-each(1 to 5, function($a) { $a * $a })", []interface{}{1, 4, 9, 16, 25}},
		{"fn:for-each(('a','b','c'), xs:string(?))", []interface{}{"a", "b", "c"}},
		{"fn:for-each(('23', '29'), xs:integer#1)", []interface{}{23, 29}},
		{"fn:for-each(('23', '29'), xs:integer(?))", []interface{}{23, 29}},
		{"(1 to 20)[fn:position() = 20]", []interface{}{20}},
		{"(1 to 20)[fn:position() = 1]", []interface{}{1}},
		{"(1 to 20)[fn:last() - 1]", []interface{}{19}},
		{"fn:min((1,4,[8,5,[7,6,(2,30,0)]]))", []interface{}{0}},
		{"fn:min([true(), false()])", []interface{}{false}},
		{"fn:min((7.5,1.0,9.9))", []interface{}{1.0}},
		{"fn:min(('c','b','a','z'))", []interface{}{"a"}},
		{"fn:min((1,2,3,4,5))", []interface{}{1}},
		{"fn:max([true(), false()])", []interface{}{true}},
		{"fn:max((7.5,1.0,9.9))", []interface{}{9.9}},
		{"fn:max(('c','b','a','z'))", []interface{}{"z"}},
		{"fn:max((1,2,3,4,5))", []interface{}{5}},
		{"fn:max((1,2,3,4,5,[97,43,201,422,[777,542,999,321]]))", []interface{}{999}},
		{"fn:avg(())", []interface{}{}},
		{"fn:avg((1,2.9,3))", []interface{}{2.3}},
		{"fn:avg((1,2,3))", []interface{}{2}},
		{"fn:avg((1.1, 2.2, 3.3))", []interface{}{2.2}},
		{"fn:sum([[1, 2], [3, 4, (6,7,[8,9,10,11])]])", []interface{}{61}},
		{"fn:sum((1.1, 2.2, 3.3))", []interface{}{6.6}},
		{"fn:sum(1 to 3)", []interface{}{6}},
		{"fn:sum(())", []interface{}{0}},
		{"fn:sum((1,2,3))", []interface{}{6}},
		{"fn:sum([1,2,3])", []interface{}{6}},
		{"fn:count([1,2,3])", []interface{}{1}},
		{"fn:count([])", []interface{}{1}},
		{"let $seq2 := (98.5, 98.3, 98.9) return fn:count($seq2[. > 100])", []interface{}{0}},
		{"let $seq2 := (98.5, 98.3, 98.9) return fn:count($seq2)", []interface{}{3}},
		{"let $seq3 := () return fn:count($seq3)", []interface{}{0}},
		{"let $seq1 := ($item1, $item2) return fn:count($seq1)", []interface{}{2}},
		{"let $seq := ('item1', 'item2', 'item3', 'item4', 'item5') return subsequence($seq,-2,-1)", []interface{}{}},
		{"let $seq := ('item1', 'item2', 'item3', 'item4', 'item5') return subsequence($seq,-2,0)", []interface{}{}},
		{"let $seq := ('item1', 'item2', 'item3', 'item4', 'item5') return subsequence($seq,-2,5)", []interface{}{"item1", "item2"}},
		{"let $seq := ('item1', 'item2', 'item3', 'item4', 'item5') return subsequence($seq,1,5)", []interface{}{"item1", "item2", "item3", "item4", "item5"}},
		{"let $seq := ('item1', 'item2', 'item3', 'item4', 'item5') return subsequence($seq,0,5)", []interface{}{"item1", "item2", "item3", "item4"}},
		{"let $seq := ('item1', 'item2', 'item3', 'item4', 'item5') return subsequence($seq, -1)", []interface{}{"item1", "item2", "item3", "item4", "item5"}},
		{"let $seq := ('item1', 'item2', 'item3', 'item4', 'item5') return subsequence($seq, 4)", []interface{}{"item4", "item5"}},
		{"let $seq := ('item1', 'item2', 'item3', 'item4', 'item5') return subsequence($seq, 3, 2)", []interface{}{"item3", "item4"}},
		{"let $abc := ('a', 'b', 'c') return reverse(([1,2,3],[4,5,6]))", []interface{}{[]int{4, 5, 6}, []int{1, 2, 3}}},
		{"let $abc := ('a', 'b', 'c') return fn:reverse([1,2,3])", []interface{}{[]int{1, 2, 3}}},
		{"let $abc := ('a', 'b', 'c') return fn:reverse(())", []interface{}{}},
		{"let $abc := ('a', 'b', 'c') return reverse(('hello'))", []interface{}{"hello"}},
		{"let $abc := ('a', 'b', 'c') return fn:reverse($abc)", []interface{}{"c", "b", "a"}},
		{"let $abc := ('a', 'b', 'c') return fn:insert-before($abc, 4, 'z')", []interface{}{"a", "b", "c", "z"}},
		{"let $abc := ('a', 'b', 'c') return fn:insert-before($abc, 3, 'z')", []interface{}{"a", "b", "z", "c"}},
		{"let $abc := ('a', 'b', 'c') return fn:insert-before($abc, 2, 'z')", []interface{}{"a", "z", "b", "c"}},
		{"let $abc := ('a', 'b', 'c') return fn:insert-before($abc, 0, 'z')", []interface{}{"z", "a", "b", "c"}},
		{"let $abc := ('a', 'b', 'c') return fn:insert-before($abc, 1, 'z')", []interface{}{"z", "a", "b", "c"}},
		{"fn:tail([1,2,3])", []interface{}{}},
		{"fn:tail(())", []interface{}{}},
		{"fn:tail('a')", []interface{}{}},
		{"fn:tail(('a', 'b', 'c'))", []interface{}{"b", "c"}},
		{"fn:tail(1 to 5)", []interface{}{2, 3, 4, 5}},
		{"fn:head([1,2,3])", []interface{}{[]int{1, 2, 3}}},
		{"fn:head(())", []interface{}{}},
		{"fn:head(('a', 'b', 'c'))", []interface{}{"a"}},
		{"fn:head(1 to 5)", []interface{}{1}},
		{"fn:exists('')", []interface{}{true}},
		{"fn:exists(map{})", []interface{}{true}},
		{"fn:exists([])", []interface{}{true}},
		{"fn:exists(fn:remove(('hello', 'world'), 1))", []interface{}{true}},
		{"fn:exists(fn:remove(('hello'), 1))", []interface{}{false}},
		{"let $abc := ('a', 'b', 'c') return fn:remove((), 3)", []interface{}{}},
		{"let $abc := ('a', 'b', 'c') return fn:remove($abc, 6)", []interface{}{"a", "b", "c"}},
		{"let $abc := ('a', 'b', 'c') return remove($abc, 1)", []interface{}{"b", "c"}},
		{"let $abc := ('a', 'b', 'c') return fn:remove($abc, 0)", []interface{}{"a", "b", "c"}},
		{"empty('')", []interface{}{false}},
		{"fn:empty(fn:remove(('hello', 'world'), 1))", []interface{}{false}},
		{"fn:empty(map{})", []interface{}{false}},
		{"fn:empty([])", []interface{}{false}},
		{"fn:empty((1,2,3)[10])", []interface{}{true}},
		{"fn:substring-after('abcdcba', 'abc')", []interface{}{"dcba"}},
		{"fn:substring-before((), ())", []interface{}{""}},
		{"substring-before('tattoo', 'tatto')", []interface{}{""}},
		{"substring-before('tattoo', 'ttoo')", []interface{}{"ta"}},
		{"fn:ends-with((), ())", []interface{}{true}},
		{"ends-with('tattoo', 'attoo')", []interface{}{true}},
		{"ends-with('tattoo', 'tattoo')", []interface{}{true}},
		{"fn:starts-with((), ())", []interface{}{true}},
		{"fn:starts-with('tattoo' 'att')", []interface{}{false}},
		{"fn:starts-with('tattoo' 'tat')", []interface{}{true}},
		{"fn:contains((), '')", []interface{}{true}},
		{"fn:contains((), ())", []interface{}{true}},
		{"fn:contains('', ())", []interface{}{true}},
		{"fn:contains('tattoo', 'ttt')", []interface{}{false}},
		{"fn:contains('tattoo', 't')", []interface{}{true}},
		{"fn:normalize-space('  \t\nabc\t\n')", []interface{}{"abc"}},
		{"fn:string-length(())", []interface{}{0}},
		{"fn:string-length('Harp not on that string, madam; that is past.')", []interface{}{45}},
		{"fn:string-join(1 to 5, ', ')", []interface{}{"1, 2, 3, 4, 5"}},
		{"fn:string-join((), 'separator')", []interface{}{""}},
		{
			"fn:string-join(('Blow, ', 'blow, ', 'thou ', 'winter ', 'wind!'), '')",
			[]interface{}{"Blow, blow, thou winter wind!"},
		},
		{"fn:string-join(('Now', 'is', 'the', 'time', '...'), ' ')", []interface{}{"Now is the time ..."}},
		{"fn:string-join(1 to 9)", []interface{}{"123456789"}},
		{"fn:not('false')", []interface{}{false}},
		{`fn:not(fn:true())`, []interface{}{false}},
		{`fn:not(())`, []interface{}{true}},
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
		{
			"fn:for-each-pair(('a', 'b', 'c'), ('x', 'y', 'z'), concat#2)",
			[]interface{}{"ax", "by", "cz"},
		},
		{
			"fn:for-each-pair(1 to 5, 1 to 5, function($a, $b){10*$a + $b})",
			[]interface{}{11, 22, 33, 44, 55},
		},
		{"fn:concat(1,2,3)", []interface{}{"123"}},
		{"fn:concat(1,2,3,'a')", []interface{}{"123a"}},
		{"fn:concat('un', 'grateful')", []interface{}{"ungrateful"}},
		{
			"fn:concat('Thy ', (), 'old ', 'groans', '', ' ring', ' yet', ' in', ' my', ' ancient',' ears.')",
			[]interface{}{"Thy old groans ring yet in my ancient ears."},
		},
		{"fn:concat('Ciao!',())", []interface{}{"Ciao!"}},
		{
			"fn:concat('Ingratitude, ', 'thou ', 'marble-hearted', ' fiend!')",
			[]interface{}{"Ingratitude, thou marble-hearted fiend!"},
		},
		{"fn:concat(01, 02, 03, 04, true())", []interface{}{"1234true"}},
		{"string-join((1,2,3),'a')", []interface{}{"1a2a3"}},
		{"fn:substring('motor car', 6)", []interface{}{" car"}},
		{"substring('metadata', 4, 3)", []interface{}{"ada"}},
		{"substring('12345', 1.5, 2.6)", []interface{}{"234"}},
		{"substring('12345', 0, 3)", []interface{}{"12"}},
		{"substring('12345', 5, -3)", []interface{}{""}},
		{"substring('12345', -3, 5)", []interface{}{"1"}},
		{"substring('12345', 0 div 0E0, 3)", []interface{}{""}},
		{"substring('12345', 1, 0 div 0E0)", []interface{}{""}},
		{"substring((), 1, 3)", []interface{}{""}},
		{"substring('12345', -42, 1 div 0e0)", []interface{}{"12345"}},
		{"substring('12345', -1 div 0E0, 1 div 0e0)", []interface{}{""}},
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
		{`10 || '/' || 6`, "10/6"},
	}

	for _, tt := range tests {
		seq := testEval(tt.input)
		sequence := seq.(*object.Sequence)

		for _, item := range sequence.Items {
			switch item := item.(type) {
			case *object.String:
				testStringObject(t, item, tt.expected)
			default:
				t.Errorf("Unknown item type. got=%s", item.Type())
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
				t.Errorf("Unknown item type. got=%s", item.Type())
			}
		}
	}
}

func TestPredicate(t *testing.T) {
	tests := []struct {
		input    string
		expected []interface{}
	}{
		{`(1,2,3,4)[1]`, []interface{}{1}},
		{`(1,2,3,4)[1+1]`, []interface{}{2}},
		{`(2,1,3,4)[.=2]`, []interface{}{2}},
	}

	for _, tt := range tests {
		seq := testEval(tt.input)
		sequence := seq.(*object.Sequence)
		testSequenceObject(t, sequence, tt.expected)
	}

	seq2 := testEval("(1,2,3)[4]")
	sequence2 := seq2.(*object.Sequence)
	if sequence2.Items != nil {
		t.Errorf("the result should be empty sequence")
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
		{`(map {"first": "Tom"}, map {"first": "Dick"}, map {"first": "Harry"})?first`, []interface{}{"Tom", "Dick", "Harry"}},
		{`map { "first" : "Jenna", "last" : "Scott" }?first`, []interface{}{"Jenna"}},
		{`let $weekdays := map {"Su" : "Sunday"} return $weekdays?*`, []interface{}{"Sunday"}},
		{`let $weekdays := map {"Su" : "Sunday"} return $weekdays?Su`, []interface{}{"Sunday"}},
		{`let $weekdays := map {"Su" : "Sunday"} return $weekdays("Su")`, []interface{}{"Sunday"}},
		{`map:contains(map{"abc":23, "xyz":()}, "xyz")`, []interface{}{true}},
		{`map:contains(map{"xyz":23}, "xyz")`, []interface{}{true}},
		{`map:contains(map{}, "xyz")`, []interface{}{false}},
		{
			`let $week := map{0:"Sonntag", 1:"Montag", 2:"Dienstag", 3:"Mittwoch", 4:"Donnerstag", 5:"Freitag", 6:"Samstag"} return map:contains($week, 9)`,
			[]interface{}{false},
		},
		{
			`let $week := map{0:"Sonntag", 1:"Montag", 2:"Dienstag", 3:"Mittwoch", 4:"Donnerstag", 5:"Freitag", 6:"Samstag"} return map:contains($week, 2)`,
			[]interface{}{true},
		},
		{
			`let $week := map{0:"Sonntag", 1:"Montag", 2:"Dienstag", 3:"Mittwoch", 4:"Donnerstag", 5:"Freitag", 6:"Samstag"} return map:get($week, 9)`,
			[]interface{}{},
		},
		{
			`let $week := map{0:"Sonntag", 1:"Montag", 2:"Dienstag", 3:"Mittwoch", 4:"Donnerstag", 5:"Freitag", 6:"Samstag"} return map:get($week, 4)`,
			[]interface{}{"Donnerstag"},
		},
		{`map:keys(map{1.5:"yes"})`, []interface{}{1.5}},
		{`map:keys(map{true():"yes"})`, []interface{}{true}},
		{`map:keys(map{"do":"yes"})`, []interface{}{"do"}},
		{`map:keys(map{1:"yes"})`, []interface{}{1}},
		{`map:size(((map{})))`, []interface{}{0}},
		{`map:size(map{"true":1, "false":0})`, []interface{}{2}},
		{`map:size(map{})`, []interface{}{0}},
		{`map{"a":1}?a`, []interface{}{1}},
		{`map{"a":1,"b":2,"c":3}?("a","b")`, []interface{}{1, 2}},
		{
			`let $books := map {
				"book": map {
					"title": "Data on the Web",
					"year": 2000,
					"author": [
						map {
							"last": "Abiteboul",
							"first": "Serge"
						},
						map {
							"last": "Buneman",
							"first": "Peter"
						},
						map {
							"last": "Suciu",
							"first": "Dan"
						}
					],
					"publisher": "Morgan Kaufmann Publishers",
					"price": 39.95
				}
			} return $books("book")("author")(1)("last")`,
			[]interface{}{"Abiteboul"},
		},
	}

	for _, tt := range tests {
		seq := testEval(tt.input)
		sequence := seq.(*object.Sequence)
		testSequenceObject(t, sequence, tt.expected)
	}
}

func TestMapExpr2(t *testing.T) {
	tests := []struct {
		input    string
		keys     []interface{}
		expected []interface{}
	}{
		{
			`let $week := map{0:"Sonntag", 1:"Montag", 2:"Dienstag", 3:"Mittwoch", 4:"Donnerstag", 5:"Freitag", 6:"Samstag"} return map:merge(($week, map{6:"Sonnabend"}), map{"duplicates":"combine"})`,
			[]interface{}{0, 1, 2, 3, 4, 5, 6},
			[]interface{}{"Sonntag", "Montag", "Dienstag", "Mittwoch", "Donnerstag", "Freitag", []interface{}{"Samstag", "Sonnabend"}},
		},
		{
			`let $week := map{0:"Sonntag", 1:"Montag", 2:"Dienstag", 3:"Mittwoch", 4:"Donnerstag", 5:"Freitag", 6:"Samstag"} return map:merge(($week, map{6:"Sonnabend"}), map{"duplicates":"use-first"})`,
			[]interface{}{0, 1, 2, 3, 4, 5, 6},
			[]interface{}{"Sonntag", "Montag", "Dienstag", "Mittwoch", "Donnerstag", "Freitag", "Samstag"},
		},
		{
			`let $week := map{0:"Sonntag", 1:"Montag", 2:"Dienstag", 3:"Mittwoch", 4:"Donnerstag", 5:"Freitag", 6:"Samstag"} return map:merge(($week, map{6:"Sonnabend"}), map{"duplicates":"use-last"})`,
			[]interface{}{0, 1, 2, 3, 4, 5, 6},
			[]interface{}{"Sonntag", "Montag", "Dienstag", "Mittwoch", "Donnerstag", "Freitag", "Sonnabend"},
		},
		{
			`let $week := map{0:"Sonntag", 1:"Montag", 2:"Dienstag", 3:"Mittwoch", 4:"Donnerstag", 5:"Freitag", 6:"Samstag"} return map:merge(($week, map{7:"Unbekannt"}))`,
			[]interface{}{0, 1, 2, 3, 4, 5, 6, 7},
			[]interface{}{"Sonntag", "Montag", "Dienstag", "Mittwoch", "Donnerstag", "Freitag", "Samstag", "Unbekannt"},
		},
		{
			`map:merge((map:entry(0, "no"), map:entry(1, "yes")))`,
			[]interface{}{0, 1},
			[]interface{}{"no", "yes"},
		},
		{
			`map:merge(())`,
			[]interface{}{},
			[]interface{}{},
		},
		{
			`let $week := map{0:"Sonntag", 1:"Montag", 2:"Dienstag", 3:"Mittwoch", 4:"Donnerstag", 5:"Freitag", 6:"Samstag"} return map:remove($week, (0, 6 to 7))`,
			[]interface{}{1, 2, 3, 4, 5},
			[]interface{}{"Montag", "Dienstag", "Mittwoch", "Donnerstag", "Freitag"},
		},
		{
			`let $week := map{0:"Sonntag", 1:"Montag", 2:"Dienstag", 3:"Mittwoch", 4:"Donnerstag", 5:"Freitag", 6:"Samstag"} return map:remove($week, 4)`,
			[]interface{}{0, 1, 2, 3, 5, 6},
			[]interface{}{"Sonntag", "Montag", "Dienstag", "Mittwoch", "Freitag", "Samstag"},
		},
		{
			`let $week := map{0:"Sonntag", 1:"Montag", 2:"Dienstag", 3:"Mittwoch", 4:"Donnerstag", 5:"Freitag", 6:"Samstag"} return map:remove($week,23)`,
			[]interface{}{0, 1, 2, 3, 4, 5, 6},
			[]interface{}{"Sonntag", "Montag", "Dienstag", "Mittwoch", "Donnerstag", "Freitag", "Samstag"},
		},
		{
			`map:entry("Su", "Sunday")`,
			[]interface{}{"Su"},
			[]interface{}{"Sunday"},
		},
		{
			`let $week := map{0:"Sonntag", 1:"Montag", 2:"Dienstag", 3:"Mittwoch", 4:"Donnerstag", 5:"Freitag", 6:"Samstag"} return map:put($week, -1, "Unbekannt")`,
			[]interface{}{0, 1, 2, 3, 4, 5, 6, -1},
			[]interface{}{"Sonntag", "Montag", "Dienstag", "Mittwoch", "Donnerstag", "Freitag", "Samstag", "Unbekannt"},
		},
		{
			`let $week := map{0:"Sonntag", 1:"Montag", 2:"Dienstag", 3:"Mittwoch", 4:"Donnerstag", 5:"Freitag", 6:"Samstag"} return map:put($week, 6, "Sonnabend")`,
			[]interface{}{0, 1, 2, 3, 4, 5, 6},
			[]interface{}{"Sonntag", "Montag", "Dienstag", "Mittwoch", "Donnerstag", "Freitag", "Sonnabend"},
		},
	}

	for _, tt := range tests {
		seq := testEval(tt.input)
		sequence := seq.(*object.Sequence)
		m := sequence.Items[0].(*object.Map)
		testMapObject(t, m, tt.keys, tt.expected)
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

func TestSequenceTypes(t *testing.T) {
	seq := testEval("5 cast as xs:string")
	sequence := seq.(*object.Sequence)
	if len(sequence.Items) != 1 {
		t.Errorf("wrong number of items. got=%d, expected=1", len(sequence.Items))
	}
	item, ok := sequence.Items[0].(*object.String)
	if !ok {
		t.Errorf("casted item type should be string")
	}
	if item.Value() != "5" {
		t.Errorf("casted item value should be 5. got=%s", item.Value())
	}

	seq2 := testEval("1.23 cast as xs:string")
	sequence2 := seq2.(*object.Sequence)
	if len(sequence2.Items) != 1 {
		t.Errorf("wrong number of items. got=%d, expected=1", len(sequence2.Items))
	}
	item2, ok := sequence2.Items[0].(*object.String)
	if !ok {
		t.Errorf("casted item type should be string")
	}
	if item2.Value() != "1.23" {
		t.Errorf("casted item value should be 1.23. got=%s", item2.Value())
	}

	seq3 := testEval("'1' castable as xs:boolean")
	sequence3 := seq3.(*object.Sequence)
	if len(sequence3.Items) != 1 {
		t.Errorf("wrong number of items. got=%d, expected=1", len(sequence3.Items))
	}
	item3, ok := sequence3.Items[0].(*object.Boolean)
	if !ok {
		t.Errorf("the result type should be boolean")
	}
	if !item3.Value() {
		t.Errorf("the result value should be true")
	}

	seq4 := testEvalXML2("//age[.=25] castable as xs:boolean")
	sequence4 := seq4.(*object.Sequence)
	if len(sequence4.Items) != 1 {
		t.Errorf("wrong number of items. got=%d, expected=1", len(sequence4.Items))
	}
	item4, ok := sequence4.Items[0].(*object.Boolean)
	if !ok {
		t.Errorf("the result type should be boolean")
	}
	if item4.Value() {
		t.Errorf("the result value should be false")
	}

	seq5 := testEvalXML2("'5' instance of xs:string")
	sequence5 := seq5.(*object.Sequence)
	if len(sequence5.Items) != 1 {
		t.Errorf("wrong number of items. got=%d, expected=1", len(sequence5.Items))
	}
	item5, ok := sequence5.Items[0].(*object.Boolean)
	if !ok {
		t.Errorf("the result type should be boolean")
	}
	if !item5.Value() {
		t.Errorf("the result value should be true")
	}

	seq6 := testEvalXML2("5 instance of xs:decimal")
	sequence6 := seq6.(*object.Sequence)
	if len(sequence6.Items) != 1 {
		t.Errorf("wrong number of items. got=%d, expected=1", len(sequence6.Items))
	}
	item6, ok := sequence6.Items[0].(*object.Boolean)
	if !ok {
		t.Errorf("the result type should be boolean")
	}
	if !item6.Value() {
		t.Errorf("the result value should be true")
	}

	seq7 := testEvalXML2("('a') instance of xs:string")
	sequence7 := seq7.(*object.Sequence)
	if len(sequence7.Items) != 1 {
		t.Errorf("wrong number of items. got=%d, expected=1", len(sequence7.Items))
	}
	item7, ok := sequence7.Items[0].(*object.Boolean)
	if !ok {
		t.Errorf("the result type should be boolean")
	}
	if !item7.Value() {
		t.Errorf("the result value should be true")
	}

	seq8 := testEvalXML2("('a','b','c') instance of xs:string")
	sequence8 := seq8.(*object.Sequence)
	if len(sequence8.Items) != 1 {
		t.Errorf("wrong number of items. got=%d, expected=1", len(sequence8.Items))
	}
	item8, ok := sequence8.Items[0].(*object.Boolean)
	if !ok {
		t.Errorf("the result type should be boolean")
	}
	if item8.Value() {
		t.Errorf("the result value should be false")
	}

	seq9 := testEvalXML2("('a','b','c') instance of xs:string+")
	sequence9 := seq9.(*object.Sequence)
	if len(sequence9.Items) != 1 {
		t.Errorf("wrong number of items. got=%d, expected=1", len(sequence9.Items))
	}
	item9, ok := sequence9.Items[0].(*object.Boolean)
	if !ok {
		t.Errorf("the result type should be boolean")
	}
	if !item9.Value() {
		t.Errorf("the result value should be true")
	}

	seq10 := testEvalXML2("('a','b','c') instance of xs:string*")
	sequence10 := seq10.(*object.Sequence)
	if len(sequence10.Items) != 1 {
		t.Errorf("wrong number of items. got=%d, expected=1", len(sequence10.Items))
	}
	item10, ok := sequence10.Items[0].(*object.Boolean)
	if !ok {
		t.Errorf("the result type should be boolean")
	}
	if !item10.Value() {
		t.Errorf("the result value should be true")
	}

	seq11 := testEvalXML2("() instance of xs:string?")
	sequence11 := seq11.(*object.Sequence)
	if len(sequence11.Items) != 1 {
		t.Errorf("wrong number of items. got=%d, expected=1", len(sequence11.Items))
	}
	item11, ok := sequence11.Items[0].(*object.Boolean)
	if !ok {
		t.Errorf("the result type should be boolean")
	}
	if !item11.Value() {
		t.Errorf("the result value should be true")
	}

	seq12 := testEvalXML2("(1,2,3) instance of xs:integer+")
	sequence12 := seq12.(*object.Sequence)
	if len(sequence12.Items) != 1 {
		t.Errorf("wrong number of items. got=%d, expected=1", len(sequence12.Items))
	}
	item12, ok := sequence12.Items[0].(*object.Boolean)
	if !ok {
		t.Errorf("the result type should be boolean")
	}
	if !item12.Value() {
		t.Errorf("the result value should be true")
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
	if !bif.IsError(seq) {
		t.Errorf(seq.Inspect())
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
	if node38.Inspect() != "Attr{cover:paperback}" {
		t.Errorf("wrong attr value. got=%s, expected='Attr{cover:paperback}'", node38.Inspect())
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

	seq63 := testEvalXML2("//office[1]/@location lt //office[2]/@location")
	sequence63 := seq63.(*object.Sequence)
	if len(sequence63.Items) != 1 {
		t.Errorf("wrong number of items. got=%d, expected=1", len(sequence63.Items))
	}
	item63 := sequence63.Items[0].(*object.Boolean)
	if item63.Value() {
		t.Errorf("the result should be [false]")
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
	if node4.Inspect() != "Attr{category:web}" {
		t.Errorf("wrong attribute value. got=%s, expected='Attr{category:web}'", node4.Inspect())
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
	if node43.Inspect() != "Attr{category:1}" {
		t.Errorf("wrong attribute value. got=%s, expected='Attr{category:1}'", node43.Inspect())
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
	if node44.Inspect() != "Attr{category:2}" {
		t.Errorf("wrong attribute value. got=%s, expected='Attr{category:2}'", node44.Inspect())
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
	if attr49.Key() != "category" || attr49.Inspect() != "Attr{category:2}" {
		t.Errorf("expected attr: %s", attr49.Inspect())
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
	if attr50.Key() != "category" || attr50.Inspect() != "Attr{category:1}" {
		t.Errorf("expected attr: %s", attr50.Inspect())
	}

	seq51 := testEvalXML2("//age[.=25]")
	sequence51 := seq51.(*object.Sequence)
	if len(sequence51.Items) != 1 {
		t.Errorf("wrong number of items. got=%d, expected=1", len(sequence51.Items))
	}
	item51, ok := sequence51.Items[0].(*object.BaseNode)
	if !ok {
		t.Errorf("node type should be BaseNode")
	}
	if item51.Text() != "25" {
		t.Errorf("node text should be 25. got=%s", item51.Text())
	}

	seq52 := testEvalXML2("//age[.>'25']")
	sequence52 := seq52.(*object.Sequence)
	if len(sequence52.Items) != 4 {
		t.Errorf("wrong number of items. got=%d, expected=4", len(sequence52.Items))
	}

	seq53 := testEvalXML2("//age[.>25]")
	sequence53 := seq53.(*object.Sequence)
	if len(sequence53.Items) != 4 {
		t.Errorf("wrong number of items. got=%d, expected=4", len(sequence53.Items))
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

func TestBIF(t *testing.T) {
	seq := testEvalXML2("//employee/node-name()")
	sequence := seq.(*object.Sequence)
	if len(sequence.Items) != 5 {
		t.Errorf("wrong number of items. got=%d, expected=5", len(sequence.Items))
	}

	seq2 := testEvalXML2("//employee/*/node-name()")
	sequence2 := seq2.(*object.Sequence)
	if len(sequence2.Items) != 15 {
		t.Errorf("wrong number of items. got=%d, expected=15", len(sequence2.Items))
	}
	item2 := sequence2.Items[0].(*object.String)
	if item2.Value() != "first_name" {
		t.Errorf("node name should be first_name")
	}

	seq3 := testEvalXML2("//employee/string()")
	sequence3 := seq3.(*object.Sequence)
	if len(sequence3.Items) != 5 {
		t.Errorf("wrong number of items. got=%d, expected=5", len(sequence3.Items))
	}

	seq4 := testEvalXML2("//employee/data()")
	sequence4 := seq4.(*object.Sequence)
	if len(sequence4.Items) != 5 {
		t.Errorf("wrong number of items. got=%d, expected=5", len(sequence4.Items))
	}
	item4 := sequence4.Items[0].(*object.String)
	if item4.Value() != "ChoiJack25" {
		t.Errorf("first item value should be ChoiJack25. got=%s", item4.Value())
	}

	seq5 := testEvalXML2("//employee/string('haha')")
	sequence5 := seq5.(*object.Sequence)
	if len(sequence5.Items) != 5 {
		t.Errorf("wrong number of items. got=%d, expected=5", len(sequence5.Items))
	}
	item5 := sequence5.Items[0].(*object.String)
	if item5.Value() != "haha" {
		t.Errorf("first item value should be haha. got=%s", item5.Value())
	}

	seq6 := testEvalXML2("/base-uri()")
	sequence6 := seq6.(*object.Sequence)
	if len(sequence6.Items) != 1 {
		t.Errorf("wrong number of items. got=%d, expected=1", len(sequence6.Items))
	}

	seq7 := testEvalXML2("fn:ceiling(1.1)")
	item7 := seq7.(*object.Sequence).Items[0].(*object.Decimal)
	if item7.Value() != 2 {
		t.Errorf("result value should be 2. got=%f", item7.Value())
	}

	seq8 := testEvalXML2("floor(1.8e0)")
	item8 := seq8.(*object.Sequence).Items[0].(*object.Double)
	if item8.Value() != 1 {
		t.Errorf("result value should be 1. got=%f", item8.Value())
	}

	seq9 := testEvalXML2("fn:round(1.5)")
	item9 := seq9.(*object.Sequence).Items[0].(*object.Decimal)
	if item9.Value() != 2 {
		t.Errorf("result value should be 2. got=%f", item9.Value())
	}

	seq10 := testEvalXML2("round-half-to-even(12.5)")
	item10 := seq10.(*object.Sequence).Items[0].(*object.Decimal)
	if item10.Value() != 12 {
		t.Errorf("result value should be 12. got=%f", item10.Value())
	}

	seq11 := testEvalXML2("/number()")
	sequence11 := seq11.(*object.Sequence)
	if len(sequence11.Items) != 1 {
		t.Errorf("wrong number of items. got=%d, expected=1", len(sequence11.Items))
	}
	item11 := sequence11.Items[0].(*object.Double)
	if !math.IsNaN(item11.Value()) {
		t.Errorf("item value should be NaN")
	}

	seq12 := testEvalXML2("//employee/fn:number(1)")
	sequence12 := seq12.(*object.Sequence)
	if len(sequence12.Items) != 5 {
		t.Errorf("wrong number of items. got=%d, expected=5", len(sequence12.Items))
	}
	item12 := sequence12.Items[0].(*object.Double)
	if item12.Value() != 1 {
		t.Errorf("first item value should be 1")
	}

	seq13 := testEvalXML2("//age/number()")
	sequence13 := seq13.(*object.Sequence)
	if len(sequence13.Items) != 5 {
		t.Errorf("wrong number of items. got=%d, expected=5", len(sequence13.Items))
	}
	item13 := sequence13.Items[0].(*object.Double)
	if item13.Value() != 25 {
		t.Errorf("first item value should be 25")
	}

	seq14 := testEvalXML2("(//age/number(), number(1))")
	sequence14 := seq14.(*object.Sequence)
	if len(sequence14.Items) != 6 {
		t.Errorf("wrong number of items. got=%d, expected=6", len(sequence14.Items))
	}
	item14_1 := sequence14.Items[0].(*object.Double)
	if item14_1.Value() != 25 {
		t.Errorf("first item value should be 25")
	}
	item14_2 := sequence14.Items[5].(*object.Double)
	if item14_2.Value() != 1 {
		t.Errorf("last item value should be 1")
	}

	seq15 := testEvalXML2("(//age/number(), number(1), (5,6, (7,8)))")
	sequence15 := seq15.(*object.Sequence)
	if len(sequence15.Items) != 10 {
		t.Errorf("wrong number of items. got=%d, expected=10", len(sequence15.Items))
	}

	seq16 := testEvalXML2("(//age/number(), number(1), //employee/number())")
	sequence16 := seq16.(*object.Sequence)
	if len(sequence16.Items) != 11 {
		t.Errorf("wrong number of items. got=%d, expected=11", len(sequence16.Items))
	}
	item16_1 := sequence16.Items[5].(*object.Double)
	if item16_1.Value() != 1 {
		t.Errorf("item value should be 1")
	}
	item16_2 := sequence16.Items[10].(*object.Double)
	if !math.IsNaN(item16_2.Value()) {
		t.Errorf("item value should be nan")
	}

	seq17 := testEvalXML2("math:sqrt(2)")
	item17 := seq17.(*object.Sequence).Items[0].(*object.Double)
	if fmt.Sprintf("%.3f", item17.Value()) != "1.414" {
		t.Errorf("item value should be 1.414. got=%.3f", item17.Value())
	}

	seq18 := testEvalXML2("let $abc := ('a', 'b', '') return fn:boolean($abc)")
	if !bif.IsError(seq18) {
		t.Errorf("the result should be error")
	}

	seq19 := testEvalXML2("let $abc := ('a', 'b', '') return fn:boolean([])")
	if !bif.IsError(seq19) {
		t.Errorf("the result should be error")
	}

	seq20 := testEvalXML2("let $abc := ('a', 'b', '') return fn:boolean($abc[1])")
	item20, ok := seq20.(*object.Sequence).Items[0].(*object.Boolean)
	if !ok {
		t.Errorf("the result type should be boolean")
	}
	if !item20.Value() {
		t.Errorf("the result value should be true")
	}

	seq21 := testEvalXML2("let $abc := ('a', 'b', '') return fn:boolean($abc[0])")
	item21, ok := seq21.(*object.Sequence).Items[0].(*object.Boolean)
	if !ok {
		t.Errorf("the result type should be boolean")
	}
	if item21.Value() {
		t.Errorf("the result value should be false")
	}

	seq22 := testEvalXML2("let $abc := ('a', 'b', '') return fn:boolean($abc[3])")
	item22, ok := seq22.(*object.Sequence).Items[0].(*object.Boolean)
	if !ok {
		t.Errorf("the result type should be boolean")
	}
	if item22.Value() {
		t.Errorf("the result value should be false")
	}

	seq23 := testEvalXML2("//office/string-length()")
	sequence23 := seq23.(*object.Sequence)
	if len(sequence23.Items) != 2 {
		t.Errorf("wrong number of items. got=%d, expected=2", len(sequence23.Items))
	}

	seq24 := testEvalXML2("//age/string-length()")
	sequence24 := seq24.(*object.Sequence)
	if len(sequence24.Items) != 5 {
		t.Errorf("wrong number of items. got=%d, expected=5", len(sequence24.Items))
	}

	seq25 := testEvalXML2("//office/normalize-space()")
	sequence25 := seq25.(*object.Sequence)
	if len(sequence25.Items) != 2 {
		t.Errorf("wrong number of items. got=%d, expected=2", len(sequence25.Items))
	}
	item25 := sequence25.Items[0].(*object.String)
	if item25.Value() != "Choi Jack 25 Lee Hwa 30" {
		t.Errorf("the result value should be 'Choi Jack 25 Lee Hwa 30'. got=%s", item25.Value())
	}

	seq26 := testEvalXML2("//last()")
	sequence26 := seq26.(*object.Sequence)
	if len(sequence26.Items) != 72 {
		t.Errorf("wrong number of items. got=%d, expected=72", len(sequence26.Items))
	}

	seq27 := testEvalXML2("//age/last()")
	sequence27 := seq27.(*object.Sequence)
	if len(sequence27.Items) != 5 {
		t.Errorf("wrong number of items. got=%d, expected=5", len(sequence27.Items))
	}
	item27 := sequence27.Items[0].(*object.Integer)
	if item27.Value() != 5 {
		t.Errorf("first item value should be 5")
	}

	seq28 := testEvalXML2("//age[position()=1]")
	sequence28 := seq28.(*object.Sequence)
	if len(sequence28.Items) != 5 {
		t.Errorf("wrong number of items. got=%d, expected=5", len(sequence28.Items))
	}
}

func TestPathExpr2(t *testing.T) {
	seq := testEvalXML2(".//company")
	sequence := seq.(*object.Sequence)
	if len(sequence.Items) != 1 {
		t.Errorf("wrong number of items. got=%d, expected=1", len(sequence.Items))
	}
	node := sequence.Items[0].(*object.BaseNode)
	if node.Tree().Data != "company" {
		t.Errorf("selected node name should be company. got=%s", node.Tree().Data)
	}

	seq2 := testEvalXML2("./.")
	sequence2 := seq2.(*object.Sequence)
	if len(sequence2.Items) != 1 {
		t.Errorf("wrong number of items. got=%d, expected=1", len(sequence2.Items))
	}
	node2 := sequence2.Items[0].(*object.BaseNode)
	if node2.Type() != object.DocumentNodeType {
		t.Errorf("selected node type should be DocumentNode")
	}

	seq3 := testEvalXML2(".//company/./office")
	sequence3 := seq3.(*object.Sequence)
	if len(sequence3.Items) != 2 {
		t.Errorf("wrong number of items. got=%d, expected=2", len(sequence3.Items))
	}

	seq4 := testEvalXML2(".//company/../.")
	sequence4 := seq4.(*object.Sequence)
	if len(sequence4.Items) != 1 {
		t.Errorf("wrong number of items. got=%d, expected=1", len(sequence4.Items))
	}
	node4 := sequence4.Items[0].(*object.BaseNode)
	if node4.Tree().Data != "body" {
		t.Errorf("selected node name should be body")
	}
}

func testEval(input string) object.Item {
	l := lexer.New(input)
	p := parser.New(l)
	xpath := p.ParseXPath()
	ctx := object.NewContext()

	if len(p.Errors()) > 0 {
		var sb strings.Builder
		for _, e := range p.Errors() {
			sb.WriteString(e.Error())
		}
		return bif.NewError(sb.String())
	}

	docFunc := bif.F["fn:doc"]
	err := docFunc(ctx, bif.NewString("testdata/quotes-1.html"))
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

	if len(p.Errors()) > 0 {
		var sb strings.Builder
		for _, e := range p.Errors() {
			sb.WriteString(e.Error())
		}
		return bif.NewError(sb.String())
	}

	docFunc := bif.F["fn:doc"]
	err := docFunc(ctx, bif.NewString("testdata/company.xml"))
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

	if len(p.Errors()) > 0 {
		var sb strings.Builder
		for _, e := range p.Errors() {
			sb.WriteString(e.Error())
		}
		return bif.NewError(sb.String())
	}

	docFunc := bif.F["fn:doc"]
	err := docFunc(ctx, bif.NewString("testdata/company_2.xml"))
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
		t.Errorf("Unknown item type. got=%s", item.Type())
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

// just checking map value, map key should be one type of primitive
func testMapObject(t *testing.T, m *object.Map, key, expected interface{}) {
	ex := expected.([]interface{})

	switch ke := key.(type) {
	case []interface{}:
		for i, it := range ke {
			switch e := it.(type) {
			case int:
				key := bif.NewInteger(e).HashKey()
				pair, ok := m.Pairs[key]
				if !ok {
					t.Errorf("key [%d] doesn't exiest in map", e)
				}

				if exi, ok := ex[i].([]interface{}); ok {
					if !bif.IsSeq(pair.Value) {
						t.Errorf("value should be sequence type")
					}
					seq := pair.Value.(*object.Sequence)
					if len(seq.Items) != len(exi) {
						t.Errorf("wrong number of items")
					}
					for j, eexi := range exi {
						if !bif.IsSameAtomic(seq.Items[j], eexi) {
							t.Errorf("value doesn't match: got=%v, expected=%v", seq.Items[j], ex[i])
						}
					}
				} else if !bif.IsSameAtomic(pair.Value, ex[i]) {
					t.Errorf("value doesn't match: got=%v, expected=%v", pair.Value, ex[i])
				}
			case float64:
				key := bif.NewDecimal(e).HashKey()
				pair, ok := m.Pairs[key]
				if !ok {
					t.Errorf("key [%f] doesn't exiest in map", e)
				}

				if exi, ok := ex[i].([]interface{}); ok {
					if !bif.IsSeq(pair.Value) {
						t.Errorf("value should be sequence type")
					}
					seq := pair.Value.(*object.Sequence)
					if len(seq.Items) != len(exi) {
						t.Errorf("wrong number of items")
					}
					for j, eexi := range exi {
						if !bif.IsSameAtomic(seq.Items[j], eexi) {
							t.Errorf("value doesn't match: got=%v, expected=%v", seq.Items[j], ex[i])
						}
					}
				} else if !bif.IsSameAtomic(pair.Value, ex[i]) {
					t.Errorf("value doesn't match: got=%v, expected=%v", pair.Value, ex[i])
				}
			case string:
				key := bif.NewString(e).HashKey()
				pair, ok := m.Pairs[key]
				if !ok {
					t.Errorf("key [%s] doesn't exiest in map", e)
				}

				if exi, ok := ex[i].([]interface{}); ok {
					if !bif.IsSeq(pair.Value) {
						t.Errorf("value should be sequence type")
					}
					seq := pair.Value.(*object.Sequence)
					if len(seq.Items) != len(exi) {
						t.Errorf("wrong number of items")
					}
					for j, eexi := range exi {
						if !bif.IsSameAtomic(seq.Items[j], eexi) {
							t.Errorf("value doesn't match: got=%v, expected=%v", seq.Items[j], ex[i])
						}
					}
				} else if !bif.IsSameAtomic(pair.Value, ex[i]) {
					t.Errorf("value doesn't match: got=%v, expected=%v", pair.Value, ex[i])
				}
			case bool:
				key := bif.NewBoolean(e).HashKey()
				pair, ok := m.Pairs[key]
				if !ok {
					t.Errorf("key [%v] doesn't exiest in map", e)
				}

				if exi, ok := ex[i].([]interface{}); ok {
					if !bif.IsSeq(pair.Value) {
						t.Errorf("value should be sequence type")
					}
					seq := pair.Value.(*object.Sequence)
					if len(seq.Items) != len(exi) {
						t.Errorf("wrong number of items")
					}
					for j, eexi := range exi {
						if !bif.IsSameAtomic(seq.Items[j], eexi) {
							t.Errorf("value doesn't match: got=%v, expected=%v", seq.Items[j], ex[i])
						}
					}
				} else if !bif.IsSameAtomic(pair.Value, ex[i]) {
					t.Errorf("value doesn't match: got=%v, expected=%v", pair.Value, ex[i])
				}
			}
		}
	}
}

func testArrayObject(t *testing.T, item object.Item, expected interface{}) {
	itemArr, ok := item.(*object.Array)
	if !ok {
		t.Errorf("item type should be object.Array")
	}

	switch e := expected.(type) {
	case []int:
		if len(e) != len(itemArr.Items) {
			t.Errorf("array length not match. got=%d, expected=%d", len(itemArr.Items), len(e))
		}
		for i, it := range itemArr.Items {
			it, ok := it.(*object.Integer)
			if !ok {
				t.Errorf("item type should be integer")
			}
			if e[i] != it.Value() {
				t.Errorf("item value not match. got=%d, expected=%d", it.Value(), e[i])
			}
		}
	case []string:
		if len(e) != len(itemArr.Items) {
			t.Errorf("array length not match. got=%d, expected=%d", len(itemArr.Items), len(e))
		}
		for i, it := range itemArr.Items {
			it, ok := it.(*object.String)
			if !ok {
				t.Errorf("item type should be string")
			}
			if e[i] != it.Value() {
				t.Errorf("item value not match. got=%s, expected=%s", it.Value(), e[i])
			}
		}
	case []float64:
		if len(e) != len(itemArr.Items) {
			t.Errorf("array length not match. got=%d, expected=%d", len(itemArr.Items), len(e))
		}
		for i, it := range itemArr.Items {
			it, ok := it.(*object.Decimal)
			if !ok {
				t.Errorf("item type should be decimal")
			}
			if e[i] != it.Value() {
				t.Errorf("item value not match. got=%f, expected=%f", it.Value(), e[i])
			}
		}
	default:
		t.Errorf("cannot compare array")
	}
}

func testSequenceObject(t *testing.T, item object.Item, expected []interface{}) {
	switch item := item.(type) {
	case *object.Sequence:
		if len(item.Items) != len(expected) {
			t.Errorf("length of the item must be the same. got=%d, want=%d", len(item.Items), len(expected))
		}
		for i := 0; i < len(item.Items); i++ {
			switch it := item.Items[i].(type) {
			case *object.Integer:
				testNumberObject(t, it, expected[i])
			case *object.Decimal:
				testNumberObject(t, it, expected[i])
			case *object.Double:
				testNumberObject(t, it, expected[i])
			case *object.String:
				testStringObject(t, it, expected[i])
			case *object.Boolean:
				if it.Value() != expected[i] {
					t.Errorf("got=%t, expected=%t", it.Value(), expected[i])
				}
			case *object.Array:
				testArrayObject(t, it, expected[i])
			default:
				t.Errorf("Unknown item type. got=%s, %s", it.Type(), it.Inspect())
			}
		}
	default:
		t.Errorf("item type must object.String. got=%s", item.Type())
	}
}
