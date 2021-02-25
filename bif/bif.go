package bif

import (
	"fmt"

	"github.com/zzossig/xpath/object"
)

// Builtins defined in https://www.w3.org/TR/xpath-functions-31/
var Builtins = map[string]object.Func{
	// "op:numeric-add":   numericAdd,
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
	}

	return pnames
}

// BuiltinPTypes ..
func BuiltinPTypes(name string, num int) []object.Type {
	var ptypes []object.Type

	switch name {
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
				if args[i].Type() != object.MapType {
					isMatch = false
				}
			case object.ArrayType:
				if args[i].Type() != object.ArrayType {
					isMatch = false
				}
			case object.NodeType:
				if !IsNode(args[i]) {
					isMatch = false
				}
			case object.StringType:
				if args[i].Type() != object.StringType {
					isMatch = false
				}
			case object.BooleanType:
				if args[i].Type() != object.BooleanType {
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
	if item != nil {
		return item.Type() == object.ErrorType
	}
	return false
}

// IsSeq ..
func IsSeq(item object.Item) bool {
	if _, ok := item.(*object.Sequence); ok {
		return true
	}
	return false
}

// IsPlaceholder ..
func IsPlaceholder(item object.Item) bool {
	if _, ok := item.(*object.Sequence); ok {
		return true
	}
	return false
}

// IsNumeric ..
func IsNumeric(item object.Item) bool {
	if item.Type() == object.IntegerType ||
		item.Type() == object.DecimalType ||
		item.Type() == object.DoubleType {
		return true
	}
	return false
}

// IsFunction ..
func IsFunction(item object.Item) bool {
	if item.Type() == object.FuncType {
		return true
	}
	return false
}

// IsAnyAtomic ..
func IsAnyAtomic(item object.Item) bool {
	if item.Type() == object.DoubleType ||
		item.Type() == object.DecimalType ||
		item.Type() == object.IntegerType ||
		item.Type() == object.StringType ||
		item.Type() == object.BooleanType {
		return true
	}
	return false
}

// IsNode ..
func IsNode(item object.Item) bool {
	if item.Type() == object.NodeType ||
		item.Type() == object.DocType ||
		item.Type() == object.ElemType ||
		item.Type() == object.AttrType ||
		item.Type() == object.PIType ||
		item.Type() == object.CommentType ||
		item.Type() == object.NSNodeType ||
		item.Type() == object.TextType {
		return true
	}
	return false
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
