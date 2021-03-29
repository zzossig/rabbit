package bif

import (
	"fmt"
	"strconv"

	"github.com/zzossig/xpath/ast"
	"github.com/zzossig/xpath/object"
	"github.com/zzossig/xpath/token"
)

// F ...
var F = map[string]object.Func{
	// 2
	"fn:node-name": fnNodeName,
	"fn:string":    fnString,
	"fn:data":      fnData,
	"fn:base-uri":  fnBaseURI,

	// 4.2
	"op:numeric-add":            numericAdd,
	"op:numeric-subtract":       numericSubtract,
	"op:numeric-multiply":       numericMultiply,
	"op:numeric-divide":         numericDivide,
	"op:numeric-integer-divide": numericIntegerDivide,
	"op:numeric-mod":            numericMod,
	"op:numeric-unary-plus":     numericUnaryPlus,
	"op:numeric-unary-minus":    numericUnaryMinus,

	// 4.3
	"op:numeric-equal":        opNumericEqual,
	"op:numeric-less-than":    opNumericLessThan,
	"op:numeric-greater-than": opNumericGreaterThan,

	// 4.4
	"fn:abs":                fnAbs,
	"fn:ceiling":            fnCeiling,
	"fn:floor":              fnFloor,
	"fn:round":              fnRound,
	"fn:round-half-to-even": fnRoundHTE,

	// 4.5
	"fn:number": fnNumber,

	// 4.8
	"math:pi":    mathPI,
	"math:exp":   mathExp,
	"math:exp2":  mathExp2,
	"math:log":   mathLog,
	"math:log2":  mathLog2,
	"math:log10": mathLog10,
	"math:pow":   mathPow,
	"math:sqrt":  mathSqrt,
	"math:sin":   mathSin,
	"math:cos":   mathCos,
	"math:tan":   mathTan,
	"math:asin":  mathAsin,
	"math:acos":  mathAcos,
	"math:atan":  mathAtan,
	"math:atan2": mathAtan2,

	// 5.4
	"fn:concat":          fnConcat,
	"fn:string-join":     fnStringJoin,
	"fn:substring":       fnSubstring,
	"fn:string-length":   fnStringLength,
	"fn:normalize-space": fnNormalizeSpace,
	"fn:upper-case":      fnUpperCase,
	"fn:lower-case":      fnLowerCase,

	// 5.5
	"fn:contains":         fnContains,
	"fn:starts-with":      fnStartsWith,
	"fn:ends-with":        fnEndsWith,
	"fn:substring-before": fnSubstringBefore,
	"fn:substring-after":  fnSubstringAfter,

	// 7
	"fn:true":    fnTrue,
	"fn:false":   fnFalse,
	"fn:boolean": fnBoolean,
	"fn:not":     fnNot,

	"op:boolean-equal":        opBooleanEqual,
	"op:boolean-less-than":    opBooleanLessThan,
	"op:boolean-greater-than": opBooleanGreaterThan,

	// 14.1
	"fn:empty":         fnEmpty,
	"fn:exists":        fnExists,
	"fn:head":          fnHead,
	"fn:tail":          fnTail,
	"fn:insert-before": fnInsertBefore,
	"fn:remove":        fnRemove,
	"fn:reverse":       fnReverse,
	"fn:subsequence":   fnSubsequence,

	// 14
	"fn:doc": fnDoc,

	// 16.2
	"fn:for-each-pair": fnForEachPair,
}

// NewError ..
func NewError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

// NewString ..
func NewString(s string) *object.String {
	str := &object.String{}
	str.SetValue(s)
	return str
}

// NewBoolean ..
func NewBoolean(b bool) *object.Boolean {
	boolean := &object.Boolean{}
	boolean.SetValue(b)
	return boolean
}

// NewInteger ..
func NewInteger(i int) *object.Integer {
	integer := &object.Integer{}
	integer.SetValue(i)
	return integer
}

// NewDecimal ..
func NewDecimal(d float64) *object.Decimal {
	decimal := &object.Decimal{}
	decimal.SetValue(d)
	return decimal
}

// NewDouble ..
func NewDouble(d float64) *object.Double {
	double := &object.Double{}
	double.SetValue(d)
	return double
}

// NewSequence ..
func NewSequence(items ...object.Item) *object.Sequence {
	seq := &object.Sequence{}
	if len(items) > 0 {
		seq.Items = append(seq.Items, items...)
	}
	return seq
}

// IsError ..
func IsError(item object.Item) bool {
	if item == nil {
		return false
	}
	return item.Type() == object.ErrorType
}

// IsItem ..
func IsItem(item object.Item) bool {
	if item == nil {
		return false
	}
	return IsAnyAtomic(item) || IsNode(item) || IsAnyFunc(item)
}

// IsSeq ..
func IsSeq(item object.Item) bool {
	if item == nil {
		return false
	}
	return item.Type() == object.SequenceType
}

// IsSeqEmpty ..
func IsSeqEmpty(item object.Item) bool {
	if item == nil {
		return false
	}
	if item.Type() != object.SequenceType {
		return false
	}
	return item.Inspect() == "()"
}

// IsPlaceholder ..
func IsPlaceholder(item object.Item) bool {
	if item == nil {
		return false
	}
	return item.Type() == object.PholderType
}

// IsNumeric ..
func IsNumeric(item object.Item) bool {
	if item == nil {
		return false
	}
	return item.Type() == object.IntegerType ||
		item.Type() == object.DecimalType ||
		item.Type() == object.DoubleType
}

// IsFunc ..
func IsFunc(item object.Item) bool {
	if item == nil {
		return false
	}
	return item.Type() == object.FuncType
}

// IsAnyAtomic ..
func IsAnyAtomic(item object.Item) bool {
	if item == nil {
		return false
	}
	return item.Type() == object.DoubleType ||
		item.Type() == object.DecimalType ||
		item.Type() == object.IntegerType ||
		item.Type() == object.StringType ||
		item.Type() == object.BooleanType
}

// IsAnyFunc ..
func IsAnyFunc(item object.Item) bool {
	if item == nil {
		return false
	}
	return item.Type() == object.FuncType ||
		item.Type() == object.MapType ||
		item.Type() == object.ArrayType
}

// IsNode ..
func IsNode(item object.Item) bool {
	if item == nil {
		return false
	}
	return item.Type() == object.TextNodeType ||
		item.Type() == object.DocumentNodeType ||
		item.Type() == object.ElementNodeType ||
		item.Type() == object.CommentNodeType ||
		item.Type() == object.AttributeNodeType
}

// IsMap ..
func IsMap(item object.Item) bool {
	if item == nil {
		return false
	}
	return item.Type() == object.MapType
}

// IsArray ..
func IsArray(item object.Item) bool {
	if item == nil {
		return false
	}
	return item.Type() == object.ArrayType
}

// IsString ..
func IsString(item object.Item) bool {
	if item == nil {
		return false
	}
	return item.Type() == object.StringType
}

// IsBoolean ..
func IsBoolean(item object.Item) bool {
	if item == nil {
		return false
	}
	return item.Type() == object.BooleanType
}

// IsItemSeq ..
func IsItemSeq(item object.Item) bool {
	seq, ok := item.(*object.Sequence)
	if !ok {
		return false
	}

	for _, i := range seq.Items {
		if !IsItem(i) {
			return false
		}
	}

	return true
}

