package http

import (
	"testing"

	. "github.com/DigiStratum/GoLib/Testing"
)

// Factory Function

func TestThat_HttpStatus_HttpStatusFromCode_ReturnsCorrectStatuses(t *testing.T) {
	// Test a sampling of status codes
	if !ExpectEqual(STATUS_OK, HttpStatusFromCode(200), t) {
		return
	}
	if !ExpectEqual(STATUS_CREATED, HttpStatusFromCode(201), t) {
		return
	}
	if !ExpectEqual(STATUS_NOT_FOUND, HttpStatusFromCode(404), t) {
		return
	}
	if !ExpectEqual(STATUS_INTERNAL_SERVER_ERROR, HttpStatusFromCode(500), t) {
		return
	}
}

func TestThat_HttpStatus_HttpStatusFromCode_ReturnsUnknownForInvalidCode(t *testing.T) {
	// Test an invalid code
	if !ExpectEqual(STATUS_UNKNOWN, HttpStatusFromCode(999), t) {
		return
	}
}

// GetHttpStatusCode

func TestThat_HttpStatus_GetHttpStatusCode_ReturnsCorrectCodes(t *testing.T) {
	// Test a sampling of statuses
	if !ExpectInt(200, STATUS_OK.GetHttpStatusCode(), t) {
		return
	}
	if !ExpectInt(201, STATUS_CREATED.GetHttpStatusCode(), t) {
		return
	}
	if !ExpectInt(404, STATUS_NOT_FOUND.GetHttpStatusCode(), t) {
		return
	}
	if !ExpectInt(500, STATUS_INTERNAL_SERVER_ERROR.GetHttpStatusCode(), t) {
		return
	}
}

func TestThat_HttpStatus_GetHttpStatusCode_ReturnsZeroForUnknown(t *testing.T) {
	// Test unknown status
	if !ExpectInt(0, STATUS_UNKNOWN.GetHttpStatusCode(), t) {
		return
	}
}

// ToString

func TestThat_HttpStatus_ToString_ReturnsCorrectStrings(t *testing.T) {
	// Test a sampling of statuses
	if !ExpectString("OK", STATUS_OK.ToString(), t) {
		return
	}
	if !ExpectString("CREATED", STATUS_CREATED.ToString(), t) {
		return
	}
	if !ExpectString("NOT FOUND", STATUS_NOT_FOUND.ToString(), t) {
		return
	}
	if !ExpectString("INTERNAL SERVER ERROR", STATUS_INTERNAL_SERVER_ERROR.ToString(), t) {
		return
	}
}

func TestThat_HttpStatus_ToString_ReturnsUnknownForInvalidStatus(t *testing.T) {
	// Test unknown status
	if !ExpectString("UNKNOWN STATUS CODE", STATUS_UNKNOWN.ToString(), t) {
		return
	}
}

// Status category tests

func TestThat_HttpStatus_IsStatus1xx_ReturnsTrueFor1xxStatusesOnly(t *testing.T) {
	// Test 1xx
	if !ExpectTrue(STATUS_CONTINUE.IsStatus1xx(), t) {
		return
	}
	if !ExpectTrue(STATUS_SWITCHING_PROTOCOLS.IsStatus1xx(), t) {
		return
	}

	// Test non-1xx
	if !ExpectFalse(STATUS_OK.IsStatus1xx(), t) {
		return
	}
	if !ExpectFalse(STATUS_NOT_FOUND.IsStatus1xx(), t) {
		return
	}
}

func TestThat_HttpStatus_IsStatus2xx_ReturnsTrueFor2xxStatusesOnly(t *testing.T) {
	// Test 2xx
	if !ExpectTrue(STATUS_OK.IsStatus2xx(), t) {
		return
	}
	if !ExpectTrue(STATUS_CREATED.IsStatus2xx(), t) {
		return
	}
	if !ExpectTrue(STATUS_ACCEPTED.IsStatus2xx(), t) {
		return
	}

	// Test non-2xx
	if !ExpectFalse(STATUS_CONTINUE.IsStatus2xx(), t) {
		return
	}
	if !ExpectFalse(STATUS_NOT_FOUND.IsStatus2xx(), t) {
		return
	}
}

func TestThat_HttpStatus_IsStatus3xx_ReturnsTrueFor3xxStatusesOnly(t *testing.T) {
	// Test 3xx
	if !ExpectTrue(STATUS_MOVED_PERMANENTLY.IsStatus3xx(), t) {
		return
	}
	if !ExpectTrue(STATUS_FOUND.IsStatus3xx(), t) {
		return
	}
	if !ExpectTrue(STATUS_TEMPORARY_REDIRECT.IsStatus3xx(), t) {
		return
	}

	// Test non-3xx
	if !ExpectFalse(STATUS_OK.IsStatus3xx(), t) {
		return
	}
	if !ExpectFalse(STATUS_NOT_FOUND.IsStatus3xx(), t) {
		return
	}
}

func TestThat_HttpStatus_IsStatus4xx_ReturnsTrueFor4xxStatusesOnly(t *testing.T) {
	// Test 4xx
	if !ExpectTrue(STATUS_BAD_REQUEST.IsStatus4xx(), t) {
		return
	}
	if !ExpectTrue(STATUS_NOT_FOUND.IsStatus4xx(), t) {
		return
	}
	if !ExpectTrue(STATUS_UNAUTHORIZED.IsStatus4xx(), t) {
		return
	}

	// Test non-4xx
	if !ExpectFalse(STATUS_OK.IsStatus4xx(), t) {
		return
	}
	if !ExpectFalse(STATUS_INTERNAL_SERVER_ERROR.IsStatus4xx(), t) {
		return
	}
}

func TestThat_HttpStatus_IsStatus5xx_ReturnsTrueFor5xxStatusesOnly(t *testing.T) {
	// Test 5xx
	if !ExpectTrue(STATUS_INTERNAL_SERVER_ERROR.IsStatus5xx(), t) {
		return
	}
	if !ExpectTrue(STATUS_SERVICE_UNAVAILABLE.IsStatus5xx(), t) {
		return
	}
	if !ExpectTrue(STATUS_GATEWAY_TIMEOUT.IsStatus5xx(), t) {
		return
	}

	// Test non-5xx
	if !ExpectFalse(STATUS_OK.IsStatus5xx(), t) {
		return
	}
	if !ExpectFalse(STATUS_NOT_FOUND.IsStatus5xx(), t) {
		return
	}
}
