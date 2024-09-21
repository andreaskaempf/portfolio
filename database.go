// database.go
//
// Data model for the portfolio system, including structure definitions for
// all tables, and functions to retrieve or update data in the database.
// All database functions should be in this file.

package main

import (
	"database/sql"
	"fmt"
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
	Price    float64   // price on this date, in our local currency
	PriceX   float64   // price on this date, in the stock's currency
	Comments string    // any comments
}

// Get price by price ID
func getPrice(pid int) *Price {

	// Connect to database
	db := dbConnect()
	defer db.Close()

	// Find price, return nil if not found
	p := Price{}
	q := "select id, stock_id, pdate, price, pricex, comments from price where id = $1"
	err := db.QueryRow(q, pid).Scan(&p.Id, &p.Stock, &p.Date, &p.Price, &p.PriceX, &p.Comments)
	if err != nil {
		return nil
	}

	return &p
}

// Get all prices for a stock, sorted by ascending date
func getPrices(sid int) []Price {

	// Connect to database
	db := dbConnect()
	defer db.Close()

	// Execute query to get all prices for this stock, in date order
	rows, err := db.Query("select id, pdate, price, pricex, comments from price where stock_id = $1 order by pdate", sid)
	if err != nil {
		panic("getPrices query: " + err.Error())
	}
	defer rows.Close()

	// Collect into a list
	pp := []Price{}
	var ds string // buffer for reading date
	for rows.Next() {
		p := Price{}
		err := rows.Scan(&p.Id, &ds, &p.Price, &p.PriceX, &p.Comments)
		if err != nil {
			panic("getPrices next: " + err.Error())
		}
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
		q := "insert into price(stock_id, pdate, price, pricex, comments) values ($1, $2, $3, $4, $5)"
		_, err = db.Exec(q, p.Stock, p.Date, p.Price, p.PriceX, p.Comments)
	} else {
		q := "update price set pdate = $1, price = $2, pricex = $3, comments = $4 where id = $5"
		_, err = db.Exec(q, p.Date, p.Price, p.PriceX, p.Comments, p.Id)
	}

	// Check for error
	if err != nil {
		panic("addUpdatePrice: " + err.Error())
	}
}

//----------------------------------------------------------------//
//                      BUY/SELL TRANSACTIONS                     //
//----------------------------------------------------------------//

// Record format for one transaction
// Amount and fees are in local currency
// TODO: How do we back out the currency value vs. price?
// TODO: Add comments
type Transaction struct {
	Id     int       // ID of the transaction
	Stock  int       // ID of the stock
	Date   time.Time // the date for this transaction
	Q      float64   // the number of shares
	Amount float64   // the total amount paid, including fees
	Fees   float64   // commission or fees paid
}

// Get a list of all transactions, for a stock if argument is nonzero
func getTransactions(sid int) []Transaction {

	// Connect to database
	db := dbConnect()
	defer db.Close()

	// Execute query to get all transactions
	var err error
	var rows *sql.Rows
	if sid > 0 {
		rows, err = db.Query("select id, stock_id, tdate, q, amount, fees from trans where stock_id == $1 order by tdate", sid)
	} else {
		rows, err = db.Query("select id, stock_id, tdate, q, amount, fees from trans order by tdate")
	}
	if err != nil {
		panic("getTransactions query: " + err.Error())
	}
	defer rows.Close()

	// Collect into a list
	tt := []Transaction{}
	for rows.Next() {
		t := Transaction{}
		var ds string
		err := rows.Scan(&t.Id, &t.Stock, &ds, &t.Q, &t.Amount, &t.Fees)
		if err != nil {
			panic("getTransactions next: " + err.Error())
		}
		t.Date = parseDate(ds)
		tt = append(tt, t)
	}
	if rows.Err() != nil {
		panic("getTransactions exit: " + err.Error())
	}

	// Return list
	return tt
}

