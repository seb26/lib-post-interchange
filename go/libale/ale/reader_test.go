package ale

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"lib-post-interchange/libale/types"
)

func TestRead(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
		check   func(*testing.T, *types.Object)
	}{
		{
			name: "valid ale data with all fields",
			input: `Heading
FIELD_DELIM	TABS
VIDEO_FORMAT	1080
AUDIO_FORMAT	48kHz
FPS	23.98
FILM_FORMAT	35 mm
TAPE	A001

Column
Name	Scene	Take

Data
A001	1	1
A001	1	2
`,
			wantErr: false,
			check: func(t *testing.T, obj *types.Object) {
				if obj == nil {
					t.Fatal("Expected non-nil Object")
				}

				// Check header fields
				if obj.FieldDelimiter.GetValue() != "TABS" {
					t.Errorf("FieldDelimiter = %v, want TABS", obj.FieldDelimiter.GetValue())
				}
				if obj.VideoFormat.GetValue() != "1080" {
					t.Errorf("VideoFormat = %v, want 1080", obj.VideoFormat.GetValue())
				}
				if obj.AudioFormat.GetValue() != "48kHz" {
					t.Errorf("AudioFormat = %v, want 48kHz", obj.AudioFormat.GetValue())
				}
				if obj.FPS.GetValue() != "23.98" {
					t.Errorf("FPS = %v, want 23.98", obj.FPS.GetValue())
				}
				if obj.FilmFormat.GetValue() != "35 mm" {
					t.Errorf("FilmFormat = %v, want 35 mm", obj.FilmFormat.GetValue())
				}
				if obj.Tape.GetValue() != "A001" {
					t.Errorf("Tape = %v, want A001", obj.Tape.GetValue())
				}

				// Check columns
				expectedColumns := []string{"Name", "Scene", "Take"}
				if len(obj.Columns) != len(expectedColumns) {
					t.Errorf("Got %d columns, want %d", len(obj.Columns), len(expectedColumns))
				}
				for i, col := range obj.Columns {
					if col.Name != expectedColumns[i] {
						t.Errorf("Column[%d].Name = %v, want %v", i, col.Name, expectedColumns[i])
					}
					if col.Order != i {
						t.Errorf("Column[%d].Order = %v, want %v", i, col.Order, i)
					}
				}

				// Check rows
				if len(obj.Rows) != 2 {
					t.Errorf("Got %d rows, want 2", len(obj.Rows))
				}
			},
		},
		{
			name: "missing header section",
			input: `Column
Name	Scene	Take

Data
A001	1	1
`,
			wantErr: true,
		},
		{
			name: "missing column section",
			input: `Heading
FIELD_DELIM	TABS

Data
A001	1	1
`,
			wantErr: true,
		},
		{
			name: "missing data section",
			input: `Heading
FIELD_DELIM	TABS

Column
Name	Scene	Take
`,
			wantErr: true,
		},
		{
			name: "empty data section",
			input: `Heading
FIELD_DELIM	TABS

Column
Name	Scene	Take

Data
`,
			wantErr: true,
		},
		{
			name: "mismatched columns and data",
			input: `Heading
FIELD_DELIM	TABS

Column
Name	Scene

Data
A001	1	1
`,
			wantErr: true,
		},
		{
			name:    "empty input",
			input:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obj, err := Read(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Read() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && tt.check != nil {
				tt.check(t, obj)
			}
		})
	}
}

