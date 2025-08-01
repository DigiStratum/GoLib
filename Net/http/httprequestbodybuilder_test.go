// filepath: /Users/skelly/Documents/GoProjects/GoLib/Net/http/httprequestbodybuilder_test.go
package http

import (
	"net/http"
	"net/url"
	"strings"
	"testing"

	. "github.com/DigiStratum/GoLib/Testing"
)

// Factory Function

func TestThat_HttpRequestBodyBuilder_NewHttpRequestBodyBuilder_ReturnsInstance(t *testing.T) {
	// Setup & Test
	sut := NewHttpRequestBodyBuilder()

	// Verify
	if !ExpectNonNil(sut, t) {
		return
	}
	if !ExpectNonNil(sut.requestBody, t) {
		return
	}
	if !ExpectNonNil(sut.requestBody.body, t) {
		return
	}
	if !ExpectTrue(sut.requestBody.IsEmpty(), t) {
		return
	}
}

func TestThat_HttpRequestBodyBuilder_NewHttpRequestBodyBuilder_ImplementsInterface(t *testing.T) {
	// Setup & Test
	var sut HttpRequestBodyBuilderIfc = NewHttpRequestBodyBuilder() // This assignment verifies interface compliance

	// Verify
	if !ExpectNonNil(sut, t) {
		return
	}
}

// NewHttpRequestBodyBuilderFromRequest Tests

func TestThat_HttpRequestBodyBuilder_NewHttpRequestBodyBuilderFromRequest_ReturnsNilForNilRequest(t *testing.T) {
	// Setup & Test
	sut := NewHttpRequestBodyBuilderFromRequest(nil)

	// Verify
	if !ExpectNil(sut, t) {
		return
	}
}

func TestThat_HttpRequestBodyBuilder_NewHttpRequestBodyBuilderFromRequest_HandlesEmptyFormData(t *testing.T) {
	// Setup
	request := &http.Request{Method: "GET"}

	// Test
	sut := NewHttpRequestBodyBuilderFromRequest(request)

	// Verify
	if !ExpectNonNil(sut, t) {
		return
	}
	if !ExpectNonNil(sut.requestBody, t) {
		return
	}
	if !ExpectTrue(sut.requestBody.IsEmpty(), t) {
		return
	}
}

func TestThat_HttpRequestBodyBuilder_NewHttpRequestBodyBuilderFromRequest_ProcessesFormData(t *testing.T) {
	// Setup
	formData := map[string][]string{
		"name":    {"John Doe"},
		"age":     {"30"},
		"hobbies": {"reading", "coding", "hiking"},
	}
	request := createTestRequest("POST", formData)

	// Test
	sut := NewHttpRequestBodyBuilderFromRequest(request)

	// Verify
	if !ExpectNonNil(sut, t) {
		return
	}
	if !ExpectNonNil(sut.requestBody, t) {
		return
	}

	// Check each form field was properly set
	if !ExpectTrue(sut.requestBody.Has("name"), t) {
		return
	}
	nameValues := sut.requestBody.Get("name")
	if !ExpectNonNil(nameValues, t) {
		return
	}
	if !ExpectInt(1, len(*nameValues), t) {
		return
	}
	if !ExpectString("John Doe", (*nameValues)[0], t) {
		return
	}

	if !ExpectTrue(sut.requestBody.Has("age"), t) {
		return
	}
	ageValues := sut.requestBody.Get("age")
	if !ExpectNonNil(ageValues, t) {
		return
	}
	if !ExpectInt(1, len(*ageValues), t) {
		return
	}
	if !ExpectString("30", (*ageValues)[0], t) {
		return
	}

	if !ExpectTrue(sut.requestBody.Has("hobbies"), t) {
		return
	}
	hobbiesValues := sut.requestBody.Get("hobbies")
	if !ExpectNonNil(hobbiesValues, t) {
		return
	}
	if !ExpectInt(3, len(*hobbiesValues), t) {
		return
	}
	if !ExpectString("reading", (*hobbiesValues)[0], t) {
		return
	}
	if !ExpectString("coding", (*hobbiesValues)[1], t) {
		return
	}
	if !ExpectString("hiking", (*hobbiesValues)[2], t) {
		return
	}
}

