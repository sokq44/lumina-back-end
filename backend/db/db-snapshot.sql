CREATE DATABASE IF NOT EXISTS articles;

USE articles;

CREATE TABLE
    IF NOT EXISTS users (
        id VARCHAR(36) PRIMARY KEY NOT NULL DEFAULT (uuid ()),
        username VARCHAR(50) NOT NULL,
        email VARCHAR(255) NOT NULL,
        password VARCHAR(64) NOT NULL,
        verified BOOLEAN NOT NULL DEFAULT FALSE
    );

CREATE TABLE
    IF NOT EXISTS email_validation (
        id VARCHAR(36) PRIMARY KEY NOT NULL DEFAULT (uuid ()),
        token VARCHAR(128) NOT NULL,
        expires DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        user_id VARCHAR(36) NOT NULL,
        CONSTRAINT fk_users_email_validation FOREIGN KEY (user_id) REFERENCES users (id)
    );