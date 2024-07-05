package fetcher

import (
	"encoding/json"
	"log"

	"github.com/berberapan/dota-work/internal/utils"
)

func GetListOfMatches(teamId string) []TeamHistory {
	data := getTeamMatchHistoryApi(teamId)
	listOfMatches := []TeamHistory{}
	err := json.Unmarshal(data, &listOfMatches)
	if err != nil {
		log.Println(err)
	}
	return listOfMatches
}

func GetMatchData(matchID string) MatchData {
	data := getMatchDataApi(matchID)
	matchData := MatchData{}
	err := json.Unmarshal(data, &matchData)
	if err != nil {
		log.Println(err)
	}
	return matchData
}

func GetListOfMatchesFromDate(fromDate, teamId string) []TeamHistory {
	matchHistorySlice := []TeamHistory{}
	unixTime := utils.ConvertDateToUnix(fromDate)
	listOfMatches := GetListOfMatches(teamId)
	for _, match := range listOfMatches {
		if match.StartTime < unixTime {
			break
		}
		matchHistorySlice = append(matchHistorySlice, match)
	}
	return matchHistorySlice
}
