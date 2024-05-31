package tui

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/simplesmentemat/stealth-grid-cli/pkg/model"
)

// InitModel initializes the application model with a list of items.
//
// This function serves as a wrapper around the model.InitModel function, providing a convenient
// way to initialize the application's user interface model with the specified items.
//
// Parameters:
//   - items: A slice of list.Item representing the items to be displayed in the list.
//
// Returns:
//   - model.Model: The initialized application model with the provided list items.
func InitModel(items []list.Item) model.Model {
	return model.InitModel(items)
}
