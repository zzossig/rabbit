package main

import (
	"fmt"
)

type AA struct{}

func (aa *AA) myFunc() {
	fmt.Println("haha")
}

type BB = AA

func main() {
	myVar := BB{}
	myVar.myFunc()
}
