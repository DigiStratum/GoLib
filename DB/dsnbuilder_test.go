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

func TestThat_SetUser_AddsUserToDSN(t *testing.T) {
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

func TestThat_SetPasswd_AddsPasswdToDSNWhenUserSpecified(t *testing.T) {
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

func TestThat_SetPasswd_OmitsPasswdFromDSNWithoutUserSpecified(t *testing.T) {
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

func TestThat_SetNet_AddsNetToDSN(t *testing.T) {
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

func TestThat_SetAddr_AddsAddrToDSN_WhenNetSupplied(t *testing.T) {
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

func TestThat_SetDBName_AddsDBNameToDSN(t *testing.T) {
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

func TestThat_SetParams_AddsParamsToDSN(t *testing.T) {
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

func TestThat_SetCollation_AddsCollationToDSN(t *testing.T) {
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

func TestThat_SetLoc_AddsLocToDSN_When_NonUTCLocation(t *testing.T) {
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

func TestThat_SetMaxAllowedPacket_AddsMaxAllowedPacketToDSN(t *testing.T) {
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

func TestThat_SetServerPubKey_AddsPubKeyToDSN(t *testing.T) {
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

//SetTLSConfig(tlsConfig string) *DSNBuilder
type myStruct struct{}
func (myStruct) Read(b []byte) (n int, err error) {
	for i := range b { b[i] = 0 }
	return len(b), nil
}
func TestThat_SetTLSConfig_AddsTLSConfigToDSN(t *testing.T) {
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

//SetTimeout(timeout time.Duration) *DSNBuilder
//SetReadTimeout(readTimeout time.Duration) *DSNBuilder
//SetWriteTimeout(writeTimeout time.Duration) *DSNBuilder
//SetAllowAllFiles(allowAllFiles bool) *DSNBuilder
//SetAllowCleartextPasswords(allowCleartextPasswords bool) *DSNBuilder
//SetAllowNativePasswords(allowNativePasswords bool) *DSNBuilder
//SetAllowOldPasswords(allowOldPasswords bool) *DSNBuilder
//SetCheckConnLiveness(checkConnLiveness bool) *DSNBuilder
//SetClientFoundRows(clientFoundRows bool) *DSNBuilder
//SetColumnsWithAlias(columnsWithAlias bool) *DSNBuilder
//SetInterpolateParams(interpolateParams bool) *DSNBuilder
//SetMultiStatements(multiStatements bool) *DSNBuilder
//SetParseTime(parseTime bool) *DSNBuilder
//SetRejectReadOnly(rejectReadOnly bool) *DSNBuilder


