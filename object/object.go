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

// Type ::= ErrorType
func (e *Error) Type() Type { return ErrorType }

// Inspect ::= "Error: " + e.Message
func (e *Error) Inspect() string { return "ERROR: " + e.Message }

// Placeholder is an item that is represents ?(question token) when doing evaluation
type Placeholder struct{}

// Type ::= PholderType
func (p *Placeholder) Type() Type { return PholderType }

// Inspect ::= ?
func (p *Placeholder) Inspect() string { return "?" }

// Varref is an item that is represents $var when doing evaluation
type Varref struct {
	Name ast.EQName
}

// Type ::= VarrefType
func (v *Varref) Type() Type { return VarrefType }

// Inspect ::= EQName.Value()
func (v *Varref) Inspect() string { return fmt.Sprintf("$%s", v.Name.Value()) }

// Sequence is an ordered collection of zero or more items.
type Sequence struct {
	Items []Item
}

// Type ::= SequenceType
func (s *Sequence) Type() Type { return SequenceType }

// Inspect ::= (...)
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

// Type ::= MapType
func (m *Map) Type() Type { return MapType }

// Inspect ::= map{...}
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

// Type ::= ArrayType
func (a *Array) Type() Type { return ArrayType }

// Inspect ::= [...]
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

// Type ::= FuncType
func (fi *FuncInline) Type() Type { return FuncType }

// Inspect ::= function
func (fi *FuncInline) Inspect() string { return "function" }

// FuncPartial ::= ns:bif(?,...)
type FuncPartial struct {
	Name ast.EQName
	Args []Item
	PCnt int
	*Context
	*Func
}

// Type ::= FuncType
func (fp *FuncPartial) Type() Type      { return FuncType }
func (fp *FuncPartial) Inspect() string { return "function" }

// Integer is an item that is represents int data-type
type Integer struct {
	value int
}

// Type ::= IntegerType
func (i *Integer) Type() Type { return IntegerType }

// Inspect ::= %d
func (i *Integer) Inspect() string { return fmt.Sprintf("%d", i.value) }

// SetValue is setter for the Integer
func (i *Integer) SetValue(v int) { i.value = v }

// Value is getter for the Integer
func (i *Integer) Value() int { return i.value }

// HashKey used as a map key
func (i *Integer) HashKey() HashKey {
	return HashKey{Type: i.Type(), Value: uint64(i.value)}
}

// Decimal is an item that is represents float64 data-type
// number token that is not contains `e` or `E` is evaluated to Decimal
type Decimal struct {
	value float64
}

// Type ::= DecimalType
func (d *Decimal) Type() Type { return DecimalType }

// Inspect ::= %f
func (d *Decimal) Inspect() string { return fmt.Sprintf("%f", d.value) }

// SetValue is setter for the Decimal
func (d *Decimal) SetValue(v float64) { d.value = v }

// Value is getter for the Decimal
func (d *Decimal) Value() float64 { return d.value }

// HashKey used as a map key
func (d *Decimal) HashKey() HashKey {
	return HashKey{Type: d.Type(), Value: math.Float64bits(d.value)}
}

// Double is an item that is represents float64 data-type
// number token that contains `e` or `E` is evaluated to Double
type Double struct {
	value float64
}

// Type ::= DoubleType
func (d *Double) Type() Type { return DoubleType }

// Inspect ::= %e
func (d *Double) Inspect() string { return fmt.Sprintf("%e", d.value) }

// SetValue is setter for the Double
func (d *Double) SetValue(v float64) { d.value = v }

// Value is getter for the Double
func (d *Double) Value() float64 { return d.value }

// HashKey used as a map key
func (d *Double) HashKey() HashKey {
	return HashKey{Type: d.Type(), Value: math.Float64bits(d.value)}
}

// Boolean is an item that is represents bool data-type
type Boolean struct {
	value bool
}

// Type ::= BooleanType
func (b *Boolean) Type() Type { return BooleanType }

// Inspect ::= %t
func (b *Boolean) Inspect() string { return fmt.Sprintf("%t", b.value) }

// SetValue is setter for the Boolean
func (b *Boolean) SetValue(v bool) { b.value = v }

// Value is getter for the Boolean
func (b *Boolean) Value() bool { return b.value }

// HashKey used as a map key
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

// Type ::= StringType
func (s *String) Type() Type { return StringType }

// Inspect ::= string
func (s *String) Inspect() string { return s.value }

// SetValue is setter for the String
func (s *String) SetValue(v string) { s.value = v }

