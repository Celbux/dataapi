package tools

import (
	"encoding/csv"
	"os"
	"strings"
)

/*
 * ReadFile(file, method)
 * Default Method: ReadModeSingle (Single is stored in the 0th element)
 * ReadModeSingle returns a single string, including all \n
 * ReadModeSingleCollapsed returns a single string, with all \n stripped away
 * ReadModeMultiline returns an array of string, split on \n (each line)
 */
const (
	ReadModeSingle = iota
	ReadModeSingleCollapsed
	ReadModeMultiline
)
func ReadFile(file os.File, method... int) ([]string, error) {

	// Read the input file into a string
	lines := csv.NewReader(&file)
	lines.FieldsPerRecord = -1
	rawData, err := lines.ReadAll()
	if err != nil {
		return nil, err
	}
	defer file.Close()
	data := ""
	for _, row := range rawData {
		data += strings.Join(row, ",") + "\n"
	}

	// Return the data in the format specified
	if method == nil || method[0] == ReadModeSingle {
		return []string {data}, nil
	}
	if method[0] == ReadModeSingleCollapsed {
		return []string { strings.Replace(data, "\n", "", -1) }, nil
	}
	if method[0] == ReadModeMultiline {
		return strings.Split(data, "\n"), nil
	}

	// Default return is ReadModeSingle but this will never get called
	return []string {data}, nil

}
