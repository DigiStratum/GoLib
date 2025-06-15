package http

import (
	"testing"

	obj "github.com/DigiStratum/GoLib/Object"

	. "github.com/DigiStratum/GoLib/Testing"
)

// Factory Functions - Standard Responses

func TestThat_HttpResponse_NewHttpResponseStandard_ReturnsCorrectResponse(t *testing.T) {
	// Setup
	expectedStatus := STATUS_OK
	expectedBody := "test body"
	expectedContentType := "text/plain"

	// Test
	sut := NewHttpResponseStandard(expectedStatus, &expectedBody, expectedContentType)

	// Verify
	if !ExpectNonNil(sut, t) {
		return
	}
	if !ExpectEqual(expectedStatus, sut.GetStatus(), t) {
		return
	}
	if !ExpectString(expectedBody, *sut.GetBody(), t) {
		return
	}
	actual := sut.GetHeaders().Get("content-type")
	if !ExpectNonNil(actual, t) {
		return
	}
	if !ExpectTrue((len(*actual) > 0), t) {
		return
	}
	if !ExpectString(expectedContentType, (*actual)[0], t) {
		return
	}
}

func TestThat_HttpResponse_NewHttpReponseCode_ReturnsCodeOnlyResponse(t *testing.T) {
	// Setup
	expectedStatus := STATUS_NOT_FOUND

	// Test
	sut := NewHttpReponseCode(expectedStatus)

	// Verify
	if !ExpectNonNil(sut, t) {
		return
	}
	if !ExpectEqual(expectedStatus, sut.GetStatus(), t) {
		return
	}
	if !ExpectString("", *sut.GetBody(), t) {
		return
	}
}

func TestThat_HttpResponse_NewHttpResponseSimpleJson_ReturnsJsonResponse(t *testing.T) {
	// Setup & Test
	sut := NewHttpResponseSimpleJson(STATUS_OK)

	// Verify
	if !ExpectNonNil(sut, t) {
		return
	}
	if !ExpectEqual(STATUS_OK, sut.GetStatus(), t) {
		return
	}
	actual := sut.GetHeaders().Get("content-type")
	if !ExpectNonNil(actual, t) {
		return
	}
	if !ExpectTrue((len(*actual) > 0), t) {
		return
	}
	if !ExpectString("application/json", (*actual)[0], t) {
		return
	}
	// FIXME: non-existent expectation func; add or refactor this?
	//
	//	if !ExpectContains("\"msg\":", *sut.GetBody(), t) {
	//		return
	//	}
}

func TestThat_HttpResponse_NewHttpResponseSimpleJson_ReturnsErrorJsonResponse(t *testing.T) {
	// Setup & Test
	sut := NewHttpResponseSimpleJson(STATUS_BAD_REQUEST)

	// Verify
	if !ExpectNonNil(sut, t) {
		return
	}
	if !ExpectEqual(STATUS_BAD_REQUEST, sut.GetStatus(), t) {
		return
	}
	actual := sut.GetHeaders().Get("content-type")
	if !ExpectNonNil(actual, t) {
		return
	}
	if !ExpectTrue((len(*actual) > 0), t) {
		return
	}
	if !ExpectString("application/json", (*actual)[0], t) {
		return
	}
	// FIXME: non-existent expectation func; add or refactor this?
	//if !ExpectContains("\"error\":", *sut.GetBody(), t) {
	//	return
	//}
}

func TestThat_HttpResponse_NewHttpResponseError_ReturnsErrorResponse(t *testing.T) {
	// Setup & Test
	sut := NewHttpResponseError(STATUS_INTERNAL_SERVER_ERROR)

	// Verify
	if !ExpectNonNil(sut, t) {
		return
	}
	if !ExpectEqual(STATUS_INTERNAL_SERVER_ERROR, sut.GetStatus(), t) {
		return
	}
	actual := sut.GetHeaders().Get("content-type")
	if !ExpectNonNil(actual, t) {
		return
	}
	if !ExpectTrue((len(*actual) > 0), t) {
		return
	}
	if !ExpectString("application/json", (*actual)[0], t) {
		return
	}
}

