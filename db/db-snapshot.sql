CREATE DATABASE IF NOT EXISTS articles;

USE articles;

CREATE TABLE
    IF NOT EXISTS users
(
    id        VARCHAR(36) PRIMARY KEY NOT NULL DEFAULT (uuid()),
    username  VARCHAR(50)             NOT NULL,
    email     VARCHAR(255)            NOT NULL,
    image_url VARCHAR(255)            NOT NULL,
    password  VARCHAR(64)             NOT NULL,
    verified  BOOLEAN                 NOT NULL DEFAULT FALSE
);

CREATE TABLE
    IF NOT EXISTS email_verification
(
    id      VARCHAR(36) PRIMARY KEY NOT NULL DEFAULT (uuid()),
    token   VARCHAR(128)            NOT NULL,
    expires DATETIME                NOT NULL DEFAULT CURRENT_TIMESTAMP,
    user_id VARCHAR(36)             NOT NULL,

    CONSTRAINT fk_users_email_verification FOREIGN KEY (user_id) REFERENCES users (id),
    CONSTRAINT email_verification_unique_user_id UNIQUE (user_id)
);

CREATE TABLE
    IF NOT EXISTS refresh_tokens
(
    id      VARCHAR(36) PRIMARY KEY NOT NULL DEFAULT (uuid()),
    token   VARCHAR(255)            NOT NULL,
    expires DATETIME                NOT NULL DEFAULT CURRENT_TIMESTAMP,
    user_id VARCHAR(36)             NOT NULL,

    CONSTRAINT fk_users_refresh_tokens FOREIGN KEY (user_id) REFERENCES users (id),
    CONSTRAINT refresh_tokens_unique_user_id UNIQUE (user_id)
);

CREATE TABLE
    IF NOT EXISTS password_change
(
    id      VARCHAR(36) PRIMARY KEY NOT NULL DEFAULT (uuid()),
    token   VARCHAR(128)            NOT NULL,
    expires DATETIME                NOT NULL DEFAULT CURRENT_TIMESTAMP,
    user_id VARCHAR(36)             NOT NULL,

    CONSTRAINT fk_users_password_change FOREIGN KEY (user_id) REFERENCES users (id),
    CONSTRAINT password_change_unique_user_id UNIQUE (user_id)
);

CREATE TABLE
    IF NOT EXISTS secrets
(
    id      VARCHAR(36) PRIMARY KEY NOT NULL DEFAULT (uuid()),
    secret  VARCHAR(64)             NOT NULL,
    expires DATETIME                NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS articles
(
    id         VARCHAR(36) PRIMARY KEY NOT NULL DEFAULT (uuid()),
    title      VARCHAR(36)             NOT NULL,
    content    LONGTEXT                NOT NULL,
    user_id    VARCHAR(36)             NOT NULL,
    banner_url VARCHAR(255)            NOT NULL,
    public     BOOL                    NOT NULL DEFAULT FALSE,
    created_at DATETIME                NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_users_articles FOREIGN KEY (user_id) REFERENCES users (id)
)