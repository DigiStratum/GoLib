package http

import (
	"testing"

	. "github.com/DigiStratum/GoLib/Testing"
)

func TestThat_RequestContext_NewRequestContext_ReturnsNonNil(t *testing.T) {
	// Test
	ctx := NewRequestContext()

	// Verify
	ExpectNonNil(ctx, t)
}

func TestThat_RequestContext_GetSetRequestId_Works(t *testing.T) {
	// Setup
	ctx := NewRequestContext()

	// Test
	ctx.SetRequestId("test-uuid-123")
	actual := ctx.GetRequestId()

	// Verify
	ExpectString("test-uuid-123", actual, t)
}

func TestThat_RequestContext_GetSetServerPath_Works(t *testing.T) {
	// Setup
	ctx := NewRequestContext()

	// Test
	ctx.SetServerPath("/api")
	actual := ctx.GetServerPath()

	// Verify
	ExpectString("/api", actual, t)
}

func TestThat_RequestContext_GetSetModulePath_Works(t *testing.T) {
	// Setup
	ctx := NewRequestContext()

	// Test
	ctx.SetModulePath("/users")
	actual := ctx.GetModulePath()

	// Verify
	ExpectString("/users", actual, t)
}

func TestThat_RequestContext_GetSetPrefixPath_Works(t *testing.T) {
	// Setup
	ctx := NewRequestContext()

	// Test
	ctx.SetPrefixPath("/api/users")
	actual := ctx.GetPrefixPath()

	// Verify
	ExpectString("/api/users", actual, t)
}

func TestThat_RequestContext_GetSetModuleId_Works(t *testing.T) {
	// Setup
	ctx := NewRequestContext()

	// Test
	ctx.SetModuleId("module-abc")
	actual := ctx.GetModuleId()

	// Verify
	ExpectString("module-abc", actual, t)
}

func TestThat_HttpRequest_GetContext_ReturnsNonNilLazily(t *testing.T) {
	// Setup
	request := NewHttpRequestBuilder(METHOD_GET, "http://example.com/test").GetHttpRequest()

	// Test - GetContext should lazily initialize and return non-nil
	ctx := request.GetContext()

	// Verify
	ExpectNonNil(ctx, t)
}

func TestThat_HttpRequest_SetContext_Works(t *testing.T) {
	// Setup
	request := NewHttpRequestBuilder(METHOD_GET, "http://example.com/test").GetHttpRequest()
	ctx := NewRequestContext()
	ctx.SetRequestId("my-request-id")

	// Test
	request.SetContext(ctx)
	actualCtx := request.GetContext()

	// Verify
	ExpectString("my-request-id", actualCtx.GetRequestId(), t)
}

func TestThat_HttpRequestBuilder_SetContext_Works(t *testing.T) {
	// Setup
	ctx := NewRequestContext()
	ctx.SetRequestId("builder-context-id")

	// Test
	request := NewHttpRequestBuilder(METHOD_POST, "http://example.com/api").
		SetContext(ctx).
		GetHttpRequest()

	// Verify
	ExpectString("builder-context-id", request.GetContext().GetRequestId(), t)
}

func TestThat_HttpRequest_GetBuilder_PreservesContext(t *testing.T) {
	// Setup
	ctx := NewRequestContext()
	ctx.SetRequestId("original-id")
	ctx.SetServerPath("/server")
	ctx.SetModulePath("/module")

	request := NewHttpRequestBuilder(METHOD_GET, "http://example.com/test").
		SetContext(ctx).
		GetHttpRequest()

	// Test - GetBuilder should preserve the context
	builder := request.GetBuilder()
	rebuilt := builder.GetHttpRequest()

	// Verify
	ExpectString("original-id", rebuilt.GetContext().GetRequestId(), t)
	ExpectString("/server", rebuilt.GetContext().GetServerPath(), t)
	ExpectString("/module", rebuilt.GetContext().GetModulePath(), t)
}
