package graphql

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/simplesmentemat/stealth-grid-cli/pkg/config"
	"github.com/spf13/viper"
)

// Mock server for testing
func mockServer() *httptest.Server {
	handler := http.NewServeMux()
	handler.HandleFunc("/file-download/events/grid/series/2620066", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/zip")
		w.Write([]byte("zip content"))
	})
	return httptest.NewServer(handler)
}

func TestFetchData(t *testing.T) {
	startTime := time.Now().Add(-24 * time.Hour)
	endTime := time.Now()
	data, err := FetchData("3", startTime, endTime)
	if err != nil {
		t.Fatalf("Failed to fetch data: %v", err)
	}
	if data == nil {
		t.Fatalf("Expected non-nil data")
	}
}

func TestDownloadJSON(t *testing.T) {
	server := mockServer()
	defer server.Close()

	// Mocking the API URL
	config.APIURL = server.URL

	// Set up a temporary config file with an API key for testing
	configPath, err := os.CreateTemp("", "config.yaml")
	if err != nil {
		t.Fatalf("Failed to create temp config file: %v", err)
	}
	defer os.Remove(configPath.Name())

	// Check and print the config type and content for debugging
	fmt.Printf("Config type: %s\n", viper.GetViper().ConfigFileUsed())
	fmt.Printf("Config content: %s\n", viper.AllSettings())

	// Set the environment variable for the config file path
	os.Setenv("CONFIG_PATH", configPath.Name())

	// Initialize config
	if err := config.InitConfig(); err != nil {
		t.Fatalf("Failed to initialize config: %v", err)
	}

	err = os.MkdirAll("/tmp", os.ModePerm)
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}

	DownloadJSON("2620066", "/tmp")
	if _, err := os.Stat("/tmp/2620066.zip"); os.IsNotExist(err) {
		t.Fatalf("Expected file 2620066.zip to be created, but it does not exist")
	}
}
