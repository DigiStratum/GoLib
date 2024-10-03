package config

import(
	"testing"

	"GoLib/Data"

	. "GoLib/Testing"
)

// Interface

func TestThat_Config_NewConfigurable_Returns_Instance(t *testing.T) {
	// Setup
	var sut ConfigurableIfc = NewConfigurable() // Verifies that result satisfies IFC

	// Verify
	if ! ExpectNonNil(sut, t) { return }
}

// DeclareConfigItems

func TestThat_Config_DeclareConfigItems_returns_self(t *testing.T) {
	// Setup
	sut := NewConfigurable()

	// Test
	actual := sut.DeclareConfigItems()

	// Verify
	if ! ExpectNonNil(actual, t) { return }
	if ! ExpectInt(0, len(actual.declared), t) { return }
}

func TestThat_Config_DeclareConfigItems_declares_items_by_selector(t *testing.T) {
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

func TestThat_Config_Configure_Returns_error_when_started(t *testing.T) {
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

func TestThat_Config_Configure_Replaces_Config_data(t *testing.T) {
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

func TestThat_Config_No_MissingConfigs_when_nothing_declared(t *testing.T) {
	// Setup
	sut := NewConfigurable()

	// Verify
	if ! ExpectInt(0, len(sut.GetMissingConfigs()), t) { return }
	if ! ExpectFalse(sut.HasMissingConfigs(), t) { return }
}

func TestThat_Config_Has_MissingConfigs_when_required_but_missing(t *testing.T) {
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

func TestThat_Config_No_MissingConfigs_when_required_and_provided(t *testing.T) {
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



