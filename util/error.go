package util

import (
	"fmt"
	"io"
	"os"
	"strings"
)

const (
	// DefaultErrorExitCode is 1 for exit code
	DefaultErrorExitCode = 1
)

var fatalErrHandler = fatal

// fatal prints the message (if provided) and then exits.
func fatal(output io.Writer, msg string, code int) {
	if len(msg) > 0 {
		// add newline if needed
		if !strings.HasSuffix(msg, "\n") {
			msg += "\n"
		}
		fmt.Fprint(output, msg)
	}
	os.Exit(code)
}

// ExitWithMsg prints a user friendly error and exits with a non-zero exit code.
func ExitWithMsg(output io.Writer, format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	fatal(output, msg, DefaultErrorExitCode)
}

// CheckErr prints a user friendly error and exits with a non-zero
// exit code. Unrecognized errors will be printed with an "error: " prefix.
func CheckErr(output io.Writer, err error) {
	checkErr(output, err, fatalErrHandler)
}

// checkErr formats a given error as a string and calls the passed handleErr
// func with that string and an exit code.
func checkErr(output io.Writer, err error, handleErr func(io.Writer, string, int)) {
	if err == nil {
		return
	}
	msg := err.Error()
	if !strings.HasPrefix(msg, "error: ") {
		msg = fmt.Sprintf("error: %s", msg)
	}
	handleErr(output, msg, DefaultErrorExitCode)
}
