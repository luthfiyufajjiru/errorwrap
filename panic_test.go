package errorwrap

import (
	"fmt"
	"testing"

	"github.com/c2fo/testify/assert"
)

func panicTrigger() {
	panic(errFoo)
}

func panicIndexTrigger() {
	a := []string{"a", "v"}
	c := a[2]
	fmt.Println(c)
}

func recoverTest() (err error) {
	defer func() {
		r := recover()
		err = PostRecover("recovered err", r)
	}()
	panicTrigger()
	return
}

func recoverIndexTest() (err error) {
	defer func() {
		r := recover()
		err = PostRecover("recovered err", r)
	}()
	panicIndexTrigger()
	return
}

func TestPanic(t *testing.T) {
	e := recoverTest()
	e1 := recoverIndexTest()
	assert.True(t, Is(e, errFoo))
	assert.True(t, Is(e1, ErrIndexOutOfRange))
}
