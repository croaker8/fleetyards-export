package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// getFieldList - Get list of fields to export
func getFieldList(filePath string) ([]string, error) {

	list := make([]string, 0, 50)

	// open the file
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Error opening export field list file '%s': %s\n", filePath, err)
		return nil, err
	}
	defer file.Close()

	// read lines ignoring empty or comment lines
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if len(line) > 0 && line[0:1] != "#" {
			list = append(list, line)
		}
	}

	// check for error in scanning
	if err := scanner.Err(); err != nil {
		fmt.Printf("Error scanning lines from export field list file '%s': %s\n", filePath, err)
		return nil, err
	}

	return list, nil
}
