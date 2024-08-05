-- schema.sql
--
-- Database schema for the portfolio system.
-- To create a new database:
-- sqlite3 data.db
-- .read schema.sql
-- Ctrl-D to exit

-- A stock or fund
--CREATE TABLE stock (
--    id integer primary key, 
--    code text, 
--    name text, 
--    currency text);
--create index stock_id on stock(id);
-- create index stock_code on stock(code);

-- A daily price for a stock
--CREATE TABLE price (
--    id integer primary key, 
--    stock_id integer,
--    pdate date,
--    price float,
--    comments text);
--create index price_id on price(id);
--create index price_stock_id on price(stock_id);

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
