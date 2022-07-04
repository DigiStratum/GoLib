package objectcollection

import(
	"fmt"
	"testing"

	obj "github.com/DigiStratum/GoLib/Object"

        . "github.com/DigiStratum/GoLib/Testing"
)

func TestThat_ObjectCollection_NewObjectCollection_ReturnsSomething(t *testing.T) {
	// Test
	sut := NewObjectCollection()

	// Verify
	ExpectNonNil(sut, t)
}

func TestThat_ObjectCollection_HasObject_ReturnsFalse_ForMissingPath(t *testing.T) {
	// Setup
	sut := NewObjectCollection()

	// Test
	actual := sut.HasObject("")

	// Verify
	ExpectFalse(actual, t)
}

func TestThat_ObjectCollection_GetObject_ReturnsNil_ForMissingPath(t *testing.T) {
	// Setup
	sut := NewObjectCollection()

	// Test
	actual := sut.GetObject("")

	// Verify
	ExpectNil(actual, t)
}

func TestThat_ObjectCollection_PutObject_ReturnsError_ForNilObject(t *testing.T) {
	// Setup
	sut := NewObjectCollection()

	// Test
	actual := sut.PutObject("", nil)

	// Verify
	ExpectError(actual, t)
}

func TestThat_ObjectCollection_PutHasGetObject_Works_ForGoodPathObject(t *testing.T) {
	// Setup
	sut := NewObjectCollection()
	objs := make([]*obj.Object, 3)
	for i := 0; i < 3; i++ {
		objs[i] = obj.NewObject()
		expectedContent := fmt.Sprintf("object %d", i)
		objs[i].SetContent(&expectedContent)
	}

	// Test & Verify
	for i := 0; i < 3; i++ {
		path := fmt.Sprintf("path %d", i)
		err := sut.PutObject(path, objs[i])
		ExpectNoError(err, t)
		hasIt := sut.HasObject(path)
		ExpectTrue(hasIt, t)
		actual := sut.GetObject(path)
		ExpectNonNil(actual, t)
		actualContent := actual.GetContent()
		ExpectNonNil(actualContent, t)
		expectedContent := fmt.Sprintf("object %d", i)
		ExpectString(expectedContent, *actualContent, t)
	}
}


func TestThat_ObjectCollection_GetIterator_ReturnsIteratorWithNoResults_When_Empty(t *testing.T) {
	// Setup
	sut := NewObjectCollection()

	// Test
	it := sut.GetIterator()
	actual := it()

	// Verify
	ExpectNil(actual, t)
}

func TestThat_ObjectCollection_GetIterator_ReturnsIteratorWithResults_When_Populated(t *testing.T) {
	// Setup
	sut := NewObjectCollection()
	objs := make([]*obj.Object, 3)
	for i := 0; i < 3; i++ {
		objs[i] = obj.NewObject()
		expectedContent := fmt.Sprintf("object %d", i)
		objs[i].SetContent(&expectedContent)
		sut.PutObject(expectedContent, objs[i])
	}

	// Test
	foundObjs := make(map[string]*obj.Object)
	it := sut.GetIterator()
	for actual := it(); actual != nil; actual = it() {
		actualPOP, ok := actual.(PathObjectPair)
		ExpectTrue(ok, t)
		foundObjs[actualPOP.Path] = actualPOP.Obj
		ExpectNonNil(actualPOP.Obj, t)
		actualContent := actualPOP.Obj.GetContent()
		ExpectString(actualPOP.Path, *actualContent, t)
	}

	// Verify
	ExpectInt(len(objs), len(foundObjs), t)
}

