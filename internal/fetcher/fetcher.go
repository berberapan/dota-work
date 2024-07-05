package fetcher

import (
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

func getEnvVariable(key string) string {
	if err := godotenv.Load("../.env"); err != nil {
		log.Fatalf("Unable to load .env file. \nError message: %s", err)
	}
	return os.Getenv(key)
}

func fetchDataFromApi(url string) []byte {
	response, err := http.Get(url)
	if err != nil {
		log.Fatalf("Unable to fetch data from API./nError: %s", err)
	}
	body, err := io.ReadAll(response.Body)
	response.Body.Close()
	if response.StatusCode > 299 {
		log.Fatalf("Response failed with status code: %d and\nbody: %s\n", response.StatusCode, body)
	}
	if err != nil {
		log.Fatal(err)
	}
	return body
}

func getTeamMatchHistoryApi(teamID string) []byte {
	envString := getEnvVariable("URL_TEAM_API")
	url := strings.Replace(envString, "{team_id}", teamID, 1)
	return fetchDataFromApi(url)
}

func getMatchDataApi(matchID string) []byte {
	envString := getEnvVariable("URL_MATCH_API")
	url := strings.Replace(envString, "{match_id}", matchID, 1)
	return fetchDataFromApi(url)
}
