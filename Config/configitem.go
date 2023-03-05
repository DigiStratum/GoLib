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

	CanValidate() bool
	ValidateWith(validateFunc func (value string) bool) *configItem
	Validate(value string) bool

	// TODO: Are there any useful type-conversion helpers for capture funcs?
}

// TODO: Add support for type? Everything is a string coming in...
type configItem struct {
	name		string
	isRequired	bool
	captureFunc	func (value string) error
	validateFunc	func (value string) bool
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

// Capture
// -----------------------------------------------

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

// Validation
// -----------------------------------------------

func (r *configItem) CanValidate() bool {
	return nil != r.validateFunc
}

func (r *configItem) ValidateWith(validateFunc func (value string) bool) *configItem {
	r.validateFunc = validateFunc
	return r
}

func (r *configItem) Validate(value string) bool {
	// You're calling Validate on a ConfigItem that can't be validated? Fail.
	if ! r.CanValidate() { return false }
	return r.validateFunc(value)
}

