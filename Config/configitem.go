package config

import (
	"fmt"
)

type ConfigItemIfc interface {
	GetName() string
	SetRequired() *configItem
	IsRequired() bool

	CanCapture() bool
	CaptureWith(captureFunc func (value string) error) *configItem
	Capture(value string) error

	// TODO: Are there any useful type-conversion helpers for capture funcs?
}

// TODO: Add support for type? Everything is a string coming in...
type configItem struct {
	name		string
	isRequired	bool
	captureFunc	func (value string) error
}

// -------------------------------------------------------------------------------------------------
// Factory functions
// -------------------------------------------------------------------------------------------------

func NewConfigItem(name string) *configItem {
	return &configItem{
		name:		name,
	}
}

// -------------------------------------------------------------------------------------------------
// ConfigItemIfc
// -------------------------------------------------------------------------------------------------

func (r *configItem) GetName() string {
	return r.name
}

func (r *configItem) SetRequired() *configItem {
	r.isRequired = true
	return r
}

func (r *configItem) IsRequired() bool {
	return r.isRequired
}

func (r *configItem) CanCapture() bool {
	return nil != r.captureFunc
}

func (r *configItem) CaptureWith(captureFunc func (value string) error) *configItem {
	r.captureFunc = captureFunc
	return r
}

func (r *configItem) Capture(value string) error {
	if ! r.CanCapture() {
		return fmt.Errorf(
			"No Capture function is set for configItem: %s",
			r.GetName(),
		)
	}

	return r.captureFunc(value)
}

