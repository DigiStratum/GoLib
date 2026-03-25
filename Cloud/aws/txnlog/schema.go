package txnlog

import (
	"time"
)

// LogEntry represents a single transaction log entry (schema v1).
// Optional fields use omitempty to reduce log size.
// Pointer types are used for optional struct fields to enable proper omitempty behavior.
type LogEntry struct {
	Version       int           `json:"v"`
	Timestamp     time.Time     `json:"ts"`
	App           AppInfo       `json:"app"`
	Lambda        *LambdaInfo   `json:"lambda,omitempty"`
	Env           string        `json:"env"`
	SessionID     string        `json:"session_id,omitempty"`
	CorrelationID string        `json:"correlation_id,omitempty"`
	TenantID      string        `json:"tenant_id,omitempty"`
	UserIDHash    string        `json:"user_id_hash,omitempty"`
	Request       RequestInfo   `json:"request"`
	Response      *ResponseInfo `json:"response,omitempty"`
	Perf          *PerfInfo     `json:"perf,omitempty"`
	Source        *SourceInfo   `json:"source,omitempty"`
}

// AppInfo identifies the application.
type AppInfo struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// LambdaInfo identifies the Lambda function (if applicable).
type LambdaInfo struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

// RequestInfo describes the incoming request.
type RequestInfo struct {
	Method     string `json:"method,omitempty"`
	Path       string `json:"path,omitempty"`
	ResourceID string `json:"resource_id,omitempty"`
	EventType  string `json:"event_type,omitempty"`
}

// ResponseInfo describes the response.
type ResponseInfo struct {
	Status     int    `json:"status,omitempty"`
	Type       string `json:"type,omitempty"`
	Bytes      int    `json:"bytes,omitempty"`
	DurationMs int64  `json:"duration_ms,omitempty"`
}

// PerfInfo contains performance metrics.
type PerfInfo struct {
	ColdStart   bool  `json:"cold_start,omitempty"`
	InitMs      int64 `json:"init_ms,omitempty"`
	DBQueryMs   int64 `json:"db_query_ms,omitempty"`
	ExternalMs  int64 `json:"external_ms,omitempty"`
}

// SourceInfo contains optional source context.
type SourceInfo struct {
	IP        string `json:"ip,omitempty"`
	UserAgent string `json:"user_agent,omitempty"`
	Referer   string `json:"referer,omitempty"`
}

// IsEmpty returns true if the LogEntry has no meaningful data set.
func (e *LogEntry) IsEmpty() bool {
	return e.App.ID == "" && e.Request.Method == "" && e.Request.EventType == ""
}

// newLogEntry creates a LogEntry with version 1 and current timestamp.
func newLogEntry() *LogEntry {
	return &LogEntry{
		Version:   1,
		Timestamp: time.Now().UTC(),
	}
}
