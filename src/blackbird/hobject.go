// Copyright 2026 Uday Tiwari. All rights reserved.
// Use of this source code is governed by MIT
// license that can be found in the LICENSE file.

package blackbird

import (
	"errors"

	"github.com/udaycmd/hamlet/src/frontend/token"
	"github.com/udaycmd/hamlet/src/hmach/runtime"
)

type (
	HObject interface {
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
		Equals(HObject) bool
	}

	// HProc is a signature for the callable procedures/functions.
	HProc = func(args ...HObject) (ret HObject, err error)

	Callable interface {
		HObject

		// Call calls the Object and returns an another Object or Runtime Error.
		Call(...HObject) (HObject, error)
	}

	Iterable interface {
		HObject

		// Iterator returns a [IObject].
		Iterator() IObject
	}

	// IObject represents a [Iterable] HObject which can be iterated
	// via a loop like mechanism.
	IObject interface {
		HObject

		// Tells if there are more elements to iterate.
		Next() bool

		// Returns the key of the current element in iteration
		K() HObject

		// Returns the value of the current element in iteration
		V() HObject
	}

	// RObject is the root Runtime Object. Every other Object just
	// Overrides behaviour from this base implementation of [HObject].
	RObject struct{}
)

func (o *RObject) Value() any                                         { panic(runtime.ErrIllegalOperation) }
func (o *RObject) Type() string                                       { panic(runtime.ErrIllegalOperation) }
func (o *RObject) Operate(op token.Tok, rhs HObject) (HObject, error) { return nil, nil }
func (o *RObject) Falsy() bool                                        { return false }
func (o *RObject) Copy() HObject                                      { return nil }
func (o *RObject) Equals(a HObject) bool                              { return o == a }

// Empty Object
type EmptyObject struct {
	RObject
}

func (o *EmptyObject) Value() any                                         { return nil }
func (o *EmptyObject) Type() string                                       { return "<notype>" }
func (o *EmptyObject) Operate(op token.Tok, rhs HObject) (HObject, error) { return nil, nil }
func (o *EmptyObject) Falsy() bool                                        { return true }
func (o *EmptyObject) Copy() HObject                                      { return nil }
func (o *EmptyObject) Equals(_ HObject) bool                              { return true }

// Boolean Object
type BooleanObject struct {
	RObject
	Val bool
}

var (
	True  = &BooleanObject{Val: true}
	False = &BooleanObject{Val: false}
)

func (o *BooleanObject) Value() any            { return o.Val }
func (o *BooleanObject) Type() string          { return "<bool>" }
func (o *BooleanObject) Falsy() bool           { return !o.Val }
func (o *BooleanObject) Copy() HObject         { return &BooleanObject{Val: o.Val} }
func (o *BooleanObject) Equals(a HObject) bool { return o == a }

// Integer Object
type IntegerObject struct {
	RObject
	Val int64
}

func (o *IntegerObject) Value() any    { return o.Val }
func (o *IntegerObject) Type() string  { return "<int>" }
func (o *IntegerObject) Falsy() bool   { return o.Val == 0 }
func (o *IntegerObject) Copy() HObject { return &IntegerObject{Val: o.Val} }
func (o *IntegerObject) Equals(a HObject) bool {
	if t, ok := a.(*IntegerObject); ok {
		return o.Val == t.Val
	}
	return false
}
func (o *IntegerObject) Operate(op token.Tok, rhs HObject) (HObject, error) {
	switch rhs := rhs.(type) {
	case *IntegerObject:
		switch op {
		case token.PLUS:
			r := o.Val + rhs.Val
			if r == o.Val {
				return o, nil
			}

			return &IntegerObject{Val: r}, nil
		case token.MINUS:
			r := o.Val - rhs.Val
			if r == o.Val {
				return o, nil
			}

			return &IntegerObject{Val: r}, nil
		case token.STAR:
			r := o.Val * rhs.Val
			if r == o.Val {
				return o, nil
			}

			return &IntegerObject{Val: r}, nil
		case token.SLASH:
			r := o.Val / rhs.Val
			if r == o.Val {
				return o, nil
			}

			return &IntegerObject{Val: r}, nil
		case token.PERCENT:
			r := o.Val % rhs.Val
			if r == o.Val {
				return o, nil
			}

			return &IntegerObject{Val: r}, nil
		case token.AMPERSAND:
			r := o.Val & rhs.Val
			if r == o.Val {
				return o, nil
			}

			return &IntegerObject{Val: r}, nil
		case token.PIPE:
			r := o.Val | rhs.Val
			if r == o.Val {
				return o, nil
			}

			return &IntegerObject{Val: r}, nil
		case token.XOR:
			r := o.Val ^ rhs.Val
			if r == o.Val {
				return o, nil
			}

			return &IntegerObject{Val: r}, nil
		case token.LEFT_SHIFT:
			r := o.Val << uint64(rhs.Val)
			if r == o.Val {
				return o, nil
			}

			return &IntegerObject{Val: r}, nil
		case token.RIGHT_SHIFT:
			r := o.Val >> uint64(rhs.Val)
			if r == o.Val {
				return o, nil
			}

			return &IntegerObject{Val: r}, nil
		case token.LESS:
			if o.Val < rhs.Val {
				return True, nil
			}

			return False, nil
		case token.GREATER:
			if o.Val > rhs.Val {
				return True, nil
			}

			return False, nil
		case token.LESS_EQ:
			if o.Val <= rhs.Val {
				return True, nil
			}

			return False, nil
		case token.GREATER_EQ:
			if o.Val >= rhs.Val {
				return True, nil
			}

			return False, nil
		}
	}

	return nil, errors.New(runtime.ErrIllegalOperation)
}

// Hamlet Function/Procedure Object
type HamletProcObject struct {
	RObject
	Name string
	Proc HProc
}

func (o *HamletProcObject) Value() any                            { return o.Proc }
func (o *HamletProcObject) Type() string                          { return "<procedure>" }
func (o *HamletProcObject) Copy() HObject                         { return &HamletProcObject{Name: o.Name, Proc: o.Proc} }
func (o *HamletProcObject) Equals(_ HObject) bool                 { return false }
func (o *HamletProcObject) Call(args ...HObject) (HObject, error) { return o.Proc(args...) }
