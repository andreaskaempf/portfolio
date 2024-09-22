package main

import (
	"fmt"
	"text/template"
	"time"

	"github.com/gin-gonic/gin"
)

// Default menu
var menu = []string{"Portfolio", "Stocks", "Cash", "Currencies"}

// List of currency codes (TODO: in database)
var currencies = []string{"EUR", "USD", "CHF", "GBP", "NZD", "AUD"}

// List of cash transaction types
var cashTypes = []string{"Deposit", "Withdrawal"}

// Home currency (TODO: in database)
var homeCurrency = currencies[0]

// Last date entered on a transaction this session
var lastTransDate time.Time

func main() {

	// Set the last time entered to now
	lastTransDate = time.Now()

	// Create router, define custom functions
	r := gin.Default()
	r.FuncMap = template.FuncMap{
		"add":       func(a, b float64) float64 { return a + b },
		"sub":       func(a, b float64) float64 { return a - b },
		"mul":       func(a, b float64) float64 { return a * b },
		"div":       func(a, b float64) float64 { return a / b },
		"fmtDate":   formatDate,
		"fmtAmount": formatFloat,
	}

	// Initialize templates and location of static files
	r.LoadHTMLGlob("templates/*")
	r.Static("/static", "./static")

	// Route for home page with portfolio
	r.GET("/", showPortfolio)
	r.GET("/Portfolio", showPortfolio)

	// Routes for stocks
	r.GET("/Home", showStocks)
	r.GET("/Stocks", showStocks)
	r.GET("/stock/:id", showStock)
	r.GET("/edit_stock/:id", editStock)
	r.POST("/update_stock", saveStock)
	r.GET("/delete_stock/:id", delStock)

	// Routes for stock prices
	r.GET("/edit_price/:pid", editPrice)
	r.POST("/update_price", updatePrice)
	r.GET("/get_prices/:sid", getPricesJSON)

	// Routes for buy/sell transactions
	r.GET("/edit_transaction/:tid", editTransaction)
	r.POST("/update_transaction", saveTransaction)

	// Routes for dividends
	r.GET("/edit_dividend/:did", editDividend)
	r.POST("/update_dividend", saveDividend)

	// Cash pages
	r.GET("/Cash", showCashPage)
	r.GET("/cash/:id", showCash)
	r.GET("/edit_cash/:id", editCash)
	r.POST("/update_cash", saveCash)
	r.GET("/delete_cash/:id", delCash)

	// Routes for currencies and rates
	r.GET("/Currencies", showCurrencies)
	r.GET("/currency/:id", showCurrency)
	r.GET("/edit_currency/:id", editCurrency)
	r.POST("/update_currency", saveCurrency)
	r.GET("/delete_currency/:id", delCurrency)
	r.GET("/edit_rate/:rid", editRate)
	r.POST("/update_rate", updateRate)

	// Start server
	fmt.Println("Running on http://localhost:8080")
	r.Run() // for different port: ":8222")
}
