// Package TestRunner provides utilities for running Go tests across a project.
//
// TestRunner offers functionality to find test packages, run vetting and tests,
// and collect test results in a structured way.
package TestRunner

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// Options configures how the TestRunner behaves
type Options struct {
	// Verbose enables detailed output
	Verbose bool

	// SkipVet skips running 'go vet' on packages
	SkipVet bool

	// FailFast stops on first test failure
	FailFast bool

	// TestPattern is a regex to filter tests
	TestPattern string

	// Parallel runs tests in parallel when possible
	Parallel bool

	// Count sets the -count flag for go test
	Count int

	// CoverMode enables code coverage with the specified mode
	CoverMode string

	// IgnoreDirs specifies directories to skip
	IgnoreDirs []string
}

// DefaultOptions returns the default options for the TestRunner
func DefaultOptions() *Options {
	return &Options{
		Verbose:    true,
		SkipVet:    false,
		FailFast:   false,
		Count:      1,
		IgnoreDirs: []string{"vendor", ".git"},
	}
}

// TestResult represents the result of a test run
type TestResult struct {
	Package string
	Success bool
	Output  string
	Error   error
	Time    time.Duration
}

// TestRunnerIfc defines the interface for running tests across a Go project
type TestRunnerIfc interface {
	// FindGoDirectories finds all directories containing Go files
	FindGoDirectories(rootDir string) ([]string, error)

	// FindTestPackages finds all directories containing test files
	FindTestPackages(rootDir string) ([]string, error)

	// RunVet runs 'go vet' on each directory
	RunVet(dirs []string) error

	// RunTests runs 'go test' on each directory
	RunTests(dirs []string) error

	// RunTestPackages runs go test on each package
	RunTestPackages(packages []string) (bool, error)

	// RunAll runs the entire test suite
	RunAll(rootDir string) error

	// PrintSummary prints a summary of test results
	PrintSummary(totalTime time.Duration)

	// HasFailures returns true if any tests failed
	HasFailures() bool

	// GetOptions returns the current options
	GetOptions() *Options
}

// TestRunner runs tests across a Go project
type TestRunner struct {
	Options     *Options
	VetResults  []TestResult
	TestResults []TestResult
	mu          sync.Mutex
}

// NewTestRunner creates a new TestRunner with the given options
func NewTestRunner(options *Options) *TestRunner {
	if options == nil {
		options = DefaultOptions()
	}
	return &TestRunner{
		Options: options,
	}
}

// GetOptions returns the current options
func (r *TestRunner) GetOptions() *Options {
	if r == nil {
		return nil
	}
	return r.Options
}

// FindGoDirectories finds all directories containing Go files
func (r *TestRunner) FindGoDirectories(rootDir string) ([]string, error) {
	if r == nil {
		return nil, fmt.Errorf("nil receiver")
	}

	var dirs []string
	dirMap := make(map[string]bool)

	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip ignored directories
		if info.IsDir() {
			base := filepath.Base(path)

			// Skip dirs that start with _ (Go build convention)
			if strings.HasPrefix(base, "_") {
				return filepath.SkipDir
			}

			// Skip explicitly ignored dirs
			for _, ignoreDir := range r.Options.IgnoreDirs {
				if base == ignoreDir {
					return filepath.SkipDir
				}
			}

			return nil
		}

		// Process Go files
		if strings.HasSuffix(path, ".go") {
			dir := filepath.Dir(path)
			if !dirMap[dir] {
				dirMap[dir] = true
				dirs = append(dirs, dir)
			}
		}
		return nil
	})

	return dirs, err
}

// FindTestPackages finds all directories containing test files
func (r *TestRunner) FindTestPackages(rootDir string) ([]string, error) {
	if r == nil {
		return nil, fmt.Errorf("nil receiver")
	}

	var packages []string
	packageMap := make(map[string]bool)

	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip ignored directories
		if info.IsDir() {
			base := filepath.Base(path)

			// Skip dirs that start with _ (Go build convention)
			if strings.HasPrefix(base, "_") {
				return filepath.SkipDir
			}

			// Skip explicitly ignored dirs
			for _, ignoreDir := range r.Options.IgnoreDirs {
				if base == ignoreDir {
					return filepath.SkipDir
				}
			}

			return nil
		}

		// Look for test files
		if strings.HasSuffix(path, "_test.go") {
			dir := filepath.Dir(path)
			if !packageMap[dir] {
				packageMap[dir] = true
				packages = append(packages, dir)
			}
		}

		return nil
	})

	return packages, err
}