// IsNodeSeq ..
func IsNodeSeq(item object.Item) bool {
	seq, ok := item.(*object.Sequence)
	if !ok {
		return false
	}

	for _, i := range seq.Items {
		if !IsNode(i) {
			return false
		}
	}

	return true
}

// IsAnyFuncSeq ..
func IsAnyFuncSeq(item object.Item) bool {
	seq, ok := item.(*object.Sequence)
	if !ok {
		return false
	}

	for _, i := range seq.Items {
		if !IsAnyFunc(i) {
			return false
		}
	}

	return true
}

// IsMapSeq ..
func IsMapSeq(item object.Item) bool {
	seq, ok := item.(*object.Sequence)
	if !ok {
		return false
	}

	for _, i := range seq.Items {
		if !IsMap(i) {
			return false
		}
	}

	return true
}

// IsArraySeq ..
func IsArraySeq(item object.Item) bool {
	seq, ok := item.(*object.Sequence)
	if !ok {
		return false
	}

	for _, i := range seq.Items {
		if !IsArray(i) {
			return false
		}
	}

	return true
}

// IsAtomicSeq ..
func IsAtomicSeq(item object.Item) bool {
	seq, ok := item.(*object.Sequence)
	if !ok {
		return false
	}

	for _, i := range seq.Items {
		if !IsAnyAtomic(i) {
			return false
		}
	}

	return true
}

// IsDocNodeSeq ..
func IsDocNodeSeq(item object.Item) bool {
	seq, ok := item.(*object.Sequence)
	if !ok {
		return false
	}

	for _, i := range seq.Items {
		if i.Type() != object.DocumentNodeType {
			return false
		}
	}

	return true
}

// IsElemNodeSeq ..
func IsElemNodeSeq(item object.Item) bool {
	seq, ok := item.(*object.Sequence)
	if !ok {
		return false
	}

	for _, i := range seq.Items {
		if i.Type() != object.ElementNodeType {
			return false
		}
	}

	return true
}

// IsAttrNodeSeq ..
func IsAttrNodeSeq(item object.Item) bool {
	seq, ok := item.(*object.Sequence)
	if !ok {
		return false
	}

	for _, i := range seq.Items {
		if i.Type() != object.AttributeNodeType {
			return false
		}
	}

	return true
}

// IsCommNodeSeq ..
func IsCommNodeSeq(item object.Item) bool {
	seq, ok := item.(*object.Sequence)
	if !ok {
		return false
	}

	for _, i := range seq.Items {
		if i.Type() != object.CommentNodeType {
			return false
		}
	}

	return true
}

// IsTextNodeSeq ..
func IsTextNodeSeq(item object.Item) bool {
	seq, ok := item.(*object.Sequence)
	if !ok {
		return false
	}

	for _, i := range seq.Items {
		if i.Type() != object.TextNodeType {
			return false
		}
	}

	return true
}

// IsCastable ..
func IsCastable(tg object.Item, ty object.Type) object.Item {
	switch tg := tg.(type) {
	case *object.Sequence:
		if len(tg.Items) != 1 {
			return NewError("wrong number of sequence items. got=%d, expected=1", len(tg.Items))
		}
		return IsCastable(tg.Items[0], ty)
	case *object.Double:
		switch ty {
		case object.DoubleType:
			fallthrough
		case object.DecimalType:
			fallthrough
		case object.IntegerType:
			fallthrough
		case object.StringType:
			fallthrough
		case object.BooleanType:
			return NewBoolean(true)
		}
		return NewBoolean(false)
	case *object.Decimal:
		switch ty {
		case object.DoubleType:
			fallthrough
		case object.DecimalType:
			fallthrough
		case object.IntegerType:
			fallthrough
		case object.StringType:
			fallthrough
		case object.BooleanType:
			return NewBoolean(true)
		}
		return NewBoolean(false)
	case *object.Integer:
		switch ty {
		case object.DoubleType:
			fallthrough
		case object.DecimalType:
			fallthrough
		case object.IntegerType:
			fallthrough
		case object.StringType:
			fallthrough
		case object.BooleanType:
			return NewBoolean(true)
		}
		return NewBoolean(false)
	case *object.String:
		switch ty {
		case object.DoubleType:
			fallthrough
		case object.DecimalType:
			if _, err := strconv.ParseFloat(tg.Value(), 64); err == nil {
				return NewBoolean(true)
			}
			return NewBoolean(false)
		case object.IntegerType:
			if _, err := strconv.ParseInt(tg.Value(), 0, 64); err == nil {
				return NewBoolean(true)
			}
			return NewBoolean(false)
		case object.StringType:
			return NewBoolean(true)
		case object.BooleanType:
			if tg.Value() == "1" ||
				tg.Value() == "true" ||
				tg.Value() == "0" ||
				tg.Value() == "false" {
				return NewBoolean(true)
			}
			return NewBoolean(false)
		}
		return NewBoolean(false)
	case *object.Boolean:
		switch ty {
		case object.DoubleType:
			fallthrough
		case object.DecimalType:
			fallthrough
		case object.IntegerType:
			fallthrough
		case object.StringType:
			fallthrough
		case object.BooleanType:
			return NewBoolean(true)
		}
		return NewBoolean(false)
	case *object.BaseNode:
		switch ty {
		case object.DoubleType:
			fallthrough
		case object.DecimalType:
			if _, err := strconv.ParseFloat(tg.Text(), 64); err == nil {
				return NewBoolean(true)
			}
			return NewBoolean(false)
		case object.IntegerType:
			if _, err := strconv.ParseInt(tg.Text(), 0, 64); err == nil {
				return NewBoolean(true)
			}
			return NewBoolean(false)
		case object.StringType:
			return NewBoolean(true)
		case object.BooleanType:
			if tg.Text() == "1" ||
				tg.Text() == "true" ||
				tg.Text() == "0" ||
				tg.Text() == "false" {
				return NewBoolean(true)
			}
			return NewBoolean(false)
		}
		return NewBoolean(false)
	case *object.AttrNode:
		switch ty {
		case object.DoubleType:
			fallthrough
		case object.DecimalType:
			if _, err := strconv.ParseFloat(tg.Text(), 64); err == nil {
				return NewBoolean(true)
			}
			return NewBoolean(false)
		case object.IntegerType:
			if _, err := strconv.ParseInt(tg.Text(), 0, 64); err == nil {
				return NewBoolean(true)
			}
			return NewBoolean(false)
		case object.StringType:
			return NewBoolean(true)
		case object.BooleanType:
			if tg.Text() == "1" ||
				tg.Text() == "true" ||
				tg.Text() == "0" ||
				tg.Text() == "false" {
				return NewBoolean(true)
			}
			return NewBoolean(false)
		}
		return NewBoolean(false)
	}

	return NewBoolean(false)
}

