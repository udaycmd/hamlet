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

func Colorize(text, color string) string {
	return color + text + Reset
}

func help(w io.Writer) {
	fmt.Fprintln(w, "Hamlet", Colorize(HamletVersion, Yellow))
	fmt.Fprintln(w, "Copyright (C) 2025 Uday Tiwari")
	fmt.Fprintln(w, "Usage: hamlet [options] ... [file] [fileargs] ...")
	fmt.Fprintln(w, "Options:")
	fmt.Fprintln(w, "-version:		print the Hamlet version number and exit")
	fmt.Fprintln(w, "-nocheck:		disable typechecking")
}

func hamlet_main() error {
	v := flag.Bool("version", false, "print the Hamlet version number and exit")
	flag.Parse()

	if *v {
		help(os.Stdout)
	}

	return nil
}

func main() {
	if status := hamlet_main(); status != nil {
		help(os.Stderr)
		fmt.Fprintf(os.Stderr, "%s: %s\n", Colorize("Error", Red), status.Error())
		os.Exit(1)
	}

	// safe exit with success
	os.Exit(0)
}