func TestReadFile(t *testing.T) {
	tests := []struct {
		name     string
		filepath string
		wantErr  bool
		check    func(*testing.T, *types.Object)
	}{
		{
			name:     "non-existent file",
			filepath: "../../../samples/ALE/nonexistent.ale",
			wantErr:  true,
		},
		{
			name:     "valid sample file",
			filepath: "../../../samples/ALE/A001R1AA_AVID.ale",
			wantErr:  false,
			check: func(t *testing.T, obj *types.Object) {
				if obj == nil {
					t.Fatal("Expected non-nil Object")
				}

				// Check header fields
				if obj.FieldDelimiter.GetValue() != "TABS" {
					t.Errorf("FieldDelimiter = %v, want TABS", obj.FieldDelimiter.GetValue())
				}

				// Check we have at least one column
				if len(obj.Columns) == 0 {
					t.Error("Expected at least one column")
					return
				}

				// Check column order is preserved
				for i, col := range obj.Columns {
					if col.Order != i {
						t.Errorf("Column[%d].Order = %v, want %v", i, col.Order, i)
					}
				}

				// Check we have at least one row
				if len(obj.Rows) == 0 {
					t.Error("Expected at least one row")
					return
				}

				// Check first row has values for all columns
				firstRow := obj.Rows[0]
				if len(firstRow.ValueMap) != len(obj.Columns) {
					t.Errorf("First row has %d values, want %d", len(firstRow.ValueMap), len(obj.Columns))
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obj, err := ReadFile(tt.filepath)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && tt.check != nil {
				tt.check(t, obj)
			}
		})
	}
}

func TestReadTSVData(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    [][]string
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid tsv data",
			input: `A001	1	1
A001	1	2`,
			want: [][]string{
				{"A001", "1", "1"},
				{"A001", "1", "2"},
			},
			wantErr: false,
		},
		{
			name:    "empty input",
			input:   "",
			want:    nil,
			wantErr: true,
			errMsg:  "ale: [2.3] empty input: no data provided",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := readTSVData(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("readTSVData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err.Error() != tt.errMsg {
				t.Errorf("readTSVData() error = %v, want %v", err.Error(), tt.errMsg)
				return
			}
			if !tt.wantErr {
				if len(got) != len(tt.want) {
					t.Errorf("readTSVData() got %v rows, want %v", len(got), len(tt.want))
					return
				}
				for i := range got {
					if len(got[i]) != len(tt.want[i]) {
						t.Errorf("Row %d: got %v columns, want %v", i, len(got[i]), len(tt.want[i]))
						continue
					}
					for j := range got[i] {
						if got[i][j] != tt.want[i][j] {
							t.Errorf("Row %d, Col %d: got %v, want %v", i, j, got[i][j], tt.want[i][j])
						}
					}
				}
			}
		})
	}
}

func TestMakeRowsFromDataRows(t *testing.T) {
	columns := []types.Column{
		{Name: "Name", Order: 0},
		{Name: "Scene", Order: 1},
	}

	tests := []struct {
		name    string
		rows    [][]string
		columns []types.Column
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid rows",
			rows: [][]string{
				{"A001", "1"},
				{"A001", "2"},
			},
			columns: columns,
			wantErr: false,
		},
		{
			name: "mismatched columns",
			rows: [][]string{
				{"A001", "1", "extra"},
			},
			columns: columns,
			wantErr: true,
			errMsg:  "ale: [2.4] row has mismatched column count: row 0 has 3 columns, expected 2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := makeRowsFromDataRows(tt.rows, tt.columns)
			if (err != nil) != tt.wantErr {
				t.Errorf("makeRowsFromDataRows() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err.Error() != tt.errMsg {
				t.Errorf("makeRowsFromDataRows() error = %v, want %v", err.Error(), tt.errMsg)
				return
			}
			if !tt.wantErr {
				if len(got) != len(tt.rows) {
					t.Errorf("makeRowsFromDataRows() got %v rows, want %v", len(got), len(tt.rows))
				}
			}
		})
	}
}

func TestReadAllSampleFiles(t *testing.T) {
	// Find all .ale files in samples directory
	pattern := "../../../samples/ALE/*/*.ale"
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

			// Read file content
			content, err := os.ReadFile(filePath)
			if err != nil {
				result.err = err
				t.Errorf("Failed to read file: %v", err)
				return
			}

			// Try to parse the file
			obj, err := Read(string(content))
			if err != nil {
				result.err = err
				t.Errorf("Failed to parse file: %v", err)
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
