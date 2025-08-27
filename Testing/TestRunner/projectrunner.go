// filepath: /Users/skelly/Documents/GoProjects/GoLib/TestRunner/projectrunner.go
package TestRunner

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"
)

// ProjectTestResult represents the result of running tests in a package
type ProjectTestResult struct {
	// Package name or directory
	Package string
	// Whether the test succeeded
	Success bool
	// Error if the test failed
	Error error
	// Time taken to run the test
	Time time.Duration
}

// ProjectRunnerConfig configures how the project-level test runner behaves
type ProjectRunnerConfig struct {
	// Root directory to search for test runners
	RootDir string
	// Run tests in parallel
	Parallel bool
	// Verbose output
	Verbose bool
	// Generate coverage reports
	Coverage bool
	// Coverage mode
	CoverageMode string
	// Test pattern to filter tests
	TestPattern string
	// Run legacy shell scripts instead of Go test runners
	UseLegacyScripts bool
	// Also run original shell scripts alongside Go runners
	RunShellScripts bool
	// Debug output
	Debug bool
	// Directories to ignore
	IgnoreDirs []string
	// Name of the project
	ProjectName string
	// Packages to test (if nil, all packages will be tested)
	Packages []string
}

// DefaultProjectConfig returns the default configuration for a project runner
func DefaultProjectConfig(projectName string) *ProjectRunnerConfig {
	// Get the current directory
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting current directory: %v\n", err)
		os.Exit(1)
	}

	return &ProjectRunnerConfig{
		RootDir:          currentDir,
		Parallel:         false,
		Verbose:          true,
		Coverage:         false,
		CoverageMode:     "atomic",
		TestPattern:      "",
		UseLegacyScripts: false,
		RunShellScripts:  false,
		Debug:            false,
		IgnoreDirs:       []string{"vendor", ".git", "built", "bin", "tools"},
		ProjectName:      projectName,
	}
}

// FindFiles finds all files with the given name in the directory tree
func FindFiles(root, fileName string, ignoreDirs []string) ([]string, error) {
	var files []string

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip ignored directories
		if info.IsDir() {
			base := filepath.Base(path)
			for _, ignoreDir := range ignoreDirs {
				if base == ignoreDir {
					return filepath.SkipDir
				}
			}
			return nil
		}

		// Check if this is a file we're looking for
		if filepath.Base(path) == fileName {
			files = append(files, path)
		}

		return nil
	})

	return files, err
}

// RunProjectTests finds and executes all test runners in the project
func RunProjectTests(config *ProjectRunnerConfig) ([]ProjectTestResult, error) {
	// Make config.RootDir absolute if it isn't already
	absRootDir, err := filepath.Abs(config.RootDir)
	if err != nil {
		return nil, fmt.Errorf("error getting absolute path: %w", err)
	}
	config.RootDir = absRootDir

	// Print header
	fmt.Printf("=== %s Test Runner ===\n", config.ProjectName)
	fmt.Printf("Root directory: %s\n", config.RootDir)

	var results []ProjectTestResult
	var mutex sync.Mutex
	var wg sync.WaitGroup

	// Find all test runners
	var runners []string
	if config.UseLegacyScripts {
		// Find all 't' shell scripts
		runners, _ = FindFiles(config.RootDir, "t", config.IgnoreDirs)
	} else {
		// Find all runtests.go files
		runners, _ = FindFiles(config.RootDir, "runtests.go", config.IgnoreDirs)
	}

	// Debug output
	if config.Debug {
		fmt.Printf("Found %d test runners:\n", len(runners))
		for _, runner := range runners {
			fmt.Printf("  %s\n", runner)
		}
	} else {
		fmt.Printf("Found %d test runners\n", len(runners))
	}

	// Create a channel to track errors
	errChan := make(chan error, len(runners))

	// Function to run a single test
	runTest := func(runnerPath string) {
		defer wg.Done()

		// Extract directory and package name
		dir := filepath.Dir(runnerPath)
		packageName := filepath.Base(dir)

		// Skip if this is the root directory
		if dir == config.RootDir {
			return
		}

		if config.Verbose {
			fmt.Printf("Running tests for package: %s\n", packageName)
		}

		start := time.Now()
		var cmd *exec.Cmd

		if config.UseLegacyScripts {
			// Run the shell script
			cmd = exec.Command("/bin/bash", runnerPath)
		} else {
			// Build command to run the Go test runner
			args := []string{"run", runnerPath, "-parent=true"}

			if config.Verbose {
				args = append(args, "-v")
			}

			if config.Debug {
				args = append(args, "-debug")
			}

			if config.RunShellScripts {
				args = append(args, "-shell")
			}

			if config.Coverage {
				args = append(args, "-covermode="+config.CoverageMode)
			}

			if config.TestPattern != "" {
				args = append(args, "-run="+config.TestPattern)
			}

			cmd = exec.Command("go", args...)
		}

		// Set working directory
		cmd.Dir = dir

		// Connect to stdout/stderr
		if config.Verbose {
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
		}

		// Run the test
		err := cmd.Run()
		elapsed := time.Since(start)

		// Record result
		result := ProjectTestResult{
			Package: packageName,
			Success: err == nil,
			Error:   err,
			Time:    elapsed,
		}

		mutex.Lock()
		results = append(results, result)
		mutex.Unlock()

		if err != nil {
			errChan <- fmt.Errorf("tests failed in package %s: %v", packageName, err)
		}
	}

	// Function to ask user if they want to continue
	promptToContinue := func() bool {
		fmt.Print("Continue to iterate? (y/n): ")
		var response string
		fmt.Scanln(&response)
		return response == "y" || response == "Y" || response == ""
	}

	// Run tests in parallel or sequentially
	if config.Parallel {
		for _, runner := range runners {
			wg.Add(1)
			go runTest(runner)
		}
	} else {
		for i, runner := range runners {
			wg.Add(1)
			runTest(runner)

			// Check if we had an error
			select {
			case err := <-errChan:
				// If we're not in verbose mode, or the user doesn't want to continue after an error, return
				if !config.Verbose || !promptToContinue() {
					return results, err
				}
			default:
				// No error, but still ask if not the last runner
				if i < len(runners)-1 && !promptToContinue() {
					return results, nil
				}
			}
		}
	}

	// Wait for all tests to complete
	wg.Wait()
	close(errChan)

	// Check for errors
	for err := range errChan {
		// Just return the first error
		return results, err
	}

	return results, nil
}

// PrintTestSummary prints a summary of test results
func PrintTestSummary(results []ProjectTestResult, totalTime time.Duration) {
	passed := 0
	failed := 0

	for _, result := range results {
		if result.Success {
			passed++
		} else {
			failed++
		}
	}

	fmt.Println("\n=== Test Summary ===")
	fmt.Printf("Total packages: %d\n", len(results))
	fmt.Printf("Passed: %d\n", passed)
	fmt.Printf("Failed: %d\n", failed)
	fmt.Printf("Total time: %s\n", totalTime)

	if failed > 0 {
		fmt.Println("\n=== Failed Packages ===")
		for _, result := range results {
			if !result.Success {
				fmt.Printf("%s (%s): %v\n", result.Package, result.Time, result.Error)
			}
		}
	}
}

// HasTestFailures returns true if any tests failed
func HasTestFailures(results []ProjectTestResult) bool {
	for _, result := range results {
		if !result.Success {
			return true
		}
	}
	return false
}
