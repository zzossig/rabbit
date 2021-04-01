package bif

import (
	"strings"

	"github.com/zzossig/rabbit/object"
)

func fnNodeName(ctx *object.Context, args ...object.Item) object.Item {
	if len(args) > 0 {
		return NewError("too many parameters for function call: fn:node-name")
	}

	seq := &object.Sequence{}

	for _, n := range ctx.CNode {
		if n.Type() == object.ElementNodeType {
			seq.Items = append(seq.Items, NewString(n.Tree().Data))
		}
	}

	return seq
}

// include whitespace
func fnString(ctx *object.Context, args ...object.Item) object.Item {
	if len(args) > 1 {
		return NewError("too many parameters for function call: fn:string")
	}

	if len(args) == 1 {
		return CastType(args[0], object.StringType)
	}

	seq := &object.Sequence{}

	if len(ctx.CNode) > 0 {
		for _, n := range ctx.CNode {
			texts := collectText(nil, n)
			str := combineTexts(texts)
			seq.Items = append(seq.Items, NewString(str))
		}
		return seq
	}

	return NewError("context node is undefined")
}

// exclude whitespace
func fnData(ctx *object.Context, args ...object.Item) object.Item {
	if len(args) > 1 {
		return NewError("too many parameters for function call: fn:data")
	}

	if len(args) == 1 {
		return CastType(args[0], object.StringType)
	}

	seq := &object.Sequence{}

	if len(ctx.CNode) > 0 {
		for _, n := range ctx.CNode {
			texts := collectText(nil, n)
			str := combineTextsTrim(texts)
			seq.Items = append(seq.Items, NewString(str))
		}
		return seq
	}

	if ctx.Doc != nil {
		texts := collectText(nil, ctx.Doc)
		str := combineTextsTrim(texts)
		seq.Items = append(seq.Items, NewString(str))
		return seq
	}

	return NewError("context node is undefined")
}

func fnBaseURI(ctx *object.Context, args ...object.Item) object.Item {
	return NewString(ctx.BaseURI)
}

func collectText(texts []string, n object.Node) []string {
	if n.Type() == object.ElementNodeType || n.Type() == object.AttributeNodeType {
		texts = append(texts, n.Text())
	}
	for c := n.FirstChild(); c != nil; c = c.NextSibling() {
		texts = collectText(texts, c)
	}
	return texts
}

func combineTexts(texts []string) string {
	var sb strings.Builder
	for _, text := range texts {
		sb.WriteString(text)
	}
	return sb.String()
}

func combineTextsTrim(texts []string) string {
	var sb strings.Builder
	for _, text := range texts {
		sb.WriteString(strings.TrimSpace(text))
	}
	return sb.String()
}

func combineTextsNormalize(texts []string) string {
	var trimedTexts []string
	for _, text := range texts {
		str := strings.TrimSpace(text)
		if str != "" {
			trimedTexts = append(trimedTexts, str)
		}
	}
	return strings.Join(trimedTexts, " ")
}
