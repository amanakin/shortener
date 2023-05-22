CREATE SCHEMA IF NOT EXISTS shortener;

-- TODO: add users and grants

-- TODO: add ID and expiration date
CREATE TABLE IF NOT EXISTS shortener.urls (
    original_url VARCHAR(255) NOT NULL UNIQUE,
    short_url VARCHAR(255) NOT NULL UNIQUE
);