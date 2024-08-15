package main

import (
	//"fmt"
	//"io"
	//"os"
	//"time"
	//"net/http"
	"text/template"

	"github.com/gin-gonic/gin"
)

// Default menu
var menu = []string{"Portfolio", "Stocks", "Currencies"}

// List of currency codes (TODO: in database)
var currencies = []string{"EUR", "USD", "GBP", "NZD", "AUD"}

// Home currency (TODO: in database)
var homeCurrency = currencies[0]

func main() {

	// Create router, define custom functions
	r := gin.Default()
	r.FuncMap = template.FuncMap{"mul": func(a, b float64) float64 {
		return a * b
	}, "add": func(a, b float64) float64 {
		return a + b
	}}

	// Initialize templates and location of static files
	r.LoadHTMLGlob("templates/*")
	r.Static("/static", "./static")

	// Route for home page with portfolio
	r.GET("/", showPortfolio)
	r.GET("/Portfolio", showPortfolio)

	// Routes for stocks and prices
	r.GET("/Home", showStocks)
	r.GET("/Stocks", showStocks)
	r.GET("/stock/:id", showStock)
	r.GET("/edit_stock/:id", editStock)
	r.POST("/update_stock", saveStock)
	r.GET("/delete_stock/:id", delStock)
	r.GET("/edit_price/:pid", editPrice)
	r.POST("/update_price", updatePrice)

	// Routes for transactions
	r.GET("/edit_transaction/:tid", editTransaction)
	r.POST("/update_transaction", saveTransaction)

	// Routes for currencies and rates
	r.GET("/Currencies", showCurrencies)
	r.GET("/currency/:id", showCurrency)
	r.GET("/edit_currency/:id", editCurrency)
	r.POST("/update_currency", saveCurrency)
	r.GET("/delete_currency/:id", delCurrency)
	r.GET("/edit_rate/:rid", editRate)
	r.POST("/update_rate", updateRate)

	// Start server
	r.Run() // ":8222")
}
