package config

import(
	//"fmt"
	//"strings"
	"testing"

	. "GoLib/Testing"
)

// Interface

func TestThat_Config_Newonfig_ReturnsInstance(t *testing.T) {
	// Setup
	var sut ConfigIfc = NewConfig() // Verifies that result satisfies IFC

	// Verify
	if ! ExpectNonNil(sut, t) { return }
}


