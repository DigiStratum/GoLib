package logwriter

// Default LogWriter
type DefaultLogWriter struct {}

// -------------------------------------------------------------------------------------------------
// Factory Functions
// -------------------------------------------------------------------------------------------------

// Make a new LogWriter for default golang logger
func NewDefaultLogWriter() DefaultLogWriter {
	return DefaultLogWriter{}
}

// -------------------------------------------------------------------------------------------------
// LogWriterIfc Public Interface
// -------------------------------------------------------------------------------------------------

func (r DefaultLogWriter) Log(message string) {
	log.Println(message)
}
