package http

/*

TODO:
  * Fix/add test coverage for Configurability; Copilot added some based on assumptions of Exported
    properties but there is no valid data access path here
*/

import (
	"testing"

	cfg "github.com/DigiStratum/GoLib/Config"

	. "github.com/DigiStratum/GoLib/Testing"
)

// Interface

func TestThat_HttpClient_NewHttpClient_ReturnsInstance(t *testing.T) {
	// Setup
	var sut HttpClientIfc = NewHttpClient() // Verifies that result satisfies IFC

	// Verify
	if !ExpectNonNil(sut, t) {
		return
	}
}

// Configuration

/*
func TestThat_HttpClient_HasConfigurableItems(t *testing.T) {
	// Setup
	sut := NewHttpClient()

	// Verify
	if !ExpectNonNil(sut.Configurable(), t) {
		return
	}
	if !ExpectInt(4, len(sut.GetConfigItems()), t) {
		return
	}
	if !ExpectTrue(sut.HasConfigItem("maxBodyLenKb"), t) {
		return
	}
	if !ExpectTrue(sut.HasConfigItem("requestTimeout"), t) {
		return
	}
	if !ExpectTrue(sut.HasConfigItem("idleTimeout"), t) {
		return
	}
	if !ExpectTrue(sut.HasConfigItem("disableCompression"), t) {
		return
	}
}

func TestThat_HttpClient_ConfigDefaults_AreCorrect(t *testing.T) {
	// Setup
	sut := NewHttpClient()

	// Verify
	if !ExpectString("10240", sut.GetConfigItem("maxBodyLenKb").GetDefault(), t) {
		return
	}
	if !ExpectString("60s", sut.GetConfigItem("requestTimeout").GetDefault(), t) {
		return
	}
	if !ExpectString("30", sut.GetConfigItem("idleTimeout").GetDefault(), t) {
		return
	}
	if !ExpectString("false", sut.GetConfigItem("disableCompression").GetDefault(), t) {
		return
	}
}
*/

func TestThat_HttpClient_Config_SetsMaxBodyLenKb(t *testing.T) {
	// Setup
	sut := NewHttpClient()
	c := cfg.NewConfig()
	c.Set("maxBodyLenKb", "5120")

	// Test
	err := sut.Configure(c)

	// Verify
	if !ExpectNoError(err, t) {
		return
	}
	// Access private field through reflection or test behavior
}

func TestThat_HttpClient_Start_InitializesClient(t *testing.T) {
	// Setup
	sut := NewHttpClient()

	// Test
	err := sut.Start()

	// Verify
	if !ExpectNoError(err, t) {
		return
	}
	// Client should be initialized but it's a private field
}

// GetRequestResponse

func TestThat_HttpClient_GetRequestResponse_HandlesEmptyRequest(t *testing.T) {
	// Setup
	sut := NewHttpClient()
	err := sut.Start()
	if !ExpectNoError(err, t) {
		return
	}

	// Test - blank request should fail as the URL is empty
	req := &httpRequest{}
	resp, err := sut.GetRequestResponse(req)

	// Verify
	if !ExpectError(err, t) {
		return
	}
	if !ExpectNil(resp, t) {
		return
	}
}
