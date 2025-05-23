package main

import (
	"fmt"
	"os"
	"os/user"

	"github.com/boatnoah/monkey/repl"
)

func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Hello %s! This is the monkey programming language\n", user.Username)
	repl.Start(os.Stdin, os.Stdout)
}
