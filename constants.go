package errorwrap

import "errors"

const (
	UntrackedOrigin string = "untracked error: origin was not invoked by ErrorTrace"
)

var (
	ErrIndexOutOfRange = errors.New("index out of range")
)