// Get one transaction by id
func getTransaction(tid int) *Transaction {

	// Connect to database
	db := dbConnect()
	defer db.Close()

	// Find and read transaction, return nil if not found
	t := Transaction{}
	var ds string
	q := "select id, stock_id, tdate, q, amount, fees from trans where id = $1"
	err := db.QueryRow(q, tid).Scan(&t.Id, &t.Stock, &ds, &t.Q, &t.Amount, &t.Fees)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	t.Date = parseDate(ds)

	return &t
}

// Update an existing transaction, or add new
func addUpdateTransaction(t *Transaction) {

	// Connect to database
	db := dbConnect()
	defer db.Close()

	// Attempt insert or update
	var err error
	if t.Id == 0 {
		q := "insert into trans(stock_id, tdate, q, amount, fees) values ($1, $2, $3, $4, $5)"
		_, err = db.Exec(q, t.Stock, formatDate(t.Date), t.Q, t.Amount, t.Fees)
	} else {
		q := "update trans set tdate = $1, q = $2, amount = $3, fees = $4 where id = $5"
		_, err = db.Exec(q, formatDate(t.Date), t.Q, t.Amount, t.Fees, t.Id)
	}

	// Check for error
	if err != nil {
		panic("addUpdateTransaction: " + err.Error())
	}
}

// Delete a transaction by ID
// TODO: also delete all child records
func deleteTransaction(tid int) {

	db := dbConnect()
	defer db.Close()

	_, err := db.Exec("delete from trans where id = $1", tid)
	if err != nil {
		panic("deleteStock: " + err.Error())
	}
}

//----------------------------------------------------------------//
//                            DIVIDENDS                           //
//----------------------------------------------------------------//

// Record format for one dividend
type Dividend struct {
	Id       int       // ID of the transaction
	Stock    int       // ID of the stock
	Date     time.Time // the date for this transaction
	Amount   float64   // the total amount paid, including fees
	Comments string
}

// Get a list of all dividends for a stock, or for all stocks if ID is 0
func getDividends(sid int) []Dividend {

	// Connect to database
	db := dbConnect()
	defer db.Close()

	// Execute query to get all transactions
	var err error
	var rows *sql.Rows
	q := "select id, stock_id, tdate, amount, comments from dividend"
	if sid > 0 {
		q += " where stock_id == $1 order by tdate"
		rows, err = db.Query(q, sid)
	} else {
		q += " order by tdate"
		rows, err = db.Query(q)
	}
	if err != nil {
		panic("getDividends query: " + err.Error())
	}
	defer rows.Close()

	// Collect into a list
	dd := []Dividend{}
	for rows.Next() {
		d := Dividend{}
		var ds string
		err := rows.Scan(&d.Id, &d.Stock, &ds, &d.Amount, &d.Comments)
		if err != nil {
			panic("getDividends next: " + err.Error())
		}
		d.Date = parseDate(ds)
		dd = append(dd, d)
	}
	if rows.Err() != nil {
		panic("getDividends exit: " + err.Error())
	}

	// Return list
	return dd
}

// Get one dividend by id
func getDividend(did int) *Dividend {

	// Connect to database
	db := dbConnect()
	defer db.Close()

	// Find and read transaction, return nil if not found
	d := Dividend{}
	var ds string
	q := "select id, stock_id, tdate, amount, comments from dividend where id = $1"
	err := db.QueryRow(q, did).Scan(&d.Id, &d.Stock, &ds, &d.Amount, &d.Comments)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	d.Date = parseDate(ds)

	return &d
}

