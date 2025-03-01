package libale

import (
	"bufio"
	"encoding/csv"
	"os"
	"strings"
)

// ReadFile provides the main entry point for loading ALE data from the filesystem.
func ReadFile(filepath string) (*ALEObject, error) {
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

// Read serves as the primary interface for parsing ALE data from any string source.
// It coordinates the parsing process and ensures proper initialization of the ALEObject.
func Read(input string) (*ALEObject, error) {
	aleHeaderFields, aleColumns, aleRows, err := read(input)
	if err != nil {
		return nil, err
	}
	ale := ALEObject{
		HeaderFields: aleHeaderFields,
		Columns:      aleColumns,
		Rows:         aleRows,
	}
	ale = AssignHeaderFieldsToObject(ale)
	return &ale, nil
}

// read is the core parsing function that handles the ALE file format's three-part structure:
// headers, columns, and data. It enforces the format's rules and extracts structured data.
func read(input string) ([]ALEField, []ALEColumn, []ALERow, error) {
	var headerFields []ALEField
	var columns []ALEColumn

	scanner := bufio.NewScanner(strings.NewReader(input))

	// First line should be "Heading"
	if !scanner.Scan() || scanner.Text() != "Heading" {
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

	// Parse header fields
	fieldsArray, err := readTSVData(headerData.String())
	if err != nil {
		return nil, nil, nil, ErrParseFailedHeader
	}
	for _, field := range fieldsArray {
		if len(field) != 2 {
			continue
		}
		key := field[0]
		value := field[1]
		constructor, err := ToType(key)
		if err != nil {
			return nil, nil, nil, ErrFieldInvalidHeader
		}
		headerFields = append(headerFields, constructor(value))
	}

	// Next line should be "Column"
	if !scanner.Scan() || scanner.Text() != "Column" {
		return nil, nil, nil, ErrSectionMissingColumn
	}

	// Read column names
	if !scanner.Scan() {
		return nil, nil, nil, ErrSectionIncompleteColumn
	}
	columnsLine := scanner.Text()
	columnsArray, err := readTSVDataFirstLine(columnsLine)
	if err != nil {
		return nil, nil, nil, ErrParseFailedColumns
	}
	for index, column := range columnsArray {
		columns = append(columns, makeALEColumn(column, index))
	}

	// Skip empty line
	if !scanner.Scan() || scanner.Text() != "" {
		return nil, nil, nil, ErrFormatMalformedColumn
	}

	// Next line should be "Data"
	if !scanner.Scan() || scanner.Text() != "Data" {
		return nil, nil, nil, ErrSectionMissingData
	}

	// Read all data rows
	var dataBuilder strings.Builder
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		dataBuilder.WriteString(line + "\n")
	}

	if err := scanner.Err(); err != nil {
		return nil, nil, nil, ErrParseFailedContent
	}

	// Parse data rows
	dataRows, err := readTSVData(dataBuilder.String())
	if err != nil {
		return nil, nil, nil, ErrParseFailedData
	}
	rows := makeALERowsFromDataRows(dataRows, columns)

	return headerFields, columns, rows, nil
}

// GetHeader provides access to the ALE's header fields while maintaining encapsulation
// of the internal ALEObject structure.
func (ale *ALEObject) GetHeader() []ALEField {
	return ale.HeaderFields
}

// GetColumns provides access to the ALE's column definitions while maintaining encapsulation
// of the internal ALEObject structure.
func (ale *ALEObject) GetColumns() []ALEColumn {
	return ale.Columns
}

// GetRows provides access to the ALE's data rows while maintaining encapsulation
// of the internal ALEObject structure.
func (ale *ALEObject) GetRows() []ALERow {
	return ale.Rows
}

// readTSVData handles the parsing of tab-separated value data
func readTSVData(input string) ([][]string, error) {
	reader := csv.NewReader(strings.NewReader(input))
	reader.Comma = '\t'
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}
	return records, nil
}

// readTSVDataFirstLine uses encoding/csv's Reader but only first line
func readTSVDataFirstLine(input string) ([]string, error) {
	reader := csv.NewReader(strings.NewReader(input))
	reader.Comma = '\t'
	records, err := reader.Read()
	if err != nil {
		return nil, err
	}
	return records, nil
}

// makeALEColumn is a constructor for ALEColumn.
func makeALEColumn(name string, order int) ALEColumn {
	return ALEColumn{Name: name, Order: order}
}

// makeALEValue is a constructor for ALEValueString.
func makeALEValue(column ALEColumn, value string) ALEValueString {
	return ALEValueString{Column: column, Value: value}
}

// makeALERow is a constructor for ALERow.
func makeALERow(row []string, columns []ALEColumn) ALERow {
	var aleRow ALERow
	aleRow.Columns = columns
	aleRow.ValueMap = make(map[ALEColumn]ALEValueString)
	for cell_index, value := range row {
		column := columns[cell_index]
		aleValue := makeALEValue(column, value)
		aleRow.ValueMap[column] = aleValue
	}
	return aleRow
}

// makeALERowsFromDataRows is a constructor for ALERow, iterating over multiple data rows
func makeALERowsFromDataRows(rows [][]string, columns []ALEColumn) []ALERow {
	var aleRows []ALERow
	for row_index, row := range rows {
		aleRow := makeALERow(row, columns)
		aleRow.Order = row_index
		aleRows = append(aleRows, aleRow)
	}
	return aleRows
}
