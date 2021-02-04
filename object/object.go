package object

import (
	"fmt"
	"hash/fnv"
	"math"
	"strings"
)

// Item ..
type Item interface {
	Type() Type
	Inspect() string
}

// BIF represents Built-In Function
type BIF func(args ...Item) Item

// Type represents Item Type
type Type string

// Item Types
const (
	NilType   Type = "nil"
	ErrorType Type = "error"

	IntegerType Type = "int"
	DecimalType Type = "decimal"
	DoubleType  Type = "double"
	BooleanType Type = "bool"
	StringType  Type = "string"

	FuncType    Type = "func"
	BuiltinType Type = "bif"

	MapType   Type = "map"
	ArrayType Type = "array"

	XsAnyAtomicType      Type = "xs:anyAtomicType"
	XsUntypedAtomic      Type = "xs:untypedAtomic"
	XsDateTime           Type = "xs:dateTime"
	XsDateTimeStamp      Type = "xs:dateTimeStamp"
	XsDate               Type = "xs:date"
	XsTime               Type = "xs:time"
	XsDuration           Type = "xs:duration"
	XsYearMonthDuration  Type = "xs:yearMonthDuration"
	XsDayTimeDuration    Type = "xs:dayTimeDuration"
	XsFloat              Type = "xs:float"
	XsDouble             Type = "xs:double"
	XsDecimal            Type = "xs:decimal"
	XsInteger            Type = "xs:integer"
	XsNonPositiveInteger Type = "xs:nonPositiveInteger"
	XsNegativeInteger    Type = "xs:negativeInteger"
	XsLong               Type = "xs:long"
	XsInt                Type = "xs:int"
	XsShort              Type = "xs:short"
	XsByte               Type = "xs:byte"
	XsNonNegativeInteger Type = "xs:nonNegativeInteger"
	XsUnsignedLong       Type = "xs:unsignedLong"
	XsUnsignedInt        Type = "xs:unsignedInt"
	XsUnsignedShort      Type = "xs:unsignedShort"
	XsUnsignedByte       Type = "xs:unsignedByte"
	XsPositiveInteger    Type = "xs:positiveInteger"
	XsGYearMonth         Type = "xs:gYearMonth"
	XsGYear              Type = "xs:gYear"
	XsGMonthDay          Type = "xs:gMonthDay"
	XsGDay               Type = "xs:gDay"
	XsGMonth             Type = "xs:gMonth"
	XsString             Type = "xs:string"
	XsNormalizedString   Type = "xs:normalizedString"
	XsToken              Type = "xs:token"
	XsLanguage           Type = "xs:language"
	XsNMTOKEN            Type = "xs:NMTOKEN"
	XsName               Type = "xs:Name"
	XsNCName             Type = "xs:NCName"
	XsID                 Type = "xs:ID"
	XsIDREF              Type = "xs:IDREF"
	XsENTITY             Type = "xs:ENTITY"
	XsBoolean            Type = "xs:boolean"
	XsBase64Binary       Type = "xs:base64Binary"
	XsHexBinary          Type = "xs:hexBinary"
	XsAnyURI             Type = "xs:anyURI"
	XsQName              Type = "xs:QName"
	XsNOTATION           Type = "xs:NOTATION"
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

func (s *Sequence) Type() Type { return ArrayType }
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

// String ..
type String struct {
	Value string
}

func (s *String) Type() Type      { return StringType }
func (s *String) Inspect() string { return s.Value }

// HashKey ..
func (s *String) HashKey() HashKey {
	h := fnv.New64a()
	h.Write([]byte(s.Value))

	return HashKey{Type: s.Type(), Value: h.Sum64()}
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

	sb.WriteString("{")
	sb.WriteString(strings.Join(pairs, ", "))
	sb.WriteString("}")

	return sb.String()
}
