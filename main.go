package main

import (
	//"fmt"
	//"io"
	//"os"
	//"time"
	//"net/http"

	"github.com/gin-gonic/gin"
)

// Default menu
var menu = []string{"Home", "Currencies", "Admin"}

// List of currency codes (TODO: in database)
var currencies = []string{"EUR", "USD", "GBP", "NZD", "AUD"}

func main() {

	// Create router, initialize templates and location of static files
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	r.Static("/static", "./static")

	// Routes for stocks and prices
	r.GET("/", showStocks)
	r.GET("/Home", showStocks)
	r.GET("/stocks", showStocks)
	r.GET("/stock/:id", showStock)
	r.GET("/edit_stock/:id", editStock)
	r.POST("/update_stock", saveStock)
	r.GET("/delete_stock/:id", delStock)
	r.GET("/edit_price/:pid", editPrice)
	r.POST("/update_price", updatePrice)

	// Routes for currencies and rates
	r.GET("/Currencies", showCurrencies)
	r.GET("/currency/:id", showCurrency)
	r.GET("/edit_currency/:id", editCurrency)
	r.POST("/update_currency", saveCurrency)
	r.GET("/delete_currency/:id", delCurrency)
	//r.GET("/edit_rate/:rid", editRate)
	//r.POST("/update_rate", updateRate)

	// Start server
	r.Run() // ":8222")
}
