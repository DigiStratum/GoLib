package logwriter

import(
	"os"
	"fmt"
	"testing"

	. "github.com/DigiStratum/GoLib/Testing"
)

func TestThat_NewStdOutLogWriter_ReturnsLogWriter(t *testing.T) {
	// Test
	var sut *StdOutLogWriter = NewStdOutLogWriter()

	// Verify
	ExpectNonNil(sut, t)
}

func TestThat_StdOutLogWriter_Log_SendsMessageToStdOut(t *testing.T) {
	// Setup

	// Intercept StdOut to capture and examine output, then restore the original pipe
	// ref: https://eli.thegreenplace.net/2020/faking-stdin-and-stdout-in-go/
	r, w, err1 := os.Pipe()
	ExpectNoError(err1, t)
	origStdout := os.Stdout
	os.Stdout = w

	var sut *StdOutLogWriter = NewStdOutLogWriter()
	buf := make([]byte, 1024)
	expectedString := "just one string of bits among so many"

	// Test
	sut.Log(expectedString)
	actualBytesRead, err2 := r.Read(buf)
	actualString := string(buf[:actualBytesRead])

	// Verify
	os.Stdout = origStdout
	ExpectNoError(err2, t)
	ExpectTrue(actualBytesRead > 0, t)
	ExpectString(fmt.Sprintf("%s\n", expectedString), actualString, t)
}


