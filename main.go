package main

import (
	"fmt"
	"os"
	"os/user"

	"github.com/kiki-ki/go-monkey/repl"
)

func main() {
	u, err := user.Current()
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("Hello %s! This is the kiki-ki's Monkey programing language!\n", u.Name)
	fmt.Print("Feel free to type in commands.\n")
	repl.Start(os.Stdin, os.Stdout)
}
