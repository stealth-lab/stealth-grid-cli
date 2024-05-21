package export

import (
	"encoding/csv"
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/table"
)

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
}
