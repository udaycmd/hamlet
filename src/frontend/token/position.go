// Copyright 2026 Uday Tiwari. All rights reserved.
// Use of this source code is governed by MIT
// license that can be found in the LICENSE file.

package token

import (
	"fmt"
	"sort"
)

// [Position] represents a position in the SourceManager.
type Position int

// [SrcPos] contains the positon of token in the souce code
type SrcPos struct {
	FileName string
	Offset   int
	Line     int
	Column   int
}

const (
	InvalidPos Position = 0
)

// Returns true if position is valid
func (p Position) IsValid() bool {
	return p != InvalidPos
}

// Returns true if srcpos is valid
func (sp *SrcPos) IsValid() bool {
	return sp.Line > 0
}

// Returns a empty string if [SrcPos] is invalid
func (sp *SrcPos) String() string {
	if sp.IsValid() {
		s := sp.FileName
		s += fmt.Sprintf("(%d, %d)", sp.Line, sp.Column)
		return s
	}

	return ""
}

type SourceManager struct {
	Files    []*SourceHandle // list of files in order of insertion
	LastLook *SourceHandle   // last file looked up (cache)
	Base     int             // base for the next file to come
}

func NewSourceManager() *SourceManager {
	return &SourceManager{Base: 1}
}

// size of all files in the source manager in bytes
func (sm *SourceManager) Len() uint64 {
	var bytes uint64
	for _, sh := range sm.Files {
		bytes += uint64(sh.Len)
	}

	return bytes
}

// AddFile adds a new file in the [SourceManager]
func (sm *SourceManager) AddFile(filename string, base, size int) *SourceHandle {
	if base < 0 {
		base = sm.Base
	}

	if base < sm.Base || size < 0 {
		panic("illegal file base or size")
	}

	sh := &SourceHandle{
		manager: sm,
		Name:    filename,
		Base:    base,
		Len:     size,
		Lines:   []int{0},
	}

	base += size + 1 // because EOF also gets a position

	if base < 0 {
		panic("offset overflow")
	}

	// add the file to the manager
	sm.Base = base
	sm.Files = append(sm.Files, sh)
	sm.LastLook = sh
	return sh
}

// [SourceManager.SrcPos] returns the [SrcPos]
// inside the file in which p resides or returns [nil]
func (sm *SourceManager) SrcPos(p Position) *SrcPos {
	if h := sm.getHandle(p); h != nil {
		return h.SrcPos(p)
	}

	return &SrcPos{Line: -1}
}

// Helper function for [SourceManager.File]
func (sm *SourceManager) getHandle(p Position) *SourceHandle {
	// check for valid position
	if !p.IsValid() {
		return nil
	}

	// cache check
	sh := sm.LastLook
	if sh != nil && sh.Base <= int(p) && int(p) <= sh.Base+sh.Len {
		return sh
	}

	// cache miss, do a binary search in [SourceManager.Files]
	i := sort.Search(len(sm.Files), func(i int) bool { return sm.Files[i].Base > int(p) }) - 1

	if i >= 0 {
		sh := sm.Files[i]

		// sh.Base <= p by definition of search algorithm
		if int(p) <= sh.Base+sh.Len {
			sm.LastLook = sh
			return sh
		}
	}

	return nil
}

type SourceHandle struct {
	Name  string // name of the file with extension, like "file.hm"
	Base  int    // provides base positon for file i.e, [Base..Base+Len]
	Len   int    // size of the file
	Lines []int  // Lines contains the offset of the first character for each line (first entry is always 0)

	manager *SourceManager // handle manager of the file
}

// LineCount returns the current number of lines
func (sh *SourceHandle) LineCount() int {
	return len(sh.Lines)
}

// AddLine adds a new line to [SourceHandle.Lines]
func (sh *SourceHandle) AddLine(offset int) {
	l := len(sh.Lines)
	if (l == 0 || sh.Lines[l-1] < offset) && offset < sh.Len {
		sh.Lines = append(sh.Lines, offset)
	}
}

// [SourceManager.TapePos] returns the position in the [SourceManager].
func (sh *SourceHandle) TapePos(offset int) Position {
	if offset > sh.Len {
		panic("invalid file offset")
	}

	return Position(sh.Base + offset)
}

// [SourceHandle.Offset] return the handle's offset from the tape's offset
func (sh *SourceHandle) Offset(p Position) int {
	if int(p) < sh.Base || int(p) > sh.Base+sh.Len {
		panic("illegal tape position")
	}

	return int(p) - sh.Base
}

// [SourceHandle.SrcPos] return a srcpos of the received handle
func (sh *SourceHandle) SrcPos(p Position) *SrcPos {
	srcPos := &SrcPos{}
	if p.IsValid() {
		if int(p) < sh.Base || int(p) > sh.Base+sh.Len {
			panic("illegal position")
		}

		offset := int(p) - sh.Base
		srcPos.Offset = offset
		srcPos.FileName = sh.Name

		i := sort.Search(len(sh.Lines), func(i int) bool { return sh.Lines[i] > offset }) - 1
		srcPos.Column = offset - sh.Lines[i] + 1
		srcPos.Line = i + 1
	}

	return srcPos
}