// Set Method

func TestThat_HttpRequestBodyBuilder_Set_AddsNewNameValuePair(t *testing.T) {
	// Setup
	sut := NewHttpRequestBodyBuilder()

	// Test
	result := sut.Set("param", "value")

	// Verify
	if !ExpectNonNil(result, t) {
		return
	}
	if !ExpectTrue(sut.requestBody.Has("param"), t) {
		return
	}
	values := sut.requestBody.Get("param")
	if !ExpectNonNil(values, t) {
		return
	}
	if !ExpectInt(1, len(*values), t) {
		return
	}
	if !ExpectString("value", (*values)[0], t) {
		return
	}
}

func TestThat_HttpRequestBodyBuilder_Set_AddsMultipleValuesForSameName(t *testing.T) {
	// Setup
	sut := NewHttpRequestBodyBuilder()

	// Test
	sut.Set("param", "value1")
	sut.Set("param", "value2")

	// Verify
	if !ExpectTrue(sut.requestBody.Has("param"), t) {
		return
	}
	values := sut.requestBody.Get("param")
	if !ExpectNonNil(values, t) {
		return
	}
	if !ExpectInt(2, len(*values), t) {
		return
	}
	if !ExpectString("value1", (*values)[0], t) {
		return
	}
	if !ExpectString("value2", (*values)[1], t) {
		return
	}
}

func TestThat_HttpRequestBodyBuilder_Set_AddsMultipleValuesInSingleCall(t *testing.T) {
	// Setup
	sut := NewHttpRequestBodyBuilder()

	// Test
	sut.Set("param", "value1", "value2", "value3")

	// Verify
	if !ExpectTrue(sut.requestBody.Has("param"), t) {
		return
	}
	values := sut.requestBody.Get("param")
	if !ExpectNonNil(values, t) {
		return
	}
	if !ExpectInt(3, len(*values), t) {
		return
	}
	if !ExpectString("value1", (*values)[0], t) {
		return
	}
	if !ExpectString("value2", (*values)[1], t) {
		return
	}
	if !ExpectString("value3", (*values)[2], t) {
		return
	}
}

// Merge Method

func TestThat_HttpRequestBodyBuilder_Merge_HandlesNilInput(t *testing.T) {
	// Setup
	sut := NewHttpRequestBodyBuilder()
	sut.Set("existing", "value")

	// Test
	result := sut.Merge(nil)

	// Verify
	if !ExpectNonNil(result, t) {
		return
	}
	// Should still have the existing value
	if !ExpectTrue(sut.requestBody.Has("existing"), t) {
		return
	}
}

func TestThat_HttpRequestBodyBuilder_Merge_AddsFromOtherBody(t *testing.T) {
	// Setup
	sut := NewHttpRequestBodyBuilder()
	sut.Set("param1", "value1")

	otherBody := NewHttpRequestBodyBuilder()
	otherBody.Set("param2", "value2")

	// Test
	result := sut.Merge(otherBody.GetHttpRequestBody())

	// Verify
	if !ExpectNonNil(result, t) {
		return
	}
	if !ExpectTrue(sut.requestBody.Has("param1"), t) {
		return
	}
	if !ExpectTrue(sut.requestBody.Has("param2"), t) {
		return
	}

	values1 := sut.requestBody.Get("param1")
	if !ExpectNonNil(values1, t) {
		return
	}
	if !ExpectInt(1, len(*values1), t) {
		return
	}
	if !ExpectString("value1", (*values1)[0], t) {
		return
	}

	values2 := sut.requestBody.Get("param2")
	if !ExpectNonNil(values2, t) {
		return
	}
	if !ExpectInt(1, len(*values2), t) {
		return
	}
	if !ExpectString("value2", (*values2)[0], t) {
		return
	}
}

