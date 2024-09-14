-- schema.sql
--
-- Database schema for the portfolio system.
-- To create a new database:
-- sqlite3 data.db
-- .read schema.sql
-- Ctrl-D to exit

-- A stock or fund
CREATE TABLE stock (
    id integer primary key, 
    code text, 
    name text, 
    currency text);
create index stock_id on stock(id);
 create index stock_code on stock(code);

-- A daily price for a stock
CREATE TABLE price (
    id integer primary key, 
    stock_id integer,
    pdate date,
    price float,
    comments text);
create index price_id on price(id);
create index price_stock_id on price(stock_id);

-- A currency
CREATE TABLE currency (
    id integer primary key, 
    code text, 
    name text);
create index currency_id on currency(id);
create index currency_code on currency(code);

-- A daily rate for a currency, i.e., multiplier to get value in home currency
CREATE TABLE currency_rate (
    id integer primary key, 
    currency_id integer,
    rdate date,
    rate float);
create index rate_id on currency_rate(id);

-- A buy/sell transaction
CREATE TABLE trans (
    id integer primary key, 
    stock_id integer,
    tdate date,
    q float,
    amount float,
    fees float,
    comments text);
create index trans_id on trans(id);
create index trans_stock_id on trans(stock_id);

-- A dividend received for a stock, assumed to be an aggregate
-- amount (rather than per share), and in local currency (even if
-- the stock is in a foreign currency)
CREATE TABLE dividend (
    id integer primary key, 
    stock_id integer,
    tdate date,
    amount float,
    comments text);
create index div_id on dividend(id);
create index div_stock_id on dividend(stock_id);

-- A cash transaction (does not include buy/sell as these are implicit)
CREATE TABLE cash (
    id integer primary key, 
    tdate date,
    ttype text, -- deposit, withdraw, dividend
    amount float,
    comments text);
create index cash_id on cash(id);

