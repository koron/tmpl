package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/koron/tmpl"
)

func main() {
	var (
		dataGo string
	)
	flag.StringVar(&dataGo, "data-go", "", "load data from go source code")
	flag.Parse()

	var dataFunc tmpl.DataFunc
	if dataGo != "" {
		tmpl.SourceGosrc = dataGo
		dataFunc = tmpl.LoadGosrc
	}

	err := tmpl.Execute(dataFunc, os.Stdout, flag.Args()...)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
