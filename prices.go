package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Add a price for a stock
func editPrice(c *gin.Context) {

	// Get Stock ID from query string if provided (only required for adding new prices)
	// If there is no price ID, expect a stock ID as query parameter
	var sid int
	sid_, ok := c.GetQuery("sid")
	if ok {
		sid = parseInt(sid_)
		if sid < 0 {
			c.String(http.StatusNotFound, "Invalid stock ID")
			return
		}
	}

	// Get the Price ID. If this is zero, means to add a new price to
	// the current Stock ID (see below), otherwise get the price to edit.
	pid := parseInt(c.Param("pid"))
	var p Price
	if pid < 0 {
		c.String(http.StatusNotFound, "Invalid price ID")
		return
	} else if pid == 0 { // create a new price
		if sid <= 0 {
			c.String(http.StatusNotFound, "Cannot add price without stock ID")
			return
		}
		newPrice := 0.0 // TODO: default price sould be most recent one
		p = Price{Id: 0, Date: time.Now(), Stock: sid, Price: newPrice, Comments: ""}
	} else { // get existing price
		pp := getPrice(pid)
		if pp == nil {
			c.String(http.StatusNotFound, "Price not found")
			return
		}
		p = *pp
		sid = p.Stock
	}

	// Show the form to edit price
	c.HTML(http.StatusOK, "edit_price.html",
		gin.H{"p": p, "ds": formatDate(p.Date), "menu": menu, "current": "Stocks"})
}

// Create or update a price
func updatePrice(c *gin.Context) {

	// Get stock and price ID (latter will be 0 to add a price)
	sid_, ok := c.GetPostForm("sid")
	if !ok {
		c.String(http.StatusOK, "savePrice: Missing stock ID")
		return
	}
	pid_, ok := c.GetPostForm("pid")
	if !ok {
		c.String(http.StatusOK, "savePrice: Missing price ID")
		return
	}
	sid := parseInt(sid_)
	pid := parseInt(pid_)
	if sid < 0 || pid < 0 {
		c.String(http.StatusOK, "savePrice: Invalid stock or price ID")
		return
	}

	// Get the price or create "blank" one
	p := &Price{Stock: sid}
	if pid > 0 {
		p = getPrice(pid)
		if p == nil {
			c.String(http.StatusOK, "savePrice: price not found")
			return
		}
	}

	// Update the price with the form inputs
	date, _ := c.GetPostForm("date")
	p.Date = parseDate(date)
	price, _ := c.GetPostForm("price")
	p.Price = parseFloat(price)
	p.Comments, _ = c.GetPostForm("comments")

	// Some validation
	// TODO: don't use StatusOK for errors
	// TODO: flash in form?
	if p.Price <= 0 {
		c.String(http.StatusNotFound, "Price must be positive")
		return
	}
	if p.Date.Year() < 2000 {
		c.String(http.StatusNotFound, "Invalid or missing date")
		return
	}

	// Create or update price
	addUpdatePrice(p)

	// Go back to the stock page
	c.Redirect(http.StatusFound, fmt.Sprintf("/stock/%d", sid))
}
