package main

import (
	"fmt"
	"os"

	"github.com/koron/tmpl/tmpl"
)

func main() {
	err := tmpl.Execute(os.Stdin, os.Stdout, os.Args[1:]...)
	if err != nil {
		fmt.Println(os.Stderr, err)
		os.Exit(1)
	}
}
