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
	log.Println("Request for GetTeamData received.")

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

	log.Printf("Params received - Team ID %s, Date %s, Count %s", id, date, count)

	var compiledData fetcher.CompiledTeamData
	var processErr error

	if date != "" {
		compiledData, processErr = fetcher.GetCachedOrProcessedData(id, date, 0)
	} else if count != "" {
		conCount, convErr := strconv.Atoi(count)
		if convErr != nil {
			http.Error(w, "Couldn't handle count given.", http.StatusBadRequest)
			return
		}
		compiledData, processErr = fetcher.GetCachedOrProcessedData(id, "", conCount)
	}

	if processErr != nil {
		if processErr.Error() == "No matches found" {
			http.Error(w, "No matches found.", http.StatusNotFound)
			return
		} else {
			http.Error(w, "Error processing data: "+processErr.Error(), http.StatusInternalServerError)
		}
		return
	}

	jsonData, err := utils.DataStructToJson(compiledData)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error converting data to JSON", http.StatusInternalServerError)
		return
	}

	log.Println("Sending back response")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

func GetTournamentSchedule(w http.ResponseWriter, r *http.Request) {
	log.Println("Request for GetTournamentSchedule received.")

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}
	url := r.FormValue("url")
	if url == "" {
		http.Error(w, "Missing required parameter Liquipedia tournament url", http.StatusBadRequest)
		return
	}
	leagueCode := r.FormValue("league-code")
	if leagueCode == "" {
		http.Error(w, "Missing required parameter league code", http.StatusBadRequest)
	}
	tournamentName := r.FormValue("tournament")

	log.Printf("Params received - URL %s, LeagueCode %s, Tournament name %s", url, leagueCode, tournamentName)

	data := fetcher.GetScheduleOfTournament(url, leagueCode, tournamentName)
	jsonData, err := utils.DataStructToJson(data)
	if err != nil {
		log.Println("Error marshaling JSON:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	log.Println("Sending back response")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	w.Write(jsonData)
}
