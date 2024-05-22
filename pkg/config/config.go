package config

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

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

func InitConfig() error {
	configPath, err := getConfigPath()
	if err != nil {
		return fmt.Errorf("error getting configuration file path: %v", err)
	}
	viper.SetConfigFile(configPath)

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

func GetAPIKey() string {
	return strings.TrimSpace(viper.GetString("api_key"))
}
