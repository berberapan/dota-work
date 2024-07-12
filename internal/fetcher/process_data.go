package fetcher

import (
	"encoding/json"
	"log"
	"strconv"

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

func GetListOfMatchesFromCount(num int, teamId string) []TeamHistory {
	matchHistorySlice := []TeamHistory{}
	listOfMatches := GetListOfMatches(teamId)
	for n := range num {
		matchHistorySlice = append(matchHistorySlice, listOfMatches[n])
	}
	return matchHistorySlice
}

func GetMatchDataFromTeamHistorySlice(matches []TeamHistory) []MatchData {
	sliceOfMatches := []MatchData{}
	for _, match := range matches {
		matchID := strconv.Itoa(match.MatchId)
		sliceOfMatches = append(sliceOfMatches, GetMatchData(matchID))
	}
	return sliceOfMatches
}

func ProcessMatchData(matches []MatchData) CompiledTeamData {
}
