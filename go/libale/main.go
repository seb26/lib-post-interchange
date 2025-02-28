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
	Read()
	ReadBytes()
	Write()
	WriteBytes()

	Rows() []ALERow
}

// ReadFile() takes an input filepath, reads it and calls Read() to return an ALE object
func ReadFile(filepath string) (ALE, error) {
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

func Read(input string) (ALE, error) {
	pattern := regexp.MustCompile(ALEHeadingWordPattern)
	match := pattern.FindStringSubmatchIndex(input)
	if match != nil {
		//fields := input[match[2]:match[3]]
		//columns := input[match[10]:match[11]]
		data_first_byte := match[12]
		_, err := readTSVData(input[data_first_byte:])
		if err != nil {
			return nil, err
		}
		return nil, err
	} else {
		return nil, fmt.Errorf("unrecognised ALE header")
		// PRINT DEBUG 1kB of the header
	}
}

// readCSVData() takes a string input and returns a 2D slice of strings
func readTSVData(input string) ([][]string, error) {
	reader := csv.NewReader(strings.NewReader(input))
	reader.Comma = '\t'
	return reader.ReadAll()
}
