package bif

import (
	"fmt"

	"github.com/zzossig/xpath/ast"
	"github.com/zzossig/xpath/object"
)

// Builtins defined in https://www.w3.org/TR/xpath-functions-31/
var Builtins = map[string]object.Func{
	"op:numeric-add":            numericAdd,
	"op:numeric-subtract":       numericSubtract,
	"op:numeric-multiply":       numericMultiply,
	"op:numeric-divide":         numericDivide,
	"op:numeric-integer-divide": numericIntegerDivide,
	"op:numeric-mod":            numericMod,
	"op:numeric-unary-plus":     numericUnaryPlus,
	"op:numeric-unary-minus":    numericUnaryMinus,
	"op:numeric-equal":          numericEqual,
	"op:numeric-less-than":      numericLessThan,
	"op:numeric-greater-than":   numericGreaterThan,

	"fn:doc":           doc,
	"fn:abs":           abs,
	"fn:concat":        concat,
	"fn:for-each-pair": forEachPair,
	"fn:upper-case":    upperCase,
	"fn:lower-case":    lowerCase,
	"fn:boolean":       boolean,
}

// BuiltinPNames ..
func BuiltinPNames(name string, num int) []string {
	var pnames []string

	switch name {
	default:
		pnames = append(pnames, "arg")
	case "fn:concat":
		for i := 1; i <= num; i++ {
			pnames = append(pnames, fmt.Sprintf("arg%d", i))
		}
	case "fn:for-each-pair":
		pnames = append(pnames, []string{"seq1", "seq2", "action"}...)
	case "fn:doc":
		pnames = append(pnames, "uri")
	}

	return pnames
}

// BuiltinPTypes ..
func BuiltinPTypes(name string, num int) []object.Type {
	var ptypes []object.Type

	switch name {
	case "op:numeric-add":
		fallthrough
	case "op:numeric-subtract":
		fallthrough
	case "op:numeric-multiply":
		fallthrough
	case "op:numeric-divide":
		fallthrough
	case "op:numeric-integer-divide":
		fallthrough
	case "op:numeric-mod":
		fallthrough
	case "op:numeric-equal":
		fallthrough
	case "op:numeric-less-than":
		fallthrough
	case "op:numeric-greater-than":
		ptypes = append(ptypes, []object.Type{object.NumericType, object.NumericType}...)
	case "op:numeric-unary-plus":
		fallthrough
	case "op:numeric-unary-minus":
		ptypes = append(ptypes, object.NumericType)
	case "fn:doc":
		ptypes = append(ptypes, object.StringType)
	case "fn:abs":
		ptypes = append(ptypes, object.NumericType)
	case "fn:lower-case":
		ptypes = append(ptypes, object.StringType)
	case "fn:upper-case":
		ptypes = append(ptypes, object.StringType)
	case "fn:boolean":
		ptypes = append(ptypes, object.ItemType)
	case "fn:for-each-pair":
		ptypes = append(ptypes, []object.Type{object.ItemType, object.ItemType, object.FuncType}...)
	}

	return ptypes
}

