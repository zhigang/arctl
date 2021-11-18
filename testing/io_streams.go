package testing

import (
	"bytes"

	"github.com/zhigang/arctl/printers"
)

// NewTestIOStreams returns a valid IOStreams and in, out, errout buffers for unit tests
func NewTestIOStreams() (printers.IOStreams, *bytes.Buffer, *bytes.Buffer, *bytes.Buffer) {
	in := &bytes.Buffer{}
	out := &bytes.Buffer{}
	errOut := &bytes.Buffer{}

	return printers.IOStreams{
		In:     in,
		Out:    out,
		ErrOut: errOut,
	}, in, out, errOut
}
