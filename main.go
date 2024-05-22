package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/simplesmentemat/stealth-grid-cli/pkg/config"
	"github.com/simplesmentemat/stealth-grid-cli/pkg/model"
	"github.com/simplesmentemat/stealth-grid-cli/pkg/tui"
)

func main() {
	err := config.InitConfig()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	items := []list.Item{
		model.Item{TitleText: "League of Legends", DescriptionText: "ID: 3", ID: "3"},
		model.Item{TitleText: "Valorant", DescriptionText: "ID: 6", ID: "6"},
		model.Item{TitleText: "CS 2", DescriptionText: "ID: 28", ID: "28"},
	}

	p := tea.NewProgram(tui.InitModel(items), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		panic(err)
	}
}
