package baseball

import (
	"super-fantasy-api/models"
)

// Batter represents a FantasyPros batter projection
// Based on CSV: "Player","Team","Positions","AB","R","HR","RBI","SB","AVG","OBP","H","2B","3B","BB","SO","SLG","OPS"
type FantasyProsBatter struct {
	Name        string  `bson:"name" json:"name"`                 // Player
	Team        string  `bson:"team" json:"team"`                 // Team
	Positions   string  `bson:"positions" json:"positions"`       // Positions
	AtBats      float64 `bson:"at_bats" json:"at_bats"`           // AB
	Runs        float64 `bson:"runs" json:"runs"`                 // R
	HomeRuns    float64 `bson:"home_runs" json:"home_runs"`       // HR
	RBI         float64 `bson:"rbi" json:"rbi"`                   // RBI
	StolenBases float64 `bson:"stolen_bases" json:"stolen_bases"` // SB
	AVG         float64 `bson:"avg" json:"avg"`                   // AVG
	OBP         float64 `bson:"obp" json:"obp"`                   // OBP
	Hits        float64 `bson:"hits" json:"hits"`                 // H
	Doubles     float64 `bson:"doubles" json:"doubles"`           // 2B
	Triples     float64 `bson:"triples" json:"triples"`           // 3B
	Walks       float64 `bson:"walks" json:"walks"`               // BB
	Strikeouts  float64 `bson:"strikeouts" json:"strikeouts"`     // SO
	SLG         float64 `bson:"slg" json:"slg"`                   // SLG
	OPS         float64 `bson:"ops" json:"ops"`                   // OPS
	Year        string  `bson:"year" json:"year"`                 // year
}

// Pitcher represents a FantasyPros pitcher projection
// Based on CSV: "Player","Team","Positions","IP","K","W","SV","ERA","WHIP","ER","H","BB","HR","G","GS","L","CG"
type FanstasyProsPitcher struct {
	Name            string  `bson:"name" json:"name"`                           // Player
	Team            string  `bson:"team" json:"team"`                           // Team
	Positions       string  `bson:"positions" json:"positions"`                 // Positions
	InningsPitched  float64 `bson:"innings_pitched" json:"innings_pitched"`     // IP
	Strikeouts      float64 `bson:"strikeouts" json:"strikeouts"`               // K
	Wins            float64 `bson:"wins" json:"wins"`                           // W
	Saves           float64 `bson:"saves" json:"saves"`                         // SV
	ERA             float64 `bson:"era" json:"era"`                             // ERA
	WHIP            float64 `bson:"whip" json:"whip"`                           // WHIP
	EarnedRuns      float64 `bson:"earned_runs" json:"earned_runs"`             // ER
	HitsAllowed     float64 `bson:"hits_allowed" json:"hits_allowed"`           // H
	Walks           float64 `bson:"walks" json:"walks"`                         // BB
	HomeRunsAllowed float64 `bson:"home_runs_allowed" json:"home_runs_allowed"` // HR
	Games           float64 `bson:"games" json:"games"`                         // G
	GamesStarted    float64 `bson:"games_started" json:"games_started"`         // GS
	Losses          float64 `bson:"losses" json:"losses"`                       // L
	CompleteGames   float64 `bson:"complete_games" json:"complete_games"`       // CG
	Year            string  `bson:"year" json:"year"`                           // year
}

// CalculateFantasyProsBatterPoints converts FantasyPros batter projections to fantasy points using league settings
func CalculateFantasyProsBatterPoints(player FantasyProsBatter, settings models.LeagueSettings) models.PlayerProjection {
	totalPoints := 0.0
	totalPoints += player.Runs * settings.Batting.RunsScored
	// Total Bases: H = 1B + 2B + 3B + HR, but HR already counted separately in sample data
	// Since Singles isn't provided, approximate Total Bases using Hits and extra bases
	totalPoints += ((player.Hits - player.Doubles - player.Triples - player.HomeRuns) * settings.Batting.TotalBases) + (player.Doubles * (2 * settings.Batting.TotalBases)) + (player.Triples * (3 * settings.Batting.TotalBases)) + (player.HomeRuns * (4 * settings.Batting.TotalBases))
	totalPoints += player.RBI * settings.Batting.RunsBattedIn
	totalPoints += player.Walks * settings.Batting.Walks
	totalPoints += player.Strikeouts * settings.Batting.Strikeouts
	totalPoints += player.StolenBases * settings.Batting.StolenBases

	return models.PlayerProjection{
		PlayerName:  player.Name,
		TotalPoints: totalPoints,
	}
}

// CalculateFantasyProsPitcherPoints converts FantasyPros pitcher projections to fantasy points using league settings
func CalculateFantasyProsPitcherPoints(player FanstasyProsPitcher, settings models.LeagueSettings) models.PlayerProjection {
	totalPoints := 0.0
	totalPoints += player.Strikeouts * settings.Pitching.Strikeouts
	totalPoints += player.InningsPitched * settings.Pitching.InningsPitched
	totalPoints += player.HitsAllowed * settings.Pitching.HitsAllowed
	totalPoints += player.EarnedRuns * settings.Pitching.EarnedRuns
	totalPoints += player.Walks * settings.Pitching.WalksIssued
	totalPoints += player.Wins * settings.Pitching.Wins
	totalPoints += player.Losses * settings.Pitching.Losses
	totalPoints += player.Saves * settings.Pitching.Saves

	return models.PlayerProjection{
		PlayerName:  player.Name,
		TotalPoints: totalPoints,
	}
}
