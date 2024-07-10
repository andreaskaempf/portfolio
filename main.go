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

func main() {

	// Create router, initialize templates and location of static files
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	r.Static("/static", "./static")

	// Define routes
	r.GET("/", showStocks)
	r.GET("/Home", showStocks)
	r.GET("/stock/:id", showStock)

	// Start server
	r.Run() // ":8222")
}
