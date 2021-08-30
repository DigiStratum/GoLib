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

TODO: Add support to connect log output to a file, database, or API (event stream), etc

*/

import (
	"fmt"
	"log"
	"time"
	"strings"
	"errors"
)

// LogWriter interface
type LogWriterIfc interface {
	Log(message string)
}

// StdOut LogWriter
type stdOutLogWriter struct {}

// Default LogWriter
type defaultLogWriter struct {}

// Factory Functions

// Make a new LogWriter for StdOut
func NewStdOutLogWriter() stdOutLogWriter {
	return stdOutLogWriter{}
}

// Make a new LogWriter for default golang logger
func NewDefaultLogWriter() defaultLogWriter {
	return defaultLogWriter{}
}

// -------------------------------------------------------------------------------------------------
// LogWriterIfc Public Interface
// -------------------------------------------------------------------------------------------------

func (r stdOutLogWriter) Log(message string) {
	fmt.Println(message)
}

func (r defaultLogWriter) Log(message string) {
	log.Println(message)
}

type logLevel uint

// Log levels
const (
	CRAZY logLevel = iota	// Crazy output: data structures, dumps, ASCII art, you name it
	TRACE			// Where in the code base are we, and how were we called?
	DEBUG			// What is our state and other helpful things for trouble shooting?
	INFO			// What functional contract operation is running?
	WARN			// What possible problem do we see that may need a human response?
	ERROR			// What definite problem is there that will degrade functionality/performance?
	FATAL			// What fundamental problem is there that is considered do or die?
)

type LoggerIfc interface {
}

type logger struct {
	threadId	string			// Quasi-distinct threadId to filter log output by thread
	minLogLevel	logLevel		// The minimum logging level
	logWriter	LogWriterIfc		// The LogWriter we are going to use (TODO: Add support for multiple)
	logTimestamp	bool			// Conditionally disable timestamps on the log output (consumer may do this for us)
	logLevelLabels	map[logLevel]string	// Convert a given logLevel to a readable string
}

var loggerInstance logger

// Automagically set up our singleton
func init() {
	loggerInstance = *NewLogger()
}

// Get our singleton instance
func GetLogger() *logger {
	return &loggerInstance
}

// Get a new (non-singleton)instance
func NewLogger() *logger {
	logLevelLabels := make(map[logLevel]string)
	logLevelLabels[CRAZY] = "CRAZY"
	logLevelLabels[TRACE] = "TRACE"
	logLevelLabels[DEBUG] = "DEBUG"
	logLevelLabels[INFO] = "INFO"
	logLevelLabels[WARN] = "WARN"
	logLevelLabels[ERROR] = "ERROR"
	logLevelLabels[FATAL] = "FATAL"
	newLogger := logger {
		threadId:	fmt.Sprintf("%d", time.Now().UTC().UnixNano()),
		minLogLevel:	INFO,
		logWriter:	NewStdOutLogWriter(),
		logTimestamp:	true,
		logLevelLabels:	logLevelLabels,
	}
	return &newLogger
}

// -------------------------------------------------------------------------------------------------
// LoggerIfc Public Interface
// -------------------------------------------------------------------------------------------------

// Set the threadId to something meaningful to the caller
func (r *logger) SetThreadId(threadId string) {
	r.threadId = threadId
}

// Set the minimum log level
func (r *logger) SetMinLogLevel(level string) {
	switch (strings.ToLower(strings.TrimSpace(level))) {
		case "crazy": r.minLogLevel = CRAZY; return
		case "trace": r.minLogLevel = TRACE; return
		case "debug": r.minLogLevel = DEBUG; return
		case "info": r.minLogLevel = INFO; return
		case "warn": r.minLogLevel = WARN; return
		case "error": r.minLogLevel = ERROR; return
		case "fatal": r.minLogLevel = FATAL; return
	}
	r.Error(fmt.Sprintf("Logger: SetMinLogLevel(): Unrecognized Log Level requested: '%s'", level))
}

// Replace the current LogWriter with something more to our liking
// TODO: Add support for multiple LogWriter's so that we can send logs to more than one place
func (r *logger) SetLogWriter(logWriter LogWriterIfc) {
	r.logWriter = logWriter
}

// Set the logTimestamp state (defaults to true to enable timestamps in logger output)
func (r *logger) LogTimestamp(logTimestamp bool) {
	r.logTimestamp = logTimestamp
}

// Log some output
func (r logger) log(level logLevel, format string, a ...interface{}) error {
	msg := fmt.Sprintf(format, a...)
	logMsg := fmt.Sprintf("%5s %s", r.logLevelLabels[level], msg)
	if level >= r.minLogLevel {
		// Send the log message to our LogWriter
		timestamp := ""
		if r.logTimestamp {
			timestamp = fmt.Sprintf("%s ", time.Now().Format(time.RFC3339))
		}
		r.logWriter.Log(fmt.Sprintf(
			"%sthread:%s %s",
			timestamp,
			r.threadId,
			logMsg,
		))
	}
	// Wrap level (WARN|ERROR|FATAL)+message in an error as a code
	// reduction convenience to any caller wanting to return it
	if level >= WARN { return errors.New(logMsg) }
	return nil
}

// Log CRAZY output
func (r logger) Crazy(format string, a ...interface{}) error {
	return r.log(CRAZY, format, a...)
}

// Log TRACE output
func (r logger) Trace(format string, a ...interface{}) error {
	return r.log(TRACE, format, a...)
}

// Log DEBUG output
func (r logger) Debug(format string, a ...interface{}) error {
	return r.log(DEBUG, format, a...)
}

// Log INFO output
func (r logger) Info(format string, a ...interface{}) error {
	return r.log(INFO, format, a...)
}

// Log WARN output
func (r logger) Warn(format string, a ...interface{}) error {
	return r.log(WARN, format, a...)
}

// Log ERROR output
func (r logger) Error(format string, a ...interface{}) error {
	return r.log(ERROR, format, a...)
}

// Log FATAL output (caller should exit/panic after this)
func (r logger) Fatal(format string, a ...interface{}) error {
	return r.log(FATAL, format, a...)
}
