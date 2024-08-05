package main

import (
	"fmt"
	"net/http"

	"strings"

	"github.com/gin-gonic/gin"
)

// Show list of all currencies
func showCurrencies(c *gin.Context) {

	// Get a list of currencies
	currencies := getCurrencies()

	// Show page
	c.HTML(http.StatusOK, "currencies.html",
		gin.H{"currencies": currencies, "menu": menu})
}

// Page to show one currency
func showCurrency(c *gin.Context) {

	// Parse the ID and get the currency
	cid := parseInt(c.Param("id"))
	cur := getCurrency(cid)
	if cur == nil {
		c.String(http.StatusNotFound, fmt.Sprintf("Currency %d not found", cid))
		return
	}

	// Get all rates for this currency
	rates := getRates(cid)

	// Show page
	c.HTML(http.StatusOK, "currency.html",
		gin.H{"cur": cur, "rates": rates, "menu": menu})
}

// Show form to edit a currency (including a new one)
func editCurrency(c *gin.Context) {

	// Get currency ID (will be 0 to add an currency)
	cid := parseInt(c.Param("id"))
	if cid < 0 {
		c.String(http.StatusNotFound, "Invalid currency ID")
		return
	}

	// Get the currency or create "blank" currency
	cur := &Currency{}
	if cid > 0 {
		cur = getCurrency(cid)
		if cur == nil {
			c.String(http.StatusNotFound, "Currency not found")
			return
		}
	}

	// Show the form to edit currency
	c.HTML(http.StatusOK, "edit_currency.html",
		gin.H{"cur": cur, "menu": menu})
}

// Process form to update or add an currency
func saveCurrency(c *gin.Context) {

	// Get currency ID (will be 0 to add a currency)
	cid_, ok := c.GetPostForm("cid")
	if !ok {
		c.String(http.StatusNotFound, "saveCurrency: Missing currency ID")
		return
	}
	cid := parseInt(cid_)
	if cid < 0 {
		c.String(http.StatusNotFound, "saveCurrency: Invalid currency ID")
		return
	}

	// Get the currency or create "blank" currency
	cur := &Currency{}
	if cid > 0 {
		cur = getCurrency(cid)
		if cur == nil {
			c.String(http.StatusNotFound, "saveCurrency: currency not found")
			return
		}
	}

	// Update the currency with the form inputs
	cur.Code, _ = c.GetPostForm("code")
	cur.Name, _ = c.GetPostForm("name")

	// Some validation
	// TODO: make sure can't create a currency that already eixsts,
	// including renaming a currency to name of an existing one
	cur.Code = strings.TrimSpace(cur.Code)
	cur.Name = strings.TrimSpace(cur.Name)
	if len(cur.Code) == 0 || len(cur.Name) == 0 {
		c.String(http.StatusNotFound, "Invalid inputs: cannot be blank")
		return
	}

	// Create or update currency in database
	addUpdateCurrency(cur)

	// Go back to currencies page
	c.Redirect(http.StatusFound, "/Currencies")

}

// Delete currency: ask for confirmation first
func delCurrency(c *gin.Context) {

	// Get the currency (URL positional param)
	cid := parseInt(c.Param("id"))
	cur := getCurrency(cid)
	if cid <= 0 || cur == nil {
		c.String(http.StatusNotFound, "Currency not found")
		return
	}

	// Ask for confirmation, or go ahead and delete if confirmed
	confirm, _ := c.GetQuery("confirm")
	if confirm == "" { // no confirmation, show form
		c.HTML(http.StatusOK, "del_currency.html",
			gin.H{"cur": cur, "menu": menu})
	} else if confirm == "yes" { // confirmed, delete currency
		deleteCurrency(cid)
		c.Redirect(http.StatusFound, "/Currencies")
	} else { // confirmation denied, back to currency page
		c.Redirect(http.StatusFound, fmt.Sprintf("/currency/%d", cid))
	}
}
