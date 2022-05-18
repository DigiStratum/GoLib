package db

import(
	"time"
	"testing"
	"math/big"
	"crypto/rsa"
	"crypto/tls"

	"github.com/go-sql-driver/mysql"
	. "github.com/DigiStratum/GoLib/Testing"
)

func TestThat_BuildDSN_ReturnsSomething(t *testing.T) {
	// Setup
	var sut *DSNBuilder

	// Test
	sut = BuildDSN()

	// Verify
	ExpectNonNil(sut, t)
}

func TestThat_DSNBuilder_SetUser_AddsUserToDSN(t *testing.T) {
	// Setup
	var sut *DSNBuilder = BuildDSN()
	var actual *DSN
	var err error

	// Test
	actual, err = sut.SetUser("testuser").Build()

	// Verify
	ExpectNonNil(actual, t)
	ExpectNoError(err, t)
	ExpectString("testuser@/", actual.ToString(), t)
}

func TestThat_DSNBuilder_SetPasswd_AddsPasswdToDSNWhenUserSpecified(t *testing.T) {
	// Setup
	var sut *DSNBuilder = BuildDSN()
	var actual *DSN
	var err error

	// Test
	actual, err = sut.SetUser("testuser").SetPasswd("testpass").Build()

	// Verify
	ExpectNonNil(actual, t)
	ExpectNoError(err, t)
	ExpectString("testuser:testpass@/", actual.ToString(), t)
}

func TestThat_DSNBuilder_SetPasswd_OmitsPasswdFromDSNWithoutUserSpecified(t *testing.T) {
	// Setup
	var sut *DSNBuilder = BuildDSN()
	var actual *DSN
	var err error

	// Test
	actual, err = sut.SetPasswd("testpass").Build()

	// Verify
	ExpectNonNil(actual, t)
	ExpectNoError(err, t)
	ExpectString("/", actual.ToString(), t)
}

func TestThat_DSNBuilder_SetNet_AddsNetToDSN(t *testing.T) {
	// Setup
	var sut *DSNBuilder = BuildDSN()
	var actual *DSN
	var err error

	// Test
	actual, err = sut.SetNet("tcp").Build()

	// Verify
	ExpectNonNil(actual, t)
	ExpectNoError(err, t)
	ExpectString("tcp/", actual.ToString(), t)
}

func TestThat_DSNBuilder_SetAddr_AddsAddrToDSN_WhenNetSupplied(t *testing.T) {
	// Setup
	var sut *DSNBuilder = BuildDSN()
	var actual *DSN
	var err error

	// Test
	actual, err = sut.SetNet("tcp").SetAddr("1.2.3.4:3306").Build()

	// Verify
	ExpectNonNil(actual, t)
	ExpectNoError(err, t)
	ExpectString("tcp(1.2.3.4:3306)/", actual.ToString(), t)
}

func TestThat_DSNBuilder_SetDBName_AddsDBNameToDSN(t *testing.T) {
	// Setup
	var sut *DSNBuilder = BuildDSN()
	var actual *DSN
	var err error

	// Test
	actual, err = sut.SetDBName("bogusname").Build()

	// Verify
	ExpectNonNil(actual, t)
	ExpectNoError(err, t)
	ExpectString("/bogusname", actual.ToString(), t)
}

func TestThat_DSNBuilder_SetParams_AddsParamsToDSN(t *testing.T) {
	// Setup
	var sut *DSNBuilder = BuildDSN()
	var actual *DSN
	var err error
	params := make(map[string]string)
	params["paramname"] = "paramvalue"

	// Test
	actual, err = sut.SetParams(params).Build()

	// Verify
	ExpectNonNil(actual, t)
	ExpectNoError(err, t)
	ExpectString("/?paramname=paramvalue", actual.ToString(), t)
}

func TestThat_DSNBuilder_SetCollation_AddsCollationToDSN(t *testing.T) {
	// Setup
	var sut *DSNBuilder = BuildDSN()
	var actual *DSN
	var err error

	// Test
	actual, err = sut.SetCollation("utf8_general_ci").Build()

	// Verify
	ExpectNonNil(actual, t)
	ExpectNoError(err, t)
	ExpectString("/?collation=utf8_general_ci", actual.ToString(), t)
}

func TestThat_DSNBuilder_SetLoc_AddsLocToDSN_When_NonUTCLocation(t *testing.T) {
	// Setup
	var sut *DSNBuilder = BuildDSN()
	var actual *DSN
	var err error

	// Test
	actual, err = sut.SetLoc(time.Local).Build()

	// Verify
	ExpectNonNil(actual, t)
	ExpectNoError(err, t)
	ExpectString("/?loc=Local", actual.ToString(), t)
}

