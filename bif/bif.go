package bif

import (
	"fmt"

	"github.com/zzossig/xpath/object"
)

// Builtins defined in https://www.w3.org/TR/xpath-functions-31/
var Builtins = map[string]object.Func{
	"abs":           abs,
	"concat":        concat,
	"for-each-pair": forEachPair,
	"upper-case":    upperCase,
	"lower-case":    lowerCase,
	"boolean":       boolean,
}

// NewError ..
func NewError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

// IsError ..
func IsError(item object.Item) bool {
	if item != nil {
		return item.Type() == object.ErrorType
	}
	return false
}

func isSeq(item object.Item) bool {
	if _, ok := item.(*object.Sequence); ok {
		return true
	}
	return false
}

func isPlaceholder(item object.Item) bool {
	if _, ok := item.(*object.Sequence); ok {
		return true
	}
	return false
}

// IsEQ ..
func IsEQ(left, right object.Item) object.Item {
	if leftVal, ok := left.(*object.Integer); ok {
		switch rightVal := right.(type) {
		case *object.Integer:
			return &object.Boolean{Value: leftVal.Value == rightVal.Value}
		case *object.Decimal:
			return &object.Boolean{Value: float64(leftVal.Value) == rightVal.Value}
		case *object.Double:
			return &object.Boolean{Value: float64(leftVal.Value) == rightVal.Value}
		default:
			return NewError("cannot convert %s to number", right.Type())
		}
	} else if leftVal, ok := left.(*object.Decimal); ok {
		switch rightVal := right.(type) {
		case *object.Integer:
			return &object.Boolean{Value: leftVal.Value == float64(rightVal.Value)}
		case *object.Decimal:
			return &object.Boolean{Value: leftVal.Value == rightVal.Value}
		case *object.Double:
			return &object.Boolean{Value: leftVal.Value == rightVal.Value}
		default:
			return NewError("cannot convert %s to number", right.Type())
		}
	} else if leftVal, ok := left.(*object.Double); ok {
		switch rightVal := right.(type) {
		case *object.Integer:
			return &object.Boolean{Value: leftVal.Value == float64(rightVal.Value)}
		case *object.Decimal:
			return &object.Boolean{Value: leftVal.Value == rightVal.Value}
		case *object.Double:
			return &object.Boolean{Value: leftVal.Value == rightVal.Value}
		default:
			return NewError("cannot convert %s to number", right.Type())
		}
	} else if leftVal, ok := left.(*object.String); ok {
		switch rightVal := right.(type) {
		case *object.String:
			return &object.Boolean{Value: leftVal.Value == rightVal.Value}
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
			return &object.Boolean{Value: leftVal.Value != rightVal.Value}
		case *object.Decimal:
			return &object.Boolean{Value: float64(leftVal.Value) != rightVal.Value}
		case *object.Double:
			return &object.Boolean{Value: float64(leftVal.Value) != rightVal.Value}
		default:
			return NewError("cannot convert %s to number", right.Type())
		}
	} else if leftVal, ok := left.(*object.Decimal); ok {
		switch rightVal := right.(type) {
		case *object.Integer:
			return &object.Boolean{Value: leftVal.Value != float64(rightVal.Value)}
		case *object.Decimal:
			return &object.Boolean{Value: leftVal.Value != rightVal.Value}
		case *object.Double:
			return &object.Boolean{Value: leftVal.Value != rightVal.Value}
		default:
			return NewError("cannot convert %s to number", right.Type())
		}
	} else if leftVal, ok := left.(*object.Double); ok {
		switch rightVal := right.(type) {
		case *object.Integer:
			return &object.Boolean{Value: leftVal.Value != float64(rightVal.Value)}
		case *object.Decimal:
			return &object.Boolean{Value: leftVal.Value != rightVal.Value}
		case *object.Double:
			return &object.Boolean{Value: leftVal.Value != rightVal.Value}
		default:
			return NewError("cannot convert %s to number", right.Type())
		}
	} else if leftVal, ok := left.(*object.String); ok {
		switch rightVal := right.(type) {
		case *object.String:
			return &object.Boolean{Value: leftVal.Value != rightVal.Value}
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
			return &object.Boolean{Value: leftVal.Value < rightVal.Value}
		case *object.Decimal:
			return &object.Boolean{Value: float64(leftVal.Value) < rightVal.Value}
		case *object.Double:
			return &object.Boolean{Value: float64(leftVal.Value) < rightVal.Value}
		default:
			return NewError("cannot convert %s to number", right.Type())
		}
	} else if leftVal, ok := left.(*object.Decimal); ok {
		switch rightVal := right.(type) {
		case *object.Integer:
			return &object.Boolean{Value: leftVal.Value < float64(rightVal.Value)}
		case *object.Decimal:
			return &object.Boolean{Value: leftVal.Value < rightVal.Value}
		case *object.Double:
			return &object.Boolean{Value: leftVal.Value < rightVal.Value}
		default:
			return NewError("cannot convert %s to number", right.Type())
		}
	} else if leftVal, ok := left.(*object.Double); ok {
		switch rightVal := right.(type) {
		case *object.Integer:
			return &object.Boolean{Value: leftVal.Value < float64(rightVal.Value)}
		case *object.Decimal:
			return &object.Boolean{Value: leftVal.Value < rightVal.Value}
		case *object.Double:
			return &object.Boolean{Value: leftVal.Value < rightVal.Value}
		default:
			return NewError("cannot convert %s to number", right.Type())
		}
	} else if leftVal, ok := left.(*object.String); ok {
		switch rightVal := right.(type) {
		case *object.String:
			return &object.Boolean{Value: leftVal.Value < rightVal.Value}
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
			return &object.Boolean{Value: leftVal.Value <= rightVal.Value}
		case *object.Decimal:
			return &object.Boolean{Value: float64(leftVal.Value) <= rightVal.Value}
		case *object.Double:
			return &object.Boolean{Value: float64(leftVal.Value) <= rightVal.Value}
		default:
			return NewError("cannot convert %s to number", right.Type())
		}
	} else if leftVal, ok := left.(*object.Decimal); ok {
		switch rightVal := right.(type) {
		case *object.Integer:
			return &object.Boolean{Value: leftVal.Value <= float64(rightVal.Value)}
		case *object.Decimal:
			return &object.Boolean{Value: leftVal.Value <= rightVal.Value}
		case *object.Double:
			return &object.Boolean{Value: leftVal.Value <= rightVal.Value}
		default:
			return NewError("cannot convert %s to number", right.Type())
		}
	} else if leftVal, ok := left.(*object.Double); ok {
		switch rightVal := right.(type) {
		case *object.Integer:
			return &object.Boolean{Value: leftVal.Value <= float64(rightVal.Value)}
		case *object.Decimal:
			return &object.Boolean{Value: leftVal.Value <= rightVal.Value}
		case *object.Double:
			return &object.Boolean{Value: leftVal.Value <= rightVal.Value}
		default:
			return NewError("cannot convert %s to number", right.Type())
		}
	} else if leftVal, ok := left.(*object.String); ok {
		switch rightVal := right.(type) {
		case *object.String:
			return &object.Boolean{Value: leftVal.Value <= rightVal.Value}
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
			return &object.Boolean{Value: leftVal.Value > rightVal.Value}
		case *object.Decimal:
			return &object.Boolean{Value: float64(leftVal.Value) > rightVal.Value}
		case *object.Double:
			return &object.Boolean{Value: float64(leftVal.Value) > rightVal.Value}
		default:
			return NewError("cannot convert %s to number", right.Type())
		}
	} else if leftVal, ok := left.(*object.Decimal); ok {
		switch rightVal := right.(type) {
		case *object.Integer:
			return &object.Boolean{Value: leftVal.Value > float64(rightVal.Value)}
		case *object.Decimal:
			return &object.Boolean{Value: leftVal.Value > rightVal.Value}
		case *object.Double:
			return &object.Boolean{Value: leftVal.Value > rightVal.Value}
		default:
			return NewError("cannot convert %s to number", right.Type())
		}
	} else if leftVal, ok := left.(*object.Double); ok {
		switch rightVal := right.(type) {
		case *object.Integer:
			return &object.Boolean{Value: leftVal.Value > float64(rightVal.Value)}
		case *object.Decimal:
			return &object.Boolean{Value: leftVal.Value > rightVal.Value}
		case *object.Double:
			return &object.Boolean{Value: leftVal.Value > rightVal.Value}
		default:
			return NewError("cannot convert %s to number", right.Type())
		}
	} else if leftVal, ok := left.(*object.String); ok {
		switch rightVal := right.(type) {
		case *object.String:
			return &object.Boolean{Value: leftVal.Value > rightVal.Value}
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
			return &object.Boolean{Value: leftVal.Value >= rightVal.Value}
		case *object.Decimal:
			return &object.Boolean{Value: float64(leftVal.Value) >= rightVal.Value}
		case *object.Double:
			return &object.Boolean{Value: float64(leftVal.Value) >= rightVal.Value}
		default:
			return NewError("cannot convert %s to number", right.Type())
		}
	} else if leftVal, ok := left.(*object.Decimal); ok {
		switch rightVal := right.(type) {
		case *object.Integer:
			return &object.Boolean{Value: leftVal.Value >= float64(rightVal.Value)}
		case *object.Decimal:
			return &object.Boolean{Value: leftVal.Value >= rightVal.Value}
		case *object.Double:
			return &object.Boolean{Value: leftVal.Value >= rightVal.Value}
		default:
			return NewError("cannot convert %s to number", right.Type())
		}
	} else if leftVal, ok := left.(*object.Double); ok {
		switch rightVal := right.(type) {
		case *object.Integer:
			return &object.Boolean{Value: leftVal.Value >= float64(rightVal.Value)}
		case *object.Decimal:
			return &object.Boolean{Value: leftVal.Value >= rightVal.Value}
		case *object.Double:
			return &object.Boolean{Value: leftVal.Value >= rightVal.Value}
		default:
			return NewError("cannot convert %s to number", right.Type())
		}
	} else if leftVal, ok := left.(*object.String); ok {
		switch rightVal := right.(type) {
		case *object.String:
			return &object.Boolean{Value: leftVal.Value >= rightVal.Value}
		default:
			return NewError("cannot compare %s and %s", left.Type(), right.Type())
		}
	}
	return NewError("cannot compare %s and %s", left.Type(), right.Type())
}
