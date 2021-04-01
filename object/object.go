package object

import (
	"fmt"
	"hash/fnv"
	"math"
	"strings"

	"github.com/zzossig/rabbit/ast"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

// Item ::= node | function(*) | xs:anyAtomicType
type Item interface {
	Type() Type
	Inspect() string
}

// Node ::= document, element, attribute, comment, namespace, processing-instruction, text
type Node interface {
	Item
	Tree() *html.Node
	SetTree(t *html.Node)
	Self() *html.Node
	Parent() Node
	FirstChild() Node
	LastChild() Node
	PrevSibling() Node
	NextSibling() Node
	Text() string
}

// Error is an item that is represents error when doing evaluation
type Error struct {
	Message string
}

func (e *Error) Type() Type      { return ErrorType }
func (e *Error) Inspect() string { return "ERROR: " + e.Message }

// Placeholder is an item that is represents ?(question token) when doing evaluation
type Placeholder struct{}

func (p *Placeholder) Type() Type      { return PholderType }
func (p *Placeholder) Inspect() string { return "?" }

// Varref is an item that is represents $var when doing evaluation
type Varref struct {
	Name ast.EQName
}

func (v *Varref) Type() Type      { return VarrefType }
func (v *Varref) Inspect() string { return fmt.Sprintf("$%s", v.Name.Value()) }

// Sequence is an ordered collection of zero or more items.
type Sequence struct {
	Items []Item
}

func (s *Sequence) Type() Type { return SequenceType }
func (s *Sequence) Inspect() string {
	var sb strings.Builder

	sb.WriteString("(")
	for i, item := range s.Items {
		sb.WriteString(item.Inspect())
		if i < len(s.Items)-1 {
			sb.WriteString(", ")
		}
	}
	sb.WriteString(")")

	return sb.String()
}

// Hasher implemented in atomic types: integer, decimal, double, boolean, string
// So, atomic types are used as a Map key
type Hasher interface {
	HashKey() HashKey
}

// HashKey is used as a Map key
type HashKey struct {
	Type
	Value uint64
}

// Pair contains key, value pair
type Pair struct {
	Key   Item
	Value Item
}

// Map is an item that is represents map data-type
type Map struct {
	Pairs map[HashKey]Pair
}

func (m *Map) Type() Type { return MapType }
func (m *Map) Inspect() string {
	var sb strings.Builder

	pairs := []string{}
	for _, pair := range m.Pairs {
		pairs = append(pairs, fmt.Sprintf("%s: %s", pair.Key.Inspect(), pair.Value.Inspect()))
	}

	sb.WriteString("map")
	sb.WriteString("{")
	sb.WriteString(strings.Join(pairs, ", "))
	sb.WriteString("}")

	return sb.String()
}

// Array is an item that is represents array data-type
type Array struct {
	Items []Item
}

func (a *Array) Type() Type { return ArrayType }
func (a *Array) Inspect() string {
	var sb strings.Builder

	sb.WriteString("[")
	for i, item := range a.Items {
		sb.WriteString(item.Inspect())
		if i < len(a.Items)-1 {
			sb.WriteString(", ")
		}
	}
	sb.WriteString("]")

	return sb.String()
}

// FuncNamed ::= ns:bif#n
type FuncNamed struct {
	Name ast.EQName
	Num  int
	*Context
	*Func
}

func (fn *FuncNamed) Type() Type      { return FuncType }
func (fn *FuncNamed) Inspect() string { return fmt.Sprintf("%s#%d", fn.Name.Value(), fn.Num) }

// FuncInline ::= function() {}
type FuncInline struct {
	PL   *ast.ParamList
	Body *ast.EnclosedExpr
	Fn   Ev
	*Context
}

func (fi *FuncInline) Type() Type      { return FuncType }
func (fi *FuncInline) Inspect() string { return "function" }

// FuncPartial ::= ns:bif(?,...)
type FuncPartial struct {
	Name ast.EQName
	Args []Item
	PCnt int
	*Context
	*Func
}

func (fp *FuncPartial) Type() Type      { return FuncType }
func (fp *FuncPartial) Inspect() string { return "function" }

// Integer is an item that is represents int data-type
type Integer struct {
	value int
}

func (i *Integer) Type() Type      { return IntegerType }
func (i *Integer) Inspect() string { return fmt.Sprintf("%d", i.value) }
func (i *Integer) SetValue(v int)  { i.value = v }
func (i *Integer) Value() int      { return i.value }
func (i *Integer) HashKey() HashKey {
	return HashKey{Type: i.Type(), Value: uint64(i.value)}
}

// Decimal is an item that is represents float64 data-type
// number token that is not contains `e` or `E` is evaluated to Decimal
type Decimal struct {
	value float64
}

func (d *Decimal) Type() Type         { return DecimalType }
func (d *Decimal) Inspect() string    { return fmt.Sprintf("%f", d.value) }
func (d *Decimal) SetValue(v float64) { d.value = v }
func (d *Decimal) Value() float64     { return d.value }
func (d *Decimal) HashKey() HashKey {
	return HashKey{Type: d.Type(), Value: math.Float64bits(d.value)}
}

// Double is an item that is represents float64 data-type
// number token that contains `e` or `E` is evaluated to Double
type Double struct {
	value float64
}

func (d *Double) Type() Type         { return DoubleType }
func (d *Double) Inspect() string    { return fmt.Sprintf("%e", d.value) }
func (d *Double) SetValue(v float64) { d.value = v }
func (d *Double) Value() float64     { return d.value }
func (d *Double) HashKey() HashKey {
	return HashKey{Type: d.Type(), Value: math.Float64bits(d.value)}
}

// Boolean is an item that is represents bool data-type
type Boolean struct {
	value bool
}

func (b *Boolean) Type() Type      { return BooleanType }
func (b *Boolean) Inspect() string { return fmt.Sprintf("%t", b.value) }
func (b *Boolean) SetValue(v bool) { b.value = v }
func (b *Boolean) Value() bool     { return b.value }
func (b *Boolean) HashKey() HashKey {
	var value uint64

	if b.value {
		value = 1
	} else {
		value = 0
	}

	return HashKey{Type: b.Type(), Value: value}
}

// String is an item that is represents string data-type
type String struct {
	value string
}

func (s *String) Type() Type        { return StringType }
func (s *String) Inspect() string   { return s.value }
func (s *String) SetValue(v string) { s.value = v }
func (s *String) Value() string     { return s.value }
func (s *String) HashKey() HashKey {
	h := fnv.New64a()
	h.Write([]byte(s.value))

	return HashKey{Type: s.Type(), Value: h.Sum64()}
}

// BaseNode ::= ElementNode | TextNode | DocumentNode | CommentNode | DoctypeNode
// BaseNode just wraps *html.Node
type BaseNode struct {
	tree *html.Node
}

func (bn *BaseNode) Type() Type {
	switch bn.tree.Type {
	case html.ElementNode:
		return ElementNodeType
	case html.TextNode:
		return TextNodeType
	case html.DocumentNode:
		return DocumentNodeType
	case html.CommentNode:
		return CommentNodeType
	case html.DoctypeNode:
		return DoctypeNodeType
	default:
		return RawNodeType
	}
}
func (bn *BaseNode) Inspect() string         { return bn.tree.Data }
func (bn *BaseNode) Tree() *html.Node        { return bn.tree }
func (bn *BaseNode) SetTree(tree *html.Node) { bn.tree = tree }
func (bn *BaseNode) Self() *html.Node        { return bn.tree }
func (bn *BaseNode) Parent() Node {
	if bn.tree.Parent != nil {
		return &BaseNode{bn.tree.Parent}
	}
	return nil
}
func (bn *BaseNode) FirstChild() Node {
	if bn.tree.FirstChild != nil {
		return &BaseNode{bn.tree.FirstChild}
	}
	return nil
}
func (bn *BaseNode) LastChild() Node {
	if bn.tree.LastChild != nil {
		return &BaseNode{bn.tree.LastChild}
	}
	return nil
}
func (bn *BaseNode) PrevSibling() Node {
	if bn.tree.PrevSibling != nil {
		return &BaseNode{bn.tree.PrevSibling}
	}
	return nil
}
func (bn *BaseNode) NextSibling() Node {
	if bn.tree.NextSibling != nil {
		return &BaseNode{bn.tree.NextSibling}
	}
	return nil
}
func (bn *BaseNode) Attr() []Node {
	if len(bn.tree.Attr) > 0 {
		var nodes []Node
		for _, a := range bn.tree.Attr {
			nodes = append(nodes, &AttrNode{bn.tree, a})
		}
		return nodes
	}
	return nil
}
func (bn *BaseNode) Text() string {
	if bn.Type() == CommentNodeType || bn.Type() == TextNodeType {
		return bn.Tree().Data
	}
	for c := bn.FirstChild(); c != nil; c = bn.NextSibling() {
		if c.Type() == TextNodeType {
			return c.Tree().Data
		}
	}
	return ""
}

// AttrNode ::= AttributeNode
// Attribute node is not exist in the golang.org/x/net/html package
// so the struct field is different from the BaseNode.
// AttrNode is basically, a child of ElementNode.
type AttrNode struct {
	parent *html.Node
	attr   html.Attribute
}

func (an *AttrNode) Type() Type                  { return AttributeNodeType }
func (an *AttrNode) Inspect() string             { return an.attr.Val }
func (an *AttrNode) Key() string                 { return an.attr.Key }
func (an *AttrNode) Attr() html.Attribute        { return an.attr }
func (an *AttrNode) SetAttr(attr html.Attribute) { an.attr = attr }
func (an *AttrNode) SetTree(p *html.Node)        { an.parent = p }
func (an *AttrNode) Tree() *html.Node            { return an.parent }
func (an *AttrNode) Self() *html.Node {
	var n *html.Node

	n.Type = html.NodeType(7)
	n.Data = an.Inspect()
	n.DataAtom = atom.Lookup([]byte(an.Inspect()))
	n.Attr = append(n.Attr, an.Attr())
	n.Parent = an.Tree()
	if an.PrevSibling() != nil {
		n.PrevSibling = an.PrevSibling().Self()
	}
	if an.NextSibling() != nil {
		n.NextSibling = an.NextSibling().Self()
	}

	return n
}
func (an *AttrNode) FirstChild() Node { return nil }
func (an *AttrNode) LastChild() Node  { return nil }
func (an *AttrNode) PrevSibling() Node {
	for i, a := range an.parent.Attr {
		if a.Key == an.Key() && a.Val == an.Inspect() {
			if i <= 0 {
				return nil
			}
			return &AttrNode{parent: an.parent, attr: an.parent.Attr[i-1]}
		}
	}
	return nil
}
func (an *AttrNode) NextSibling() Node {
	for i, a := range an.parent.Attr {
		if a.Key == an.Key() && a.Val == an.Inspect() {
			if i >= len(an.parent.Attr)-1 {
				return nil
			}
			return &AttrNode{parent: an.parent, attr: an.parent.Attr[i+1]}
		}
	}
	return nil
}
func (an *AttrNode) Parent() Node {
	if an.parent != nil {
		return &BaseNode{an.parent}
	}
	return nil
}
func (an *AttrNode) Text() string { return an.attr.Val }
