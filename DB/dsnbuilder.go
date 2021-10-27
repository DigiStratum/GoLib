package db

/*

ref: https://github.com/go-sql-driver/mysql/blob/master/dsn.go
*/

import (
	"time"

	"github.com/go-sql-driver/mysql"

	cfg "github.com/DigiStratum/GoLib/Config"
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
	dsnConfig		*mysql.Config
}

// -------------------------------------------------------------------------------------------------
// Factory functions
// -------------------------------------------------------------------------------------------------

func BuildDSN() *DSNBuilder {
	return &DSNBuilder{
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
		v := config.Get(key)
		if nil == v { continue }
		value := *v
		boolValue := config.GetBool(key)

		switch (key) {
			case "User": r.SetUser(value)
			case "Passwd": r.SetPasswd(value)
			case "Net": r.SetNet(value)
			case "Addr": r.SetAddr(value)
			case "DBName": r.SetDBName(value)
			case "Collation": r.SetCollation(value)
			case "Loc":
				loc, err := time.LoadLocation(value)
				if nil != err { return err }
				r.SetLoc(loc)
			case "MaxAllowedPacket":
				int64Value := config.GetInt64(key);
				if nil == int64Value { continue }
				r.SetMaxAllowedPacket(int(*int64Value))
			case "ServerPubKey": r.SetServerPubKey(value)
			case "TLSConfig": r.SetTLSConfig(value)
			case "Timeout":
				dur, err := time.ParseDuration(value)
				if nil != err { return err }
				r.SetTimeout(dur)
			case "ReadTimeout":
				dur, err := time.ParseDuration(value)
				if nil != err { return err }
				r.SetReadTimeout(dur)
			case "WriteTimeout":
				dur, err := time.ParseDuration(value)
				if nil != err { return err }
				r.SetWriteTimeout(dur)
			case "AllowAllFiles": r.SetAllowAllFiles(boolValue)
			case "AllowCleartextPasswords": r.SetAllowCleartextPasswords(boolValue)
			case "AllowNativePasswords": r.SetAllowNativePasswords(boolValue)
			case "AllowOldPasswords": r.SetAllowOldPasswords(boolValue)
			case "CheckConnLiveness": r.SetCheckConnLiveness(boolValue)
			case "ClientFoundRows": r.SetClientFoundRows(boolValue)
			case "ColumnsWithAlias": r.SetColumnsWithAlias(boolValue)
			case "InterpolateParams": r.SetInterpolateParams(boolValue)
			case "MultiStatements": r.SetMultiStatements(boolValue)
			case "ParseTime": r.SetParseTime(boolValue)
			case "RejectReadOnly": r.SetRejectReadOnly(boolValue)
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
func (r *DSNBuilder) SetUser(user string) *DSNBuilder {
	r.dsnConfig.User = user;
	return r
}

// Password (requires User)
func (r *DSNBuilder) SetPasswd(passwd string) *DSNBuilder {
	r.dsnConfig.Passwd = passwd;
	return r
}

// Network type
func (r *DSNBuilder) SetNet(net string) *DSNBuilder {
	r.dsnConfig.Net = net;
	return r
}

// Network address (requires Net)
func (r *DSNBuilder) SetAddr(addr string) *DSNBuilder {
	r.dsnConfig.Addr = addr;
	return r
}

// Database name
func (r *DSNBuilder) SetDBName(name string) *DSNBuilder {
	r.dsnConfig.DBName = name;
	return r
}

// Connection parameters
func (r *DSNBuilder) SetParams(params map[string]string) *DSNBuilder {
	r.dsnConfig.Params = params;
	return r
}

// Connection collation
func (r *DSNBuilder) SetCollation(collation string) *DSNBuilder {
	r.dsnConfig.Collation = collation;
	return r
}

// Location for time.Time values
func (r *DSNBuilder) SetLoc(loc *time.Location) *DSNBuilder {
	r.dsnConfig.Loc = loc;
	return r
}

// Max packet size allowed
func (r *DSNBuilder) SetMaxAllowedPacket(maxAllowedPacket int) *DSNBuilder {
	r.dsnConfig.MaxAllowedPacket = maxAllowedPacket;
	return r
}

// Server public key name
func (r *DSNBuilder) SetServerPubKey(serverPubKey string) *DSNBuilder {
	r.dsnConfig.ServerPubKey = serverPubKey;
	return r
}

// TLS configuration name
func (r *DSNBuilder) SetTLSConfig(tlsConfig string) *DSNBuilder {
	r.dsnConfig.TLSConfig = tlsConfig;
	return r
}

// Dial timeout
func (r *DSNBuilder) SetTimeout(timeout time.Duration) *DSNBuilder {
	r.dsnConfig.Timeout = timeout;
	return r
}

// I/O read timeout
func (r *DSNBuilder) SetReadTimeout(readTimeout time.Duration) *DSNBuilder {
	r.dsnConfig.ReadTimeout = readTimeout;
	return r
}

// I/O write timeout
func (r *DSNBuilder) SetWriteTimeout(writeTimeout time.Duration) *DSNBuilder {
	r.dsnConfig.WriteTimeout = writeTimeout;
	return r
}

// Allow all files to be used with LOAD DATA LOCAL INFILE
func (r *DSNBuilder) SetAllowAllFiles(allowAllFiles bool) *DSNBuilder {
	r.dsnConfig.AllowAllFiles = allowAllFiles;
	return r
}

// Allows the cleartext client side plugin
func (r *DSNBuilder) SetAllowCleartextPasswords(allowCleartextPasswords bool) *DSNBuilder {
	r.dsnConfig.AllowCleartextPasswords = allowCleartextPasswords;
	return r
}

// Allows the native password authentication method
func (r *DSNBuilder) SetAllowNativePasswords(allowNativePasswords bool) *DSNBuilder {
	r.dsnConfig.AllowNativePasswords = allowNativePasswords;
	return r
}

// Allows the old insecure password method
func (r *DSNBuilder) SetAllowOldPasswords(allowOldPasswords bool) *DSNBuilder {
	r.dsnConfig.AllowOldPasswords = allowOldPasswords;
	return r
}

// Check connections for liveness before using them
func (r *DSNBuilder) SetCheckConnLiveness(checkConnLiveness bool) *DSNBuilder {
	r.dsnConfig.CheckConnLiveness = checkConnLiveness;
	return r
}

// Return number of matching rows instead of rows changed
func (r *DSNBuilder) SetClientFoundRows(clientFoundRows bool) *DSNBuilder {
	r.dsnConfig.ClientFoundRows = clientFoundRows;
	return r
}

// Prepend table alias to column names
func (r *DSNBuilder) SetColumnsWithAlias(columnsWithAlias bool) *DSNBuilder {
	r.dsnConfig.ColumnsWithAlias = columnsWithAlias;
	return r
}

// Interpolate placeholders into query string
func (r *DSNBuilder) SetInterpolateParams(interpolateParams bool) *DSNBuilder {
	r.dsnConfig.InterpolateParams = interpolateParams;
	return r
}

// Allow multiple statements in one query
func (r *DSNBuilder) SetMultiStatements(multiStatements bool) *DSNBuilder {
	r.dsnConfig.MultiStatements = multiStatements;
	return r
}

// Parse time values to time.Time
func (r *DSNBuilder) SetParseTime(parseTime bool) *DSNBuilder {
	r.dsnConfig.ParseTime = parseTime;
	return r
}

// Reject read-only connections
func (r *DSNBuilder) SetRejectReadOnly(rejectReadOnly bool) *DSNBuilder {
	r.dsnConfig.RejectReadOnly = rejectReadOnly;
	return r
}

func (r DSNBuilder) Build() (*DSN, error) {
	dsn := r.dsnConfig.FormatDSN()
	return NewDSN(dsn)
}
