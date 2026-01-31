// Copyright 2025 Uday Tiwari. All rights reserved.
// Use of this source code is governed by MIT
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/udaycmd/hamlet/src/lexer"
	"github.com/udaycmd/hamlet/src/parser"
	"github.com/udaycmd/hamlet/src/token"
)

var (
	showHelp      bool
	showVersion   bool
	maxShowErrors int
)

func printVersion() {
	fmt.Println("Hamlet Interpreter", HamletVersion)
}

func help() {
	printVersion()
	fmt.Println("Copyright (C) 2026 Uday Tiwari")
	fmt.Println("Usage: hamlet [option] ... [file] [filearg] ...")
	fmt.Println("Options:")
	flag.PrintDefaults()
}

func hamlet_main() error {
	flag.Usage = help

	if showHelp {
		flag.Usage()
		return nil
	}

	if showVersion {
		printVersion()
		return nil
	}

	input := flag.Arg(0)
	if input == "" {
		return fmt.Errorf("no input provided")
	}

	data, err := os.ReadFile(input)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("%s: No such file found", input)
		}

		return err
	}

	// TODO: Change this
	test := token.NewSourceManager()
	testFile := test.AddFile(input, -1, len(data))
	p := parser.NewParser(testFile, data, 10, lexer.NoAsi, os.Stdout)
	p.Parse()

	return nil
}

func main() {
	if err := hamlet_main(); err != nil {
		fmt.Fprintf(os.Stderr, "\033[1;31merror: \033[0m%v\n", err)
		os.Exit(1)
	}
}

// all flags
func init() {
	flag.BoolVar(&showVersion, "version", false, "print the Hamlet version number and exit")
	flag.BoolVar(&showHelp, "help", false, "print help")
	flag.IntVar(&maxShowErrors, "max-errors", 5, "maximum numbers of errors to show")
	flag.Parse()
}
