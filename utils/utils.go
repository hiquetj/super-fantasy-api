package utils

import "fmt"

// parseFloat converts a string to float64
func ParseFloat(s string) float64 {
	var f float64
	fmt.Sscanf(s, "%f", &f)
	return f
}

// parseInt converts a string to int
func ParseInt(s string) int {
	var i int
	fmt.Sscanf(s, "%d", &i)
	return i
}
