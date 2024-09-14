// Cash: deposit, withdraw, dividends, buy/sell

package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Show page with cash balance and transaction history
func showCashPage(c *gin.Context) {

	// Get cash transactions up to today
	today := time.Now()
	trans := getAllCash(today) // including "virtual" buy/sell

	// TODO: get cash value today, just sum of table above
	bal := 100.0

	// Show page
	c.HTML(http.StatusOK, "cash.html",
		gin.H{"d": today, "transactions": trans, "balance": bal, "menu": menu})
}

// Page to show one cash transaction
func showCash(c *gin.Context) {

	// Parse the ID and get the cash transaction
	tid := parseInt(c.Param("id"))
	t := getCashTransaction(tid)
	if t == nil {
		c.String(http.StatusOK, fmt.Sprintf("Cash %d not found", tid))
		return
	}

	// Show page
	c.HTML(http.StatusOK, "cash_trans.html",
		gin.H{"c": t, "menu": menu})
}

// Show form to edit/create a cash transaction
func editCash(c *gin.Context) {

	// Get cash ID (will be 0 to add)
	tid := parseInt(c.Param("id"))
	if tid < 0 {
		c.String(http.StatusOK, "Invalid cash ID")
		return
	}

	// Get the cash transaction or create "blank" cash
	t := &Cash{}
	if tid > 0 {
		t = getCashTransaction(tid)
		if t == nil {
			c.String(http.StatusOK, "Cash not found")
			return
		}
	} else {
		t.Date = time.Now()
		t.Type = cashTypes[0]
	}

	// Adjust withdrawal amounts to be positive
	if t.Type == "Withdrawal" {
		t.Amount *= -1.0
	}

	// Show the form to edit cash
	c.HTML(http.StatusOK, "edit_cash.html",
		gin.H{"c": t, "types": cashTypes, "menu": menu})
}

// Process form to update or add a cash transaction
func saveCash(c *gin.Context) {

	// Get cash ID (will be 0 to add a cash)
	tid_, ok := c.GetPostForm("id")
	if !ok {
		c.String(http.StatusOK, "saveCash: Missing cash ID")
		return
	}
	tid := parseInt(tid_)
	if tid < 0 {
		c.String(http.StatusOK, "saveCash: Invalid cash ID")
		return
	}

	// Get the cash or create "blank" cash
	t := &Cash{}
	if tid > 0 {
		t = getCashTransaction(tid)
		if t == nil {
			c.String(http.StatusOK, "saveCash: cash not found")
			return
		}
	}

	// Update the cash with the form inputs
	ds, _ := c.GetPostForm("date")
	t.Date = parseDate(ds)
	t.Type, _ = c.GetPostForm("type")
	amt, _ := c.GetPostForm("amount")
	t.Amount = parseFloat(amt)
	t.Comments, _ = c.GetPostForm("comments")

	// Some validation
	if !validDate(t.Date) || t.Amount == 0 {
		c.String(http.StatusOK, "Invalid date, or amount is zero")
		return
	}

	// If a withdrawal, make amount negative
	if t.Type == "Withdrawal" {
		t.Amount *= -1.0
	}

	// Create or update transaction in database
	addUpdateCash(t)

	// Go back to cash page or list
	if tid == 0 {
		c.Redirect(http.StatusFound, "/Cash")
	} else {
		c.Redirect(http.StatusFound, fmt.Sprintf("/cash/%d", tid))
	}
}

// Delete cash: ask for confirmation first
func delCash(c *gin.Context) {

	// Get the cash (URL positional param)
	tid := parseInt(c.Param("id"))
	t := getCashTransaction(tid)
	if tid <= 0 || t == nil {
		c.String(http.StatusOK, "Cash not found")
		return
	}

	// Ask for confirmation, or go ahead and delete if confirmed
	confirm, _ := c.GetQuery("confirm")
	if confirm == "" { // no confirmation, show form
		c.HTML(http.StatusOK, "del_cash.html", gin.H{"t": t, "menu": menu})
	} else if confirm == "yes" { // confirmed, delete cash
		deleteCash(tid)
		c.Redirect(http.StatusFound, "/Cash")
	} else { // confirmation denied, back to cash page
		c.Redirect(http.StatusFound, fmt.Sprintf("/cash/%d", tid))
	}
}

// Get cash transactions up to a particular date, including
// "virtual" buy/sell and dividends
func getAllCash(d time.Time) []Cash {

	// Get all explicit transactions, e.g., deposits & withdrawals
	// TODO: filter by date
	cc := getCashTransactions()

	// TODO: add all buy/sell transactions

	// TODO: add dividends

	// TODO: sort by date
	return cc
}