// CheckBuiltinPTypes ..
func CheckBuiltinPTypes(fname string, args []object.Item) object.Item {
	if fname == "fn:concat" {
		for _, arg := range args {
			if !IsAnyAtomic(arg) {
				return NewError("wrong type of argument in concat function: %s", arg.Type())
			}
		}
	} else {
		types := BuiltinPTypes(fname, len(args))

		for i, t := range types {
			isMatch := true

			switch t {
			case object.AnyAtomicType:
				if !IsAnyAtomic(args[i]) {
					isMatch = false
				}
			case object.NumericType:
				if !IsNumeric(args[i]) {
					isMatch = false
				}
			case object.FuncType:
				if !IsFunction(args[i]) {
					isMatch = false
				}
			case object.MapType:
				if !IsMap(args[i]) {
					isMatch = false
				}
			case object.ArrayType:
				if IsArray(args[i]) {
					isMatch = false
				}
			case object.NodeType:
				if !IsNode(args[i]) {
					isMatch = false
				}
			case object.StringType:
				if IsString(args[i]) {
					isMatch = false
				}
			case object.BooleanType:
				if IsBoolean(args[i]) {
					isMatch = false
				}
			case object.DoubleType, object.DecimalType:
				if args[i].Type() != object.DoubleType && args[i].Type() != object.DecimalType {
					isMatch = false
				}
			case object.IntegerType:
				if args[i].Type() != object.IntegerType {
					isMatch = false
				}
			}

			if !isMatch {
				return NewError("wrong type of argument in %s function: %s", fname, args[i].Type())
			}
		}
	}

	return object.NIL
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

// IsError ..
func IsError(item object.Item) bool {
	if item == nil {
		return false
	}
	return item.Type() == object.ErrorType
}

// IsSeq ..
func IsSeq(item object.Item) bool {
	if item == nil {
		return false
	}
	return item.Type() == object.SequenceType
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

// IsFunction ..
func IsFunction(item object.Item) bool {
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
		default:
			return NewError("cannot convert %s to number", right.Type())
		}
	} else if leftVal, ok := left.(*object.Decimal); ok {
		switch rightVal := right.(type) {
		case *object.Integer:
			return NewBoolean(leftVal.Value() == float64(rightVal.Value()))
		case *object.Decimal:
			return NewBoolean(leftVal.Value() == rightVal.Value())
		case *object.Double:
			return NewBoolean(leftVal.Value() == rightVal.Value())
		default:
			return NewError("cannot convert %s to number", right.Type())
		}
	} else if leftVal, ok := left.(*object.Double); ok {
		switch rightVal := right.(type) {
		case *object.Integer:
			return NewBoolean(leftVal.Value() == float64(rightVal.Value()))
		case *object.Decimal:
			return NewBoolean(leftVal.Value() == rightVal.Value())
		case *object.Double:
			return NewBoolean(leftVal.Value() == rightVal.Value())
		default:
			return NewError("cannot convert %s to number", right.Type())
		}
	} else if leftVal, ok := left.(*object.String); ok {
		switch rightVal := right.(type) {
		case *object.String:
			return NewBoolean(leftVal.Value() == rightVal.Value())
		default:
			return NewError("cannot compare %s and %s", left.Type(), right.Type())
		}
	}
	return NewError("cannot compare %s and %s", left.Type(), right.Type())
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
		default:
			return NewError("cannot convert %s to number", right.Type())
		}
	} else if leftVal, ok := left.(*object.Decimal); ok {
		switch rightVal := right.(type) {
		case *object.Integer:
			return NewBoolean(leftVal.Value() != float64(rightVal.Value()))
		case *object.Decimal:
			return NewBoolean(leftVal.Value() != rightVal.Value())
		case *object.Double:
			return NewBoolean(leftVal.Value() != rightVal.Value())
		default:
			return NewError("cannot convert %s to number", right.Type())
		}
	} else if leftVal, ok := left.(*object.Double); ok {
		switch rightVal := right.(type) {
		case *object.Integer:
			return NewBoolean(leftVal.Value() != float64(rightVal.Value()))
		case *object.Decimal:
			return NewBoolean(leftVal.Value() != rightVal.Value())
		case *object.Double:
			return NewBoolean(leftVal.Value() != rightVal.Value())
		default:
			return NewError("cannot convert %s to number", right.Type())
		}
	} else if leftVal, ok := left.(*object.String); ok {
		switch rightVal := right.(type) {
		case *object.String:
			return NewBoolean(leftVal.Value() != rightVal.Value())
		default:
			return NewError("cannot compare %s and %s", left.Type(), right.Type())
		}
	}
	return NewError("cannot compare %s and %s", left.Type(), right.Type())
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
		default:
			return NewError("cannot convert %s to number", right.Type())
		}
	} else if leftVal, ok := left.(*object.Decimal); ok {
		switch rightVal := right.(type) {
		case *object.Integer:
			return NewBoolean(leftVal.Value() < float64(rightVal.Value()))
		case *object.Decimal:
			return NewBoolean(leftVal.Value() < rightVal.Value())
		case *object.Double:
			return NewBoolean(leftVal.Value() < rightVal.Value())
		default:
			return NewError("cannot convert %s to number", right.Type())
		}
	} else if leftVal, ok := left.(*object.Double); ok {
		switch rightVal := right.(type) {
		case *object.Integer:
			return NewBoolean(leftVal.Value() < float64(rightVal.Value()))
		case *object.Decimal:
			return NewBoolean(leftVal.Value() < rightVal.Value())
		case *object.Double:
			return NewBoolean(leftVal.Value() < rightVal.Value())
		default:
			return NewError("cannot convert %s to number", right.Type())
		}
	} else if leftVal, ok := left.(*object.String); ok {
		switch rightVal := right.(type) {
		case *object.String:
			return NewBoolean(leftVal.Value() < rightVal.Value())
		default:
			return NewError("cannot compare %s and %s", left.Type(), right.Type())
		}
	}
	return NewError("cannot compare %s and %s", left.Type(), right.Type())
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
		default:
			return NewError("cannot convert %s to number", right.Type())
		}
	} else if leftVal, ok := left.(*object.Decimal); ok {
		switch rightVal := right.(type) {
		case *object.Integer:
			return NewBoolean(leftVal.Value() <= float64(rightVal.Value()))
		case *object.Decimal:
			return NewBoolean(leftVal.Value() <= rightVal.Value())
		case *object.Double:
			return NewBoolean(leftVal.Value() <= rightVal.Value())
		default:
			return NewError("cannot convert %s to number", right.Type())
		}
	} else if leftVal, ok := left.(*object.Double); ok {
		switch rightVal := right.(type) {
		case *object.Integer:
			return NewBoolean(leftVal.Value() <= float64(rightVal.Value()))
		case *object.Decimal:
			return NewBoolean(leftVal.Value() <= rightVal.Value())
		case *object.Double:
			return NewBoolean(leftVal.Value() <= rightVal.Value())
		default:
			return NewError("cannot convert %s to number", right.Type())
		}
	} else if leftVal, ok := left.(*object.String); ok {
		switch rightVal := right.(type) {
		case *object.String:
			return NewBoolean(leftVal.Value() <= rightVal.Value())
		default:
			return NewError("cannot compare %s and %s", left.Type(), right.Type())
		}
	}
	return NewError("cannot compare %s and %s", left.Type(), right.Type())
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
		default:
			return NewError("cannot convert %s to number", right.Type())
		}
	} else if leftVal, ok := left.(*object.Decimal); ok {
		switch rightVal := right.(type) {
		case *object.Integer:
			return NewBoolean(leftVal.Value() > float64(rightVal.Value()))
		case *object.Decimal:
			return NewBoolean(leftVal.Value() > rightVal.Value())
		case *object.Double:
			return NewBoolean(leftVal.Value() > rightVal.Value())
		default:
			return NewError("cannot convert %s to number", right.Type())
		}
	} else if leftVal, ok := left.(*object.Double); ok {
		switch rightVal := right.(type) {
		case *object.Integer:
			return NewBoolean(leftVal.Value() > float64(rightVal.Value()))
		case *object.Decimal:
			return NewBoolean(leftVal.Value() > rightVal.Value())
		case *object.Double:
			return NewBoolean(leftVal.Value() > rightVal.Value())
		default:
			return NewError("cannot convert %s to number", right.Type())
		}
	} else if leftVal, ok := left.(*object.String); ok {
		switch rightVal := right.(type) {
		case *object.String:
			return NewBoolean(leftVal.Value() > rightVal.Value())
		default:
			return NewError("cannot compare %s and %s", left.Type(), right.Type())
		}
	}
	return NewError("cannot compare %s and %s", left.Type(), right.Type())
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
		default:
			return NewError("cannot convert %s to number", right.Type())
		}
	} else if leftVal, ok := left.(*object.Decimal); ok {
		switch rightVal := right.(type) {
		case *object.Integer:
			return NewBoolean(leftVal.Value() >= float64(rightVal.Value()))
		case *object.Decimal:
			return NewBoolean(leftVal.Value() >= rightVal.Value())
		case *object.Double:
			return NewBoolean(leftVal.Value() >= rightVal.Value())
		default:
			return NewError("cannot convert %s to number", right.Type())
		}
	} else if leftVal, ok := left.(*object.Double); ok {
		switch rightVal := right.(type) {
		case *object.Integer:
			return NewBoolean(leftVal.Value() >= float64(rightVal.Value()))
		case *object.Decimal:
			return NewBoolean(leftVal.Value() >= rightVal.Value())
		case *object.Double:
			return NewBoolean(leftVal.Value() >= rightVal.Value())
		default:
			return NewError("cannot convert %s to number", right.Type())
		}
	} else if leftVal, ok := left.(*object.String); ok {
		switch rightVal := right.(type) {
		case *object.String:
			return NewBoolean(leftVal.Value() >= rightVal.Value())
		default:
			return NewError("cannot compare %s and %s", left.Type(), right.Type())
		}
	}
	return NewError("cannot compare %s and %s", left.Type(), right.Type())
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
		return true
	}
	return false
}

