// database.go
//
// Data model for the portfolio system, including structure definitions for
// all tables, and functions to retrieve or update data in the database.
// All database functions should be in this file.

package main

import (
	"database/sql"
	//"fmt"

	//"sort"
	//"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// Connect to database
// Don't forget to "defer db.Close() after calling this
func dbConnect() *sql.DB {
	db, err := sql.Open("sqlite3", "data.db")
	if err != nil {
		panic("dbConnect: " + err.Error())
	}
	return db
}

//----------------------------------------------------------------//
//                              STOCKS                            //
//----------------------------------------------------------------//

// Any type of security, including shares and funds
// Stocks: code, type, description, currency

// Record format for one stock
type Stock struct {
	Id       int
	Code     string
	Name     string
	Currency string
}

// Get a list of all stocks, in alphabetical order
func getStocks() []Stock {

	// Connect to database
	db := dbConnect()
	defer db.Close()

	// Execute query to get all stocks, in alphabetical order
	rows, err := db.Query("select id, code, name, currency from stock order by code")
	if err != nil {
		panic("getStocks query: " + err.Error())
	}
	defer rows.Close()

	// Collect into a list
	ss := []Stock{}
	for rows.Next() {
		s := Stock{}
		err := rows.Scan(&s.Id, &s.Code, &s.Name, &s.Currency)
		if err != nil {
			panic("getStocks next: " + err.Error())
		}
		ss = append(ss, s)
	}
	if rows.Err() != nil {
		panic("getStocks exit: " + err.Error())
	}

	// Return list
	return ss
}

// Get one stock by id
func getStock(sid int) *Stock {

	// Connect to database
	db := dbConnect()
	defer db.Close()

	// Find stock, return nil if not found
	s := Stock{}
	q := "select id, code, name, currency from stock where id = $1"
	err := db.QueryRow(q, sid).Scan(&s.Id, &s.Code, &s.Name, &s.Currency)
	if err != nil {
		return nil
	}

	return &s
}

// Update an existing stock, or add new
func addUpdateStock(s *Stock) {

	// Connect to database
	db := dbConnect()
	defer db.Close()

	// Attempt insert or update
	var err error
	if s.Id == 0 {
		q := "insert into stock(code, name, currency) values ($1, $2, $3)"
		_, err = db.Exec(q, s.Code, s.Name, s.Currency)
	} else {
		q := "update stock set code = $1, name = $2, currency = $3 where id = $4"
		_, err = db.Exec(q, s.Code, s.Name, s.Currency, s.Id)
	}

	// Check for error
	if err != nil {
		panic("addUpdateStock: " + err.Error())
	}
}

// Delete a stock by ID
// TODO: also delete all child records
func deleteStock(sid int) {

	db := dbConnect()
	defer db.Close()

	_, err := db.Exec("delete from stock where id = $1", sid)
	if err != nil {
		panic("deleteStock: " + err.Error())
	}
}

//----------------------------------------------------------------//
//                              PRICES                            //
//----------------------------------------------------------------//

// Price of a stock on a certain date, in its currency

// Record format for a price
type Price struct {
	Id       int       // the id of this price record
	Date     time.Time // the date for this price
	Stock    int       // id of the stock this price is for
	Price    float64   // price on this date, in the stock's currency
	Comments string    // any comments
}

// Get the prices for a stock
func getPrices(sid int) []Price {

	// Connect to database
	db := dbConnect()
	defer db.Close()

	// Execute query to get all prices for this stock, in date order
	rows, err := db.Query("select id, pdate, price, comments from price where stock_id = $1 order by date", sid)
	if err != nil {
		panic("getPrices query: " + err.Error())
	}
	defer rows.Close()

	// Collect into a list
	pp := []Price{}
	var ds string // buffer for reading date
	for rows.Next() {
		p := Price{}
		err := rows.Scan(&p.Id, &ds, &p.Price, &p.Comments)
		if err != nil {
			panic("getPrices next: " + err.Error())
		}
		// TODO: convert date
		pp = append(pp, p)
	}
	if rows.Err() != nil {
		panic("getPricess exit: " + err.Error())
	}

	// Return list
	return pp
}
