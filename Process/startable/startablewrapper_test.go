package startable

import (
	"fmt"
	"testing"
)

// Factory Function Tests

func TestThat_MakeStartable_ReturnsInstance(t *testing.T) {
	// Test
	sut := MakeStartable(func() error { return nil }, func() {})

	// Verify
	if sut == nil {
		t.Error("Expected non-nil instance")
	}
}

func TestThat_MakeStartable_AcceptsNilFuncs(t *testing.T) {
	// Test
	sut := MakeStartable(nil, nil)

	// Verify
	if sut == nil {
		t.Error("Expected non-nil instance even with nil funcs")
	}
}

// Start Tests

func TestThat_StartableWrapper_Start_CallsStartFunc(t *testing.T) {
	// Setup
	called := false
	sut := MakeStartable(
		func() error { called = true; return nil },
		func() {},
	)

	// Test
	err := sut.Start()

	// Verify
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !called {
		t.Error("Expected start func to be called")
	}
}

func TestThat_StartableWrapper_Start_ReturnsErrorFromStartFunc(t *testing.T) {
	// Setup
	sut := MakeStartable(
		func() error { return fmt.Errorf("start failed") },
		func() {},
	)

	// Test
	err := sut.Start()

	// Verify
	if err == nil {
		t.Error("Expected error from start func")
	}
	if err.Error() != "start failed" {
		t.Errorf("Expected 'start failed', got '%s'", err.Error())
	}
}

func TestThat_StartableWrapper_Start_ReturnsErrorForNilStartFunc(t *testing.T) {
	// Setup
	sut := MakeStartable(nil, func() {})

	// Test
	err := sut.Start()

	// Verify
	if err == nil {
		t.Error("Expected error for nil start func")
	}
}

func TestThat_StartableWrapper_Start_IsIdempotent(t *testing.T) {
	// Setup
	callCount := 0
	sut := MakeStartable(
		func() error { callCount++; return nil },
		func() {},
	)

	// Test - call Start twice
	sut.Start()
	sut.Start()

	// Verify - start func should only be called once
	if callCount != 1 {
		t.Errorf("Expected start func to be called once, got %d", callCount)
	}
}

// IsStarted Tests

func TestThat_StartableWrapper_IsStarted_ReturnsFalseInitially(t *testing.T) {
	// Setup
	sut := MakeStartable(func() error { return nil }, func() {})

	// Verify
	if sut.IsStarted() {
		t.Error("Expected IsStarted to be false initially")
	}
}

func TestThat_StartableWrapper_IsStarted_ReturnsTrueAfterStart(t *testing.T) {
	// Setup
	sut := MakeStartable(func() error { return nil }, func() {})

	// Test
	sut.Start()

	// Verify
	if !sut.IsStarted() {
		t.Error("Expected IsStarted to be true after Start")
	}
}

func TestThat_StartableWrapper_IsStarted_ReturnsFalseAfterFailedStart(t *testing.T) {
	// Setup
	sut := MakeStartable(func() error { return fmt.Errorf("fail") }, func() {})

	// Test
	sut.Start()

	// Verify
	if sut.IsStarted() {
		t.Error("Expected IsStarted to be false after failed Start")
	}
}

// Stop Tests

func TestThat_StartableWrapper_Stop_CallsStopFunc(t *testing.T) {
	// Setup
	called := false
	sut := MakeStartable(
		func() error { return nil },
		func() { called = true },
	)

	// Test
	sut.Stop()

	// Verify
	if !called {
		t.Error("Expected stop func to be called")
	}
}

// AddStartables Tests

func TestThat_StartableWrapper_AddStartables_ReturnsNil(t *testing.T) {
	// Setup
	sut := MakeStartable(func() error { return nil }, func() {})

	// Test - wrapper doesn't support nested startables
	result := sut.AddStartables(NewStartable())

	// Verify
	if result != nil {
		t.Error("Expected nil from AddStartables on wrapper")
	}
}

// Interface Implementation

func TestThat_StartableWrapper_ImplementsStartableIfc(t *testing.T) {
	var _ StartableIfc = MakeStartable(nil, nil)
}
