CREATE TABLE
    IF NOT EXISTS customers (
        id SERIAL PRIMARY KEY,
        name VARCHAR(250) NOT NULL,
        created_at TIME
        WITH
            TIME ZONE NOT NULL
    );