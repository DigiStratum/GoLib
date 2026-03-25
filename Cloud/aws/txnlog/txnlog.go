// Package txnlog provides centralized transaction logging for the DS ecosystem.
//
// All DS apps use this package to write structured transaction logs to a shared
// CloudWatch log group for monitoring, alerting, and auditing.
//
// Basic usage:
//
//	config := txnlog.Config{
//	    LogGroup: "/ds/ecosystem/transactions",
//	    AppID:    "dskanban",
//	    AppName:  "DS Projects",
//	    Env:      "prod",
//	}
//	logger, err := txnlog.New(config)
//	if err != nil {
//	    // Handle error - will fallback to local logging
//	}
//
//	// Store in context for request handling
//	ctx = txnlog.WithLogger(ctx, logger)
//
//	// Later, in request handler:
//	logger := txnlog.FromContext(ctx)
//	logger.SetRequest("POST", "/api/issues", "1700")
//	defer logger.Complete()
//
//	// ... handle request ...
//
//	logger.SetResponse(201, "Issue", 1234)
package txnlog

import (
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws/session"
)

// Config holds configuration for the transaction logger.
type Config struct {
	// LogGroup is the CloudWatch log group name (e.g., "/ds/ecosystem/transactions")
	LogGroup string

	// AppID is the application identifier (e.g., "dskanban")
	AppID string

	// AppName is the human-readable application name (e.g., "DS Projects")
	AppName string

	// Env is the environment (e.g., "prod", "dev", "staging")
	Env string

	// LambdaID is the Lambda function identifier (optional, auto-detected from AWS_LAMBDA_FUNCTION_NAME)
	LambdaID string

	// LambdaName is the Lambda function display name (optional)
	LambdaName string

	// Session is an optional AWS session. If nil, a default session is created.
	Session *session.Session

	// Client is an optional CloudWatch client. If set, Session is ignored.
	// Useful for testing with mock clients.
	Client CloudWatchClient

	// LocalFallback determines behavior when CloudWatch is unavailable.
	// If true (default), logs to stdout instead of failing.
	LocalFallback bool
}

// Logger handles transaction logging.
type Logger struct {
	entry  *LogEntry
	client CloudWatchClient
	start  time.Time
	config Config
}

// New creates a new Logger with the given configuration.
// If CloudWatch is unavailable and LocalFallback is true (default),
// it falls back to local logging.
func New(config Config) (*Logger, error) {
	// Set default for LocalFallback
	if !config.LocalFallback {
		config.LocalFallback = true
	}

	// Auto-detect Lambda environment
	if config.LambdaID == "" {
		config.LambdaID = os.Getenv("AWS_LAMBDA_FUNCTION_NAME")
	}

	// Initialize log entry with app info
	entry := newLogEntry()
	entry.App = AppInfo{
		ID:   config.AppID,
		Name: config.AppName,
	}
	entry.Env = config.Env
	if config.LambdaID != "" || config.LambdaName != "" {
		entry.Lambda = &LambdaInfo{
			ID:   config.LambdaID,
			Name: config.LambdaName,
		}
	}

	logger := &Logger{
		entry:  entry,
		start:  time.Now(),
		config: config,
	}

	// Use provided client if available
	if config.Client != nil {
		logger.client = config.Client
		return logger, nil
	}

	// Generate log stream name (per-instance, time-based)
	logStream := generateLogStreamName(config.AppID, config.Env)

	// Try to create CloudWatch client
	client, err := NewCloudWatchClient(config.Session, config.LogGroup, logStream)
	if err != nil {
		if config.LocalFallback {
			log.Printf("[txnlog] CloudWatch unavailable, using local fallback: %v", err)
			logger.client = NewLocalFallbackClient()
			return logger, nil
		}
		return nil, err
	}

	logger.client = client
	return logger, nil
}

