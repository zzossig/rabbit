package bif

import (
	"bufio"
	"net/http"
	"os"

	"github.com/zzossig/xpath/object"
	"golang.org/x/net/html"
)

func doc(ctx *object.Context, args ...object.Item) object.Item {
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
	}

	if err != nil {
		return NewError("cannot retrieve resource %s", uri)
	}
	return nil
}
