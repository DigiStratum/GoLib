package events

import (
	"encoding/json"
	"testing"
)

// Factory Function Tests

func TestThat_NewEvent_ReturnsInstance(t *testing.T) {
	// Setup
	props := map[string]string{"key": "value"}

	// Test
	sut := NewEvent(props)

	// Verify
	if sut == nil {
		t.Error("Expected non-nil instance")
	}
}

func TestThat_NewEvent_ReturnsInstanceWithNilProps(t *testing.T) {
	// Test
	sut := NewEvent(nil)

	// Verify
	if sut == nil {
		t.Error("Expected non-nil instance even with nil properties")
	}
}

func TestThat_NewEvent_ReturnsInstanceWithEmptyProps(t *testing.T) {
	// Test
	sut := NewEvent(map[string]string{})

	// Verify
	if sut == nil {
		t.Error("Expected non-nil instance with empty properties")
	}
}

// ToJson Tests

func TestThat_Event_ToJson_ReturnsValidJson(t *testing.T) {
	// Setup
	props := map[string]string{
		"event_type": "test",
		"data":       "hello",
	}
	sut := NewEvent(props)

	// Test
	jsonStr, err := sut.ToJson()

	// Verify
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}
	if jsonStr == nil {
		t.Error("Expected non-nil JSON string")
		return
	}

	// Verify it's valid JSON by unmarshaling
	var result map[string]string
	if err := json.Unmarshal([]byte(*jsonStr), &result); err != nil {
		t.Errorf("Invalid JSON produced: %v", err)
		return
	}

	// Verify content
	if result["event_type"] != "test" {
		t.Errorf("Expected event_type 'test', got '%s'", result["event_type"])
	}
	if result["data"] != "hello" {
		t.Errorf("Expected data 'hello', got '%s'", result["data"])
	}
}

func TestThat_Event_ToJson_ReturnsEmptyObjectForEmptyProps(t *testing.T) {
	// Setup
	sut := NewEvent(map[string]string{})

	// Test
	jsonStr, err := sut.ToJson()

	// Verify
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}
	if jsonStr == nil {
		t.Error("Expected non-nil JSON string")
		return
	}
	if *jsonStr != "{}" {
		t.Errorf("Expected '{}', got '%s'", *jsonStr)
	}
}

func TestThat_Event_ToJson_ReturnsNullForNilProps(t *testing.T) {
	// Setup
	sut := NewEvent(nil)

	// Test
	jsonStr, err := sut.ToJson()

	// Verify
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}
	if jsonStr == nil {
		t.Error("Expected non-nil JSON string")
		return
	}
	if *jsonStr != "null" {
		t.Errorf("Expected 'null', got '%s'", *jsonStr)
	}
}

func TestThat_Event_ToJson_HandlesSpecialCharacters(t *testing.T) {
	// Setup
	props := map[string]string{
		"message": "hello \"world\" with\nnewlines",
	}
	sut := NewEvent(props)

	// Test
	jsonStr, err := sut.ToJson()

	// Verify
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}

	// Verify valid JSON with proper escaping
	var result map[string]string
	if err := json.Unmarshal([]byte(*jsonStr), &result); err != nil {
		t.Errorf("Invalid JSON produced: %v", err)
	}
}

// Interface Implementation Tests

func TestThat_Event_ImplementsEventIfc(t *testing.T) {
	// This test verifies the Event struct implements EventIfc
	var _ EventIfc = &Event{}
	var _ EventIfc = NewEvent(nil)
}
