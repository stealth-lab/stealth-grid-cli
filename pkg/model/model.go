// Package model provides the main application model and its functionalities.
// It includes data structures, commands, and functions to manage the application state
// and handle user interactions.
package model

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"time"
	"unicode"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/simplesmentemat/stealth-grid-cli/pkg/export"
	"github.com/simplesmentemat/stealth-grid-cli/pkg/graphql"
	"github.com/sqweek/dialog"
)

// Item represents a list item with a title, description, and ID.
// It is used in the application's list model to display selectable items.
type Item struct {
	TitleText       string // TitleText is the title of the item.
	DescriptionText string // DescriptionText is the description of the item.
	ID              string // ID is the unique identifier of the item.
}

// FilterValue returns the title text for filtering purposes.
//
// This method implements the list.Item interface.
func (i Item) FilterValue() string { return i.TitleText }

// Title returns the title text of the item.
//
// This method implements the list.Item interface.
func (i Item) Title() string { return i.TitleText }

// Description returns the description text of the item.
//
// This method implements the list.Item interface.
func (i Item) Description() string { return i.DescriptionText }

// State represents the different states of the application.
// It is used to manage and track the current state of the application.
type State int

const (
	// SelectGame indicates that the application is in the state where the user selects a game.
	SelectGame State = iota

	// EnterStartDays indicates that the application is in the state where the user enters the start days.
	EnterStartDays

	// EnterEndDays indicates that the application is in the state where the user enters the end days.
	EnterEndDays

	// ShowTable indicates that the application is in the state where the table with series data is displayed.
	ShowTable

	// SelectSeries indicates that the application is in the state where the user selects a series.
	SelectSeries

	// Downloading indicates that the application is in the state where data is being downloaded.
	Downloading
	SelectDownloadOption
)

// Model represents the main application model.
type Model struct {
	ListModel         list.Model
	Table             table.Model
	Spinner           spinner.Model
	ErrMsg            string
	CurrentState      State
	Loading           bool
	SelectedID        string
	Data              []table.Row
	StartDays         string
	EndDays           string
	DownloadOption    string
	DownloadOptions   []list.Item
	DownloadListModel list.Model
}

// BaseStyle defines the base style for the application.
//
// BaseStyle sets the border style to normal and the border color to a shade of gray (color 240).
var BaseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

// InitModel initializes the application model with a list of items.
//
// This function sets up the list of items, configures the spinner,
// and sets the initial state of the application to SelectGame.
//
// Parameters:
//   - items: A slice of list.Item representing the items to be displayed in the list.
//
// Returns:
//   - Model: A Model struct initialized with the provided list items, a configured spinner,
//     and the initial application state set to SelectGame.
func InitModel(items []list.Item) Model {
	const (
		defaultWidth = 40 // defaultWidth is the default width of the list.
		listHeight   = 20 // listHeight is the height of the list.
	)

	// Initialize the list with the provided items, a default delegate, and specified dimensions.
	l := list.New(items, list.NewDefaultDelegate(), defaultWidth, listHeight)
	l.Title = "Stealth Grid - Select a Game"

	// Initialize and configure the spinner.
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	options := []list.Item{}

	dl := list.New(options, list.NewDefaultDelegate(), defaultWidth, listHeight)
	dl.Title = "Select Download Option"

	return Model{
		ListModel:         l,
		Spinner:           s,
		CurrentState:      SelectGame,
		DownloadOptions:   options,
		DownloadListModel: dl,
	}
}

// Init initializes the application.
//
// This function sets up the initial command for the application,
// which in this case is the spinner tick command to start the spinner animation.
//
// Returns:
//   - tea.Cmd: A command that starts the spinner tick command.
func (m Model) Init() tea.Cmd {
	return m.Spinner.Tick
}

// fetchDataCmd fetches data for the specified title ID within the given time range.
//
// This function creates a command that fetches data from a GraphQL API for a specified
// title ID and time range. It returns the result as a tea.Msg. If an error occurs during
// the data fetch, the error message is returned.
//
// Parameters:
//   - titleID: A string representing the ID of the title to query for.
//   - startTime: A time.Time object representing the start time of the query range.
//   - endTime: A time.Time object representing the end time of the query range.
//
// Returns:
//   - tea.Cmd: A command that fetches the data and returns a tea.Msg containing the result or an error message.
func fetchDataCmd(titleID string, startTime, endTime time.Time) tea.Cmd {
	return func() tea.Msg {
		result, err := graphql.FetchData(titleID, startTime, endTime)
		if err != nil {
			return err.Error()
		}
		return result
	}
}

