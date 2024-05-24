
# Stealth Grid CLI

Stealth Grid CLI is a command-line interface (CLI) tool designed to fetch, display, and export data related to esports tournaments and series from the Grid API. This tool allows users to select a game, specify a date range, view the resulting data in a table, and export the data to a CSV file. Additionally, users can download related data as a ZIP file.

## Features

- **Game Selection**: Choose from a list of available games to view their series data.
- **Date Range Filtering**: Specify the start and end days to filter the series data within a particular date range.
- **Data Display**: View the series data in a tabular format with columns for Start Time, Series ID, Tournament, Blue Team, and Red Team.
- **Data Export**: Export the displayed data to a CSV file at a user-specified location.
- **Data Download**: Download detailed data for a selected series as a ZIP file to a user-specified directory.
- **Interactive UI**: Navigate through the application using keyboard controls for an interactive experience.

## Installation

Download the MSI installer from the [releases](https://github.com/simplesmentemat/stealth-grid-cli/releases) page and run it to install the Stealth Grid CLI on your system.

## Usage

Once installed, you can start the Stealth Grid CLI by opening your terminal and running the following command:

```sh
stealth
```

### Main Menu

1. **Select Game**: Use the arrow keys to navigate through the list of games and press `Enter` to select a game.
2. **Enter Start Days**: Enter the number of past days to include (e.g., 10) and press `Enter`.
3. **Enter End Days**: Enter the number of future days to include (e.g., 1) and press `Enter`.
4. **Show Table**: The series data will be displayed in a table format. Use the arrow keys to navigate through the table.

### Export Data

While viewing the table, press `e` to export the displayed data to a CSV file. A dialog will prompt you to select the location and filename for the CSV file.

### Download Data

To download detailed data for a selected series:
1. Navigate to the desired row in the table.
2. Press `Enter` to select the series.
3. A dialog will prompt you to select the directory where the ZIP file will be saved.

### Key Controls

- `q` or `Ctrl+C`: Quit the application.
- `Enter`: Confirm selection or proceed to the next step.
- `e`: Export data to CSV.
- `Backspace`: Delete the last character when entering start or end days.
- `Up/Down Arrow`: Navigate through lists and tables.

## Dependencies

The Stealth Grid CLI utilizes several Go packages to provide its functionality:
- `github.com/charmbracelet/bubbles/list`
- `github.com/charmbracelet/bubbles/spinner`
- `github.com/charmbracelet/bubbles/table`
- `github.com/charmbracelet/bubbletea`
- `github.com/charmbracelet/lipgloss`
- `github.com/sqweek/dialog`

## Contributing

Contributions are welcome! Please fork the repository and submit pull requests with your changes.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

---