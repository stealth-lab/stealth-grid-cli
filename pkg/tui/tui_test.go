package tui

import (
	"testing"

	"github.com/charmbracelet/bubbles/list"
	"github.com/simplesmentemat/stealth-grid-cli/pkg/model"
)

func TestInitModel(t *testing.T) {
	items := []list.Item{
		model.Item{TitleText: "Game 1", DescriptionText: "Description 1", ID: "1"},
		model.Item{TitleText: "Game 2", DescriptionText: "Description 2", ID: "2"},
	}
	model := InitModel(items)
	if model.ListModel.Title != "Stealth Grid - Select a Game" {
		t.Fatalf("Expected title to be 'Stealth Grid - Select a Game', but got %s", model.ListModel.Title)
	}
}
