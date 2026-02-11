package serializable

import (
	"testing"

	. "github.com/DigiStratum/GoLib/Testing"
	enc "github.com/DigiStratum/GoLib/Data/transcoder/encodingscheme"
	xc "github.com/DigiStratum/GoLib/Data/transcoder"
)

// -------------------------------------------------------------------------------------------------
// Factory Function Tests
// -------------------------------------------------------------------------------------------------

func TestThat_NewSerializer_ReturnsNewSerializer_WithTranscoder(t *testing.T) {
	// Setup
	transcoder := xc.NewTranscoder()
	transcoder.SetEncoderScheme(enc.NewEncodingSchemeBase64())
	transcoder.SetDecoderScheme(enc.NewEncodingSchemeBase64())

	// Test
	sut := NewSerializer(transcoder)

	// Verify
	ExpectNonNil(sut, t)
}

func TestThat_NewSerializer_ReturnsNewSerializer_WithNilTranscoder(t *testing.T) {
	// Test
	sut := NewSerializer(nil)

	// Verify
	ExpectNonNil(sut, t)
}

// -------------------------------------------------------------------------------------------------
// Serialize Method Tests
// -------------------------------------------------------------------------------------------------

func TestThat_Serializer_Serialize_ReturnsSerializedString_WithValidInput(t *testing.T) {
	// Setup
	transcoder := xc.NewTranscoder()
	transcoder.SetEncoderScheme(enc.NewEncodingSchemeBase64())
	transcoder.SetDecoderScheme(enc.NewEncodingSchemeBase64())
	sut := NewSerializer(transcoder)
	data := `{"key":"value"}`
	typeName := "Object"

	// Test
	result, err := sut.Serialize(&data, typeName)

	// Verify
	ExpectNoError(err, t)
	ExpectNonNil(result, t)
	if nil == result { return }
	if len(*result) == 0 {
		t.Error("Expected non-empty serialized string")
	}
}

func TestThat_Serializer_Serialize_ReturnsError_WithNilTranscoder(t *testing.T) {
	// Setup
	sut := NewSerializer(nil)
	data := `{"key":"value"}`
	typeName := "Object"

	// Test
	result, err := sut.Serialize(&data, typeName)

	// Verify
	ExpectError(err, t)
	ExpectNil(result, t)
}

func TestThat_Serializer_Serialize_ReturnsError_WithWrongEncodingScheme(t *testing.T) {
	// Setup
	transcoder := xc.NewTranscoder()
	transcoder.SetEncoderScheme(enc.NewEncodingSchemeHTMLEscape())  // Wrong scheme
	transcoder.SetDecoderScheme(enc.NewEncodingSchemeHTMLEscape())
	sut := NewSerializer(transcoder)
	data := `{"key":"value"}`
	typeName := "Object"

	// Test
	result, err := sut.Serialize(&data, typeName)

	// Verify
	ExpectError(err, t)
	ExpectNil(result, t)
}

func TestThat_Serializer_Serialize_HandlesEmptyString(t *testing.T) {
	// Setup
	transcoder := xc.NewTranscoder()
	transcoder.SetEncoderScheme(enc.NewEncodingSchemeBase64())
	transcoder.SetDecoderScheme(enc.NewEncodingSchemeBase64())
	sut := NewSerializer(transcoder)
	data := ""
	typeName := "Object"

	// Test
	result, err := sut.Serialize(&data, typeName)

	// Verify
	ExpectNoError(err, t)
	ExpectNonNil(result, t)
}

func TestThat_Serializer_Serialize_HandlesNilData(t *testing.T) {
	// Setup
	transcoder := xc.NewTranscoder()
	transcoder.SetEncoderScheme(enc.NewEncodingSchemeBase64())
	transcoder.SetDecoderScheme(enc.NewEncodingSchemeBase64())
	sut := NewSerializer(transcoder)
	typeName := "Object"

	// Test
	result, err := sut.Serialize(nil, typeName)

	// Verify
	ExpectError(err, t)
	ExpectNil(result, t)
}

func TestThat_Serializer_Serialize_HandlesSpecialCharacters(t *testing.T) {
	// Setup
	transcoder := xc.NewTranscoder()
	transcoder.SetEncoderScheme(enc.NewEncodingSchemeBase64())
	transcoder.SetDecoderScheme(enc.NewEncodingSchemeBase64())
	sut := NewSerializer(transcoder)
	
	testCases := []struct {
		name string
		data string
	}{
		{"newlines", "line1\nline2\nline3"},
		{"tabs", "col1\tcol2\tcol3"},
		{"quotes", `{"message":"It's \"quoted\""}`},
		{"unicode", "Hello 世界 🌍"},
		{"special chars", `!@#$%^&*()_+-={}[]|\:;"'<>,.?/~` + "`"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Test
			result, err := sut.Serialize(&tc.data, "Object")

			// Verify
			ExpectNoError(err, t)
			ExpectNonNil(result, t)
		})
	}
}

