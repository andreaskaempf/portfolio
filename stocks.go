// Stocks, including buy/sell transactions and dividends

package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

//-----------------------------------------------------------------//
//                             STOCKS                              //
//-----------------------------------------------------------------//

// Show list of all stocks
func showStocks(c *gin.Context) {

	// Get a list of all stocks, including not held
	//stocks := getStocks()
	today := time.Now()
	holdings := getPortfolio(today, false)

	// Show page
	c.HTML(http.StatusOK, "stocks.html",
		gin.H{"holdings": holdings, "menu": menu, "current": "Stocks"})
}

// Page to show one stock
func showStock(c *gin.Context) {

	// Parse the ID and get the stock
	sid := parseInt(c.Param("id"))
	s := getStock(sid)
	if s == nil {
		c.String(http.StatusOK, fmt.Sprintf("Stock %d not found", sid))
		return
	}

	// Get all transactions, dividends and prices for this stock
	prices := getPrices(sid)
	transactions := getTransactions(sid)
	dividends := getDividends(sid)

	// Count up the number of units held
	var units float64
	for _, t := range transactions {
		units += t.Q
	}

	// Show page
	c.HTML(http.StatusOK, "stock.html",
		gin.H{"s": s, "transactions": transactions, "units": units,
			"prices": prices, "dividends": dividends, "home": homeCurrency,
			"menu": menu, "current": "Stocks"})
}

// Show form to edit a stock (including a new one)
func editStock(c *gin.Context) {

	// Get stock ID (will be 0 to add an stock)
	sid := parseInt(c.Param("id"))
	if sid < 0 {
		c.String(http.StatusOK, "Invalid stock ID")
		return
	}

	// Get the stock or create "blank" stock
	s := &Stock{}
	if sid > 0 {
		s = getStock(sid)
		if s == nil {
			c.String(http.StatusOK, "Stock not found")
			return
		}
	}

	// Show the form to edit stock
	c.HTML(http.StatusOK, "edit_stock.html",
		gin.H{"s": s, "currencies": currencies,
			"menu": menu, "current": "Stocks"})
}

// Process form to update or add an stock
func saveStock(c *gin.Context) {

	// Get stock ID (will be 0 to add a stock)
	sid_, ok := c.GetPostForm("id")
	if !ok {
		c.String(http.StatusOK, "saveStock: Missing stock ID")
		return
	}
	sid := parseInt(sid_)
	if sid < 0 {
		c.String(http.StatusOK, "saveStock: Invalid stock ID")
		return
	}

	// Get the stock or create "blank" stock
	s := &Stock{}
	if sid > 0 {
		s = getStock(sid)
		if s == nil {
			c.String(http.StatusOK, "saveStock: stock not found")
			return
		}
	}

	// Update the stock with the form inputs
	s.Code, _ = c.GetPostForm("code")
	s.Name, _ = c.GetPostForm("name")
	s.Currency, _ = c.GetPostForm("currency")

	// Some validation
	s.Code = strings.TrimSpace(s.Code)
	s.Name = strings.TrimSpace(s.Name)
	s.Currency = strings.TrimSpace(s.Currency)
	if len(s.Code) == 0 || len(s.Name) == 0 {
		c.String(http.StatusOK, "Invalid inputs: cannot be blank")
		return
	}

	// Create or update person database
	addUpdateStock(s)

	// Go back to stocks page or list
	if sid == 0 {
		c.Redirect(http.StatusFound, "/Stocks")
	} else {
		c.Redirect(http.StatusFound, fmt.Sprintf("/stock/%d", sid))
	}
}

// Delete stock: ask for confirmation first
func delStock(c *gin.Context) {

	// Get the stock (URL positional param)
	sid := parseInt(c.Param("id"))
	s := getStock(sid)
	if sid <= 0 || s == nil {
		c.String(http.StatusOK, "Stock not found")
		return
	}

	// Ask for confirmation, or go ahead and delete if confirmed
	confirm, _ := c.GetQuery("confirm")
	if confirm == "" { // no confirmation, show form
		c.HTML(http.StatusOK, "del_stock.html",
			gin.H{"s": s, "menu": menu, "current": "Stocks"})
	} else if confirm == "yes" { // confirmed, delete stock
		deleteStock(sid)
		c.Redirect(http.StatusFound, "/Stocks")
	} else { // confirmation denied, back to stock page
		c.Redirect(http.StatusFound, fmt.Sprintf("/stock/%d", sid))
	}
}

//-----------------------------------------------------------------//
//                      BUY/SELL TRANSACTIONS                      //
//-----------------------------------------------------------------//

// Edit/create a transaction for a stock. Edits the transaction if an ID
// is provided. If zero, adds a transaction for the stock ID expected in
// the query string.
// Show form to edit a stock (including a new one)
func editTransaction(c *gin.Context) {

	// Get transaction ID (will be 0 to add)
	tid := parseInt(c.Param("tid"))
	if tid < 0 {
		c.String(http.StatusNotFound, "Invalid transaction ID")
		return
	}

	// If transaction ID is zero, create a "blank" transaction, otherwise
	// get the transaction. Blank transaction requires a stock to be provided
	// as query string.
	var sid int
	t := &Transaction{}
	if tid == 0 {
		sid_, _ := c.GetQuery("sid")
		sid = parseInt(sid_)
		if sid <= 0 {
			c.String(http.StatusNotFound, "Missing stock ID, required for adding transaction")
			return
		}
		t = &Transaction{Stock: sid, Date: lastTransDate}
	} else {
		t = getTransaction(tid)
		if t == nil {
			c.String(http.StatusNotFound, fmt.Sprintf("Transaction %d not found", tid))
			return
		}
		sid = t.Stock
	}

	// Get the stock as well
	s := getStock(sid)
	if s == nil {
		c.String(http.StatusNotFound, "Stock not found")
		return
	}

	// Show the form to edit transaction
	c.HTML(http.StatusOK, "edit_transaction.html",
		gin.H{"t": t, "s": s, "menu": menu, "current": "Stocks"})
}

