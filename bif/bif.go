package bif

import (
	"fmt"

	"github.com/zzossig/xpath/object"
)

// Builtins ..
var Builtins = map[string]object.Func{
	"abs":           abs,
	"concat":        concat,
	"for-each-pair": forEachPair,
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