func TestThat_DSNBuilder_SetMaxAllowedPacket_AddsMaxAllowedPacketToDSN(t *testing.T) {
	// Setup
	var sut *DSNBuilder = BuildDSN()
	var actual *DSN
	var err error

	// Test
	actual, err = sut.SetMaxAllowedPacket(333).Build()

	// Verify
	ExpectNonNil(actual, t)
	ExpectNoError(err, t)
	ExpectString("/?maxAllowedPacket=333", actual.ToString(), t)
}

func TestThat_DSNBuilder_SetServerPubKey_AddsPubKeyToDSN(t *testing.T) {
	// Setup
	var sut *DSNBuilder = BuildDSN()
	var actual *DSN
	var err error
	var modulus *big.Int = big.NewInt(int64(333333))
	var pubkey rsa.PublicKey = rsa.PublicKey{ N: modulus, E: 333 }
	mysql.RegisterServerPubKey("boguspubkey", &pubkey)

	// Test
	actual, err = sut.SetServerPubKey("boguspubkey").Build()

	// Verify
	ExpectNonNil(actual, t)
	ExpectNoError(err, t)
	ExpectString("/?serverPubKey=boguspubkey", actual.ToString(), t)
}

type myStruct struct{}
func (myStruct) Read(b []byte) (n int, err error) {
	for i := range b { b[i] = 0 }
	return len(b), nil
}

func TestThat_DSNBuilder_SetTLSConfig_AddsTLSConfigToDSN(t *testing.T) {
	// Setup
	var sut *DSNBuilder = BuildDSN()
	var actual *DSN
	var err error
	tlsConfig := &tls.Config{ Rand: myStruct{} }
	mysql.RegisterTLSConfig("bogustlsconfig", tlsConfig)

	// Test
	actual, err = sut.SetTLSConfig("bogustlsconfig").Build()

	// Verify
	ExpectNonNil(actual, t)
	ExpectNoError(err, t)
	ExpectString("/?tls=bogustlsconfig", actual.ToString(), t)
}

func TestThat_DSNBuilder_SetTimeout_AddsTimeoutToDSN(t *testing.T) {
	// Setup
	var sut *DSNBuilder = BuildDSN()
	var actual *DSN
	var err error
	timeout := time.Duration(333)

	// Test
	actual, err = sut.SetTimeout(timeout).Build()

	// Verify
	ExpectNonNil(actual, t)
	ExpectNoError(err, t)
	ExpectString("/?timeout=333ns", actual.ToString(), t)
}

func TestThat_DSNBuilder_SetReadTimeout_AddsReadTimeoutToDSN(t *testing.T) {
	// Setup
	var sut *DSNBuilder = BuildDSN()
	var actual *DSN
	var err error
	timeout := time.Duration(333)

	// Test
	actual, err = sut.SetReadTimeout(timeout).Build()

	// Verify
	ExpectNonNil(actual, t)
	ExpectNoError(err, t)
	ExpectString("/?readTimeout=333ns", actual.ToString(), t)
}

func TestThat_DSNBuilder_SetWriteTimeout_AddsWriteTimeoutToDSN(t *testing.T) {
	// Setup
	var sut *DSNBuilder = BuildDSN()
	var actual *DSN
	var err error
	timeout := time.Duration(333)

	// Test
	actual, err = sut.SetWriteTimeout(timeout).Build()

	// Verify
	ExpectNonNil(actual, t)
	ExpectNoError(err, t)
	ExpectString("/?writeTimeout=333ns", actual.ToString(), t)
}

func TestThat_DSNBuilder_SetAllowAllFiles_AddsAllowAllFilesToDSN(t *testing.T) {
	// Setup
	var sut *DSNBuilder = BuildDSN()
	var actual *DSN
	var err error

	// Test
	// Default is false, so override with true to make it show up
	actual, err = sut.SetAllowAllFiles(true).Build()

	// Verify
	ExpectNonNil(actual, t)
	ExpectNoError(err, t)
	ExpectString("/?allowAllFiles=true", actual.ToString(), t)
}

func TestThat_DSNBuilder_SetAllowCleartextPasswords_AddsAllowCleartextPasswordsToDSN(t *testing.T) {
	// Setup
	var sut *DSNBuilder = BuildDSN()
	var actual *DSN
	var err error

	// Test
	// Default is false, so override with true to make it show up
	actual, err = sut.SetAllowCleartextPasswords(true).Build()

	// Verify
	ExpectNonNil(actual, t)
	ExpectNoError(err, t)
	ExpectString("/?allowCleartextPasswords=true", actual.ToString(), t)
}