// CastType ..
func CastType(tg object.Item, ty object.Type) object.Item {
	bl := IsCastable(tg, ty)
	if IsError(bl) {
		return bl
	}

	blObj := bl.(*object.Boolean)
	if !blObj.Value() {
		return NewError("cannot convert %s with value %s to %s", tg.Type(), tg.Inspect(), ty)
	}

	switch tg := tg.(type) {
	case *object.Sequence:
		return CastType(tg.Items[0], ty)
	case *object.Double:
		switch ty {
		case object.DoubleType:
			return tg
		case object.DecimalType:
			return NewDecimal(tg.Value())
		case object.IntegerType:
			return NewInteger(int(tg.Value()))
		case object.StringType:
			return NewString(strconv.FormatFloat(tg.Value(), 'f', -1, 64))
		case object.BooleanType:
			return NewBoolean(tg.Value() != 0)
		}
	case *object.Decimal:
		switch ty {
		case object.DoubleType:
			return NewDouble(tg.Value())
		case object.DecimalType:
			return tg
		case object.IntegerType:
			return NewInteger(int(tg.Value()))
		case object.StringType:
			return NewString(strconv.FormatFloat(tg.Value(), 'f', -1, 64))
		case object.BooleanType:
			return NewBoolean(tg.Value() != 0)
		}
	case *object.Integer:
		switch ty {
		case object.DoubleType:
			return NewDouble(float64(tg.Value()))
		case object.DecimalType:
			return NewDecimal(float64(tg.Value()))
		case object.IntegerType:
			return tg
		case object.StringType:
			return NewString(strconv.FormatInt(int64(tg.Value()), 10))
		case object.BooleanType:
			return NewBoolean(tg.Value() != 0)
		}
	case *object.String:
		switch ty {
		case object.DoubleType:
			if i, err := strconv.ParseFloat(tg.Value(), 64); err == nil {
				return NewDouble(i)
			}
		case object.DecimalType:
			if i, err := strconv.ParseFloat(tg.Value(), 64); err == nil {
				return NewDecimal(i)
			}
		case object.IntegerType:
			if i, err := strconv.ParseInt(tg.Value(), 0, 64); err == nil {
				return NewInteger(int(i))
			}
		case object.StringType:
			return tg
		case object.BooleanType:
			if tg.Value() == "0" || tg.Value() == "false" {
				return NewBoolean(false)
			}
			if tg.Value() == "1" || tg.Value() == "true" {
				return NewBoolean(true)
			}
		}
	case *object.Boolean:
		switch ty {
		case object.DoubleType:
			if tg.Value() {
				return NewDouble(1)
			}
			return NewDouble(0)
		case object.DecimalType:
			if tg.Value() {
				return NewDecimal(1)
			}
			return NewDecimal(0)
		case object.IntegerType:
			if tg.Value() {
				return NewInteger(1)
			}
			return NewInteger(0)
		case object.StringType:
			if tg.Value() {
				return NewString("true")
			}
			return NewString("false")
		case object.BooleanType:
			return tg
		}
	case *object.BaseNode:
		switch ty {
		case object.DoubleType:
			if i, err := strconv.ParseFloat(tg.Text(), 64); err == nil {
				return NewDouble(i)
			}
		case object.DecimalType:
			if i, err := strconv.ParseFloat(tg.Text(), 64); err == nil {
				return NewDecimal(i)
			}
		case object.IntegerType:
			if i, err := strconv.ParseInt(tg.Text(), 0, 64); err == nil {
				return NewInteger(int(i))
			}
		case object.StringType:
			return NewString(tg.Text())
		case object.BooleanType:
			if tg.Text() == "0" || tg.Text() == "false" {
				return NewBoolean(false)
			}
			if tg.Text() == "1" || tg.Text() == "true" {
				return NewBoolean(true)
			}
		}
	case *object.AttrNode:
		switch ty {
		case object.DoubleType:
			if i, err := strconv.ParseFloat(tg.Text(), 64); err == nil {
				return NewDouble(i)
			}
		case object.DecimalType:
			if i, err := strconv.ParseFloat(tg.Text(), 64); err == nil {
				return NewDecimal(i)
			}
		case object.IntegerType:
			if i, err := strconv.ParseInt(tg.Text(), 0, 64); err == nil {
				return NewInteger(int(i))
			}
		case object.StringType:
			return NewString(tg.Text())
		case object.BooleanType:
			if tg.Text() == "0" || tg.Text() == "false" {
				return NewBoolean(false)
			}
			if tg.Text() == "1" || tg.Text() == "true" {
				return NewBoolean(true)
			}
		}
	}

	return NewError("cannot convert %s with value %s to %s", tg.Type(), tg.Inspect(), ty)
}

// IsPrecede ..
func IsPrecede(n1, n2 object.Node, src *object.BaseNode) object.Item {
	for c := src.FirstChild(); c != nil; c = c.NextSibling() {
		c := c.(*object.BaseNode)

		if n1.Type() == object.AttributeNodeType && n2.Type() == object.AttributeNodeType {
			n1 := n1.(*object.AttrNode)
			n2 := n2.(*object.AttrNode)

			for _, a := range src.Attr() {
				a := a.(*object.AttrNode)
				if n1.Tree() == a.Tree() && n1.Key() == a.Key() {
					return NewBoolean(true)
				}
				if n2.Tree() == a.Tree() && n2.Key() == a.Key() {
					return NewBoolean(false)
				}
			}
		} else if n1.Type() == object.AttributeNodeType {
			n1 := n1.(*object.AttrNode)

			for _, a := range src.Attr() {
				a := a.(*object.AttrNode)
				if n1.Tree() == a.Tree() && n1.Key() == a.Key() {
					if n1.Tree() != n2.Tree() {
						return NewBoolean(true)
					} else {
						return NewBoolean(false)
					}

				}
			}
		} else if n2.Type() == object.AttributeNodeType {
			n2 := n2.(*object.AttrNode)

			for _, a := range src.Attr() {
				a := a.(*object.AttrNode)
				if n2.Tree() == a.Tree() && n2.Key() == a.Key() {
					if n1.Tree() != n2.Tree() {
						return NewBoolean(false)
					} else {
						return NewBoolean(true)
					}

				}
			}
		} else {
			if n2.Tree() == c.Tree() {
				return NewBoolean(false)
			}
			if n1.Tree() == c.Tree() {
				return NewBoolean(true)
			}
		}

		if c.FirstChild() != nil {
			result := IsPrecede(n1, n2, c)
			if result != object.NIL {
				return result
			}
		}
	}
	return object.NIL
}

