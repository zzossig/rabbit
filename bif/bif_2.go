package bif

import "github.com/zzossig/xpath/object"

func fnNodeName(ctx *object.Context, args ...object.Item) object.Item {
	if len(args) > 0 {
		NewError("wrong number of arguments. got=%d, expected=0", len(args))
	}

	seq := &object.Sequence{}

	for _, n := range ctx.CNode {
		if n.Type() == object.ElementNodeType {
			seq.Items = append(seq.Items, NewString(n.Tree().Data))
		}
	}

	return seq
}

func fnString(ctx *object.Context, args ...object.Item) object.Item {
	// if len(args) > 1 {
	// 	NewError("wrong number of arguments. got=%d, expected=0 or 1", len(args))
	// } else if len(args) == 1 {

	// } else {

	// }

	return nil
}
