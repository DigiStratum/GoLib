package stringset

import(
	"testing"

	. "github.com/DigiStratum/GoLib/Testing"
)

func TestThat_NewStringSet_ReturnsSomething(t *testing.T) {
	// Test
	sut := NewStringSet()

	// Verify
	ExpectNonNil(sut, t)
}

func TestThat_StringSet_Size_Returns0_WhenEmpty(t *testing.T) {
	// Setup
	sut := NewStringSet()

	// Test
	actual := sut.Size()

	// Verify
	ExpectInt(0, actual, t)
	ExpectTrue(sut.IsEmpty(), t)
}

func TestThat_StringSet_Set_SetsString(t *testing.T) {
	// Setup
	sut := NewStringSet()
	expected := "sillystring!"

	// Test
	sut.Set(expected)
	actual := sut.Size()

	// Verify
	ExpectInt(1, actual, t)
	ExpectTrue(sut.Has(expected), t)
	ExpectFalse(sut.Has("missingstring"), t)
}

func TestThat_StringSet_SetAll_SetsMultipleStrings(t *testing.T) {
	// Setup
	sut := NewStringSet()
	ss1 := []string{"silly", "string", "spectacular"}
	expected := len(ss1)

	// Test
	sut.SetAll(&ss1)
	actual := sut.Size()

	// Verify
	ExpectInt(expected, actual, t)
	ExpectTrue(sut.HasAll(&ss1), t)
	ss1 = append(ss1, "missingstring")
	ExpectFalse(sut.HasAll(&ss1), t)
	ss2 := []string{ "missingstring" }
	ExpectFalse(sut.HasAny(&ss2), t)
}

func TestThat_StringSet_Merge_MergesMultipleStrings(t *testing.T) {
	// Setup
	ss1 := []string{"silly", "string"}
	sut := NewStringSet()
	sut.SetAll(&ss1)

	ss2 := []string{"spectacular", "sillystring!"}
	mergeSet := NewStringSet()
	mergeSet.SetAll(&ss2)

	expected := len(ss1) + len(ss2)

	// Test
	sut.Merge(mergeSet)
	actual := sut.Size()

	// Verify
	ExpectInt(expected, actual, t)
	ExpectTrue(sut.HasAll(&ss1), t)
	ExpectTrue(sut.HasAll(&ss2), t)
}

func TestThat_StringSet_Drop_DropsString(t *testing.T) {
	// Setup
	sut := NewStringSet()
	unexpected := "sillystring!"
	sut.Set(unexpected)

	// Test
	sut.Drop(unexpected)
	actual := sut.Size()

	// Verify
	ExpectInt(0, actual, t)
	ExpectFalse(sut.Has(unexpected), t)
}

func TestThat_StringSet_DropAll_DropsAllStrings(t *testing.T) {
	// Setup
	ss1 := []string{"silly", "string"}
	ss2 := []string{"spectacular", "sillystring!"}
	sut := NewStringSet()
	sut.SetAll(&ss1)
	sut.SetAll(&ss2)
	expected := len(ss1)

	// Test
	sut.DropAll(&ss2)
	actual := sut.Size()

	// Verify
	ExpectInt(expected, actual, t)
	ExpectTrue(sut.HasAll(&ss1), t)
	ExpectFalse(sut.HasAll(&ss2), t)
	ExpectFalse(sut.HasAny(&ss2), t)
}

func TestThat_StringSet_ToArray_ReturnsPopulatedArray_WhenNonEmpty(t *testing.T) {
	// Setup
	sut := NewStringSet()
	expected := []string{"silly", "string", "spectacular"}
	sut.SetAll(&expected)

	// Test
	actual := sut.ToArray()

	// Verify
	ExpectNonNil(actual, t)
	ExpectTrue(sut.HasAll(actual), t)
	ExpectInt(len(expected), len(*actual), t)
}

func TestThat_StringSet_ToArray_ReturnsEmptyArray_WhenEmpty(t *testing.T) {
	// Setup
	sut := NewStringSet()

	// Test
	actual := sut.ToArray()

	// Verify
	ExpectNonNil(actual, t)
	ExpectInt(0, len(*actual), t)
}

func TestThat_StringSet_GetIterator_ReturnsIteratorFunction_ThatWorksYay(t *testing.T) {
	// Setup
	sut := NewStringSet()
	expected := []string{"silly", "string", "spectacular"}
	sut.SetAll(&expected)

	// Test
	it := sut.GetIterator()
	actual := make([]string, 0)
	for namei := it(); namei != nil; namei = it() {
		namep := namei.(*string)
		if nil != namep {
			actual = append(actual, *namep)
		}
	}

	// Verify
	ExpectInt(len(expected), len(actual), t)
	ExpectTrue(sut.HasAll(&actual), t)
}
