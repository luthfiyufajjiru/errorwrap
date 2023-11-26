package errorwrap

import (
	"errors"
	"fmt"
	"testing"

	"github.com/c2fo/testify/assert"
)

func TestTrace(t *testing.T) {
	e0 := errors.New("this is standard error")
	e1 := fmt.Errorf("error occured at modul x: %w", e0)
	e2 := New("Error occured", e1)
	e3 := New("Catch an error", e2)
	e4 := New("another basic err", nil)
	e5 := New("another high error", fmt.Errorf("err occured: %w", e3))
	assert.Len(t, TraceSites(e0), 0)
	assert.Len(t, TraceSites(e2), 2)
	assert.Len(t, TraceSites(e3), 3)
	assert.Len(t, TraceSites(e4), 1)
	assert.Len(t, TraceSites(e5), 4)
}
