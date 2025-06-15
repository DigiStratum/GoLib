package http

import (
	"testing"

	"github.com/DigiStratum/GoLib/Data/metadata"

	. "github.com/DigiStratum/GoLib/Testing"
)

// Factory Function

func TestThat_HttpRequest_NewHttpRequest_ReturnsInstance(t *testing.T) {
	// Setup & Test
	var sut HttpRequestIfc = NewHttpRequestBuilder().GetHttpRequest() // Verifies that result satisfies IFC

	// Verify
	if !ExpectNonNil(sut, t) {
		return
	}
}

// Getters/Setters

func TestThat_HttpRequest_GetURL_ReturnsEmptyWhenNotSet(t *testing.T) {
	// Setup
	sut := NewHttpRequestBuilder().GetHttpRequest()

	// Test
	actual := sut.GetURL()

	// Verify
	if !ExpectString("", actual, t) {
		return
	}
}

func TestThat_HttpRequest_GetMethod_ReturnsDefaultWhenNotSet(t *testing.T) {
	// Setup
	sut := NewHttpRequestBuilder().GetHttpRequest()

	// Test
	actual := sut.GetMethod()

	// Verify
	if !ExpectEqual(METHOD_GET, actual, t) {
		return
	}
}

func TestThat_HttpRequest_GetHeaders_ReturnsEmptyHeadersWhenNotSet(t *testing.T) {
	// Setup
	sut := NewHttpRequestBuilder().GetHttpRequest()

	// Test
	actual := sut.GetHeaders()

	// Verify
	if !ExpectNonNil(actual, t) {
		return
	}
	if !ExpectInt(0, len(*actual.GetNames()), t) {
		return
	}
}

// Builder functionality - testing the request builder creates proper request objects

func TestThat_HttpRequest_BuilderCreatesRequestWithURL(t *testing.T) {
	// Setup
	builder := NewHttpRequestBuilder()
	testUrl := "https://example.com/test"

	// Test
	sut := builder.SetURL(testUrl).GetHttpRequest()

	// Verify
	if !ExpectNonNil(sut, t) {
		return
	}
	if !ExpectEqual(testUrl, sut.GetURL(), t) {
		return
	}
	if !ExpectEqual("https", sut.GetScheme(), t) {
		return
	}
	if !ExpectEqual("example.com", sut.GetHost(), t) {
		return
	}
	if !ExpectEqual("/test", sut.GetURI(), t) {
		return
	}
}

func TestThat_HttpRequest_BuilderCreatesRequestWithMethod(t *testing.T) {
	// Setup
	builder := NewHttpRequestBuilder()

	// Test
	sut := builder.SetMethod(METHOD_POST).GetHttpRequest()

	// Verify
	if !ExpectNonNil(sut, t) {
		return
	}
	if !ExpectEqual(METHOD_POST, sut.GetMethod(), t) {
		return
	}
}

func TestThat_HttpRequest_BuilderCreatesRequestWithBody(t *testing.T) {
	// Setup
	builder := NewHttpRequestBuilder()
	testBody := "test body content"

	// Test
	sut := builder.SetBody(&testBody).GetHttpRequest()

	// Verify
	if !ExpectNonNil(sut, t) {
		return
	}
	if !ExpectString(testBody, *sut.GetBody(), t) {
		return
	}
}

func TestThat_HttpRequest_BuilderCreatesRequestWithHeaders(t *testing.T) {
	// Setup
	builder := NewHttpRequestBuilder()
	headers := NewHttpHeadersBuilder()
	headers.Set("Content-Type", "application/json")

	// Test
	sut := builder.SetHeaders(
		headers.GetHttpHeaders(),
	).GetHttpRequest()

	// Verify
	if !ExpectNonNil(sut, t) {
		return
	}
	if !ExpectNonNil(sut.GetHeaders(), t) {
		return
	}
	actual := sut.GetHeaders().Get("Content-Type")
	if !ExpectNonNil(actual, t) {
		return
	}
	if !ExpectInt(1, len(*actual), t) {
		return
	}
	if !ExpectString("application/json", (*actual)[0], t) {
		return
	}
}

func TestThat_HttpRequest_BuilderCreatesRequestWithQueryParameters(t *testing.T) {
	// Setup
	builder := NewHttpRequestBuilder()
	queryParams := metadata.NewMetadataBuilder().Set("param1", "value1").GetMetadata()

	// Test
	sut := builder.SetQueryParameters(queryParams).GetHttpRequest()

	// Verify
	if !ExpectNonNil(sut, t) {
		return
	}
	if !ExpectNonNil(sut.GetQueryParameters(), t) {
		return
	}
	actual := sut.GetQueryParameters().Get("param1")
	if !ExpectNonNil(actual, t) {
		return
	}
	if !ExpectString("value1", *actual, t) {
		return
	}
}

// Complex builder scenarios

func TestThat_HttpRequest_Builder_CreatesCompleteRequest(t *testing.T) {
	// Setup
	builder := NewHttpRequestBuilder()
	testUrl := "https://example.com/api/users?sort=desc"
	testBody := "{\"name\":\"test\"}"

	headers := NewHttpHeadersBuilder().
		Set("Content-Type", "application/json").
		Set("Authorization", "Bearer token123")

	// Test
	sut := builder.
		SetURL(testUrl).
		SetMethod(METHOD_POST).
		SetBody(&testBody).
		SetHeaders(headers.GetHttpHeaders()).GetHttpRequest()

	// Verify
	if !ExpectNonNil(sut, t) {
		return
	}
	if !ExpectString(testUrl, sut.GetURL(), t) {
		return
	}
	if !ExpectEqual(METHOD_POST, sut.GetMethod(), t) {
		return
	}
	if !ExpectString(testBody, *sut.GetBody(), t) {
		return
	}
	actual1 := sut.GetHeaders().Get("Content-Type")
	if !ExpectNonNil(actual1, t) {
		return
	}
	if !ExpectString("application/json", (*actual1)[0], t) {
		return
	}
	actual2 := sut.GetHeaders().Get("Authorization")
	if !ExpectNonNil(actual2, t) {
		return
	}
	if !ExpectString("Bearer token123", (*actual2)[0], t) {
		return
	}
	if !ExpectString("sort=desc", sut.GetQueryString(), t) {
		return
	}
}
