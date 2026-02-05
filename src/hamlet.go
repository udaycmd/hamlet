// Copyright 2025 Uday Tiwari. All rights reserved.
// Use of this source code is governed by MIT
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"io"
	"os"
)

const (
	HamletVersion string = "0.1.0-dev"
)

var (
	flags         *flag.FlagSet
	showHelp      bool
	showVersion   bool
	maxShowErrors int

	printVersion = func() {
		fmt.Println("Hamlet", HamletVersion)
	}

	help = func() {
		flags.SetOutput(os.Stderr)

		printVersion()
		fmt.Println("Copyright (C) 2026 Uday Tiwari")
		fmt.Println("Usage: hamlet [option] ... [file] [filearg] ...")
		fmt.Println("Options:")
		flags.PrintDefaults()
	}
)

func driver() error {
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

	return nil
}

func main() {
	if err := driver(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// all flags
func init() {
	flags = &flag.FlagSet{}
	flags.SetOutput(io.Discard)

	flags.BoolVar(&showVersion, "version", false, "display product version")
	flags.BoolVar(&showHelp, "help", false, "display help")
	flags.IntVar(&maxShowErrors, "max-errors", 5, "maximum number of errors to show")
}
