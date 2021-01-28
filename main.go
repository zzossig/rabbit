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
	fmt.Println("eq" == "EQ")
}
