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

	// Define routes
	r.GET("/", showStocks)
	r.GET("/stocks", showStocks)
	r.GET("/stock/:id", showStock)
	r.GET("/edit_stock/:id", editStock)
	r.POST("/update_stock", saveStock)
	r.GET("/delete_stock/:id", delStock)

	// Start server
	r.Run() // ":8222")
}
