package logger

import(
	"testing"

	. "github.com/DigiStratum/GoLib/Testing"
)

func TestThat_LogLevel_StringCoversionsMatch(t *testing.T) {
	// Setup
	var actualCrazy, actualTrace, actualDebug, actualInfo, actualWarn, actualError, actualFatal *LogLevel
	var errCrazy, errTrace, errDebug, errInfo, errWarn, errError, errFatal error

	// Test
	actualCrazy, errCrazy = StringToLogLevel("CRAZY")
	actualTrace, errTrace = StringToLogLevel("TRACE")
	actualDebug, errDebug = StringToLogLevel("DEBUG")
	actualInfo, errInfo = StringToLogLevel("INFO")
	actualWarn, errWarn = StringToLogLevel("WARN")
	actualError, errError = StringToLogLevel("ERROR")
	actualFatal, errFatal = StringToLogLevel("FATAL")

	// Verify
	ExpectNonNil(actualCrazy, t)
	ExpectNoError(errCrazy, t)
	ExpectString("CRAZY", actualCrazy.ToString(), t)

	ExpectNonNil(actualTrace, t)
	ExpectNoError(errTrace, t)
	ExpectString("TRACE", actualTrace.ToString(), t)
	ExpectNonNil(actualDebug, t)
	ExpectNoError(errDebug, t)
	ExpectString("DEBUG", actualDebug.ToString(), t)
	ExpectNonNil(actualInfo, t)
	ExpectNoError(errInfo, t)
	ExpectString("INFO", actualInfo.ToString(), t)
	ExpectNonNil(actualWarn, t)
	ExpectNoError(errWarn, t)
	ExpectString("WARN", actualWarn.ToString(), t)
	ExpectNonNil(actualError, t)
	ExpectNoError(errError, t)
	ExpectString("ERROR", actualError.ToString(), t)
	ExpectNonNil(actualFatal, t)
	ExpectNoError(errFatal, t)
	ExpectString("FATAL", actualFatal.ToString(), t)
}