// IsEQ ..
func IsEQ(left, right object.Item) object.Item {
	if leftVal, ok := left.(*object.Integer); ok {
		switch rightVal := right.(type) {
		case *object.Integer:
			return NewBoolean(leftVal.Value() == rightVal.Value())
		case *object.Decimal:
			return NewBoolean(float64(leftVal.Value()) == rightVal.Value())
		case *object.Double:
			return NewBoolean(float64(leftVal.Value()) == rightVal.Value())
		case *object.BaseNode:
			if i, err := strconv.ParseInt(rightVal.Text(), 0, 64); err == nil {
				return NewBoolean(leftVal.Value() == int(i))
			}
		case *object.AttrNode:
			if i, err := strconv.ParseInt(rightVal.Text(), 0, 64); err == nil {
				return NewBoolean(leftVal.Value() == int(i))
			}
		}
	} else if leftVal, ok := left.(*object.Decimal); ok {
		switch rightVal := right.(type) {
		case *object.Integer:
			return NewBoolean(leftVal.Value() == float64(rightVal.Value()))
		case *object.Decimal:
			return NewBoolean(leftVal.Value() == rightVal.Value())
		case *object.Double:
			return NewBoolean(leftVal.Value() == rightVal.Value())
		case *object.BaseNode:
			if i, err := strconv.ParseFloat(rightVal.Text(), 64); err == nil {
				return NewBoolean(leftVal.Value() == i)
			}
		case *object.AttrNode:
			if i, err := strconv.ParseFloat(rightVal.Text(), 64); err == nil {
				return NewBoolean(leftVal.Value() == i)
			}
		}
	} else if leftVal, ok := left.(*object.Double); ok {
		switch rightVal := right.(type) {
		case *object.Integer:
			return NewBoolean(leftVal.Value() == float64(rightVal.Value()))
		case *object.Decimal:
			return NewBoolean(leftVal.Value() == rightVal.Value())
		case *object.Double:
			return NewBoolean(leftVal.Value() == rightVal.Value())
		case *object.BaseNode:
			if i, err := strconv.ParseFloat(rightVal.Text(), 64); err == nil {
				return NewBoolean(leftVal.Value() == i)
			}
		case *object.AttrNode:
			if i, err := strconv.ParseFloat(rightVal.Text(), 64); err == nil {
				return NewBoolean(leftVal.Value() == i)
			}
		}
	} else if leftVal, ok := left.(*object.String); ok {
		switch rightVal := right.(type) {
		case *object.Integer:
			if i, err := strconv.ParseInt(leftVal.Value(), 0, 64); err == nil {
				return NewBoolean(int(i) == rightVal.Value())
			}
		case *object.Decimal:
			if i, err := strconv.ParseFloat(leftVal.Value(), 64); err == nil {
				return NewBoolean(i == rightVal.Value())
			}
		case *object.Double:
			if i, err := strconv.ParseFloat(leftVal.Value(), 64); err == nil {
				return NewBoolean(i == rightVal.Value())
			}
		case *object.String:
			return NewBoolean(leftVal.Value() == rightVal.Value())
		case *object.BaseNode:
			return NewBoolean(leftVal.Value() == rightVal.Text())
		case *object.AttrNode:
			return NewBoolean(leftVal.Value() == rightVal.Text())
		}
	} else if leftVal, ok := left.(*object.BaseNode); ok {
		switch rightVal := right.(type) {
		case *object.Integer:
			if i, err := strconv.ParseInt(leftVal.Text(), 0, 64); err == nil {
				return NewBoolean(int(i) == rightVal.Value())
			}
		case *object.Decimal:
			if i, err := strconv.ParseFloat(leftVal.Text(), 64); err == nil {
				return NewBoolean(i == rightVal.Value())
			}
		case *object.Double:
			if i, err := strconv.ParseFloat(leftVal.Text(), 64); err == nil {
				return NewBoolean(i == rightVal.Value())
			}
		case *object.BaseNode:
			return NewBoolean(leftVal.Text() == rightVal.Text())
		case *object.AttrNode:
			return NewBoolean(leftVal.Text() == rightVal.Text())
		case *object.String:
			return NewBoolean(leftVal.Text() == rightVal.Value())
		}
	} else if leftVal, ok := left.(*object.AttrNode); ok {
		switch rightVal := right.(type) {
		case *object.Integer:
			if i, err := strconv.ParseInt(leftVal.Text(), 0, 64); err == nil {
				return NewBoolean(int(i) == rightVal.Value())
			}
		case *object.Decimal:
			if i, err := strconv.ParseFloat(leftVal.Text(), 64); err == nil {
				return NewBoolean(i == rightVal.Value())
			}
		case *object.Double:
			if i, err := strconv.ParseFloat(leftVal.Text(), 64); err == nil {
				return NewBoolean(i == rightVal.Value())
			}
		case *object.AttrNode:
			return NewBoolean(leftVal.Text() == rightVal.Text())
		case *object.BaseNode:
			return NewBoolean(leftVal.Text() == rightVal.Text())
		case *object.String:
			return NewBoolean(leftVal.Text() == rightVal.Value())
		}
	}

	return NewError("cannot compare: %s, %s", left.Inspect(), right.Inspect())
}

// IsNE ..
func IsNE(left, right object.Item) object.Item {
	if leftVal, ok := left.(*object.Integer); ok {
		switch rightVal := right.(type) {
		case *object.Integer:
			return NewBoolean(leftVal.Value() != rightVal.Value())
		case *object.Decimal:
			return NewBoolean(float64(leftVal.Value()) != rightVal.Value())
		case *object.Double:
			return NewBoolean(float64(leftVal.Value()) != rightVal.Value())
		case *object.BaseNode:
			if i, err := strconv.ParseInt(rightVal.Text(), 0, 64); err == nil {
				return NewBoolean(leftVal.Value() != int(i))
			}
		case *object.AttrNode:
			if i, err := strconv.ParseInt(rightVal.Text(), 0, 64); err == nil {
				return NewBoolean(leftVal.Value() != int(i))
			}
		}
	} else if leftVal, ok := left.(*object.Decimal); ok {
		switch rightVal := right.(type) {
		case *object.Integer:
			return NewBoolean(leftVal.Value() != float64(rightVal.Value()))
		case *object.Decimal:
			return NewBoolean(leftVal.Value() != rightVal.Value())
		case *object.Double:
			return NewBoolean(leftVal.Value() != rightVal.Value())
		case *object.BaseNode:
			if i, err := strconv.ParseFloat(rightVal.Text(), 64); err == nil {
				return NewBoolean(leftVal.Value() != i)
			}
		case *object.AttrNode:
			if i, err := strconv.ParseFloat(rightVal.Text(), 64); err == nil {
				return NewBoolean(leftVal.Value() != i)
			}
		}
	} else if leftVal, ok := left.(*object.Double); ok {
		switch rightVal := right.(type) {
		case *object.Integer:
			return NewBoolean(leftVal.Value() != float64(rightVal.Value()))
		case *object.Decimal:
			return NewBoolean(leftVal.Value() != rightVal.Value())
		case *object.Double:
			return NewBoolean(leftVal.Value() != rightVal.Value())
		case *object.BaseNode:
			if i, err := strconv.ParseFloat(rightVal.Text(), 64); err == nil {
				return NewBoolean(leftVal.Value() != i)
			}
		case *object.AttrNode:
			if i, err := strconv.ParseFloat(rightVal.Text(), 64); err == nil {
				return NewBoolean(leftVal.Value() != i)
			}
		}
	} else if leftVal, ok := left.(*object.String); ok {
		switch rightVal := right.(type) {
		case *object.Integer:
			if i, err := strconv.ParseInt(leftVal.Value(), 0, 64); err == nil {
				return NewBoolean(int(i) != rightVal.Value())
			}
		case *object.Decimal:
			if i, err := strconv.ParseFloat(leftVal.Value(), 64); err == nil {
				return NewBoolean(i != rightVal.Value())
			}
		case *object.Double:
			if i, err := strconv.ParseFloat(leftVal.Value(), 64); err == nil {
				return NewBoolean(i != rightVal.Value())
			}
		case *object.String:
			return NewBoolean(leftVal.Value() != rightVal.Value())
		case *object.BaseNode:
			return NewBoolean(leftVal.Value() != rightVal.Text())
		case *object.AttrNode:
			return NewBoolean(leftVal.Value() != rightVal.Text())
		}
	} else if leftVal, ok := left.(*object.BaseNode); ok {
		switch rightVal := right.(type) {
		case *object.Integer:
			if i, err := strconv.ParseInt(leftVal.Text(), 0, 64); err == nil {
				return NewBoolean(int(i) != rightVal.Value())
			}
		case *object.Decimal:
			if i, err := strconv.ParseFloat(leftVal.Text(), 64); err == nil {
				return NewBoolean(i != rightVal.Value())
			}
		case *object.Double:
			if i, err := strconv.ParseFloat(leftVal.Text(), 64); err == nil {
				return NewBoolean(i != rightVal.Value())
			}
		case *object.BaseNode:
			return NewBoolean(leftVal.Text() != rightVal.Text())
		case *object.AttrNode:
			return NewBoolean(leftVal.Text() != rightVal.Text())
		case *object.String:
			return NewBoolean(leftVal.Text() != rightVal.Value())
		}
	} else if leftVal, ok := left.(*object.AttrNode); ok {
		switch rightVal := right.(type) {
		case *object.Integer:
			if i, err := strconv.ParseInt(leftVal.Text(), 0, 64); err == nil {
				return NewBoolean(int(i) != rightVal.Value())
			}
		case *object.Decimal:
			if i, err := strconv.ParseFloat(leftVal.Text(), 64); err == nil {
				return NewBoolean(i != rightVal.Value())
			}
		case *object.Double:
			if i, err := strconv.ParseFloat(leftVal.Text(), 64); err == nil {
				return NewBoolean(i != rightVal.Value())
			}
		case *object.AttrNode:
			return NewBoolean(leftVal.Text() != rightVal.Text())
		case *object.BaseNode:
			return NewBoolean(leftVal.Text() != rightVal.Text())
		case *object.String:
			return NewBoolean(leftVal.Text() != rightVal.Value())
		}
	}

	return NewError("cannot compare: %s, %s", left.Inspect(), right.Inspect())
}

