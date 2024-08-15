// Home page, with current portfolio and its current value

package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Show table of holdings with current value and return since purchase
func showPortfolio(c *gin.Context) {

	// Get portfolio for today
	today := time.Now()
	holdings := getPortfolio(today)

	// Show page
	c.HTML(http.StatusOK, "portfolio.html",
		gin.H{"d": today, "holdings": holdings, "menu": menu})
}

// Portfolio holding a particular date
type Holding struct {
	Stock  Stock   // the asset held
	Units  float64 // quantity held
	Price  float64 // average price in local currency
	Cost   float64 // total cost in home currency
	Value  float64 // current value, in home currency
	Return float64 // percentage return since purchase
}

// Get holdings on a particular date
func getPortfolio(d time.Time) []Holding {

	// Get all stocks, including those never or no longer held, and
	// go through transactions to determine holdings for each stock,
	// and total cost of the holdings
	holdings := []Holding{}
	for _, s := range getStocks() {

		// Accumulate holdings and average cost
		var q, cost float64
		for _, t := range getTransactions(s.Id) {
			// Consider purchases before date, or sales after date
			if t.Q > 0 && (sameDate(t.Date, d) || earlier(t.Date, d)) {
				q += t.Q
				cost += t.Amount - t.Fees
			} else if t.Q < 0 && (sameDate(t.Date, d) || later(t.Date, d)) {
				q += t.Q              // add negative to reduce balance
				avgCost := cost / q   // average cost paid so far
				cost -= avgCost * t.Q // ??? correct?
			}
		}

		// If any of this stock currently held, calculate current value and return
		// and add it to list
		if q != 0 { // should never be negative, but just in case ...
			curValue := q * stockValue(s.Id, d)
			pcntUp := ((curValue - cost) / cost) * 100.0
			h := Holding{Stock: s, Units: q, Cost: cost, Value: curValue, Return: pcntUp}
			holdings = append(holdings, h)
		}
	}
	return holdings
}

// Value of a stock on a date, in home currency
func stockValue(sid int, d time.Time) float64 {

	// Get the stock
	stock := getStock(sid)
	if stock == nil {
		panic("stockValue: stock not found")
		return 0
	}

	// Get the (approximate) price of the stock on given date
	ts := TimeSeries{}
	for _, p := range getPrices(sid) {
		ts = append(ts, TimeSeriesPoint{p.Date, p.Price})
	}
	price := priceOn(ts, d)

	// If not in home currency, get exchange rate on that date
	exchangeRate := 1.0 // will be 1 if already in home currency
	if stock.Currency != homeCurrency {
		ts := TimeSeries{}
		cur := getCurrencyCode(stock.Currency)
		if cur == nil {
			panic(fmt.Sprintf("stockValue: currency \"%s\" not found", stock.Currency))
		}
		for _, x := range getRates(cur.Id) {
			ts = append(ts, TimeSeriesPoint{x.Date, x.Rate})
		}
		exchangeRate = priceOn(ts, d)
	}

	// Return value of the stock
	return price * exchangeRate
}
