# GitHub Copilot Instructions for GoLib Project

This document provides guidelines for code generation and assistance in the GoLib project. It covers Go best practices and project-specific conventions for this reusable library collection.

## Project Structure

GoLib is a collection of reusable Go packages that provide common functionality across various domains:
- Cache management
- Time and chronology utilities
- Configuration handling
- Data structures and type management
- Database connections
- Dependency injection
- Event handling
- File I/O operations
- Logging
- Networking
- Object models
- Process management
- Testing utilities

## Custom Conventions

### Package Documentation
- **Every package must include a README.md file**
- **README.md should cover:**
  - **A clear summary of the package's purpose**
  - **Design principles and architectural decisions**
  - **Key abstractions and interfaces**
  - **Integration points with other packages**
  - **Usage examples for the primary functions**
```markdown
<!-- Example README.md structure -->
# Package Name

## Purpose
Brief description of what this package does and its role within the larger system.

## Design
Overview of the design principles, patterns used, and architectural decisions.

## Key Components
* ComponentA - Description of what it does and why it exists
* ComponentB - Description of what it does and why it exists

## Integration
How this package integrates with other parts of the system.

## Usage Examples
```go
// Simple example showing primary usage
```
```

### Receiver Functions
- **Always use the identifier "r" as the receiver for functions attached to a type definition or structure**
```go
// CORRECT
func (r *MyType) MyMethod() {}

// INCORRECT
func (m *MyType) MyMethod() {}
```
- **Use pointer receiver for mutable (write) operations, copy for immutable (read) operations**

### Interface and Implementation Pattern
- **Accept interfaces, return structs** (except by exception)
- **ONE exported struct+interface per source file** will make the code easier to read (with exceptions)
- **Implement io.Closer interface** for any class that opens/manages precious and/or external resources

### Interface Naming
- **Use "Ifc" suffix for interfaces, not "Interface"**
- **Always export an interface that declares receiver function prototypes**
```go
// CORRECT
type DocumentStoreIfc interface {
    Create(ctx context.Context, document *Document) error
}

// INCORRECT
type DocumentStoreInterface interface {}
```

### Factory Functions
- **Always use the prefix "New" for factory functions, not "Create"**
```go
// CORRECT
func NewDocumentStore(config map[string]interface{}) (DocumentStoreIfc, error) {}

// INCORRECT
func CreateDocumentStore(config map[string]interface{}) (DocumentStoreIfc, error) {}
```

### Data Encapsulation
- **Do not export data types unless necessary**
- **Use the interface for integration points instead of concrete types**
- **Generate getter/setter/accessor functions for non-exported properties**
```go
// CORRECT
type documentStore struct {
    connectionPool *connectionPool
}

func (r *documentStore) GetConnectionPool() ConnectionPoolIfc {
    return r.connectionPool
}

// INCORRECT
type DocumentStore struct {
    ConnectionPool *ConnectionPool // directly exposed
}
```

### Nil Handling
- **Handle r == nil receivers** in all methods
```go
func (r *MyType) MyMethod() error {
    if r == nil {
        return errors.New("nil receiver")
    }
    // method implementation
}
```

### Concurrency
- **Prefer go-routine+channel for concurrency orchestration, over mutex**
- **Use mutex lock, semaphore, channels, etc. for mutable (write) operations**

### Error Handling
- **Use fmt.Errorf() instead of errors.New(fmt.Sprintf())**
- **Don't produce error log output from library functions** where it can be left to the consumer
```go
// CORRECT
func (r *Cache) Get(key string) (interface{}, error) {
    if value, ok := r.items[key]; ok {
        return value, nil
    }
    return nil, fmt.Errorf("key %s not found", key)
}

// INCORRECT
func (r *Cache) Get(key string) interface{} {
    if value, ok := r.items[key]; ok {
        return value
    }
    log.Printf("Error: key %s not found", key)
    return nil
}
```

### Logging
- **Log Trace() messages** to track entry into library functions with calling arguments as appropriate

### Testing
- **Generate unit test coverage for all functions**
- **Test both exported and non-exported functions and properties**
- **Each test should verify a single aspect of functionality**
- **Use `-count=1` flag to bypass test execution cache**
- **Generate coverage reports with `go tool cover -html=.test_coverage.txt`**

## Go Best Practices

### Error Handling
- Return errors rather than using panic
- Check errors immediately after function calls that can produce them
- Wrap errors with context when propagating up the call stack
- Use sentinel errors for expected error conditions that callers might want to check for

### Context Usage
- Pass context.Context as the first parameter for functions that:
  - Make network calls
  - Access external systems
  - Perform long-running operations
- Respect context cancellation in all operations
- Don't store context in structs

### Concurrency
- Use channels for communication, not just synchronization
- Prefer sync.Mutex over channels for simple state protection
- Always close channels when no more values will be sent
- Handle potential panics in goroutines with recover()

### Documentation
- Document all exported functions, types, and constants
- Include usage examples for complex or commonly used functionality
- Document non-obvious behaviors or edge cases

### Testing
- Use table-driven tests for functions with multiple test cases
- Avoid using global state in tests
- Mock external dependencies for unit tests