package fetcher

import (
	"encoding/json"
	"log"
	"strconv"
	"strings"

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
	if num > len(listOfMatches) {
		num = len(listOfMatches)
	}
	for i := 0; i < num; i++ {
		matchHistorySlice = append(matchHistorySlice, listOfMatches[i])
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

func ProcessMatchData(matches []MatchData, teamId string) CompiledTeamData {
	id, err := strconv.Atoi(teamId)
	if err != nil {
		log.Println("Unable to convert Team ID from string to int")
	}
	compiledData := CompiledTeamData{
		TeamID:          id,
		NumberOfMatches: len(matches),
	}

	for _, match := range matches {
		isRadiant := id == match.RadiantTeamID
		isWin := (isRadiant && match.RadiantWin) || (!isRadiant && !match.RadiantWin)

		updateMatchData(&compiledData.AllMatches, match, isRadiant)

		if isWin {
			updateMatchData(&compiledData.WonMatches, match, isRadiant)
		} else {
			updateMatchData(&compiledData.LostMatches, match, isRadiant)
		}

		if isRadiant {
			updateMatchData(&compiledData.RadiantMatches, match, isRadiant)
		} else {
			updateMatchData(&compiledData.DireMatches, match, isRadiant)
		}
	}

	calculateTeamData(&compiledData)

	return compiledData
}

func updateMatchData(data *CompiledMatchData, match MatchData, isRadiant bool) {
	data.Durations = append(data.Durations, match.Duration)
	if isRadiant {
		data.TeamTowers = append(data.TeamTowers, utils.ConvertBinIntToCount(match.TowerStatusDire, 11))
		data.TeamBarracks = append(data.TeamBarracks, utils.ConvertBinIntToCount(match.BarracksStatusDire, 6))
		data.TeamScores = append(data.TeamScores, match.RadiantScore)
	} else {
		data.TeamTowers = append(data.TeamTowers, utils.ConvertBinIntToCount(match.TowerStatusRadiant, 11))
		data.TeamBarracks = append(data.TeamBarracks, utils.ConvertBinIntToCount(match.BarracksStatusRadiant, 6))
		data.TeamScores = append(data.TeamScores, match.DireScore)
	}
	totalTowers := utils.ConvertBinIntToCount(match.TowerStatusDire, 11) + utils.ConvertBinIntToCount(match.TowerStatusRadiant, 11)
	data.TowerTotals = append(data.TowerTotals, totalTowers)
	totalBarracks := utils.ConvertBinIntToCount(match.BarracksStatusDire, 6) + utils.ConvertBinIntToCount(match.BarracksStatusRadiant, 6)
	data.BarrackTotals = append(data.BarrackTotals, totalBarracks)
	totalScore := match.DireScore + match.RadiantScore
	data.ScoreTotals = append(data.ScoreTotals, totalScore)

	firstTower := -1
	firstBarrack := -1
	firstBlood := -1
	firstRoshan := -1
	totalRoshan := 0
	teamRoshan := 0

	for _, objective := range match.Objectives {

		if objective.Type == "CHAT_MESSAGE_FIRSTBLOOD" {
			if objective.PlayerSlot < 5 && isRadiant {
				firstBlood = 1
			} else if objective.PlayerSlot > 5 && !isRadiant {
				firstBlood = 1
			} else {
				firstBlood = 0
			}
		}

		if firstTower == -1 {
			if key, ok := objective.Key.(string); ok && strings.Contains(key, "tower") {
				switch {
				case strings.Contains(key, "bad") && isRadiant:
					firstTower = 1
				case strings.Contains(key, "good") && !isRadiant:
					firstTower = 1
				default:
					firstTower = 0
				}
			}
		}

		if firstBarrack == -1 {
			if key, ok := objective.Key.(string); ok && strings.Contains(key, "rax") {
				switch {
				case strings.Contains(key, "bad") && isRadiant:
					firstBarrack = 1
				case strings.Contains(key, "good") && !isRadiant:
					firstBarrack = 1
				default:
					firstBarrack = 0
				}
			}
		}

		if objective.Type == "CHAT_MESSAGE_ROSHAN_KILL" {
			switch {
			case objective.Team == 2 && isRadiant:
				teamRoshan++
				totalRoshan++
				if firstRoshan == -1 {
					firstRoshan = 1
				}
			case objective.Team == 3 && !isRadiant:
				teamRoshan++
				totalRoshan++
				if firstRoshan == -1 {
					firstRoshan = 1
				}
			default:
				totalRoshan++
				if firstRoshan == -1 {
					firstRoshan = 0
				}
			}
		}
	}

	data.FirstTower = append(data.FirstTower, firstTower)
	data.FirstBarrack = append(data.FirstBarrack, firstBarrack)
	data.FirstBlood = append(data.FirstBlood, firstBlood)
	data.FirstRoshan = append(data.FirstRoshan, firstRoshan)
	data.TeamRoshans = append(data.TeamRoshans, teamRoshan)
	data.RoshanTotals = append(data.RoshanTotals, totalRoshan)
}

func calculateTeamData(data *CompiledTeamData) {
	calculateMatchData(&data.AllMatches)
	calculateMatchData(&data.WonMatches)
	calculateMatchData(&data.LostMatches)
	calculateMatchData(&data.RadiantMatches)
	calculateMatchData(&data.DireMatches)
}

func calculateMatchData(data *CompiledMatchData) {
	data.NumberOfMaps = len(data.Durations)
	if len(data.Durations) > 0 {
		data.AverageDuration = utils.ConvertSecondsToFormattedString(utils.CalculateAverageFromFloatSlice(data.Durations))
		data.MedianDuration = utils.ConvertSecondsToFormattedString(utils.CalculateMedianFromFloatSlice(data.Durations))
	}
	if len(data.TeamTowers) > 0 {
		data.TeamTowersAverage = utils.CalculateAverageFromIntSlice(data.TeamTowers)
		data.TeamTowersMedian = utils.CalculateMedianFromIntSlice(data.TeamTowers)
		data.TeamTowersMode = utils.IntSliceToString(utils.CalculateModeFromIntSlice(data.TeamTowers))
	}
	if len(data.TowerTotals) > 0 {
		data.TotalTowersAverage = utils.CalculateAverageFromIntSlice(data.TowerTotals)
		data.TotalTowersMedian = utils.CalculateMedianFromIntSlice(data.TowerTotals)
		data.TotalTowersMode = utils.IntSliceToString(utils.CalculateModeFromIntSlice(data.TowerTotals))
	}
	if len(data.TeamBarracks) > 0 {
		data.TeamBarracksAverage = utils.CalculateAverageFromIntSlice(data.TeamBarracks)
		data.TeamBarracksMedian = utils.CalculateMedianFromIntSlice(data.TeamBarracks)
		data.TeamBarracksMode = utils.IntSliceToString(utils.CalculateModeFromIntSlice(data.TeamBarracks))
	}
	if len(data.BarrackTotals) > 0 {
		data.TotalBarracksAverage = utils.CalculateAverageFromIntSlice(data.BarrackTotals)
		data.TotalBarracksMedian = utils.CalculateMedianFromIntSlice(data.BarrackTotals)
		data.TotalBarracksMode = utils.IntSliceToString(utils.CalculateModeFromIntSlice(data.BarrackTotals))
	}
	if len(data.TeamScores) > 0 {
		data.TeamScoreAverage = utils.CalculateAverageFromIntSlice(data.TeamScores)
		data.TeamScoreMedian = utils.CalculateMedianFromIntSlice(data.TeamScores)
		data.TeamScoreMode = utils.IntSliceToString(utils.CalculateModeFromIntSlice(data.TeamScores))
	}
	if len(data.ScoreTotals) > 0 {
		data.TotalScoreAverage = utils.CalculateAverageFromIntSlice(data.ScoreTotals)
		data.TotalScoreMedian = utils.CalculateMedianFromIntSlice(data.ScoreTotals)
		data.TotalScoreMode = utils.IntSliceToString(utils.CalculateModeFromIntSlice(data.ScoreTotals))
	}
	if len(data.FirstTower) > 0 {
		data.FirstTowerPct = utils.CalculatePercentageFromIntSlice(data.FirstTower)
	}
	if len(data.FirstBarrack) > 0 {
		data.FirstBarrackPct = utils.CalculatePercentageFromIntSlice(data.FirstBarrack)
	}
	if len(data.FirstBlood) > 0 {
		data.FirstBloodPct = utils.CalculatePercentageFromIntSlice(data.FirstBlood)
	}
	if len(data.FirstRoshan) > 0 {
		data.FirstRoshanPct = utils.CalculatePercentageFromIntSlice(data.FirstRoshan)
	}
	if len(data.TeamRoshans) > 0 {
		data.TeamRoshanAverage = utils.CalculateAverageFromIntSlice(data.TeamRoshans)
		data.TeamRoshanMedian = utils.CalculateMedianFromIntSlice(data.TeamRoshans)
		data.TeamRoshanMode = utils.IntSliceToString(utils.CalculateModeFromIntSlice(data.TeamRoshans))
	}
	if len(data.RoshanTotals) > 0 {
		data.TotalRoshansAverage = utils.CalculateAverageFromIntSlice(data.RoshanTotals)
		data.TotalRoshansMedian = utils.CalculateMedianFromIntSlice(data.RoshanTotals)
		data.TotalRoshansMode = utils.IntSliceToString(utils.CalculateModeFromIntSlice(data.RoshanTotals))
	}
}
