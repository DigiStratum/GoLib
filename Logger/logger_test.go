package logger

import(
	"fmt"
	"testing"

	. "github.com/DigiStratum/GoLib/Testing"
)

var LastMessage string
type mockLogWriter struct { }

func (r mockLogWriter) Log(format string, a ...interface{}) {
	LastMessage = fmt.Sprintf(format, a...)
}

func TestThat_GetLogger_ReturnsSomething(t *testing.T) {
	// Test
	var sut *Logger = GetLogger()

	// Verify
	ExpectNonNil(sut, t)
}

func TestThat_NewLogger_ReturnsSomething(t *testing.T) {
	// Test
	var sut *Logger = NewLogger("newstream")

	// Verify
	ExpectNonNil(sut, t)
}

func TestThat_Logger_SetLogWriter_ReplacesStdOutWithMock(t *testing.T) {
	// Setup
	var sut *Logger = GetLogger()
	mockWriter := mockLogWriter{}
	expectedMessage := "test message"

	// Test
	sut.SetLogWriter(mockWriter)
	sut.Error("%s", expectedMessage)

	// Verify
	// Actual: 2022-05-19T08:03:24-07:00 thread:1652972604971037495 ERROR test message
	ExpectMatch("^\\d{4}-\\d{2}-\\d{2}T\\d{2}:\\d{2}:\\d{2}-\\d{2}:\\d{2}\\s*thread:\\d+\\s+ERROR test message$", LastMessage, t)
}

func TestThat_Logger_LogTimestamp_EliminatesTimestampFromMessages(t *testing.T) {
	// Setup
	var sut *Logger = GetLogger()
	mockWriter := mockLogWriter{}
	expectedMessage := "test message"
	sut.SetLogWriter(mockWriter)
	sut.LogTimestamp(false)

	// Test
	sut.Error("%s", expectedMessage)

	// Verify
	// Actual: thread:1652972604971037495 ERROR test message
	ExpectMatch("^thread:\\d+\\s+ERROR test message$", LastMessage, t)
}


func TestThat_Logger_DefaultMinLogLevel_SuppressesLogLevelsBelowDefault(t *testing.T) {
	// Setup
	var sut *Logger = getMockedLogger()
	expectedMessage := "test message"

	// Test
	sut.Debug("%s", expectedMessage)
	actualMessageBelow := LastMessage
	sut.Info("%s", expectedMessage)
	actualMessageAt := LastMessage
	sut.Error("%s", expectedMessage)
	actualMessageAbove := LastMessage

	// Verify
	ExpectEmptyString(actualMessageBelow, t)
	ExpectNonEmptyString(actualMessageAt, t)
	ExpectNonEmptyString(actualMessageAbove, t)
	// Actual: 2022-05-19T08:03:24-07:00 thread:1652972604971037495 ERROR|INFO test message
	ExpectMatch("^.*INFO test message$", actualMessageAt, t)
	ExpectMatch("^.*ERROR test message$", actualMessageAbove, t)
}

func TestThat_Logger_SetMinLogLevel_PassessAllLogLevels_WhenAtLowestSetting(t *testing.T) {
	// Setup
	var sut *Logger = getMockedLogger()
	expectedMessage := "test message"

	// Test / Verify
	sut.SetMinLogLevel(CRAZY)
	sut.Fatal("%s", expectedMessage)
	ExpectNonEmptyString(LastMessage, t)
	sut.Error("%s", expectedMessage)
	ExpectNonEmptyString(LastMessage, t)
	sut.Warn("%s", expectedMessage)
	ExpectNonEmptyString(LastMessage, t)
	sut.Info("%s", expectedMessage)
	ExpectNonEmptyString(LastMessage, t)
	sut.Debug("%s", expectedMessage)
	ExpectNonEmptyString(LastMessage, t)
	sut.Trace("%s", expectedMessage)
	ExpectNonEmptyString(LastMessage, t)
	sut.Crazy("%s", expectedMessage)
	ExpectNonEmptyString(LastMessage, t)
}

func TestThat_Logger_SetMinLogLevel_SuppressesLowerLogLevels(t *testing.T) {
	// Setup
	var sut *Logger = getMockedLogger()
	expectedMessage := "test message"

	// Test / Verify
	sut.SetMinLogLevel(FATAL)
	sut.Error("%s", expectedMessage)
	ExpectEmptyString(LastMessage, t)
	sut.Warn("%s", expectedMessage)
	ExpectEmptyString(LastMessage, t)
	sut.Info("%s", expectedMessage)
	ExpectEmptyString(LastMessage, t)
	sut.Debug("%s", expectedMessage)
	ExpectEmptyString(LastMessage, t)
	sut.Trace("%s", expectedMessage)
	ExpectEmptyString(LastMessage, t)
	sut.Crazy("%s", expectedMessage)
	ExpectEmptyString(LastMessage, t)
	sut.Fatal("%s", expectedMessage)
	ExpectNonEmptyString(LastMessage, t)
}

func TestThat_Logger_ErrorsReturnedOverWarnLevel(t *testing.T) {
	// Setup
	var sut *Logger = getMockedLogger()
	expectedMessage := "test message"
	var err error

	// Test / Verify
	err = sut.Fatal("%s", expectedMessage)
	ExpectError(err, t)
	err = sut.Error("%s", expectedMessage)
	ExpectError(err, t)
	err = sut.Warn("%s", expectedMessage)
	ExpectError(err, t)
	err = sut.Info("%s", expectedMessage)
	ExpectNoError(err, t)
	err = sut.Debug("%s", expectedMessage)
	ExpectNoError(err, t)
	err = sut.Trace("%s", expectedMessage)
	ExpectNoError(err, t)
	err = sut.Crazy("%s", expectedMessage)
	ExpectNoError(err, t)
}

func getMockedLogger() *Logger {
	var sut *Logger = GetLogger()
	sut.SetLogWriter(mockLogWriter{})
	LastMessage = ""
	return sut
}

