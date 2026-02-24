// Copyright 2026 Uday Tiwari. All rights reserved.
// Use of this source code is governed by MIT
// license that can be found in the LICENSE file.

package utils

import (
	"reflect"
	"testing"
)

// --- Test Helpers ---

func Fail(t *testing.T, format string, args ...any) {
	t.Helper()
	t.Errorf(format, args...)
}

func ExpectEqual(t *testing.T, x, y any, format string, args ...any) {
	t.Helper()
	if !reflect.DeepEqual(x, y) {
		Fail(t, format, args...)
	}
}
