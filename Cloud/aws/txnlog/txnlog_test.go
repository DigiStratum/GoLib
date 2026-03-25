package txnlog

import (
	"context"
	"encoding/json"
	"testing"
	"time"
)

// mockClient records log entries for testing.
type mockClient struct {
	entries []*LogEntry
}

func (m *mockClient) WriteLog(entry *LogEntry) error {
	m.entries = append(m.entries, entry)
	return nil
}

func TestNew(t *testing.T) {
	mock := &mockClient{}
	config := Config{
		LogGroup:   "/ds/ecosystem/transactions",
		AppID:      "testapp",
		AppName:    "Test App",
		Env:        "test",
		LambdaID:   "testapp-api-test",
		LambdaName: "api",
		Client:     mock,
	}

	logger, err := New(config)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	if logger == nil {
		t.Fatal("New() returned nil logger")
	}

	if logger.entry.App.ID != "testapp" {
		t.Errorf("App.ID = %q, want %q", logger.entry.App.ID, "testapp")
	}

	if logger.entry.App.Name != "Test App" {
		t.Errorf("App.Name = %q, want %q", logger.entry.App.Name, "Test App")
	}

	if logger.entry.Env != "test" {
		t.Errorf("Env = %q, want %q", logger.entry.Env, "test")
	}

	if logger.entry.Lambda.ID != "testapp-api-test" {
		t.Errorf("Lambda.ID = %q, want %q", logger.entry.Lambda.ID, "testapp-api-test")
	}
}

func TestLoggerSetRequest(t *testing.T) {
	mock := &mockClient{}
	logger := NewWithClient(Config{AppID: "test", Env: "test"}, mock)

	logger.SetRequest("POST", "/api/issues", "1700")

	if logger.entry.Request.Method != "POST" {
		t.Errorf("Request.Method = %q, want %q", logger.entry.Request.Method, "POST")
	}
	if logger.entry.Request.Path != "/api/issues" {
		t.Errorf("Request.Path = %q, want %q", logger.entry.Request.Path, "/api/issues")
	}
	if logger.entry.Request.ResourceID != "1700" {
		t.Errorf("Request.ResourceID = %q, want %q", logger.entry.Request.ResourceID, "1700")
	}
}

func TestLoggerSetResponse(t *testing.T) {
	mock := &mockClient{}
	logger := NewWithClient(Config{AppID: "test", Env: "test"}, mock)

	logger.SetResponse(201, "Issue", 1234)

	if logger.entry.Response.Status != 201 {
		t.Errorf("Response.Status = %d, want %d", logger.entry.Response.Status, 201)
	}
	if logger.entry.Response.Type != "Issue" {
		t.Errorf("Response.Type = %q, want %q", logger.entry.Response.Type, "Issue")
	}
	if logger.entry.Response.Bytes != 1234 {
		t.Errorf("Response.Bytes = %d, want %d", logger.entry.Response.Bytes, 1234)
	}
}

func TestLoggerSetSession(t *testing.T) {
	mock := &mockClient{}
	logger := NewWithClient(Config{AppID: "test", Env: "test"}, mock)

	logger.SetSession("sess123", "corr456")

	if logger.entry.SessionID != "sess123" {
		t.Errorf("SessionID = %q, want %q", logger.entry.SessionID, "sess123")
	}
	if logger.entry.CorrelationID != "corr456" {
		t.Errorf("CorrelationID = %q, want %q", logger.entry.CorrelationID, "corr456")
	}
}

func TestLoggerSetTenant(t *testing.T) {
	mock := &mockClient{}
	logger := NewWithClient(Config{AppID: "test", Env: "test"}, mock)

	logger.SetTenant("tenant001")

	if logger.entry.TenantID != "tenant001" {
		t.Errorf("TenantID = %q, want %q", logger.entry.TenantID, "tenant001")
	}
}

func TestLoggerSetUser(t *testing.T) {
	mock := &mockClient{}
	logger := NewWithClient(Config{AppID: "test", Env: "test"}, mock)

	logger.SetUser("sha256:abc123")

	if logger.entry.UserIDHash != "sha256:abc123" {
		t.Errorf("UserIDHash = %q, want %q", logger.entry.UserIDHash, "sha256:abc123")
	}
}

