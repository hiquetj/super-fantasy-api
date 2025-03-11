package handlers

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"super-fantasy-api/data/baseball"
	"super-fantasy-api/db"
	"super-fantasy-api/models"
	"super-fantasy-api/utils"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

// CalculateBaseballProjections handles fetching projections from MongoDB
func CalculateBaseballProjections(c *gin.Context) {
	// Get the CSV file from the form
	file, _, err := c.Request.FormFile("csv")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get CSV file: " + err.Error()})
		return
	}
	defer file.Close()

	// Read CSV into a buffer
	var buf bytes.Buffer
	_, err = io.Copy(&buf, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read CSV: " + err.Error()})
		return
	}

	// Parse CSV
	reader := csv.NewReader(strings.NewReader(buf.String()))
	records, err := reader.ReadAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse CSV: " + err.Error()})
		return
	}

	// Get league settings from form field
	settingsStr := c.Request.FormValue("settings") // Use FormValue instead of PostForm
	if settingsStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing or empty settings field"})
		return
	}

	// Debug: Log the raw settings string to verify
	c.Request.ParseMultipartForm(10 << 20) // Parse form if not already done (10MB max)
	c.Writer.WriteString("Debug: settingsStr = " + settingsStr + "\n")

	var request models.ProjectionRequest
	if err := json.Unmarshal([]byte(settingsStr), &request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid league settings format: " + err.Error()})
		return
	}

	// Process CSV records into projections (assuming Batter format for now)
	var projections []models.PlayerProjection
	for i, record := range records {
		if i == 0 {
			continue // Skip header
		}

		// Dispatch to appropriate save function based on source and position
		switch request.Source {
		case "fangraphs":
			switch request.Position {
			case "batter":
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
					Year:           request.Year,
					Position:       request.Position,
				}

				playerProjection := baseball.CalculateBatterPoints(player, request.Settings)
				projections = append(projections, playerProjection)
			case "pitcher":
				player := baseball.FangraphsPitcher{
					Rank:              utils.ParseInt(record[0]),
					Name:              record[1],
					Team:              record[2],
					Wins:              utils.ParseFloat(record[3]),
					Losses:            utils.ParseFloat(record[4]),
					ERA:               utils.ParseFloat(record[5]),
					Games:             utils.ParseFloat(record[6]),
					GamesStarted:      utils.ParseFloat(record[7]),
					Saves:             utils.ParseFloat(record[8]),
					Holds:             utils.ParseFloat(record[9]),
					BlownSaves:        utils.ParseFloat(record[10]),
					InningsPitched:    utils.ParseFloat(record[11]),
					TotalBattersFaced: utils.ParseFloat(record[12]),
					HitsAllowed:       utils.ParseFloat(record[13]),
					RunsAllowed:       utils.ParseFloat(record[14]),
					EarnedRuns:        utils.ParseFloat(record[15]),
					HomeRunsAllowed:   utils.ParseFloat(record[16]),
					Walks:             utils.ParseFloat(record[17]),
					IntWalks:          utils.ParseFloat(record[18]),
					HitByPitch:        utils.ParseFloat(record[19]),
					Strikeouts:        utils.ParseFloat(record[20]),
					Year:              request.Year,
					Position:          request.Position,
				}

				playerProjection := baseball.CalculatePitcherPoints(player, request.Settings)
				projections = append(projections, playerProjection)
			default:
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid position: must be 'batter' or 'pitcher'"})
				return
			}
		case "fantasypros":
			switch request.Position {
			case "batter":
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
					Year:        request.Year,
					Position:    request.Position,
				}
				playerProjection := baseball.CalculateFantasyProsBatterPoints(player, request.Settings)
				projections = append(projections, playerProjection)
			case "pitcher":
				player := baseball.FantasyProsPitcher{
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
					Year:            request.Year,
					Position:        request.Position,
				}
				playerProjection := baseball.CalculateFantasyProsPitcherPoints(player, request.Settings)
				projections = append(projections, playerProjection)
			default:
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid position: must be 'batter' or 'pitcher'"})
				return
			}
		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid source: must be 'fangraphs' or 'fantasypros'"})
			return
		}
	}

	// Return projections as JSON
	c.JSON(http.StatusOK, gin.H{
		"projections":       projections,
		"projection_source": request.ProjectionName,
	})
}

