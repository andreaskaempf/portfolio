// Utility functions

package main

import (
	"strconv"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
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

// Format a float as 2 dec places with thousands separator
// Source: https://stackoverflow.com/questions/13020308/how-to-fmt-printf-an-integer-with-thousands-comma
func formatFloat(n float64) string {
	p := message.NewPrinter(language.English)
	return p.Sprintf("%.2f", n)
}
