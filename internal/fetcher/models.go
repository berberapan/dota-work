package fetcher

type MatchData struct {
	Objectives            []Objective `json:"objectives"`
	StartTime             int         `json:"start_time"`
	Duration              float64     `json:"duration"`
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
	Time       int         `json:"time"`
	Type       string      `json:"type"`
	Key        interface{} `json:"key,omitempty"`
	Team       int         `json:"team,omitempty"`
	PlayerSlot int         `json:"player_slot,omitempty"`
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
	LostMatches     CompiledMatchData `json:"lost_matches"`
	RadiantMatches  CompiledMatchData `json:"radiant_matches"`
	DireMatches     CompiledMatchData `json:"dire_matches"`
}

type CompiledMatchData struct {
	AverageDuration      string    `json:"average_duration"`
	MedianDuration       string    `json:"median_duration"`
	TeamTowersAverage    float64   `json:"team_towers_average"`
	TeamTowersMedian     float64   `json:"team_towers_median"`
	TotalTowersAverage   float64   `json:"total_towers_average"`
	TotalTowersMedian    float64   `json:"total_towers_median"`
	FirstTowerPct        float64   `json:"first_tower_pct"`
	TeamBarracksAverage  float64   `json:"team_barracks_average"`
	TeamBarracksMedian   float64   `json:"team_barracks_median"`
	TotalBarracksAverage float64   `json:"total_barracks_average"`
	TotalBarracksMedian  float64   `json:"total_barracks_median"`
	FirstBarrackPct      float64   `json:"first_barrack_pct"`
	TeamScoreAverage     float64   `json:"team_score_average"`
	TeamScoreMedian      float64   `json:"team_score_median"`
	TotalScoreAverage    float64   `json:"total_score_average"`
	TotalScoreMedian     float64   `json:"total_score_median"`
	FirstBloodPct        float64   `json:"first_blood_pct"`
	TeamRoshanAverage    float64   `json:"team_roshan_average"`
	TeamRoshanMedian     float64   `json:"team_roshan_median"`
	TotalRoshansAverage  float64   `json:"total_roshans_average"`
	TotalRoshansMedian   float64   `json:"total_roshans_median"`
	FirstRoshanPct       float64   `json:"first_roshan_pct"`
	Durations            []float64 `json:"-"`
	TeamTowers           []int     `json:"-"`
	TowerTotals          []int     `json:"-"`
	TeamBarracks         []int     `json:"-"`
	BarrackTotals        []int     `json:"-"`
	TeamScores           []int     `json:"-"`
	ScoreTotals          []int     `json:"-"`
	FirstTower           []int     `json:"-"`
	FirstBarrack         []int     `json:"-"`
	FirstBlood           []int     `json:"-"`
	TeamRoshans          []int     `json:"-"`
	RoshanTotals         []int     `json:"-"`
	FirstRoshan          []int     `json:"-"`
}

type MatchScheduleData struct {
	TeamA      string `json:"team_a"`
	TeamB      string `json:"team_b"`
	Tournament string `json:"tournament"`
	LeagueCode int    `json:"league_code"`
	BestOf     int    `json:"best_of"`
	Date       string `json:"date"`
	StartTime  string `json:"start_time"`
	Game       string `json:"game"`
}
