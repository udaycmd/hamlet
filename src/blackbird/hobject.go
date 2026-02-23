// Copyright 2026 Uday Tiwari. All rights reserved.
// Use of this source code is governed by MIT
// license that can be found in the LICENSE file.

package blackbird

import (
	"errors"
	"math"

	"github.com/udaycmd/hamlet/src/frontend/token"
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

func (o *RObject) Value() any                                         { panic("TODO: change this") }
func (o *RObject) Type() string                                       { panic("TODO: change this") }
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
	case *FloatObject:
		switch op {
		case token.PLUS:
			return &FloatObject{Val: float64(o.Val) + rhs.Val}, nil
		case token.MINUS:
			return &FloatObject{Val: float64(o.Val) - rhs.Val}, nil
		case token.STAR:
			return &FloatObject{Val: float64(o.Val) * rhs.Val}, nil
		case token.SLASH:
			return &FloatObject{Val: float64(o.Val) / rhs.Val}, nil
		case token.LESS:
			if float64(o.Val) < rhs.Val {
				return True, nil
			}

			return False, nil
		case token.GREATER:
			if float64(o.Val) > rhs.Val {
				return True, nil
			}

			return False, nil
		case token.LESS_EQ:
			if float64(o.Val) <= rhs.Val {
				return True, nil
			}

			return False, nil
		case token.GREATER_EQ:
			if float64(o.Val) >= rhs.Val {
				return True, nil
			}

			return False, nil
		}
	case *CharObject:
		switch op {
		case token.PLUS:
			return &CharObject{Val: rune(o.Val) + rhs.Val}, nil
		case token.MINUS:
			return &CharObject{Val: rune(o.Val) - rhs.Val}, nil
		case token.LESS:
			if o.Val < int64(rhs.Val) {
				return True, nil
			}

			return False, nil
		case token.GREATER:
			if o.Val > int64(rhs.Val) {
				return True, nil
			}

			return False, nil
		case token.LESS_EQ:
			if o.Val <= int64(rhs.Val) {
				return True, nil
			}

			return False, nil
		case token.GREATER_EQ:
			if o.Val >= int64(rhs.Val) {
				return True, nil
			}

			return False, nil
		}
	}

	return nil, errors.New("TODO: change this")
}

// Real/Floating Point Number Object
type FloatObject struct {
	RObject
	Val float64
}

