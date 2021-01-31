package main

import "fmt"

type bb interface {
	abc()
}

type cc struct {
	bb
}

func (c *cc) abc() {}

type aa struct {
	str string
}

func (a *aa) abc() {}

func main() {
	v1 := cc{}
	fmt.Println(v1.bb == nil)
	v1.bb = &aa{}
	fmt.Println(v1.bb == nil)
}
