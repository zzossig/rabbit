package main

import (
	"fmt"
	"strconv"
)

func main() {
	f := 64.2345
	s := strconv.FormatFloat(f, 'g', 1, 64)
	fmt.Println(s)
	s = strconv.FormatFloat(f, 'f', -1, 64)
	fmt.Println(s)
	fmt.Println(strconv.FormatInt(int64(66), 10))
}