// NewWithClient creates a Logger with a specific CloudWatch client.
// Useful for testing or custom implementations.
func NewWithClient(config Config, client CloudWatchClient) *Logger {
	entry := newLogEntry()
	entry.App = AppInfo{
		ID:   config.AppID,
		Name: config.AppName,
	}
	entry.Env = config.Env
	if config.LambdaID != "" || config.LambdaName != "" {
		entry.Lambda = &LambdaInfo{
			ID:   config.LambdaID,
			Name: config.LambdaName,
		}
	}

	return &Logger{
		entry:  entry,
		client: client,
		start:  time.Now(),
		config: config,
	}
}

// Clone creates a new Logger with the same config but fresh entry.
// Useful for creating per-request loggers from a shared config.
func (l *Logger) Clone() *Logger {
	return NewWithClient(l.config, l.client)
}

// SetRequest sets request information.
func (l *Logger) SetRequest(method, path, resourceID string) {
	l.entry.Request.Method = method
	l.entry.Request.Path = path
	l.entry.Request.ResourceID = resourceID
}

// SetEventType sets the event type for non-HTTP transactions.
func (l *Logger) SetEventType(eventType string) {
	l.entry.Request.EventType = eventType
}

// SetResponse sets response information.
func (l *Logger) SetResponse(status int, dataType string, bytes int) {
	l.ensureResponse()
	l.entry.Response.Status = status
	l.entry.Response.Type = dataType
	l.entry.Response.Bytes = bytes
}

// SetSession sets session and correlation IDs.
func (l *Logger) SetSession(sessionID, correlationID string) {
	l.entry.SessionID = sessionID
	l.entry.CorrelationID = correlationID
}

// SetTenant sets tenant information.
func (l *Logger) SetTenant(tenantID string) {
	l.entry.TenantID = tenantID
}

// SetUser sets user information (hashed).
func (l *Logger) SetUser(userIDHash string) {
	l.entry.UserIDHash = userIDHash
}

// SetPerf sets performance metrics.
func (l *Logger) SetPerf(coldStart bool, initMs, dbQueryMs, externalMs int64) {
	l.entry.Perf = &PerfInfo{
		ColdStart:  coldStart,
		InitMs:     initMs,
		DBQueryMs:  dbQueryMs,
		ExternalMs: externalMs,
	}
}

// SetColdStart marks this as a cold start request.
func (l *Logger) SetColdStart(coldStart bool) {
	l.ensurePerf()
	l.entry.Perf.ColdStart = coldStart
}

// SetSource sets source information (IP, user agent, etc.).
func (l *Logger) SetSource(ip, userAgent, referer string) {
	l.entry.Source = &SourceInfo{
		IP:        ip,
		UserAgent: userAgent,
		Referer:   referer,
	}
}

// ensureResponse ensures the Response pointer is initialized.
func (l *Logger) ensureResponse() {
	if l.entry.Response == nil {
		l.entry.Response = &ResponseInfo{}
	}
}

// ensurePerf ensures the Perf pointer is initialized.
func (l *Logger) ensurePerf() {
	if l.entry.Perf == nil {
		l.entry.Perf = &PerfInfo{}
	}
}

// Entry returns the current log entry for inspection or modification.
func (l *Logger) Entry() *LogEntry {
	return l.entry
}

// Complete calculates duration and writes the log entry.
// This should be called at the end of request processing, typically via defer.
func (l *Logger) Complete() {
	// Calculate duration
	duration := time.Since(l.start)
	l.ensureResponse()
	l.entry.Response.DurationMs = duration.Milliseconds()

	// Write to CloudWatch (or fallback)
	if err := l.client.WriteLog(l.entry); err != nil {
		log.Printf("[txnlog] failed to write log: %v", err)
	}
}

// generateLogStreamName creates a unique log stream name.
func generateLogStreamName(appID, env string) string {
	// Format: appid/env/YYYY/MM/DD/instance-id
	// Using timestamp as instance ID for simplicity
	now := time.Now().UTC()
	return now.Format(appID + "/" + env + "/2006/01/02/" + "150405.000")
}
