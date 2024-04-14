package fileio

import(
	"testing"
	"runtime"
	"path/filepath"

	. "github.com/DigiStratum/GoLib/Testing"
)

func TestThat_NewFile_ReturnsSomething(t *testing.T) {
	// Test
	sut := NewFile("missingfile.txt")

	// Verify
	ExpectNonNil(sut, t)
}

func TestThat_NewFile_GetPath_Returns_OriginalPath(t *testing.T) {
	// Setup
	expected := "missingfile.txt"
	sut := NewFile(expected)

	// Test
	actual := sut.GetPath()

	// Verify
	ExpectString(expected, actual, t)
}

// TODO: Figure out some incantation that throws an error instead of resolved path
func TestThat_NewFile_GetAbsPath_Returns_GoodPath(t *testing.T) {
	// Setup
	_, filename, _, _ := runtime.Caller(0)
	expected := filepath.Dir(filename) + "/missingfile.txt"
	sut := NewFile("./missingfile.txt")

	// Test
	actual, _ := sut.GetAbsPath()

	// Verify
	ExpectNonNil(actual, t)
	ExpectString(expected, *actual, t)
}

func TestThat_NewFile_Exists_Returns_False_ForMissingFile(t *testing.T) {
	// Setup
	sut := NewFile("missingfile.txt")

	// Test
	actual := sut.Exists()

	// Verify
	ExpectFalse(actual, t)
}