// -------------------------------------------------------------------------------------------------
// Deserialize Method Tests
// -------------------------------------------------------------------------------------------------

func TestThat_Serializer_Deserialize_ReturnsError_WithNilData(t *testing.T) {
	// Setup
	transcoder := xc.NewTranscoder()
	transcoder.SetEncoderScheme(enc.NewEncodingSchemeBase64())
	transcoder.SetDecoderScheme(enc.NewEncodingSchemeBase64())
	sut := NewSerializer(transcoder)

	// Test
	result, err := sut.Deserialize(nil)

	// Verify
	ExpectError(err, t)
	ExpectNil(result, t)
}

func TestThat_Serializer_Deserialize_ReturnsError_WithMalformedData(t *testing.T) {
	// Setup
	transcoder := xc.NewTranscoder()
	transcoder.SetEncoderScheme(enc.NewEncodingSchemeBase64())
	transcoder.SetDecoderScheme(enc.NewEncodingSchemeBase64())
	sut := NewSerializer(transcoder)

	testCases := []struct {
		name string
		data string
	}{
		{"empty string", ""},
		{"random text", "just random text"},
		{"missing prefix", "[j64:T2JqZWN0:eyJrZXkiOiJ2YWx1ZSJ9]"},
		{"missing suffix", "ser[j64:T2JqZWN0:eyJrZXkiOiJ2YWx1ZSJ9"},
		{"malformed format", "ser[invalid]"},
		{"missing method", "ser[:T2JqZWN0:eyJrZXkiOiJ2YWx1ZSJ9]"},
		{"missing type", "ser[j64::eyJrZXkiOiJ2YWx1ZSJ9]"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Test
			result, err := sut.Deserialize(&tc.data)

			// Verify
			ExpectError(err, t)
			ExpectNil(result, t)
		})
	}
}

func TestThat_Serializer_Deserialize_ReturnsError_WithUnsupportedMethod(t *testing.T) {
	// Setup
	transcoder := xc.NewTranscoder()
	transcoder.SetEncoderScheme(enc.NewEncodingSchemeBase64())
	transcoder.SetDecoderScheme(enc.NewEncodingSchemeBase64())
	sut := NewSerializer(transcoder)
	data := "ser[xyz:T2JqZWN0:eyJrZXkiOiJ2YWx1ZSJ9]"

	// Test
	result, err := sut.Deserialize(&data)

	// Verify
	ExpectError(err, t)
	ExpectNil(result, t)
}

func TestThat_Serializer_Deserialize_ReturnsError_WithNilTranscoder(t *testing.T) {
	// Setup
	sut := NewSerializer(nil)
	data := "ser[j64:T2JqZWN0:eyJrZXkiOiJ2YWx1ZSJ9]"

	// Test
	result, err := sut.Deserialize(&data)

	// Verify
	ExpectError(err, t)
	ExpectNil(result, t)
}

func TestThat_Serializer_Deserialize_ReturnsError_WithWrongDecoderScheme(t *testing.T) {
	// Setup
	transcoder := xc.NewTranscoder()
	transcoder.SetEncoderScheme(enc.NewEncodingSchemeHTMLEscape())
	transcoder.SetDecoderScheme(enc.NewEncodingSchemeHTMLEscape())  // Wrong scheme
	sut := NewSerializer(transcoder)
	data := "ser[j64:T2JqZWN0:eyJrZXkiOiJ2YWx1ZSJ9]"

	// Test
	result, err := sut.Deserialize(&data)

	// Verify
	ExpectError(err, t)
	ExpectNil(result, t)
}

func TestThat_Serializer_Deserialize_ReturnsError_WithInvalidBase64(t *testing.T) {
	// Setup
	transcoder := xc.NewTranscoder()
	transcoder.SetEncoderScheme(enc.NewEncodingSchemeBase64())
	transcoder.SetDecoderScheme(enc.NewEncodingSchemeBase64())
	sut := NewSerializer(transcoder)
	data := "ser[j64:InvalidBase64!!!:InvalidBase64!!!]"

	// Test
	result, err := sut.Deserialize(&data)

	// Verify
	ExpectError(err, t)
	ExpectNil(result, t)
}

func TestThat_Serializer_Deserialize_ReturnsError_WithWrongTypeName(t *testing.T) {
	// Setup
	transcoder := xc.NewTranscoder()
	transcoder.SetEncoderScheme(enc.NewEncodingSchemeBase64())
	transcoder.SetDecoderScheme(enc.NewEncodingSchemeBase64())
	sut := NewSerializer(transcoder)
	// "WrongType" in base64 instead of "Object"
	wrongType := "V3JvbmdUeXBl"
	data := "ser[j64:" + wrongType + ":eyJrZXkiOiJ2YWx1ZSJ9]"

	// Test
	result, err := sut.Deserialize(&data)

	// Verify
	ExpectError(err, t)
	ExpectNil(result, t)
}