// UploadCSV handles CSV upload and saves to MongoDB based on source and position
func UploadCSV(c *gin.Context) {
	// Get the CSV file from the form
	file, _, err := c.Request.FormFile("csv")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get CSV file: " + err.Error()})
		return
	}
	defer file.Close()

	// Read CSV into a buffer
	var buf bytes.Buffer
	_, err = io.Copy(&buf, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read CSV: " + err.Error()})
		return
	}

	// Get upload settings from form field
	settingsStr := c.Request.FormValue("settings")
	if settingsStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing or empty settings field"})
		return
	}

	var request models.UploadRequest
	if err := json.Unmarshal([]byte(settingsStr), &request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid settings format: " + err.Error()})
		return
	}

	// Dispatch to appropriate save function based on source and position
	switch request.Source {
	case "fangraphs":
		switch request.Position {
		case "batter":
			err = db.SaveFanGraphsBatterCSV(buf.String(), request.Year, request.Suffix, request.Position)
		case "pitcher":
			err = db.SaveFanGraphsPitcherCSV(buf.String(), request.Year, request.Suffix, request.Position)
		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid position: must be 'batter' or 'pitcher'"})
			return
		}
	case "fantasypros":
		switch request.Position {
		case "batter":
			err = db.SaveFantasyProsBatterCSV(buf.String(), request.Year, request.Position)
		case "pitcher":
			err = db.SaveFantasyProsPitcherCSV(buf.String(), request.Year, request.Position)
		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid position: must be 'batter' or 'pitcher'"})
			return
		}
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid source: must be 'fangraphs' or 'fantasypros'"})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save CSV: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "CSV uploaded and saved successfully"})
}