// downloadDataCmd downloads data for the specified series ID to the specified directory.
//
// This function creates a command that downloads a ZIP file containing data for a specified
// series ID and saves it to the given directory. It returns a message indicating the download
// status.
//
// Parameters:
//   - seriesID: A string representing the ID of the series to download the data for.
//   - directory: A string representing the directory where the ZIP file will be saved.
//
// Returns:
//   - tea.Cmd: A command that downloads the data and returns a tea.Msg indicating the download status.
func downloadDataCmd(seriesID string, option string) tea.Cmd {
	return func() tea.Msg {
		directory, err := dialog.Directory().Title("Select Download Directory").Browse()
		if err != nil || directory == "" {
			return "Download cancelled or directory not selected"
		}

		if option == "events-grid-compressed" {
			err := graphql.DownloadJSON(seriesID, directory)
			if err != nil {
				return fmt.Sprintf("Error downloading JSON: %v", err)
			}
		} else {
			err := graphql.DownloadGame(seriesID, option, directory)
			if err != nil {
				return fmt.Sprintf("Error downloading ROFL for game %s: %v", option, err)
			}
		}

		return "Download complete"
	}
}

// Update handles messages and updates the application state.
//
// This function processes incoming messages and updates the application state accordingly.
// It handles various message types such as key presses, data fetch results, and download status.
//
// Parameters:
//   - msg: A tea.Msg representing the incoming message to be handled.
//
// Returns:
//   - tea.Model: The updated model.
//   - tea.Cmd: A command to be executed, if any.
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
		} else if msg != "" {
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

// handleKeyMsg handles key messages.
//
// This function processes keyboard input messages and updates the application state accordingly.
// It handles various key presses such as quitting the application, confirming selections,
// exporting data, navigating the list or table, and other default key actions.
//
// Parameters:
//   - msg: A tea.KeyMsg representing the key message to be handled.
//
// Returns:
//   - tea.Model: The updated model.
//   - tea.Cmd: A command to be executed, if any.
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
		return m, tea.ClearScreen
	case "backspace":
		return m.handleBackspaceKey()
	case "up", "down":
		if m.CurrentState == SelectGame || m.CurrentState == ShowTable || m.CurrentState == SelectDownloadOption {
			var cmd tea.Cmd
			if m.CurrentState == SelectGame {
				m.ListModel, cmd = m.ListModel.Update(msg)
			} else if m.CurrentState == ShowTable {
				m.Table, cmd = m.Table.Update(msg)
			} else if m.CurrentState == SelectDownloadOption {
				m.DownloadListModel, cmd = m.DownloadListModel.Update(msg)
			}
			return m, cmd
		}
	default:
		return m.handleDefaultKey(msg.String())
	}
	return m, nil
}

// handleEnterKey handles the 'enter' key press.
//
// This function processes the 'enter' key press based on the current state of the application.
// It performs actions such as selecting a game, confirming date ranges, and more.
//
// Returns:
//   - tea.Model: The updated model.
//   - tea.Cmd: A command to be executed, if any.
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
		m.CurrentState = SelectDownloadOption
		m.SelectedID = selectedRow[1]

		roflCount, hasJSON, err := graphql.FetchGameList(m.SelectedID)
		if err != nil {
			m.ErrMsg = fmt.Sprintf("Error fetching game list: %v", err)
			return m, nil
		}

		file, err := os.Create("output.txt")
		if err != nil {
			fmt.Println("Error creating file:", err)
		}
		defer file.Close()

		_, err = file.WriteString(fmt.Sprintf("roflCount: %d\nhasJSON: %v\n", roflCount, hasJSON))
		if err != nil {
			fmt.Println("Error writing to file:", err)
		}

		var options []list.Item
		if hasJSON {
			options = append(options, Item{TitleText: "Download JSON", ID: "events-grid-compressed"})
		}
		for i := 1; i <= roflCount; i++ {
			options = append(options, Item{TitleText: fmt.Sprintf("Download Game %d", i), ID: strconv.Itoa(i)})
		}
		m.DownloadOptions = options
		m.DownloadListModel.SetItems(options)

		return m, nil
	case SelectDownloadOption:
		selectedOption := m.DownloadListModel.SelectedItem().(Item)
		m.DownloadOption = selectedOption.ID
		m.CurrentState = Downloading
		m.Loading = true
		return m, tea.Batch(tea.ClearScreen, downloadDataCmd(m.SelectedID, m.DownloadOption), m.Spinner.Tick)
	case Downloading:
		m.Loading = false
		m.CurrentState = ShowTable
		return m, tea.ClearScreen
	case SelectSeries:
		m.Loading = true
		directory, err := dialog.Directory().Title("Select Download Directory").Browse()
		if err != nil || directory == "" {
			m.Loading = false
			m.CurrentState = ShowTable
			return m, tea.ClearScreen
		}

		return m, tea.Batch(tea.ClearScreen, downloadDataCmd(m.SelectedID, m.DownloadOption), m.Spinner.Tick)
	}
	return m, nil
}

