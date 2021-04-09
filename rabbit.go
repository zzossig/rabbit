package rabbit

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/zzossig/rabbit/bif"
	"github.com/zzossig/rabbit/eval"
	"github.com/zzossig/rabbit/lexer"
	"github.com/zzossig/rabbit/object"
	"github.com/zzossig/rabbit/parser"
	"github.com/zzossig/rabbit/repl"
	"golang.org/x/net/html"
)

// XPath is a base object to evaluate xpath expressions.
// xpath field is xpath expression that is saved when using Eval method.
// context is a context that contains a document and context node.
// SetDoc function saves a document to the context field.
// evaled field is set when calling Eval method.
// object.Item is a custom data type used in rabbit language.
// You can convert object.Item to a golang data type using Data or Nodes method.
// errors field is collected errors while parsing and evaluating
type XPath struct {
	xpath   string
	context *object.Context
	evaled  object.Item
	errors  []error
}

// New creates new xpath object.
func New() *XPath {
	return &XPath{context: object.NewContext()}
}

// SetDoc set document to a context.
// if document is not set in a context, node related xpath expressions are not going to work.
// input param can be url or local filepath.
func (x *XPath) SetDoc(input string) *XPath {
	initContext(x.context)
	x.xpath = ""

	f := bif.F["fn:doc"]
	err := f(x.context, bif.NewString(input))
	if err != nil {
		x.errors = append(x.errors, fmt.Errorf(err.Inspect()))
	}
	return x
}

// SetDocR is another version of SetDoc.
func (x *XPath) SetDocR(r *http.Response) *XPath {
	initContext(x.context)
	x.xpath = ""
	defer r.Body.Close()

	nr := bufio.NewReader(r.Body)
	parsedHTML, err := html.Parse(nr)
	if err != nil {
		x.errors = append(x.errors, err)
	}
	parsedHTML.Type = html.DocumentNode

	docNode := &object.BaseNode{}
	docNode.SetTree(parsedHTML)
	x.context.Doc = docNode
	x.context.CNode = []object.Node{x.context.Doc}

	return x
}

// SetDocN is another version of SetDoc.
func (x *XPath) SetDocN(n *html.Node) *XPath {
	initContext(x.context)
	x.xpath = ""

	docNode := &object.BaseNode{}
	docNode.SetTree(n)
	x.context.Doc = docNode
	x.context.CNode = []object.Node{x.context.Doc}

	return x
}

// SetDocS is another version of SetDoc.
func (x *XPath) SetDocS(s string) *XPath {
	initContext(x.context)
	x.xpath = ""

	nr := strings.NewReader(s)
	parsedHTML, err := html.Parse(nr)
	if err != nil {
		x.errors = append(x.errors, err)
	}
	parsedHTML.Type = html.DocumentNode

	docNode := &object.BaseNode{}
	docNode.SetTree(parsedHTML)
	x.context.Doc = docNode
	x.context.CNode = []object.Node{x.context.Doc}

	return x
}

// Eval evaluates a xpath expression and save the result to evaled field.
func (x *XPath) Eval(input string) *XPath {
	if len(x.errors) > 0 {
		return x
	}

	if x.xpath != "" && input != "" && input[0] != '/' {
		x.xpath += "/" + input
	} else {
		x.xpath += input
	}

	l := lexer.New(input)
	p := parser.New(l)
	px := p.ParseXPath()

	if len(p.Errors()) != 0 {
		x.errors = append(x.errors, p.Errors()...)
		return x
	}

	e := eval.Eval(px, x.context)
	if bif.IsError(e) {
		x.errors = append(x.errors, fmt.Errorf(e.Inspect()))
		return x
	}
	x.evaled = e

	return x
}

