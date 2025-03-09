package db

import (
	"context"
	"encoding/csv"
	"fmt"
	"strings"
	"time"

	"super-fantasy-api/data/baseball"
	"super-fantasy-api/utils"

	"go.mongodb.org/mongo-driver/bson"
)

// SaveBatterCSV parses and saves generic batter CSV data to MongoDB
func SaveFanGraphsBatterCSV(csvData string, year string) error {
	reader := csv.NewReader(strings.NewReader(csvData))
	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("failed to parse CSV: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Optionally clear existing batter data with a specific source if needed
	_, err = MongoInstance.Collection.DeleteMany(ctx, bson.M{"source": "batter-stats"})
	if err != nil {
		return fmt.Errorf("failed to clear batter data: %v", err)
	}

	var documents []interface{}
	for i, record := range records {
		if i == 0 {
			continue // Skip header
		}
		player := baseball.FangraphsBatter{
			Rank:           utils.ParseInt(record[0]),
			Name:           record[1],
			Team:           record[2],
			Games:          utils.ParseFloat(record[3]),
			AtBats:         utils.ParseFloat(record[4]),
			PlateApps:      utils.ParseFloat(record[5]),
			Hits:           utils.ParseFloat(record[6]),
			Singles:        utils.ParseFloat(record[7]),
			Doubles:        utils.ParseFloat(record[8]),
			Triples:        utils.ParseFloat(record[9]),
			HomeRuns:       utils.ParseFloat(record[10]),
			Runs:           utils.ParseFloat(record[11]),
			RBI:            utils.ParseFloat(record[12]),
			Walks:          utils.ParseFloat(record[13]),
			IntWalks:       utils.ParseFloat(record[14]),
			Strikeouts:     utils.ParseFloat(record[15]),
			HitByPitch:     utils.ParseFloat(record[16]),
			SacFlies:       utils.ParseFloat(record[17]),
			SacHits:        utils.ParseFloat(record[18]),
			StolenBases:    utils.ParseFloat(record[19]),
			CaughtStealing: utils.ParseFloat(record[20]),
			AVG:            utils.ParseFloat(record[21]),
			Year:           year,
		}
		documents = append(documents, player)
	}

	_, err = MongoInstance.Collection.InsertMany(ctx, documents)
	if err != nil {
		return fmt.Errorf("failed to insert documents: %v", err)
	}
	return nil
}

// SavePitcherCSV parses and saves pitcher CSV data to MongoDB
func SaveFanGraphsPitcherCSV(csvData string, year string) error {
	reader := csv.NewReader(strings.NewReader(csvData))
	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("failed to parse CSV: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	_, err = MongoInstance.Collection.DeleteMany(ctx, bson.M{"source": "pitcher-stats"})
	if err != nil {
		return fmt.Errorf("failed to clear pitcher data: %v", err)
	}

	var documents []interface{}
	for i, record := range records {
		if i == 0 {
			continue // Skip header
		}
		player := baseball.FangraphsPitcher{
			Rank:              utils.ParseInt(record[0]),    // #
			Name:              record[1],                    // Name
			Team:              record[2],                    // Team
			Wins:              utils.ParseFloat(record[3]),  // W
			Losses:            utils.ParseFloat(record[4]),  // L
			ERA:               utils.ParseFloat(record[5]),  // ERA
			Games:             utils.ParseFloat(record[6]),  // G
			GamesStarted:      utils.ParseFloat(record[7]),  // GS
			Saves:             utils.ParseFloat(record[8]),  // SV
			Holds:             utils.ParseFloat(record[9]),  // HLD
			BlownSaves:        utils.ParseFloat(record[10]), // BS
			InningsPitched:    utils.ParseFloat(record[11]), // IP
			TotalBattersFaced: utils.ParseFloat(record[12]), // TBF
			HitsAllowed:       utils.ParseFloat(record[13]), // H
			RunsAllowed:       utils.ParseFloat(record[14]), // R
			EarnedRuns:        utils.ParseFloat(record[15]), // ER
			HomeRunsAllowed:   utils.ParseFloat(record[16]), // HR
			Walks:             utils.ParseFloat(record[17]), // BB
			IntWalks:          utils.ParseFloat(record[18]), // IBB
			HitByPitch:        utils.ParseFloat(record[19]), // HBP
			Strikeouts:        utils.ParseFloat(record[20]), // SO
			Year:              year,
		}
		documents = append(documents, player)
	}

	_, err = MongoInstance.Collection.InsertMany(ctx, documents)
	if err != nil {
		return fmt.Errorf("failed to insert documents: %v", err)
	}
	return nil
}

// SaveFantasyProsBatterCSV saves FantasyPros batter CSV data to MongoDB
func SaveFantasyProsBatterCSV(csvData string, year string) error {
	reader := csv.NewReader(strings.NewReader(csvData))
	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("failed to parse CSV: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	_, err = MongoInstance.Collection.DeleteMany(ctx, bson.M{"source": "fantasypros-batter"})
	if err != nil {
		return fmt.Errorf("failed to clear FantasyPros batter data: %v", err)
	}

	var documents []interface{}
	for i, record := range records {
		if i == 0 {
			continue // Skip header
		}
		player := baseball.FantasyProsBatter{
			Name:        record[0],
			Team:        record[1],
			Positions:   record[2],
			AtBats:      utils.ParseFloat(record[3]),
			Runs:        utils.ParseFloat(record[4]),
			HomeRuns:    utils.ParseFloat(record[5]),
			RBI:         utils.ParseFloat(record[6]),
			StolenBases: utils.ParseFloat(record[7]),
			AVG:         utils.ParseFloat(record[8]),
			OBP:         utils.ParseFloat(record[9]),
			Hits:        utils.ParseFloat(record[10]),
			Doubles:     utils.ParseFloat(record[11]),
			Triples:     utils.ParseFloat(record[12]),
			Walks:       utils.ParseFloat(record[13]),
			Strikeouts:  utils.ParseFloat(record[14]),
			SLG:         utils.ParseFloat(record[15]),
			OPS:         utils.ParseFloat(record[16]),
			Year:        year,
		}
		documents = append(documents, player)
	}

	_, err = MongoInstance.Collection.InsertMany(ctx, documents)
	if err != nil {
		return fmt.Errorf("failed to insert documents: %v", err)
	}
	return nil
}

// SaveFantasyProsPitcherCSV saves FantasyPros pitcher CSV data to MongoDB
func SaveFantasyProsPitcherCSV(csvData string, year string) error {
	reader := csv.NewReader(strings.NewReader(csvData))
	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("failed to parse CSV: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	_, err = MongoInstance.Collection.DeleteMany(ctx, bson.M{"source": "fantasypros-pitcher"})
	if err != nil {
		return fmt.Errorf("failed to clear FantasyPros pitcher data: %v", err)
	}

	var documents []interface{}
	for i, record := range records {
		if i == 0 {
			continue // Skip header
		}
		player := baseball.FanstasyProsPitcher{
			Name:            record[0],
			Team:            record[1],
			Positions:       record[2],
			InningsPitched:  utils.ParseFloat(record[3]),
			Strikeouts:      utils.ParseFloat(record[4]),
			Wins:            utils.ParseFloat(record[5]),
			Saves:           utils.ParseFloat(record[6]),
			ERA:             utils.ParseFloat(record[7]),
			WHIP:            utils.ParseFloat(record[8]),
			EarnedRuns:      utils.ParseFloat(record[9]),
			HitsAllowed:     utils.ParseFloat(record[10]),
			Walks:           utils.ParseFloat(record[11]),
			HomeRunsAllowed: utils.ParseFloat(record[12]),
			Games:           utils.ParseFloat(record[13]),
			GamesStarted:    utils.ParseFloat(record[14]),
			Losses:          utils.ParseFloat(record[15]),
			CompleteGames:   utils.ParseFloat(record[16]),
			Year:            year,
		}
		documents = append(documents, player)
	}

	_, err = MongoInstance.Collection.InsertMany(ctx, documents)
	if err != nil {
		return fmt.Errorf("failed to insert documents: %v", err)
	}
	return nil
}
