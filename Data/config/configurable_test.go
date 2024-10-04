package config

import(
	"fmt"
	"testing"

	"GoLib/Data"

	. "GoLib/Testing"
)

// Interface

func TestThat_Configurable_NewConfigurable_Returns_Instance(t *testing.T) {
	// Setup
	var sut ConfigurableIfc = NewConfigurable() // Verifies that result satisfies IFC

	// Verify
	if ! ExpectNonNil(sut, t) { return }
}

// DeclareConfigItems

func TestThat_Configurable_DeclareConfigItems_returns_self(t *testing.T) {
	// Setup
	sut := NewConfigurable()

	// Test
	actual := sut.DeclareConfigItems()

	// Verify
	if ! ExpectNonNil(actual, t) { return }
	if ! ExpectInt(0, len(actual.declared), t) { return }
}

func TestThat_Configurable_DeclareConfigItems_declares_items_by_selector(t *testing.T) {
	// Setup
	sut := NewConfigurable()

	// Test
	actual := sut.DeclareConfigItems(
		NewConfigItem("selector1"),
		NewConfigItem("selector2"),
		NewConfigItem("selector1"), // <- Note: this duplicate should replace original
	)

	// Verify
	if ! ExpectInt(2, len(actual.declared), t) { return }
}

// Configure

func TestThat_Configurable_Configure_Returns_error_when_started(t *testing.T) {
	// Setup
	sut := NewConfigurable()

	// Test
	actualErr1 := sut.Configure(nil)
	sut.Start()
	actualErr2 := sut.Configure(nil)

	// Verify
	if ! ExpectNoError(actualErr1, t) { return }
	if ! ExpectError(actualErr2, t) { return }
}

func TestThat_Configurable_Configure_Replaces_Configurable_data(t *testing.T) {
	// Setup
	sut := NewConfigurable()
	initialConfig := NewConfig()
	initialConfig.PrepareArray().AppendArrayValue(data.NewString("initial"))
	sut.Configure(initialConfig)
	expectedKey := "new"
	expectedValue := "replacement"
	newConfig := NewConfig()
	newConfig.PrepareObject().SetObjectProperty(expectedKey, data.NewString(expectedValue))

	// Test
	sut.Configure(newConfig)

	// Verify
	if ! ExpectTrue(sut.config.IsObject(), t) { return }
	if ! ExpectTrue(sut.config.HasObjectProperty(expectedKey), t) { return }
	if ! ExpectString(expectedValue, sut.config.GetObjectProperty(expectedKey).GetString(), t) { return }
}

// GetMissingConfigs | HasMissingConfigs

func TestThat_Configurable_No_MissingConfigs_when_nothing_declared(t *testing.T) {
	// Setup
	sut := NewConfigurable()

	// Verify
	if ! ExpectInt(0, len(sut.GetMissingConfigs()), t) { return }
	if ! ExpectFalse(sut.HasMissingConfigs(), t) { return }
}

func TestThat_Configurable_Has_MissingConfigs_when_required_but_missing(t *testing.T) {
	// Setup
	sut := NewConfigurable().
		DeclareConfigItems(
			NewConfigItem("c1"),
			NewConfigItem("c2").SetRequired(),
			NewConfigItem("c3").SetRequired(),
			NewConfigItem("c4").SetRequired(),
		)
	config := NewConfig()
	config.PrepareObject().
		SetObjectProperty("c3", data.NewNull())
	sut.Configure(config)

	// Verify
	if ! ExpectInt(2, len(sut.GetMissingConfigs()), t) { return }
	if ! ExpectTrue(sut.HasMissingConfigs(), t) { return }
}

func TestThat_Configurable_No_MissingConfigs_when_required_and_provided(t *testing.T) {
	// Setup
	sut := NewConfigurable().
		DeclareConfigItems(
			NewConfigItem("c1"),
			NewConfigItem("c2").SetRequired(),
			NewConfigItem("c3").SetRequired(),
			NewConfigItem("c4").SetRequired(),
		)
	config := NewConfig()
	config.PrepareObject().
		SetObjectProperty("c2", data.NewNull()).
		SetObjectProperty("c3", data.NewNull()).
		SetObjectProperty("c4", data.NewNull())
	sut.Configure(config)

	// Verify
	if ! ExpectInt(0, len(sut.GetMissingConfigs()), t) { return }
	if ! ExpectFalse(sut.HasMissingConfigs(), t) { return }
}

// Start

func TestThat_Configurable_Start_returns_no_error_by_default(t *testing.T) {
	// Setup
	sut := NewConfigurable()

	// Verify
	if ! ExpectNoError(sut.Start(), t) { return }
	if ! ExpectNoError(sut.Start(), t) { return } // Also no error when starting a second time
}

func TestThat_Configurable_Start_returns_error_with_missing_configs(t *testing.T) {
	// Setup
	sut := NewConfigurable().
		DeclareConfigItems(
			NewConfigItem("c1"),
			NewConfigItem("c2").SetRequired(),
			NewConfigItem("c3").SetRequired(),
			NewConfigItem("c4").SetRequired(),
		)
	config := NewConfig()
	config.PrepareObject().
		SetObjectProperty("c3", data.NewNull())
	sut.Configure(config)

	// Verify
	if ! ExpectError(sut.Start(), t) { return }
}

func TestThat_Configurable_Start_returns_error_for_validation_errors(t *testing.T) {
	// Setup
	sut := NewConfigurable().
		DeclareConfigItems(
			NewConfigItem("c1").ValidateWith(func (dv data.DataValueIfc) error {
				return fmt.Errorf("fail!")
			}),
		)
	config := NewConfig()
	config.PrepareObject().
		SetObjectProperty("c1", data.NewNull())
	sut.Configure(config)

	// Verify
	if ! ExpectError(sut.Start(), t) { return }
}

func TestThat_Configurable_Start_returns_error_for_capture_errors(t *testing.T) {
	// Setup
	sut := NewConfigurable().
		DeclareConfigItems(
			NewConfigItem("c1").CaptureWith(func (dv data.DataValueIfc) error {
				return fmt.Errorf("fail!")
			}),
		)
	config := NewConfig()
	config.PrepareObject().
		SetObjectProperty("c1", data.NewNull())
	sut.Configure(config)

	// Verify
	if ! ExpectError(sut.Start(), t) { return }
}



// GetConfig

func TestThat_Configurable_Returns_nil_when_not_started(t *testing.T) {
	// Setup
	sut := NewConfigurable()

	// Verify
	if ! ExpectNil(sut.GetConfig(), t) { return }
}

func TestThat_Configurable_Returns_config_when_started(t *testing.T) {
	// Setup
	sut := NewConfigurable()
	expectedKey := "prop"
	expectedValue := "value"
	cfg := NewConfig()
	cfg.PrepareObject().
		SetObjectProperty(expectedKey, data.NewString(expectedValue))
	sut.Configure(cfg)
	sut.Start()

	// Test
	actual := sut.GetConfig()

	// Verify
	if ! ExpectNonNil(actual, t) { return }
	if ! ExpectTrue(actual.IsObject(), t) { return }
	if ! ExpectNonNil(actual.GetObjectProperty(expectedKey), t) { return }
	if ! ExpectTrue(actual.GetObjectProperty(expectedKey).IsString(), t) { return }
	if ! ExpectString(expectedValue, actual.GetObjectProperty(expectedKey).GetString(), t) { return }
}


