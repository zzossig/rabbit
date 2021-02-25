package object

import (
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
	TypeANT     ast.QName
}

func (e *Element) node() {}

// Attribute ..
type Attribute struct {
	NodeName    string
	StringValue string
	Parent      *Node
	TypeANT     ast.QName
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
