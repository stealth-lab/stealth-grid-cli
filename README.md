# Stealth GRID CLI

Welcome to the **Stealth GRID CLI**! This tool helps you fetch, display, and export data related to esports tournaments and series from the GRID API.

![License](https://img.shields.io/badge/license-MIT-blue.svg)
![Version](https://img.shields.io/badge/version-0.0.1-brightgreen.svg)

## Behind Stealth GRID CLI

- **GRID**: The leading esports data platform with exclusive global rights to distribute official data and video streams for select events and tournaments.
- **Stealth**: Utilizes GRID data to address key questions in esports by analyzing large volumes of data in real-time, creating agile and adaptable models.
- **CLI**: Initially developed for internal use, this CLI has been open-sourced to assist developers working with GRID data.

## Usage

Start the CLI by opening your terminal and running:
```sh
stealth
```

## Features
- **Game Selection**: Choose from a list of available games to view their series data.
- **Date Range Filtering**: Filter series data by specifying start and end days.
- **Data Display**: View series data in a table with columns for Start Time, Series ID, Tournament, Team One, and Team Two.
- **Data Export**: Export displayed data to a CSV file at a user-specified location.
- **Data Download**: Download detailed data for a selected series as a ZIP file to a user-specified directory.
- **Interactive UI**: Navigate through the application using keyboard controls for an interactive experience.

## Main Menu
1. **Select Game**: Use the arrow keys to navigate through the list of games and press `Enter` to select a game.
2. **Enter Start Days**: Enter the number of past days to include (e.g., 10) and press `Enter`.
3. **Enter End Days**: Enter the number of future days to include (e.g., 1) and press `Enter`.
4. **Show Table**: The series data will be displayed in a table format. Use the arrow keys to navigate through the table.

## Key Controls
- `q` or `Ctrl+C`: Quit the application.
- `Enter`: Confirm selection or proceed to the next step.
- `e`: Export data to CSV.
- `Backspace`: Delete the last character when entering start or end days.
- `Up/Down Arrow`: Navigate through lists and tables.

## Export Data
While viewing the table, press `e` to export the displayed data to a CSV file. A dialog will prompt you to select the location and filename for the CSV file.

## Download Data
To download detailed data for a selected series:

1. Navigate to the desired row in the table.
2. Press `Enter` to select the series.
3. A dialog will prompt you to select the directory where the ZIP file will be saved.

## Contributing
We welcome contributions! Please fork the repository and submit pull requests with your changes. For major changes, please open an issue first to discuss what you would like to change.

1. Fork the repository.
2. Create your feature branch (`git checkout -b feature/YourFeature`).
3. Commit your changes (`git commit -m 'Add some feature'`).
4. Push to the branch (`git push origin feature/YourFeature`).
5. Open a pull request.

## Todo
- [ ] MSI installer.

## License
This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
