package object

import (
	"math"

	"github.com/zzossig/xpath/bif"
)

// https://www.w3.org/TR/xpath-datamodel-31/

// Node ..
type Node interface {
	node()
}

// Atomic : t - type, v - value
type Atomic struct {
	t Type
	v Value
}

func (a *Atomic) Type() Type   { return a.t }
func (a *Atomic) Value() Value { return a.v }
func (a *Atomic) SetType(t Type) {
	a.t = t
}
func (a *Atomic) SetValue(v interface{}) Error {
	switch a.t {
	case ByteType:
		v, ok := v.(int)
		if !ok {
			return bif.NewError("Cannot convert %v to int", v)
		}
		if v > math.MaxInt8 {

		}
	case ShortType:
	case IntType:
	case LongType:
	case IntegerType:
	case DecimalType:
	case DoubleType:
	case BooleanType:
	case StringType:
	default:
		a.t = UntypedAtomicType
		a.v = v
	}
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
