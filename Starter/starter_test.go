package starter

import (
	"fmt"
	"testing"
)

// NOTE: This package is DEPRECATED - use Process/startable instead

// Factory Function Tests

func TestThat_NewStartable_ReturnsInstance(t *testing.T) {
	sut := NewStartable()
	if sut == nil {
		t.Error("Expected non-nil instance")
	}
}

func TestThat_NewStartable_AcceptsInitialStartables(t *testing.T) {
	child := NewStartable()
	sut := NewStartable(child)
	if sut == nil {
		t.Error("Expected non-nil instance")
	}
	if len(sut.startables) != 1 {
		t.Errorf("Expected 1 startable, got %d", len(sut.startables))
	}
}

// AddStartables Tests

func TestThat_Startable_AddStartables_AddsMultiple(t *testing.T) {
	sut := NewStartable()
	result := sut.AddStartables(NewStartable(), NewStartable())
	if result != sut {
		t.Error("Expected fluent return")
	}
	if len(sut.startables) != 2 {
		t.Errorf("Expected 2 startables, got %d", len(sut.startables))
	}
}

// Start/IsStarted Tests

func TestThat_Startable_IsStarted_ReturnsFalseInitially(t *testing.T) {
	sut := NewStartable()
	if sut.IsStarted() {
		t.Error("Expected false initially")
	}
}

func TestThat_Startable_Start_SetsIsStarted(t *testing.T) {
	sut := NewStartable()
	err := sut.Start()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !sut.IsStarted() {
		t.Error("Expected true after Start")
	}
}

func TestThat_Startable_Start_IsIdempotent(t *testing.T) {
	sut := NewStartable()
	sut.Start()
	err := sut.Start()
	if err != nil {
		t.Errorf("Unexpected error on second Start: %v", err)
	}
}

func TestThat_Startable_Start_StartsChildren(t *testing.T) {
	child := NewStartable()
	sut := NewStartable(child)
	sut.Start()
	if !child.IsStarted() {
		t.Error("Expected child to be started")
	}
}

func TestThat_Startable_Start_PropagatesChildError(t *testing.T) {
	child := MakeStartable(func() error { return fmt.Errorf("fail") })
	sut := NewStartable(child)
	err := sut.Start()
	if err == nil {
		t.Error("Expected error from child")
	}
}

// MakeStartable (wrapper) Tests

func TestThat_MakeStartable_ReturnsInstance(t *testing.T) {
	sut := MakeStartable(func() error { return nil })
	if sut == nil {
		t.Error("Expected non-nil instance")
	}
}

func TestThat_StartableWrapper_Start_CallsStartFunc(t *testing.T) {
	called := false
	sut := MakeStartable(func() error { called = true; return nil })
	sut.Start()
	if !called {
		t.Error("Expected start func to be called")
	}
}

func TestThat_StartableWrapper_Start_ReturnsError(t *testing.T) {
	sut := MakeStartable(func() error { return fmt.Errorf("fail") })
	err := sut.Start()
	if err == nil {
		t.Error("Expected error")
	}
}

func TestThat_StartableWrapper_Start_ReturnsErrorForNilFunc(t *testing.T) {
	sut := MakeStartable(nil)
	err := sut.Start()
	if err == nil {
		t.Error("Expected error for nil func")
	}
}

func TestThat_StartableWrapper_Start_IsIdempotent(t *testing.T) {
	count := 0
	sut := MakeStartable(func() error { count++; return nil })
	sut.Start()
	sut.Start()
	if count != 1 {
		t.Errorf("Expected 1 call, got %d", count)
	}
}

func TestThat_StartableWrapper_AddStartables_ReturnsNil(t *testing.T) {
	sut := MakeStartable(func() error { return nil })
	result := sut.AddStartables(NewStartable())
	if result != nil {
		t.Error("Expected nil from AddStartables on wrapper")
	}
}

// Interface tests

func TestThat_Startable_ImplementsStartableIfc(t *testing.T) {
	var _ StartableIfc = NewStartable()
}

func TestThat_StartableWrapper_ImplementsStartableIfc(t *testing.T) {
	var _ StartableIfc = MakeStartable(nil)
}