func TestThat_HttpResponse_NewHttpResponseOk_ReturnsOkResponse(t *testing.T) {
	// Setup
	expectedBody := "success"
	expectedContentType := "text/plain"

	// Test
	sut := NewHttpResponseOk(&expectedBody, expectedContentType)

	// Verify
	if !ExpectNonNil(sut, t) {
		return
	}
	if !ExpectEqual(STATUS_OK, sut.GetStatus(), t) {
		return
	}
	if !ExpectString(expectedBody, *sut.GetBody(), t) {
		return
	}
	actual := sut.GetHeaders().Get("content-type")
	if !ExpectNonNil(actual, t) {
		return
	}
	if !ExpectTrue((len(*actual) > 0), t) {
		return
	}
	if !ExpectString(expectedContentType, (*actual)[0], t) {
		return
	}
}

func TestThat_HttpResponse_NewHttpResponseErrorJson_ReturnsCorrectErrorJson(t *testing.T) {
	// Setup
	expectedStatus := STATUS_NOT_FOUND
	expectedMessage := "Resource not found"

	// Test
	sut := NewHttpResponseErrorJson(expectedStatus, expectedMessage)

	// Verify
	if !ExpectNonNil(sut, t) {
		return
	}
	if !ExpectEqual(expectedStatus, sut.GetStatus(), t) {
		return
	}
	actual := sut.GetHeaders().Get("content-type")
	if !ExpectNonNil(actual, t) {
		return
	}
	if !ExpectTrue((len(*actual) > 0), t) {
		return
	}
	if !ExpectString("application/json", (*actual)[0], t) {
		return
	}
	// FIXME: non-existent expectation func; add or refactor this?
	//if !ExpectContains(expectedMessage, *sut.GetBody(), t) {
	//	return
	//}
}

func TestThat_HttpResponse_NewHttpResponseObject_ReturnsObjectResponse(t *testing.T) {
	// Setup
	content := "test content"
	object := obj.NewObject()
	object.SetContent(&content)

	// Test
	sut := NewHttpResponseObject(object, "/test.txt")

	// Verify
	if !ExpectNonNil(sut, t) {
		return
	}
	if !ExpectEqual(STATUS_OK, sut.GetStatus(), t) {
		return
	}
	actual := sut.GetHeaders().Get("content-type")
	if !ExpectNonNil(actual, t) {
		return
	}
	if !ExpectTrue((len(*actual) > 0), t) {
		return
	}
	if !ExpectString("text/plain", (*actual)[0], t) {
		return
	}
	if !ExpectString(content, *sut.GetBody(), t) {
		return
	}
}

func TestThat_HttpResponse_NewHttpResponseObjectCacheable_ReturnsCacheableResponse(t *testing.T) {
	// Setup
	content := "test content"
	object := obj.NewObject()
	object.SetContent(&content)
	maxAge := 3600

	// Test
	sut := NewHttpResponseObjectCacheable(object, "/test.html", maxAge)

	// Verify
	if !ExpectNonNil(sut, t) {
		return
	}
	if !ExpectEqual(STATUS_OK, sut.GetStatus(), t) {
		return
	}
	actual := sut.GetHeaders().Get("content-type")
	if !ExpectNonNil(actual, t) {
		return
	}
	if !ExpectTrue((len(*actual) > 0), t) {
		return
	}
	if !ExpectString("text/plain", (*actual)[0], t) {
		return
	}
	if !ExpectTrue(sut.GetHeaders().Has("cache-control"), t) {
		return
	}
	actual = sut.GetHeaders().Get("cache-control")
	if !ExpectTrue((len(*actual) > 0), t) {
		return
	}
	if !ExpectEqual("max-age=3600", (*actual)[0], t) {
		return
	}
}

