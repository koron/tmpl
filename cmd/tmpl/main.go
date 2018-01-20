package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/koron/tmpl"
)

func main() {
	var (
		dataGo  string
		outFile string
	)
	flag.StringVar(&dataGo, "data-go", "", "load data from go source code")
	flag.StringVar(&outFile, "o", "", "filename to output")
	flag.Parse()

	var dataFunc tmpl.DataFunc
	if dataGo != "" {
		tmpl.SourceGosrc = dataGo
		dataFunc = tmpl.LoadGosrc
	}

	var w io.Writer
	if outFile != "" {
		f, err := os.Create(outFile)
		if err != nil {
			fmt.Fprintln(os.Stderr, "failed to open output file: %s", err)
			os.Exit(1)
		}
		w = f
	} else {
		w = os.Stdout
	}

	err := tmpl.Execute(dataFunc, os.Stdout, flag.Args()...)
	if wc, ok := w.(io.Closer); ok {
		wc.Close()
	}

	if err != nil {
		if outFile != "" {
			os.Remove(outFile)
		}
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
