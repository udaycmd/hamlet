// Copyright 2025 Uday Tiwari. All rights reserved.
// Use of this source code is governed by MIT
// license that can be found in the LICENSE file.

package hmach

type Opcode uint8

const (
	Nop Opcode = iota + 100
	Dup
	Halt
	Push
	Pop
	Ret
)