// IsLT ..
func IsLT(left, right object.Item) object.Item {
	if leftVal, ok := left.(*object.Integer); ok {
		switch rightVal := right.(type) {
		case *object.Integer:
			return NewBoolean(leftVal.Value() < rightVal.Value())
		case *object.Decimal:
			return NewBoolean(float64(leftVal.Value()) < rightVal.Value())
		case *object.Double:
			return NewBoolean(float64(leftVal.Value()) < rightVal.Value())
		case *object.BaseNode:
			if i, err := strconv.ParseInt(rightVal.Text(), 0, 64); err == nil {
				return NewBoolean(leftVal.Value() < int(i))
			}
		case *object.AttrNode:
			if i, err := strconv.ParseInt(rightVal.Text(), 0, 64); err == nil {
				return NewBoolean(leftVal.Value() < int(i))
			}
		}
	} else if leftVal, ok := left.(*object.Decimal); ok {
		switch rightVal := right.(type) {
		case *object.Integer:
			return NewBoolean(leftVal.Value() < float64(rightVal.Value()))
		case *object.Decimal:
			return NewBoolean(leftVal.Value() < rightVal.Value())
		case *object.Double:
			return NewBoolean(leftVal.Value() < rightVal.Value())
		case *object.BaseNode:
			if i, err := strconv.ParseFloat(rightVal.Text(), 64); err == nil {
				return NewBoolean(leftVal.Value() < i)
			}
		case *object.AttrNode:
			if i, err := strconv.ParseFloat(rightVal.Text(), 64); err == nil {
				return NewBoolean(leftVal.Value() < i)
			}
		}
	} else if leftVal, ok := left.(*object.Double); ok {
		switch rightVal := right.(type) {
		case *object.Integer:
			return NewBoolean(leftVal.Value() < float64(rightVal.Value()))
		case *object.Decimal:
			return NewBoolean(leftVal.Value() < rightVal.Value())
		case *object.Double:
			return NewBoolean(leftVal.Value() < rightVal.Value())
		case *object.BaseNode:
			if i, err := strconv.ParseFloat(rightVal.Text(), 64); err == nil {
				return NewBoolean(leftVal.Value() < i)
			}
		case *object.AttrNode:
			if i, err := strconv.ParseFloat(rightVal.Text(), 64); err == nil {
				return NewBoolean(leftVal.Value() < i)
			}
		}
	} else if leftVal, ok := left.(*object.String); ok {
		switch rightVal := right.(type) {
		case *object.Integer:
			if i, err := strconv.ParseInt(leftVal.Value(), 0, 64); err == nil {
				return NewBoolean(int(i) < rightVal.Value())
			}
		case *object.Decimal:
			if i, err := strconv.ParseFloat(leftVal.Value(), 64); err == nil {
				return NewBoolean(i < rightVal.Value())
			}
		case *object.Double:
			if i, err := strconv.ParseFloat(leftVal.Value(), 64); err == nil {
				return NewBoolean(i < rightVal.Value())
			}
		case *object.String:
			return NewBoolean(leftVal.Value() < rightVal.Value())
		case *object.BaseNode:
			return NewBoolean(leftVal.Value() < rightVal.Text())
		case *object.AttrNode:
			return NewBoolean(leftVal.Value() < rightVal.Text())
		}
	} else if leftVal, ok := left.(*object.BaseNode); ok {
		switch rightVal := right.(type) {
		case *object.Integer:
			if i, err := strconv.ParseInt(leftVal.Text(), 0, 64); err == nil {
				return NewBoolean(int(i) < rightVal.Value())
			}
		case *object.Decimal:
			if i, err := strconv.ParseFloat(leftVal.Text(), 64); err == nil {
				return NewBoolean(i < rightVal.Value())
			}
		case *object.Double:
			if i, err := strconv.ParseFloat(leftVal.Text(), 64); err == nil {
				return NewBoolean(i < rightVal.Value())
			}
		case *object.BaseNode:
			return NewBoolean(leftVal.Text() < rightVal.Text())
		case *object.AttrNode:
			return NewBoolean(leftVal.Text() < rightVal.Text())
		case *object.String:
			return NewBoolean(leftVal.Text() < rightVal.Value())
		}
	} else if leftVal, ok := left.(*object.AttrNode); ok {
		switch rightVal := right.(type) {
		case *object.Integer:
			if i, err := strconv.ParseInt(leftVal.Text(), 0, 64); err == nil {
				return NewBoolean(int(i) < rightVal.Value())
			}
		case *object.Decimal:
			if i, err := strconv.ParseFloat(leftVal.Text(), 64); err == nil {
				return NewBoolean(i < rightVal.Value())
			}
		case *object.Double:
			if i, err := strconv.ParseFloat(leftVal.Text(), 64); err == nil {
				return NewBoolean(i < rightVal.Value())
			}
		case *object.AttrNode:
			return NewBoolean(leftVal.Text() < rightVal.Text())
		case *object.BaseNode:
			return NewBoolean(leftVal.Text() < rightVal.Text())
		case *object.String:
			return NewBoolean(leftVal.Text() < rightVal.Value())
		}
	}

	return NewError("cannot compare: %s, %s", left.Inspect(), right.Inspect())
}

