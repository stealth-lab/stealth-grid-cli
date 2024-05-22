package model

import (
	"strconv"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/simplesmentemat/stealth-grid-cli/pkg/export"
	"github.com/simplesmentemat/stealth-grid-cli/pkg/graphql"
)

type Item struct {
	TitleText       string
	DescriptionText string
	ID              string
}

func (i Item) FilterValue() string { return i.TitleText }

func (i Item) Title() string {
	return i.TitleText
}

func (i Item) Description() string {
	return i.DescriptionText
}

type State int

const (
	SelectGame State = iota
	EnterStartDays
	EnterEndDays
	ShowTable
	SelectSeries
	Downloading
)

type Model struct {
	ListModel    list.Model
	Table        table.Model
	Spinner      spinner.Model
	ErrMsg       string
	CurrentState State
	Loading      bool
	SelectedID   string
	Data         []table.Row
	StartDays    string
	EndDays      string
}

var BaseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

func InitModel(items []list.Item) Model {
	const defaultWidth = 40
	const listHeight = 20
	l := list.New(items, list.NewDefaultDelegate(), defaultWidth, listHeight)
	l.Title = "Select a Game"

	s := spinner.New()
	s.Spinner = spinner.Dot

	return Model{ListModel: l, Spinner: s, CurrentState: SelectGame}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func fetchDataCmd(titleID string, startTime, endTime time.Time) tea.Cmd {
	return func() tea.Msg {
		result, err := graphql.FetchData(titleID, startTime, endTime)
		if err != nil {
			return err.Error()
		}
		return result
	}
}

func downloadDataCmd(seriesID string) tea.Cmd {
	return func() tea.Msg {
		graphql.DownloadJSON(seriesID)
		return "Download complete"
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKeyMsg(msg)

	case map[string]interface{}:
		return m.handleDataMsg(msg)

	case string:
		if msg == "Download complete" {
			m.CurrentState = ShowTable
			m.Loading = false
			return m, tea.Batch(tea.ClearScreen, m.Spinner.Tick)
		} else {
			m.ErrMsg = msg
		}
		return m, nil

	case spinner.TickMsg:
		if m.Loading {
			var cmd tea.Cmd
			m.Spinner, cmd = m.Spinner.Update(msg)
			cmds = append(cmds, cmd)
		}
	}

	if m.CurrentState == ShowTable && !m.Loading {
		var cmd tea.Cmd
		m.Table, cmd = m.Table.Update(msg)
		cmds = append(cmds, cmd)
	} else {
		var cmd tea.Cmd
		m.ListModel, cmd = m.ListModel.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m *Model) handleKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "ctrl+c":
		return m, tea.Quit
	case "enter":
		return m.handleEnterKey()
	case "e":
		if m.CurrentState == ShowTable {
			export.ExportData(m.Data)
		}
	case "backspace":
		return m.handleBackspaceKey()
	case "up", "down":
		if m.CurrentState == SelectGame || m.CurrentState == ShowTable {
			var cmd tea.Cmd
			if m.CurrentState == SelectGame {
				m.ListModel, cmd = m.ListModel.Update(msg)
			} else if m.CurrentState == ShowTable {
				m.Table, cmd = m.Table.Update(msg)
			}
			return m, cmd
		}
	default:
		return m.handleDefaultKey(msg.String())
	}
	return m, nil
}

func (m *Model) handleEnterKey() (tea.Model, tea.Cmd) {
	switch m.CurrentState {
	case SelectGame:
		selectedItem := m.ListModel.SelectedItem().(Item)
		m.SelectedID = selectedItem.ID
		m.CurrentState = EnterStartDays
		return m, nil
	case EnterStartDays:
		m.CurrentState = EnterEndDays
		return m, nil
	case EnterEndDays:
		startDays, _ := strconv.Atoi(m.StartDays)
		endDays, _ := strconv.Atoi(m.EndDays)
		startTime := time.Now().Add(time.Duration(-startDays) * 24 * time.Hour)
		endTime := time.Now().Add(time.Duration(endDays) * 24 * time.Hour)
		m.Loading = true
		m.CurrentState = ShowTable
		return m, tea.Batch(tea.ClearScreen, fetchDataCmd(m.SelectedID, startTime, endTime), m.Spinner.Tick)
	case ShowTable:
		selectedRow := m.Table.SelectedRow()
		m.CurrentState = Downloading
		m.SelectedID = selectedRow[1]
		m.Loading = true
		return m, tea.Batch(tea.ClearScreen, downloadDataCmd(m.SelectedID), m.Spinner.Tick)
	case Downloading:
		m.Loading = false
		m.CurrentState = ShowTable
		return m, tea.ClearScreen
	case SelectSeries:
		m.Loading = true
		return m, tea.Batch(tea.ClearScreen, downloadDataCmd(m.SelectedID), m.Spinner.Tick)
	}
	return m, nil
}

