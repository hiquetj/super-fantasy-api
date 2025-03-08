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
func SaveBatterCSV(csvData string) error {
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
		player := baseball.Batter{
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
func SavePitcherCSV(csvData string) error {
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
		player := baseball.Pitcher{
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
		}
		documents = append(documents, player)
	}

	_, err = MongoInstance.Collection.InsertMany(ctx, documents)
	if err != nil {
		return fmt.Errorf("failed to insert documents: %v", err)
	}
	return nil
}
