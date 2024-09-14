// Cash: deposit, withdraw, dividends, buy/sell

package main

import (
	//"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Show table of holdings with current value and return since purchase
func showCash(c *gin.Context) {

	// Get cash transactions up to today
	today := time.Now()
	trans := getCashTransactions(today)

	// TODO: get cash value today
	bal := 100.0

	// Show page
	c.HTML(http.StatusOK, "cash.html",
		gin.H{"d": today, "transactions": trans, "balance": bal, "menu": menu})
}

// A cash transaction
type CashTransaction struct {
	Date      time.Time
	Amount    float64 // pos = deposit/sell, neg = withdraw/buy
	TransType string  // transaction type
	Comments  string  // comments about the transaction
}

// Get cash transactions up to a particular date
func getCashTransactions(d time.Time) []CashTransaction {

	// Get all explicit transactions, e.g., deposits, withdrawals, dividends
	trans := []CashTransaction{}
	//for _, s := range getCash() {}

	// TODO: add all buy/sell transactions

	// TODO: sort by date
	return trans
}
