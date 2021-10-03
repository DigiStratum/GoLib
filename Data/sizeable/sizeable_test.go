package sizeable

import(
	"testing"

	. "github.com/DigiStratum/GoLib/Testing"
)

const FIXED_SIZE = 333

type empty_target struct {
	buffer	[10]int
}

type regular_target struct {
	Buffer	[50]int
}

type sizeable_target struct {
}

func (r sizeable_target) Size() int64 {
	return FIXED_SIZE
}

func TestThat_Sizeable_Interface_Matches(t *testing.T) {
	// Setup
	var sut interface{}
	sut = sizeable_target{}

	// Verify
	_, ok := sut.(SizeableIfc)
	ExpectTrue(ok, t)
}

func TestThat_Size_Func_Uses_Sizeable_Interface_Size(t *testing.T) {
	// Setup
	sut := sizeable_target{}

	// Verify
	ExpectInt64(FIXED_SIZE, Size(sut), t)
}

func TestThat_Size_Func_Returns_0_For_Structs_Without_Exported_Fields(t *testing.T) {
	// Setup
	sut := empty_target{}

	// Verify
	ExpectInt64(0, Size(sut), t)
}

func TestThat_Size_Func_Calculates_Size_For_Structs_With_Exported_Fields(t *testing.T) {
	// Setup
	sut := regular_target{}

	// Verify - forgive the magic number - this is a serialized byte size for the buffer of 50 ints
	ExpectInt64(121, Size(sut), t)
}