func TestThat_DSNBuilder_SetAllowNativePasswords_AddsAllowNativePasswordsToDSN(t *testing.T) {
	// Setup
	var sut *DSNBuilder = BuildDSN()
	var actual *DSN
	var err error

	// Test
	// Default is true, so override with false to make it show up
	actual, err = sut.SetAllowNativePasswords(false).Build()

	// Verify
	ExpectNonNil(actual, t)
	ExpectNoError(err, t)
	ExpectString("/?allowNativePasswords=false", actual.ToString(), t)
}

func TestThat_DSNBuilder_SetAllowOldPasswords_AddsAllowOldPasswordsToDSN(t *testing.T) {
	// Setup
	var sut *DSNBuilder = BuildDSN()
	var actual *DSN
	var err error

	// Test
	// Default is false, so override with true to make it show up
	actual, err = sut.SetAllowOldPasswords(true).Build()

	// Verify
	ExpectNonNil(actual, t)
	ExpectNoError(err, t)
	ExpectString("/?allowOldPasswords=true", actual.ToString(), t)
}

func TestThat_DSNBuilder_SetCheckConnLiveness_AddsCheckConnLivenessToDSN(t *testing.T) {
	// Setup
	var sut *DSNBuilder = BuildDSN()
	var actual *DSN
	var err error

	// Test
	// Default is true, so override with false to make it show up
	actual, err = sut.SetCheckConnLiveness(false).Build()

	// Verify
	ExpectNonNil(actual, t)
	ExpectNoError(err, t)
	ExpectString("/?checkConnLiveness=false", actual.ToString(), t)
}

func TestThat_DSNBuilder_SetClientFoundRows_AddsClientFoundRowsToDSN(t *testing.T) {
	// Setup
	var sut *DSNBuilder = BuildDSN()
	var actual *DSN
	var err error

	// Test
	// Default is false, so override with true to make it show up
	actual, err = sut.SetClientFoundRows(true).Build()

	// Verify
	ExpectNonNil(actual, t)
	ExpectNoError(err, t)
	ExpectString("/?clientFoundRows=true", actual.ToString(), t)
}

func TestThat_DSNBuilder_SetColumnsWithAlias_AddsColumnsWithAliasToDSN(t *testing.T) {
	// Setup
	var sut *DSNBuilder = BuildDSN()
	var actual *DSN
	var err error

	// Test
	// Default is false, so override with true to make it show up
	actual, err = sut.SetColumnsWithAlias(true).Build()

	// Verify
	ExpectNonNil(actual, t)
	ExpectNoError(err, t)
	ExpectString("/?columnsWithAlias=true", actual.ToString(), t)
}

func TestThat_DSNBuilder_SetInterpolateParams_AddsInterpolateParamsToDSN(t *testing.T) {
	// Setup
	var sut *DSNBuilder = BuildDSN()
	var actual *DSN
	var err error

	// Test
	// Default is false, so override with true to make it show up
	actual, err = sut.SetInterpolateParams(true).Build()

	// Verify
	ExpectNonNil(actual, t)
	ExpectNoError(err, t)
	ExpectString("/?interpolateParams=true", actual.ToString(), t)
}

func TestThat_DSNBuilder_SetMultiStatements_AddsMultiStatementsToDSN(t *testing.T) {
	// Setup
	var sut *DSNBuilder = BuildDSN()
	var actual *DSN
	var err error

	// Test
	// Default is false, so override with true to make it show up
	actual, err = sut.SetMultiStatements(true).Build()

	// Verify
	ExpectNonNil(actual, t)
	ExpectNoError(err, t)
	ExpectString("/?multiStatements=true", actual.ToString(), t)
}

func TestThat_DSNBuilder_SetParseTime_AddsParseTimeToDSN(t *testing.T) {
	// Setup
	var sut *DSNBuilder = BuildDSN()
	var actual *DSN
	var err error

	// Test
	// Default is false, so override with true to make it show up
	actual, err = sut.SetParseTime(true).Build()

	// Verify
	ExpectNonNil(actual, t)
	ExpectNoError(err, t)
	ExpectString("/?parseTime=true", actual.ToString(), t)
}

func TestThat_DSNBuilder_SetRejectReadOnly_AddsRejectReadOnlyToDSN(t *testing.T) {
	// Setup
	var sut *DSNBuilder = BuildDSN()
	var actual *DSN
	var err error

	// Test
	// Default is false, so override with true to make it show up
	actual, err = sut.SetRejectReadOnly(true).Build()

	// Verify
	ExpectNonNil(actual, t)
	ExpectNoError(err, t)
	ExpectString("/?rejectReadOnly=true", actual.ToString(), t)
}



