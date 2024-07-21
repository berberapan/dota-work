package fetcher

import (
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/berberapan/dota-work/internal/utils"
)

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
	envString := utils.GetEnvVariable("URL_TEAM_API")
	url := strings.Replace(envString, "{team_id}", teamID, 1)
	return fetchDataFromApi(url)
}

func getMatchDataApi(matchID string) []byte {
	envString := utils.GetEnvVariable("URL_MATCH_API")
	url := strings.Replace(envString, "{match_id}", matchID, 1)
	return fetchDataFromApi(url)
}

func fetchDataFromLiquipediaGzip(url string) ([]byte, error) {
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("Unable to create request.\n%w", err)
	}
	request.Header.Set("User-Agent", "Personal project, github.com/berberapan/dota-work")
	request.Header.Set("Accept-Encoding", "gzip")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("Unable to make request.\n%w", err)
	}
	defer response.Body.Close()

	gzipReader, err := gzip.NewReader(response.Body)
	if err != nil {
		log.Println("Unable to create gzip reader.\n", err)
		return nil, fmt.Errorf("Unable to create gzip reader.\n%w", err)
	}
	defer gzipReader.Close()

	body, err := io.ReadAll(gzipReader)
	if err != nil {
		return nil, fmt.Errorf("Unable to read gzip response from source.\n%w", err)
	}
	return body, nil
}

func getLiquipediaData(tournament string) []byte {
	params := url.Values{}
	params.Add("action", "parse")
	params.Add("page", tournament)
	params.Add("prop", "wikitext")
	params.Add("format", "json")
	envString := utils.GetEnvVariable("URL_LIQUIPEDIA")
	url := envString + "?" + params.Encode()

	data, err := fetchDataFromLiquipediaGzip(url)
	if err != nil {
		log.Println(err)
	}
	return data
}
