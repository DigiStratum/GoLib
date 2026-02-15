# GoLib Test Plan

This document describes the testing approach, coverage requirements, and test organization for GoLib.

## Test Overview

GoLib uses Go's standard testing framework with the following characteristics:

- **Test Files**: 68 test files across the codebase
- **Source Files**: 199 Go source files
- **Test Ratio**: ~34% (test files to source files)

## Testing Philosophy

1. **Unit Tests First** - Each package should have comprehensive unit tests
2. **Interface Mocking** - Use interfaces to enable testability
3. **Test Isolation** - Tests should not depend on external resources
4. **Deterministic Results** - Use injectable time sources for time-dependent code

## Test Organization

Tests are co-located with source files following Go conventions:

```
Package/
├── component.go
├── component_test.go
├── helper.go
└── helper_test.go
```

## Package Test Coverage

### Cache (`Cache/`)

| Test File | Coverage Area |
|-----------|---------------|
| `cache_test.go` | Core cache operations, LRU eviction, expiration |

**Test Cases:**
- Set/Get operations
- Key existence checks
- Expiration behavior
- Size limit enforcement
- Count limit enforcement
- LRU eviction order
- Concurrent access safety
- Flush/Close operations

### Logger (`Logger/`)

| Test File | Coverage Area |
|-----------|---------------|
| `logger_test.go` | Logging at all levels, formatting |
| `loglevel_test.go` | Log level parsing and comparison |
| `logwriter/stdoutlogwriter_test.go` | StdOut writer output |

**Test Cases:**
- Log level filtering
- Message formatting with arguments
- Timestamp inclusion/exclusion
- Prefixed logger creation
- Singleton vs new instance behavior
- Error return for WARN+ levels

### Config (`Config/`)

| Test File | Coverage Area |
|-----------|---------------|
| (via Data/hashmap tests) | HashMap operations |

**Test Cases:**
- Key/value operations
- Subset extraction by prefix
- Config merging
- Reference dereferencing (`%key%` substitution)
- Recursive reference resolution
- Circular reference detection

### HTTP Client (`Net/http/`)

| Test File | Coverage Area |
|-----------|---------------|
| `httpclient_test.go` | HTTP client send operations |
| `httprequest_test.go` | Request construction |
| `httprequestbodybuilder_test.go` | Body building |
| `httpheaders_test.go` | Header manipulation |
| `httpheadersbuilder_test.go` | Header building (if exists) |
| `httpresponse_test.go` | Response parsing |
| `httpstatus_test.go` | Status code handling |
| `requestcontext_test.go` | Request context |

**Test Cases:**
- GET/POST/PUT/DELETE methods
- Header setting and retrieval
- Request body encoding
- Response status code parsing
- Response body reading
- Request builder fluent API

### Database (`DB/MySQL/`)

| Test File | Coverage Area |
|-----------|---------------|
| `connection_test.go` | Connection lifecycle |
| `connectionpool_test.go` | Pool management |
| `connectioncommon_test.go` | Shared test utilities |
| `leasedconnection_test.go` | Lease pattern |
| `leasedconnections_test.go` | Multiple leases |
| `pooledconnection_test.go` | Pooled connection wrapper |
| `mysqlconnectionfactory_test.go` | Factory pattern |
| `query_test.go` | Query construction and execution |
| `result_test.go` | Result handling |
| `resultrow_test.go` | Row data access |
| `resultset_test.go` | Result set iteration |
| `sqlquery_test.go` | SQL query building |

**Test Cases:**
- Connection establishment
- Connection pooling
- Lease acquire/release
- Query parameter binding
- Result set iteration
- NULL value handling
- Connection timeout
- Pool exhaustion behavior

**Mocking:**
Uses `github.com/DATA-DOG/go-sqlmock` for database mocking.

### Chrono (`Chrono/`)

**Test Cases:**
- Timestamp creation
- Time arithmetic (Add)
- Expiration checking
- Unix timestamp conversion
- Mock time source injection

### Object Store (`Object/store/`)

**Test Cases per Store Type:**

**S3 Store:**
- Get object
- Put object
- Delete object
- List objects
- Handle missing objects
- Handle permission errors

