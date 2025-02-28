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
func ReadFile(filepath string) (int, error) {
	_, err := os.Stat(filepath)
	if os.IsNotExist(err) {
		return 0, err
	}
	data, err := os.ReadFile(filepath)
	if err != nil {
		return 0, err
	}
	dataString := string(data)
	return Read(dataString)
}

func Read(input string) (int, error) {
	_, _, _, err := read(input)
	if err != nil {
		return 0, err
	}
	//ale := &ALEObject{
	//	HeaderFields: aleHeaderFields,
	//	Columns:      aleColumns,
	//	Rows:         aleRows,
	//}
	//newAle := ALE(ale)
	return 0, nil
}

func read(input string) (*[]ALEHeaderField, *[]ALEColumn, *[]ALERow, error) {
	var columns []ALEColumn

	// Find the ALE header
	pattern := regexp.MustCompile(ALEHeadingWordPattern)
	match := pattern.FindStringSubmatchIndex(input)

	if match == nil {
		return nil, nil, nil, fmt.Errorf("unrecognised ALE header")
	}

	// Find the index positions for 'fields' and 'columns'
	fieldsId := pattern.SubexpIndex("fields")
	columnsId := pattern.SubexpIndex("columns")

	fieldsStart, fieldsEnd := match[2*fieldsId], match[2*fieldsId+1]
	fieldsList := input[fieldsStart:fieldsEnd]
	columnsStart, columnsEnd := match[2*columnsId], match[2*columnsId+1]
	columnsList := input[columnsStart:columnsEnd]

	fieldsArray, err := readTSVData(fieldsList)
	if err != nil {
		return nil, nil, nil, err
	}
	var fieldsMap = make(map[string]string)
	// Map these header fields as plain strings
	for _, field := range fieldsArray {
		fieldsMap[field[0]] = field[1]
	}
	// Interpret
	var headerFields []ALEHeaderField
	for key, value := range fieldsMap {
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
		columns = append(columns, *makeALEColumn(column, index))
	}
	fmt.Println("Columns:", columns)
	return &headerFields, &columns, nil, nil
}

// GetHeader() returns the header fields of the ALE object
func (ale *ALEObject) GetHeader() []ALEHeaderField {
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
func makeALEColumn(input string, order int) *ALEColumn {
	return &ALEColumn{Name: input, Order: order}
}

// makeALEValue() takes a string input and returns an ALEBaseValue object
func makeALEValue(input string) ALEValueString {
	return ALEValueString{Value: input}
}

// makeALERow() takes a string input and returns an ALERow object
