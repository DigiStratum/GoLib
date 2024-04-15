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

func TestThat_File_GetPath_Returns_OriginalPath(t *testing.T) {
	// Setup
	expected := "missingfile.txt"
	sut := NewFile(expected)

	// Test
	actual := sut.GetPath()

	// Verify
	ExpectString(expected, actual, t)
}

// TODO: Figure out some incantation that throws an error instead of resolved path
func TestThat_File_GetAbsPath_Returns_GoodPath(t *testing.T) {
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

func TestThat_File_Exists_Returns_False_ForMissingFile(t *testing.T) {
	// Setup
	sut := NewFile("missingfile.txt")

	// Test
	actual := sut.Exists()

	// Verify
	ExpectFalse(actual, t)
}

func TestThat_File_Exists_Returns_True_ForGoodFile(t *testing.T) {
	// Setup
	sut := NewFile("testfile.txt")

	// Test
	actual := sut.Exists()

	// Verify
	ExpectTrue(actual, t)
}

func TestThat_File_GetName_Returns_Name_ForGoodFile(t *testing.T) {
	// Setup
	expected := "testfile.txt"
	sut := NewFile(expected)

	// Test
	actual, err := sut.GetName()

	// Verify
	ExpectNoError(err, t)
	ExpectNonNil(actual, t)
	ExpectString(expected, *actual, t)
}

func TestThat_File_GetName_Returns_Error_ForMissingFile(t *testing.T) {
	// Setup
	expected := "missingfile.txt"
	sut := NewFile(expected)

	// Test
	actual, err := sut.GetName()

	// Verify
	ExpectError(err, t)
	ExpectNil(actual, t)
}

func TestThat_File_GetSize_Returns_Error_ForMissingFile(t *testing.T) {
	// Setup
	expected := "missingfile.txt"
	sut := NewFile(expected)

	// Test
	actual, err := sut.GetSize()

	// Verify
	ExpectError(err, t)
	ExpectNil(actual, t)
}

func TestThat_File_GetSize_Returns_Number_ForGoodFile(t *testing.T) {
	// Setup
	sut := NewFile("testfile.txt")

	// Test
	actual, err := sut.GetSize()

	// Verify
	ExpectNoError(err, t)
	ExpectNonNil(actual, t)
	ExpectInt64(21, *actual, t)
}

func TestThat_File_GetMode_Returns_Error_ForMissingFile(t *testing.T) {
	// Setup
	expected := "missingfile.txt"
	sut := NewFile(expected)

	// Test
	actual, err := sut.GetMode()

	// Verify
	ExpectError(err, t)
	ExpectNil(actual, t)
}

func TestThat_File_GetMode_Returns_FileMode_ForGoodFile(t *testing.T) {
	// Setup
	sut := NewFile("testfile.txt")

	// Test
	actual, err := sut.GetMode()

	// Verify
	ExpectNoError(err, t)
	ExpectNonNil(actual, t)
}

func TestThat_File_GetModTime_Returns_Error_ForMissingFile(t *testing.T) {
	// Setup
	expected := "missingfile.txt"
	sut := NewFile(expected)

	// Test
	actual, err := sut.GetModTime()

	// Verify
	ExpectError(err, t)
	ExpectNil(actual, t)
}

func TestThat_File_GetModTime_Returns_Time_ForGoodFile(t *testing.T) {
	// Setup
	sut := NewFile("testfile.txt")

	// Test
	actual, err := sut.GetModTime()

	// Verify
	ExpectNoError(err, t)
	ExpectNonNil(actual, t)
}

func TestThat_File_IsDir_Returns_Error_ForMissingFile(t *testing.T) {
	// Setup
	expected := "missingdir"
	sut := NewFile(expected)

	// Test
	_, err := sut.IsDir()

	// Verify
	ExpectError(err, t)
}

func TestThat_File_IsDir_Returns_True_ForGoodFileDir(t *testing.T) {
	// Setup
	sut := NewFile(".")

	// Test
	actual, err := sut.IsDir()

	// Verify
	ExpectNoError(err, t)
	ExpectTrue(actual, t)
}

func TestThat_File_GetSys_Returns_Error_ForMissingFile(t *testing.T) {
	// Setup
	expected := "missingdir"
	sut := NewFile(expected)

	// Test
	_, err := sut.GetSys()

	// Verify
	ExpectError(err, t)
}

func TestThat_File_GetSys_Returns_Something_ForGoodFileDir(t *testing.T) {
	// Setup
	sut := NewFile("testfile.txt")

	// Test
	actual, err := sut.GetSys()

	// Verify
	ExpectNoError(err, t)
	ExpectNonNil(actual, t)
}

