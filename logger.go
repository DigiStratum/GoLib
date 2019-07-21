// DigiStratum GoLib - Logger
package golib

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
type LogWriter interface {
	Log(message string)
}

// DStdOut LogWriter
type stdOutLogWriter struct {}

// Make a new LogWriter for StdOut
func NewStdOutLogWriter() LogWriter {
	return stdOutLogWriter{}
}

//  Satisfies LogWriter interface
func (lw stdOutLogWriter) Log(message string) {
	fmt.Println(message)
}

// Default LogWriter
type defaultLogWriter struct {}

// Make a new LogWriter for default golang logger
func NewDefaultLogWriter() LogWriter {
	return defaultLogWriter{}
}

//  Satisfies LogWriter interface
func (lw defaultLogWriter) Log(message string) {
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

type logger struct {
	threadId	string		// Quasi-distinct threadId to filter log output by thread
	minLogLevel	logLevel	// The minimum logging level
	logWriter	LogWriter	// The LogWriter we are going to use (TODO: Add support for multiple)
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

// Get a new instance
func NewLogger() *logger {
	// We would be galactically unlucky to get two threads that start at the same nano-second...
	newLogger := logger {
		threadId:	fmt.Sprintf("%d", time.Now().UTC().UnixNano()),
		minLogLevel:	INFO,
		logWriter:	NewStdOutLogWriter(),
	}
	return &newLogger
}

// Set the threadId to something meaningful to the caller
func (l *logger) SetThreadId(threadId string) {
	l.threadId = threadId
}

// Set the minimum log level
func (l *logger) SetMinLogLevel(level string) {
	switch (strings.ToLower(strings.TrimSpace(level))) {
		case "crazy": l.minLogLevel = CRAZY; return
		case "trace": l.minLogLevel = TRACE; return
		case "debug": l.minLogLevel = DEBUG; return
		case "info": l.minLogLevel = INFO; return
		case "warn": l.minLogLevel = WARN; return
		case "error": l.minLogLevel = ERROR; return
		case "fatal": l.minLogLevel = FATAL; return
	}
	l.Error(fmt.Sprintf("Logger: SetMinLogLevel(): Unrecognized Log Level requested: '%s'", level))
}

// Replace the current LogWriter with something more to our liking
// TODO: Add support for multiple LogWriter's so that we can send logs to more than one place
func (l *logger) SetLogWriter(logWriter LogWriter) {
	l.logWriter = logWriter
}

// Log some output
// Wrap level+msg in an error as a code-reduction convenience to any caller wanting to return it
func (l *logger) log(level logLevel, msg string) error {
	var prefix string
	switch (level) {
		case CRAZY: prefix = "CRAZY"
		case TRACE: prefix = "TRACE"
		case DEBUG: prefix = "DEBUG"
		case INFO: prefix = "INFO"
		case WARN: prefix = "WARN"
		case ERROR: prefix = "ERROR"
		case FATAL: prefix = "FATAL"
	}
	logMsg := fmt.Sprintf("%5s %s", prefix, msg)
	if l.WillLogForLevel(level) {
		// Send the log message to our LogWriter
		t := time.Now()
		l.logWriter.Log(fmt.Sprintf(
			"%s thread:%s %s",
			t.Format(time.RFC3339),
			l.threadId,
			logMsg,
		))
	}
	return errors.New(logMsg)
}

// Will this logger log for the specified level?
func (l *logger) WillLogForLevel(lvl logLevel) bool {
	return lvl >= l.minLogLevel;
}


// Log CRAZY output
func (l *logger) Crazy(msg string) error {
	return l.log(CRAZY, msg)
}

// Log TRACE output
func (l *logger) Trace(msg string) error {
	return l.log(TRACE, msg)
}

// Log DEBUG output
func (l *logger) Debug(msg string) error {
	return l.log(DEBUG, msg)
}

// Log INFO output
func (l *logger) Info(msg string) error {
	return l.log(INFO, msg)
}

// Log WARN output
func (l *logger) Warn(msg string) error {
	return l.log(WARN, msg)
}

// Log ERROR output
func (l *logger) Error(msg string) error {
	return l.log(ERROR, msg)
}

// Log FATAL output (caller should exit/panic after this)
func (l *logger) Fatal(msg string) error {
	return l.log(FATAL, msg)
}

