package rabbit

import (
	"fmt"

	"github.com/zzossig/rabbit/object"
	"golang.org/x/net/html"
)

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

func convertString(item object.Item) []string {
	var s []string
	switch item := item.(type) {
	case *object.Sequence:
		for _, i := range item.Items {
			ss := convertString(i)
			s = append(s, ss...)
		}
	case *object.Array:
		for _, i := range item.Items {
			ss := convertString(i)
			s = append(s, ss...)
		}
	case *object.BaseNode:
		s = append(s, item.Text())
	case *object.AttrNode:
		s = append(s, item.Attr().Val)
	default:
		s = append(s, item.Inspect())
	}
	return s
}

func initContext(ctx *object.Context) {
	ctx.CSize = 0
	ctx.CPos = 0
	ctx.CAxis = ""
	ctx.CItem = nil
	ctx.CNode = []object.Node{}
}

func copyContext(ctx *object.Context) *object.Context {
	c := object.NewContext()
	c.Doc = ctx.Doc
	return c
}

func copyContextN(ctx *object.Context, n object.Node) *object.Context {
	c := object.NewContext()
	c.Doc = ctx.Doc
	c.CNode = append(c.CNode, n)
	return c
}
