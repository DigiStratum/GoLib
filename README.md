# GoLib
Go Library code, generally reusable bits

## TODO
* Convert all method receivers from pointer to non-pointer (value based) for methods that do not modify the value (?)
* Test coverage for all library packages; test runner that covers all packages
* Documentation/README, SDK site
* Examples
* Thread-safe mutexes, semaphores, etc. for concurrency everywhere
* Standardized error handling; use fmt.Errorf("%d", num) instead of errors.New(fmt.Sprintf("%d", num))