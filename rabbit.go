package rabbit

import (
	"fmt"
	"os"

	"github.com/zzossig/rabbit/bif"
	"github.com/zzossig/rabbit/eval"
	"github.com/zzossig/rabbit/lexer"
	"github.com/zzossig/rabbit/object"
	"github.com/zzossig/rabbit/parser"
	"github.com/zzossig/rabbit/repl"
	"golang.org/x/net/html"
)

// XPath is a base object to evaluate xpath expressions.
// input is xpath expression that is saved when you are using Eval method.
// context is a context that contains a document.
// SetDoc function saves a document to the context.
// evaled field is set when you call Eval method.
// evaled type is an object.Item which is a custom data type used in rabbit language.
// You can convert evaled type to a golang data type using Data or Nodes method.
// errors field is a collected errors while parsing and evaluating
type XPath struct {
	input   string
	context *object.Context
	evaled  object.Item
	errors  []error
}

// New creates new xpath object.
func New() *XPath {
	return &XPath{context: object.NewContext()}
}

// SetDoc set document to a context.
// if document is not set in a context, node related xpath expressions not going to work.
// input param can be url or local filepath.
func (x *XPath) SetDoc(input string) *XPath {
	f := bif.F["fn:doc"]

	err := f(x.context, bif.NewString(input))
	if err != nil {
		x.errors = append(x.errors, fmt.Errorf(err.Inspect()))
	}
	return x
}

// Eval evaluates a xpath expression and save the result to evaled field.
func (x *XPath) Eval(input string) *XPath {
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

// Data convert evaled field to []interface{}
func (x *XPath) Data() []interface{} {
	initContext(x.context)

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

// Nodes convert evaled field to []*html.Node
func (x *XPath) Nodes() []*html.Node {
	initContext(x.context)

	if x.evaled == nil {
		x.errors = append(x.errors, fmt.Errorf("cannot convert item since evaled field is nil"))
		return nil
	}

	e, err := convertNode(x.evaled)
	if err != nil {
		x.errors = append(x.errors, err)
		return nil
	}

	return e
}

// Raw returns evaled field
func (x *XPath) Raw() object.Item {
	initContext(x.context)
	return x.evaled
}

// Errors returns errors field
func (x *XPath) Errors() []error {
	return x.errors
}

// String returns input field
func (x *XPath) String() string {
	return x.input
}

// CLI is a command line interface
func (x *XPath) CLI() {
	repl.Start(os.Stdin, os.Stdout, x.context)
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

func convertNode(item object.Item) ([]*html.Node, error) {
	switch item := item.(type) {
	case *object.Sequence:
		nodes := make([]*html.Node, 0, len(item.Items))
		for _, i := range item.Items {
			switch i := i.(type) {
			case *object.Sequence:
				n, err := convertNode(i)
				if err != nil {
					return nodes, fmt.Errorf("unknown node type: %s", i.Type())
				}
				nodes = append(nodes, n...)
			case *object.Array:
				n, err := convertNode(i)
				if err != nil {
					return nodes, fmt.Errorf("unknown node type: %s", i.Type())
				}
				nodes = append(nodes, n...)
			case *object.BaseNode:
				nodes = append(nodes, i.Self())
			case *object.AttrNode:
				nodes = append(nodes, i.Self())
			default:
				return nodes, fmt.Errorf("unknown node type: %s", i.Type())
			}
		}
		return nodes, nil
	case *object.Array:
		nodes := make([]*html.Node, 0, len(item.Items))
		for _, i := range item.Items {
			switch i := i.(type) {
			case *object.Sequence:
				n, err := convertNode(i)
				if err != nil {
					return nodes, fmt.Errorf("unknown node type: %s", i.Type())
				}
				nodes = append(nodes, n...)
			case *object.Array:
				n, err := convertNode(i)
				if err != nil {
					return nodes, fmt.Errorf("unknown node type: %s", i.Type())
				}
				nodes = append(nodes, n...)
			case *object.BaseNode:
				nodes = append(nodes, i.Self())
			case *object.AttrNode:
				nodes = append(nodes, i.Self())
			default:
				return nodes, fmt.Errorf("unknown node type: %s", i.Type())
			}
		}
		return nodes, nil
	case *object.BaseNode:
		return []*html.Node{item.Self()}, nil
	case *object.AttrNode:
		return []*html.Node{item.Self()}, nil
	default:
		return []*html.Node{}, fmt.Errorf("unknown node type: %s", item.Type())
	}
}

func initContext(ctx *object.Context) {
	ctx.CSize = 0
	ctx.CPos = 0
	ctx.CAxis = ""
	ctx.CItem = nil
	ctx.CNode = []object.Node{}
}