func TestLoggerSetPerf(t *testing.T) {
	mock := &mockClient{}
	logger := NewWithClient(Config{AppID: "test", Env: "test"}, mock)

	logger.SetPerf(true, 100, 50, 25)

	if !logger.entry.Perf.ColdStart {
		t.Error("Perf.ColdStart = false, want true")
	}
	if logger.entry.Perf.InitMs != 100 {
		t.Errorf("Perf.InitMs = %d, want %d", logger.entry.Perf.InitMs, 100)
	}
	if logger.entry.Perf.DBQueryMs != 50 {
		t.Errorf("Perf.DBQueryMs = %d, want %d", logger.entry.Perf.DBQueryMs, 50)
	}
	if logger.entry.Perf.ExternalMs != 25 {
		t.Errorf("Perf.ExternalMs = %d, want %d", logger.entry.Perf.ExternalMs, 25)
	}
}

func TestLoggerSetSource(t *testing.T) {
	mock := &mockClient{}
	logger := NewWithClient(Config{AppID: "test", Env: "test"}, mock)

	logger.SetSource("192.168.1.1", "Mozilla/5.0", "https://example.com")

	if logger.entry.Source.IP != "192.168.1.1" {
		t.Errorf("Source.IP = %q, want %q", logger.entry.Source.IP, "192.168.1.1")
	}
	if logger.entry.Source.UserAgent != "Mozilla/5.0" {
		t.Errorf("Source.UserAgent = %q, want %q", logger.entry.Source.UserAgent, "Mozilla/5.0")
	}
	if logger.entry.Source.Referer != "https://example.com" {
		t.Errorf("Source.Referer = %q, want %q", logger.entry.Source.Referer, "https://example.com")
	}
}

func TestLoggerComplete(t *testing.T) {
	mock := &mockClient{}
	logger := NewWithClient(Config{
		AppID:   "testapp",
		AppName: "Test App",
		Env:     "test",
	}, mock)

	logger.SetRequest("GET", "/api/health", "")
	logger.SetResponse(200, "", 0)

	// Simulate some work
	time.Sleep(10 * time.Millisecond)

	logger.Complete()

	if len(mock.entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(mock.entries))
	}

	entry := mock.entries[0]
	if entry.Response.DurationMs < 10 {
		t.Errorf("DurationMs = %d, expected >= 10", entry.Response.DurationMs)
	}
}

func TestLoggerClone(t *testing.T) {
	mock := &mockClient{}
	original := NewWithClient(Config{
		AppID:   "testapp",
		AppName: "Test App",
		Env:     "test",
	}, mock)

	original.SetTenant("tenant1")

	cloned := original.Clone()

	// Cloned should have fresh entry
	if cloned.entry.TenantID != "" {
		t.Errorf("Clone should have fresh entry, TenantID = %q", cloned.entry.TenantID)
	}

	// Cloned should share the same client
	cloned.SetRequest("POST", "/api/test", "")
	cloned.Complete()

	if len(mock.entries) != 1 {
		t.Fatalf("expected 1 entry from cloned logger, got %d", len(mock.entries))
	}
}

func TestContextIntegration(t *testing.T) {
	mock := &mockClient{}
	logger := NewWithClient(Config{AppID: "test", Env: "test"}, mock)

	ctx := context.Background()

	// Initially no logger
	if FromContext(ctx) != nil {
		t.Error("FromContext should return nil for empty context")
	}

	// Add logger to context
	ctx = WithLogger(ctx, logger)

	// Retrieve logger
	retrieved := FromContext(ctx)
	if retrieved != logger {
		t.Error("FromContext should return the same logger")
	}
}

func TestMustFromContextPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("MustFromContext should panic when no logger in context")
		}
	}()

	ctx := context.Background()
	MustFromContext(ctx)
}