func ExportPlayerPointsCSV(c *gin.Context) {
	// Get league settings from form field
	settingsStr := c.Request.FormValue("settings")
	if settingsStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing or empty settings field"})
		return
	}

	var request models.ProjectionRequest
	if err := json.Unmarshal([]byte(settingsStr), &request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid league settings format: " + err.Error()})
		return
	}

	// Query all documents from the Baseball collection
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cursor, err := db.MongoInstance.Collection.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query MongoDB: " + err.Error()})
		return
	}
	defer cursor.Close(ctx)

	// Map to store player points: key is "Name:Position", value is a map of source to points
	playerPoints := make(map[string]map[string]float64)

	// Process each document
	for cursor.Next(ctx) {
		var doc bson.M
		if err := cursor.Decode(&doc); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode document: " + err.Error()})
			return
		}

		rawName, _ := doc["name"].(string)
		name := utils.NormalizeName(rawName)
		source, _ := doc["source"].(string)
		var position string
		var points float64

		// Determine position and calculate points based on source
		switch source {
		case "fantasypros":
			if _, isBatter := doc["at_bats"]; isBatter {
				position = "Batter"
				player := baseball.FantasyProsBatter{
					Name:        name,
					Team:        utils.GetString(doc, "team"),
					Positions:   utils.GetString(doc, "positions"),
					AtBats:      utils.GetFloat64(doc, "at_bats"),
					Runs:        utils.GetFloat64(doc, "runs"),
					HomeRuns:    utils.GetFloat64(doc, "home_runs"),
					RBI:         utils.GetFloat64(doc, "rbi"),
					StolenBases: utils.GetFloat64(doc, "stolen_bases"),
					AVG:         utils.GetFloat64(doc, "avg"),
					OBP:         utils.GetFloat64(doc, "obp"),
					Hits:        utils.GetFloat64(doc, "hits"),
					Doubles:     utils.GetFloat64(doc, "doubles"),
					Triples:     utils.GetFloat64(doc, "triples"),
					Walks:       utils.GetFloat64(doc, "walks"),
					Strikeouts:  utils.GetFloat64(doc, "strikeouts"),
					SLG:         utils.GetFloat64(doc, "slg"),
					OPS:         utils.GetFloat64(doc, "ops"),
					Year:        utils.GetString(doc, "year"),
					Source:      source,
					Position:    position,
				}
				proj := baseball.CalculateFantasyProsBatterPoints(player, request.Settings)
				points = proj.TotalPoints
			} else if _, isPitcher := doc["innings_pitched"]; isPitcher {
				position = "Pitcher"
				player := baseball.FantasyProsPitcher{
					Name:            name,
					Team:            utils.GetString(doc, "team"),
					Positions:       utils.GetString(doc, "positions"),
					InningsPitched:  utils.GetFloat64(doc, "innings_pitched"),
					Strikeouts:      utils.GetFloat64(doc, "strikeouts"),
					Wins:            utils.GetFloat64(doc, "wins"),
					Saves:           utils.GetFloat64(doc, "saves"),
					ERA:             utils.GetFloat64(doc, "era"),
					WHIP:            utils.GetFloat64(doc, "whip"),
					EarnedRuns:      utils.GetFloat64(doc, "earned_runs"),
					HitsAllowed:     utils.GetFloat64(doc, "hits_allowed"),
					Walks:           utils.GetFloat64(doc, "walks"),
					HomeRunsAllowed: utils.GetFloat64(doc, "home_runs_allowed"),
					Games:           utils.GetFloat64(doc, "games"),
					GamesStarted:    utils.GetFloat64(doc, "games_started"),
					Losses:          utils.GetFloat64(doc, "losses"),
					CompleteGames:   utils.GetFloat64(doc, "complete_games"),
					Year:            utils.GetString(doc, "year"),
					Source:          source,
					Position:        position,
				}
				proj := baseball.CalculateFantasyProsPitcherPoints(player, request.Settings)
				points = proj.TotalPoints
			}
		case "fangraphs_atc":
			position = "Pitcher"
			player := baseball.FangraphsPitcher{
				Rank:              utils.GetInt(doc, "rank"),
				Name:              name,
				Team:              utils.GetString(doc, "team"),
				Wins:              utils.GetFloat64(doc, "wins"),
				Losses:            utils.GetFloat64(doc, "losses"),
				ERA:               utils.GetFloat64(doc, "era"),
				Games:             utils.GetFloat64(doc, "games"),
				GamesStarted:      utils.GetFloat64(doc, "games_started"),
				Saves:             utils.GetFloat64(doc, "saves"),
				Holds:             utils.GetFloat64(doc, "holds"),
				BlownSaves:        utils.GetFloat64(doc, "blown_saves"),
				InningsPitched:    utils.GetFloat64(doc, "innings_pitched"),
				TotalBattersFaced: utils.GetFloat64(doc, "total_batters_faced"),
				HitsAllowed:       utils.GetFloat64(doc, "hits_allowed"),
				RunsAllowed:       utils.GetFloat64(doc, "runs_allowed"),
				EarnedRuns:        utils.GetFloat64(doc, "earned_runs"),
				HomeRunsAllowed:   utils.GetFloat64(doc, "home_runs_allowed"),
				Walks:             utils.GetFloat64(doc, "walks"),
				IntWalks:          utils.GetFloat64(doc, "int_walks"),
				HitByPitch:        utils.GetFloat64(doc, "hit_by_pitch"),
				Strikeouts:        utils.GetFloat64(doc, "strikeouts"),
				Year:              utils.GetString(doc, "year"),
				Source:            source,
				Position:          position,
			}
			proj := baseball.CalculatePitcherPoints(player, request.Settings)
			points = proj.TotalPoints
		case "fangraphs_batx":
			position = "Batter"
			player := baseball.FangraphsBatter{
				Rank:           utils.GetInt(doc, "rank"),
				Name:           name,
				Team:           utils.GetString(doc, "team"),
				Games:          utils.GetFloat64(doc, "games"),
				AtBats:         utils.GetFloat64(doc, "at_bats"),
				PlateApps:      utils.GetFloat64(doc, "plate_apps"),
				Hits:           utils.GetFloat64(doc, "hits"),
				Singles:        utils.GetFloat64(doc, "singles"),
				Doubles:        utils.GetFloat64(doc, "doubles"),
				Triples:        utils.GetFloat64(doc, "triples"),
				HomeRuns:       utils.GetFloat64(doc, "home_runs"),
				Runs:           utils.GetFloat64(doc, "runs"),
				RBI:            utils.GetFloat64(doc, "rbi"),
				Walks:          utils.GetFloat64(doc, "walks"),
				IntWalks:       utils.GetFloat64(doc, "int_walks"),
				Strikeouts:     utils.GetFloat64(doc, "strikeouts"),
				HitByPitch:     utils.GetFloat64(doc, "hit_by_pitch"),
				SacFlies:       utils.GetFloat64(doc, "sac_flies"),
				SacHits:        utils.GetFloat64(doc, "sac_hits"),
				StolenBases:    utils.GetFloat64(doc, "stolen_bases"),
				CaughtStealing: utils.GetFloat64(doc, "caught_stealing"),
				AVG:            utils.GetFloat64(doc, "avg"),
				Year:           utils.GetString(doc, "year"),
				Source:         source,
				Position:       position,
			}
			proj := baseball.CalculateBatterPoints(player, request.Settings)
			points = proj.TotalPoints
		case "fangraphs_steamer":
			if _, isBatter := doc["at_bats"]; isBatter {
				position = "Batter"
				player := baseball.FangraphsBatter{
					Rank:           utils.GetInt(doc, "rank"),
					Name:           name,
					Team:           utils.GetString(doc, "team"),
					Games:          utils.GetFloat64(doc, "games"),
					AtBats:         utils.GetFloat64(doc, "at_bats"),
					PlateApps:      utils.GetFloat64(doc, "plate_apps"),
					Hits:           utils.GetFloat64(doc, "hits"),
					Singles:        utils.GetFloat64(doc, "singles"),
					Doubles:        utils.GetFloat64(doc, "doubles"),
					Triples:        utils.GetFloat64(doc, "triples"),
					HomeRuns:       utils.GetFloat64(doc, "home_runs"),
					Runs:           utils.GetFloat64(doc, "runs"),
					RBI:            utils.GetFloat64(doc, "rbi"),
					Walks:          utils.GetFloat64(doc, "walks"),
					IntWalks:       utils.GetFloat64(doc, "int_walks"),
					Strikeouts:     utils.GetFloat64(doc, "strikeouts"),
					HitByPitch:     utils.GetFloat64(doc, "hit_by_pitch"),
					SacFlies:       utils.GetFloat64(doc, "sac_flies"),
					SacHits:        utils.GetFloat64(doc, "sac_hits"),
					StolenBases:    utils.GetFloat64(doc, "stolen_bases"),
					CaughtStealing: utils.GetFloat64(doc, "caught_stealing"),
					AVG:            utils.GetFloat64(doc, "avg"),
					Year:           utils.GetString(doc, "year"),
					Source:         source,
					Position:       position,
				}
				proj := baseball.CalculateBatterPoints(player, request.Settings)
				points = proj.TotalPoints
			} else if _, isPitcher := doc["innings_pitched"]; isPitcher {
				position = "Pitcher"
				player := baseball.FangraphsPitcher{
					Rank:              utils.GetInt(doc, "rank"),
					Name:              name,
					Team:              utils.GetString(doc, "team"),
					Wins:              utils.GetFloat64(doc, "wins"),
					Losses:            utils.GetFloat64(doc, "losses"),
					ERA:               utils.GetFloat64(doc, "era"),
					Games:             utils.GetFloat64(doc, "games"),
					GamesStarted:      utils.GetFloat64(doc, "games_started"),
					Saves:             utils.GetFloat64(doc, "saves"),
					Holds:             utils.GetFloat64(doc, "holds"),
					BlownSaves:        utils.GetFloat64(doc, "blown_saves"),
					InningsPitched:    utils.GetFloat64(doc, "innings_pitched"),
					TotalBattersFaced: utils.GetFloat64(doc, "total_batters_faced"),
					HitsAllowed:       utils.GetFloat64(doc, "hits_allowed"),
					RunsAllowed:       utils.GetFloat64(doc, "runs_allowed"),
					EarnedRuns:        utils.GetFloat64(doc, "earned_runs"),
					HomeRunsAllowed:   utils.GetFloat64(doc, "home_runs_allowed"),
					Walks:             utils.GetFloat64(doc, "walks"),
					IntWalks:          utils.GetFloat64(doc, "int_walks"),
					HitByPitch:        utils.GetFloat64(doc, "hit_by_pitch"),
					Strikeouts:        utils.GetFloat64(doc, "strikeouts"),
					Year:              utils.GetString(doc, "year"),
					Source:            source,
				}
				proj := baseball.CalculatePitcherPoints(player, request.Settings)
				points = proj.TotalPoints
			}
		default:
			continue // Skip unknown sources
		}

		// Key for map: "Name:Position"
		key := fmt.Sprintf("%s:%s", name, position)
		if _, exists := playerPoints[key]; !exists {
			playerPoints[key] = make(map[string]float64)
		}

		// Map source to column name
		var column string
		switch source {
		case "fantasypros":
			column = "FantasyPros"
		case "fangraphs_atc":
			column = "FangraphsATC"
		case "fangraphs_batx":
			column = "FangraphsBatX"
		case "fangraphs_steamer":
			column = "Steamer"
		}
		playerPoints[key][column] = points
	}

	if err := cursor.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cursor error: " + err.Error()})
		return
	}

	// Generate CSV
	var csvBuf bytes.Buffer
	writer := csv.NewWriter(&csvBuf)

	// Write headers (added "Aggregate")
	headers := []string{"Player", "Position", "FantasyPros", "FangraphsATC", "FangraphsBatX", "Steamer", "Aggregate"}
	if err := writer.Write(headers); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to write CSV headers: " + err.Error()})
		return
	}

	// Write rows with aggregate
	for key, pointsMap := range playerPoints {
		parts := strings.Split(key, ":")
		name := parts[0]
		position := parts[1]

		// Calculate aggregate
		var sum float64
		var count int
		if val := pointsMap["FantasyPros"]; val != 0 {
			sum += val
			count++
		}
		if position == "Pitcher" && pointsMap["FangraphsATC"] != 0 {
			sum += pointsMap["FangraphsATC"]
			count++
		}
		if position == "Batter" && pointsMap["FangraphsBatX"] != 0 {
			sum += pointsMap["FangraphsBatX"]
			count++
		}
		if val := pointsMap["Steamer"]; val != 0 {
			sum += val
			count++
		}
		aggregate := 0.0
		if count > 0 {
			aggregate = sum / float64(count)
		}

		row := []string{
			name,
			position,
			fmt.Sprintf("%.1f", pointsMap["FantasyPros"]),
			fmt.Sprintf("%.1f", pointsMap["FangraphsATC"]),
			fmt.Sprintf("%.1f", pointsMap["FangraphsBatX"]),
			fmt.Sprintf("%.1f", pointsMap["Steamer"]),
			fmt.Sprintf("%.1f", aggregate),
		}
		// Replace "0.0" with "" for missing values (except Aggregate)
		for i := 2; i < len(row)-1; i++ { // Skip Aggregate column
			if row[i] == "0.0" {
				row[i] = ""
			}
		}
		if err := writer.Write(row); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to write CSV row: " + err.Error()})
			return
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to flush CSV writer: " + err.Error()})
		return
	}

	// Set response headers and send CSV
	c.Header("Content-Disposition", "attachment; filename=player_points.csv")
	c.Header("Content-Type", "text/csv")
	c.Data(http.StatusOK, "text/csv", csvBuf.Bytes())
}