// Process form to update or add a transaction
func saveTransaction(c *gin.Context) {

	// Get transaction and stock ID (tid will be 0 to add)
	tid_, ok1 := c.GetPostForm("tid")
	sid_, ok2 := c.GetPostForm("sid")
	tid := parseInt(tid_)
	sid := parseInt(sid_)
	if !ok1 || !ok2 || sid < 0 || tid < 0 {
		c.String(http.StatusOK, "saveStock: Missing or invalid stock and transaction IDs")
		return
	}

	// Get the transaction or create a "blank" one
	t := &Transaction{Id: tid, Stock: sid}
	if tid > 0 {
		t = getTransaction(tid)
		if t == nil {
			c.String(http.StatusNotFound, "saveTransaction: not found")
			return
		}
	}

	// Update the stock with the form inputs
	ds, _ := c.GetPostForm("date")
	q, _ := c.GetPostForm("q")
	amount, _ := c.GetPostForm("amount")
	fees, _ := c.GetPostForm("fees")

	// Convert and validate fields
	t.Date = parseDate(ds)
	t.Q = parseFloat(q)
	t.Amount = parseFloat(amount)
	t.Fees = parseFloat(fees)
	if t.Date.Year() < 2000 || t.Q == 0 || t.Amount <= 0 || t.Fees < 0 {
		c.String(http.StatusOK, "Invalid inputs")
		return
	}

	// Create or update person database
	addUpdateTransaction(t)

	// Remember the last transaction date for next entry
	lastTransDate = t.Date

	// Go back to stock page
	c.Redirect(http.StatusFound, fmt.Sprintf("/stock/%d", sid))
}

//-----------------------------------------------------------------//
//                           DIVIDENDS                             //
//-----------------------------------------------------------------//

// Edit/create a dividend for a stock. Edits the dividend if an ID
// is provided. If zero, adds a transaction for the stock ID expected in
// the query string.
func editDividend(c *gin.Context) {

	// Get dividend ID (will be 0 to add)
	did := parseInt(c.Param("did"))
	if did < 0 {
		c.String(http.StatusNotFound, "Invalid dividend ID")
		return
	}

	// If dividend ID is zero, create a "blank" dividend, otherwise
	// get the dividend. Blank dividend requires a stock to be provided
	// as query string.
	var sid int
	d := &Dividend{}
	if did == 0 {
		sid_, _ := c.GetQuery("sid")
		sid = parseInt(sid_)
		if sid <= 0 {
			c.String(http.StatusNotFound, "Missing stock ID, required for adding dividend")
			return
		}
		d = &Dividend{Stock: sid, Date: lastTransDate} // TODO: why not just reuse blank dividend?
		d.Comments = "From statement"
	} else {
		d = getDividend(did)
		if d == nil {
			c.String(http.StatusNotFound, fmt.Sprintf("Dividend %d not found", did))
			return
		}
		sid = d.Stock
	}

	// Get the stock as well
	s := getStock(sid)
	if s == nil {
		c.String(http.StatusNotFound, "Stock not found")
		return
	}

	// Show the form to edit dividend
	c.HTML(http.StatusOK, "edit_dividend.html",
		gin.H{"d": d, "s": s, "menu": menu, "current": "Stocks"})
}

// Process form to update or add a transaction
func saveDividend(c *gin.Context) {

	// Get dividend and stock ID (did will be 0 to add)
	did_, ok1 := c.GetPostForm("did")
	sid_, ok2 := c.GetPostForm("sid")
	did := parseInt(did_)
	sid := parseInt(sid_)
	if !ok1 || !ok2 || sid < 0 || did < 0 {
		c.String(http.StatusOK, "saveDividend: Missing or invalid stock and dividend IDs")
		return
	}

	// Get the dividend or create a "blank" one
	d := &Dividend{Id: did, Stock: sid}
	if did > 0 {
		d = getDividend(did)
		if d == nil {
			c.String(http.StatusNotFound, "saveDividend: not found")
			return
		}
	}

	// Update the dividend with the form inputs
	ds, _ := c.GetPostForm("date")
	amount, _ := c.GetPostForm("amount")
	d.Comments, _ = c.GetPostForm("comments")

	// Convert and validate fields
	d.Date = parseDate(ds)
	d.Amount = parseFloat(amount)
	if !validDate(d.Date) || d.Amount <= 0 {
		c.String(http.StatusOK, "Invalid inputs")
		return
	}

	// Create or update person database
	addUpdateDividend(d)

	// Remember the last transaction date for next entry
	lastTransDate = d.Date

	// Go back to stock page
	c.Redirect(http.StatusFound, fmt.Sprintf("/stock/%d", sid))
}
