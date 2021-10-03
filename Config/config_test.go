package config

import(
	"fmt"
	"testing"

	. "github.com/DigiStratum/GoLib/Testing"
)

func TestThat_Config_NewConfig_ReturnsSomething(t *testing.T) {
	// Setup
	sut := NewConfig()

	// Verify
	ExpectNonNil(sut, t)
}

func TestThat_Config_MergeConfig_ChangesNothing_WhenProvidedEmptyConfig(t *testing.T) {
	// Setup
	sut := NewConfig()
	sut.Set("boguskey", "bogusvalue")
	mergeCfg := NewConfig()

	// Test
	sut.MergeConfig(mergeCfg)

	// Verify
	ExpectInt(1, sut.Size(), t)
}

func TestThat_Config_MergeConfig_AddsSome_WhenProvidedPopulatedConfig(t *testing.T) {
	// Setup
	sut := NewConfig()
	sut.Set("boguskey", "bogusvalue")
	mergeCfg := NewConfig()
	numKeys := 10
	for k := 0; k < numKeys; k++ {
		mergeCfg.Set(fmt.Sprintf("key%d", k), fmt.Sprintf("value%d", k))
	}

	// Test
	sut.MergeConfig(mergeCfg)

	// Verify
	ExpectInt(numKeys + 1, sut.Size(), t)
}

func TestThat_Config_GetSubset_ReturnsPopulatedConfig_WhenProvidedGoodPrefix(t *testing.T) {
	// Setup
	sut := NewConfig()
	sut.Set("boguskey", "bogusvalue")
	numKeys := 10
	for k := 0; k < numKeys; k++ {
		sut.Set(fmt.Sprintf("key%d", k), fmt.Sprintf("value%d", k))
	}

	// Test
	subCfg := sut.GetSubset("k")

	// Verify
	ExpectInt(numKeys, subCfg.Size(), t)
}

func TestThat_Config_GetSubset_ReturnsEmptyConfig_WhenProvidedBadPrefix(t *testing.T) {
	// Setup
	sut := NewConfig()
	sut.Set("boguskey", "bogusvalue")
	numKeys := 10
	for k := 0; k < numKeys; k++ {
		sut.Set(fmt.Sprintf("key%d", k), fmt.Sprintf("value%d", k))
	}

	// Test
	subCfg := sut.GetSubset("invalidprefix")

	// Verify
	ExpectInt(0, subCfg.Size(), t)
}

func TestThat_Config_GetInverseSubset_ReturnsPopulatedConfig_WhenProvidedGoodPrefix(t *testing.T) {
	// Setup
	sut := NewConfig()
	sut.Set("boguskey", "bogusvalue")
	numKeys := 10
	for k := 0; k < numKeys; k++ {
		sut.Set(fmt.Sprintf("key%d", k), fmt.Sprintf("value%d", k))
	}

	// Test
	subCfg := sut.GetInverseSubset("k")

	// Verify
	ExpectInt(1, subCfg.Size(), t)
}

func TestThat_Config_GetInverseSubset_ReturnsEmptyConfig_WhenProvidedBadPrefix(t *testing.T) {
	// Setup
	sut := NewConfig()
	sut.Set("boguskey", "bogusvalue")
	numKeys := 10
	for k := 0; k < numKeys; k++ {
		sut.Set(fmt.Sprintf("key%d", k), fmt.Sprintf("value%d", k))
	}

	// Test
	subCfg := sut.GetInverseSubset("invalidprefix")

	// Verify
	ExpectInt(numKeys+1, subCfg.Size(), t)
}

func TestThat_Config_DereferenceString_ChangesNothing_WhenProvidedSimpleString(t *testing.T) {
	// Setup
	sut := NewConfig()
	sut.Set("boguskey", "bogusvalue")

	// Test
	src := "plainstring"
	res := sut.DereferenceString(src)

	// Verify
	ExpectString(src, *res, t)
}

func TestThat_Config_DereferenceString_SubstitutesValues_WhenProvidedStringWithKeys(t *testing.T) {
	// Setup
	sut := NewConfig()
	sut.Set("k1", "v1")
	sut.Set("k2", "v2")

	// Test
	res := sut.DereferenceString("--%k1%--%k2%--")

	// Verify
	ExpectString("--v1--v2--", *res, t)
}

