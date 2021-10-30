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
	fmt.Printf("Konitiwa %s! This is kiki-ki's Monkey programing language!\n\n", u.Name)
	fmt.Println("Usage:")
	fmt.Printf("\tHelp: 'h' or 'help'\n")
	fmt.Printf("\tEscape: 'q' or 'exit'\n\n")
	repl.Start(os.Stdin, os.Stdout)
}
