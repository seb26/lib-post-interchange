// Package testutil provides shared test utilities for ALE parsing
package testutil

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"lib-post-interchange/libale/types"
)

// ReadFileTestFunc represents a function that reads an ALE file and returns a parsed object
type ReadFileTestFunc func(string) (*types.Object, error)

// TestALEFiles runs comprehensive tests on all ALE files in the samples directory
func TestALEFiles(t *testing.T, readFileFn ReadFileTestFunc) {
	// Find all .ale files in samples directory
	pattern := "../../../../samples/ALE/*/*.ale"
	matches, err := filepath.Glob(pattern)
	if err != nil {
		t.Fatalf("Failed to glob sample files: %v", err)
	}
	if len(matches) == 0 {
		t.Fatalf("No .ale files found in samples directory")
	}

	// Track test results
	type testResult struct {
		path   string
		passed bool
		err    error
	}
	results := make([]testResult, 0, len(matches))

	// Test each sample file
	for _, filePath := range matches {
		result := testResult{path: filePath}

		t.Run(filepath.Base(filePath), func(t *testing.T) {
			// Verify file exists and is readable
			info, err := os.Stat(filePath)
			if err != nil {
				result.err = err
				t.Errorf("Failed to stat file: %v", err)
				return
			}
			if info.Size() == 0 {
				result.err = fmt.Errorf("file is empty")
				t.Error(result.err)
				return
			}

			// Try to read and parse the file
			obj, err := readFileFn(filePath)
			if err != nil {
				result.err = err
				t.Errorf("Failed to read file: %v", err)
				return
			}
			if obj == nil {
				result.err = fmt.Errorf("expected non-nil Object for valid file")
				t.Error(result.err)
				return
			}

			// Basic validation of the parsed object
			if len(obj.Columns) == 0 {
				result.err = fmt.Errorf("expected at least one column")
				t.Error(result.err)
				return
			}
			if len(obj.Rows) == 0 {
				result.err = fmt.Errorf("expected at least one row")
				t.Error(result.err)
				return
			}
			if obj.FieldDelimiter.GetValue() == "" {
				result.err = fmt.Errorf("expected non-empty field delimiter")
				t.Error(result.err)
				return
			}

			// Validate column order
			for i, col := range obj.Columns {
				if col.Order != i {
					result.err = fmt.Errorf("column[%d].Order = %v, want %v", i, col.Order, i)
					t.Error(result.err)
					return
				}
			}

			// Validate row data
			for i, row := range obj.Rows {
				if len(row.ValueMap) != len(obj.Columns) {
					result.err = fmt.Errorf("row %d has %d values, want %d", i, len(row.ValueMap), len(obj.Columns))
					t.Error(result.err)
					return
				}
			}

			// If we got here, all checks passed
			result.passed = true
		})

		results = append(results, result)
	}

	// Print summary
	t.Logf("\nTest Summary:")
	t.Logf("------------")
	passed := 0
	for _, r := range results {
		status := "❌ FAILED"
		if r.passed {
			status = "✅ PASSED"
			passed++
		}
		t.Logf("%s: %s", status, filepath.Base(r.path))
		if r.err != nil {
			t.Logf("   Error: %v", r.err)
		}
	}
	t.Logf("------------")
	t.Logf("Total: %d, Passed: %d, Failed: %d", len(results), passed, len(results)-passed)
}
