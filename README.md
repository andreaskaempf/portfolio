# Portfolio System

This will be a simple portfolio management system for personal investors. It will allow you to
define assets such as stocks, funds, etc., and record transactions to buy/sell shares. You
can update the prices, record dividends, and update currency exchange rates. The tool will then
allow you to see to the total value of your portfolio over time, and report the ROI coming from
asset value changes, dividend income, and currency movements.

This is written 100% in Go, including the front end, which is just static HTML.

To build and run:
```
go mod init portfolio
go build
./portfolio
```

Then, you can browse to http://localhost:8080

Note that you need to create a database first, using SQLite3:

```
sqlite3 data.db
.read schema.sql
.exit
```

AK, July & August 2024
