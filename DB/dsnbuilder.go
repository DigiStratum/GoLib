package db

/*

ref: https://github.com/go-sql-driver/mysql/blob/master/dsn.go
*/

import (
	//"fmt"

	"github.com/go-sql-driver/mysql"
)

type DSNBuilderIfc interface {
	Configure(config cfg.ConfigIfc) error
	SetUser(user string) *DSNBuilder
	SetPasswd(passwd string) *DSNBuilder
	SetNet(net string) *DSNBuilder
	SetAddr(addr string) *DSNBuilder
	SetDBName(name string) *DSNBuilder
	SetParams(params map[string]string) *DSNBuilder
	SetCollation(collation string) *DSNBuilder
	SetLoc(loc *time.Location) *DSNBuilder
	SetMaxAllowedPacket(maxAllowedPacket int) *DSNBuilder
	SetServerPubKey(serverPubKey string) *DSNBuilder
	SetTLSConfig(tlsConfig string) *DSNBuilder
	SetTimeout(timeout time.Duration) *DSNBuilder
	SetReadTimeout(readTimeout time.Duration) *DSNBuilder
	SetWriteTimeout(writeTimeout time.Duration) *DSNBuilder
	SetAllowAllFiles(allowAllFiles bool) *DSNBuilder
	SetAllowCleartextPasswords(allowCleartextPasswords bool) *DSNBuilder
	SetAllowNativePasswords(allowNativePasswords bool) *DSNBuilder
	SetAllowOldPasswords(allowOldPasswords bool) *DSNBuilder
	SetCheckConnLiveness(checkConnLiveness bool) *DSNBuilder
	SetClientFoundRows(clientFoundRows bool) *DSNBuilder
	SetColumnsWithAlias(columnsWithAlias bool) *DSNBuilder
	SetInterpolateParams(interpolateParams bool) *DSNBuilder
	SetMultiStatements(multiStatements bool) *DSNBuilder
	SetParseTime(parseTime bool) *DSNBuilder
	SetRejectReadOnly(rejectReadOnly bool) *DSNBuilder
	Build() (*DSN, error)
}

type DSNBuilder struct {
	dsnConfig		mysql.Config
}

func BuildDSN() *DSNBuilder {
	return &DNSBuilder{
		dsnConfig:	mysql.NewConfig(),
	}
}

// -------------------------------------------------------------------------------------------------
// ConfigurableIfc Public Interface
// -------------------------------------------------------------------------------------------------

func (r *DSNBuilder) Configure(config cfg.ConfigIfc) error {
	keys := config.GetKeys()
	params := make(map[string]string)
	for _, key := range keys {
		value := config.Get(key)
		if nil == value { continue }
		boolValue := config.GetBool(key)

		switch (key) {
			"User": r.SetUser(value)
			"Passwd": r.SetPasswd(value)
			"Net": r.SetNet(value)
			"Addr": r.SetAddr(value)
			"DBName": r.SetDBName(value)
			"Collation": r.SetCollation(value)
			"Loc": r.SetLoc(value) // *time.Location
			"MaxAllowedPacket":
				int64Value := config.GetIn64(key);
				if nil == int64Value { continue }
				r.SetMaxAllowedPacket(int(*int64Value))
			"ServerPubKey": r.SetServerPubKey(value)
			"TLSConfig": r.SetTLSConfig(value)
			"Timeout": r.SetTimeout(value)
				// time.Duration
			"ReadTimeout": r.SetReadTimeout(value)
				// time.Duration
			"WriteTimeout": r.SetWriteTimeout(value)
				// time.Duration
			"AllowAllFiles": r.SetAllowAllFiles(boolValue)
			"AllowCleartextPasswords": r.SetAllowCleartextPasswords(boolValue)
			"AllowNativePasswords": r.SetAllowNativePasswords(boolValue)
			"AllowOldPasswords": r.SetAllowOldPasswords(boolValue)
			"CheckConnLiveness": r.SetCheckConnLiveness(boolValue)
			"ClientFoundRows": r.SetClientFoundRows(boolValue)
			"ColumnsWithAlias": r.SetColumnsWithAlias(boolValue)
			"InterpolateParams": r.SetInterpolateParams(boolValue)
			"MultiStatements": r.SetMultiStatements(boolValue)
			"ParseTime": r.SetParseTime(boolValue)
			"RejectReadOnly": r.SetRejectReadOnly(boolValue)
			default:
				// Anything else to be treated as a name-value Param
				params[key] = value
		}
	}
	if len(params) > 0 { r.SetParams(params) }

	return nil
}

