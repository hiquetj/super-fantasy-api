package handlers

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"super-fantasy-api/data/baseball"
	"super-fantasy-api/db"
	"super-fantasy-api/models"
	"super-fantasy-api/utils"

	"github.com/gin-gonic/gin"
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
				}
				playerProjection := baseball.CalculateFantasyProsBatterPoints(player, request.Settings)
				projections = append(projections, playerProjection)
			case "pitcher":
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
					Year:            request.Year,
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
			err = db.SaveFanGraphsBatterCSV(buf.String(), request.Year)
		case "pitcher":
			err = db.SaveFanGraphsPitcherCSV(buf.String(), request.Year)
		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid position: must be 'batter' or 'pitcher'"})
			return
		}
	case "fantasypros":
		switch request.Position {
		case "batter":
			err = db.SaveFantasyProsBatterCSV(buf.String(), request.Year)
		case "pitcher":
			err = db.SaveFantasyProsPitcherCSV(buf.String(), request.Year)
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
