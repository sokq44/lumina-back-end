CREATE DATABASE IF NOT EXISTS lumina;

USE lumina;

CREATE TABLE
    IF NOT EXISTS users (
        id VARCHAR(36) PRIMARY KEY NOT NULL DEFAULT (uuid ()),
        username VARCHAR(50) NOT NULL,
        email VARCHAR(255) NOT NULL,
        image_url VARCHAR(255) NOT NULL,
        password VARCHAR(64) NOT NULL,
        verified BOOLEAN NOT NULL DEFAULT FALSE
    );

CREATE TABLE
    IF NOT EXISTS email_verification (
        id VARCHAR(36) PRIMARY KEY NOT NULL DEFAULT (uuid ()),
        token VARCHAR(128) NOT NULL,
        expires DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        user_id VARCHAR(36) NOT NULL,
        CONSTRAINT fk_users_email_verification FOREIGN KEY (user_id) REFERENCES users (id),
        CONSTRAINT email_verification_unique_user_id UNIQUE (user_id)
    );

CREATE TABLE
    IF NOT EXISTS refresh_tokens (
        id VARCHAR(36) PRIMARY KEY NOT NULL DEFAULT (uuid ()),
        token VARCHAR(255) NOT NULL,
        expires DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        user_id VARCHAR(36) NOT NULL,
        CONSTRAINT fk_users_refresh_tokens FOREIGN KEY (user_id) REFERENCES users (id),
        CONSTRAINT refresh_tokens_unique_user_id UNIQUE (user_id)
    );

CREATE TABLE
    IF NOT EXISTS password_change (
        id VARCHAR(36) PRIMARY KEY NOT NULL DEFAULT (uuid ()),
        token VARCHAR(128) NOT NULL,
        expires DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        user_id VARCHAR(36) NOT NULL,
        CONSTRAINT fk_users_password_change FOREIGN KEY (user_id) REFERENCES users (id),
        CONSTRAINT password_change_unique_user_id UNIQUE (user_id)
    );

CREATE TABLE
    IF NOT EXISTS email_change (
        id VARCHAR(36) PRIMARY KEY NOT NULL DEFAULT (uuid ()),
        token VARCHAR(128) NOT NULL,
        new_email VARCHAR(255) NOT NULL,
        expires DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        user_id VARCHAR(36) NOT NULL,
        CONSTRAINT fk_users_email_change FOREIGN KEY (user_id) REFERENCES users (id),
        CONSTRAINT email_change_unique_user_id UNIQUE (user_id)
    );

CREATE TABLE
    IF NOT EXISTS secrets (
        id VARCHAR(36) PRIMARY KEY NOT NULL DEFAULT (uuid ()),
        secret VARCHAR(64) NOT NULL,
        expires DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
    );

CREATE TABLE
    IF NOT EXISTS articles (
        id VARCHAR(36) PRIMARY KEY NOT NULL DEFAULT (uuid ()),
        title VARCHAR(36) NOT NULL,
        content LONGTEXT NOT NULL,
        user_id VARCHAR(36) NOT NULL,
        banner_url VARCHAR(255) NOT NULL,
        public BOOL NOT NULL DEFAULT FALSE,
        created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        CONSTRAINT fk_users_articles FOREIGN KEY (user_id) REFERENCES users (id)
    );

CREATE TABLE
    IF NOT EXISTS comments (
        id VARCHAR(36) PRIMARY KEY NOT NULL DEFAULT (uuid ()),
        user_id VARCHAR(36) NOT NULL,
        content VARCHAR(255) NOT NULL,
        created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        last_modified DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        CONSTRAINT fk_users_comments FOREIGN KEY (user_id) REFERENCES users (id)
    );

CREATE TABLE
    IF NOT EXISTS discussions (
        id VARCHAR(36) PRIMARY KEY NOT NULL DEFAULT (uuid ()),
        created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
    );

CREATE TABLE
    IF NOT EXISTS discussions_comments (
        id VARCHAR(36) PRIMARY KEY NOT NULL DEFAULT (uuid ()),
        discussion_id VARCHAR(36) NOT NULL,
        comment_id VARCHAR(36) NOT NULL,
        CONSTRAINT fk_discussions_discussions_comments FOREIGN KEY (discussion_id) REFERENCES discussions (id),
        CONSTRAINT fk_comments_discussions_comments FOREIGN KEY (comment_id) REFERENCES comments (id)
    );

CREATE TABLE
    IF NOT EXISTS articles_comments (
        id VARCHAR(36) PRIMARY KEY NOT NULL DEFAULT (uuid ()),
        article_id VARCHAR(36) NOT NULL,
        comment_id VARCHAR(36) NOT NULL,
        CONSTRAINT fk_articles_articles_comments FOREIGN KEY (article_id) REFERENCES articles (id),
        CONSTRAINT fk_comments_articles_comments FOREIGN KEY (comment_id) REFERENCES comments (id)
    );

CREATE TABLE
    IF NOT EXISTS articles_discussions (
        id VARCHAR(36) PRIMARY KEY NOT NULL DEFAULT (uuid ()),
        article_id VARCHAR(36) NOT NULL,
        discussion_id VARCHAR(36) NOT NULL,
        CONSTRAINT fk_articles_articles_discussions FOREIGN KEY (article_id) REFERENCES articles (id),
        CONSTRAINT fk_discussions_articles_discussions FOREIGN KEY (discussion_id) REFERENCES discussions (id)
    );