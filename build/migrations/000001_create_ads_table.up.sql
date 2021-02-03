CREATE TABLE IF NOT EXISTS ads (
    id serial PRIMARY KEY,
    name VARCHAR(200) NOT NULL,
    description VARCHAR(1000) NOT NULL,
    price numeric NOT NULL,
    created_at timestamp DEFAULT current_timestamp
);