**MySQL Store:**
- CRUD operations
- Path parsing
- Error handling

**DynamoDB Store:**
- Item get/put
- Key construction
- Attribute mapping

**Secrets Store:**
- Secret retrieval
- Secret caching

### Testing Utilities (`Testing/`)

The `Testing/` package provides test helpers:

- **mocks/**: Mock implementations for interfaces
- **TestRunner/**: Test harness utilities

## Running Tests

### All Tests

```bash
cd GoLib
go test ./...
```

### Single Package

```bash
go test ./Cache/...
go test ./Logger/...
go test ./Net/http/...
```

### With Coverage

```bash
# Generate coverage report
go test -coverprofile=coverage.txt ./...

# View coverage by package
go tool cover -func=coverage.txt

# Generate HTML report
go tool cover -html=coverage.txt -o coverage.html
```

### Skip Cache

Use `-count=1` to bypass Go's test result caching:

```bash
go test -count=1 ./Cache/...
```

### Verbose Output

```bash
go test -v ./...
```

## Coverage Targets

| Package | Target | Notes |
|---------|--------|-------|
| Cache | 80%+ | Core functionality |
| Logger | 80%+ | Core functionality |
| Config | 70%+ | Critical dereferencing logic |
| Net/http | 80%+ | HTTP client/server |
| DB/MySQL | 80%+ | Database operations |
| Object/store | 70%+ | Integration-heavy |
| Chrono | 90%+ | Simple utilities |
| Data | 80%+ | Data structures |

## Test Patterns

### Table-Driven Tests

```go
func TestCache_Get(t *testing.T) {
    tests := []struct {
        name     string
        key      string
        value    interface{}
        expected interface{}
    }{
        {"string value", "key1", "value1", "value1"},
        {"int value", "key2", 42, 42},
        {"nil on miss", "missing", nil, nil},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            c := cache.NewCache()
            if tt.value != nil {
                c.Set(tt.key, tt.value)
            }
            result := c.Get(tt.key)
            if result != tt.expected {
                t.Errorf("got %v, want %v", result, tt.expected)
            }
        })
    }
}
```

### Mock Time Source

```go
type MockTimeSource struct {
    currentTime time.Time
}

func (m *MockTimeSource) Now() chrono.TimeStampIfc {
    return chrono.NewTimeStampFromTime(m.currentTime)
}

func (m *MockTimeSource) Advance(d time.Duration) {
    m.currentTime = m.currentTime.Add(d)
}

func TestCache_Expiration(t *testing.T) {
    mockTime := &MockTimeSource{currentTime: time.Now()}
    c := cache.NewCache()
    c.SetTimeSource(mockTime)
    
    c.Set("key", "value")
    c.SetExpires("key", mockTime.Now().Add(60))
    
    // Advance past expiration
    mockTime.Advance(61 * time.Second)
    
    if c.Get("key") != nil {
        t.Error("expected expired item to return nil")
    }
}
```

### Database Mocking

```go
func TestQuery_Execute(t *testing.T) {
    db, mock, err := sqlmock.New()
    if err != nil {
        t.Fatalf("failed to create mock: %v", err)
    }
    defer db.Close()
    
    rows := sqlmock.NewRows([]string{"id", "name"}).
        AddRow(1, "Alice").
        AddRow(2, "Bob")
    
    mock.ExpectQuery("SELECT").WillReturnRows(rows)
    
    // Test query execution...
}
```

## Continuous Integration

Tests should run on every commit:

```yaml
# .github/workflows/test.yml
name: Tests
on: [push, pull_request]
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v4
        with:
          go-version: '1.24'
      - run: go test -race -coverprofile=coverage.txt ./...
      - uses: codecov/codecov-action@v3
```

## Known Test Gaps

Areas needing additional coverage:

1. **Object Store backends** - Integration tests with mocked AWS services
2. **WebUI templating** - Magic tag substitution edge cases
3. **Dependencies** - DI framework usage patterns
4. **Process** - Runnable lifecycle management
5. **Starter** - Application startup/shutdown sequences

## Test Maintenance

- Review test coverage quarterly
- Add tests for bug fixes (regression tests)
- Update tests when interfaces change
- Remove obsolete tests for deprecated code
