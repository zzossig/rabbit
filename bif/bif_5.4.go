package bif

import (
	"strings"

	"github.com/zzossig/xpath/object"
)

func concat(args ...object.Item) object.Item {
	var sb strings.Builder
	for _, arg := range args {
		sb.WriteString(arg.Inspect())
	}
	return &object.String{Value: sb.String()}
}
