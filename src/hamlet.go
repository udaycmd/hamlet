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
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Yellow = "\033[33m"
)

var (
	flags = flag.FlagSet{}
)

func Colorize(text, color string) string {
	return color + text + Reset
}

func printVersion(w io.Writer) {
	fmt.Fprintln(w, "Hamlet Interpreter", Colorize(HamletVersion, Yellow))
}

func help(w io.Writer) {
	flags.SetOutput(w)

	printVersion(w)
	fmt.Fprintln(w, "Copyright (C) 2025 Uday Tiwari")
	fmt.Fprintln(w, "Usage: hamlet [option] ... [file] [filearg] ...")
	fmt.Fprintln(w, "Options:")
	flags.PrintDefaults()
}

func hamlet_main(args []string) error {
	v := flags.Bool("version", false, "print the Hamlet version number and exit")
	nc := flags.Bool("nocheck", false, "disable typechecking")
	h := flags.Bool("help", false, "print help")

	flags.SetOutput(io.Discard)
	err := flags.Parse(args)

	if err != nil {
		if err == flag.ErrHelp {
			help(os.Stdout)
		}
		return err
	}

	if *h {
		help(os.Stdout)
		return nil
	}

	if *v {
		printVersion(os.Stdout)
	}

	if *nc {
	}

	return nil
}

func main() {
	if status := hamlet_main(os.Args[1:]); status != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", Colorize("Error", Red), status.Error())
		os.Exit(1)
	}

	// safe exit with success
	os.Exit(0)
}
