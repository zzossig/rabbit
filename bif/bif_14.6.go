package bif

import (
	"bufio"
	"net/http"
	"os"
	"path/filepath"

	"github.com/zzossig/rabbit/object"
	"golang.org/x/net/html"
)

func fnDoc(ctx *object.Context, args ...object.Item) object.Item {
	if len(args) > 1 {
		return NewError("too many parameters for function call: fn:doc")
	}
	if len(args) < 1 {
		return NewError("too few parameters for function call: fn:doc")
	}

	uri, ok := args[0].(*object.String)
	if !ok {
		return NewError("cannot match item type with required type")
	}

	docNode := &object.BaseNode{}

	if file, err := os.Open(uri.Value()); err == nil {
		defer file.Close()

		buf := bufio.NewReader(file)
		parsedHTML, err := html.Parse(buf)
		if err != nil {
			return NewError(err.Error())
		}
		parsedHTML.Type = html.DocumentNode

		docNode.SetTree(parsedHTML)
		ctx.Doc = docNode
		ctx.CNode = []object.Node{ctx.Doc}

		path, err := os.Getwd()
		if err != nil {
			return NewError(err.Error())
		}
		ctx.BaseURI = filepath.Join(path, uri.Value())
	}

	if resp, err := http.Get(uri.Value()); err == nil {
		defer resp.Body.Close()

		buf := bufio.NewReader(resp.Body)
		parsedHTML, err := html.Parse(buf)
		if err != nil {
			return NewError(err.Error())
		}
		parsedHTML.Type = html.DocumentNode

		docNode.SetTree(parsedHTML)
		ctx.Doc = docNode
		ctx.CNode = []object.Node{ctx.Doc}
		ctx.BaseURI = uri.Value()
	}
	return nil
}
