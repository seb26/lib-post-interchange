package ale

import (
	"fmt"
	"strings"
	"testing"

	"ale/types"
)

func TestWriteMatchesSample(t *testing.T) {
	// Read the sample file
	ale, err := ReadFile("internal/testdata/basic.ale")
	if err != nil {
		t.Fatalf("Failed to read sample file: %v", err)
	}

	// Write the ALE object back to string
	output, err := Write(ale)
	if err != nil {
		t.Fatalf("Failed to write ALE object: %v", err)
	}

	// Read the output back as an ALE object to compare structures
	outputAle, err := Read(output)
	if err != nil {
		t.Fatalf("Failed to read written output: %v", err)
	}

	// Compare the two ALE objects
	compareALEObjects(t, ale, outputAle)
}

func TestWriteColumnOrderPreservation(t *testing.T) {
	// Create ALE object with specifically ordered columns
	ale := &types.Object{
		HeaderFields: []types.Field{
			types.FieldDelimiter{BaseField: types.BaseField{Key: "FIELD_DELIM", Value: "TABS"}},
		},
		Columns: []types.Column{
			{Name: "Last", Order: 2},
			{Name: "First", Order: 0},
			{Name: "Middle", Order: 1},
		},
		Rows: []types.Row{
			{
				Columns: []types.Column{
					{Name: "Last", Order: 2},
					{Name: "First", Order: 0},
					{Name: "Middle", Order: 1},
				},
				ValueMap: map[types.Column]types.Value{
					{Name: "Last", Order: 2}: types.StringValue{
						Column: types.Column{Name: "Last", Order: 2},
						Value:  "C",
					},
					{Name: "First", Order: 0}: types.StringValue{
						Column: types.Column{Name: "First", Order: 0},
						Value:  "A",
					},
					{Name: "Middle", Order: 1}: types.StringValue{
						Column: types.Column{Name: "Middle", Order: 1},
						Value:  "B",
					},
				},
			},
		},
	}

	// Write and read back
	output, err := Write(ale)
	if err != nil {
		t.Fatalf("Failed to write ALE object: %v", err)
	}

	// Verify the raw output has columns in correct order
	lines := strings.Split(output, "\n")
	var columnLine string
	for _, line := range lines {
		if strings.Contains(line, "First\tMiddle\tLast") {
			columnLine = line
			break
		}
	}
	if columnLine == "" {
		t.Fatal("Column line not found in output")
	}
	expectedOrder := "First\tMiddle\tLast"
	if columnLine != expectedOrder {
		t.Errorf("Column order not preserved in output:\nExpected: %q\nGot: %q", expectedOrder, columnLine)
	}

	// Read back and verify structure
	outputAle, err := Read(output)
	if err != nil {
		t.Fatalf("Failed to read written output: %v", err)
	}

	// Verify column order in parsed object
	for i, expectedName := range []string{"First", "Middle", "Last"} {
		if outputAle.Columns[i].Name != expectedName {
			t.Errorf("Column %d: expected name %q, got %q", i, expectedName, outputAle.Columns[i].Name)
		}
		if outputAle.Columns[i].Order != i {
			t.Errorf("Column %q: expected order %d, got %d", expectedName, i, outputAle.Columns[i].Order)
		}
	}
}

