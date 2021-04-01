package rabbit

import (
	"fmt"

	"github.com/zzossig/rabbit/bif"
	"github.com/zzossig/rabbit/eval"
	"github.com/zzossig/rabbit/lexer"
	"github.com/zzossig/rabbit/object"
	"github.com/zzossig/rabbit/parser"
)

type xpath struct {
	input   string
	context *object.Context
	evaled  object.Item
	errors  []error
}

// New creates new xpath object.
func New() *xpath {
	return &xpath{context: object.NewContext()}
}

// SetDoc set document to a context.
// if document is not set in a context, node related xpath expressions not going to work.
// input param can be url or local filepath.
func (x *xpath) SetDoc(input string) *xpath {
	f := bif.F["fn:doc"]

	err := f(x.context, bif.NewString(input))
	if err != nil {
		x.errors = append(x.errors, fmt.Errorf(err.Inspect()))
	}
	return x
}

// Eval evaluates a xpath expression and save the result to xpath.
func (x *xpath) Eval(input string) *xpath {
	if len(x.errors) > 0 {
		return x
	}

	x.input = input
	l := lexer.New(input)
	p := parser.New(l)
	px := p.ParseXPath()

	if len(p.Errors()) != 0 {
		x.errors = append(x.errors, p.Errors()...)
		return x
	}

	e := eval.Eval(px, x.context)
	if bif.IsError(e) {
		x.errors = append(x.errors, fmt.Errorf(e.Inspect()))
		return x
	}
	x.evaled = e

	return x
}

// Data convert evaled field to a golang data type
func (x *xpath) Data() []interface{} {
	if x.evaled == nil {
		x.errors = append(x.errors, fmt.Errorf("cannot convert item since evaled field is nil"))
		return nil
	}

	e, err := convert(x.evaled)
	if err != nil {
		x.errors = append(x.errors, err)
		return nil
	}

	// evaled field always a sequence type so converted value always be a []interface{} type
	return e.([]interface{})
}

func convert(item object.Item) (interface{}, error) {
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
		k, err := convert(pair.Key)
		if err != nil {
			return nil, err
		}
		v, err := convert(pair.Value)
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
		v, err := convert(v)
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
		v, err := convert(v)
		if err != nil {
			return nil, err
		}
		ss = append(ss, v)
	}
	return ss, nil
}
