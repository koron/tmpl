package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/koron/tmpl"
)

func main() {
	flag.Parse()

	err := tmpl.Execute(nil, os.Stdout, flag.Args()...)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
