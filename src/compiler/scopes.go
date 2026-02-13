// Copyright 2026 Uday Tiwari. All rights reserved.
// Use of this source code is governed by MIT
// license that can be found in the LICENSE file.

package compiler

type Scope uint8

const (
	Module Scope = iota
	Local
	Closure
)
