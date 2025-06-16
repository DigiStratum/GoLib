package http

import (
	"testing"

	. "github.com/DigiStratum/GoLib/Testing"
)

// Factory Function

func TestThat_HttpHeaders_NewHttpHeaders_ReturnsInstance(t *testing.T) {
	// Setup & Test
	var sut HttpHeadersIfc = &httpHeaders{} // Verifies that result satisfies IFC

	// Verify
	if !ExpectNonNil(sut, t) {
		return
	}
	if !ExpectTrue(sut.IsEmpty(), t) {
		return
	}
}

// Has, IsEmpty, Size

func TestThat_HttpHeaders_Has_ReturnsFalseForMissingHeader(t *testing.T) {
	// Setup
	sut := &httpHeaders{}

	// Test & Verify
	if !ExpectFalse(sut.Has("Content-Type"), t) {
		return
	}
}

func TestThat_HttpHeaders_Has_ReturnsTrueForExistingHeader(t *testing.T) {
	// Setup
	sut := NewHttpHeadersBuilder().
		Set("Content-Type", "application/json").
		GetHttpHeaders()

	// Test & Verify
	if !ExpectTrue(sut.Has("Content-Type"), t) {
		return
	}
}

func TestThat_HttpHeaders_IsEmpty_ReturnsTrueWhenEmpty(t *testing.T) {
	// Setup
	sut := &httpHeaders{}

	// Test & Verify
	if !ExpectTrue(sut.IsEmpty(), t) {
		return
	}
}

func TestThat_HttpHeaders_IsEmpty_ReturnsFalseWhenNotEmpty(t *testing.T) {
	// Setup
	sut := NewHttpHeadersBuilder().
		Set("Content-Type", "application/json").
		GetHttpHeaders()

	// Test & Verify
	if !ExpectFalse(sut.IsEmpty(), t) {
		return
	}
}

func TestThat_HttpHeaders_Size_ReturnsCorrectSizeForEmptyHeaders(t *testing.T) {
	// Setup
	sut := &httpHeaders{}

	// Test
	size := sut.Size()

	// Verify - empty headers return 2 for the final double newline
	if !ExpectInt(2, size, t) {
		return
	}
}

func TestThat_HttpHeaders_Size_ReturnsCorrectSizeForSingleHeader(t *testing.T) {
	// Setup
	sut := NewHttpHeadersBuilder().
		Set("X-Test", "value").
		GetHttpHeaders()

	// Test
	size := sut.Size()

	// The way Size() calculates it's specifically:
	// len("value") + 4 (for colon-space and space+semicolon) + 1 (for newline) (plus a final newline for the last one!)
	// which gives 5 + 4 + 1 = 10
	if !ExpectInt(11, size, t) {
		return
	}
}

func TestThat_HttpHeaders_Size_ReturnsCorrectSizeForMultipleHeaders(t *testing.T) {
	// Setup
	sut := NewHttpHeadersBuilder().
		Set("X-Test", "value").
		Set("Content-Type", "application/json").
		GetHttpHeaders()

	// Test
	size := sut.Size()

	// Using Size() calculation method:
	// "value" = 5 + 4 + 1 = 10
	// "application/json" = 16 + 4 + 1 = 21 (plus a final newline for the last one!)
	// Total = 32
	if !ExpectInt(32, size, t) {
		return
	}
}

// GetNames, Get

func TestThat_HttpHeaders_GetNames_ReturnsAllHeaderNames(t *testing.T) {
	// Setup
	sut := NewHttpHeadersBuilder().
		Set("X-Test", "value").
		Set("Content-Type", "application/json").
		GetHttpHeaders()

	// Test
	names := sut.GetNames()

	// Verify
	if !ExpectNonNil(names, t) {
		return
	}
	if !ExpectInt(2, len(*names), t) {
		return
	}

	// Check both names exist (order might vary)
	found1, found2 := false, false
	for _, name := range *names {
		if name == "X-Test" {
			found1 = true
		}
		if name == "Content-Type" {
			found2 = true
		}
	}
	if !ExpectTrue(found1, t) {
		return
	}
	if !ExpectTrue(found2, t) {
		return
	}
}

func TestThat_HttpHeaders_Get_ReturnsNilForMissingHeader(t *testing.T) {
	// Setup
	sut := &httpHeaders{}

	// Test
	values := sut.Get("X-Missing")

	// Verify
	if !ExpectNil(values, t) {
		return
	}
}

func TestThat_HttpHeaders_Get_ReturnsValuesForExistingHeader(t *testing.T) {
	// Setup
	sut := NewHttpHeadersBuilder().
		Set("X-Test", "value").
		GetHttpHeaders()

	// Test
	values := sut.Get("X-Test")

	// Verify
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

// ToMap

func TestThat_HttpHeaders_ToMap_ReturnsMapCopy(t *testing.T) {
	// Setup
	sut := NewHttpHeadersBuilder().
		Set("X-Test", "value").
		GetHttpHeaders()

	// Test
	headerMap := sut.ToMap()

	// Verify
	if !ExpectNonNil(headerMap, t) {
		return
	}
	if !ExpectInt(1, len(*headerMap), t) {
		return
	}

	values, ok := (*headerMap)["X-Test"]
	if !ExpectTrue(ok, t) {
		return
	}
	if !ExpectInt(1, len(values), t) {
		return
	}
	if !ExpectString("value", values[0], t) {
		return
	}
}

// GetAcceptableLanguages

func TestThat_HttpHeaders_GetAcceptableLanguages_ReturnsNilWhenNoLanguages(t *testing.T) {
	// Setup
	sut := &httpHeaders{}

	// Test
	languages := sut.GetAcceptableLanguages()

	// Verify
	if !ExpectNil(languages, t) {
		return
	}
}

func TestThat_HttpHeaders_GetAcceptableLanguages_ReturnsLanguagesList(t *testing.T) {
	// Setup
	acceptLanguageValue := "en-US,en;q=0.9,fr;q=0.8"

	sut := NewHttpHeadersBuilder().
		Set("Accept-Language", acceptLanguageValue).
		GetHttpHeaders()

	// Test
	languages := sut.GetAcceptableLanguages()

	// Verify
	if !ExpectNonNil(languages, t) {
		return
	}
	if !ExpectInt(3, len(*languages), t) {
		return
	}
	if !ExpectString("en-US", (*languages)[0], t) {
		return
	}
	if !ExpectString("en", (*languages)[1], t) {
		return
	}
	if !ExpectString("fr", (*languages)[2], t) {
		return
	}
}

// Builder Tests

func TestThat_HttpHeaders_BuilderCreatesHeaders(t *testing.T) {
	// Setup & Test
	sut := NewHttpHeadersBuilder().
		Set("X-Test", "value").
		Set("Content-Type", "application/json").
		GetHttpHeaders()

	// Verify
	if !ExpectNonNil(sut, t) {
		return
	}
	if !ExpectFalse(sut.IsEmpty(), t) {
		return
	}
	if !ExpectTrue(sut.Has("X-Test"), t) {
		return
	}
	if !ExpectTrue(sut.Has("Content-Type"), t) {
		return
	}
}

func TestThat_HttpHeaders_BuilderAddsMultipleValues(t *testing.T) {
	// Setup & Test
	sut := NewHttpHeadersBuilder().
		Set("X-Test", "value1").
		Set("X-Test", "value2").
		GetHttpHeaders()

	// Verify
	if !ExpectNonNil(sut, t) {
		return
	}
	values := sut.Get("X-Test")
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
