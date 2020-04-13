// Package assert provides functions for testing.
package assert

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
)

// Equal makes the test as failed using default formatting if got is not equal to want.
func Equal(t testing.TB, want, got interface{}, args ...interface{}) {
	t.Helper()
	if !reflect.DeepEqual(want, got) {
		msg := fmt.Sprint(args...)
		t.Errorf("%s\n%s", msg, cmp.Diff(want, got))
	}
}

// NotEqual makes the test as failed using default formatting if got is equal to want.
func NotEqual(t testing.TB, want, got interface{}, args ...interface{}) {
	t.Helper()
	if reflect.DeepEqual(want, got) {
		msg := fmt.Sprint(args...)
		t.Errorf("%s\nUnexpected: <%#v>", msg, want)
	}
}

// T makes the test as failed using default formatting if ok is false.
func T(t testing.TB, ok bool, args ...interface{}) {
	t.Helper()
	if !ok {
		msg := fmt.Sprint(args...)
		t.Errorf("%s\nFailure", msg)
	}
}
