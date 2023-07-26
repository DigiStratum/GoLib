package config

/*
A support structure for Configuration Item declaration/handling

We can supply ConfigItem's to NewConfigurable so that Configure() can be passed the custom handlers needed.

TODO:
 * Need test coverage for this mess
 * Add support for type? All incoming config values are arbitrary strings, no type enforcement
*/

import (
	"fmt"
)

type CaptureFunc func (value string) error
type SubsetFunc func (config ConfigIfc) error
type ValidateFunc func (value string) bool

type ConfigItemIfc interface {
	GetName() string
	SetRequired() *configItem
	IsRequired() bool

	CanCapture() bool
	CaptureWith(captureFunc CaptureFunc) *configItem
	Capture(value string) error

	CanCaptureSubset() bool
	CaptureSubsetWith(subsetFunc SubsetFunc) *configItem
	CaptureSubset(config ConfigIfc) error

	CanValidate() bool
	ValidateWith(validateFunc ValidateFunc) *configItem
	Validate(value string) bool

	// TODO: Are there any useful type-conversion helpers for capture funcs?
}

type configItem struct {
	name		string
	isRequired	bool
	captureFunc	CaptureFunc
	validateFunc	ValidateFunc
	subsetFunc	SubsetFunc
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

func (r *configItem) CaptureWith(captureFunc CaptureFunc) *configItem {
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

// Subsets
// -----------------------------------------------

func (r *configItem) CanCaptureSubset() bool {
	return nil != r.subsetFunc
}

func (r *configItem) CaptureSubsetWith(subsetFunc SubsetFunc) *configItem {
	r.subsetFunc = subsetFunc
	return r
}

func (r *configItem) CaptureSubset(config ConfigIfc) error {
	if ! r.CanCaptureSubset() {
		return fmt.Errorf(
			"No Subset Capture function is set for configItem: %s",
			r.GetName(),
		)
	}

	return r.subsetFunc(config)
}

// Validation
// -----------------------------------------------

func (r *configItem) CanValidate() bool {
	return nil != r.validateFunc
}

func (r *configItem) ValidateWith(validateFunc ValidateFunc) *configItem {
	r.validateFunc = validateFunc
	return r
}

func (r *configItem) Validate(value string) bool {
	// You're calling Validate on a ConfigItem that can't be validated? Fail.
	if ! r.CanValidate() { return false }
	return r.validateFunc(value)
}

