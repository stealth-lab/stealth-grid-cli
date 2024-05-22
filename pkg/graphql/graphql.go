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

type QueryVariables struct {
	StartTime   string `json:"startTime"`
	EndTime     string `json:"endTime"`
	AfterCursor string `json:"afterCursor"`
	TitleIDs    string `json:"titleIds"`
}

type GraphQLRequest struct {
	Query     string         `json:"query"`
	Variables QueryVariables `json:"variables"`
}

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

	req, err := http.NewRequest("POST", "https://api.grid.gg/central-data/graphql", bytes.NewBuffer(reqBody))
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

func DownloadJSON(serieID string, directory string) {
	url := fmt.Sprintf("https://api.grid.gg/file-download/events/grid/series/%s", serieID)
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
