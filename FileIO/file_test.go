package fileio

import(
	"testing"
	"runtime"
	"path/filepath"

	. "github.com/DigiStratum/GoLib/Testing"
)

const TEST_RESOURCE_DIR = "res_test"
const TEST_FILE = "testfile.txt"
const MISSING_FILE = "missingfile.txt"

func TestThat_NewFile_Returns_FileIfc(t *testing.T) {
	// Test
	var sut FileIfc = NewFile(MISSING_FILE)

	// Verify
	ExpectNonNil(sut, t)
}

func TestThat_File_GetPath_Returns_OriginalPath(t *testing.T) {
	// Setup
	sut := NewFile(MISSING_FILE)

	// Test
	actual := sut.GetPath()

	// Verify
	ExpectString(MISSING_FILE, actual, t)
}

/*
// There does not appear to be a code path to inject an error in to system libraries that causes this to fail
// ref: https://stackoverflow.com/questions/16742331/how-to-mock-abstract-filesystem-in-go
func TestThat_File_GetAbsPath_Returns_Error_ForBadFilepath(t *testing.T) {
	// Setup
	sut := NewFile("")

	// Test
	actual, err := sut.GetAbsPath()

	// Verify
	ExpectNil(actual, t)
	ExpectError(err, t)
}
*/

func TestThat_File_GetAbsPath_Returns_GoodPath(t *testing.T) {
	// Setup
	_, filename, _, _ := runtime.Caller(0)
	expected := filepath.Join(filepath.Dir(filename), MISSING_FILE)
	sut := NewFile(filepath.Join(".", MISSING_FILE))

	// Test
	actual, err := sut.GetAbsPath()

	// Verify
	ExpectNonNil(actual, t)
	ExpectNoError(err, t)
	ExpectString(expected, *actual, t)
}

func TestThat_File_Exists_Returns_False_ForMissingFile(t *testing.T) {
	// Setup
	sut := NewFile(MISSING_FILE)

	// Test
	actual := sut.Exists()

	// Verify
	ExpectFalse(actual, t)
}

func TestThat_File_Exists_Returns_True_ForGoodFile(t *testing.T) {
	// Setup
	sut := NewFile(filepath.Join(TEST_RESOURCE_DIR, TEST_FILE))

	// Test
	actual := sut.Exists()

	// Verify
	ExpectTrue(actual, t)
}

func TestThat_File_GetName_Returns_Name_ForGoodFile(t *testing.T) {
	// Setup
	sut := NewFile(filepath.Join(TEST_RESOURCE_DIR, TEST_FILE))

	// Test
	actual, err := sut.GetName()

	// Verify
	ExpectNoError(err, t)
	ExpectNonNil(actual, t)
	ExpectString(TEST_FILE, *actual, t)
}

func TestThat_File_GetName_Returns_Error_ForMissingFile(t *testing.T) {
	// Setup
	sut := NewFile(MISSING_FILE)

	// Test
	actual, err := sut.GetName()

	// Verify
	ExpectError(err, t)
	ExpectNil(actual, t)
}

func TestThat_File_GetSize_Returns_Error_ForMissingFile(t *testing.T) {
	// Setup
	sut := NewFile(MISSING_FILE)

	// Test
	actual, err := sut.GetSize()

	// Verify
	ExpectError(err, t)
	ExpectNil(actual, t)
}

func TestThat_File_GetSize_Returns_Number_ForGoodFile(t *testing.T) {
	// Setup
	sut := NewFile(filepath.Join(TEST_RESOURCE_DIR, TEST_FILE))

	// Test
	actual, err := sut.GetSize()

	// Verify
	ExpectNoError(err, t)
	ExpectNonNil(actual, t)
	ExpectInt64(21, *actual, t)
}

func TestThat_File_GetMode_Returns_Error_ForMissingFile(t *testing.T) {
	// Setup
	sut := NewFile(MISSING_FILE)

	// Test
	actual, err := sut.GetMode()

	// Verify
	ExpectError(err, t)
	ExpectNil(actual, t)
}

func TestThat_File_GetMode_Returns_FileMode_ForGoodFile(t *testing.T) {
	// Setup
	sut := NewFile(filepath.Join(TEST_RESOURCE_DIR, TEST_FILE))

	// Test
	actual, err := sut.GetMode()

	// Verify
	ExpectNoError(err, t)
	ExpectNonNil(actual, t)
}

func TestThat_File_GetModTime_Returns_Error_ForMissingFile(t *testing.T) {
	// Setup
	sut := NewFile(MISSING_FILE)

	// Test
	actual, err := sut.GetModTime()

	// Verify
	ExpectError(err, t)
	ExpectNil(actual, t)
}

func TestThat_File_GetModTime_Returns_Time_ForGoodFile(t *testing.T) {
	// Setup
	sut := NewFile(filepath.Join(TEST_RESOURCE_DIR, TEST_FILE))

	// Test
	actual, err := sut.GetModTime()

	// Verify
	ExpectNoError(err, t)
	ExpectNonNil(actual, t)
}

func TestThat_File_IsDir_Returns_Error_ForMissingFile(t *testing.T) {
	// Setup
	sut := NewFile(MISSING_FILE)

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
	sut := NewFile(MISSING_FILE)

	// Test
	_, err := sut.GetSys()

	// Verify
	ExpectError(err, t)
}

func TestThat_File_GetSys_Returns_Something_ForGoodFileDir(t *testing.T) {
	// Setup
	sut := NewFile(filepath.Join(TEST_RESOURCE_DIR, TEST_FILE))

	// Test
	actual, err := sut.GetSys()

	// Verify
	ExpectNoError(err, t)
	ExpectNonNil(actual, t)
}

