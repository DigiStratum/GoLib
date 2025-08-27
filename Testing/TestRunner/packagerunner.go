// filepath: /Users/skelly/Documents/GoProjects/GoLib/TestRunner/packagerunner.go
package TestRunner

import (
	"fmt"
	"os"
	"path/filepath"
)

// PackageTestRunner provides package-specific test configuration
type PackageTestRunner struct {
	// Package directory
	Dir string
	// Package name
	Name string
	// Custom test setup
	Setup func() error
	// Custom test teardown
	Teardown func() error
	// Package-specific test options
	Options *Options
}

// DefaultPackageRunner creates a default package test runner configuration
func DefaultPackageRunner() *PackageTestRunner {
	// Get the current directory
	dir, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting current directory: %v\n", err)
		os.Exit(1)
	}

	// Extract package name from directory
	name := filepath.Base(dir)

	// Return default configuration
	return &PackageTestRunner{
		Dir:  dir,
		Name: name,
		// Default options - customize as needed
		Options: &Options{
			Verbose:  true,
			SkipVet:  false,
			FailFast: false,
			Count:    1,
			// Ignore directories that shouldn't be tested
			IgnoreDirs: []string{"vendor", ".git", "testdata"},
		},
	}
}

// RunPackageTests runs tests for a specific package
func RunPackageTests(runner *PackageTestRunner, verbose, debug bool, runShellScript bool, coverMode, testPattern string, isParent bool) error {
	// Get absolute path to package directory
	packageDir, err := filepath.Abs(runner.Dir)
	if err != nil {
		return fmt.Errorf("error getting absolute path: %w", err)
	}

	// Configure test runner options
	options := runner.Options
	if options == nil {
		options = DefaultOptions()
	}

	// Apply command-line flags
	options.Verbose = verbose
	options.TestPattern = testPattern
	options.CoverMode = coverMode

	// Debug info
	if debug {
		fmt.Printf("Package test runner for: %s\n", runner.Name)
		fmt.Printf("Directory: %s\n", packageDir)
		fmt.Printf("Options: %+v\n", options)
	}

	// Run setup if defined
	if runner.Setup != nil {
		if err := runner.Setup(); err != nil {
			return fmt.Errorf("setup failed: %w", err)
		}
	}

	// Create and run the test runner
	testRunner := NewTestRunner(options)
	err = testRunner.RunAll(packageDir)

	// Run teardown if defined
	if runner.Teardown != nil {
		if teardownErr := runner.Teardown(); teardownErr != nil {
			// Don't override the original error
			if err == nil {
				err = teardownErr
			} else {
				err = fmt.Errorf("%v; teardown failed: %w", err, teardownErr)
			}
		}
	}

	// Handle errors
	if err != nil || testRunner.HasFailures() {
		if !isParent {
			fmt.Fprintf(os.Stderr, "Tests failed: %v\n", err)
		}
		return err
	}

	if !isParent {
		fmt.Println("All tests passed!")
	}

	return nil
}
