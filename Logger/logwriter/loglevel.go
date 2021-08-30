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

func init() {
	logLevels := make(map[LogLevel]string)
	logLevels[CRAZY] = "CRAZY"
	logLevels[TRACE] = "TRACE"
	logLevels[DEBUG] = "DEBUG"
	logLevels[INFO] = "INFO"
	logLevels[WARN] = "WARN"
	logLevels[ERROR] = "ERROR"
	logLevels[FATAL] = "FATAL"
}

// -------------------------------------------------------------------------------------------------
// Factory Functions
// -------------------------------------------------------------------------------------------------

func StringToLogLevel(logLevelStr string) (*LogLevel, error) {
	switch logLevelStr {
		case "CRAZY": return ll := CRAZY; return &ll, nil
		case "TRACE": return ll := TRACE; return &ll, nil
		case "DEBUG": return ll := DEBUG; return &ll, nil
		case "INFO": return ll := INFO; return &ll, nil
		case "WARN": return ll := WARN; return &ll, nil
		case "ERROR": return ll := ERROR; return &ll, nil
		case "FATAL": return ll := FATAL; return &ll, nil
	}
	return nil, fmt.Errorf("Specifier is not a valid LogLevel [%s]", logLevelStr)
}

// -------------------------------------------------------------------------------------------------
// LogLevelIfc Public Interface
// -------------------------------------------------------------------------------------------------

func (r LogLevel) ToString() string {
	if r >= logLevelEnd { return "" }
	return logLevels[r]
}
