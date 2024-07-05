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