// Update an existing dividend, or add new
func addUpdateDividend(d *Dividend) {

	// Connect to database
	db := dbConnect()
	defer db.Close()

	// Attempt insert or update
	var err error
	if d.Id == 0 {
		q := "insert into dividend(stock_id, tdate, amount, comments) values ($1, $2, $3, $4)"
		_, err = db.Exec(q, d.Stock, formatDate(d.Date), d.Amount, d.Comments)
	} else {
		q := "update dividend set tdate = $1, amount = $2, comments = $3 where id = $4"
		_, err = db.Exec(q, formatDate(d.Date), d.Amount, d.Comments, d.Id)
	}

	// Check for error
	if err != nil {
		panic("addUpdateTransaction: " + err.Error())
	}
}

// Delete a dividend by ID
func deleteDividend(did int) {

	db := dbConnect()
	defer db.Close()

	_, err := db.Exec("delete from dividend where id = $1", did)
	if err != nil {
		panic("deleteDividend: " + err.Error())
	}
}

//----------------------------------------------------------------//
//                        CASH TRANSACTIONS                       //
//----------------------------------------------------------------//

// Record format for one cash transaction
type Cash struct {
	Id       int       // ID of the transaction
	Date     time.Time // the date for this transaction
	Type     string    // "deposit", "withdraw"
	Amount   float64   // + for deposit, - for withdraw
	Comments string
}

// Get a list of all cash transactions
func getCashTransactions() []Cash {

	// Connect to database
	db := dbConnect()
	defer db.Close()

	// Execute query to get all transactions
	var err error
	var rows *sql.Rows
	rows, err = db.Query("select id, tdate, ttype, amount, comments from cash order by tdate")
	if err != nil {
		panic("getCashTransactions query: " + err.Error())
	}
	defer rows.Close()

	// Collect into a list
	cc := []Cash{}
	for rows.Next() {
		c := Cash{}
		var ds string
		err := rows.Scan(&c.Id, &ds, &c.Type, &c.Amount, &c.Comments)
		if err != nil {
			panic("getCashTransactions next: " + err.Error())
		}
		c.Date = parseDate(ds)
		cc = append(cc, c)
	}
	if rows.Err() != nil {
		panic("getCashTransactions exit: " + err.Error())
	}

	// Return list
	return cc
}

// Get one cash transaction by id
func getCashTransaction(tid int) *Cash {

	// Connect to database
	db := dbConnect()
	defer db.Close()

	// Find and read transaction, return nil if not found
	c := Cash{}
	var ds string
	q := "select id, tdate, ttype, amount, comments from cash where id = $1"
	err := db.QueryRow(q, tid).Scan(&c.Id, &ds, &c.Type, &c.Amount, &c.Comments)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	c.Date = parseDate(ds)

	return &c
}

// Update an existing transaction, or add new
func addUpdateCash(t *Cash) {

	// Connect to database
	db := dbConnect()
	defer db.Close()

	// Attempt insert or update
	var err error
	if t.Id == 0 {
		q := "insert into cash(tdate, ttype, amount, comments) values ($1, $2, $3, $4)"
		_, err = db.Exec(q, formatDate(t.Date), t.Type, t.Amount, t.Comments)
	} else {
		q := "update cash set tdate = $1, ttype = $2, amount = $3, comments = $4 where id = $5"
		_, err = db.Exec(q, formatDate(t.Date), t.Type, t.Amount, t.Comments, t.Id)
	}

	// Check for error
	if err != nil {
		panic("addUpdateCash: " + err.Error())
	}
}

// Delete a cash transaction by ID
func deleteCash(tid int) {

	db := dbConnect()
	defer db.Close()

	_, err := db.Exec("delete from cash where id = $1", tid)
	if err != nil {
		panic("deleteCash: " + err.Error())
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

// Get one currency by code
func getCurrencyCode(code string) *Currency {

	// Connect to database
	db := dbConnect()
	defer db.Close()

	// Find currency, return nil if not found
	cur := Currency{}
	q := "select id, code, name from currency where code = $1"
	err := db.QueryRow(q, code).Scan(&cur.Id, &cur.Code, &cur.Name)
	if err != nil {
		fmt.Println("getCurrencyCode: " + err.Error())
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
