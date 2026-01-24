// Copyright 2025 Uday Tiwari. All rights reserved.
// Use of this source code is governed by MIT
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
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
	fmt.Println("Copyright (C) 2025 Uday Tiwari")
	fmt.Println("Usage: hamlet [option] ... [file] [filearg] ...")
	fmt.Println("Options:")
	flag.PrintDefaults()
}

func hamlet_main() error {
	if showHelp {
		help()
		return nil
	}

	if showVersion {
		printVersion()
		return nil
	}

	input := flag.Arg(0)
	data, err := os.ReadFile(input)
	if err != nil {
		return err
	}

	input, err = filepath.Abs(input)
	if err != nil {
		return err
	}

	// TODO: Change this
	fmt.Print(input, "\n", string(data))

	return nil
}

func main() {
	if err := hamlet_main(); err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", "\033[31mError\033[0m", err)
		os.Exit(1)
	}
}

// all flags
func init() {
	flag.BoolVar(&showVersion, "version", false, "print the Hamlet version number and exit")
	flag.BoolVar(&showHelp, "help", false, "print help")
	flag.IntVar(&maxShowErrors, "maxerr", 5, "maximum numbers of errors to show")
	flag.Parse()
}