// handleBackspaceKey handles the 'backspace' key press.
//
// This function processes the 'backspace' key press to delete the last character
// in the StartDays or EndDays fields based on the current state of the application.
//
// Returns:
//   - tea.Model: The updated model.
//   - tea.Cmd: A command to be executed, if any.
func (m *Model) handleBackspaceKey() (tea.Model, tea.Cmd) {
	if m.CurrentState == EnterStartDays && len(m.StartDays) > 0 {
		m.StartDays = m.StartDays[:len(m.StartDays)-1]
	} else if m.CurrentState == EnterEndDays && len(m.EndDays) > 0 {
		m.EndDays = m.EndDays[:len(m.EndDays)-1]
	}
	return m, nil
}

// handleDefaultKey handles default key actions.
//
// This function processes other key presses that are not explicitly handled by the previous cases.
// It performs default actions based on the specific key pressed.
//
// Parameters:
//   - key: A string representing the key that was pressed.
//
// Returns:
//   - tea.Model: The updated model.
//   - tea.Cmd: A command to be executed, if any.
func (m *Model) handleDefaultKey(key string) (tea.Model, tea.Cmd) {
	if !unicode.IsDigit([]rune(key)[0]) {
		return m, nil
	}

	if m.CurrentState == EnterStartDays {
		m.StartDays += key
	} else if m.CurrentState == EnterEndDays {
		m.EndDays += key
	}
	return m, nil
}

// handleDataMsg handles data messages.
//
// This function processes incoming data messages, updates the application state with the
// retrieved data, and constructs a table to display the series information. It performs
// several checks to ensure the data is valid and sets error messages if any issues are found.
//
// Parameters:
//   - msg: A map[string]interface{} representing the data message to be handled.
//
// Returns:
//   - tea.Model: The updated model.
//   - tea.Cmd: A command to be executed, if any.
//
// Processing Steps:
//  1. Clear any existing error message and set loading to false.
//  2. Extract the data from the message. If the data is not found, set an error message.
//  3. Extract the series data from the data map. If the series data is not found, set an error message.
//  4. Extract the edges array from the series map. If the edges array is not found, set an error message.
//  5. Iterate over the edges array to extract relevant data for each series, including the start time, series ID,
//     tournament name, and team names. Ensure there are at least two teams.
//  6. Construct table rows from the extracted data and append them to the rows slice.
//  7. Sort the rows by start time in ascending order.
//  8. Define the table columns.
//  9. Create a new table with the specified columns, rows, and styles.
// 10. Define the table styles for the headers and selected rows.
// 11. Update the model with the new table and data.
// 12. Return the updated model and no additional command.
//
// Error Handling:
//   - The function includes checks to ensure the data is valid at each step, setting error messages if any issues are found.
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

	sort.SliceStable(rows, func(i, j int) bool {
		timeI, _ := time.Parse(time.RFC3339, rows[i][0])
		timeJ, _ := time.Parse(time.RFC3339, rows[j][0])
		return timeI.Before(timeJ)
	})

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

// View returns the current view of the application.
//
// This function constructs and returns the string representation of the current view based on the application's
// state. It handles various states such as selecting a game, entering start and end days, showing the table,
// and downloading data. If there is an error message, it returns the error message.
//
// Returns:
//   - string: The current view of the application.
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
			return BaseStyle.Render(fmt.Sprintf("\n\n   %s Loading data, please wait...  \n\n", m.Spinner.View()))
		}
		return BaseStyle.Render(m.Table.View()) + "\nPress 'e' to export data, or press Enter to select a series."
	case SelectDownloadOption:
		return BaseStyle.Render(m.DownloadListModel.View())
	case Downloading:
		if m.Loading {
			return BaseStyle.Render(fmt.Sprintf("\n\n   %s Downloading data, please wait...  \n\n", m.Spinner.View()))
		}
		return BaseStyle.Render(m.Table.View())
	case SelectSeries:
		return BaseStyle.Render("Press Enter to download the selected series.")
	}
	return ""
}
