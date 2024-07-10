// Utility functions

package main

import (
	"strconv"
)

// Returns -1 if could not be parsed
func parseInt(s string) int {
	n, err := strconv.Atoi(s)
	if err != nil {
		return -1
	} else {
		return n
	}
}

// Returns -1 if could not be parsed
func parseFloat(s string) float64 {
	n, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return -1
	} else {
		return n
	}
}
