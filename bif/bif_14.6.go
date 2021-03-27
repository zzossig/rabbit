package bif

import (
	"bufio"
	"net/http"
	"os"
	"path/filepath"

	"github.com/zzossig/xpath/object"
	"golang.org/x/net/html"
)

func fnDoc(ctx *object.Context, args ...object.Item) object.Item {
	uri := args[0].(*object.String)
	docNode := &object.BaseNode{}
	var err error

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
		ctx.BaseURI = uri.Value()
	}

	if err != nil {
		return NewError("cannot retrieve resource %s", uri)
	}
	return nil
}
