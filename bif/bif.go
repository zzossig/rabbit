package bif

import (
	"fmt"

	"github.com/zzossig/xpath/object"
)

// predefined
var (
	NIL   = &object.Nil{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
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