// IsLE ..
func IsLE(left, right object.Item) object.Item {
	if leftVal, ok := left.(*object.Integer); ok {
		switch rightVal := right.(type) {
		case *object.Integer:
			return NewBoolean(leftVal.Value() <= rightVal.Value())
		case *object.Decimal:
			return NewBoolean(float64(leftVal.Value()) <= rightVal.Value())
		case *object.Double:
			return NewBoolean(float64(leftVal.Value()) <= rightVal.Value())
		case *object.BaseNode:
			if i, err := strconv.ParseInt(rightVal.Text(), 0, 64); err == nil {
				return NewBoolean(leftVal.Value() <= int(i))
			}
		case *object.AttrNode:
			if i, err := strconv.ParseInt(rightVal.Text(), 0, 64); err == nil {
				return NewBoolean(leftVal.Value() <= int(i))
			}
		}
	} else if leftVal, ok := left.(*object.Decimal); ok {
		switch rightVal := right.(type) {
		case *object.Integer:
			return NewBoolean(leftVal.Value() <= float64(rightVal.Value()))
		case *object.Decimal:
			return NewBoolean(leftVal.Value() <= rightVal.Value())
		case *object.Double:
			return NewBoolean(leftVal.Value() <= rightVal.Value())
		case *object.BaseNode:
			if i, err := strconv.ParseFloat(rightVal.Text(), 64); err == nil {
				return NewBoolean(leftVal.Value() <= i)
			}
		case *object.AttrNode:
			if i, err := strconv.ParseFloat(rightVal.Text(), 64); err == nil {
				return NewBoolean(leftVal.Value() <= i)
			}
		}
	} else if leftVal, ok := left.(*object.Double); ok {
		switch rightVal := right.(type) {
		case *object.Integer:
			return NewBoolean(leftVal.Value() <= float64(rightVal.Value()))
		case *object.Decimal:
			return NewBoolean(leftVal.Value() <= rightVal.Value())
		case *object.Double:
			return NewBoolean(leftVal.Value() <= rightVal.Value())
		case *object.BaseNode:
			if i, err := strconv.ParseFloat(rightVal.Text(), 64); err == nil {
				return NewBoolean(leftVal.Value() <= i)
			}
		case *object.AttrNode:
			if i, err := strconv.ParseFloat(rightVal.Text(), 64); err == nil {
				return NewBoolean(leftVal.Value() <= i)
			}
		}
	} else if leftVal, ok := left.(*object.String); ok {
		switch rightVal := right.(type) {
		case *object.Integer:
			if i, err := strconv.ParseInt(leftVal.Value(), 0, 64); err == nil {
				return NewBoolean(int(i) <= rightVal.Value())
			}
		case *object.Decimal:
			if i, err := strconv.ParseFloat(leftVal.Value(), 64); err == nil {
				return NewBoolean(i <= rightVal.Value())
			}
		case *object.Double:
			if i, err := strconv.ParseFloat(leftVal.Value(), 64); err == nil {
				return NewBoolean(i <= rightVal.Value())
			}
		case *object.String:
			return NewBoolean(leftVal.Value() <= rightVal.Value())
		case *object.BaseNode:
			return NewBoolean(leftVal.Value() <= rightVal.Text())
		case *object.AttrNode:
			return NewBoolean(leftVal.Value() <= rightVal.Text())
		}
	} else if leftVal, ok := left.(*object.BaseNode); ok {
		switch rightVal := right.(type) {
		case *object.Integer:
			if i, err := strconv.ParseInt(leftVal.Text(), 0, 64); err == nil {
				return NewBoolean(int(i) <= rightVal.Value())
			}
		case *object.Decimal:
			if i, err := strconv.ParseFloat(leftVal.Text(), 64); err == nil {
				return NewBoolean(i <= rightVal.Value())
			}
		case *object.Double:
			if i, err := strconv.ParseFloat(leftVal.Text(), 64); err == nil {
				return NewBoolean(i <= rightVal.Value())
			}
		case *object.BaseNode:
			return NewBoolean(leftVal.Text() <= rightVal.Text())
		case *object.AttrNode:
			return NewBoolean(leftVal.Text() <= rightVal.Text())
		case *object.String:
			return NewBoolean(leftVal.Text() <= rightVal.Value())
		}
	} else if leftVal, ok := left.(*object.AttrNode); ok {
		switch rightVal := right.(type) {
		case *object.Integer:
			if i, err := strconv.ParseInt(leftVal.Text(), 0, 64); err == nil {
				return NewBoolean(int(i) <= rightVal.Value())
			}
		case *object.Decimal:
			if i, err := strconv.ParseFloat(leftVal.Text(), 64); err == nil {
				return NewBoolean(i <= rightVal.Value())
			}
		case *object.Double:
			if i, err := strconv.ParseFloat(leftVal.Text(), 64); err == nil {
				return NewBoolean(i <= rightVal.Value())
			}
		case *object.AttrNode:
			return NewBoolean(leftVal.Text() <= rightVal.Text())
		case *object.BaseNode:
			return NewBoolean(leftVal.Text() <= rightVal.Text())
		case *object.String:
			return NewBoolean(leftVal.Text() <= rightVal.Value())
		}
	}

	return NewError("cannot compare: %s, %s", left.Inspect(), right.Inspect())
}