func TestThat_Config_Dereference_SubstitutesValuesInOurConfig_WhenProvidedReferenceConfigWithMatchingKeys(t *testing.T) {
	// Setup
	sut := NewConfig()
	ref := NewConfig()
	numKeys := 3
	for k := 0; k < numKeys; k++ {
		sut.Set(fmt.Sprintf("key%d", k), fmt.Sprintf("--%%ref%d%%--", k))
		ref.Set(fmt.Sprintf("ref%d", k), fmt.Sprintf("value%d", k))
	}

	// Test
	res := sut.Dereference(ref)

	// Verify
	ExpectInt(numKeys, res, t)
	for k := 0; k < numKeys; k++ {
		key := sut.Get(fmt.Sprintf("key%d", k))
		ExpectString(fmt.Sprintf("--value%d--", k), *key, t)
	}
}

func TestThat_Config_Dereference_ChangesNothing_WhenProvidedNothing(t *testing.T) {
	// Setup
	sut := NewConfig()
	var ref ConfigIfc
	numKeys := 3
	for k := 0; k < numKeys; k++ {
		sut.Set(fmt.Sprintf("key%d", k), fmt.Sprintf("--%%ref%d%%--", k))
	}

	// Test
	res := sut.Dereference(ref)

	// Verify
	ExpectInt(0, res, t)
	for k := 0; k < numKeys; k++ {
		key := sut.Get(fmt.Sprintf("key%d", k))
		ExpectString(fmt.Sprintf("--%%ref%d%%--", k), *key, t)
	}
}

func TestThat_Config_DereferenceAll_ChangesNothing_WhenProvidedNothing(t *testing.T) {
	// Setup
	sut := NewConfig()
	var ref ConfigIfc
	numKeys := 3
	for k := 0; k < numKeys; k++ {
		sut.Set(fmt.Sprintf("key%d", k), fmt.Sprintf("--%%ref%d%%--", k))
	}

	// Test
	res := sut.DereferenceAll(ref)

	// Verify
	ExpectInt(0, res, t)
	for k := 0; k < numKeys; k++ {
		key := sut.Get(fmt.Sprintf("key%d", k))
		ExpectString(fmt.Sprintf("--%%ref%d%%--", k), *key, t)
	}
}

func TestThat_Config_DereferenceAll_SubstitutesValuesInOurConfig_WhenProvidedReferenceConfigsWithMatchingKeys(t *testing.T) {
	// Setup
	sut := NewConfig()
	ref1 := NewConfig()
	ref2 := NewConfig()
	numKeys := 3
	for k := 0; k < numKeys; k++ {
		sut.Set(fmt.Sprintf("key1%d", k), fmt.Sprintf("--%%ref1%d%%--", k))
		ref1.Set(fmt.Sprintf("ref1%d", k), fmt.Sprintf("value1%d", k))
		sut.Set(fmt.Sprintf("key2%d", k), fmt.Sprintf("--%%ref2%d%%--", k))
		ref2.Set(fmt.Sprintf("ref2%d", k), fmt.Sprintf("value2%d", k))
	}

	// Test
	res := sut.DereferenceAll(ref1, ref2)

	// Verify
	ExpectInt(numKeys * 2, res, t)
	for k := 0; k < numKeys; k++ {
		key := sut.Get(fmt.Sprintf("key1%d", k))
		ExpectString(fmt.Sprintf("--value1%d--", k), *key, t)
		key = sut.Get(fmt.Sprintf("key2%d", k))
		ExpectString(fmt.Sprintf("--value2%d--", k), *key, t)
	}
}

func TestThat_Config_DereferenceLoop_ChangesNothing_WhenProvidedNothing(t *testing.T) {
	// Setup
	sut := NewConfig()
	var ref ConfigIfc
	numKeys := 3
	for k := 0; k < numKeys; k++ {
		sut.Set(fmt.Sprintf("key%d", k), fmt.Sprintf("--%%ref%d%%--", k))
	}

	// Test
	res := sut.DereferenceLoop(10, ref)

	// Verify
	ExpectFalse(res, t)
	for k := 0; k < numKeys; k++ {
		key := sut.Get(fmt.Sprintf("key%d", k))
		ExpectString(fmt.Sprintf("--%%ref%d%%--", k), *key, t)
	}
}

