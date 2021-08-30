// DigiStratum GoLib - Logger
package logger

/*

This simple logging class provides a standardized interface to produce log output using log levels
as a means of filtering what types of log output get produced. The log level may be changed at any
time and is soft-configurable. Currently the output is limited to StdOut, however the intent is to
provide a connector interface to redirect output to one or more other locations.

There are two ways to get a Logger instance, one returns our own singleton, the other returns a new
instance. This way you may use the singleton throughout your application without reinitializing or
passing it around all over the place, or create a new Logger with a separate configuration and do
just that, as needed.

TODO:
 * Add support to connect log output to a file, database, or API (event stream), etc

*/

import (
	"fmt"
	"log"
	"time"
	"strings"
	"errors"
)

type LoggerIfc interface {
	SetMinLogLevel(minLogLevel LogLevel)
	SetLogWriter(logWriter LogWriterIfc)
	LogTimestamp(logTimestamp bool)
	Crazy(format string, a ...interface{}) error
	Trace(format string, a ...interface{}) error
	Debug(format string, a ...interface{}) error
	Info(format string, a ...interface{}) error
	Warn(format string, a ...interface{}) error
	Error(format string, a ...interface{}) error
	Fatal(format string, a ...interface{}) error
}

type Logger struct {
	streamId	string			// Quasi-distinct streamId to filter log output by thread
	minLogLevel	LogLevel		// The minimum logging level
	logWriter	LogWriterIfc		// The LogWriter we are going to use (TODO: Add support for multiple)
	logTimestamp	bool			// Conditionally disable timestamps on the log output (consumer may do this for us)
}

// -------------------------------------------------------------------------------------------------
// Initialization - Singleton
// -------------------------------------------------------------------------------------------------

var loggerInstance Logger

// Automagically set up our singleton
func init() {
	// Default log streamId is our instantiation timestamp
	streamId := fmt.Sprintf("%d", time.Now().UTC().UnixNano())
	loggerInstance = *NewLogger(streamId)
	LogTimestamp(logTimestamp bool)
}

// Get our singleton instance
func GetLogger() *Logger {
	return &loggerInstance
}

// -------------------------------------------------------------------------------------------------
// Factory Functions
// -------------------------------------------------------------------------------------------------

// Get a new (non-singleton)instance with a streamId meaningful to the caller (such as its runtime thread id)
func NewLogger(streamId string) *Logger {
	newLogger := Logger{
		streamId:	streamId,
		minLogLevel:	INFO,
		logWriter:	NewStdOutLogWriter(),
		logTimestamp:	true,
	}
	return &newLogger
}

// -------------------------------------------------------------------------------------------------
// LoggerIfc Public Interface
// -------------------------------------------------------------------------------------------------

// Set the minimum log level
func (r *Logger) SetMinLogLevel(minLogLevel LogLevel) {
	r.minLogLevel = minLogLevel
}

// Replace the current LogWriter with something more to our liking
// TODO: Add support for multiple LogWriter's so that we can send logs to more than one place
func (r *Logger) SetLogWriter(logWriter LogWriterIfc) {
	r.logWriter = logWriter
}

// Set the logTimestamp state (defaults to true to enable timestamps in logger output)
func (r *Logger) LogTimestamp(logTimestamp bool) {
	r.logTimestamp = logTimestamp
}

// Log CRAZY output
func (r Logger) Crazy(format string, a ...interface{}) error {
	return r.log(CRAZY, format, a...)
}

// Log TRACE output
func (r Logger) Trace(format string, a ...interface{}) error {
	return r.log(TRACE, format, a...)
}

// Log DEBUG output
func (r Logger) Debug(format string, a ...interface{}) error {
	return r.log(DEBUG, format, a...)
}

// Log INFO output
func (r Logger) Info(format string, a ...interface{}) error {
	return r.log(INFO, format, a...)
}

// Log WARN output
func (r Logger) Warn(format string, a ...interface{}) error {
	return r.log(WARN, format, a...)
}

// Log ERROR output
func (r Logger) Error(format string, a ...interface{}) error {
	return r.log(ERROR, format, a...)
}

// Log FATAL output (caller should exit/panic after this)
func (r Logger) Fatal(format string, a ...interface{}) error {
	return r.log(FATAL, format, a...)
}

// -------------------------------------------------------------------------------------------------
// Logger Private Interface
// -------------------------------------------------------------------------------------------------

// Log some output
func (r Logger) log(level LogLevel, format string, a ...interface{}) error {
	msg := fmt.Sprintf(format, a...)
	logMsg := fmt.Sprintf("%5s %s", level.ToString(), msg)
	if level >= r.minLogLevel {
		// Send the log message to our LogWriter
		timestamp := ""
		if r.logTimestamp {
			timestamp = fmt.Sprintf("%s ", time.Now().Format(time.RFC3339))
		}
		r.logWriter.Log(fmt.Sprintf(
			"%sthread:%s %s",
			timestamp,
			r.streamId,
			logMsg,
		))
	}
	// Wrap level (WARN|ERROR|FATAL)+message in an error as a code
	// reduction convenience to any caller wanting to return it
	if level >= WARN { return errors.New(logMsg) }
	return nil
}
