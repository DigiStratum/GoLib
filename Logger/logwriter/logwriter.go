package logwriter

type LogWriterIfc interface {
	Log(format string, a ...interface{})
}
