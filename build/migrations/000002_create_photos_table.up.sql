CREATE TABLE IF NOT EXISTS photos (
    id serial PRIMARY KEY,
    ad_id int NOT NULL,
    url VARCHAR(256) NOT NULL,
    FOREIGN KEY (ad_id) REFERENCES ads (id) ON DELETE CASCADE
);

