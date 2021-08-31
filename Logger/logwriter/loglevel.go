package logger

type LogLevelIfc interface {
	ToString() string
}

type LogLevel uint

// Log levels
const (
	CRAZY LogLevel = iota	// Crazy output: data structures, dumps, ASCII art, you name it
	TRACE			// Where in the code base are we, and how were we called?
	DEBUG			// What is our state and other helpful things for trouble shooting?
	INFO			// What functional contract operation is running?
	WARN			// What possible problem do we see that may need a human response?
	ERROR			// What definite problem is there that will degrade functionality/performance?
	FATAL			// What fundamental problem is there that is considered do or die?
	logLevelEnd
)

// -------------------------------------------------------------------------------------------------
// Initialization
// -------------------------------------------------------------------------------------------------

var logLevels map[LogLevel]string
var logLabels map[string]LogLevel

func init() {
	logLevels := make(map[LogLevel]string)
	logLevels[CRAZY] = "CRAZY"
	logLevels[TRACE] = "TRACE"
	logLevels[DEBUG] = "DEBUG"
	logLevels[INFO] = "INFO"
	logLevels[WARN] = "WARN"
	logLevels[ERROR] = "ERROR"
	logLevels[FATAL] = "FATAL"

	logLabels := make(map[string]LogLevel])
	logLabels["CRAZY"] = CRAZY
	logLabels["TRACE"] = TRACE
	logLabels["DEBUG"] = DEBUG
	logLabels["INFO"] = INFO
	logLabels["WARN"] = WARN
	logLabels["ERROR"] = ERROR
	logLabels["FATAL"] = FATAL

}

// -------------------------------------------------------------------------------------------------
// Factory Functions
// -------------------------------------------------------------------------------------------------

func StringToLogLevel(logLevelStr string) (*LogLevel, error) {
	if logLevel, ok := logLabels[logLevelStr]; ok { return &logLevel, nil }
	return nil, fmt.Errorf("Specifier is not a valid LogLevel [%s]", logLevelStr)
}

// -------------------------------------------------------------------------------------------------
// LogLevelIfc Public Interface
// -------------------------------------------------------------------------------------------------

func (r LogLevel) ToString() string {
	if r >= logLevelEnd { return "" }
	return logLevels[r]
}
