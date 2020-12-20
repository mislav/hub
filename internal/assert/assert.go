// Package assert provides functions for testing.
package assert

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
)

// Equal fails the test if the provided values do not pass the deep equality test
func Equal(t testing.TB, want, got interface{}, args ...interface{}) {
	t.Helper()
	if !reflect.DeepEqual(want, got) {
		msg := fmt.Sprint(args...)
		t.Errorf("%s\n%s", msg, cmp.Diff(want, got))
	}
}

// NotEqual is the negative assertion of Equal
func NotEqual(t testing.TB, want, got interface{}, args ...interface{}) {
	t.Helper()
	if reflect.DeepEqual(want, got) {
		msg := fmt.Sprint(args...)
		t.Errorf("%s\nUnexpected: <%#v>", msg, want)
	}
}

// Nil checks whether the value is nil
func Nil(t testing.TB, got interface{}) {
	t.Helper()
	if got != nil {
		t.Errorf("expected nil, got: %v", got)
	}
}

// NotNil is the negative assertion of Nil
func NotNil(t testing.TB, got interface{}) {
	t.Helper()
	if got == nil {
		t.Error("did not expect nil")
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
