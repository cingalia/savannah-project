CREATE TABLE orders (
	id SERIAL PRIMARY KEY,
    customer_id integer REFERENCES customers (id),
	item_ids VARCHAR(50) NOT NULL,
    summary VARCHAR(100) NOT NULL,
    total_price integer NOT NULL,
	created_at TIMESTAMP NOT NULL DEFAULT NOW()
);