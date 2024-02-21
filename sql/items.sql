CREATE TABLE items (
	id SERIAL PRIMARY KEY,
	name VARCHAR(50) NOT NULL,
    description VARCHAR(100) NOT NULL,
    price integer NOT NULL,
    quantity integer NOT NULL,
	created_at TIMESTAMP NOT NULL DEFAULT NOW()
);