// Home page, with current portfolio and its current value

package main

import (
	//"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Show table of holdings with current value and return since purchase
func showPortfolio(c *gin.Context) {

	// Get portfolio holdings for today
	today := time.Now()
	holdings := getPortfolio(today, true)

	// Get cash value today
	var cash float64
	for _, c := range getAllCash(today) {
		cash += c.Amount
	}

	// Show page
	c.HTML(http.StatusOK, "portfolio.html",
		gin.H{"d": today, "holdings": holdings, "cash": cash,
			"menu": menu, "current": "Portfolio"})
}

// Portfolio holding a particular date
type Holding struct {
	Stock     Stock   // the asset held
	Units     float64 // quantity held
	UnitCost  float64 // avg price paid per unit
	CurPrice  float64 // current price in local currency
	TotCost   float64 // price paid cost in home currency
	CurValue  float64 // current value, in home currency
	Dividends float64 // total dividends received from this stock
	Return    float64 // percentage return since purchase
}

// Get holdings on a particular date, optionally only those held on that date
func getPortfolio(d time.Time, heldNow bool) []Holding {

	// Get all stocks, including those never or no longer held, and
	// go through transactions to determine holdings for each stock,
	// and total cost of the holdings
	holdings := []Holding{}
	for _, s := range getStocks() {

		// Accumulate holdings and average cost, up to a certain date
		var q, cost float64
		for _, t := range getTransactions(s.Id) {
			// Consider purchases before date, or sales after date
			if later(t.Date, d) {
				continue
			} else if t.Q > 0 { // purchase
				q += t.Q
				cost += t.Amount - t.Fees
			} else if t.Q < 0 { // sale
				q += t.Q              // add negative to reduce balance
				avgCost := cost / q   // average cost paid so far
				cost -= avgCost * t.Q // ??? correct?
			}
		}

		// Accumulate dividends
		var totDividends float64
		for _, d := range getDividends(s.Id) {
			totDividends += d.Amount
		}

		// If any of this stock currently held, calculate current value and return
		// and add it to list
		if q != 0 || !heldNow { // should never be negative, but just in case ...
			unitCost := cost / q            // TODO: risk of /0?
			curPrice := stockValue(s.Id, d) // current price
			curValue := q * curPrice
			gain := (curPrice-unitCost)*q + totDividends
			pcntUp := gain / cost * 100.0
			h := Holding{Stock: s, Units: q, UnitCost: unitCost, CurPrice: curPrice,
				CurValue: curValue, Dividends: totDividends, Return: pcntUp}
			holdings = append(holdings, h)
		}
	}
	return holdings
}

// Value of a stock on a date, just uses the last price before
// or on the date
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
	price := latestPriceAt(ts, d)

	// If not in home currency, get exchange rate on that date
	/*exchangeRate := 1.0 // will be 1 if already in home currency
	if stock.Currency != homeCurrency {
		ts := TimeSeries{}
		cur := getCurrencyCode(stock.Currency)
		if cur == nil {
			panic(fmt.Sprintf("stockValue: currency \"%s\" not found", stock.Currency))
		}
		for _, x := range getRates(cur.Id) {
			ts = append(ts, TimeSeriesPoint{x.Date, x.Rate})
		}
		exchangeRate = latestPriceAt(ts, d)
	}*/

	// Return value of the stock
	return price //* exchangeRate
}

// Units held of a stock on a certain date
func unitsHeld(sid int, d time.Time) float64 {
	var q float64
	for _, t := range getTransactions(sid) {
		if !later(t.Date, d) {
			q += t.Q
		}
	}
	return q
}
