package logwriter

// StdOut LogWriter

import (
	"fmt"
)

type StdOutLogWriter struct {}

// -------------------------------------------------------------------------------------------------
// Factory Functions
// -------------------------------------------------------------------------------------------------

// Make a new LogWriter for StdOut
func NewStdOutLogWriter() *StdOutLogWriter {
	return &StdOutLogWriter{}
}

// -------------------------------------------------------------------------------------------------
// LogWriterIfc Public Interface
// -------------------------------------------------------------------------------------------------

func (r StdOutLogWriter) Log(format string, a ...interface{}) {
	fmt.Println(fmt.Sprintf(format, a...))
}

