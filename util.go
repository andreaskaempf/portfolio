// Utility functions

package main

import (
	//"fmt"
	//"sort"
	"strconv"
	//"strings"
	//"time"
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

/*
// Dates come back as "yyyy-mm-ddT00:00:00Z", so strip off the time
func cleanDate(d string) string {
	if len(d) == 20 && d[10] == 'T' {
		return d[:10]
	} else {
		return d
	}
}

// Extract keys from a map, sorted
func getSortedKeys(dict map[string]int) []string {
	keys := []string{}
	for k, _ := range dict {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })
	return keys
}

// Is string in a list? (NOT USED?)
func in(s string, list []string) bool {
	for _, x := range list {
		if x == s {
			return true
		}
	}
	return false
}

// Parse a "yyyy-mm-dd" date, return true if valid
func validDate(s string) bool {
	_, err := time.Parse("2006-01-02", s)
	return err == nil
}

// Parse a "hh:mm" time, return true if valid
func validTime(s string) bool {
	_, err := time.Parse("15:04", s)
	return err == nil
}

// Check if a valid e-mail address, does not check letters, since
// these might be international
func validEmail(email string) bool {
	if len(email) < 5 {
		return false
	}
	var at, period int
	for i, c := range email {
		if c == '@' { // position of at, only one allowed
			if at > 0 {
				return false
			}
			at = i
		} else if c == '.' {
			period = i
		}
	}

	// At sign must be present, can't be at beginning or end
	if at == 0 || at == len(email)-1 {
		return false
	}

	// Last period must be present, after the at sign, and not at the end
	if period == 0 || period < at || period == len(email)-1 {
		return false
	}

	// Otherwise looks okay
	return true
}

// Try to guess first & last name from an e-mail. If the part before
// at sign includes a period or underscore, assume these are first and
// last names. Otherwise, assume the whole thing is first name and
// use "Unknown" for last name.
func guessNames(email string) (string, string) {

	// Find the text before the at sign
	at := strings.Index(email, "@")
	if at <= 0 {
		return email, "(Unknown)"
	}
	email = email[:at]

	// Try splitting at various delimiters
	for i := 1; i < len(email)-1; i++ {
		if email[i] == '.' || email[i] == '_' {
			fname := email[:i]
			lname := email[i+1:]
			return fname, lname
		}
	}

	// Return text before at as first name if no delimiter found
	return email, "(Unknown)"
}

// Given a string of "yyyy-mm-dd" format, add one day to it
func addDay(ds string) string {

	// Parse the date
	d := parseDate(ds)
	if d.Year() < 2000 {
		fmt.Println("addDay: invalid date", ds)
		return ds
	}

	// Add a day
	d = d.AddDate(0, 0, 1)

	// Return as "yyyy-mm-dd"
	return formatDate(d)
}

// Parse a "yyyy-mm-dd" date, return Jan 1 1900 if error in source date
func parseDate(ds string) time.Time {
	if ds == "" {
		return time.Date(1900, 1, 1, 0, 0, 0, 0, time.Now().Location())
	}
	t, err := time.Parse("2006-1-2", ds)
	if err != nil {
		fmt.Println("parseDate: invalid date", ds, err.Error())
		return time.Date(1900, 1, 1, 0, 0, 0, 0, time.Now().Location())
	}
	return t
}

// Format a date as "yyyy-mm-dd"
func formatDate(d time.Time) string {
	return d.Format("2006-1-2")
}

// Given a date string, return "4 d" if 4 days ago, etc.
func age(ds string) string {

	// If empty, return empty string
	if len(ds) == 0 {
		return ""
	}

	// Parse the "yyyy-mm-dd hh:mm"
	d, err := time.Parse("2006-1-2 15:04", ds)
	if err != nil {
		fmt.Println("age: invalid date time", d, err.Error())
		return "(err)"
	}

	// Calculate the age from now
	now := time.Now()
	delta := now.Sub(d)

	// Return appropriate string
	if delta.Hours() < 1 {
		return fmt.Sprintf("%dm", int(delta.Minutes()))
	} else if delta.Hours() <= 24 {
		return fmt.Sprintf("%dh", int(delta.Hours()))
	} else {
		return fmt.Sprintf("%dd", int(delta.Hours()/24))
	}
}

// Convert a time "hh:mm" to floating point,  e.g., "12:30" -> 12.5
func timeToDec(t string) float64 {

	// Separate hours and mins
	hm := strings.Split(t, ":")
	if len(hm) != 2 {
		fmt.Println("Invalid time string, returning 12.00:", t)
		return 12.0
	}

	// Parse hours and mins
	h := parseInt(hm[0])
	m := parseInt(hm[1])
	if h < 0 || h > 23 || m < 0 || m > 59 {
		fmt.Println("Invalid time values, returning 12.00:", t)
		return 12.0
	}

	// Return as floating point
	return float64(h) + float64(m)/60.0
}
*/
