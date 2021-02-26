package object

import (
	"fmt"
	"hash/fnv"
	"math"
	"strings"

	"github.com/zzossig/xpath/ast"
)

// Item ..
type Item interface {
	Type() Type
	Inspect() string
}

// Node ..
type Node interface {
	Item
	node()
}

// predefined
var (
	NIL   = &Nil{}
	TRUE  = &Boolean{true}
	FALSE = &Boolean{false}
)

// Nil ..
type Nil struct{}

func (n *Nil) Type() Type      { return NilType }
func (n *Nil) Inspect() string { return "nil" }

// Error ..
type Error struct {
	Message string
}

func (e *Error) Type() Type      { return ErrorType }
func (e *Error) Inspect() string { return "ERROR: " + e.Message }

// Placeholder ..
type Placeholder struct{}

func (p *Placeholder) Type() Type      { return PholderType }
func (p *Placeholder) Inspect() string { return "?" }

// Varref ..
type Varref struct {
	Name ast.EQName
}

func (v *Varref) Type() Type      { return VarrefType }
func (v *Varref) Inspect() string { return fmt.Sprintf("$%s", v.Name.Value()) }

// Sequence ..
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

// Hasher ..
type Hasher interface {
	HashKey() HashKey
}

// HashKey ..
type HashKey struct {
	Type
	Value uint64
}

// Pair ..
type Pair struct {
	Key   Item
	Value Item
}

// Map ..
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

// Array ..
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

// FuncNamed ..
type FuncNamed struct {
	Name ast.EQName
	Num  int
	*Context
}

func (fn *FuncNamed) Type() Type      { return FuncType }
func (fn *FuncNamed) Inspect() string { return fmt.Sprintf("%s#%d", fn.Name.Value(), fn.Num) }

// FuncInline ..
type FuncInline struct {
	PL   *ast.ParamList
	Body *ast.EnclosedExpr
	*Context
}

func (fi *FuncInline) Type() Type      { return FuncType }
func (fi *FuncInline) Inspect() string { return "function" }

// FuncPartial ..
type FuncPartial struct {
	Name   ast.EQName
	Args   []Item
	PNames []string
	PCnt   int
	*Context
	*Func
}

func (fp *FuncPartial) Type() Type      { return FuncType }
func (fp *FuncPartial) Inspect() string { return "function" }

// Integer ..
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

// Decimal ..
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

// Double ..
type Double struct {
	value float64
}

func (d *Double) Type() Type         { return DecimalType }
func (d *Double) Inspect() string    { return fmt.Sprintf("%e", d.value) }
func (d *Double) SetValue(v float64) { d.value = v }
func (d *Double) Value() float64     { return d.value }
func (d *Double) HashKey() HashKey {
	return HashKey{Type: d.Type(), Value: math.Float64bits(d.value)}
}

// Boolean ..
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

// String ..
type String struct {
	value string
}

func (s *String) Type() Type        { return StringType }
func (s *String) Inspect() string   { return fmt.Sprintf("%s", s.value) }
func (s *String) SetValue(v string) { s.value = v }
func (s *String) Value() string     { return s.value }
func (s *String) HashKey() HashKey {
	h := fnv.New64a()
	h.Write([]byte(s.value))

	return HashKey{Type: s.Type(), Value: h.Sum64()}
}

type DocNode struct {
	Parent, FirstChild, LastChild *Node

	NodeName ast.EQName
	Children []*Node
	Attr     []*Node
}

func (dn *DocNode) node()           {}
func (dn *DocNode) Type() Type      { return NodeType }
func (dn *DocNode) Inspect() string { return dn.NodeName.Value() }

type ElemNode struct {
	Parent, FirstChild, LastChild, PrevSibling, NextSibling *Node

	NodeName ast.EQName
	Children []*Node
	Attr     []*AttrNode
}

func (en *ElemNode) node()           {}
func (en *ElemNode) Type() Type      { return NodeType }
func (en *ElemNode) Inspect() string { return en.NodeName.Value() }

type AttrNode struct {
	Parent, PrevSibling, NextSibling *Node

	NodeName ast.EQName
	Data     string
}

func (an *AttrNode) node()           {}
func (an *AttrNode) Type() Type      { return NodeType }
func (an *AttrNode) Inspect() string { return an.Data }

type TextNode struct {
	Parent, PrevSibling, NextSibling *Node

	Content string
}

func (tn *TextNode) node()           {}
func (tn *TextNode) Type() Type      { return NodeType }
func (tn *TextNode) Inspect() string { return tn.Content }

type CommentNode struct {
	Parent, PrevSibling, NextSibling *Node

	Content string
}

func (cn *CommentNode) node()           {}
func (cn *CommentNode) Type() Type      { return NodeType }
func (cn *CommentNode) Inspect() string { return cn.Content }
