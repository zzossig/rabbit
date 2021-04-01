package rabbit

import (
	"fmt"
	"strings"

	"github.com/zzossig/rabbit/bif"
	"github.com/zzossig/rabbit/eval"
	"github.com/zzossig/rabbit/lexer"
	"github.com/zzossig/rabbit/object"
	"github.com/zzossig/rabbit/parser"
)

var context *object.Context

// XPath function evaluates xpath expression with a context.
// The ctx param can be nil. In that case,  the global context is used.
// If global context also nil, you can't use node-related xpath expressions.
func XPath(ctx *object.Context, input string) (object.Item, error) {
	l := lexer.New(input)
	p := parser.New(l)
	px := p.ParseXPath()

	if len(p.Errors()) != 0 {
		return nil, createError(p.Errors())
	}

	var e object.Item
	if ctx != nil {
		e = eval.Eval(px, ctx)
	} else if context != nil {
		e = eval.Eval(px, context)
	} else {
		e = eval.Eval(px, nil)
	}

	if bif.IsError(e) {
		return nil, fmt.Errorf(e.Inspect())
	}

	return e, nil
}

// NewContext create new context for using evaluate xpath expressions.
// input param can be a url or local file path.
// returned context is used as the XPath function's first parameter.
func NewContext(input string) (*object.Context, error) {
	ctx := object.NewContext()
	f := bif.F["fn:doc"]

	err := f(ctx, bif.NewString(input))
	if err != nil {
		return nil, fmt.Errorf(err.Inspect())
	}

	return ctx, nil
}

// SetGlobalConetxt set global context so when you call XPath function without ctx param, global context is used as default.
func SetGlobalConetxt(input string) error {
	ctx := object.NewContext()
	f := bif.F["fn:doc"]

	err := f(ctx, bif.NewString(input))
	if err != nil {
		return fmt.Errorf(err.Inspect())
	}

	context = ctx
	return nil
}

// Convert convert object.Item to a golang data type
func Convert(item object.Item) (interface{}, error) {
	switch item := item.(type) {
	case *object.Integer:
		return item.Value(), nil
	case *object.Decimal:
		return item.Value(), nil
	case *object.Double:
		return item.Value(), nil
	case *object.Boolean:
		return item.Value(), nil
	case *object.String:
		return item.Value(), nil
	case *object.BaseNode:
		return item.Tree(), nil
	case *object.AttrNode:
		return item.Self(), nil
	case *object.Map:
		return convertMap(item)
	case *object.Array:
		return convertArray(item)
	case *object.Sequence:
		return convertSequence(item)
	}
	return nil, fmt.Errorf("cannot convert item: %v", item)
}

func convertMap(m *object.Map) (map[interface{}]interface{}, error) {
	mm := make(map[interface{}]interface{}, len(m.Pairs))
	for _, pair := range m.Pairs {
		k, err := Convert(pair.Key)
		if err != nil {
			return nil, err
		}
		v, err := Convert(pair.Value)
		if err != nil {
			return nil, err
		}
		mm[k] = v
	}
	return mm, nil
}

func convertArray(a *object.Array) ([]interface{}, error) {
	aa := make([]interface{}, 0, len(a.Items))
	for _, v := range a.Items {
		v, err := Convert(v)
		if err != nil {
			return nil, err
		}
		aa = append(aa, v)
	}
	return aa, nil
}

func convertSequence(s *object.Sequence) ([]interface{}, error) {
	ss := make([]interface{}, 0, len(s.Items))
	for _, v := range s.Items {
		v, err := Convert(v)
		if err != nil {
			return nil, err
		}
		ss = append(ss, v)
	}
	return ss, nil
}

func createError(errors []error) error {
	var sb strings.Builder
	for i, e := range errors {
		sb.WriteString(e.Error())
		if i < len(errors)-1 {
			sb.WriteString("\n")
		}
	}
	return fmt.Errorf(sb.String())
}
