CREATE TABLE campaigns (
    id SERIAL PRIMARY KEY,
    name TEXT,
    status TEXT DEFAULT  'Paused',
    type TEXT,
    budget INT,
    created_on TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO campaigns (name, type,budget)
VALUES ('Generate Leads', 'Sponsored Ads',12000);