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
	//"time"

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

// Stocks: code, type, description, currency

// Record format for one stock
type Stock struct {
	Id       int
	Code     string
	Name     string
	Currency string
}

// Get a list of all events, most recent at the top
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

/*
// Update an existing event, or add new
func addUpdateEvent(e *Event) {

	// Connect to database
	db := dbConnect()
	defer db.Close()

	// Attempt insert or update
	var err error
	if e.Id == 0 {
		q := "insert into event(sdate, edate, name, location, description, active) values ($1, $2, $3, $4, $5, $6)"
		_, err = db.Exec(q, e.SDate, e.EDate, e.Name, e.Location, e.Description, e.Active)
	} else {
		q := "update event set sdate = $1, edate = $2, name = $3, location = $4, description = $5, active = $6 where id = $7"
		_, err = db.Exec(q, e.SDate, e.EDate, e.Name, e.Location, e.Description, e.Active, e.Id)
	}

	// Check for error
	if err != nil {
		panic("addUpdateEvent: " + err.Error())
	}
}

// Delete an event by ID
// TODO: also delete all sessions, assignments, etc.
func deleteEvent(eid int) {

	db := dbConnect()
	defer db.Close()

	//_, err1 := db.Exec("delete from assignment where person_id = $1", pid)
	_, err2 := db.Exec("delete from event where id = $1", eid)
	if err2 != nil { //}|| err2 != nil {
		panic("deleteEvent: " + err2.Error()) // + err2.Error())
	}
}
*/
