package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func CalculateHockeyProjections(c *gin.Context) {
	// Similar implementation for hockey
	c.JSON(http.StatusOK, gin.H{"message": "Hockey projections not implemented yet"})
}