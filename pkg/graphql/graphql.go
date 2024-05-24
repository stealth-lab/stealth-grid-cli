package graphql

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/simplesmentemat/stealth-grid-cli/pkg/config"
)

// QueryVariables represents the variables for the GraphQL query.
type QueryVariables struct {
	StartTime   string `json:"startTime"`
	EndTime     string `json:"endTime"`
	AfterCursor string `json:"afterCursor"`
	TitleIDs    string `json:"titleIds"`
}

// GraphQLRequest represents the structure of a GraphQL request.
type GraphQLRequest struct {
	Query     string         `json:"query"`
	Variables QueryVariables `json:"variables"`
}

// FetchData fetches data from a GraphQL API given a title ID and a time range.
//
// This function constructs a GraphQL query to fetch series data from the API
// "https://api.grid.gg/central-data/graphql" based on the provided title ID and
// time range. The data is retrieved using a POST request and is returned as a
// map. If any error occurs during the process, it is returned.
//
// Parameters:
//   - titleID: A string representing the ID of the title to query for. This is used
//     to filter the series based on the specific title.
//   - startTime: A time.Time object representing the start time of the query range.
//     This is converted to RFC3339 format and used in the query to filter series
//     that start on or after this time.
//   - endTime: A time.Time object representing the end time of the query range.
//     This is converted to RFC3339 format and used in the query to filter series
//     that end on or before this time.
//
// Returns:
//   - A map[string]interface{} containing the query results if the request is
//     successful. The map includes information about the series such as total
//     count, page info, and details about each series.
//   - An error if the request fails at any point. Errors can occur during JSON
//     marshalling of the request, creation of the HTTP request, sending the HTTP
//     request, or decoding the JSON response.
func FetchData(titleID string, startTime, endTime time.Time) (map[string]interface{}, error) {
	variables := QueryVariables{
		StartTime:   startTime.Format(time.RFC3339),
		EndTime:     endTime.Format(time.RFC3339),
		AfterCursor: "",
		TitleIDs:    titleID,
	}

	query := `query GetAllSeriesInNext24Hours($startTime: String, $endTime: String, $afterCursor: Cursor, $titleIds: [ID!]) {
		allSeries(first: 50, filter: {startTimeScheduled: {gte: $startTime, lte: $endTime}, titleIds: {in: $titleIds}}, orderBy: StartTimeScheduled, after: $afterCursor) {
			totalCount
			pageInfo {
				hasPreviousPage
				hasNextPage
				startCursor
				endCursor
			}
			edges {
				cursor
				node {
					id
					tournament {
						nameShortened
						name
						id
					}
					startTimeScheduled
					format {
						nameShortened
					}
					teams {
						baseInfo {
							name
							id
						}
					}
				}
			}
		}
	}`

	graphQLReq := GraphQLRequest{
		Query:     query,
		Variables: variables,
	}

	reqBody, err := json.Marshal(graphQLReq)
	if err != nil {
		return nil, fmt.Errorf("error marshalling GraphQL request: %v", err)
	}

	req, err := http.NewRequest("POST", config.APIURL+"/central-data/graphql", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	apiKey := config.GetAPIKey()
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("x-api-key", apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request to server: %v", err)
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error decoding JSON response: %v", err)
	}

	return result, nil
}

// DownloadJSON downloads a ZIP file for a given series ID from the specified API.
//
// This function constructs a URL to download a ZIP file related to the specified
// series ID. It sends an HTTP GET request to the URL and handles the response,
// saving the ZIP file locally. If any error occurs during the process, it logs
// the error and terminates.
//
// Parameters:
//   - serieID: A string representing the ID of the series to download the ZIP file for.
//     This ID is used to construct the download URL.
//
// The function performs the following steps:
//  1. Constructs the download URL using the provided series ID.
//  2. Creates an HTTP GET request to the constructed URL.
//  3. Sets the necessary headers (including the API key) for the request.
//  4. Sends the request using an HTTP client and handles the response.
//  5. Checks if the response status code is OK (200). If not, logs an error and terminates.
//  6. Creates a file to save the downloaded ZIP content.
//  7. Copies the content from the response body to the created file.
//  8. Logs a success message if the file is saved successfully, or an error message if any step fails.
func DownloadJSON(serieID string, directory string) {
	url := fmt.Sprintf("%s/file-download/events/grid/series/%s", config.APIURL, serieID)
	fmt.Println(directory)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Printf("Erro ao criar solicitação: %v\n", err)
		return
	}

	apiKey := config.GetAPIKey()
	req.Header.Add("x-api-key", apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Erro ao baixar o ZIP: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Erro: Código de status %d\n", resp.StatusCode)
		return
	}

	filePath := filepath.Join(directory, fmt.Sprintf("%s.zip", serieID))
	file, err := os.Create(filePath)
	if err != nil {
		fmt.Printf("Erro ao criar o arquivo: %v\n", err)
		return
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		fmt.Printf("Erro ao salvar o ZIP no arquivo: %v\n", err)
		return
	}
}
