package baseball

import (
	"super-fantasy-api/models"
)

// struct is based on batx rankings in fangraphs
// https://www.fangraphs.com/projections?type=thebatx&stats=bat&pos=all&team=0&players=0&lg=all&z=1741170599&pageitems=30&statgroup=standard&fantasypreset=dashboard
type FangraphsBatter struct {
	Rank           int     `bson:"rank" json:"rank"`                       // # (position in list)
	Name           string  `bson:"name" json:"name"`                       // Name
	Team           string  `bson:"team" json:"team"`                       // Team
	Games          float64 `bson:"games" json:"games"`                     // G
	AtBats         float64 `bson:"at_bats" json:"at_bats"`                 // AB
	PlateApps      float64 `bson:"plate_apps" json:"plate_apps"`           // PA
	Hits           float64 `bson:"hits" json:"hits"`                       // H
	Singles        float64 `bson:"singles" json:"singles"`                 // 1B
	Doubles        float64 `bson:"doubles" json:"doubles"`                 // 2B
	Triples        float64 `bson:"triples" json:"triples"`                 // 3B
	HomeRuns       float64 `bson:"home_runs" json:"home_runs"`             // HR
	Runs           float64 `bson:"runs" json:"runs"`                       // R
	RBI            float64 `bson:"rbi" json:"rbi"`                         // RBI
	Walks          float64 `bson:"walks" json:"walks"`                     // BB
	IntWalks       float64 `bson:"int_walks" json:"int_walks"`             // IBB (Intentional Walks)
	Strikeouts     float64 `bson:"strikeouts" json:"strikeouts"`           // SO
	HitByPitch     float64 `bson:"hit_by_pitch" json:"hit_by_pitch"`       // HBP
	SacFlies       float64 `bson:"sac_flies" json:"sac_flies"`             // SF
	SacHits        float64 `bson:"sac_hits" json:"sac_hits"`               // SH
	StolenBases    float64 `bson:"stolen_bases" json:"stolen_bases"`       // SB
	CaughtStealing float64 `bson:"caught_stealing" json:"caught_stealing"` // CS
	AVG            float64 `bson:"avg" json:"avg"`                         // AVG
	Year           string  `bson:"year" json:"year"`                       // year
}

// struct is based on atc rankings in fangraphs
// https://www.fangraphs.com/fantasy-tools/auction-calculator?teams=12&lg=MLB&dollars=260&mb=1&mp=20&msp=5&mrp=5&type=pit&players=&proj=atc&split=&points=c%7C1%2C2%2C3%2C4%2C5%2C7%7C0%2C13%2C14%2C2%2C3%2C4%2C6&rep=0&drp=0&pp=SS%2C2B%2C3B%2COF%2C1B%2CC&pos=1%2C1%2C1%2C1%2C4%2C1%2C0%2C0%2C1%2C1%2C3%2C2%2C4%2C5%2C0&sort=&view=0
type FangraphsPitcher struct {
	Rank              int     `bson:"rank" json:"rank"`                               // #
	Name              string  `bson:"name" json:"name"`                               // Name
	Team              string  `bson:"team" json:"team"`                               // Team
	Wins              float64 `bson:"wins" json:"wins"`                               // W
	Losses            float64 `bson:"losses" json:"losses"`                           // L
	ERA               float64 `bson:"era" json:"era"`                                 // ERA
	Games             float64 `bson:"games" json:"games"`                             // G
	GamesStarted      float64 `bson:"games_started" json:"games_started"`             // GS
	Saves             float64 `bson:"saves" json:"saves"`                             // SV
	Holds             float64 `bson:"holds" json:"holds"`                             // HLD
	BlownSaves        float64 `bson:"blown_saves" json:"blown_saves"`                 // BS
	InningsPitched    float64 `bson:"innings_pitched" json:"innings_pitched"`         // IP
	TotalBattersFaced float64 `bson:"total_batters_faced" json:"total_batters_faced"` // TBF
	HitsAllowed       float64 `bson:"hits_allowed" json:"hits_allowed"`               // H
	RunsAllowed       float64 `bson:"runs_allowed" json:"runs_allowed"`               // R
	EarnedRuns        float64 `bson:"earned_runs" json:"earned_runs"`                 // ER
	HomeRunsAllowed   float64 `bson:"home_runs_allowed" json:"home_runs_allowed"`     // HR
	Walks             float64 `bson:"walks" json:"walks"`                             // BB
	IntWalks          float64 `bson:"int_walks" json:"int_walks"`                     // IBB
	HitByPitch        float64 `bson:"hit_by_pitch" json:"hit_by_pitch"`               // HBP
	Strikeouts        float64 `bson:"strikeouts" json:"strikeouts"`                   // SO
	Year              string  `bson:"year" json:"year"`                               // year
}

// CalculatePoints converts FanGraphs projections to fantasy points using league settings
func CalculateBatterPoints(player FangraphsBatter, settings models.LeagueSettings) models.PlayerProjection {
	// TODO: add conditionalizing for other types of league settings
	// for now, we leave what we use
	totalPoints := 0.0
	totalPoints += player.Runs * settings.Batting.RunsScored
	totalPoints += (player.Singles * settings.Batting.TotalBases) + (player.Doubles * (2 * settings.Batting.TotalBases)) + (player.Triples * (3 * settings.Batting.TotalBases)) + (player.HomeRuns * (4 * settings.Batting.TotalBases))
	totalPoints += player.RBI * settings.Batting.RunsBattedIn
	totalPoints += player.Walks * settings.Batting.Walks
	totalPoints += player.Strikeouts * settings.Batting.Strikeouts
	totalPoints += player.StolenBases * settings.Batting.StolenBases

	return models.PlayerProjection{
		PlayerName:  player.Name,
		TotalPoints: totalPoints,
	}
}

// CalculatePoints converts FanGraphs projections to fantasy points using league settings
func CalculatePitcherPoints(player FangraphsPitcher, settings models.LeagueSettings) models.PlayerProjection {
	// TODO: add conditionalizing for other types of league settings
	// for now, we leave what we use
	totalPoints := 0.0
	totalPoints += player.Strikeouts * settings.Pitching.Strikeouts
	totalPoints += player.InningsPitched * settings.Pitching.InningsPitched
	totalPoints += player.HitsAllowed * settings.Pitching.HitsAllowed
	totalPoints += player.EarnedRuns * settings.Pitching.EarnedRuns
	totalPoints += player.Walks * settings.Pitching.WalksIssued
	totalPoints += player.Wins * settings.Pitching.Wins
	totalPoints += player.Losses * settings.Pitching.Losses
	totalPoints += player.Saves * settings.Pitching.Saves
	totalPoints += player.Holds * settings.Pitching.Holds

	return models.PlayerProjection{
		PlayerName:  player.Name,
		TotalPoints: totalPoints,
	}
}
