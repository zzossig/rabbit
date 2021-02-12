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

// Func represents function type
type Func func(args ...Item) Item

// Type represents Item Type
type Type string

// Item Types
const (
	NilType     Type = "nil"
	ErrorType   Type = "error"
	PholderType Type = "?"

	IntegerType  Type = "int"
	DecimalType  Type = "decimal"
	DoubleType   Type = "double"
	BooleanType  Type = "bool"
	StringType   Type = "string"
	MapType      Type = "map"
	ArrayType    Type = "array"
	SequenceType Type = "sequence"

	FuncCallType   Type = "functionC"
	FuncNamedType  Type = "functionN"
	FuncInlineType Type = "functionI"

	NodeType Type = "node"
)

// Hasher ..
type Hasher interface {
	HashKey() HashKey
}

// HashKey ..
type HashKey struct {
	Type
	Value uint64
}

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

// Integer ..
type Integer struct {
	Value int
}

func (i *Integer) Type() Type      { return IntegerType }
func (i *Integer) Inspect() string { return fmt.Sprintf("%d", i.Value) }
func (i *Integer) HashKey() HashKey {
	return HashKey{Type: i.Type(), Value: uint64(i.Value)}
}

// Decimal ..
type Decimal struct {
	Value float64
}

func (d *Decimal) Type() Type      { return DecimalType }
func (d *Decimal) Inspect() string { return fmt.Sprintf("%f", d.Value) }
func (d *Decimal) HashKey() HashKey {
	return HashKey{Type: d.Type(), Value: math.Float64bits(d.Value)}
}

// Double ..
type Double struct {
	Value float64
}

func (d *Double) Type() Type      { return DoubleType }
func (d *Double) Inspect() string { return fmt.Sprintf("%e", d.Value) }
func (d *Double) HashKey() HashKey {
	return HashKey{Type: d.Type(), Value: math.Float64bits(d.Value)}
}

// Boolean ..
type Boolean struct {
	Value bool
}

func (b *Boolean) Type() Type      { return BooleanType }
func (b *Boolean) Inspect() string { return fmt.Sprintf("%t", b.Value) }
func (b *Boolean) HashKey() HashKey {
	var value uint64

	if b.Value {
		value = 1
	} else {
		value = 0
	}

	return HashKey{Type: b.Type(), Value: value}
}

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
type Placeholder struct {
	Value interface{}
}

func (p *Placeholder) Type() Type { return PholderType }
func (p *Placeholder) Inspect() string {
	switch v := p.Value.(type) {
	case *Integer:
		return v.Inspect()
	case *Decimal:
		return v.Inspect()
	case *Double:
		return v.Inspect()
	case *String:
		return v.Inspect()
	case *Error:
		return v.Inspect()
	default:
		return fmt.Sprintf("%s", p.Value)
	}
}

// String ..
type String struct {
	Value string
}

func (s *String) Type() Type      { return StringType }
func (s *String) Inspect() string { return fmt.Sprintf("%q", s.Value) }
func (s *String) HashKey() HashKey {
	h := fnv.New64a()
	h.Write([]byte(s.Value))

	return HashKey{Type: s.Type(), Value: h.Sum64()}
}

// FuncNamed ..
type FuncNamed struct {
	Name string
	Num  int
	*Env
}

func (fn *FuncNamed) Type() Type      { return FuncNamedType }
func (fn *FuncNamed) Inspect() string { return fmt.Sprintf("functionN: %s", fn.Name) }

// FuncInline ..
type FuncInline struct {
	Body  *ast.EnclosedExpr
	PL    *ast.ParamList
	SType *ast.SequenceType
}

func (fi *FuncInline) Type() Type      { return FuncInlineType }
func (fi *FuncInline) Inspect() string { return "functionI" }

// FuncCall ..
type FuncCall struct {
	Name string
	*Func
	*Env
}

func (fc *FuncCall) Type() Type      { return FuncCallType }
func (fc *FuncCall) Inspect() string { return "functionC" }

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
