# 🐰rabbit

[![Go Report Card](https://goreportcard.com/badge/github.com/zzossig/rabbit)](https://goreportcard.com/report/github.com/zzossig/rabbit)
> An interpreted language written in Go - XPath 3.1 implementation for HTML

XML Path Language(XPath) 3.1 is W3C recommendation since 21 march 2017.
The rabbit language is built for selecting HTML nodes with XPath syntax.

## Overview

Rabbit language is built for HTML, not for XML. Since XPath 3.1 is targeted for XML, it was not possible to implement all the concepts listed in [https://www.w3.org/TR/xpath-31/](https://www.w3.org/TR/xpath-31/). But in most cases, it is fair enough for selecting HTML nodes with rabbit language.

For example)

- `//a`
- `//div[@category='web']/preceding::node()[2]`
- `let $abc := ('a', 'b', 'c') return fn:insert-before($abc, 4, 'z')`

## Basic Usage

```go
// you can chaining xpath object. data is nil or []interface{}
data := rabbit.New().SetDoc("uri/or/filepath.txt").Eval("//a").Data()
```

```go
// if you expect evaled result is a sequence of html node, 
// use Nodes() instead of Data()
nodes := rabbit.New().SetDoc("uri/or/filepath.txt").Eval("//a").Nodes()
```

```go
// with error check
x := rabbit.New()
x.SetDoc("uri/or/filepath.txt")
if len(x.errors) > 0 {
  // ... do something with errors (the errors type is []error)
}
data := x.Eval("//a").Data()
```

```go
// without SetDoc. Since document is not set in the context, 
// node related xpath expressions are not going to work.
x := rabbit.New()
data := x.Eval("1+1").Data()
```

```go
// you can test simple xpath expressions using cli program
rabbit.New().SetDoc("uri/or/filepath.txt").CLI()
```

## Features

### What is supported

1. Primary Expressions
    - Integer(1)
    - Decimal(1.1)
    - Double(1e1)
    - String("")
    - Boolean(true, false)
    - Variable($var)
    - Context Item(.)
    - Placeholder(?)
2. Functions
    - Named Function(built in function - bif)
    - Inline Function(custom function)
    - Map
    - Array
    - Arrow operator(=>)
    - Simple Map Operator(!)
3. Path Expressions
    - Forward Step(child::, descendant::, ...)
    - Reverse Step(parent::, ...)
    - Node Test
    - Predicate([])
    - Abbreviated Syntax(@, ..)
4. Sequence Expressions(())
5. Arithmetic Expressions
    - Additive(+, -)
    - Multiplicative(*, div, idiv, mod)
    - Unary(+, -)
6. String Concatenation Expressions(||)
7. Comparison Expressions
    - Value Compare(eq, ne, lt, le, gt, ge)
    - Node Compare(is, <<, >>)
    - General Compare(=, !=, <, <=, >, >=)
8. Logical Expressions(and, or)
9. For Expressions(for)
10. Let Expressions(let)
11. Conditional Expressions(if)
12. Quantified Expressions(some, every)
13. Lookup(?)

### What is not supported

1. Namespace<br/>
Rabbit language doesn't care about prefixed tag names or xmlns attributes in tags. So, xmlns attribute is not treated as a namespace node, and a prefixed tag does not complain if no namespace for the prefix is specified in a document.

2. Limited Types<br/>
There are bunch of data type in XPath data model. You can check all the types in [https://www.w3.org/TR/xpath-datamodel-31/](https://www.w3.org/TR/xpath-datamodel-31/). Many of the types are not supported in Rabbit language and most of the data types in Rabbit language are simplified as string. It makes no sense to implement all the data type because there is no such a things like XML Schema Definition(xsd) in HTML.

3. Limited KindTest<br/>
In the XPath 3.1 document, there are 10 kinds of KindTest. But namespace-node test, processing-instruction test, schema-attribute test, schema-element test is not supported in Rabbit language because our parsing engine(/x/net/html) does not recognize them.

4. Sequence Type Check<br/>
In XPath 3.1, you can specify data types in linline function. It is looks like this.
`function($a as xs:string) as xs:string {$a}`.
This syntax is not a part of the Rabbit language. Inline function should like this.
`function($a) {$a}`.

5. Node Test with Argument<br/>
Node test with argument is not supported. For example, `element(person)`, `element(person, surgeon)`, `element(*, surgeon)`, `attribute(price)`, `attribute(*, xs:decimal)` are not allowed. But you can do `element()`, `attribute()`.

6. Whildcard Expressions<br/>
Only `*` wildcard is allowed in the Rabbit language. `NCName:*`, `*:NCName`, `BracedURILiteral*` are not supported since namespace is not a big deal in the Rabbit language.

## Notice

### Attribute node is custom *html.Node type

Rabbit language support attribute node. But /x/net/html package has no such a type(it only have 6 kinds of node) and treat attribute as a field of element node. So, in order to make attribute as a node, I had to make a custom *html.Node type. It have following fields.

- Type: html.NodeType(7).
- Parent: node(*html.Node) that is contain the attribute
- FirstChild, LastChild: `nil`
- PrevSibling, NextSibling: prev or next attribute node(*html.Node) of current one
- Data: attribute value(string).
- DataAtom: atomized Data(atom.Atom)
- Namespace: ""(empty string)
- Attr: Attr field contains only one html.Attribute item. Is has key, value pair for the attribute.

### Not well formed document will be transfromed

Rabbit language uses the /x/net/html package for parsing HTML. So, the type of the selected node will be *html.Node.
One thing that should know is that /x/net/html package is wrap a document with html, head, body tags if it is not well-formed.

For example, if your document looks like this

```html
<div>
  ...
</div>
```

/x/net/html package transforms the document to this internally.

```html
<html>
  <head></head>
  <body>
    <div>
      ...
    </div>
  </body>
</html>
```

So, in this example, XPath expression `/div` has no result because the root node is an `html`, not `div`.
Keep in mind this fact and otherwise, you can get confused.
