package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/berberapan/dota-work/internal/fetcher"
	"github.com/berberapan/dota-work/internal/utils"
)

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func GetTeamData(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}
	id := r.FormValue("team_id")
	if id == "" {
		http.Error(w, "Missing required parameter Team ID", http.StatusBadRequest)
		return
	}
	date := r.FormValue("date")
	count := r.FormValue("count")

	if date == "" && count == "" {
		http.Error(w, "Missing parameter. Date or count must to be provided.", http.StatusBadRequest)
		return
	}

	var jsonData []byte

	if date != "" {
		matchHistory := fetcher.GetListOfMatchesFromDate(date, id)
		matchSlice := fetcher.GetMatchDataFromTeamHistorySlice(matchHistory)
		compiledData := fetcher.ProcessMatchData(matchSlice, id)
		jsonData, err = utils.DataStructToJson(compiledData)
		if err != nil {
			log.Println(err)
		}

	} else if count != "" {
		conCount, err := strconv.Atoi(count)
		if err != nil {
			http.Error(w, "Couldn't handle count given.", http.StatusBadRequest)
			return
		}
		matchHistory := fetcher.GetListOfMatchesFromCount(conCount, id)
		matchSlice := fetcher.GetMatchDataFromTeamHistorySlice(matchHistory)
		compiledData := fetcher.ProcessMatchData(matchSlice, id)
		jsonData, err = utils.DataStructToJson(compiledData)
		if err != nil {
			log.Println(err)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}
