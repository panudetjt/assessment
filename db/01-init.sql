CREATE TABLE IF NOT EXISTS "expenses" (
    id SERIAL PRIMARY KEY,
    title TEXT,
    amount INT,
    note TEXT,
    tags TEXT[] 
);