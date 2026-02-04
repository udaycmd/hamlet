// Copyright 2025 Uday Tiwari. All rights reserved.
// Use of this source code is governed by MIT
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/udaycmd/hamlet/src/lexer"
	"github.com/udaycmd/hamlet/src/parser"
	"github.com/udaycmd/hamlet/src/token"
)

var (
	flags         *flag.FlagSet
	showHelp      bool
	showVersion   bool
	maxShowErrors int
)

func printVersion() {
	fmt.Println("Hamlet", HamletVersion)
}

func help() {
	flags.SetOutput(os.Stderr)

	printVersion()
	fmt.Println("Copyright (C) 2026 Uday Tiwari")
	fmt.Println("Usage: hamlet [option] ... [file] [filearg] ...")
	fmt.Println("Options:")
	flags.PrintDefaults()
}

func hamlet_main() error {
	if err := flags.Parse(os.Args[1:]); err != nil {
		if err == flag.ErrHelp {
			help()
			return nil
		}

		return err
	}

	if showHelp {
		help()
		return nil
	}

	if showVersion {
		printVersion()
		return nil
	}

	input := flags.Arg(0)
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
	_, err = p.Parse()
	fmt.Printf("\n\n%v\n\n", err)

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
	flags = &flag.FlagSet{}
	flags.SetOutput(io.Discard)

	flags.BoolVar(&showVersion, "version", false, "print the Hamlet version number and exit")
	flags.BoolVar(&showHelp, "help", false, "print help")
	flags.IntVar(&maxShowErrors, "max-errors", 5, "maximum numbers of errors to show")
}
