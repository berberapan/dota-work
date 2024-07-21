package fetcher

import (
	"encoding/json"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/berberapan/dota-work/internal/utils"
)

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

func GetScheduleOfTournament(url string, leagueCode int) []MatchScheduleData {
	schedule := []MatchScheduleData{}
	mapData := LiquipediaDataToMap(url)
	data := LiquipediaMappedDataToString(mapData)

	lines := strings.Split(data, "\n")
	isMatch := false

	bestof := 0
	teamA := ""
	teamB := ""
	date := ""
	startTime := ""
	lc := leagueCode
	tournamentName := ""
	game := "DOTA2"

	for _, line := range lines {
		if strings.Contains(line, "tickername") {
			tournamentName = line[strings.Index(line, "=")+1:]
		}
		if strings.Contains(line, "{{Match2") {
			isMatch = true
		}
		if isMatch && strings.Contains(line, "bestof") {
			bestOfString := line[strings.Index(line, "=")+1:]
			bestof, _ = strconv.Atoi(bestOfString)
		}
		if isMatch && strings.Contains(line, "opponent1") {
			startIdx := strings.LastIndex(line, "|") + 1
			endIdx := strings.Index(line, "}}")
			teamA = line[startIdx:endIdx]
		}
		if isMatch && strings.Contains(line, "opponent2") {
			startIdx := strings.LastIndex(line, "|") + 1
			endIdx := strings.Index(line, "}}")
			teamB = line[startIdx:endIdx]
		}
		if isMatch && strings.Contains(line, "date") {
			dateLine := line[strings.Index(line, "=")+1:]
			dateLineClean := strings.ReplaceAll(dateLine, ",", "")
			objects := strings.Fields(dateLineClean)
			srcYear, err := strconv.Atoi(objects[2])
			if err != nil {
				log.Printf("Unable to convert string Year to int Year.\n%v", err)
				return schedule
			}
			srcMonth := utils.ConvertStringToMonth(objects[0])
			srcDay, err := strconv.Atoi(objects[1])
			if err != nil {
				log.Printf("Unable to convert string Day to int Day.\n%v", err)
				return schedule
			}
			hhMM := strings.Split(objects[4], ":")
			srcHour, err := strconv.Atoi(hhMM[0])
			if err != nil {
				log.Printf("Unable to convert string Hour to int Hour.\n%v", err)
				return schedule
			}
			srcMin, err := strconv.Atoi(hhMM[1])
			if err != nil {
				log.Printf("Unable to convert string Hour to int Hour.\n%v", err)
				return schedule
			}
			tzObject := objects[5]
			tzAbbr := tzObject[strings.Index(tzObject, "/")+1 : strings.Index(tzObject, "}}")]
			srcLoc := utils.ConvertTzAbbrToTzString(tzAbbr)
			location, err := time.LoadLocation(srcLoc)
			if err != nil {
				log.Printf("Unable to load source timezone.\n%v", err)
				return schedule
			}
			ukTime, err := time.LoadLocation("Europe/London")
			if err != nil {
				log.Printf("Unable to load UK time.\n%v", err)
				return schedule
			}

			srcTime := time.Date(srcYear, srcMonth, srcDay, srcHour, srcMin, 0, 0, location)
			convertedTime := srcTime.In(ukTime)

			date = utils.ConvertDateToString(convertedTime)
			startTime = utils.ConvertTimeToString(convertedTime)

			match := MatchScheduleData{
				TeamA:      teamA,
				TeamB:      teamB,
				Tournament: tournamentName,
				BestOf:     bestof,
				LeagueCode: lc,
				Date:       date,
				StartTime:  startTime,
				Game:       game,
			}
			schedule = append(schedule, match)
			isMatch = false
		}
	}
	return schedule
}
