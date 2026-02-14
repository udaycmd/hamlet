// Copyright 2026 Uday Tiwari. All rights reserved.
// Use of this source code is governed by MIT
// license that can be found in the LICENSE file.

package blackbird

import "github.com/udaycmd/hamlet/src/frontend/token"

type HObject interface {
	// Value returns the Object's actual value.
	Value() any

	// Type returns the name of the Object's type.
	Type() string

	// Operate performs a Unary or Binary operation on the Object
	// depending upon the presence of the second argument to the call.
	Operate(token.Tok, HObject) (HObject, error)

	// Falsy returns true if the value of Object if falsy in nature.
	Falsy() bool

	// Copy returns a deep-copy the Object
	Copy() HObject

	// Equal returns true if the Object is deep-equal to another Object.
	Equal(HObject) bool

	// Iterable returns true if the Object is iterable by a loop.
	Iterable() bool

	// Callable returns true if the Object can be called as a function.
	Callable() bool

	// Call calls the Object and returns an another Object or Runtime Error.
	Call(...HObject) (HObject, error)

	// Native returns true if the Object is a Native function.
	Native() bool
}

// RObject is the root Runtime Object. Every other Object just
// Overrides behaviour from this base implementation of [HObject].
type RObject struct{}

func (o *RObject) Value() any {
	panic(1)
}

func (o *RObject) Type() string {
	panic(1)
}

func (o *RObject) Operate(op token.Tok, rhs HObject) (HObject, error) {
	return nil, nil
}

func (o *RObject) Falsy() bool {
	return false
}

func (o *RObject) Copy() HObject {
	return nil
}

func (o *RObject) Equal(a HObject) bool {
	return o == a
}

func (o *RObject) Iterable() bool {
	return false
}

func (o *RObject) Callable() bool {
	return false
}

func (o *RObject) Call(...HObject) (HObject, error) {
	return nil, nil
}

func (o *RObject) Native() bool {
	return false
}