func (m *Model) handleBackspaceKey() (tea.Model, tea.Cmd) {
	if m.CurrentState == EnterStartDays && len(m.StartDays) > 0 {
		m.StartDays = m.StartDays[:len(m.StartDays)-1]
	} else if m.CurrentState == EnterEndDays && len(m.EndDays) > 0 {
		m.EndDays = m.EndDays[:len(m.EndDays)-1]
	}
	return m, nil
}

func (m *Model) handleDefaultKey(key string) (tea.Model, tea.Cmd) {
	if m.CurrentState == EnterStartDays {
		m.StartDays += key
	} else if m.CurrentState == EnterEndDays {
		m.EndDays += key
	}
	return m, nil
}

func (m *Model) handleDataMsg(msg map[string]interface{}) (tea.Model, tea.Cmd) {
	m.ErrMsg = ""
	m.Loading = false
	data, ok := msg["data"].(map[string]interface{})
	if !ok {
		m.ErrMsg = "No data found"
		return m, nil
	}

	series, ok := data["allSeries"].(map[string]interface{})
	if !ok {
		m.ErrMsg = "No series found"
		return m, nil
	}

	edges, ok := series["edges"].([]interface{})
	if !ok {
		m.ErrMsg = "No edges found"
		return m, nil
	}

	var rows []table.Row
	for _, edge := range edges {
		node := edge.(map[string]interface{})["node"].(map[string]interface{})
		tournament := node["tournament"].(map[string]interface{})
		teams := node["teams"].([]interface{})

		if len(teams) < 2 {
			continue
		}

		team1 := teams[0].(map[string]interface{})["baseInfo"].(map[string]interface{})["name"].(string)
		team2 := teams[1].(map[string]interface{})["baseInfo"].(map[string]interface{})["name"].(string)

		row := table.Row{
			node["startTimeScheduled"].(string),
			node["id"].(string),
			tournament["name"].(string),
			team1,
			team2,
		}
		rows = append(rows, row)
	}

	columns := []table.Column{
		{Title: "Start Time", Width: 20},
		{Title: "Serie ID", Width: 10},
		{Title: "Tournament", Width: 20},
		{Title: "Team One", Width: 20},
		{Title: "Team Two", Width: 20},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(15),
		table.WithWidth(100),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	m.Table = t
	m.Data = rows
	return m, nil
}

func (m Model) View() string {
	if m.ErrMsg != "" {
		return m.ErrMsg
	}
	switch m.CurrentState {
	case SelectGame:
		return BaseStyle.Render(m.ListModel.View())
	case EnterStartDays:
		return BaseStyle.Render("Enter the number of past days to include (e.g., 10): " + m.StartDays)
	case EnterEndDays:
		return BaseStyle.Render("Enter the number of future days to include (e.g., 1): " + m.EndDays)
	case ShowTable:
		if m.Loading {
			return BaseStyle.Render(m.Spinner.View()) + "\nLoading data, please wait..."
		}
		return BaseStyle.Render(m.Table.View()) + "\nPress 'e' to export data, or press Enter to select a series."
	case Downloading:
		if m.Loading {
			return BaseStyle.Render(m.Spinner.View()) + "\nDownloading data, please wait..."
		}
		return BaseStyle.Render(m.Table.View())
	case SelectSeries:
		return BaseStyle.Render("Press Enter to download the selected series.")
	}
	return ""
}
