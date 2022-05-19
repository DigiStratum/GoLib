package logwriter

import (
	"fmt"
)

// StdOut LogWriter
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

func (r StdOutLogWriter) Log(message string) {
	fmt.Println(message)
}
