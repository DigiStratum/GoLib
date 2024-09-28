package config

/*
A support structure for Configuration Item declaration/handling

We can supply ConfigItem's to NewConfigurable so that Configure() can be passed the custom handlers needed.

TODO:
 * Need test coverage for this mess
 * Add support for type? All incoming config values are arbitrary strings, no type enforcement
 * Are there any useful type-conversion helpers for capture funcs?
*/

import (
	"fmt"

	"GoLib/Data"
)

type CaptureFunc func (dataValue data.DataValueIfc) error
type ValidateFunc func (dataValue data.DataValueIfc) bool

type ConfigItemIfc interface {
	GetName() string
	SetRequired() *configItem
	IsRequired() bool

	CanCapture() bool
	CaptureWith(captureFunc CaptureFunc) *configItem
	Capture(dataValue data.DataValueIfc) error

	CanValidate() bool
	ValidateWith(validateFunc ValidateFunc) *configItem
	Validate(dataValue data.DataValueIfc) bool

	SetDefault(dataValue data.DataValueIfc) *configItem
}

type configItem struct {
	name		string
	defaultValue	data.DataValueIfc
	hasDefault	bool
	isRequired	bool
	captureFunc	CaptureFunc
	validateFunc	ValidateFunc
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

func (r *configItem) SetDefault(dataValue data.DataValueIfc) *configItem {
	r.defaultValue = dataValue
	r.hasDefault = true
	return r
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

func (r *configItem) Capture(dataValue data.DataValueIfc) error {
	if ! r.CanCapture() {
		return fmt.Errorf(
			"No Capture function is set for configItem: %s",
			r.GetName(),
		)
	}

	// If the value is nil...
	if nil == dataValue {
		// ... and we have a default...
		if r.hasDefault {
			// ... use it!
			return r.captureFunc(r.defaultValue)
		}

		// ... no default, yet a value is required...
		if r.isRequired {
			// ... reject it!
			return fmt.Errorf(
				"Nil value with no Default for required configItem: %s",
				r.GetName(),
			)
		}

		// ... not required, so ignore it.
		return nil
	}

	return r.captureFunc(dataValue)
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

func (r *configItem) Validate(dataValue data.DataValueIfc) bool {
	// You're calling Validate on a ConfigItem that can't be validated? Fail!
	if ! r.CanValidate() { return false }
	return r.validateFunc(dataValue)
}

