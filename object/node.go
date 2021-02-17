package object

import (
	"fmt"
	"math"

	"github.com/zzossig/xpath/ast"
)

// https://www.w3.org/TR/xpath-datamodel-31/

// Node ..
type Node interface {
	Item
	BaseURI() string
	DocURI() string
	Attributes() []*Attribute
	Children() []*Node
	NSNodes() []*Namespace
	NodeKind() string
}

// Document ..
type Document struct {
	BaseURI     string
	DocURI      string
	StringValue string
	Children    []*Node
}

func (d *Document) node() {}

// Element ..
type Element struct {
	BaseURI     string
	NodeName    string
	StringValue string
	Attr        []Attribute
	Parent      *Node
	Children    []*Node
}

func (e *Element) node() {}

// Attribute ..
type Attribute struct {
	NodeName    string
	StringValue string
	Parent      *Node
}

func (a *Attribute) node() {}

// Text ..
type Text struct {
	Content []byte
	Parent  *Node
}

func (t *Text) node() {}

// Namespace ..
type Namespace struct {
	Prefix string
	URI    string
	Parent *Node
}

func (n *Namespace) node() {}

// PI ..
type PI struct {
	BaseURI string
	Target  string
	Content []byte
	Parent  *Node
}

func (pi *PI) node() {}

// Comment ..
type Comment struct {
	Content []byte
	Parent  *Node
}

func (c *Comment) node() {}

// QName : prefix, local
type QName struct {
	prefix ast.NCName
	local  ast.NCName
}

func (qn *QName) Prefix() string {
	return qn.prefix.Value()
}

func (qn *QName) Local() string {
	return qn.local.Value()
}

func (qn *QName) NamespaceURI(c *Context) string {
	return ""
}

// Atomic : t - type, v - value
type Atomic struct {
	t Type
	v Value
}

func (a *Atomic) Type() Type   { return a.t }
func (a *Atomic) Value() Value { return a.v }
func (a *Atomic) SetValue(v Value, t Type) error {
	switch t {
	case ByteType:
		v, ok := v.(int)
		if !ok {
			return fmt.Errorf("cannot convert %v to int", v)
		}
		if v > math.MaxInt8 || v < math.MinInt8 {
			return fmt.Errorf("max=%d, min=%d, got=%d", math.MaxInt8, math.MinInt8, v)
		}

		a.t = t
		a.v = v
		return nil
	case ShortType:
		v, ok := v.(int)
		if !ok {
			return fmt.Errorf("cannot convert %v to int", v)
		}
		if v > math.MaxInt16 || v < math.MinInt16 {
			return fmt.Errorf("max=%d, min=%d, got=%d", math.MaxInt16, math.MinInt16, v)
		}

		a.t = t
		a.v = v
		return nil
	case IntType:
		v, ok := v.(int)
		if !ok {
			return fmt.Errorf("cannot convert %v to int", v)
		}
		if v > math.MaxInt32 || v < math.MinInt32 {
			return fmt.Errorf("max=%d, min=%d, got=%d", math.MaxInt32, math.MinInt32, v)
		}

		a.t = t
		a.v = v
		return nil
	case LongType:
		v, ok := v.(int)
		if !ok {
			return fmt.Errorf("cannot convert %v to int", v)
		}
		if v > math.MaxInt64 || v < math.MinInt64 {
			return fmt.Errorf("max=%d, min=%d, got=%d", math.MaxInt64, math.MinInt64, v)
		}

		a.t = t
		a.v = v
		return nil
	case IntegerType:
		v, ok := v.(int)
		if !ok {
			return fmt.Errorf("cannot convert %v to int", v)
		}
		if v > math.MaxInt64 || v < math.MinInt64 {
			return fmt.Errorf("max=%d, min=%d, got=%d", math.MaxInt64, math.MinInt64, v)
		}

		a.t = t
		a.v = v
		return nil
	case DecimalType:
		v, ok := v.(float64)
		if !ok {
			return fmt.Errorf("cannot convert %v to float64", v)
		}
		if v > math.MaxFloat64 || v < -math.MaxFloat64 {
			return fmt.Errorf("max=%f, min=%f, got=%f", math.MaxFloat64, -math.MaxFloat64, v)
		}

		a.t = t
		a.v = v
		return nil
	case DoubleType:
		v, ok := v.(float64)
		if !ok {
			return fmt.Errorf("cannot convert %v to float64", v)
		}
		if v > math.MaxFloat64 || v < -math.MaxFloat64 {
			return fmt.Errorf("max=%f, min=%f, got=%f", math.MaxFloat64, -math.MaxFloat64, v)
		}

		a.t = t
		a.v = v
		return nil
	case BooleanType:
		v, ok := v.(bool)
		if !ok {
			return fmt.Errorf("cannot convert %v to boolean", v)
		}

		a.t = t
		a.v = v
		return nil
	case StringType:
		a.v = v
		return nil
	case UntypedAtomicType:
		a.t = t
		a.v = v
		return nil
	}
	return fmt.Errorf("cannot set value of type %s", t)
}
