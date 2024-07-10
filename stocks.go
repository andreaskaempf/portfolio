package main

import (
	"fmt"
	"net/http"

	//"strings"

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

// Show one stock
func showStock(c *gin.Context) {

	// Parse the ID and get the event
	sid := parseInt(c.Param("id"))
	s := getStock(sid)
	if sid < 1 || s == nil {
		c.String(http.StatusOK, fmt.Sprintf("Stock %d not found", sid))
		return
	}

	// Show page
	c.HTML(http.StatusOK, "stock.html",
		gin.H{"s": s, "menu": menu})
}
