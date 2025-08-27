// filepath: /Users/skelly/Documents/GoProjects/GoLib/TestRunner/testrunner.go
package TestRunner

import (
	"errors"
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

	// RunLegacyT enables running legacy 't' shell scripts
	RunLegacyT bool

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
		RunLegacyT: true,
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

// TestRunner runs tests across a Go project
type TestRunner struct {
	Options       *Options
	VetResults    []TestResult
	TestResults   []TestResult
	LegacyResults []TestResult
	mu            sync.Mutex
}

// New creates a new TestRunner with the given options
func New(options *Options) *TestRunner {
	if options == nil {
		options = DefaultOptions()
	}
	return &TestRunner{
		Options: options,
	}
}

// FindGoDirectories finds all directories containing Go files
func (tr *TestRunner) FindGoDirectories(rootDir string) ([]string, error) {
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
			for _, ignoreDir := range tr.Options.IgnoreDirs {
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

// RunVet runs 'go vet' on each directory
func (tr *TestRunner) RunVet(dirs []string) error {
	if tr.Options.SkipVet {
		if tr.Options.Verbose {
			fmt.Println("Skipping go vet")
		}
		return nil
	}

	var vetFailed bool
	for _, dir := range dirs {
		if tr.Options.Verbose {
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

		tr.mu.Lock()
		tr.VetResults = append(tr.VetResults, result)
		tr.mu.Unlock()

		if err != nil {
			vetFailed = true
			fmt.Printf("Go vet failed in %s:\n%s\n", dir, output)
			if tr.Options.FailFast {
				return fmt.Errorf("go vet failed in %s", dir)
			}
		}
	}

	if vetFailed {
		return errors.New("go vet failed in one or more packages")
	}
	return nil
}

// RunTests runs 'go test' on each directory
func (tr *TestRunner) RunTests(dirs []string) error {
	var wg sync.WaitGroup
	var testFailed bool
	errChan := make(chan error, len(dirs))

	// Process directories sequentially or in parallel
	processDir := func(dir string) {
		defer wg.Done()

		if tr.Options.Verbose {
			fmt.Printf("Testing: %s\n", dir)
		}

		args := []string{"test", "-v"}

		// Add count flag
		args = append(args, "-count", fmt.Sprintf("%d", tr.Options.Count))

		// Add test pattern if provided
		if tr.Options.TestPattern != "" {
			args = append(args, "-run", tr.Options.TestPattern)
		}

		// Add coverage if desired
		if tr.Options.CoverMode != "" {
			coverprofile := filepath.Join(dir, "coverage.out")
			args = append(args, "-covermode", tr.Options.CoverMode, "-coverprofile", coverprofile)
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

		tr.mu.Lock()
		tr.TestResults = append(tr.TestResults, result)
		tr.mu.Unlock()

		if tr.Options.Verbose {
			fmt.Println(string(output))
		}

		if err != nil {
			errChan <- fmt.Errorf("tests failed in %s", dir)
		}
	}

	if tr.Options.Parallel {
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
				if tr.Options.FailFast {
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
		if tr.Options.FailFast {
			return err
		}
	}

	if testFailed {
		return errors.New("tests failed in one or more packages")
	}
	return nil
}

// FindLegacyTScripts finds all 't' shell scripts
func (tr *TestRunner) FindLegacyTScripts(rootDir string) ([]string, error) {
	var tScripts []string

	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip ignored directories
		if info.IsDir() {
			base := filepath.Base(path)
			for _, ignoreDir := range tr.Options.IgnoreDirs {
				if base == ignoreDir {
					return filepath.SkipDir
				}
			}
			return nil
		}

		if filepath.Base(path) == "t" {
			tScripts = append(tScripts, path)
		}
		return nil
	})

	return tScripts, err
}

// RunLegacyTScripts runs all 't' shell scripts
func (tr *TestRunner) RunLegacyTScripts(rootDir string) error {
	if !tr.Options.RunLegacyT {
		if tr.Options.Verbose {
			fmt.Println("Skipping legacy 't' scripts")
		}
		return nil
	}

	scripts, err := tr.FindLegacyTScripts(rootDir)
	if err != nil {
		return fmt.Errorf("failed to find legacy 't' scripts: %w", err)
	}

	var legacyFailed bool
	for _, script := range scripts {
		if tr.Options.Verbose {
			fmt.Printf("Running legacy script: %s\n", script)
		}

		start := time.Now()
		cmd := exec.Command("/bin/bash", script)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		elapsed := time.Since(start)

		result := TestResult{
			Package: script,
			Success: err == nil,
			Error:   err,
			Time:    elapsed,
		}

		tr.mu.Lock()
		tr.LegacyResults = append(tr.LegacyResults, result)
		tr.mu.Unlock()

		if err != nil {
			legacyFailed = true
			fmt.Printf("Legacy script %s failed\n", script)
			if tr.Options.FailFast {
				return fmt.Errorf("legacy script %s failed", script)
			}
		}
	}

	if legacyFailed {
		return errors.New("one or more legacy 't' scripts failed")
	}
	return nil
}

// RunAll runs the entire test suite
func (tr *TestRunner) RunAll(rootDir string) error {
	fmt.Println("=== Starting Test Run ===")
	startTime := time.Now()

	// Find all Go directories
	dirs, err := tr.FindGoDirectories(rootDir)
	if err != nil {
		return fmt.Errorf("failed to find Go directories: %w", err)
	}

	// Run go vet
	vetErr := tr.RunVet(dirs)

	// Run tests regardless of vet errors unless FailFast is true
	var testErr error
	if vetErr == nil || !tr.Options.FailFast {
		testErr = tr.RunTests(dirs)
	}

	// Run legacy scripts regardless of previous errors unless FailFast is true
	var legacyErr error
	if (vetErr == nil && testErr == nil) || !tr.Options.FailFast {
		legacyErr = tr.RunLegacyTScripts(rootDir)
	}

	// Print summary
	tr.PrintSummary(time.Since(startTime))

	// Return first error
	if vetErr != nil {
		return vetErr
	}
	if testErr != nil {
		return testErr
	}
	return legacyErr
}

// PrintSummary prints a summary of test results
func (tr *TestRunner) PrintSummary(totalTime time.Duration) {
	fmt.Println("\n=== Test Summary ===")
	fmt.Printf("Total Time: %s\n", totalTime)

	// Vet results
	vetPassed := 0
	for _, result := range tr.VetResults {
		if result.Success {
			vetPassed++
		}
	}
	fmt.Printf("Go Vet: %d/%d passed\n", vetPassed, len(tr.VetResults))

	// Test results
	testPassed := 0
	for _, result := range tr.TestResults {
		if result.Success {
			testPassed++
		}
	}
	fmt.Printf("Go Tests: %d/%d passed\n", testPassed, len(tr.TestResults))

	// Legacy results
	legacyPassed := 0
	for _, result := range tr.LegacyResults {
		if result.Success {
			legacyPassed++
		}
	}
	fmt.Printf("Legacy Scripts: %d/%d passed\n", legacyPassed, len(tr.LegacyResults))

	// Failed tests
	if vetPassed < len(tr.VetResults) || testPassed < len(tr.TestResults) || legacyPassed < len(tr.LegacyResults) {
		fmt.Println("\n=== Failed Tests ===")

		for _, result := range tr.VetResults {
			if !result.Success {
				fmt.Printf("Vet Failed: %s\n", result.Package)
			}
		}

		for _, result := range tr.TestResults {
			if !result.Success {
				fmt.Printf("Test Failed: %s (%s)\n", result.Package, result.Time)
			}
		}

		for _, result := range tr.LegacyResults {
			if !result.Success {
				fmt.Printf("Legacy Script Failed: %s\n", result.Package)
			}
		}
	}
}

// HasFailures returns true if any tests failed
func (tr *TestRunner) HasFailures() bool {
	for _, result := range tr.VetResults {
		if !result.Success {
			return true
		}
	}

	for _, result := range tr.TestResults {
		if !result.Success {
			return true
		}
	}

	for _, result := range tr.LegacyResults {
		if !result.Success {
			return true
		}
	}

	return false
}
