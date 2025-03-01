package ale

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"strings"

	"lib-post-interchange/libale/format"
	"lib-post-interchange/libale/types"
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
		return nil, nil, nil, ErrSectionMissingHeading
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
		return nil, nil, nil, ErrSectionMissingColumn
	}

	// Skip any empty lines before column names
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		columnsLine := line
		if columnsLine == format.Data {
			return nil, nil, nil, ErrSectionIncompleteColumn
		}
		columnsArray, err := readTSVDataFirstLine(columnsLine)
		if err != nil {
			if _, ok := err.(*Error); ok {
				return nil, nil, nil, err // Pass through our custom errors
			}
			return nil, nil, nil, ErrParseFailedColumns.WithContext(fmt.Sprintf("csv error: %v", err))
		}
		for index, column := range columnsArray {
			columns = append(columns, makeColumn(column, index))
		}
		break
	}

	if len(columns) == 0 {
		return nil, nil, nil, ErrSectionIncompleteColumn
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
			return nil, nil, nil, ErrSectionMissingData
		}

		// Parse data row
		dataRow, err := readTSVDataFirstLine(line)
		if err != nil {
			if _, ok := err.(*Error); ok {
				return nil, nil, nil, err // Pass through our custom errors
			}
			return nil, nil, nil, ErrParseFailedData.WithContext(fmt.Sprintf("csv error: %v", err))
		}
		dataRows = append(dataRows, dataRow)
	}

	if !foundData || len(dataRows) == 0 {
		return nil, nil, nil, ErrSectionMissingData
	}

	rows, err := makeRowsFromDataRows(dataRows, columns)
	if err != nil {
		return nil, nil, nil, err
	}

	return headerFields, columns, rows, nil
}

// readTSVData handles the parsing of tab-separated value data
func readTSVData(input string) ([][]string, error) {
	if input == "" {
		return nil, ErrParseEmptyInput.WithContext("no data provided")
	}
	reader := csv.NewReader(strings.NewReader(input))
	reader.Comma = '\t'
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}
	if len(records) == 0 {
		return nil, ErrParseEmptyInput.WithContext("no data rows provided")
	}
	return records, nil
}

// readTSVDataFirstLine uses encoding/csv's Reader but only first line
func readTSVDataFirstLine(input string) ([]string, error) {
	if input == "" {
		return nil, ErrParseEmptyInput.WithContext("empty input string")
	}
	reader := csv.NewReader(strings.NewReader(input))
	reader.Comma = '\t'
	records, err := reader.Read()
	if err != nil {
		return nil, err
	}
	return records, nil
}

// makeColumn is a constructor for Column.
func makeColumn(name string, order int) types.Column {
	return types.Column{Name: name, Order: order}
}

// makeValue is a constructor for StringValue.
func makeValue(column types.Column, value string) types.StringValue {
	return types.StringValue{Column: column, Value: value}
}

// makeRow is a constructor for Row.
func makeRow(row []string, columns []types.Column) types.Row {
	var aleRow types.Row
	aleRow.Columns = columns
	aleRow.ValueMap = make(map[types.Column]types.Value)

	// Process all columns, padding with empty strings if row is too short
	for i, column := range columns {
		value := ""
		if i < len(row) {
			value = row[i]
		}
		aleValue := makeValue(column, value)
		aleRow.ValueMap[column] = aleValue
	}
	return aleRow
}

// makeRowsFromDataRows is a constructor for Row, iterating over multiple data rows
func makeRowsFromDataRows(rows [][]string, columns []types.Column) ([]types.Row, error) {
	var aleRows []types.Row
	for rowIndex, row := range rows {
		// Warn if row has more columns than defined
		if len(row) > len(columns) {
			fmt.Printf("Warning: row %d has %d columns, expected %d (extra columns will be ignored)\n",
				rowIndex, len(row), len(columns))
		}
		aleRow := makeRow(row, columns)
		aleRow.Order = rowIndex
		aleRows = append(aleRows, aleRow)
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
