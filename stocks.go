package main

import (
	"fmt"
	"net/http"

	"strings"

	"github.com/gin-gonic/gin"
)

// Show list of all stocks
func showStocks(c *gin.Context) {

	// Get a list of holdings
	stocks := getStocks()

	// Show page
	c.HTML(http.StatusOK, "stocks.html",
		gin.H{"stocks": stocks, "menu": menu})
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

	// Get all prices for this stock
	prices := getPrices(sid)

	// Show page
	c.HTML(http.StatusOK, "stock.html",
		gin.H{"s": s, "prices": prices, "menu": menu})
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
		gin.H{"s": s, "currencies": currencies, "menu": menu})
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
		c.Redirect(http.StatusFound, "/stocks")
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
			gin.H{"s": s, "menu": menu})
	} else if confirm == "yes" { // confirmed, delete stock
		deleteStock(sid)
		c.Redirect(http.StatusFound, "/stocks")
	} else { // confirmation denied, back to stock page
		c.Redirect(http.StatusFound, fmt.Sprintf("/stock/%d", sid))
	}
}

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
	var t := &Transaction{}
	if tid == 0 {
		sid_, _ := c.GetQuery("sid")
		sid = parseInt(sid_)
		if sid < 0 {
			c.String(http.StatusNotFound, "Missing stock ID, required for adding transaction")
			return
		}
		t = Transaction{Stock: sid}
	} else {
		t = getTransaction(tid)
		if t == nil {
			c.String(http.StatusNotFound, "Transaction not found")
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
		gin.H{"t": t, "s": s, "menu": menu})
}

// Process form to update or add an stock
func saveTransaction(c *gin.Context) {

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
		c.Redirect(http.StatusFound, "/stocks")
	} else {
		c.Redirect(http.StatusFound, fmt.Sprintf("/stock/%d", sid))
	}
}