// Value is getter for the String
func (s *String) Value() string { return s.value }

// HashKey used as a map key
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

// Type ::= ElementNodeType | TextNodeType | DocumentNodeType | CommentNodeType | DoctypeNodeType | RawNodeType
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

// Inspect ::= *html.Node.Data
func (bn *BaseNode) Inspect() string {
	switch bn.tree.Type {
	case html.ElementNode:
		return fmt.Sprintf("Elem{%s}", bn.tree.Data)
	case html.TextNode:
		return fmt.Sprintf("Text{%s}", bn.tree.Data)
	case html.DocumentNode:
		return fmt.Sprintf("Doc{%s}", bn.tree.Data)
	case html.CommentNode:
		return fmt.Sprintf("Comm{%s}", bn.tree.Data)
	case html.DoctypeNode:
		return fmt.Sprintf("Doctype{%s}", bn.tree.Data)
	}
	return bn.tree.Data
}

// Tree returns *html.Node. BaseNode is just a wrapper type for the *html.Node
func (bn *BaseNode) Tree() *html.Node { return bn.tree }

// SetTree is setter for the BaseNode
func (bn *BaseNode) SetTree(tree *html.Node) { bn.tree = tree }

// Self is getter for the BaseNode
func (bn *BaseNode) Self() *html.Node { return bn.tree }

// Parent returns parent node of the current one if exist
func (bn *BaseNode) Parent() Node {
	if bn.tree.Parent != nil {
		return &BaseNode{bn.tree.Parent}
	}
	return nil
}

// FirstChild returns first child node of the current one if exist
func (bn *BaseNode) FirstChild() Node {
	if bn.tree.FirstChild != nil {
		return &BaseNode{bn.tree.FirstChild}
	}
	return nil
}

// LastChild returns last child node of the current one if exist
func (bn *BaseNode) LastChild() Node {
	if bn.tree.LastChild != nil {
		return &BaseNode{bn.tree.LastChild}
	}
	return nil
}

// PrevSibling returns previous sibling node of the current one if exist
func (bn *BaseNode) PrevSibling() Node {
	if bn.tree.PrevSibling != nil {
		return &BaseNode{bn.tree.PrevSibling}
	}
	return nil
}

// NextSibling returns next sibling node of the current one if exist
func (bn *BaseNode) NextSibling() Node {
	if bn.tree.NextSibling != nil {
		return &BaseNode{bn.tree.NextSibling}
	}
	return nil
}

// Attr returns Attr field of element node with wrap it to AttrNode
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

// Text returns *html.Node.Data
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

// Type ::= AttributeNodeType
func (an *AttrNode) Type() Type { return AttributeNodeType }

// Inspect returns a value of the attr field
func (an *AttrNode) Inspect() string { return fmt.Sprintf("Attr{%s:%s}", an.attr.Key, an.attr.Val) }

// Key returns a key of the attr field
func (an *AttrNode) Key() string { return an.attr.Key }

// Attr returns attr field
func (an *AttrNode) Attr() html.Attribute { return an.attr }

// SetAttr is setter for the attr field
func (an *AttrNode) SetAttr(attr html.Attribute) { an.attr = attr }

// SetTree is setter for the parent field
func (an *AttrNode) SetTree(p *html.Node) { an.parent = p }

// Tree is getter for the parent field
func (an *AttrNode) Tree() *html.Node { return an.parent }

// Self returns customized *html.Node which represents Attribute node
func (an *AttrNode) Self() *html.Node {
	var n html.Node

	n.Type = html.NodeType(7)
	n.Data = an.Key()
	n.DataAtom = atom.Lookup([]byte(an.Key()))
	n.Attr = append(n.Attr, an.attr)
	n.Parent = an.Tree()
	if an.PrevSibling() != nil {
		n.PrevSibling = an.PrevSibling().Self()
	}
	if an.NextSibling() != nil {
		n.NextSibling = an.NextSibling().Self()
	}

	return &n
}

// FirstChild is not exist in AttrNode
func (an *AttrNode) FirstChild() Node { return nil }

// LastChild is not exist in AttrNode
func (an *AttrNode) LastChild() Node { return nil }

// PrevSibling returns previous sibling of the current one if exist
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

// NextSibling returns next sibling of the current one if exist
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

// Parent returns parent node of the current one if exist
func (an *AttrNode) Parent() Node {
	if an.parent != nil {
		return &BaseNode{an.parent}
	}
	return nil
}
func (an *AttrNode) Text() string { return an.attr.Val }
