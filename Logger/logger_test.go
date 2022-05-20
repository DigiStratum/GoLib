package logger

import(
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

// SetLogWriter(logWriter lw.LogWriterIfc)
func TestThat_SetLogWriter_ReplacesStdOutWithMock(t *testing.T) {
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
	ExpectMatch("^.*ERROR test message$", *LastMessage, t)
}

