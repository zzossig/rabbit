package object

// Context contains items that is used in Eval or built-in functions
// store field stores Varref as a key, Item as as value
// In example expression, let $a := 1 return $a, 'a' is a key and 1 is a value
// outer is a context to make inner scope of function.
// In example expression, let $a := 1 return function($a) {$a}
// returned function makes new context to keep the varref $a
type Context struct {
	store map[string]Item
	outer *Context
	Doc   Node
	CNode []Node
	CItem Item
	Focus
	Static
}

// Focus contains context size, context position, context axis
// CSize is used in bif - fn:last()
// CPos is used in predicate expressions
// CAxis is used to evaluate (relative)path expressions
type Focus struct {
	CSize int
	CPos  int
	CAxis string
}

// Static contains information that is available during static analysis of the expression, prior to its evaluation
type Static struct {
	BaseURI string
}

// NewContext creates a new context
func NewContext() *Context {
	s := make(map[string]Item)
	return &Context{store: s, outer: nil}
}

// NewEnclosedContext creates a new context and stores the current context in the outer field
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

// Get retrieve items from the store field or the outer field.
func (c *Context) Get(name string) (Item, bool) {
	item, ok := c.store[name]
	if !ok && c.outer != nil {
		item, ok = c.outer.Get(name)
	}
	return item, ok
}

// Set save item in the current context
func (c *Context) Set(name string, val Item) Item {
	c.store[name] = val
	return val
}
