package object

// Context ...
type Context struct {
	store map[string]Item
	outer *Context
	CItem Item
}

// NewEnclosedContext ...
func NewEnclosedContext(outer *Context) *Context {
	ctx := NewContext()
	ctx.outer = outer
	return ctx
}

// NewContext ...
func NewContext() *Context {
	s := make(map[string]Item)
	return &Context{store: s, outer: nil}
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