// IsContain ..
func IsContain(src []object.Item, target object.Item) bool {
	for _, item := range src {
		if item == target {
			return true
		}
	}
	return false
}

// IsContainN ..
func IsContainN(src []object.Node, target object.Node) bool {
	for _, item := range src {
		if item, ok := item.(*object.BaseNode); ok {
			if item.Type() != target.Type() {
				return false
			}

			target := target.(*object.BaseNode)
			if item.Tree() == target.Tree() {
				return true
			}
		}

		if item, ok := item.(*object.AttrNode); ok {
			if item.Type() != target.Type() {
				return false
			}

			target := target.(*object.AttrNode)
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

// AppendKind ..
func AppendKind(src []object.Node, target object.Node, typeID byte) []object.Node {
	switch typeID {
	case 2:
		if target.Type() == object.ElementNodeType {
			src = AppendNode(src, target)
		}
	case 3:
		if target.Type() == object.ElementNodeType {
			target := target.(*object.BaseNode)
			for _, a := range target.Attr() {
				src = AppendNode(src, a)
			}
		}
	case 7:
		if target.Type() == object.CommentNodeType {
			src = AppendNode(src, target)
		}
	case 8:
		if target.Type() == object.TextNodeType {
			src = AppendNode(src, target)
		}
	case 10:
		src = AppendNode(src, target)
		if target.Type() == object.ElementNodeType {
			target := target.(*object.BaseNode)
			for _, a := range target.Attr() {
				src = AppendNode(src, a)
			}
		}
	}
	return src
}

// WalkDescKind ..
func WalkDescKind(nodes []object.Node, n object.Node, typeID byte) []object.Node {
	for c := n.FirstChild(); c != nil; c = c.NextSibling() {
		nodes = AppendKind(nodes, c, typeID)
		if c.FirstChild() != nil {
			nodes = WalkDescKind(nodes, c, typeID)
		}
	}
	return nodes
}

// WalkDescName ..
func WalkDescName(nodes []object.Node, n object.Node, t *ast.NameTest) []object.Node {
	for c := n.FirstChild(); c != nil; c = c.NextSibling() {
		if c.Type() == object.ElementNodeType {
			switch t.TypeID {
			case 1:
				if c.Tree().Data == t.EQName.Value() {
					nodes = AppendNode(nodes, c)
				}
			case 2:
				switch t.Wildcard.TypeID {
				case 1:
					nodes = AppendNode(nodes, c)
				}
			}

		}

		if c.FirstChild() != nil {
			nodes = WalkDescName(nodes, c, t)
		}
	}
	return nodes
}
