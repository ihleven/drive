package errors

import (
	"errors"
	"fmt"
	"runtime"
	"strings"

	"github.com/palantir/stacktrace/cleanpath"
)

var CleanPath = cleanpath.RemoveGoPath

func Errorf(msg string, vals ...interface{}) error {
	return create(nil, NoCode, msg, vals...)
}

func New(code ErrorCode, msg string, vals ...interface{}) error {
	return create(nil, code, msg, vals...)
}

func Wrap(cause error, msg string, vals ...interface{}) error {
	if cause == nil {
		// Allow calling Propagate without checking whether there is error
		return nil
	}
	return create(cause, NoCode, msg, vals...)
}

func Propagate(cause error) error {
	if cause == nil {
		// Allow calling Propagate without checking whether there is error
		return nil
	}
	return create(cause, NoCode, "")
}

func Augment(cause error, code ErrorCode, msg string, vals ...interface{}) error {
	if cause == nil {
		// Allow calling PropagateWithCode without checking whether there is error
		return nil
	}
	return create(cause, code, msg, vals...)
}

func GetCode(err error) ErrorCode {
	if err, ok := err.(*stacktrace); ok {
		return err.code
	}
	return NoCode
}

type stacktrace struct {
	message  string
	cause    error
	code     ErrorCode
	file     string
	function string
	line     int
}

func create(cause error, code ErrorCode, msg string, vals ...interface{}) error {
	// If no error code specified, inherit error code from the cause.
	if code == NoCode {
		code = GetCode(cause)
	}

	err := &stacktrace{
		message: fmt.Sprintf(msg, vals...),
		cause:   cause,
		code:    code,
	}

	// Caller of create is NewError or Propagate, so user's code is 2 up.
	pc, file, line, ok := runtime.Caller(2)
	if !ok {
		return err
	}
	if CleanPath != nil {
		file = CleanPath(file)
	}
	err.file, err.line = file, line

	f := runtime.FuncForPC(pc)
	if f == nil {
		return err
	}
	err.function = shortFuncName(f)

	return err
}

/* "FuncName" or "Receiver.MethodName" */
func shortFuncName(f *runtime.Func) string {
	// f.Name() is like one of these:
	// - "github.com/palantir/shield/package.FuncName"
	// - "github.com/palantir/shield/package.Receiver.MethodName"
	// - "github.com/palantir/shield/package.(*PtrReceiver).MethodName"
	longName := f.Name()

	withoutPath := longName[strings.LastIndex(longName, "/")+1:]
	withoutPackage := withoutPath[strings.Index(withoutPath, ".")+1:]

	shortName := withoutPackage
	shortName = strings.Replace(shortName, "(", "", 1)
	shortName = strings.Replace(shortName, "*", "", 1)
	shortName = strings.Replace(shortName, ")", "", 1)

	return shortName
}

func (st *stacktrace) Error() string {
	return fmt.Sprint(st)
}

var DefaultFormat = FormatFull

// Format is the type of the two possible values of stacktrace.DefaultFormat.
type Format int

const (
	// FormatFull means format as a full stacktrace including line number information.
	FormatFull Format = iota
	// FormatBrief means Format on a single line without line number information.
	FormatBrief
)

var _ fmt.Formatter = (*stacktrace)(nil)

func (st *stacktrace) Format(f fmt.State, c rune) {
	var text string
	if f.Flag('+') && !f.Flag('#') && c == 's' { // "%+s"
		text = formatFull(st)
	} else if f.Flag('#') && !f.Flag('+') && c == 's' { // "%#s"
		text = formatBrief(st)
	} else {
		text = map[Format]func(*stacktrace) string{
			FormatFull:  formatFull,
			FormatBrief: formatBrief,
		}[DefaultFormat](st)
	}

	formatString := "%"
	// keep the flags recognized by fmt package
	for _, flag := range "-+# 0" {
		if f.Flag(int(flag)) {
			formatString += string(flag)
		}
	}
	if width, has := f.Width(); has {
		formatString += fmt.Sprint(width)
	}
	if precision, has := f.Precision(); has {
		formatString += "."
		formatString += fmt.Sprint(precision)
	}
	formatString += string(c)
	fmt.Fprintf(f, formatString, text)
}

func formatFull(st *stacktrace) string {
	var str string
	newline := func() {
		if str != "" && !strings.HasSuffix(str, "\n") {
			str += "\n"
		}
	}

	for curr, ok := st, true; ok; curr, ok = curr.cause.(*stacktrace) {
		str += curr.message

		if curr.file != "" {
			newline()
			if curr.function == "" {
				str += fmt.Sprintf(" --- at %v:%v ---", curr.file, curr.line)
			} else {
				str += fmt.Sprintf(" --- at %v:%v (%v) ---", curr.file, curr.line, curr.function)
			}
		}

		if curr.cause != nil {
			newline()
			if cause, ok := curr.cause.(*stacktrace); !ok {
				str += "Caused by: "
				str += curr.cause.Error()
			} else if cause.message != "" {
				str += "Caused by: "
			}
		}
	}

	return str
}

func formatBrief(st *stacktrace) string {
	var str string
	concat := func(msg string) {
		if str != "" && msg != "" {
			str += ": "
		}
		str += msg
	}

	curr := st
	for {
		concat(curr.message)
		if cause, ok := curr.cause.(*stacktrace); ok {
			curr = cause
		} else {
			break
		}
	}
	if curr.cause != nil {
		concat(curr.cause.Error())
	}
	return str
}

func RootCause(err error) error {
	for {
		//fmt.Println("/////////RootCause")

		st, ok := err.(*stacktrace)
		//fmt.Println("/////////", st.message, "///", ok, "//////////")
		if !ok {
			return err
		}
		//fmt.Println("/////////", st.cause, "//////////")
		if st.cause == nil {
			// return Errorf(st.message)
			return errors.New(st.message)
		}
		err = st.cause
	}
}
