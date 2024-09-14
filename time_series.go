// Dates and time series

package main

import (
	"fmt"
	"time"
)

// Time series point
type TimeSeriesPoint struct {
	d time.Time // date
	p float64   // price on that date
}

// A Time series is a list of date/value pairs
type TimeSeries []TimeSeriesPoint

// Parse a date in "yyyy-mm-dd" format
func parseDate(ds string) time.Time {

	// An invalid date (e.g., from parse)
	// TODO: make global, or just check for year < 2000
	InvalidDate := time.Date(1970, 1, 1, 1, 0, 0, 0, time.Local)

	// If there is a space, ditch everything after the space
	if len(ds) > 10 && (ds[10] == ' ' || ds[10] == 'T') {
		ds = ds[:10]
		//fmt.Println("ds shorted to", ds)
	}

	// Parse time and return it
	t, err := time.Parse("2006-01-02", ds)
	if err != nil {
		fmt.Println("Invalid date: ", ds)
		return InvalidDate
	}
	return t
}

// Convenience function to check if date is value
func validDate(d time.Time) bool {
	return d.Year() != 1970
}

// Format a date as "yyyy-mm-dd"
func formatDate(d time.Time) string {
	return d.Format("2006-01-02")
}

// Get price on a certain date, using first date if before, last date if after,
// or linear interpolation if between dates. Assumes time series is sorted by
// date.
func priceOn(ts TimeSeries, on time.Time) float64 {

	// Return 0 if series is empty
	if len(ts) == 0 {
		return 0
	}

	// Return first price if on or before first date
	if earlier(on, ts[0].d) || sameDate(on, ts[0].d) {
		return ts[0].p
	}

	// Return last price if on or after last date
	if later(on, ts[len(ts)-1].d) || sameDate(on, ts[len(ts)-1].d) {
		return ts[len(ts)-1].p
	}

	// Otherwise search for first date that is on or after requested date
	for i := 1; i < len(ts); i++ {

		// Exact date, return price
		if sameDate(ts[i].d, on) {
			return ts[i].p
		}

		// Price is after requested date, interpolate
		if later(ts[i].d, on) {
			d0 := ts[i-1].d                // date before
			p0 := ts[i-1].p                // price before
			d1 := ts[i].d                  // date after
			p1 := ts[i].p                  //price after
			interval := d1.Sub(d0).Hours() // interval between time series dates
			delta := on.Sub(d0).Hours()    // interval between earlier date and target
			return p0 + (p1-p0)*(delta/interval)
		}
	}

	// Should never happen
	fmt.Println("Unable to interpolate date")
	return 0
}

// Determine if two dates are the same, ignoring time
func sameDate(d1, d2 time.Time) bool {
	return d1.Year() == d2.Year() && d1.Month() == d2.Month() && d1.Day() == d2.Day()
}

// Determine if date1 is earlier than date2
func earlier(d1, d2 time.Time) bool {
	if d1.Year() < d2.Year() {
		return true
	} else if d1.Year() > d2.Year() {
		return false
	} else if d1.Month() < d2.Month() {
		return true
	} else if d1.Month() > d2.Month() {
		return false
	} else {
		return d1.Day() < d2.Day()
	}
}

// Determine if date1 is later than date2
func later(d1, d2 time.Time) bool {
	return !sameDate(d1, d2) && !earlier(d1, d2)
}
