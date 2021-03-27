package object

// Context ...
type Context struct {
	store map[string]Item
	outer *Context
	Doc   Node
	CNode []Node
	CItem Item
	Focus
	Static
}

// Focus ...
type Focus struct {
	CSize int
	CPos  int
	CAxis string
}

// Static ...
type Static struct {
	BaseURI string
}

// NewContext ...
func NewContext() *Context {
	s := make(map[string]Item)
	return &Context{store: s, outer: nil}
}

// NewEnclosedContext ...
func NewEnclosedContext(outer *Context) *Context {
	ctx := NewContext()
	ctx.outer = outer
	ctx.Doc = outer.Doc
	ctx.CNode = outer.CNode
	ctx.CItem = outer.CItem
	ctx.CSize = outer.CSize
	ctx.CAxis = outer.CAxis
	ctx.CPos = outer.CPos
	ctx.BaseURI = outer.BaseURI
	return ctx
}

// Get ...
func (c *Context) Get(name string) (Item, bool) {
	item, ok := c.store[name]
	if !ok && c.outer != nil {
		item, ok = c.outer.Get(name)
	}
	return item, ok
}

// Set ...
func (c *Context) Set(name string, val Item) Item {
	c.store[name] = val
	return val
}
