CREATE TABLE customers (
	id SERIAL PRIMARY KEY,
    firstname VARCHAR(50),
	lastname VARCHAR(50),
    phone VARCHAR(25) UNIQUE NOT NULL,
	password VARCHAR(255),
    email VARCHAR(255) UNIQUE NOT NULL,
	created_at TIMESTAMP NOT NULL, 
    last_login TIMESTAMP
);