// RunVet runs 'go vet' on each directory
func (r *TestRunner) RunVet(dirs []string) error {
	if r == nil {
		return fmt.Errorf("nil receiver")
	}

	if r.Options.SkipVet {
		if r.Options.Verbose {
			fmt.Println("Skipping go vet")
		}
		return nil
	}

	var vetFailed bool
	for _, dir := range dirs {
		if r.Options.Verbose {
			fmt.Printf("Vetting: %s\n", dir)
		}

		cmd := exec.Command("go", "vet", "./...")
		cmd.Dir = dir
		output, err := cmd.CombinedOutput()

		result := TestResult{
			Package: dir,
			Success: err == nil,
			Output:  string(output),
			Error:   err,
		}

		r.mu.Lock()
		r.VetResults = append(r.VetResults, result)
		r.mu.Unlock()

		if err != nil {
			vetFailed = true
			fmt.Printf("Go vet failed in %s:\n%s\n", dir, output)
			if r.Options.FailFast {
				return fmt.Errorf("go vet failed in %s", dir)
			}
		}
	}

	if vetFailed {
		return fmt.Errorf("go vet failed in one or more packages")
	}
	return nil
}

// RunTests runs 'go test' on each directory
func (r *TestRunner) RunTests(dirs []string) error {
	if r == nil {
		return fmt.Errorf("nil receiver")
	}

	var wg sync.WaitGroup
	var testFailed bool
	errChan := make(chan error, len(dirs))

	// Process directories sequentially or in parallel
	processDir := func(dir string) {
		defer wg.Done()

		if r.Options.Verbose {
			fmt.Printf("Testing: %s\n", dir)
		}

		args := []string{"test", "-v"}

		// Add count flag
		args = append(args, "-count", fmt.Sprintf("%d", r.Options.Count))

		// Add test pattern if provided
		if r.Options.TestPattern != "" {
			args = append(args, "-run", r.Options.TestPattern)
		}

		// Add coverage if desired
		if r.Options.CoverMode != "" {
			coverprofile := filepath.Join(dir, "coverage.out")
			args = append(args, "-covermode", r.Options.CoverMode, "-coverprofile", coverprofile)
		}

		// Add package specifier
		args = append(args, "./...")

		start := time.Now()
		cmd := exec.Command("go", args...)
		cmd.Dir = dir
		output, err := cmd.CombinedOutput()
		elapsed := time.Since(start)

		result := TestResult{
			Package: dir,
			Success: err == nil,
			Output:  string(output),
			Error:   err,
			Time:    elapsed,
		}

		r.mu.Lock()
		r.TestResults = append(r.TestResults, result)
		r.mu.Unlock()

		if r.Options.Verbose {
			fmt.Println(string(output))
		}

		if err != nil {
			errChan <- fmt.Errorf("tests failed in %s", dir)
		}
	}

	if r.Options.Parallel {
		for _, dir := range dirs {
			wg.Add(1)
			go processDir(dir)
		}
	} else {
		for _, dir := range dirs {
			wg.Add(1)
			processDir(dir)

			// If FailFast is true and we got an error, stop processing
			select {
			case err := <-errChan:
				if r.Options.FailFast {
					return err
				}
				testFailed = true
			default:
				// Continue to next directory
			}
		}
	}

	wg.Wait()
	close(errChan)

	// Check for errors
	for err := range errChan {
		testFailed = true
		if r.Options.FailFast {
			return err
		}
	}

	if testFailed {
		return fmt.Errorf("tests failed in one or more packages")
	}
	return nil
}

