package main

import (
	"fmt"
	"os"
	"os/user"

	"github.com/zzossig/rabbit/repl"
)

func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Hello %s!\nThis is the Rabbit xpath 3.1 implementation language!\n", user.Username)
	fmt.Printf("Feel free to type in commands\n")
	repl.Start(os.Stdin, os.Stdout)
}