// Evals evaluates a xpath expression and returns slice of *XPath.
func (x *XPath) Evals(input string) []*XPath {
	if len(x.errors) > 0 {
		return []*XPath{x}
	}

	if x.xpath != "" && input != "" && input[0] != '/' {
		x.xpath += "/" + input
	} else {
		x.xpath += input
	}

	l := lexer.New(input)
	p := parser.New(l)
	px := p.ParseXPath()

	if len(p.Errors()) != 0 {
		x.errors = append(x.errors, p.Errors()...)
		return []*XPath{x}
	}

	e := eval.Eval(px, x.context)
	if bif.IsError(e) {
		x.errors = append(x.errors, fmt.Errorf(e.Inspect()))
		return []*XPath{x}
	}

	result := []*XPath{}
	seq := e.(*object.Sequence)

	for _, item := range seq.Items {
		switch item := item.(type) {
		case *object.Integer:
			newX := &XPath{xpath: x.xpath, evaled: item, context: copyContext(x.context)}
			result = append(result, newX)
		case *object.Decimal:
			newX := &XPath{xpath: x.xpath, evaled: item, context: copyContext(x.context)}
			result = append(result, newX)
		case *object.Double:
			newX := &XPath{xpath: x.xpath, evaled: item, context: copyContext(x.context)}
			result = append(result, newX)
		case *object.Boolean:
			newX := &XPath{xpath: x.xpath, evaled: item, context: copyContext(x.context)}
			result = append(result, newX)
		case *object.String:
			newX := &XPath{xpath: x.xpath, evaled: item, context: copyContext(x.context)}
			result = append(result, newX)
		case *object.Map:
			newX := &XPath{xpath: x.xpath, evaled: item, context: copyContext(x.context)}
			result = append(result, newX)
		case *object.Array:
			newX := &XPath{xpath: x.xpath, evaled: item, context: copyContext(x.context)}
			result = append(result, newX)
		case *object.BaseNode:
			newX := &XPath{xpath: x.xpath, evaled: item, context: copyContextN(x.context, item)}
			result = append(result, newX)
		case *object.AttrNode:
			newX := &XPath{xpath: x.xpath, evaled: item, context: copyContextN(x.context, item)}
			result = append(result, newX)
		}
	}

	return result
}

func (x *XPath) Get() string {
	items := x.GetAll()
	if items == nil {
		return ""
	}
	return items[0]
}

func (x *XPath) GetAll() []string {
	initContext(x.context)
	x.xpath = ""

	if x.evaled == nil {
		x.errors = append(x.errors, fmt.Errorf("cannot convert item since evaled field is nil"))
		return nil
	}

	return convertString(x.evaled)
}

// Data selects first item of returned value from DataAll
func (x *XPath) Data() interface{} {
	items := x.DataAll()
	if len(items) > 0 {
		return items[0]
	}
	return nil
}

// DataAll convert evaled field to []interface{}
func (x *XPath) DataAll() []interface{} {
	initContext(x.context)
	x.xpath = ""

	if x.evaled == nil {
		x.errors = append(x.errors, fmt.Errorf("cannot convert item since evaled field is nil"))
		return nil
	}

	e, err := convert(x.evaled)
	if err != nil {
		x.errors = append(x.errors, err)
		return nil
	}

	return e.([]interface{})
}

// Node selects first item of returned value from NodeAll
func (x *XPath) Node() *html.Node {
	nodes := x.NodeAll()
	if len(nodes) > 0 {
		return nodes[0]
	}
	return nil
}

// NodeAll convert evaled field to []*html.Node
func (x *XPath) NodeAll() []*html.Node {
	initContext(x.context)
	x.xpath = ""

	if x.evaled == nil {
		x.errors = append(x.errors, fmt.Errorf("cannot convert item since evaled field is nil"))
		return nil
	}

	e, err := convertNode(x.evaled)
	if err != nil {
		x.errors = append(x.errors, err)
		return nil
	}

	return e
}

// Raw returns evaled field
func (x *XPath) Raw() object.Item {
	return x.evaled
}

// Errors returns errors field
func (x *XPath) Errors() []error {
	return x.errors
}

// String returns input field
func (x *XPath) String() string {
	return x.xpath
}

// CLI is a command line interface
func (x *XPath) CLI() {
	repl.Start(os.Stdin, os.Stdout, x.context)
}
