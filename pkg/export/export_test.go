package export

import (
	"os"
	"testing"

	"github.com/charmbracelet/bubbles/table"
)

func TestExportData(t *testing.T) {
	rows := []table.Row{
		{"2024-05-10T00:00:00Z", "1", "Tournament 1", "Team 1", "Team 2"},
	}
	ExportData(rows)
	if _, err := os.Stat("games.csv"); os.IsNotExist(err) {
		t.Fatalf("Expected file games.csv to be created, but it does not exist")
	}
}
