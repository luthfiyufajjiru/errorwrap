package errorwrap

import (
	"errors"
	"fmt"
	"regexp"
	"runtime"
	"strings"
)

type (
	ErrorTrace struct {
		Message, Site string
		Origin        error
		Messages      []string
		Sites         []string
	}
)

func (t *ErrorTrace) Error() (res string) {
	origin := "}"
	var trace *ErrorTrace
	if originOk, traceOk := t.Origin != nil, errors.As(t.Origin, &trace); originOk && !traceOk {
		origin = fmt.Sprintf(`,"origin":"%s"}`, t.Origin.Error())
	} else if originOk && traceOk {
		origin = fmt.Sprintf(`,"origin":%s}`, t.Origin.Error())
	}
	res = fmt.Sprintf(`{"message":"%s","site":"%s"%s`, t.Message, t.Site, origin)
	return
}

// If the input is ErrorTrace type it return the current message, else is err.Error()
func Message(err error) string {
	x := new(ErrorTrace)
	if errors.As(err, &x) {
		return x.Message
	}
	return err.Error()
}

func New(message string, origin error) error {
	var (
		site   string
		trace  *ErrorTrace
		result *ErrorTrace = new(ErrorTrace)
	)

	result.Message = message
	result.Origin = origin

	_, filename, line, ok := runtime.Caller(1)
	if ok {
		site = fmt.Sprintf("%s:%d", filename, line)
		result.Site = site
	}

	if originOk, traceOk := origin != nil, errors.As(origin, &trace); originOk && traceOk {
		if trace.Sites != nil {
			result.Sites = trace.Sites
			result.Messages = trace.Messages
		}
		result.Sites = append(result.Sites, site)
		result.Messages = append(result.Messages, message)
	} else if !originOk && !traceOk {
		result.Messages = []string{message}
		result.Sites = []string{site}
	} else if originOk && !traceOk {
		result.Messages = []string{origin.Error(), message}
		result.Sites = []string{UntrackedOrigin, site}
	}

	return result
}

// The trace is descending -> the oldest is on the top
func TraceMessages(err error) (res []string) {
	var x *ErrorTrace
	switch errors.As(err, &x) {
	case true:
		return x.Messages
	default:
		return
	}
}

// The trace is descending -> the oldest is on the top
func TraceSites(err error) (res []string) {
	var x *ErrorTrace
	switch errors.As(err, &x) {
	case true:
		return x.Sites
	default:
		return
	}
}

// Is function like in standard errors.IS(err, target error) but this is returning true if errors.New("foo") matched with errors.New("foo") while the standards errors.Is(...) is returning false. Then this function could matching the error trace from the grandchild to the root error.
// Caution do not misplaced the root with the trace, if root is placed in err while the trace is placed in the target it returns false
func Is(err, target error) bool {
	if target == nil || err == nil {
		return false
	}
	var x *ErrorTrace
	ok := errors.As(err, &x)
	if ok {
		if originOk := x.Origin != nil; originOk {
			if errors.Is(err, target) {
				return true
			}
			return Is(x.Origin, target)
		} else if !originOk {
			return errors.Is(err, target)
		}
	}
	validate := errors.Is(err, target) || err.Error() == target.Error()
	return validate
}

func postRecover(err error) {
	var switchVar bool
	errx, _ := err.(*ErrorTrace)
	ln := len(errx.Sites)

	if _, filename, line, ok := runtime.Caller(4); ok {
		if 1 < ln {
			match, _ := regexp.MatchString(`.runtime\/panic.go*`, filename)
			if _, _filename, _line, _ok := runtime.Caller(6); match && _ok {
				filename = _filename
				line = _line
				switchVar = true
			}
			x := fmt.Sprintf("%s:%d", filename, line)
			errx.Site = x
			if switchVar {
				errx.Sites[1] = x
			} else if !switchVar {
				errx.Sites[0] = x
			}
		}
	}

	if _, filename, line, ok := runtime.Caller(5); ok {
		if switchVar {
			errx.Sites[0] = fmt.Sprintf("%s:%d", filename, line)
		} else if !switchVar {
			errx.Sites[1] = fmt.Sprintf("%s:%d", filename, line)
		}
	}
}

func PostRecover(message string, r any) (err error) {
	switch r := r.(type) {
	case string:
		err = New(message, errors.New(r))
		postRecover(err)
	case error:
		var s runtime.Error
		if errors.As(r, &s) && strings.Contains(s.Error(), "index out of range") {
			r = fmt.Errorf("%w: %s", ErrIndexOutOfRange, s.Error())
		}
		err = New(message, r)
		postRecover(err)
	default:
		rr := fmt.Errorf("%v", r)
		err = New(message, rr)
		postRecover(err)
	}
	return
}
