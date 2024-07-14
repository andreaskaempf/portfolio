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

