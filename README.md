# GoLib
Go Library code, generally reusable bits

## TODO
* Documentation/README, SDK site
* Examples
* Standardized error handling; use fmt.Errorf("%d", num) instead of errors.New(fmt.Sprintf("%d", num))

## REFACTORING:
** Accept Interfaces, return structs (except by exception)
** Use pointer receiver for mutable (write) operations, copy for immutable (read) operations
** Use mutex lock, semaphore, channels, etc for mutable (write) operations (prefer go-routine+channel for concurrency orchestration, over mutex)
** Use r for receiver everywhere for better copy/paste; only use (*r) when necessary (actually 'this' might have been syntactically nice...)
** Clean up TODO's / FIXME's as reasonably able
** Add test overage as reasonably able
** Add documentation (godoc, readme) as reasonably able
** Add working examples as reasonably able
** Use fmt.Errorf() instead of errors.New()
** Don't produce error log output from library functions where it can be left to the consumer
** ONE exported struct+interface per source file will make the code easier to read (with exceptions)
** Log Trace() messages to track entry into library functions with calling arguments as appropriate
** Implement io.Closer interface for any class that opens/managers precious and/or external resources which should be closed/freed upon release/shutdown, etc.
*** ref: https://pkg.go.dev/io#Closer
*** ref: https://stackoverflow.com/questions/32768243/go-destructors/32781054
** Handle r == nil receivers
*** ref: https://tour.golang.org/methods/12

## TEST
go test -count=1 cache_test.go cache.go cacheitem.go
^^ the -count=1 bypasses the test execution cache to force the tests to run each time.