func TestWriteRowOrderPreservation(t *testing.T) {
	// Create ALE object with specifically ordered rows
	ale := &types.Object{
		HeaderFields: []types.Field{
			types.FieldDelimiter{BaseField: types.BaseField{Key: "FIELD_DELIM", Value: "TABS"}},
		},
		Columns: []types.Column{
			{Name: "Name", Order: 0},
			{Name: "Value", Order: 1},
		},
	}

	// Create 5 rows with distinct values
	rows := make([]types.Row, 5)
	for i := 0; i < 5; i++ {
		rows[i] = types.Row{
			Order: i,
			Columns: []types.Column{
				{Name: "Name", Order: 0},
				{Name: "Value", Order: 1},
			},
			ValueMap: map[types.Column]types.Value{
				{Name: "Name", Order: 0}: types.StringValue{
					Column: types.Column{Name: "Name", Order: 0},
					Value:  fmt.Sprintf("Name%d", i),
				},
				{Name: "Value", Order: 1}: types.StringValue{
					Column: types.Column{Name: "Value", Order: 1},
					Value:  fmt.Sprintf("Value%d", i),
				},
			},
		}
	}
	ale.Rows = rows

	// Write and read back
	output, err := Write(ale)
	if err != nil {
		t.Fatalf("Failed to write ALE object: %v", err)
	}

	// Verify the raw output has rows in correct order
	lines := strings.Split(output, "\n")
	var dataLines []string
	inDataSection := false
	for _, line := range lines {
		if line == "Data" {
			inDataSection = true
			continue
		}
		if inDataSection && line != "" {
			dataLines = append(dataLines, line)
		}
	}

	// Check raw output order
	for i, line := range dataLines {
		expected := fmt.Sprintf("Name%d\tValue%d", i, i)
		if line != expected {
			t.Errorf("Row %d: expected %q, got %q", i, expected, line)
		}
	}

	// Read back and verify structure
	outputAle, err := Read(output)
	if err != nil {
		t.Fatalf("Failed to read written output: %v", err)
	}

	// Verify row order in parsed object
	for i, row := range outputAle.Rows {
		nameVal := row.ValueMap[types.Column{Name: "Name", Order: 0}].String()
		valueVal := row.ValueMap[types.Column{Name: "Value", Order: 1}].String()
		expectedName := fmt.Sprintf("Name%d", i)
		expectedValue := fmt.Sprintf("Value%d", i)

		if nameVal != expectedName {
			t.Errorf("Row %d: expected name %q, got %q", i, expectedName, nameVal)
		}
		if valueVal != expectedValue {
			t.Errorf("Row %d: expected value %q, got %q", i, expectedValue, valueVal)
		}
		if row.Order != i {
			t.Errorf("Row with values (%q, %q): expected order %d, got %d", nameVal, valueVal, i, row.Order)
		}
	}
}

func compareALEObjects(t *testing.T, expected, actual *types.Object) {
	t.Helper()

	// Compare header fields
	if len(expected.HeaderFields) != len(actual.HeaderFields) {
		t.Errorf("Header fields length mismatch: expected %d, got %d", len(expected.HeaderFields), len(actual.HeaderFields))
		return
	}

	// Create maps for easier comparison
	expectedHeaders := make(map[string]string)
	actualHeaders := make(map[string]string)

	for _, field := range expected.HeaderFields {
		expectedHeaders[field.GetKey()] = field.GetValue()
	}
	for _, field := range actual.HeaderFields {
		actualHeaders[field.GetKey()] = field.GetValue()
	}

	for key, expectedValue := range expectedHeaders {
		if actualValue, ok := actualHeaders[key]; !ok {
			t.Errorf("Missing header field %s", key)
		} else if expectedValue != actualValue {
			t.Errorf("Header field %s value mismatch: expected %q, got %q", key, expectedValue, actualValue)
		}
	}

	// Compare columns
	if len(expected.Columns) != len(actual.Columns) {
		t.Errorf("Columns length mismatch: expected %d, got %d", len(expected.Columns), len(actual.Columns))
		return
	}

	for i, expectedCol := range expected.Columns {
		actualCol := actual.Columns[i]
		if expectedCol.Name != actualCol.Name {
			t.Errorf("Column %d name mismatch: expected %q, got %q", i, expectedCol.Name, actualCol.Name)
		}
		if expectedCol.Order != actualCol.Order {
			t.Errorf("Column %d order mismatch: expected %d, got %d", i, expectedCol.Order, actualCol.Order)
		}
	}

	// Compare rows
	if len(expected.Rows) != len(actual.Rows) {
		t.Errorf("Rows length mismatch: expected %d, got %d", len(expected.Rows), len(actual.Rows))
		return
	}

	for i, expectedRow := range expected.Rows {
		actualRow := actual.Rows[i]

		// Compare values for each column
		for _, col := range expected.Columns {
			expectedVal, expectedOk := expectedRow.ValueMap[col]
			actualVal, actualOk := actualRow.ValueMap[col]

			if !expectedOk && !actualOk {
				continue
			}
			if expectedOk != actualOk {
				t.Errorf("Row %d, column %q presence mismatch: expected %v, got %v", i, col.Name, expectedOk, actualOk)
				continue
			}
			if expectedVal.String() != actualVal.String() {
				t.Errorf("Row %d, column %q value mismatch: expected %q, got %q", i, col.Name, expectedVal.String(), actualVal.String())
			}
		}
	}
}

func TestWriteNilObject(t *testing.T) {
	_, err := Write(nil)
	if err == nil {
		t.Error("Expected error when writing nil object, got nil")
	}
	if !strings.Contains(err.Error(), "cannot write nil ALE object") {
		t.Errorf("Expected error message to contain 'cannot write nil ALE object', got %q", err.Error())
	}
}
