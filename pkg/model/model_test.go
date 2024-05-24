package model

import (
	"testing"
	"time"

	"github.com/charmbracelet/bubbles/list"
)

func TestInitModel(t *testing.T) {
	items := []list.Item{
		Item{TitleText: "Game 1", DescriptionText: "Description 1", ID: "1"},
		Item{TitleText: "Game 2", DescriptionText: "Description 2", ID: "2"},
	}
	model := InitModel(items)
	if model.ListModel.Title != "Stealth Grid - Select a Game" {
		t.Fatalf("Expected title to be 'Stealth Grid - Select a Game', but got %s", model.ListModel.Title)
	}
}

func TestFetchDataCmd(t *testing.T) {
	startTime := time.Now().Add(-24 * time.Hour)
	endTime := time.Now()
	cmd := fetchDataCmd("3", startTime, endTime)
	if cmd == nil {
		t.Fatalf("Expected fetchDataCmd to return a non-nil command")
	}
}

func TestDownloadDataCmd(t *testing.T) {
	cmd := downloadDataCmd("3", "/tmp")
	if cmd == nil {
		t.Fatalf("Expected downloadDataCmd to return a non-nil command")
	}
}
