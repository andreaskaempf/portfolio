-- schema.sql
--
-- Database schema for the portfolio system.
-- To create a new database:
-- sqlite3 data.db
-- .read schema.sql
-- Ctrl-D to exit

-- An stock or fund
CREATE TABLE stock (
    id integer primary key, 
    code text, 
    name text, 
    currency text);
create index stock_id on stock(id);