// IsGT ..
func IsGT(left, right object.Item) object.Item {
	if leftVal, ok := left.(*object.Integer); ok {
		switch rightVal := right.(type) {
		case *object.Integer:
			return NewBoolean(leftVal.Value() > rightVal.Value())
		case *object.Decimal:
			return NewBoolean(float64(leftVal.Value()) > rightVal.Value())
		case *object.Double:
			return NewBoolean(float64(leftVal.Value()) > rightVal.Value())
		case *object.BaseNode:
			if i, err := strconv.ParseInt(rightVal.Text(), 0, 64); err == nil {
				return NewBoolean(leftVal.Value() > int(i))
			}
		case *object.AttrNode:
			if i, err := strconv.ParseInt(rightVal.Text(), 0, 64); err == nil {
				return NewBoolean(leftVal.Value() > int(i))
			}
		}
	} else if leftVal, ok := left.(*object.Decimal); ok {
		switch rightVal := right.(type) {
		case *object.Integer:
			return NewBoolean(leftVal.Value() > float64(rightVal.Value()))
		case *object.Decimal:
			return NewBoolean(leftVal.Value() > rightVal.Value())
		case *object.Double:
			return NewBoolean(leftVal.Value() > rightVal.Value())
		case *object.BaseNode:
			if i, err := strconv.ParseFloat(rightVal.Text(), 64); err == nil {
				return NewBoolean(leftVal.Value() > i)
			}
		case *object.AttrNode:
			if i, err := strconv.ParseFloat(rightVal.Text(), 64); err == nil {
				return NewBoolean(leftVal.Value() > i)
			}
		}
	} else if leftVal, ok := left.(*object.Double); ok {
		switch rightVal := right.(type) {
		case *object.Integer:
			return NewBoolean(leftVal.Value() > float64(rightVal.Value()))
		case *object.Decimal:
			return NewBoolean(leftVal.Value() > rightVal.Value())
		case *object.Double:
			return NewBoolean(leftVal.Value() > rightVal.Value())
		case *object.BaseNode:
			if i, err := strconv.ParseFloat(rightVal.Text(), 64); err == nil {
				return NewBoolean(leftVal.Value() > i)
			}
		case *object.AttrNode:
			if i, err := strconv.ParseFloat(rightVal.Text(), 64); err == nil {
				return NewBoolean(leftVal.Value() > i)
			}
		}
	} else if leftVal, ok := left.(*object.String); ok {
		switch rightVal := right.(type) {
		case *object.Integer:
			if i, err := strconv.ParseInt(leftVal.Value(), 0, 64); err == nil {
				return NewBoolean(int(i) > rightVal.Value())
			}
		case *object.Decimal:
			if i, err := strconv.ParseFloat(leftVal.Value(), 64); err == nil {
				return NewBoolean(i > rightVal.Value())
			}
		case *object.Double:
			if i, err := strconv.ParseFloat(leftVal.Value(), 64); err == nil {
				return NewBoolean(i > rightVal.Value())
			}
		case *object.String:
			return NewBoolean(leftVal.Value() > rightVal.Value())
		case *object.BaseNode:
			return NewBoolean(leftVal.Value() > rightVal.Text())
		case *object.AttrNode:
			return NewBoolean(leftVal.Value() > rightVal.Text())
		}
	} else if leftVal, ok := left.(*object.BaseNode); ok {
		switch rightVal := right.(type) {
		case *object.Integer:
			if i, err := strconv.ParseInt(leftVal.Text(), 0, 64); err == nil {
				return NewBoolean(int(i) > rightVal.Value())
			}
		case *object.Decimal:
			if i, err := strconv.ParseFloat(leftVal.Text(), 64); err == nil {
				return NewBoolean(i > rightVal.Value())
			}
		case *object.Double:
			if i, err := strconv.ParseFloat(leftVal.Text(), 64); err == nil {
				return NewBoolean(i > rightVal.Value())
			}
		case *object.BaseNode:
			return NewBoolean(leftVal.Text() > rightVal.Text())
		case *object.AttrNode:
			return NewBoolean(leftVal.Text() > rightVal.Text())
		case *object.String:
			return NewBoolean(leftVal.Text() > rightVal.Value())
		}
	} else if leftVal, ok := left.(*object.AttrNode); ok {
		switch rightVal := right.(type) {
		case *object.Integer:
			if i, err := strconv.ParseInt(leftVal.Text(), 0, 64); err == nil {
				return NewBoolean(int(i) > rightVal.Value())
			}
		case *object.Decimal:
			if i, err := strconv.ParseFloat(leftVal.Text(), 64); err == nil {
				return NewBoolean(i > rightVal.Value())
			}
		case *object.Double:
			if i, err := strconv.ParseFloat(leftVal.Text(), 64); err == nil {
				return NewBoolean(i > rightVal.Value())
			}
		case *object.AttrNode:
			return NewBoolean(leftVal.Text() > rightVal.Text())
		case *object.BaseNode:
			return NewBoolean(leftVal.Text() > rightVal.Text())
		case *object.String:
			return NewBoolean(leftVal.Text() > rightVal.Value())
		}
	}

	return NewError("cannot compare: %s, %s", left.Inspect(), right.Inspect())
}

// IsGE ..
func IsGE(left, right object.Item) object.Item {
	if leftVal, ok := left.(*object.Integer); ok {
		switch rightVal := right.(type) {
		case *object.Integer:
			return NewBoolean(leftVal.Value() >= rightVal.Value())
		case *object.Decimal:
			return NewBoolean(float64(leftVal.Value()) >= rightVal.Value())
		case *object.Double:
			return NewBoolean(float64(leftVal.Value()) >= rightVal.Value())
		case *object.BaseNode:
			if i, err := strconv.ParseInt(rightVal.Text(), 0, 64); err == nil {
				return NewBoolean(leftVal.Value() >= int(i))
			}
		case *object.AttrNode:
			if i, err := strconv.ParseInt(rightVal.Text(), 0, 64); err == nil {
				return NewBoolean(leftVal.Value() >= int(i))
			}
		}
	} else if leftVal, ok := left.(*object.Decimal); ok {
		switch rightVal := right.(type) {
		case *object.Integer:
			return NewBoolean(leftVal.Value() >= float64(rightVal.Value()))
		case *object.Decimal:
			return NewBoolean(leftVal.Value() >= rightVal.Value())
		case *object.Double:
			return NewBoolean(leftVal.Value() >= rightVal.Value())
		case *object.BaseNode:
			if i, err := strconv.ParseFloat(rightVal.Text(), 64); err == nil {
				return NewBoolean(leftVal.Value() >= i)
			}
		case *object.AttrNode:
			if i, err := strconv.ParseFloat(rightVal.Text(), 64); err == nil {
				return NewBoolean(leftVal.Value() >= i)
			}
		}
	} else if leftVal, ok := left.(*object.Double); ok {
		switch rightVal := right.(type) {
		case *object.Integer:
			return NewBoolean(leftVal.Value() >= float64(rightVal.Value()))
		case *object.Decimal:
			return NewBoolean(leftVal.Value() >= rightVal.Value())
		case *object.Double:
			return NewBoolean(leftVal.Value() >= rightVal.Value())
		case *object.BaseNode:
			if i, err := strconv.ParseFloat(rightVal.Text(), 64); err == nil {
				return NewBoolean(leftVal.Value() >= i)
			}
		case *object.AttrNode:
			if i, err := strconv.ParseFloat(rightVal.Text(), 64); err == nil {
				return NewBoolean(leftVal.Value() >= i)
			}
		}
	} else if leftVal, ok := left.(*object.String); ok {
		switch rightVal := right.(type) {
		case *object.Integer:
			if i, err := strconv.ParseInt(leftVal.Value(), 0, 64); err == nil {
				return NewBoolean(int(i) >= rightVal.Value())
			}
		case *object.Decimal:
			if i, err := strconv.ParseFloat(leftVal.Value(), 64); err == nil {
				return NewBoolean(i >= rightVal.Value())
			}
		case *object.Double:
			if i, err := strconv.ParseFloat(leftVal.Value(), 64); err == nil {
				return NewBoolean(i >= rightVal.Value())
			}
		case *object.String:
			return NewBoolean(leftVal.Value() >= rightVal.Value())
		case *object.BaseNode:
			return NewBoolean(leftVal.Value() >= rightVal.Text())
		case *object.AttrNode:
			return NewBoolean(leftVal.Value() >= rightVal.Text())
		}
	} else if leftVal, ok := left.(*object.BaseNode); ok {
		switch rightVal := right.(type) {
		case *object.Integer:
			if i, err := strconv.ParseInt(leftVal.Text(), 0, 64); err == nil {
				return NewBoolean(int(i) >= rightVal.Value())
			}
		case *object.Decimal:
			if i, err := strconv.ParseFloat(leftVal.Text(), 64); err == nil {
				return NewBoolean(i >= rightVal.Value())
			}
		case *object.Double:
			if i, err := strconv.ParseFloat(leftVal.Text(), 64); err == nil {
				return NewBoolean(i >= rightVal.Value())
			}
		case *object.BaseNode:
			return NewBoolean(leftVal.Text() >= rightVal.Text())
		case *object.AttrNode:
			return NewBoolean(leftVal.Text() >= rightVal.Text())
		case *object.String:
			return NewBoolean(leftVal.Text() >= rightVal.Value())
		}
	} else if leftVal, ok := left.(*object.AttrNode); ok {
		switch rightVal := right.(type) {
		case *object.Integer:
			if i, err := strconv.ParseInt(leftVal.Text(), 0, 64); err == nil {
				return NewBoolean(int(i) >= rightVal.Value())
			}
		case *object.Decimal:
			if i, err := strconv.ParseFloat(leftVal.Text(), 64); err == nil {
				return NewBoolean(i >= rightVal.Value())
			}
		case *object.Double:
			if i, err := strconv.ParseFloat(leftVal.Text(), 64); err == nil {
				return NewBoolean(i >= rightVal.Value())
			}
		case *object.AttrNode:
			return NewBoolean(leftVal.Text() >= rightVal.Text())
		case *object.BaseNode:
			return NewBoolean(leftVal.Text() >= rightVal.Text())
		case *object.String:
			return NewBoolean(leftVal.Text() >= rightVal.Value())
		}
	}

	return NewError("cannot compare: %s, %s", left.Inspect(), right.Inspect())
}

