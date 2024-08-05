package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Add a rate for a currency
func editRate(c *gin.Context) {

	// Get Currency ID from query string if provided (only required
	// for adding new prices)
	// If there is no rate ID, expect a currency ID as query parameter
	var cid int
	cid_, ok := c.GetQuery("cid")
	if ok {
		cid = parseInt(cid_)
		if cid < 0 {
			c.String(http.StatusNotFound, "Invalid currency ID")
			return
		}
	}

	// Get the Rate ID. If this is zero, means to add a new price to
	// the current Currency ID (see below), otherwise get the price to edit.
	rid := parseInt(c.Param("rid"))
	var r Rate
	if rid < 0 {
		c.String(http.StatusNotFound, "Invalid rate ID")
		return
	} else if rid == 0 { // create a new rate
		if cid <= 0 {
			c.String(http.StatusNotFound, "Cannot add rate without currency ID")
			return
		}
		newRate := 0.0 // TODO: default rate sould be most recent one
		r = Rate{Currency: cid, Date: time.Now(), Rate: newRate}
	} else { // get existing rate
		rp := getRate(rid)
		if rp == nil {
			c.String(http.StatusNotFound, "Rate not found")
			return
		}
		r = *rp
		rid = r.Id
	}

	// Show the form to edit rate
	c.HTML(http.StatusOK, "edit_rate.html",
		gin.H{"r": r, "ds": formatDate(r.Date), "menu": menu})
}

// Create or update a rate
func updateRate(c *gin.Context) {

	// Get currency and rate ID (latter will be 0 to add a rate)
	cid_, ok := c.GetPostForm("cid")
	if !ok {
		c.String(http.StatusOK, "saveRate: Missing currency ID")
		return
	}
	rid_, ok := c.GetPostForm("rid")
	if !ok {
		c.String(http.StatusOK, "saveRate: Missing rate ID")
		return
	}
	cid := parseInt(cid_)
	rid := parseInt(rid_)
	if cid < 0 || rid < 0 {
		c.String(http.StatusOK, "saveRate: Invalid currency or rate ID")
		return
	}

	// Get the rate or create "blank" one
	r := &Rate{Currency: cid}
	if rid > 0 {
		r = getRate(rid)
		if r == nil {
			c.String(http.StatusOK, "saveRate: rate not found")
			return
		}
	}

	// Update the rate with the form inputs
	date, _ := c.GetPostForm("date")
	r.Date = parseDate(date)
	rate, _ := c.GetPostForm("rate")
	r.Rate = parseFloat(rate)

	// Some validation
	if r.Rate <= 0 {
		c.String(http.StatusNotFound, "Rate must be positive")
		return
	}
	if r.Date.Year() < 2000 {
		c.String(http.StatusNotFound, "Invalid or missing date")
		return
	}

	// Create or update rate
	addUpdateRate(r)

	// Go back to the currency page
	c.Redirect(http.StatusFound, fmt.Sprintf("/currency/%d", cid))
}
