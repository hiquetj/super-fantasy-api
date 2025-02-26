package models

type LeagueSettings struct {
	Batting struct {
		RunsScored      float64 `json:"runs_scored"`
		TotalBases      float64 `json:"total_bases"`
		RunsBattedIn    float64 `json:"runs_batted_in"`
		Walks           float64 `json:"walks"`
		Strikeouts      float64 `json:"strikeouts"`
		StolenBases     float64 `json:"stolen_bases"`
		HittingForCycle float64 `json:"hitting_for_cycle"`
	} `json:"batting"`
	Pitching struct {
		InningsPitched float64 `json:"innings_pitched"`
		HitsAllowed    float64 `json:"hits_allowed"`
		EarnedRuns     float64 `json:"earned_runs"`
		WalksIssued    float64 `json:"walks_issued"`
		Strikeouts     float64 `json:"strikeouts"`
		NoHitters      float64 `json:"no_hitters"`
		PerfectGames   float64 `json:"perfect_games"`
		Wins           float64 `json:"wins"`
		Losses         float64 `json:"losses"`
		Saves          float64 `json:"saves"`
		Holds          float64 `json:"holds"`
	} `json:"pitching"`
}

type ProjectionRequest struct {
	Settings       LeagueSettings `json:"settings"`
	ProjectionName string         `json:"projection_name"`
}

type PlayerProjection struct {
	PlayerID    string  `json:"player_id"`
	PlayerName  string  `json:"player_name"`
	TotalPoints float64 `json:"total_points"`
}