func TestThat_HttpRequestBodyBuilder_Merge_MergesValuesForSameName(t *testing.T) {
	// Setup
	sut := NewHttpRequestBodyBuilder()
	sut.Set("param", "value1")

	otherBody := NewHttpRequestBodyBuilder()
	otherBody.Set("param", "value2")

	// Test
	result := sut.Merge(otherBody.GetHttpRequestBody())

	// Verify
	if !ExpectNonNil(result, t) {
		return
	}
	if !ExpectTrue(sut.requestBody.Has("param"), t) {
		return
	}

	values := sut.requestBody.Get("param")
	if !ExpectNonNil(values, t) {
		return
	}
	if !ExpectInt(2, len(*values), t) {
		return
	}
	if !ExpectString("value1", (*values)[0], t) {
		return
	}
	if !ExpectString("value2", (*values)[1], t) {
		return
	}
}

func TestThat_HttpRequestBodyBuilder_Merge_HandlesMergeWithEmptyNames(t *testing.T) {
	// Setup
	sut := NewHttpRequestBodyBuilder()
	sut.Set("param", "value")

	// Create a mock that returns nil for GetNames
	mockBody := &mockHttpRequestBody{
		getNamesFunc: func() *[]string {
			return nil
		},
	}

	// Test
	result := sut.Merge(mockBody)

	// Verify - should keep original values and not crash
	if !ExpectNonNil(result, t) {
		return
	}
	if !ExpectTrue(sut.requestBody.Has("param"), t) {
		return
	}
}

// GetHttpRequestBody Method

func TestThat_HttpRequestBodyBuilder_GetHttpRequestBody_ReturnsBody(t *testing.T) {
	// Setup
	sut := NewHttpRequestBodyBuilder()
	sut.Set("param", "value")

	// Test
	body := sut.GetHttpRequestBody()

	// Verify
	if !ExpectNonNil(body, t) {
		return
	}
	if !ExpectTrue(body.Has("param"), t) {
		return
	}
	values := body.Get("param")
	if !ExpectNonNil(values, t) {
		return
	}
	if !ExpectInt(1, len(*values), t) {
		return
	}
	if !ExpectString("value", (*values)[0], t) {
		return
	}
}

func TestNewHttpRequestBodyBuilder(t *testing.T) {
	builder := NewHttpRequestBodyBuilder()
	if builder == nil {
		t.Fatal("Expected non-nil builder")
	}

	body := builder.GetHttpRequestBody()
	if body == nil {
		t.Fatal("Expected non-nil body")
	}

	if !body.IsEmpty() {
		t.Error("Expected empty body")
	}
}

// Helper function to create test HTTP requests with form data
func createTestRequest(method string, formData map[string][]string) *http.Request {
	form := url.Values{}
	for name, values := range formData {
		for _, value := range values {
			form.Add(name, value)
		}
	}

	formEncoded := form.Encode()
	req, _ := http.NewRequest(method, "http://example.com", strings.NewReader(formEncoded))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	return req
}

// Mock for testing

type mockHttpRequestBody struct {
	hasFunc      func(name string) bool
	getNamesFunc func() *[]string
	isEmptyFunc  func() bool
	getFunc      func(name string) *[]string
	sizeFunc     func() int
}

func (m *mockHttpRequestBody) Has(name string) bool {
	if m.hasFunc != nil {
		return m.hasFunc(name)
	}
	return false
}

func (m *mockHttpRequestBody) GetNames() *[]string {
	if m.getNamesFunc != nil {
		return m.getNamesFunc()
	}
	return nil
}

func (m *mockHttpRequestBody) IsEmpty() bool {
	if m.isEmptyFunc != nil {
		return m.isEmptyFunc()
	}
	return true
}

func (m *mockHttpRequestBody) Get(name string) *[]string {
	if m.getFunc != nil {
		return m.getFunc(name)
	}
	return nil
}

func (m *mockHttpRequestBody) Size() int {
	if m.sizeFunc != nil {
		return m.sizeFunc()
	}
	return 0
}