// RunTestPackages runs go test on each package
func (r *TestRunner) RunTestPackages(packages []string) (bool, error) {
	if r == nil {
		return false, fmt.Errorf("nil receiver")
	}

	startTime := time.Now()
	fmt.Println("=== Starting Test Run ===")

	allPassed := true
	var testCount, passCount int

	for _, pkg := range packages {
		testCount++

		if r.Options.Verbose {
			fmt.Printf("Testing package: %s\n", pkg)
		}

		// Run go vet if not skipped
		if !r.Options.SkipVet {
			if r.Options.Verbose {
				fmt.Printf("Vetting: %s\n", pkg)
			}

			cmd := exec.Command("go", "vet", "./...")
			cmd.Dir = pkg
			vetOutput, vetErr := cmd.CombinedOutput()

			if vetErr != nil {
				fmt.Printf("Go vet failed in %s:\n%s\n", pkg, vetOutput)
				allPassed = false

				if r.Options.FailFast {
					break
				}

				continue // Skip tests for this package
			}
		}

		// Build test command
		args := []string{"test", "-v"}

		// Add count flag if specified
		if r.Options.Count > 0 {
			args = append(args, "-count", fmt.Sprintf("%d", r.Options.Count))
		}

		// Add test pattern if provided
		if r.Options.TestPattern != "" {
			args = append(args, "-run", r.Options.TestPattern)
		}

		// Add coverage if desired
		if r.Options.CoverMode != "" {
			coverprofile := filepath.Join(pkg, "coverage.out")
			args = append(args, "-covermode", r.Options.CoverMode, "-coverprofile", coverprofile)
		}

		// Add package specifier
		args = append(args, "./...")

		// Run tests
		cmd := exec.Command("go", args...)
		cmd.Dir = pkg
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		err := cmd.Run()
		if err != nil {
			fmt.Printf("Tests failed in %s: %v\n", pkg, err)
			allPassed = false

			if r.Options.FailFast {
				break
			}
		} else {
			passCount++
		}
	}

	// Print summary
	duration := time.Since(startTime)
	fmt.Println("\n=== Test Summary ===")
	fmt.Printf("Total Time: %s\n", duration)
	fmt.Printf("Packages: %d/%d passed\n", passCount, testCount)

	if allPassed {
		fmt.Println("All tests passed!")
	} else {
		fmt.Println("Some tests failed.")
	}

	return allPassed, nil
}

// RunAll runs the entire test suite
func (r *TestRunner) RunAll(rootDir string) error {
	if r == nil {
		return fmt.Errorf("nil receiver")
	}

	fmt.Println("=== Starting Test Run ===")
	startTime := time.Now()

	// Find all Go directories
	dirs, err := r.FindGoDirectories(rootDir)
	if err != nil {
		return fmt.Errorf("failed to find Go directories: %w", err)
	}

	// Run go vet
	vetErr := r.RunVet(dirs)

	// Run tests regardless of vet errors unless FailFast is true
	var testErr error
	if vetErr == nil || !r.Options.FailFast {
		testErr = r.RunTests(dirs)
	}

	// Print summary
	r.PrintSummary(time.Since(startTime))

	// Return first error
	if vetErr != nil {
		return vetErr
	}
	return testErr
}

// PrintSummary prints a summary of test results
func (r *TestRunner) PrintSummary(totalTime time.Duration) {
	if r == nil {
		return
	}

	fmt.Println("\n=== Test Summary ===")
	fmt.Printf("Total Time: %s\n", totalTime)

	// Vet results
	vetPassed := 0
	for _, result := range r.VetResults {
		if result.Success {
			vetPassed++
		}
	}
	fmt.Printf("Go Vet: %d/%d passed\n", vetPassed, len(r.VetResults))

	// Test results
	testPassed := 0
	for _, result := range r.TestResults {
		if result.Success {
			testPassed++
		}
	}
	fmt.Printf("Go Tests: %d/%d passed\n", testPassed, len(r.TestResults))

	// Failed tests
	if vetPassed < len(r.VetResults) || testPassed < len(r.TestResults) {
		fmt.Println("\n=== Failed Tests ===")

		for _, result := range r.VetResults {
			if !result.Success {
				fmt.Printf("Vet Failed: %s\n", result.Package)
			}
		}

		for _, result := range r.TestResults {
			if !result.Success {
				fmt.Printf("Test Failed: %s (%s)\n", result.Package, result.Time)
			}
		}
	}
}

// HasFailures returns true if any tests failed
func (r *TestRunner) HasFailures() bool {
	if r == nil {
		return false
	}

	for _, result := range r.VetResults {
		if !result.Success {
			return true
		}
	}

	for _, result := range r.TestResults {
		if !result.Success {
			return true
		}
	}

	return false
}