func (o *FloatObject) Value() any    { return o.Val }
func (o *FloatObject) Type() string  { return "<float>" }
func (o *FloatObject) Falsy() bool   { return math.IsNaN(o.Val) }
func (o *FloatObject) Copy() HObject { return &FloatObject{Val: o.Val} }
func (o *FloatObject) Equals(a HObject) bool {
	if t, ok := a.(*FloatObject); ok {
		return o.Val == t.Val
	}
	return false
}
func (o *FloatObject) Operate(op token.Tok, rhs HObject) (HObject, error) {
	switch rhs := rhs.(type) {
	case *FloatObject:
		switch op {
		case token.PLUS:
			r := o.Val + rhs.Val
			if r == o.Val {
				return o, nil
			}

			return &FloatObject{Val: r}, nil
		case token.MINUS:
			r := o.Val - rhs.Val
			if r == o.Val {
				return o, nil
			}

			return &FloatObject{Val: r}, nil
		case token.STAR:
			r := o.Val * rhs.Val
			if r == o.Val {
				return o, nil
			}

			return &FloatObject{Val: r}, nil
		case token.SLASH:
			r := o.Val / rhs.Val
			if r == o.Val {
				return o, nil
			}

			return &FloatObject{Val: r}, nil
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
	case *IntegerObject:
		switch op {
		case token.PLUS:
			r := o.Val + float64(rhs.Val)
			if r == o.Val {
				return o, nil
			}

			return &FloatObject{Val: r}, nil
		case token.MINUS:
			r := o.Val - float64(rhs.Val)
			if r == o.Val {
				return o, nil
			}

			return &FloatObject{Val: r}, nil
		case token.STAR:
			r := o.Val * float64(rhs.Val)
			if r == o.Val {
				return o, nil
			}

			return &FloatObject{Val: r}, nil
		case token.SLASH:
			r := o.Val / float64(rhs.Val)
			if r == o.Val {
				return o, nil
			}

			return &FloatObject{Val: r}, nil
		case token.LESS:
			if o.Val < float64(rhs.Val) {
				return True, nil
			}

			return False, nil
		case token.GREATER:
			if o.Val > float64(rhs.Val) {
				return True, nil
			}

			return False, nil
		case token.LESS_EQ:
			if o.Val <= float64(rhs.Val) {
				return True, nil
			}

			return False, nil
		case token.GREATER_EQ:
			if o.Val >= float64(rhs.Val) {
				return True, nil
			}

			return False, nil
		}
	}

	return nil, errors.New("TODO: change this")
}

// Character Object
type CharObject struct {
	RObject
	Val rune
}

func (o *CharObject) Value() any    { return o.Val }
func (o *CharObject) Type() string  { return "<char>" }
func (o *CharObject) Falsy() bool   { return o.Val == 0 }
func (o *CharObject) Copy() HObject { return &CharObject{Val: o.Val} }
func (o *CharObject) Equals(a HObject) bool {
	if t, ok := a.(*CharObject); ok {
		return o.Val == t.Val
	}
	return false
}
func (o *CharObject) Operate(op token.Tok, rhs HObject) (HObject, error) {
	switch rhs := rhs.(type) {
	case *CharObject:
		switch op {
		case token.PLUS:
			r := o.Val + rhs.Val
			if r == o.Val {
				return o, nil
			}

			return &CharObject{Val: r}, nil
		case token.MINUS:
			r := o.Val - rhs.Val
			if r == o.Val {
				return o, nil
			}

			return &CharObject{Val: r}, nil
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
	case *IntegerObject:
		switch op {
		case token.PLUS:
			r := o.Val + rune(rhs.Val)
			if r == o.Val {
				return o, nil
			}

			return &CharObject{Val: r}, nil
		case token.MINUS:
			r := o.Val - rune(rhs.Val)
			if r == o.Val {
				return o, nil
			}

			return &CharObject{Val: r}, nil
		case token.LESS:
			if int64(o.Val) < rhs.Val {
				return True, nil
			}

			return False, nil
		case token.GREATER:
			if int64(o.Val) > rhs.Val {
				return True, nil
			}

			return False, nil
		case token.LESS_EQ:
			if int64(o.Val) <= rhs.Val {
				return True, nil
			}

			return False, nil
		case token.GREATER_EQ:
			if int64(o.Val) >= rhs.Val {
				return True, nil
			}

			return False, nil
		}
	}

	return nil, errors.New("TODO: change this")
}

type StrObject struct {
	RObject
	Val string
}

func (o *StrObject) Value() any    { return o.Val }
func (o *StrObject) Type() string  { return "<string>" }
func (o *StrObject) Falsy() bool   { return len(o.Val) == 0 }
func (o *StrObject) Copy() HObject { return &StrObject{Val: o.Val} }
func (o *StrObject) Equals(a HObject) bool {
	if t, ok := a.(*StrObject); ok {
		return o.Val == t.Val
	}
	return false
}
func (o *StrObject) Operate(op token.Tok, rhs HObject) (HObject, error) {
	switch rhs := rhs.(type) {
	case *StrObject:
		switch op {
		case token.PLUS:
			// TODO: check for resulting string length
			return &StrObject{Val: o.Val + rhs.Val}, nil
		}
	}

	return nil, errors.New("TODO: change this")
}
func (o *StrObject) Iterator() IObject {
	v := []rune(o.Val)

	return &StrIterator{
		str: v,
		sz:  len(v),
	}
}

// [StrIterator] represents a [StrObject] iterator.
type StrIterator struct {
	RObject
	index int
	str   []rune
	sz    int
}

func (io *StrIterator) Value() any            { return io.index }
func (io *StrIterator) Type() string          { return "<str_iterator>" }
func (io *StrIterator) Falsy() bool           { return true }
func (io *StrIterator) Copy() HObject         { return &StrIterator{index: io.index, str: io.str, sz: io.sz} }
func (io *StrIterator) Equals(_ HObject) bool { return false }
func (io *StrIterator) K() HObject            { return &IntegerObject{Val: int64(io.index)} }
func (io *StrIterator) V() HObject            { return &CharObject{Val: io.str[io.index-1]} }
func (io *StrIterator) Next() bool {
	io.index += 1
	return io.index <= io.sz
}

// Array Object
type Array struct {
	RObject
	Val []HObject
}

func (o *Array) Value() any    { return o.Val }
func (o *Array) Type() string  { return "<array>" }
func (o *Array) Falsy() bool   { return len(o.Val) == 0 }
func (o *Array) Copy() HObject { return &Array{Val: o.Val} }
func (o *Array) Equals(a HObject) bool {
	t, ok := a.(*Array)
	if !ok {
		return false
	}

	if len(o.Val) != len(t.Val) {
		return false
	}

	for i, o := range o.Val {
		if !o.Equals(t.Val[i]) {
			return false
		}
	}
	return true
}
func (o *Array) Iterator() IObject {
	return &ArrayIterator{
		objs: o.Val,
		sz:   len(o.Val),
	}
}

// [ArrayIterator] represents a [Array] iterator.
type ArrayIterator struct {
	RObject
	index int
	objs  []HObject
	sz    int
}

func (io *ArrayIterator) Value() any   { return io.index }
func (io *ArrayIterator) Type() string { return "<array_iterator>" }
func (io *ArrayIterator) Falsy() bool  { return true }
func (io *ArrayIterator) Copy() HObject {
	return &ArrayIterator{index: io.index, objs: io.objs, sz: io.sz}
}
func (io *ArrayIterator) Equals(_ HObject) bool { return false }
func (io *ArrayIterator) K() HObject            { return &IntegerObject{Val: int64(io.index - 1)} }
func (io *ArrayIterator) V() HObject            { return io.objs[io.index-1] }
func (io *ArrayIterator) Next() bool {
	io.index += 1
	return io.index <= io.sz
}

// Hamlet Function/Procedure Object
type HamletProcObject struct {
	RObject
	Name string
	proc HProc
}

func (o *HamletProcObject) Value() any                            { return o.proc }
func (o *HamletProcObject) Type() string                          { return "<procedure>" }
func (o *HamletProcObject) Copy() HObject                         { return &HamletProcObject{Name: o.Name, proc: o.proc} }
func (o *HamletProcObject) Equals(_ HObject) bool                 { return false }
func (o *HamletProcObject) Call(args ...HObject) (HObject, error) { return o.proc(args...) }
