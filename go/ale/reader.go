package ale

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"ale/errors"
	"ale/format"
	"ale/types"
)

// ReadFile reads and parses an ALE file from the filesystem.
func ReadFile(filepath string) (*types.Object, error) {
	_, err := os.Stat(filepath)
	if os.IsNotExist(err) {
		return nil, err
	}
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	dataString := string(data)
	return Read(dataString)
}

// Read parses ALE data from a string.
func Read(input string) (*types.Object, error) {
	headerFields, columns, rows, err := read(input)
	if err != nil {
		return nil, err
	}
	ale := types.Object{
		HeaderFields: headerFields,
		Columns:      columns,
		Rows:         rows,
	}
	assignHeaderFields(&ale)
	return &ale, nil
}

// read is the core parsing function that handles the ALE file format's three-part structure:
// headers, columns, and data. It enforces the format's rules and extracts structured data.
func read(input string) ([]types.Field, []types.Column, []types.Row, error) {
	var headerFields []types.Field
	var columns []types.Column

	scanner := bufio.NewScanner(strings.NewReader(input))

	// First line should be "Heading"
	if !scanner.Scan() || scanner.Text() != format.Heading {
		return nil, nil, nil, errors.ErrInputMissingHeading
	}

	// Read header fields until empty line
	var headerData strings.Builder
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			break
		}
		headerData.WriteString(line + "\n")
	}

	// Parse header fields by splitting each line on first tab
	headerLines := strings.Split(strings.TrimRight(headerData.String(), "\r\n"), "\n")
	for _, line := range headerLines {
		line = strings.TrimRight(line, "\r") // Handle any remaining \r in CRLF files
		parts := strings.SplitN(line, "\t", 2)
		if len(parts) != 2 {
			continue // Skip malformed lines
		}
		key := parts[0]
		value := strings.ReplaceAll(parts[1], "\t", " ") // Replace any tabs in value with spaces
		headerFields = append(headerFields, types.BaseField{
			Key:   key,
			Value: value,
		})
	}

	// Next line should be "Column"
	if !scanner.Scan() || scanner.Text() != format.Column {
		return nil, nil, nil, errors.ErrInputMissingColumn
	}

	// Skip any empty lines before column names
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		columnsLine := line
		if columnsLine == format.Data {
			return nil, nil, nil, errors.ErrInputIncompleteColumn
		}
		columnsArray, err := readTSVLine(columnsLine)
		if err != nil {
			if _, ok := err.(*errors.Error); ok {
				return nil, nil, nil, err // Pass through our custom errors
			}
			return nil, nil, nil, errors.ErrInputFailedColumns.WithContext(fmt.Sprintf("csv error: %v", err))
		}
		for index, column := range columnsArray {
			columns = append(columns, makeColumn(column, index))
		}
		break
	}

	if len(columns) == 0 {
		return nil, nil, nil, errors.ErrInputIncompleteColumn
	}

	// Skip empty lines until we find "Data" or data content
	foundData := false
	var dataRows [][]string
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		if line == format.Data {
			foundData = true
			continue
		}
		if !foundData {
			return nil, nil, nil, errors.ErrInputMissingData
		}

		// Parse data row
		dataRow, err := readTSVLine(line)
		if err != nil {
			if _, ok := err.(*errors.Error); ok {
				return nil, nil, nil, err // Pass through our custom errors
			}
			return nil, nil, nil, errors.ErrInputFailedData.WithContext(fmt.Sprintf("csv error: %v", err))
		}
		dataRows = append(dataRows, dataRow)
	}

	if !foundData {
		return nil, nil, nil, errors.ErrInputMissingData
	}

	if len(dataRows) == 0 {
		return nil, nil, nil, errors.ErrInputMissingData
	}

	rows, err := makeRowsFromDataRows(dataRows, columns)
	if err != nil {
		return nil, nil, nil, err
	}

	return headerFields, columns, rows, nil
}

// readTSVLine parses a TSV line into fields
func readTSVLine(input string) ([]string, error) {
	if input == "" {
		return nil, errors.ErrInputEmpty.WithContext("empty input string")
	}

	fields := make([]string, 0, 1)
	fieldStart := 0
	for i := 0; i < len(input); i++ {
		if input[i] == '\t' {
			fields = append(fields, input[fieldStart:i])
			fieldStart = i + 1
		}
	}
	// Add final field
	fields = append(fields, input[fieldStart:])

	return fields, nil
}

// makeColumn is a constructor for Column.
func makeColumn(name string, order int) types.Column {
	return types.Column{Name: name, Order: order}
}

// makeValue is a constructor for StringValue.
func makeValue(column types.Column, value string) types.StringValue {
	return types.StringValue{Column: column, Value: value}
}

// makeRow creates a Row from a slice of values and column definitions.
// If the row has fewer values than columns, the remaining columns are filled with empty strings.
// If the row has more values than columns, the extra values are ignored.
func makeRow(row []string, columns []types.Column) types.Row {
	aleRow := types.Row{
		Columns:  columns,
		ValueMap: make(map[types.Column]types.Value, len(columns)),
	}

	// Process each column in order, ensuring the column order is preserved
	for _, column := range columns {
		value := ""
		if column.Order < len(row) {
			value = row[column.Order]
		}
		aleRow.ValueMap[column] = makeValue(column, value)
	}

	return aleRow
}

// makeRowsFromDataRows creates a slice of Rows from raw data rows and column definitions.
// It handles cases where rows have more or fewer values than columns:
// - If a row has more values than columns, continue but warn that the extra values are ignored
// - If a row has fewer values than columns, the remaining columns are filled with empty strings
func makeRowsFromDataRows(rows [][]string, columns []types.Column) ([]types.Row, error) {
	if len(rows) == 0 {
		return nil, nil
	}

	aleRows := make([]types.Row, len(rows))
	for i, row := range rows {
		// Warn about extra data that will be ignored
		if len(row) > len(columns) {
			fmt.Printf("Warning: row %d has %d values, expected %d (extra values will be ignored)\n",
				i, len(row), len(columns))
		}

		aleRow := makeRow(row, columns)
		aleRow.Order = i
		aleRows[i] = aleRow
	}

	return aleRows, nil
}

// assignHeaderFields assigns header fields to their specific types in the Object.
func assignHeaderFields(ale *types.Object) {
	for _, field := range ale.HeaderFields {
		switch field.GetKey() {
		case "FIELD_DELIM":
			ale.FieldDelimiter = types.FieldDelimiter{BaseField: types.BaseField{Key: field.GetKey(), Value: field.GetValue()}}
		case "FPS":
			ale.FPS = types.FrameRate{BaseField: types.BaseField{Key: field.GetKey(), Value: field.GetValue()}}
		case "AUDIO_FORMAT":
			ale.AudioFormat = types.AudioFormat{BaseField: types.BaseField{Key: field.GetKey(), Value: field.GetValue()}}
		case "VIDEO_FORMAT":
			ale.VideoFormat = types.VideoFormat{BaseField: types.BaseField{Key: field.GetKey(), Value: field.GetValue()}}
		case "FILM_FORMAT":
			ale.FilmFormat = types.FilmFormat{BaseField: types.BaseField{Key: field.GetKey(), Value: field.GetValue()}}
		case "TAPE":
			ale.Tape = types.Tape{BaseField: types.BaseField{Key: field.GetKey(), Value: field.GetValue()}}
		}
	}
}
