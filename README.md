
# Behind Stealth GRID CLI

- GRID is the leading esports data platform with exclusive global rights to distribute official data and video streams for select events and tournaments.
- Stealth leverages this data to answer key questions in esports by analyzing large volumes in real-time to create agile, adaptable models.
- Originally developed for internal use, the CLI has been open-sourced to assist developers in accessing and utilizing this data.

## Features

- **Game Selection**: Choose from a list of available games to view their series data.
- **Date Range Filtering**: Specify the start and end days to filter the series data within a particular date range.
- **Data Display**: View the series data in a tabular format with columns for Start Time, Series ID, Tournament, Team One, and Team Two.
- **Data Export**: Export the displayed data to a CSV file at a user-specified location.
- **Data Download**: Download detailed data for a selected series as a ZIP file to a user-specified directory.
- **Interactive UI**: Navigate through the application using keyboard controls for an interactive experience.

## Usage

Start the CLI by opening your terminal and running:

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

## Contributing

Contributions are welcome! Please fork the repository and submit pull requests with your changes.

## Todo

- [ ] MSI installer.

## License
This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

---