func TestThat_Config_DereferenceLoop_ChangesNothing_WhenMaxLoopsZero(t *testing.T) {
	// Setup
	sut := NewConfig()
	ref1 := NewConfig()
	numKeys := 3
	for k := 0; k < numKeys; k++ {
		sut.Set(fmt.Sprintf("key%d", k), fmt.Sprintf("-%%ref1%d%%-", k))
		ref1.Set(fmt.Sprintf("ref1%d", k), fmt.Sprintf("-%%ref2%d%%-", k))
		ref1.Set(fmt.Sprintf("ref2%d", k), fmt.Sprintf("-value%d-", k))
	}

	// Test / Verify
	res := sut.DereferenceLoop(0, ref1)
	ExpectFalse(res, t)
	for k := 0; k < numKeys; k++ {
		key := sut.Get(fmt.Sprintf("key%d", k))
		ExpectString(fmt.Sprintf("-%%ref1%d%%-", k), *key, t)
	}
}

func TestThat_Config_DereferenceLoop_SubstitutesAllValuesInOurConfig_WhenProvidedReferenceConfigsWithMatchingKeysNoLimit(t *testing.T) {
	// Setup
	sut := NewConfig()
	ref1 := NewConfig()
	numKeys := 3
	for k := 0; k < numKeys; k++ {
		sut.Set(fmt.Sprintf("key%d", k), fmt.Sprintf("-%%ref1%d%%-", k))
		ref1.Set(fmt.Sprintf("ref1%d", k), fmt.Sprintf("-%%ref2%d%%-", k))
		ref1.Set(fmt.Sprintf("ref2%d", k), fmt.Sprintf("-value%d-", k))
	}

	// Test
	res := sut.DereferenceLoop(10, ref1)

	// Verify
	ExpectTrue(res, t)
	for k := 0; k < numKeys; k++ {
		key := sut.Get(fmt.Sprintf("key%d", k))
		ExpectString(fmt.Sprintf("---value%d---", k), *key, t)
	}
}

func TestThat_Config_DereferenceLoop_SubstitutesSomeValuesInOurConfig_WhenProvidedReferenceConfigsWithMatchingKeysLimit1(t *testing.T) {
	// Setup
	sut := NewConfig()
	ref1 := NewConfig()
	numKeys := 3
	for k := 0; k < numKeys; k++ {
		sut.Set(fmt.Sprintf("key%d", k), fmt.Sprintf("-%%ref1%d%%-", k))
		ref1.Set(fmt.Sprintf("ref1%d", k), fmt.Sprintf("-%%ref2%d%%-", k))
		ref1.Set(fmt.Sprintf("ref2%d", k), fmt.Sprintf("-value%d-", k))
	}

	// Fully dereferences, depending on the order in which keys are processes in DereferenceString()
	res := sut.DereferenceLoop(1, ref1)
	ExpectFalse(res, t)
	for k := 0; k < numKeys; k++ {
		key := sut.Get(fmt.Sprintf("key%d", k))
		ExpectString(fmt.Sprintf("--%%ref2%d%%--", k), *key, t)
	}
}

func TestThat_Config_DereferenceLoop_SubstitutesAllValuesInOurConfig_WhenProvidedReferenceConfigsWithMatchingKeysLimit2(t *testing.T) {
	// Setup
	sut := NewConfig()
	ref1 := NewConfig()
	numKeys := 3
	for k := 0; k < numKeys; k++ {
		sut.Set(fmt.Sprintf("key%d", k), fmt.Sprintf("-%%ref1%d%%-", k))
		ref1.Set(fmt.Sprintf("ref1%d", k), fmt.Sprintf("-%%ref2%d%%-", k))
		ref1.Set(fmt.Sprintf("ref2%d", k), fmt.Sprintf("-value%d-", k))
	}

	// Fully dereferences, depending on the order in which keys are processes in DereferenceString()
	res := sut.DereferenceLoop(2, ref1)
	ExpectFalse(res, t)
	for k := 0; k < numKeys; k++ {
		key := sut.Get(fmt.Sprintf("key%d", k))
		ExpectString(fmt.Sprintf("---value%d---", k), *key, t)
	}
}

func TestThat_Config_getReferenceKeysFromString_ReturnsExpectedKeys(t *testing.T) {
	// Setup
	sut := NewConfig()
	k1 := "k1"
	k2 := "k2"
	kstr := fmt.Sprintf("-%%%s%%-%%%s%%-", k1, k2)

	// Test
	res, err := sut.getReferenceKeysFromString(kstr)

	// Verify
	ExpectNil(err, t)
	ExpectInt(2, len(res), t)
	br1 := []byte(res[0])
	ExpectInt(len(k1), len(br1), t)
	ExpectInt(len(k1), len(res[0]), t)
	ExpectString(k1, res[0], t)
	br2 := []byte(res[1])
	ExpectInt(len(k2), len(br2), t)
	ExpectInt(len(k2), len(res[1]), t)
	ExpectString(k2, res[1], t)
}
