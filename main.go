package main

import (
	"fmt"
)

type ss struct {
	zz string
	abc
}

type abc struct {
	b byte
}

func main() {
	my := ss{zz: "Abc"}
	fmt.Println(my.abc.b)
}
