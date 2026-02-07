package version

import (
	"testing"
)

// Factory Function Tests

func TestThat_NewMajorMinor_ReturnsInstanceForValidVersion(t *testing.T) {
	// Test
	sut := NewMajorMinor("1.0")

	// Verify
	if sut == nil {
		t.Error("Expected non-nil instance")
		return
	}
	if sut.GetVersion() != "1.0" {
		t.Errorf("Expected version '1.0', got '%s'", sut.GetVersion())
	}
}

func TestThat_NewMajorMinor_ReturnsNilForEmptyVersion(t *testing.T) {
	// Test
	sut := NewMajorMinor("")

	// Verify
	if sut != nil {
		t.Error("Expected nil for empty version")
	}
}

func TestThat_NewMajorMinor_ReturnsNilForInvalidFormat(t *testing.T) {
	// Test with too many parts
	sut := NewMajorMinor("1.2.3.4")
	if sut != nil {
		t.Error("Expected nil for too many version parts")
	}

	// Test with non-numeric major
	sut = NewMajorMinor("abc.0")
	if sut != nil {
		t.Error("Expected nil for non-numeric major")
	}

	// Test with non-numeric minor
	sut = NewMajorMinor("1.abc")
	if sut != nil {
		t.Error("Expected nil for non-numeric minor")
	}
}

func TestThat_NewMajorMinor_HandlesMajorOnly(t *testing.T) {
	// Test
	sut := NewMajorMinor("5")

	// Verify
	if sut == nil {
		t.Error("Expected non-nil for major-only version")
		return
	}
	if sut.GetVersionMajor() != 5 {
		t.Errorf("Expected major 5, got %d", sut.GetVersionMajor())
	}
	if sut.GetVersionMinor() != 0 {
		t.Errorf("Expected minor 0, got %d", sut.GetVersionMinor())
	}
}

// Interface Tests

func TestThat_MajorMinor_GetVersion_ReturnsOriginalString(t *testing.T) {
	// Setup
	sut := NewMajorMinor("2.5")

	// Test & Verify
	if sut.GetVersion() != "2.5" {
		t.Errorf("Expected '2.5', got '%s'", sut.GetVersion())
	}
}

func TestThat_MajorMinor_GetScheme_ReturnsMAJMIN(t *testing.T) {
	// Setup
	sut := NewMajorMinor("1.0")

	// Test & Verify
	if sut.GetScheme() != "MAJMIN" {
		t.Errorf("Expected 'MAJMIN', got '%s'", sut.GetScheme())
	}
}

func TestThat_MajorMinor_GetVersionMajor_ReturnsParsedMajor(t *testing.T) {
	// Setup
	sut := NewMajorMinor("3.7")

	// Test & Verify
	if sut.GetVersionMajor() != 3 {
		t.Errorf("Expected major 3, got %d", sut.GetVersionMajor())
	}
}

func TestThat_MajorMinor_GetVersionMinor_ReturnsParsedMinor(t *testing.T) {
	// Setup
	sut := NewMajorMinor("3.7")

	// Test & Verify
	if sut.GetVersionMinor() != 7 {
		t.Errorf("Expected minor 7, got %d", sut.GetVersionMinor())
	}
}

// Compare Tests

func TestThat_MajorMinor_Compare_ReturnsZeroForEqualVersions(t *testing.T) {
	// Setup
	v1 := NewMajorMinor("1.5")
	v2 := NewMajorMinor("1.5")

	// Test
	result, err := v1.Compare(v2)

	// Verify
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if result != 0 {
		t.Errorf("Expected 0 for equal versions, got %d", result)
	}
}

func TestThat_MajorMinor_Compare_ReturnsNegativeWhenOursIsLower(t *testing.T) {
	// Test cases
	cases := []struct {
		v1, v2 string
	}{
		{"1.0", "2.0"},  // Major lower
		{"1.0", "1.5"},  // Minor lower
		{"0.9", "1.0"},  // Both lower
	}

	for _, tc := range cases {
		v1 := NewMajorMinor(tc.v1)
		v2 := NewMajorMinor(tc.v2)

		result, err := v1.Compare(v2)
		if err != nil {
			t.Errorf("Unexpected error for %s vs %s: %v", tc.v1, tc.v2, err)
		}
		if result != -1 {
			t.Errorf("Expected -1 for %s < %s, got %d", tc.v1, tc.v2, result)
		}
	}
}

func TestThat_MajorMinor_Compare_ReturnsPositiveWhenOursIsHigher(t *testing.T) {
	// Test cases
	cases := []struct {
		v1, v2 string
	}{
		{"2.0", "1.0"},  // Major higher
		{"1.5", "1.0"},  // Minor higher
		{"2.1", "1.9"},  // Both considerations
	}

	for _, tc := range cases {
		v1 := NewMajorMinor(tc.v1)
		v2 := NewMajorMinor(tc.v2)

		result, err := v1.Compare(v2)
		if err != nil {
			t.Errorf("Unexpected error for %s vs %s: %v", tc.v1, tc.v2, err)
		}
		if result != 1 {
			t.Errorf("Expected 1 for %s > %s, got %d", tc.v1, tc.v2, result)
		}
	}
}

func TestThat_MajorMinor_Compare_ReturnsErrorForMismatchedScheme(t *testing.T) {
	// Setup
	v1 := NewMajorMinor("1.0")
	v2 := NewSemVer("1.0.0")

	// Test
	_, err := v1.Compare(v2)

	// Verify
	if err == nil {
		t.Error("Expected error for mismatched scheme")
	}
}

// Nil safety

func TestThat_MajorMinor_GetVersion_ReturnsEmptyForNilReceiver(t *testing.T) {
	var sut *majmin = nil
	
	if sut.GetVersion() != "" {
		t.Error("Expected empty string for nil receiver")
	}
}
