// database.go
//
// Data model for the portfolio system, including structure definitions for
// all tables, and functions to retrieve or update data in the database.
// All database functions should be in this file.

package main

import (
	"database/sql"
	//"fmt"
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
//                          STOCK PRICES                          //
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

// Get price by price ID
func getPrice(pid int) *Price {

	// Connect to database
	db := dbConnect()
	defer db.Close()

	// Find price, return nil if not found
	p := Price{}
	q := "select id, stock_id, pdate, price, comments from price where id = $1"
	err := db.QueryRow(q, pid).Scan(&p.Id, &p.Stock, &p.Date, &p.Price, &p.Comments)
	if err != nil {
		return nil
	}

	return &p
}

// Get all prices for a stock
func getPrices(sid int) []Price {

	// Connect to database
	db := dbConnect()
	defer db.Close()

	// Execute query to get all prices for this stock, in date order
	rows, err := db.Query("select id, pdate, price, comments from price where stock_id = $1 order by pdate desc", sid)
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
		//fmt.Println(p)
		// TODO: convert date
		p.Date = parseDate(ds)
		pp = append(pp, p)
	}
	if rows.Err() != nil {
		panic("getPricess exit: " + err.Error())
	}

	// Return list
	return pp
}

// Update an existing price, or add new
func addUpdatePrice(p *Price) {

	// Connect to database
	db := dbConnect()
	defer db.Close()

	// Attempt insert or update
	var err error
	if p.Id == 0 {
		q := "insert into price(stock_id, pdate, price, comments) values ($1, $2, $3, $4)"
		_, err = db.Exec(q, p.Stock, p.Date, p.Price, p.Comments)
	} else {
		q := "update price set pdate = $1, price = $2, comments = $3 where id = $4"
		_, err = db.Exec(q, p.Date, p.Price, p.Comments, p.Id)
	}

	// Check for error
	if err != nil {
		panic("addUpdatePrice: " + err.Error())
	}
}

//----------------------------------------------------------------//
//                          CURRENCIES                            //
//----------------------------------------------------------------//

// Record format for one currency
type Currency struct {
	Id   int
	Code string
	Name string
}

// Get a list of all currencys, in alphabetical order
func getCurrencies() []Currency {

	// Connect to database
	db := dbConnect()
	defer db.Close()

	// Execute query to get all currencys, in alphabetical order
	rows, err := db.Query("select id, code, name from currency order by code")
	if err != nil {
		panic("getCurrencies query: " + err.Error())
	}
	defer rows.Close()

	// Collect into a list
	curs := []Currency{}
	for rows.Next() {
		cur := Currency{}
		err := rows.Scan(&cur.Id, &cur.Code, &cur.Name)
		if err != nil {
			panic("getCurrencies next: " + err.Error())
		}
		curs = append(curs, cur)
	}
	if rows.Err() != nil {
		panic("getCurrencies exit: " + err.Error())
	}

	// Return list
	return curs
}

// Get one currency by id
func getCurrency(id int) *Currency {

	// Connect to database
	db := dbConnect()
	defer db.Close()

	// Find currency, return nil if not found
	cur := Currency{}
	q := "select id, code, name from currency where id = $1"
	err := db.QueryRow(q, id).Scan(&cur.Id, &cur.Code, &cur.Name)
	if err != nil {
		return nil
	}

	return &cur
}

// Update an existing currency, or add new
func addUpdateCurrency(cur *Currency) {

	// Connect to database
	db := dbConnect()
	defer db.Close()

	// Attempt insert or update
	var err error
	if cur.Id == 0 {
		q := "insert into currency(code, name) values ($1, $2)"
		_, err = db.Exec(q, cur.Code, cur.Name)
	} else {
		q := "update currency set code = $1, name = $2 where id = $3"
		_, err = db.Exec(q, cur.Code, cur.Name, cur.Id)
	}

	// Check for error
	if err != nil {
		panic("addUpdateCurrency: " + err.Error())
	}
}

// Delete a currency by ID
// TODO: also delete all child records
func deleteCurrency(cid int) {

	db := dbConnect()
	defer db.Close()

	_, err := db.Exec("delete from currency where id = $1", cid)
	if err != nil {
		panic("deleteCurrency: " + err.Error())
	}
}

//----------------------------------------------------------------//
//                        CURRENCY RATES                          //
//----------------------------------------------------------------//

// Price of a currency relative to the home currency, i.e., multiply
// price by rate the get value in EUR

// Record format for a currency rate
type Rate struct {
	Id       int       // the id of this rate record
	Date     time.Time // the date for this rate
	Currency int       // id of the currency this price is for
	Rate     float64   // rate on this date
}

// Get rate by rate ID
func getRate(rid int) *Rate {

	// Connect to database
	db := dbConnect()
	defer db.Close()

	// Find rate, return nil if not found
	r := Rate{}
	q := "select id, currency_id, rdate, rate from currency_rate where id = $1"
	err := db.QueryRow(q, rid).Scan(&r.Id, &r.Currency, &r.Date, &r.Rate)
	if err != nil {
		return nil
	}

	return &r
}

// Get all rates for a currency
func getRates(cid int) []Rate {

	// Connect to database
	db := dbConnect()
	defer db.Close()

	// Execute query to get all rates, in date order
	rows, err := db.Query("select id, rdate, rate from currency_rate where currency_id = $1 order by rdate desc", cid)
	if err != nil {
		panic("getRates query: " + err.Error())
	}
	defer rows.Close()

	// Collect into a list
	rr := []Rate{}
	var ds string // buffer for reading date
	for rows.Next() {
		r := Rate{}
		err := rows.Scan(&r.Id, &ds, &r.Rate)
		if err != nil {
			panic("getRates next: " + err.Error())
		}
		r.Date = parseDate(ds)
		rr = append(rr, r)
	}
	if rows.Err() != nil {
		panic("getRatess exit: " + err.Error())
	}

	// Return list
	return rr
}

// Update an existing rate, or add new
func addUpdateRate(r *Rate) {

	// Connect to database
	db := dbConnect()
	defer db.Close()

	// Attempt insert or update
	var err error
	if r.Id == 0 {
		q := "insert into currency_rate(currency_id, rdate, rate) values ($1, $2, $3)"
		_, err = db.Exec(q, r.Currency, r.Date, r.Rate)
	} else {
		q := "update currency_rate set rdate = $1, rate = $2 where id = $3"
		_, err = db.Exec(q, r.Date, r.Rate, r.Id)
	}

	// Check for error
	if err != nil {
		panic("addUpdateRate: " + err.Error())
	}
}
