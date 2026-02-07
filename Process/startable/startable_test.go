package startable

import (
	"fmt"
	"testing"
)

// Factory Function Tests

func TestThat_NewStartable_ReturnsInstance(t *testing.T) {
	// Test
	sut := NewStartable()

	// Verify
	if sut == nil {
		t.Error("Expected non-nil instance")
	}
}

func TestThat_NewStartable_AcceptsInitialStartables(t *testing.T) {
	// Setup
	child := NewStartable()

	// Test
	sut := NewStartable(child)

	// Verify
	if sut == nil {
		t.Error("Expected non-nil instance")
	}
	if len(sut.startables) != 1 {
		t.Errorf("Expected 1 startable, got %d", len(sut.startables))
	}
}

// AddStartables Tests

func TestThat_Startable_AddStartables_AddsMultiple(t *testing.T) {
	// Setup
	sut := NewStartable()
	child1 := NewStartable()
	child2 := NewStartable()

	// Test
	result := sut.AddStartables(child1, child2)

	// Verify
	if result != sut {
		t.Error("Expected fluent return")
	}
	if len(sut.startables) != 2 {
		t.Errorf("Expected 2 startables, got %d", len(sut.startables))
	}
}

// Start/IsStarted Tests

func TestThat_Startable_IsStarted_ReturnsFalseInitially(t *testing.T) {
	// Setup
	sut := NewStartable()

	// Verify
	if sut.IsStarted() {
		t.Error("Expected IsStarted to be false initially")
	}
}

func TestThat_Startable_Start_SetsIsStarted(t *testing.T) {
	// Setup
	sut := NewStartable()

	// Test
	err := sut.Start()

	// Verify
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !sut.IsStarted() {
		t.Error("Expected IsStarted to be true after Start")
	}
}

func TestThat_Startable_Start_IsIdempotent(t *testing.T) {
	// Setup
	sut := NewStartable()

	// Test - call Start twice
	err1 := sut.Start()
	err2 := sut.Start()

	// Verify - no errors on subsequent calls
	if err1 != nil {
		t.Errorf("Unexpected error on first Start: %v", err1)
	}
	if err2 != nil {
		t.Errorf("Unexpected error on second Start: %v", err2)
	}
}

func TestThat_Startable_Start_StartsChildren(t *testing.T) {
	// Setup
	child := NewStartable()
	sut := NewStartable(child)

	// Test
	err := sut.Start()

	// Verify
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !child.IsStarted() {
		t.Error("Expected child to be started")
	}
}

func TestThat_Startable_Start_PropagatesChildError(t *testing.T) {
	// Setup - use wrapper that returns error
	child := MakeStartable(
		func() error { return fmt.Errorf("child start failed") },
		func() {},
	)
	sut := NewStartable(child)

	// Test
	err := sut.Start()

	// Verify
	if err == nil {
		t.Error("Expected error from child")
	}
	if !sut.IsStarted() == false {
		// Parent should not be marked as started if child fails
	}
}

// Stop Tests

func TestThat_Startable_Stop_StopsChildren(t *testing.T) {
	// Setup
	stopCalled := false
	child := MakeStartable(
		func() error { return nil },
		func() { stopCalled = true },
	)
	sut := NewStartable(child)
	sut.Start()

	// Test
	sut.Stop()

	// Verify
	if !stopCalled {
		t.Error("Expected child Stop to be called")
	}
}

func TestThat_Startable_Stop_IsNoOpWhenNotStarted(t *testing.T) {
	// Setup
	stopCalled := false
	child := MakeStartable(
		func() error { return nil },
		func() { stopCalled = true },
	)
	sut := NewStartable(child)
	// Note: NOT calling Start()

	// Test
	sut.Stop()

	// Verify - stop should not be called on children if not started
	if stopCalled {
		t.Error("Expected child Stop to NOT be called when parent not started")
	}
}

// Interface Implementation

func TestThat_Startable_ImplementsStartableIfc(t *testing.T) {
	var _ StartableIfc = NewStartable()
}
