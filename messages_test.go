package errorwrap

import (
	"errors"
	"fmt"
	"testing"

	"github.com/c2fo/testify/assert"
)

func TestMessages(t *testing.T) {
	e0 := errors.New("this is standard error")
	e1 := fmt.Errorf("error occured at modul x: %w", e0)
	e2 := New("Error occured", e1)
	e3 := New("Catch an error", e2)
	e4 := New("another basic err", nil)
	e5 := New("another high error", fmt.Errorf("err occured: %w", e3))
	assert.Len(t, TraceMessages(e0), 0)
	assert.Len(t, TraceMessages(e2), 2)
	assert.Len(t, TraceMessages(e3), 3)
	assert.Len(t, TraceMessages(e4), 1)
	assert.Len(t, TraceMessages(e5), 4)
	assert.Equal(t, "error occured at modul x: this is standard error", TraceMessages(e3)[0])
	assert.Equal(t, "Catch an error", TraceMessages(e3)[2])
}
