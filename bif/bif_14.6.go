package bif

import (
	"bufio"
	"net/http"
	"os"

	"github.com/zzossig/xpath/object"
	"golang.org/x/net/html"
)

func doc(args ...object.Item) object.Item {
	uri := args[0].(*object.String)
	docNode := &object.Node{}

	if file, ferr := os.Open(uri.Value()); ferr == nil {
		defer file.Close()

		buf := bufio.NewReader(file)
		parsedHTML, err := html.Parse(buf)
		if err != nil {
			return NewError(err.Error())
		}

		docNode.SetTree(parsedHTML)
		return docNode
	}

	if resp, herr := http.Get(uri.Value()); herr == nil {
		defer resp.Body.Close()

		buf := bufio.NewReader(resp.Body)
		parsedHTML, err := html.Parse(buf)
		if err != nil {
			return NewError(err.Error())
		}

		docNode.SetTree(parsedHTML)
		return docNode
	}

	return NewError("cannot retrieve resource %s", uri)
}
