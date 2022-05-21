package logger

import(
//	"fmt"
	"testing"

	. "github.com/DigiStratum/GoLib/Testing"
)

var LastMessage *string
type mockLogWriter struct { }

func (r mockLogWriter) Log(message string) {
	LastMessage = &message
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
	sut.Error(expectedMessage)

	// Verify
	ExpectNonNil(LastMessage, t)
	// Actual: 2022-05-19T08:03:24-07:00 thread:1652972604971037495 ERROR test message
	ExpectMatch("^\\d{4}-\\d{2}-\\d{2}T\\d{2}:\\d{2}:\\d{2}-\\d{2}:\\d{2}\\s*thread:\\d+\\s+ERROR test message$", *LastMessage, t)
}

func TestThat_Logger_LogTimestamp_EliminatesTimestampFromMessages(t *testing.T) {
	// Setup
	var sut *Logger = GetLogger()
	mockWriter := mockLogWriter{}
	expectedMessage := "test message"
	sut.SetLogWriter(mockWriter)
	sut.LogTimestamp(false)

	// Test
	sut.Error(expectedMessage)

	// Verify
	ExpectNonNil(LastMessage, t)
	// Actual: thread:1652972604971037495 ERROR test message
	ExpectMatch("^thread:\\d+\\s+ERROR test message$", *LastMessage, t)
}


func TestThat_Logger_DefaultMinLogLevel_SuppressesLogLevelsBelowDefault(t *testing.T) {
	// Setup
	var sut *Logger = getMockedLogger()
	expectedMessage := "test message"

	// Test
	sut.Debug(expectedMessage)
	actualMessageBelow := LastMessage
	sut.Info(expectedMessage)
	actualMessageAt := LastMessage
	sut.Error(expectedMessage)
	actualMessageAbove := LastMessage

	// Verify
	ExpectNil(actualMessageBelow, t)
	ExpectNonNil(actualMessageAt, t)
	ExpectNonNil(actualMessageAbove, t)
	// Actual: 2022-05-19T08:03:24-07:00 thread:1652972604971037495 ERROR|INFO test message
	ExpectMatch("^.*INFO test message$", *actualMessageAt, t)
	ExpectMatch("^.*ERROR test message$", *actualMessageAbove, t)
}

func TestThat_Logger_SetMinLogLevel_PassessAllLogLevels_WhenAtLowestSetting(t *testing.T) {
	// Setup
	var sut *Logger = getMockedLogger()
	expectedMessage := "test message"

	// Test / Verify
	sut.SetMinLogLevel(CRAZY)
	sut.Fatal(expectedMessage)
	ExpectNonNil(LastMessage, t)
	sut.Error(expectedMessage)
	ExpectNonNil(LastMessage, t)
	sut.Warn(expectedMessage)
	ExpectNonNil(LastMessage, t)
	sut.Info(expectedMessage)
	ExpectNonNil(LastMessage, t)
	sut.Debug(expectedMessage)
	ExpectNonNil(LastMessage, t)
	sut.Trace(expectedMessage)
	ExpectNonNil(LastMessage, t)
	sut.Crazy(expectedMessage)
	ExpectNonNil(LastMessage, t)
}

func TestThat_Logger_SetMinLogLevel_SuppressesLowerLogLevels(t *testing.T) {
	// Setup
	var sut *Logger = getMockedLogger()
	expectedMessage := "test message"

	// Test / Verify
	sut.SetMinLogLevel(FATAL)
	sut.Error(expectedMessage)
	ExpectNil(LastMessage, t)
	sut.Warn(expectedMessage)
	ExpectNil(LastMessage, t)
	sut.Info(expectedMessage)
	ExpectNil(LastMessage, t)
	sut.Debug(expectedMessage)
	ExpectNil(LastMessage, t)
	sut.Trace(expectedMessage)
	ExpectNil(LastMessage, t)
	sut.Crazy(expectedMessage)
	ExpectNil(LastMessage, t)
	sut.Fatal(expectedMessage)
	ExpectNonNil(LastMessage, t)
}

func TestThat_Logger_ErrorsReturnedOverWarnLevel(t *testing.T) {
	// Setup
	var sut *Logger = getMockedLogger()
	expectedMessage := "test message"
	var err error

	// Test / Verify
	err = sut.Fatal(expectedMessage)
	ExpectError(err, t)
	err = sut.Error(expectedMessage)
	ExpectError(err, t)
	err = sut.Warn(expectedMessage)
	ExpectError(err, t)
	err = sut.Info(expectedMessage)
	ExpectNoError(err, t)
	err = sut.Debug(expectedMessage)
	ExpectNoError(err, t)
	err = sut.Trace(expectedMessage)
	ExpectNoError(err, t)
	err = sut.Crazy(expectedMessage)
	ExpectNoError(err, t)
}

func getMockedLogger() *Logger {
	var sut *Logger = GetLogger()
	sut.SetLogWriter(mockLogWriter{})
	LastMessage = nil
	return sut
}

