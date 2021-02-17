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

// NewEnclosedEnv ...
func NewEnclosedEnv(outer *Env) *Env {
	env := NewEnv()
	env.outer = outer
	return env
}

// NewEnv ...
func NewEnv() *Env {
	s := make(map[string]Item)
	ns := map[string]string{
		"xs":    "http://www.w3.org/2001/XMLSchema",
		"fn":    "http://www.w3.org/2005/xpath-functions",
		"map":   "http://www.w3.org/2005/xpath-functions/map",
		"array": "http://www.w3.org/2005/xpath-functions/array",
		"math":  "http://www.w3.org/2005/xpath-functions/math",
		"err":   "http://www.w3.org/2005/xqt-errors",
	}

	return &Env{store: s, outer: nil, NS: ns}
}

// Env ...
type Env struct {
	store map[string]Item
	outer *Env
	CItem Item
	Args  []Item
	NS    map[string]string
	r     io.Reader
	isXML bool
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

// NewReader ..
func (e *Env) NewReader(i interface{}, isXML bool) error {
	e.isXML = isXML

	switch i := i.(type) {
	case string:
		e.r = strings.NewReader(i)
		return nil
	case []byte:
		e.r = bytes.NewReader(i)
		return nil
	case *bytes.Buffer:
		e.r = i
		return nil
	}
	return fmt.Errorf("type now allowed: %T", i)
}

// NewReaderFile ..
func (e *Env) NewReaderFile(path string, isXML bool) error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}

	p := filepath.Join(dir, path)
	b, err := ioutil.ReadFile(p)
	if err != nil {
		return err
	}

	return e.NewReader(b, isXML)
}

// NewReaderHTTP ..
func (e *Env) NewReaderHTTP(addr string, isXML bool) error {
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

	return e.NewReader(buf, isXML)
}
