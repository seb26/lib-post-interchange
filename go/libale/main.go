package libale

import (
	"encoding/csv"
	"fmt"
	"os"
	"regexp"
	"strings"
)

// Read
// Write

const ALE_HEADER_READ_FIRST_BYTES = 1024

type ALE interface {
	ALEObject
}

// ReadFile() takes an input filepath, reads it and calls Read() to return an ALE object
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

func read(input string) ([]ALEField, []ALEColumn, []ALERow, error) {
	var columns []ALEColumn
	var headerFields []ALEField

	// Find the ALE header
	pattern := regexp.MustCompile(ALEHeadingWordPattern)
	match := pattern.FindStringSubmatchIndex(input)
	if match == nil {
		return nil, nil, nil, fmt.Errorf("unrecognised ALE header")
	}
	// Find the index positions for 'fields' and 'columns'
	fieldsId := pattern.SubexpIndex("fields")
	columnsId := pattern.SubexpIndex("columns")
	dataHeaderId := pattern.SubexpIndex("data_header")
	fieldsStart, fieldsEnd := match[2*fieldsId], match[2*fieldsId+1]
	fieldsList := input[fieldsStart:fieldsEnd]
	columnsStart, columnsEnd := match[2*columnsId], match[2*columnsId+1]
	columnsList := input[columnsStart:columnsEnd]
	dataHeaderEnd := match[2*dataHeaderId+1]
	fieldsArray, err := readTSVData(fieldsList)
	if err != nil {
		return nil, nil, nil, err
	}
	// Interpret an array per line of fields, where [key, value]
	for _, field := range fieldsArray {
		key := field[0]
		value := field[1]
		constructor, err := ToType(key)
		if err != nil {
			return nil, nil, nil, err
		}
		headerFields = append(headerFields, constructor(value))
	}
	// Parse columns
	columnsArray, err := readTSVDataFirstLine(columnsList)
	if err != nil {
		return nil, nil, nil, err
	}
	for index, column := range columnsArray {
		columns = append(columns, makeALEColumn(column, index))
	}
	// Parse data
	data := strings.TrimSpace(input[dataHeaderEnd:])
	dataRows, err := readTSVData(data)
	if err != nil {
		return nil, nil, nil, err
	}
	rows := makeALERowsFromDataRows(dataRows, columns)
	return headerFields, columns, rows, nil
}

// GetHeader() returns the header fields of the ALE object
func (ale *ALEObject) GetHeader() []ALEField {
	return ale.HeaderFields
}

// GetColumns() returns the columns of the ALE object
func (ale *ALEObject) GetColumns() []ALEColumn {
	return ale.Columns
}

// GetRows() returns the rows of the ALE object
func (ale *ALEObject) GetRows() []ALERow {
	return ale.Rows
}

// readTSVData() takes a string input and returns a 2D array slice of strings
func readTSVData(input string) ([][]string, error) {
	reader := csv.NewReader(strings.NewReader(input))
	reader.Comma = '\t'
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}
	return records, nil
}

// readTSVDataFirstLine() takes a string input and returns a slice of strings
func readTSVDataFirstLine(input string) ([]string, error) {
	reader := csv.NewReader(strings.NewReader(input))
	reader.Comma = '\t'
	records, err := reader.Read()
	if err != nil {
		return nil, err
	}
	return records, nil
}

// makeALEColumn() takes a string input and returns an ALEColumn object
func makeALEColumn(name string, order int) ALEColumn {
	return ALEColumn{Name: name, Order: order}
}

// makeALEValue() takes a string input and returns an ALEBaseValue object
func makeALEValue(column ALEColumn, value string) ALEValueString {
	return ALEValueString{Column: column, Value: value}
}

// makeALERow() takes an array of row values, and an array of columns, and returns an ALERow object
func makeALERow(row []string, columns []ALEColumn) ALERow {
	var aleRow ALERow
	aleRow.Columns = columns
	aleRow.ValueMap = make(map[ALEColumn]ALEValueString) // Initialize the map
	for cell_index, value := range row {
		column := columns[cell_index]
		aleValue := makeALEValue(column, value)
		aleRow.ValueMap[column] = aleValue
	}
	return aleRow
}

// makeALERowsFromData() takes a 2D array of strings and returns a slice of ALERow objects
func makeALERowsFromDataRows(rows [][]string, columns []ALEColumn) []ALERow {
	var aleRows []ALERow
	for row_index, row := range rows {
		aleRow := makeALERow(row, columns)
		aleRow.Order = row_index
		aleRows = append(aleRows, aleRow)
	}
	return aleRows
}
