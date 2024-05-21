package graphql

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
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

/*
 */
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

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("x-api-key", "DQmZW6rVNkTT2iZMAsLpzHbuGoTzYTEmiyqLht6p")

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

/*
 */
func DownloadJSON(serieID string) {
	url := fmt.Sprintf("https://api.grid.gg/file-download/events/grid/series/%s", serieID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Printf("Erro ao criar solicitação: %v", err)
		return
	}

	req.Header.Add("x-api-key", "DQmZW6rVNkTT2iZMAsLpzHbuGoTzYTEmiyqLht6p")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Erro ao baixar o ZIP: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Erro: Código de status %d\n", resp.StatusCode)
		return
	}

	file, err := os.Create(fmt.Sprintf("%s.zip", serieID))
	if err != nil {
		fmt.Printf("Erro ao criar o arquivo: %v", err)
		return
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		fmt.Printf("Erro ao salvar o ZIP no arquivo: %v", err)
		return
	}

	fmt.Printf("Arquivo ZIP para a série %s baixado com sucesso.\n", serieID)
}
