package main

import (
	"fmt"
	//"io"
	//"os"
	//"time"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	fmt.Println("vim-go")

	// Create router, initialize templates and location of static files
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	r.Static("/static", "./static")

	// Define routes
	r.GET("/", homePage)

	// Start server
	r.Run() // ":8222")
}

// Show home page
func homePage(c *gin.Context) {

	// Get a list of holdings
	stocks := getStocks()

	// Default menu
	menu := []string{"Home", "Currencies", "Admin"}

	// Show page
	c.HTML(http.StatusOK, "stocks.html",
		gin.H{"stocks": stocks, "menu": menu})
}

// Get a list of stocks, hard coded for now
func getStocks() []Stock {
	return []Stock{
		Stock{"IBM", "IBM International", 200, 32, 0},
		Stock{"MSFT", "Microsoft", 100, 22, 0},
		Stock{"NVID", "Nvidia, Corp", 150, 270, 0},
	}
}

type Stock struct {
	Code, Description string
	Q, Price, Value   float64
}
