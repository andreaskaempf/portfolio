// Some unit tests for the portfolio application

package main

import (
	"fmt"
	"testing"
	//"time"
)

// Test date parsing and formatting
func TestDateParse(t *testing.T) {

	// Parse date
	ds := "2024-07-14"
	d := parseDate(ds)
	if d.Year() != 2024 || d.Month() != 7 || d.Day() != 14 {
		fmt.Println(ds, "=>", d)
		t.Error("Error parsing date")
	}

	// Format date to string
	s := formatDate(d)
	if s != ds {
		fmt.Println(ds, "=>", s)
		t.Error("Could not format date")
	}
}

// Test date interpolation
func TestDateInterpolation(t *testing.T) {

	// Create a series of date+prices
	dates := []string{"2024-01-01", "2024-01-03", "2024-01-31", "2024-12-31"}
	prices := []float64{100.0, 110.0, 95.0, 200}
	ts := []TimeSeriesPoint{}
	for i := 0; i < len(dates); i++ {
		p := TimeSeriesPoint{parseDate(dates[i]), prices[i]}
		ts = append(ts, p)
	}

	// Test prices on exact dates in time series
	for i := 0; i < len(ts); i++ {
		p := priceOn(ts, ts[i].d)
		if p != ts[i].p {
			t.Errorf("Invalid price %f on %s", p, ts[i].d)
		}
	}

	// Test some interpolated dates
	// TODO: add more interpolations
	dd := []string{"2023-12-31", "2020-01-01", "2025-01-01", "2032-03-31",
		"2024-01-02"}
	expected := []float64{100, 100, 200, 200, 105}
	for i := 0; i < len(dd); i++ {
		d := parseDate(dd[i])
		p := priceOn(ts, d)
		if p != expected[i] {
			t.Errorf("Invalid price %f on %s", p, dd[i])
		}
	}

}
