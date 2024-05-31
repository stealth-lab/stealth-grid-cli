// Package config provides functionality for managing the configuration of the Stealth Grid CLI application.
//
// This package handles reading, writing, and initializing configuration files,
// including retrieving the API key necessary for accessing the Grid API.
package config

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

var APIURL = "https://api.grid.gg"

// getConfigPath returns the path to the configuration file.
//
// It creates the necessary directories if they do not exist.
//
// Returns:
//   - string: The path to the configuration file.
//   - error: An error if there is any issue determining the user's home directory
//     or creating the configuration directory.
func getConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	configDir := filepath.Join(home, ".config", "stealth-grid-cli")
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		err := os.MkdirAll(configDir, 0755)
		if err != nil {
			return "", err
		}
	}
	return filepath.Join(configDir, "config.yaml"), nil
}

// InitConfig initializes the configuration by reading from or creating a config file.
//
// If the configuration file does not exist or is incomplete, it prompts the user to enter the API key
// and saves it to the configuration file.
//
// Returns:
//   - error: An error if there is any issue reading or writing the configuration file, or if the API key is not set up correctly.
func InitConfig() error {
	configPath, err := getConfigPath()
	if err != nil {
		return fmt.Errorf("error getting configuration file path: %v", err)
	}

	viper.SetConfigFile(configPath)
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Configuration not found. Please set up the API key:")
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter the API key: ")
		apiKey, _ := reader.ReadString('\n')
		apiKey = strings.TrimSpace(apiKey)

		viper.Set("api_key", apiKey)
		err = viper.WriteConfigAs(configPath)
		if err != nil {
			return fmt.Errorf("error saving configuration: %v", err)
		}
		fmt.Println("Configuration saved successfully.")
	} else {
		apiKey := viper.GetString("api_key")
		if apiKey == "" {
			return fmt.Errorf("API key is not set up correctly. Please set up the API key")
		}
	}

	return nil
}

// GetAPIKey retrieves the API key from the configuration file.
//
// It reads the API key from the configuration file managed by Viper and trims any leading or trailing whitespace.
//
// Returns:
//   - string: The API key.
func GetAPIKey() string {
	return strings.TrimSpace(viper.GetString("api_key"))
}