func TestThat_HttpResponse_NewHttpResponseRedirect_ReturnsRedirectResponse(t *testing.T) {
	// Setup
	expectedUrl := "https://example.com/new-location"

	// Test
	sut := NewHttpResponseRedirect(expectedUrl)

	// Verify
	if !ExpectNonNil(sut, t) {
		return
	}
	if !ExpectEqual(STATUS_TEMPORARY_REDIRECT, sut.GetStatus(), t) {
		return
	}
	if !ExpectTrue(sut.GetHeaders().Has("location"), t) {
		return
	}
	actual := sut.GetHeaders().Get("location")
	if !ExpectTrue((len(*actual) > 0), t) {
		return
	}
	if !ExpectString(expectedUrl, (*actual)[0], t) {
		return
	}
}

func TestThat_HttpResponse_NewHttpResponseRedirectPermanent_ReturnsPermanentRedirectResponse(t *testing.T) {
	// Setup
	expectedUrl := "https://example.com/permanent-location"

	// Test
	sut := NewHttpResponseRedirectPermanent(expectedUrl)

	// Verify
	if !ExpectNonNil(sut, t) {
		return
	}
	if !ExpectEqual(STATUS_MOVED_PERMANENTLY, sut.GetStatus(), t) {
		return
	}
	actual := sut.GetHeaders().Get("location")
	if !ExpectTrue((len(*actual) > 0), t) {
		return
	}
	if !ExpectString(expectedUrl, (*actual)[0], t) {
		return
	}
}

// Getters/Setters and Converters

func TestThat_HttpResponse_GetBinBody_ReturnsCorrectBinaryData(t *testing.T) {
	// Setup
	body := "test binary data"
	sut := &httpResponse{}
	sut.bodystring = &body

	// Test
	binBody := sut.GetBinBody()

	// Verify
	if !ExpectNonNil(binBody, t) {
		return
	}
	if !ExpectInt(len(body), len(*binBody), t) {
		return
	}
}

func TestThat_HttpResponse_GetBody_ReturnsCorrectStringData(t *testing.T) {
	// Setup
	binData := []byte{72, 101, 108, 108, 111} // "Hello" in ASCII
	sut := &httpResponse{}
	sut.bodybytes = &binData

	// Test
	strBody := sut.GetBody()

	// Verify
	if !ExpectNonNil(strBody, t) {
		return
	}
	if !ExpectString("Hello", *strBody, t) {
		return
	}
}

// Builder Test

func TestThat_HttpResponse_BuilderCreatesCompleteResponse(t *testing.T) {
	// Setup
	builder := NewHttpResponseBuilder()
	expectedStatus := STATUS_CREATED
	expectedBody := "created resource"
	headers := NewHttpHeadersBuilder().
		Set("Content-Type", "application/json").
		Set("Location", "/resource/123").
		GetHttpHeaders()

	// Test
	sut := builder.
		SetStatus(expectedStatus).
		SetBody(&expectedBody).
		SetHeaders(headers).
		GetHttpResponse()

	// Verify
	if !ExpectNonNil(sut, t) {
		return
	}
	if !ExpectEqual(expectedStatus, sut.GetStatus(), t) {
		return
	}
	if !ExpectString(expectedBody, *sut.GetBody(), t) {
		return
	}
	actual := sut.GetHeaders().Get("Content-Type")
	if !ExpectNonNil(actual, t) {
		return
	}
	if !ExpectTrue((len(*actual) > 0), t) {
		return
	}
	if !ExpectString("application/json", (*actual)[0], t) {
		return
	}
	actual = sut.GetHeaders().Get("Location")
	if !ExpectTrue((len(*actual) > 0), t) {
		return
	}
	if !ExpectString("/resource/123", (*actual)[0], t) {
		return
	}
}