// -------------------------------------------------------------------------------------------------
// Roundtrip Tests (Serialize then Deserialize)
// These tests validate the happy path for both serialization and deserialization by ensuring
// that data can be successfully serialized and then deserialized back to its original form.
// -------------------------------------------------------------------------------------------------

func TestThat_Serializer_RoundTrip_PreservesData(t *testing.T) {
	// Setup
	transcoder := xc.NewTranscoder()
	transcoder.SetEncoderScheme(enc.NewEncodingSchemeBase64())
	transcoder.SetDecoderScheme(enc.NewEncodingSchemeBase64())
	sut := NewSerializer(transcoder)

	testCases := []struct {
		name     string
		data     string
		typeName string
	}{
		{"simple json", `{"key":"value"}`, "Object"},
		{"complex json", `{"name":"test","count":123,"active":true}`, "Object"},
		{"empty json", `{}`, "Object"},
		{"json array", `["a","b","c"]`, "Object"},
		{"nested json", `{"outer":{"inner":"value"}}`, "Object"},
		{"plain text", "Hello, World!", "Object"},
		{"empty string", "", "Object"},
		{"numbers", "123456789", "Object"},
		{"special chars", "!@#$%^&*()", "Object"},
		{"unicode", "Hello 世界 🌍", "Object"},
		{"newlines", "line1\nline2\nline3", "Object"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Test - Serialize
			serialized, err := sut.Serialize(&tc.data, tc.typeName)
			if !ExpectNoError(err, t) { return }
			if !ExpectNonNil(serialized, t) { return }

			// Test - Deserialize
			deserialized, err := sut.Deserialize(serialized)
			if !ExpectNoError(err, t) { return }
			if !ExpectNonNil(deserialized, t) { return }

			// Verify - Data should match
			ExpectString(tc.data, *deserialized, t)
		})
	}
}

func TestThat_Serializer_RoundTrip_WithMultipleInstances(t *testing.T) {
	// Setup
	transcoder1 := xc.NewTranscoder()
	transcoder1.SetEncoderScheme(enc.NewEncodingSchemeBase64())
	transcoder1.SetDecoderScheme(enc.NewEncodingSchemeBase64())
	serializer1 := NewSerializer(transcoder1)

	transcoder2 := xc.NewTranscoder()
	transcoder2.SetEncoderScheme(enc.NewEncodingSchemeBase64())
	transcoder2.SetDecoderScheme(enc.NewEncodingSchemeBase64())
	serializer2 := NewSerializer(transcoder2)

	data := `{"message":"cross-instance test"}`
	typeName := "Object"

	// Test - Serialize with first instance
	serialized, err := serializer1.Serialize(&data, typeName)
	if !ExpectNoError(err, t) { return }
	if !ExpectNonNil(serialized, t) { return }

	// Test - Deserialize with second instance
	deserialized, err := serializer2.Deserialize(serialized)
	if !ExpectNoError(err, t) { return }
	if !ExpectNonNil(deserialized, t) { return }

	// Verify - Data should match
	ExpectString(data, *deserialized, t)
}

// -------------------------------------------------------------------------------------------------
// Format Validation Tests
// -------------------------------------------------------------------------------------------------

func TestThat_Serializer_Serialize_ProducesExpectedFormat(t *testing.T) {
	// Setup
	transcoder := xc.NewTranscoder()
	transcoder.SetEncoderScheme(enc.NewEncodingSchemeBase64())
	transcoder.SetDecoderScheme(enc.NewEncodingSchemeBase64())
	sut := NewSerializer(transcoder)
	data := `{"key":"value"}`
	typeName := "Object"

	// Test
	result, err := sut.Serialize(&data, typeName)

	// Verify
	if !ExpectNoError(err, t) { return }
	if !ExpectNonNil(result, t) { return }

	// Check format: "ser[{Method}:{Type}:{Data}]"
	if len(*result) < 10 {
		t.Error("Serialized string too short")
		return
	}
	if (*result)[:4] != "ser[" {
		t.Errorf("Expected prefix 'ser[', got '%s'", (*result)[:4])
	}
	if (*result)[len(*result)-1:] != "]" {
		t.Errorf("Expected suffix ']', got '%s'", (*result)[len(*result)-1:])
	}
	// Should contain "j64:" for method
	if len(*result) < 8 || (*result)[4:8] != "j64:" {
		t.Error("Expected method 'j64:' after prefix")
	}
}
