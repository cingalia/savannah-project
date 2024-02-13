CREATE TABLE orders (
	id SERIAL PRIMARY KEY,
    customer_id integer REFERENCES customers (id),
	item VARCHAR(50) NOT NULL,
    description VARCHAR(100) NOT NULL,
    price integer NOT NULL,
	created_at TIMESTAMP NOT NULL DEFAULT NOW()
);