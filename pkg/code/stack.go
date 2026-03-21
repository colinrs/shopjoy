package code

import (
	"errors"
	"fmt"
	"io"
	"runtime"
	"strconv"
	"strings"
)

const (
	// 需要排除的堆栈其 package name
	mainPkg = "main"
	testPkg = "testing"
)

// stack represents a stack of program counters.
type stack []uintptr

// Format 实现 fmt.Formatter, 自定义打印格式
func (s stack) Format(st fmt.State, verb rune) {
	if verb == 'v' && st.Flag('+') {
		var mainFlag bool
		for _, pc := range s {
			fn := runtime.FuncForPC(pc)
			if fn == nil {
				fmt.Printf("[unknown] - pc: %d\n", pc)
				continue
			}
			if strings.HasPrefix(fn.Name(), mainPkg) {
				mainFlag = true
			} else if strings.HasPrefix(fn.Name(), testPkg) || mainFlag {
				break
			}

			file, line := fn.FileLine(pc - 1)
			_, _ = io.WriteString(st, "\n| ")
			_, _ = io.WriteString(st, fn.Name())
			_, _ = io.WriteString(st, "\n| \t")
			_, _ = io.WriteString(st, file)
			_, _ = io.WriteString(st, ":")
			_, _ = io.WriteString(st, strconv.Itoa(line))
		}
		return

	}
}

// callers copy 了 github.com/pkg/errors 中获取堆栈的代码,
// 但允许自定义深度
func callers(depth int) []uintptr {
	const numFrames = 32
	var pcs [numFrames]uintptr
	n := runtime.Callers(2+depth, pcs[:])
	return pcs[0:n]
}

type CodeMsgWithStack struct {
	c   *Err
	stk stack
}

func NewCodeMsgWithStack(c *Err) *CodeMsgWithStack {
	return &CodeMsgWithStack{
		c:   c,
		stk: callers(1),
	}
}

type withStack struct {
	err error
	stk stack
}

func WithStack(err error) error {
	if err == nil {
		return nil
	}

	var e *withStack
	if errors.As(err, &e) {
		return e
	}

	return &withStack{
		err: err,
		stk: callers(1),
	}
}

func WithStackByDepth(err error, depth int) error {
	if err == nil {
		return nil
	}

	var e *withStack
	if errors.As(err, &e) {
		return e
	}

	return &withStack{
		err: err,
		stk: callers(depth),
	}
}
func (w *withStack) Error() string {
	return w.err.Error()
}

func (w *withStack) Cause() error { return w.err }

// Unwrap provides compatibility for Go 1.13 error chains.
func (w *withStack) Unwrap() error {
	return w.err
}

func (w *withStack) String() string {
	return w.Error()
}

// Format implement fmt.Formatter, add trace to zap.Error(err).
func (w *withStack) Format(s fmt.State, verb rune) {
	switch verb {
	case 'T':
		_, _ = io.WriteString(s, "*withStack")
	case 'v':
		if s.Flag('#') {
			fmt.Fprintf(s, "&serr.withStack{Error: %#v, Stack: ...}", w.err)
			return
		}

		if s.Flag('+') {
			_, _ = io.WriteString(s, w.err.Error())
			_, _ = io.WriteString(s, "\nAttached stack trace:")
			w.stk.Format(s, verb)
			return
		}
		fallthrough
	case 's':
		_, _ = io.WriteString(s, w.Error())
	case 'q':
		_, _ = io.WriteString(s, strconv.Quote(w.Error()))
	default:
		_, _ = io.WriteString(s, w.Error())
	}
}
