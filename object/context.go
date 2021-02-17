package object

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// NewEnclosedContext ...
func NewEnclosedContext(outer *Context) *Context {
	ctx := NewContext()
	ctx.outer = outer
	return ctx
}

// NewContext ...
func NewContext() *Context {
	s := make(map[string]Item)
	ns := map[string]string{
		"xs":    "http://www.w3.org/2001/XMLSchema",
		"fn":    "http://www.w3.org/2005/xpath-functions",
		"map":   "http://www.w3.org/2005/xpath-functions/map",
		"array": "http://www.w3.org/2005/xpath-functions/array",
		"math":  "http://www.w3.org/2005/xpath-functions/math",
		"err":   "http://www.w3.org/2005/xqt-errors",
	}

	return &Context{store: s, outer: nil, NS: ns}
}

// Context ...
type Context struct {
	store map[string]Item
	outer *Context
	CItem Item
	Args  []Item
	NS    map[string]string
	r     io.Reader
	isXML bool
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

// NewReader ..
func (c *Context) NewReader(i interface{}, isXML bool) error {
	c.isXML = isXML

	switch i := i.(type) {
	case string:
		c.r = strings.NewReader(i)
		return nil
	case []byte:
		c.r = bytes.NewReader(i)
		return nil
	case *bytes.Buffer:
		c.r = i
		return nil
	}
	return fmt.Errorf("type now allowed: %T", i)
}

// NewReaderFile ..
func (c *Context) NewReaderFile(path string, isXML bool) error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}

	p := filepath.Join(dir, path)
	b, err := ioutil.ReadFile(p)
	if err != nil {
		return err
	}

	return c.NewReader(b, isXML)
}

// NewReaderHTTP ..
func (c *Context) NewReaderHTTP(addr string, isXML bool) error {
	resp, err := http.Get(addr)
	if err != nil {
		return err
	}

	buf := bytes.NewBuffer([]byte(""))
	_, err = io.Copy(buf, resp.Body)
	if err != nil {
		return err
	}
	resp.Body.Close()

	return c.NewReader(buf, isXML)
}
