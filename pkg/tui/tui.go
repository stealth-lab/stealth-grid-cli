package tui

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/simplesmentemat/stealth-grid-cli/pkg/model"
)

func InitModel(items []list.Item) model.Model {
	return model.InitModel(items)
}
