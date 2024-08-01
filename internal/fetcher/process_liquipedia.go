package fetcher

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/berberapan/dota-work/internal/utils"
)

const (
	matchIndicator = "{{Match2"
	ukTimeZone     = "Europe/London"
)

var defaultTime time.Time = time.Date(2020, time.Month(1), 1, 0, 0, 0, 0, time.UTC)

func LiquipediaDataToMap(url string) map[string]interface{} {
	var result map[string]interface{}
	tournament, err := utils.LiquipediaUrlToTournamentString(url)
	if err != nil {
		log.Printf("Couldn't use given URL.\n%v", err)
	}
	data := getLiquipediaData(tournament)
	err = json.Unmarshal(data, &result)
	if err != nil {
		log.Printf("Unable to unmarshal JSON from URL.\n%v", err)
	}
	return result
}

func LiquipediaMappedDataToString(data map[string]interface{}) string {
	content, ok := data["parse"].(map[string]interface{})["wikitext"].(map[string]interface{})["*"].(string)
	if !ok {
		log.Println("Couldn't parse content of JSON. Check URL")
		return ""
	}
	return content
}

func GetScheduleOfTournament(url, leagueCode, tournament string) []MatchScheduleData {
	mapData := LiquipediaDataToMap(url)
	data := LiquipediaMappedDataToString(mapData)
	lines := strings.Split(data, "\n")
	leagueCodeInt, _ := strconv.Atoi(leagueCode)
	tournamentName := tournament

	var schedule []MatchScheduleData
	var currentMatch matchInfo
	if tournamentName == "" {
		tournamentName = extractTournamentName(lines)
	}
	isMatch := false
	bestOf := 0

	for _, line := range lines {
		if strings.Contains(line, matchIndicator) {
			isMatch = true
			currentMatch = matchInfo{prevBestOf: bestOf}
		}
		if isMatch {
			currentMatch.updateFromLine(line)
			if currentMatch.isComplete() {
				if currentMatch.dateTime != defaultTime {
					schedule = append(schedule, createMatchScheduleData(currentMatch, tournamentName, leagueCodeInt))
				}
				isMatch = false
				bestOf = currentMatch.bestOf
			}
		}
	}
	return schedule
}

type matchInfo struct {
	bestOf     int
	prevBestOf int
	teamA      string
	teamB      string
	dateTime   time.Time
}

func (m *matchInfo) updateFromLine(line string) {
	switch {
	case strings.Contains(line, "bestof"):
		m.bestOf, _ = extractIntValue(line)
	case strings.Contains(line, "opponent1"):
		m.teamA = extractTeamName(line)
	case strings.Contains(line, "opponent2"):
		m.teamB = extractTeamName(line)
	case strings.Contains(line, "|date"):
		m.dateTime = extractDateTime(line)
	}
}

func (m *matchInfo) isComplete() bool {
	if m.bestOf == 0 && m.teamA != "" && m.teamB != "" && !m.dateTime.IsZero() {
		m.bestOf = m.prevBestOf
	}
	return m.bestOf > 0 && m.teamA != "" && m.teamB != "" && !m.dateTime.IsZero()
}

func extractTeamName(line string) string {
	matchLine := line[1:]
	startIdx := strings.Index(matchLine, "|") + 1
	endIdx := strings.Index(matchLine, "}}")
	teamName := strings.TrimSpace(matchLine[startIdx:endIdx])
	if teamName == "" {
		teamName = "TBA"
	}
	return teamName
}

func extractIntValue(line string) (int, error) {
	valueStr := strings.TrimSpace(line[strings.Index(line, "=")+1:])
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return 0, fmt.Errorf("Failed to convert string to int.\n%w", err)
	}
	return value, nil
}

func extractTournamentName(lines []string) string {
	for _, line := range lines {
		if strings.Contains(line, "tickername") {
			tickerLine := line[strings.Index(line, "tickername"):]
			lastIdx := len(tickerLine)
			if strings.Contains(tickerLine, "|") {
				lastIdx = strings.Index(tickerLine, "|")
			}
			return strings.TrimSpace(tickerLine[strings.Index(tickerLine, "=")+1 : lastIdx])
		}
	}
	return ""
}

func extractDateTime(line string) time.Time {
	dateLine := strings.TrimSpace(line[strings.Index(line, "=")+1:])
	dateLineClean := strings.ReplaceAll(dateLine, ",", "")
	objects := strings.Fields(dateLineClean)
	if len(objects) == 0 {
		return defaultTime
	}
	srcYear, err := strconv.Atoi(objects[2])
	if err != nil {
		log.Printf("Unable to convert string Year to int Year.\n%v", err)
		return time.Time{}
	}
	srcMonth := utils.ConvertStringToMonth(objects[0])
	srcDay, err := strconv.Atoi(objects[1])
	if err != nil {
		log.Printf("Unable to convert string Day to int Day.\n%v", err)
		return time.Time{}
	}

	hhMM := strings.Split(objects[4], ":")
	srcHour, err := strconv.Atoi(hhMM[0])
	if err != nil {
		log.Printf("Unable to convert string Hour to int Hour.\n%v", err)
		return time.Time{}
	}
	srcMin, err := strconv.Atoi(hhMM[1])
	if err != nil {
		log.Printf("Unable to convert string Minute to int Minute.\n%v", err)
		return time.Time{}
	}

	tzObject := objects[5]
	tzAbbr := tzObject[strings.Index(tzObject, "/")+1 : strings.Index(tzObject, "}}")]
	srcLoc := utils.ConvertTzAbbrToTzString(tzAbbr)

	location, err := time.LoadLocation(srcLoc)
	if err != nil {
		log.Printf("Unable to load source timezone. Setting UTC timezone.\n%v", err)
		location = time.UTC
	}
	ukTime, err := time.LoadLocation(ukTimeZone)
	if err != nil {
		log.Printf("Unable to load UK timezone.\n%v", err)
		return defaultTime
	}
	srcTime := time.Date(srcYear, srcMonth, srcDay, srcHour, srcMin, 0, 0, location)
	return srcTime.In(ukTime)
}

func createMatchScheduleData(m matchInfo, tournamentName string, leagueCode int) MatchScheduleData {
	return MatchScheduleData{
		TeamA:      m.teamA,
		TeamB:      m.teamB,
		Tournament: tournamentName,
		BestOf:     m.bestOf,
		LeagueCode: leagueCode,
		Date:       utils.ConvertDateToString(m.dateTime),
		StartTime:  utils.ConvertTimeToString(m.dateTime),
		Game:       "DOTA2",
	}
}
