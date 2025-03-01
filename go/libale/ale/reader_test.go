package ale

import (
	"testing"

	"lib-post-interchange/libale/ale/testutil"
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
			wantErr: false,
			check: func(t *testing.T, obj *types.Object) {
				if obj == nil {
					t.Fatal("Expected non-nil Object")
				}
				if len(obj.Columns) != 2 {
					t.Errorf("Got %d columns, want 2", len(obj.Columns))
				}
				if len(obj.Rows) != 1 {
					t.Errorf("Got %d rows, want 1", len(obj.Rows))
				}
				// Check that extra data was ignored
				row := obj.Rows[0]
				if len(row.ValueMap) != 2 {
					t.Errorf("Row has %d values, want 2", len(row.ValueMap))
				}
				// Verify the values we kept
				for col, val := range row.ValueMap {
					switch col.Name {
					case "Name":
						if val.String() != "A001" {
							t.Errorf("Name = %q, want %q", val.String(), "A001")
						}
					case "Scene":
						if val.String() != "1" {
							t.Errorf("Scene = %q, want %q", val.String(), "1")
						}
					}
				}
			},
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
		check   func(t *testing.T, rows []types.Row)
	}{
		{
			name: "valid rows",
			rows: [][]string{
				{"A001", "1"},
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
			wantErr: false,
			check: func(t *testing.T, rows []types.Row) {
				if len(rows) != 1 {
					t.Errorf("Got %d rows, want 1", len(rows))
					return
				}
				row := rows[0]
				if len(row.ValueMap) != len(columns) {
					t.Errorf("Row has %d values, want %d", len(row.ValueMap), len(columns))
				}
				// Verify the values we kept
				for col, val := range row.ValueMap {
					switch col.Name {
					case "Name":
						if val.String() != "A001" {
							t.Errorf("Name = %q, want %q", val.String(), "A001")
						}
					case "Scene":
						if val.String() != "1" {
							t.Errorf("Scene = %q, want %q", val.String(), "1")
						}
					}
				}
			},
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
			if tt.check != nil {
				tt.check(t, got)
			}
		})
	}
}

func TestReadAllSampleFiles(t *testing.T) {
	testutil.TestALEFiles(t, ReadFile)
}

func TestMakeRow(t *testing.T) {
	tests := []struct {
		name    string
		row     []string
		columns []types.Column
		want    map[string]string // map of column name to expected value
	}{
		{
			name: "exact match of columns and values",
			row:  []string{"A001", "1", "2"},
			columns: []types.Column{
				{Name: "Name", Order: 0},
				{Name: "Scene", Order: 1},
				{Name: "Take", Order: 2},
			},
			want: map[string]string{
				"Name":  "A001",
				"Scene": "1",
				"Take":  "2",
			},
		},
		{
			name: "fewer values than columns",
			row:  []string{"A001", "1"},
			columns: []types.Column{
				{Name: "Name", Order: 0},
				{Name: "Scene", Order: 1},
				{Name: "Take", Order: 2},
			},
			want: map[string]string{
				"Name":  "A001",
				"Scene": "1",
				"Take":  "", // Should be padded with empty string
			},
		},
		{
			name: "more values than columns",
			row:  []string{"A001", "1", "2", "extra"},
			columns: []types.Column{
				{Name: "Name", Order: 0},
				{Name: "Scene", Order: 1},
				{Name: "Take", Order: 2},
			},
			want: map[string]string{
				"Name":  "A001",
				"Scene": "1",
				"Take":  "2",
				// "extra" value should be ignored
			},
		},
		{
			name: "many more columns than values",
			row:  []string{"A001", "1", "2"},
			columns: []types.Column{
				{Name: "Name", Order: 0},
				{Name: "Scene", Order: 1},
				{Name: "Take", Order: 2},
				{Name: "Col4", Order: 3},
				{Name: "Col5", Order: 4},
				{Name: "Col6", Order: 5},
				{Name: "Col7", Order: 6},
				{Name: "Col8", Order: 7},
				{Name: "Col9", Order: 8},
				{Name: "Col10", Order: 9},
			},
			want: map[string]string{
				"Name":  "A001",
				"Scene": "1",
				"Take":  "2",
				"Col4":  "",
				"Col5":  "",
				"Col6":  "",
				"Col7":  "",
				"Col8":  "",
				"Col9":  "",
				"Col10": "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := makeRow(tt.row, tt.columns)

			// Check that we have all expected values
			if len(got.ValueMap) != len(tt.columns) {
				t.Errorf("makeRow() got %v values, want %v", len(got.ValueMap), len(tt.columns))
			}

			// Check each value matches expected
			for col, wantVal := range tt.want {
				var found bool
				for c, v := range got.ValueMap {
					if c.Name == col {
						found = true
						if v.String() != wantVal {
							t.Errorf("makeRow() value for column %q = %q, want %q", col, v.String(), wantVal)
						}
						break
					}
				}
				if !found {
					t.Errorf("makeRow() missing value for column %q", col)
				}
			}

			// Check column order is preserved
			for i, col := range tt.columns {
				found := false
				for c := range got.ValueMap {
					if c.Name == col.Name {
						if c.Order != i {
							t.Errorf("makeRow() column %q has order %d, want %d", c.Name, c.Order, i)
						}
						found = true
						break
					}
				}
				if !found {
					t.Errorf("makeRow() missing column %q", col.Name)
				}
			}
		})
	}
}
