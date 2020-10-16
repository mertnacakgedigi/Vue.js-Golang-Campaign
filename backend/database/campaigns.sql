DROP TABLE IF EXISTS campaigns
CREATE TABLE campaigns (
    id INT PRIMARY KEY,
    name TEXT,
    status TEXT DEFAULT  'Paused',
    type TEXT,
    budget INT,
    created_at Text
);

