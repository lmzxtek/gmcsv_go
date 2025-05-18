package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
)

func main() {
	// Open the CSV file
	// csvFile, err := os.Open("data.csv")
	csvFile, err := os.Open("https://gmc.zwdk.im/download/test.txt")
	if err != nil {
		fmt.Println("Error opening CSV file:", err)
		return
	}
	defer csvFile.Close()

	// Create a new CSV reader
	csvReader := csv.NewReader(csvFile)

	// Loop through each row in the CSV file
	for {
		// Read each row from the CSV file
		row, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("Error reading CSV file:", err)
			return
		}

		// Print the row
		fmt.Println(row)
	}
}
