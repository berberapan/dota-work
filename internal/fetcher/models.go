package fetcher

type MatchData struct {
	Objectives            []Objective `json:"objectives"`
	StartTime             int         `json:"start_time"`
	Duration              int         `json:"duration"`
	RadiantWin            bool        `json:"radiant_win"`
	TowerStatusRadiant    int         `json:"tower_status_radiant"`
	TowerStatusDire       int         `json:"tower_status_dire"`
	BarracksStatusRadiant int         `json:"barracks_status_radiant"`
	BarracksStatusDire    int         `json:"barracks_status_dire"`
	RadiantScore          int         `json:"radiant_score"`
	DireScore             int         `json:"dire_score"`
	RadiantTeamID         int         `json:"radiant_team_id"`
}

type Objective struct {
	Time int         `json:"time"`
	Type string      `json:"type"`
	Key  interface{} `json:"key,omitempty"`
	Team int         `json:"team,omitempty"`
}

type TeamHistory struct {
	MatchId   int `json:"match_id"`
	StartTime int `json:"start_time"`
}

type CompiledTeamData struct {
	TeamID          int               `json:"team_id"`
	NumberOfMatches int               `json:"number_of_matches"`
	AllMatches      CompiledMatchData `json:"all_matches"`
	WonMatches      CompiledMatchData `json:"won_matches"`
	RadiantMatches  CompiledMatchData `json:"radiant_matches"`
}

type CompiledMatchData struct {
	AverageDuration      string  `json:"average_duration"`
	MedianDuration       string  `json:"median_duration"`
	TeamTowersAverage    float64 `json:"team_towers_average"`
	TotalTowersAverage   float64 `json:"total_towers_average"`
	FirstTowerPct        float64 `json:"first_tower_pct"`
	TeamBarracksAverage  float64 `json:"team_barracks_average"`
	TotalBarracksAverage float64 `json:"total_barracks_average"`
	FirstBarrackPct      float64 `json:"first_barrack_pct"`
	TeamScoreAverage     float64 `json:"team_score_average"`
	TeamScoreMedian      int     `json:"team_score_median"`
	TotalScoreAverage    float64 `json:"total_score_average"`
	TotalScoreMedian     int     `json:"total_score_median"`
	FirstBloodPct        float64 `json:"first_blood_pct"`
}
