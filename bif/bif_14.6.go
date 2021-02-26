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

	if file, ferr := os.Open(uri.Value()); ferr == nil {
		defer file.Close()

		buf := bufio.NewReader(file)
		root, err := html.Parse(buf)
		if err != nil {
			return NewError(err.Error())
		}
	}

	if resp, herr := http.Get(uri.Value()); herr == nil {
		defer resp.Body.Close()

		buf := bufio.NewReader(resp.Body)
		root, err := html.Parse(buf)
		if err != nil {
			return NewError(err.Error())
		}
	}

	return NewError("cannot retrieve resource %s", uri)
}

func buildNodeTree(doc, n *html.Node)