// IsOccurMatch ..
func IsOccurMatch(item object.Item, t token.Token) bool {
	seq, ok := item.(*object.Sequence)
	if !ok {
		return false
	}

	switch t.Type {
	case token.QUESTION:
		if seq.Items == nil || len(seq.Items) == 1 {
			return true
		}
		return false
	case token.ASTERISK:
		return true
	case token.PLUS:
		if len(seq.Items) > 0 {
			return true
		}
		return false
	default:
		if len(seq.Items) == 1 {
			return true
		}
		return false
	}
}

// IsTypeMatch ..
func IsTypeMatch(item object.Item, st *ast.SequenceType) object.Item {
	switch st.TypeID {
	case 1:
		seq, ok := item.(*object.Sequence)
		if !ok {
			return NewBoolean(false)
		}
		if seq.Items != nil {
			return NewBoolean(false)
		}
		return NewBoolean(true)
	case 2:
		oi := st.OccurrenceIndicator
		it := st.NodeTest.(*ast.ItemType)

		switch it.TypeID {
		case 1:
			kt := it.NodeTest.(*ast.KindTest)
			switch kt.TypeID {
			case 1:
				if item.Type() == object.SequenceType {
					return NewBoolean(IsOccurMatch(item, oi.Token) && IsDocNodeSeq(item))
				}
				return NewBoolean(item.Type() == object.DocumentNodeType)
			case 2:
				if item.Type() == object.SequenceType {
					return NewBoolean(IsOccurMatch(item, oi.Token) && IsElemNodeSeq(item))
				}
				return NewBoolean(item.Type() == object.ElementNodeType)
			case 3:
				if item.Type() == object.SequenceType {
					return NewBoolean(IsOccurMatch(item, oi.Token) && IsAttrNodeSeq(item))
				}
				return NewBoolean(item.Type() == object.AttributeNodeType)
			case 7:
				if item.Type() == object.SequenceType {
					return NewBoolean(IsOccurMatch(item, oi.Token) && IsCommNodeSeq(item))
				}
				return NewBoolean(item.Type() == object.CommentNodeType)
			case 8:
				if item.Type() == object.SequenceType {
					return NewBoolean(IsOccurMatch(item, oi.Token) && IsTextNodeSeq(item))
				}
				return NewBoolean(item.Type() == object.TextNodeType)
			case 10:
				if item.Type() == object.SequenceType {
					return NewBoolean(IsOccurMatch(item, oi.Token) && IsNodeSeq(item))
				}
				return NewBoolean(IsNode(item))
			case 4:
				fallthrough
			case 5:
				fallthrough
			case 6:
				fallthrough
			case 9:
				return NewError("not supported kind test")
			}
		case 2:
			if item.Type() == object.SequenceType {
				return NewBoolean(IsOccurMatch(item, oi.Token) && IsItemSeq(item))
			}
			return NewBoolean(IsItem(item))
		case 3:
			if item.Type() == object.SequenceType {
				return NewBoolean(IsOccurMatch(item, oi.Token) && IsAnyFuncSeq(item))
			}
			return NewBoolean(IsAnyFunc(item))
		case 4:
			if item.Type() == object.SequenceType {
				return NewBoolean(IsOccurMatch(item, oi.Token) && IsMapSeq(item))
			}
			return NewBoolean(IsMap(item))
		case 5:
			if item.Type() == object.SequenceType {
				return NewBoolean(IsOccurMatch(item, oi.Token) && IsArraySeq(item))
			}
			return NewBoolean(IsArray(item))
		case 6:
			if item.Type() == object.SequenceType {
				return NewBoolean(IsOccurMatch(item, oi.Token) && IsAtomicSeq(item))
			}
			return NewBoolean(IsAnyAtomic(item))
		case 7:
			pit := it.NodeTest.(*ast.ParenthesizedItemType)
			st.NodeTest = pit.NodeTest
			return IsTypeMatch(item, st)
		}
	}

	return NewBoolean(false)
}

// IsKindMatch ..
func IsKindMatch(n object.Node, typeID byte) bool {
	switch typeID {
	case 1:
		if n.Type() == object.DocumentNodeType {
			return true
		}
	case 2:
		if n.Type() == object.ElementNodeType {
			return true
		}
	case 3:
		if n.Type() == object.AttributeNodeType {
			return true
		}
	case 7:
		if n.Type() == object.CommentNodeType {
			return true
		}
	case 8:
		if n.Type() == object.TextNodeType {
			return true
		}
	case 10:
		if n.Type() != object.AttributeNodeType {
			return true
		}
	}
	return false
}

// IsContainN ..
func IsContainN(src []object.Node, target object.Node) bool {
	for _, item := range src {
		if item, ok := item.(*object.BaseNode); ok {
			target, ok := target.(*object.BaseNode)
			if !ok {
				return false
			}

			if item.Tree() == target.Tree() {
				return true
			}
		}

		if item, ok := item.(*object.AttrNode); ok {
			target, ok := target.(*object.AttrNode)
			if !ok {
				return false
			}

			if item.Tree() == target.Tree() && item.Key() == target.Key() {
				return true
			}
		}
	}
	return false
}

// AppendNode ..
func AppendNode(src []object.Node, target object.Node) []object.Node {
	if !IsContainN(src, target) {
		src = append(src, target)
	}
	return src
}

// CopyFocus ..
func CopyFocus(ctx *object.Context) *object.Focus {
	return &object.Focus{CSize: ctx.CSize, CPos: ctx.CPos, CAxis: ctx.CAxis}
}

// ReplaceFocus ..
func ReplaceFocus(ctx *object.Context, focus *object.Focus) {
	ctx.CSize = focus.CSize
	ctx.CAxis = focus.CAxis
	ctx.CPos = focus.CPos
}

// UnwrapSeq ..
func UnwrapSeq(item object.Item) []object.Item {
	if seq, ok := item.(*object.Sequence); ok {
		var items []object.Item
		for _, it := range seq.Items {
			if it.Type() == object.SequenceType {
				items = append(items, UnwrapSeq(it)...)
			} else {
				items = append(items, it)
			}
		}
		return items
	}
	return []object.Item{item}
}
