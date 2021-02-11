package object

// NewEnclosedEnv ...
func NewEnclosedEnv(outer *Env) *Env {
	env := NewEnv()
	env.outer = outer
	return env
}

// NewEnv ...
func NewEnv() *Env {
	s := make(map[string]Item)
	return &Env{store: s, outer: nil}
}

// Env ...
type Env struct {
	store map[string]Item
	outer *Env
	Args  []Item
	CItem Item
}

// Get ...
func (e *Env) Get(name string) (Item, bool) {
	item, ok := e.store[name]
	if !ok && e.outer != nil {
		item, ok = e.outer.Get(name)
	}
	return item, ok
}

// Set ...
func (e *Env) Set(name string, val Item) Item {
	e.store[name] = val
	return val
}
