CREATE TABLE orders (
	id SERIAL PRIMARY KEY,
    customerid VARCHAR(50),
	item VARCHAR(50),
    description VARCHAR(100),
    price VARCHAR(25),
	date  timestamp NOT NULL DEFAULT NOW()
);