package export

import (
	"encoding/csv"
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/table"
)

// ExportData exports the provided data to a CSV file named "games.csv".
// The CSV file will include the headers "Start Time", "Serie ID", "Tournament", "Blue Team", and "Red Team".
// Each row in the provided data will be written as a record in the CSV file.
//
// Parameters:
//
//	data: []table.Row - A slice of table rows containing the data to be exported. Each row is expected to have 5 elements:
//	  - Start Time (string): The start time of the game
//	  - Serie ID (string): The unique identifier of the series
//	  - Tournament (string): The name of the tournament
//	  - Blue Team (string): The name of the blue team
//	  - Red Team (string): The name of the red team
//
// Behavior:
//
//	The function creates a new CSV file named "games.csv" in the current working directory. It writes the headers to the file,
//	followed by each row of data. If an error occurs during file creation or writing, the function prints an error message
//	to the console. Upon successful completion, a confirmation message is printed and the function pauses for 1 second.
func ExportData(data []table.Row) {
	fileName := "games.csv"
	file, err := os.Create(fileName)
	if err != nil {
		fmt.Printf("Error creating file: %v", err)
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	headers := []string{"Start Time", "Serie ID", "Tournament", "Blue Team", "Red Team"}
	if err := writer.Write(headers); err != nil {
		fmt.Printf("Error writing headers to CSV: %v", err)
		return
	}

	for _, row := range data {
		record := []string{row[0], row[1], row[2], row[3], row[4]}
		if err := writer.Write(record); err != nil {
			fmt.Printf("Error writing record to CSV: %v", err)
			return
		}
	}
	fmt.Printf("Data successfully exported to %s\n", fileName)
	time.Sleep(1 * time.Second)
}