func TestSchemaV1JSON(t *testing.T) {
	mock := &mockClient{}
	logger := NewWithClient(Config{
		AppID:      "dskanban",
		AppName:    "DS Projects",
		Env:        "prod",
		LambdaID:   "dskanban-api-prod",
		LambdaName: "api",
	}, mock)

	logger.SetRequest("POST", "/api/issues", "1700")
	logger.SetResponse(201, "Issue", 1234)
	logger.SetSession("abc123", "req-xyz789")
	logger.SetTenant("t_001")
	logger.SetUser("sha256:abc...")
	logger.Complete()

	if len(mock.entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(mock.entries))
	}

	entry := mock.entries[0]

	// Verify schema version
	if entry.Version != 1 {
		t.Errorf("Version = %d, want 1", entry.Version)
	}

	// Marshal to JSON and verify structure
	data, err := json.Marshal(entry)
	if err != nil {
		t.Fatalf("JSON marshal error: %v", err)
	}

	// Unmarshal to map to check JSON keys
	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatalf("JSON unmarshal error: %v", err)
	}

	// Check required fields exist
	if _, ok := m["v"]; !ok {
		t.Error("JSON missing 'v' field")
	}
	if _, ok := m["ts"]; !ok {
		t.Error("JSON missing 'ts' field")
	}
	if _, ok := m["app"]; !ok {
		t.Error("JSON missing 'app' field")
	}
	if _, ok := m["request"]; !ok {
		t.Error("JSON missing 'request' field")
	}
}

func TestOmitEmpty(t *testing.T) {
	mock := &mockClient{}
	logger := NewWithClient(Config{
		AppID: "test",
		Env:   "test",
		// No Lambda info
	}, mock)

	logger.SetRequest("GET", "/api/health", "")
	// No response, session, tenant, user, perf, source
	logger.Complete()

	entry := mock.entries[0]
	data, _ := json.Marshal(entry)

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	// Lambda should be omitted (empty)
	if _, ok := m["lambda"]; ok {
		t.Error("Empty lambda should be omitted from JSON")
	}

	// session_id should be omitted
	if _, ok := m["session_id"]; ok {
		t.Error("Empty session_id should be omitted from JSON")
	}

	// tenant_id should be omitted
	if _, ok := m["tenant_id"]; ok {
		t.Error("Empty tenant_id should be omitted from JSON")
	}
}

func TestNoopClient(t *testing.T) {
	client := NewNoopClient()

	entry := &LogEntry{Version: 1}
	err := client.WriteLog(entry)

	if err != nil {
		t.Errorf("NoopClient.WriteLog() error = %v", err)
	}
}

func TestLocalFallbackClient(t *testing.T) {
	client := NewLocalFallbackClient()

	entry := &LogEntry{
		Version: 1,
		App:     AppInfo{ID: "test"},
	}

	err := client.WriteLog(entry)
	if err != nil {
		t.Errorf("LocalFallbackClient.WriteLog() error = %v", err)
	}
}

func TestSetEventType(t *testing.T) {
	mock := &mockClient{}
	logger := NewWithClient(Config{AppID: "test", Env: "test"}, mock)

	logger.SetEventType("user.created")

	if logger.entry.Request.EventType != "user.created" {
		t.Errorf("Request.EventType = %q, want %q", logger.entry.Request.EventType, "user.created")
	}
}

func TestSetColdStart(t *testing.T) {
	mock := &mockClient{}
	logger := NewWithClient(Config{AppID: "test", Env: "test"}, mock)

	logger.SetColdStart(true)

	if !logger.entry.Perf.ColdStart {
		t.Error("Perf.ColdStart = false, want true")
	}
}

func TestEntry(t *testing.T) {
	mock := &mockClient{}
	logger := NewWithClient(Config{AppID: "test", Env: "test"}, mock)

	entry := logger.Entry()

	if entry != logger.entry {
		t.Error("Entry() should return the internal entry")
	}

	// Direct modification should work
	entry.TenantID = "direct-set"

	if logger.entry.TenantID != "direct-set" {
		t.Error("Direct modification of Entry() should affect logger")
	}
}

func TestLogEntryIsEmpty(t *testing.T) {
	empty := &LogEntry{}
	if !empty.IsEmpty() {
		t.Error("Empty LogEntry should return IsEmpty() = true")
	}

	withApp := &LogEntry{App: AppInfo{ID: "test"}}
	if withApp.IsEmpty() {
		t.Error("LogEntry with App.ID should return IsEmpty() = false")
	}

	withMethod := &LogEntry{Request: RequestInfo{Method: "GET"}}
	if withMethod.IsEmpty() {
		t.Error("LogEntry with Request.Method should return IsEmpty() = false")
	}

	withEvent := &LogEntry{Request: RequestInfo{EventType: "test.event"}}
	if withEvent.IsEmpty() {
		t.Error("LogEntry with Request.EventType should return IsEmpty() = false")
	}
}
