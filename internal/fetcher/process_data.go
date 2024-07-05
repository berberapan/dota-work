package fetcher

import (
	"encoding/json"
	"log"
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
