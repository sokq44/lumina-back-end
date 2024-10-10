CREATE DATABASE IF NOT EXISTS articles;

USE articles;

CREATE TABLE
    users (
        id VARCHAR(36) PRIMARY KEY NOT NULL DEFAULT (uuid ()),
        name VARCHAR(50) NOT NULL,
        email VARCHAR(255) NOT NULL,
        password VARCHAR(64) NOT NULL,
        verified BOOLEAN NOT NULL
    )