// -------------------------------------------------------------------------------------------------
// DSNBuilderIfc Public Interface
// -------------------------------------------------------------------------------------------------

// Username
func (r *DSNBuilder) SetUser(user string) *DSNBuilder { return r }

// Password (requires User)
func (r *DSNBuilder) SetPasswd(passwd string) *DSNBuilder { return r }

// Network type
func (r *DSNBuilder) SetNet(net string) *DSNBuilder { return r }

// Network address (requires Net)
func (r *DSNBuilder) SetAddr(addr string) *DSNBuilder { return r }

// Database name
func (r *DSNBuilder) SetDBName(name string) *DSNBuilder { return r }

// Connection parameters
func (r *DSNBuilder) SetParams(params map[string]string) *DSNBuilder { return r }

// Connection collation
func (r *DSNBuilder) SetCollation(collation string) *DSNBuilder { return r }

// Location for time.Time values
func (r *DSNBuilder) SetLoc(loc *time.Location) *DSNBuilder { return r }

// Max packet size allowed
func (r *DSNBuilder) SetMaxAllowedPacket(maxAllowedPacket int) *DSNBuilder { return r }

// Server public key name
func (r *DSNBuilder) SetServerPubKey(serverPubKey string) *DSNBuilder { return r }

// TLS configuration name
func (r *DSNBuilder) SetTLSConfig(tlsConfig string) *DSNBuilder { return r }

// Dial timeout
func (r *DSNBuilder) SetTimeout(timeout time.Duration) *DSNBuilder { return r }

// I/O read timeout
func (r *DSNBuilder) SetReadTimeout(readTimeout time.Duration) *DSNBuilder { return r }

// I/O write timeout
func (r *DSNBuilder) SetWriteTimeout(writeTimeout time.Duration) *DSNBuilder { return r }

// Allow all files to be used with LOAD DATA LOCAL INFILE
func (r *DSNBuilder) SetAllowAllFiles(allowAllFiles bool) *DSNBuilder { return r }

// Allows the cleartext client side plugin
func (r *DSNBuilder) SetAllowCleartextPasswords(allowCleartextPasswords bool) *DSNBuilder { return r }

// Allows the native password authentication method
func (r *DSNBuilder) SetAllowNativePasswords(allowNativePasswords bool) *DSNBuilder { return r }

// Allows the old insecure password method
func (r *DSNBuilder) SetAllowOldPasswords(allowOldPasswords bool) *DSNBuilder { return r }

// Check connections for liveness before using them
func (r *DSNBuilder) SetCheckConnLiveness(checkConnLiveness bool) *DSNBuilder { return r }

// Return number of matching rows instead of rows changed
func (r *DSNBuilder) SetClientFoundRows(clientFoundRows bool) *DSNBuilder { return r }

// Prepend table alias to column names
func (r *DSNBuilder) SetColumnsWithAlias(columnsWithAlias bool) *DSNBuilder { return r }

// Interpolate placeholders into query string
func (r *DSNBuilder) SetInterpolateParams(interpolateParams bool) *DSNBuilder { return r }

// Allow multiple statements in one query
func (r *DSNBuilder) SetMultiStatements(multiStatements bool) *DSNBuilder { return r }

// Parse time values to time.Time
func (r *DSNBuilder) SetParseTime(parseTime bool) *DSNBuilder { return r }

// Reject read-only connections
func (r *DSNBuilder) SetRejectReadOnly(rejectReadOnly bool) *DSNBuilder { return r }

func (r DSNBuilder) Build() (*DSN, error) {
	dsn := r.dsnConfig.FormatDSN()
	return NewDSN(dsn)
}