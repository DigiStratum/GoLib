package data

import(
	"testing"

	. "GoLib/Testing"
)

func TestThat_DataType_ToString_Returns_expected_type_descriptors(t *testing.T) {
        // Verify
	if ! ExpectString("invalid", DATA_TYPE_INVALID.ToString(), t) { return }
	if ! ExpectString("null", DATA_TYPE_NULL.ToString(), t) { return }
	if ! ExpectString("boolean", DATA_TYPE_BOOLEAN.ToString(), t) { return }
	if ! ExpectString("integer", DATA_TYPE_INTEGER.ToString(), t) { return }
	if ! ExpectString("float", DATA_TYPE_FLOAT.ToString(), t) { return }
	if ! ExpectString("string", DATA_TYPE_STRING.ToString(), t) { return }
	if ! ExpectString("object", DATA_TYPE_OBJECT.ToString(), t) { return }
	if ! ExpectString("array", DATA_TYPE_ARRAY.ToString(), t) { return }
	dt := DATA_TYPE_ARRAY + 1
	if ! ExpectString("", dt.ToString(), t) { return }
}


