package config

import(
	"testing"

	. "github.com/DigiStratum/GoLib/Testing"
)

func TestThat_Config_NewConfig_ReturnsSomething(t *testing.T) {
	// Setup
	sut := NewConfig()

	// Verify
	ExpectNonNil(sut, t)
}

func TestThat_Config_MergeConfig_ChangesNothing_WhenGivenEmptyConfig(t *testing.T) {
	// Setup
	sut := NewConfig()
	sut.Set("boguskey", "bogusvalue")
	emptyCfg := NewConfig()

	// Test
	sut.MergeConfig(emptyCfg)

	// Verify
	ExpectInt(1, sut.Size(), t)
}
