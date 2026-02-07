package version

import (
	"testing"
)

// Factory Function Tests

func TestThat_NewSemVer_ReturnsInstanceForValidVersion(t *testing.T) {
	// Test
	sut := NewSemVer("1.2.3")

	// Verify
	if sut == nil {
		t.Error("Expected non-nil instance")
		return
	}
	if sut.GetVersion() != "1.2.3" {
		t.Errorf("Expected version '1.2.3', got '%s'", sut.GetVersion())
	}
}

func TestThat_NewSemVer_ReturnsNilForEmptyVersion(t *testing.T) {
	// Test
	sut := NewSemVer("")

	// Verify
	if sut != nil {
		t.Error("Expected nil for empty version")
	}
}

func TestThat_NewSemVer_ReturnsNilForInvalidFormat(t *testing.T) {
	// Test with too many parts
	sut := NewSemVer("1.2.3.4")
	if sut != nil {
		t.Error("Expected nil for too many version parts")
	}

	// Test with non-numeric parts
	sut = NewSemVer("abc.0.0")
	if sut != nil {
		t.Error("Expected nil for non-numeric major")
	}

	sut = NewSemVer("1.abc.0")
	if sut != nil {
		t.Error("Expected nil for non-numeric minor")
	}

	sut = NewSemVer("1.0.abc")
	if sut != nil {
		t.Error("Expected nil for non-numeric patch")
	}
}

func TestThat_NewSemVer_HandlesMajorOnly(t *testing.T) {
	// Test
	sut := NewSemVer("5")

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
	if sut.GetVersionPatch() != 0 {
		t.Errorf("Expected patch 0, got %d", sut.GetVersionPatch())
	}
}

func TestThat_NewSemVer_HandlesMajorMinorOnly(t *testing.T) {
	// Test
	sut := NewSemVer("5.3")

	// Verify
	if sut == nil {
		t.Error("Expected non-nil for major.minor version")
		return
	}
	if sut.GetVersionMajor() != 5 {
		t.Errorf("Expected major 5, got %d", sut.GetVersionMajor())
	}
	if sut.GetVersionMinor() != 3 {
		t.Errorf("Expected minor 3, got %d", sut.GetVersionMinor())
	}
	if sut.GetVersionPatch() != 0 {
		t.Errorf("Expected patch 0, got %d", sut.GetVersionPatch())
	}
}

// Interface Tests

func TestThat_SemVer_GetVersion_ReturnsOriginalString(t *testing.T) {
	// Setup
	sut := NewSemVer("2.5.8")

	// Test & Verify
	if sut.GetVersion() != "2.5.8" {
		t.Errorf("Expected '2.5.8', got '%s'", sut.GetVersion())
	}
}

func TestThat_SemVer_GetScheme_ReturnsSEMVER(t *testing.T) {
	// Setup
	sut := NewSemVer("1.0.0")

	// Test & Verify
	if sut.GetScheme() != "SEMVER" {
		t.Errorf("Expected 'SEMVER', got '%s'", sut.GetScheme())
	}
}

func TestThat_SemVer_GetVersionParts_ReturnParsedValues(t *testing.T) {
	// Setup
	sut := NewSemVer("3.7.12")

	// Verify
	if sut.GetVersionMajor() != 3 {
		t.Errorf("Expected major 3, got %d", sut.GetVersionMajor())
	}
	if sut.GetVersionMinor() != 7 {
		t.Errorf("Expected minor 7, got %d", sut.GetVersionMinor())
	}
	if sut.GetVersionPatch() != 12 {
		t.Errorf("Expected patch 12, got %d", sut.GetVersionPatch())
	}
}

// Compare Tests

func TestThat_SemVer_Compare_ReturnsZeroForEqualVersions(t *testing.T) {
	// Setup
	v1 := NewSemVer("1.5.3")
	v2 := NewSemVer("1.5.3")

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

func TestThat_SemVer_Compare_ReturnsNegativeWhenOursIsLower(t *testing.T) {
	// Test cases
	cases := []struct {
		v1, v2 string
	}{
		{"1.0.0", "2.0.0"},    // Major lower
		{"1.0.0", "1.5.0"},    // Minor lower
		{"1.0.0", "1.0.5"},    // Patch lower
		{"0.9.9", "1.0.0"},    // All lower
		{"1.1.9", "1.2.0"},    // Minor bump beats patch
	}

	for _, tc := range cases {
		v1 := NewSemVer(tc.v1)
		v2 := NewSemVer(tc.v2)

		result, err := v1.Compare(v2)
		if err != nil {
			t.Errorf("Unexpected error for %s vs %s: %v", tc.v1, tc.v2, err)
		}
		if result != -1 {
			t.Errorf("Expected -1 for %s < %s, got %d", tc.v1, tc.v2, result)
		}
	}
}

func TestThat_SemVer_Compare_ReturnsPositiveWhenOursIsHigher(t *testing.T) {
	// Test cases
	cases := []struct {
		v1, v2 string
	}{
		{"2.0.0", "1.0.0"},    // Major higher
		{"1.5.0", "1.0.0"},    // Minor higher
		{"1.0.5", "1.0.0"},    // Patch higher
		{"2.0.0", "1.9.9"},    // Major bump beats all
	}

	for _, tc := range cases {
		v1 := NewSemVer(tc.v1)
		v2 := NewSemVer(tc.v2)

		result, err := v1.Compare(v2)
		if err != nil {
			t.Errorf("Unexpected error for %s vs %s: %v", tc.v1, tc.v2, err)
		}
		if result != 1 {
			t.Errorf("Expected 1 for %s > %s, got %d", tc.v1, tc.v2, result)
		}
	}
}

func TestThat_SemVer_Compare_ReturnsErrorForMismatchedScheme(t *testing.T) {
	// Setup
	v1 := NewSemVer("1.0.0")
	v2 := NewMajorMinor("1.0")

	// Test
	_, err := v1.Compare(v2)

	// Verify
	if err == nil {
		t.Error("Expected error for mismatched scheme")
	}
}

// Nil safety

func TestThat_SemVer_GetVersion_ReturnsEmptyForNilReceiver(t *testing.T) {
	var sut *semver = nil
	
	if sut.GetVersion() != "" {
		t.Error("Expected empty string for nil receiver")
	}
}
