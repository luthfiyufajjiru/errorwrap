package errorwrap

import (
	"errors"
	"fmt"
	"testing"

	"github.com/c2fo/testify/assert"
)

var (
	errFoo error = errors.New("commons")
)

func TestChain(t *testing.T) {
	e0 := errors.New("this is standard error")
	e1 := fmt.Errorf("error occured at modul x: %w", e0)
	e2 := New("Error occured", e1)
	e3 := New("Catch an error", e2)
	e4 := New("another basic err", nil)
	e5 := New("another high error", fmt.Errorf("err occured: %w", e3))
	assert.True(t, Is(e2, e0))
	assert.True(t, Is(e3, e2))
	assert.True(t, Is(e5, e3))
	assert.True(t, Is(e3, e0))
	assert.True(t, Is(e5, e0))
	assert.False(t, Is(e2, e3))
	assert.False(t, Is(e0, e1))
	assert.False(t, Is(e0, e4))
	assert.False(t, Is(e0, e5))
	assert.False(t, Is(e3, e4))
	assert.False(t, Is(e5, e4))
}
