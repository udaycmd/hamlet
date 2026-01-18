// Copyright 2026 Uday Tiwari. All rights reserved.
// Use of this source code is governed by MIT
// license that can be found in the LICENSE file.

package token

// [Position] represents a position in the SourceManager.
type Position int

const (
	InvalidPos Position = 0
)

// returns true if position is valid
func (p Position) IsValid() bool {
	return p != InvalidPos
}

