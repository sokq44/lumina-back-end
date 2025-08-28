CREATE DATABASE IF NOT EXISTS lumina;

USE lumina;

CREATE TABLE
    IF NOT EXISTS users (
        id VARCHAR(36) PRIMARY KEY NOT NULL DEFAULT (uuid ()),
        username VARCHAR(50) NOT NULL DEFAULT '',
        email VARCHAR(255) NOT NULL DEFAULT '',
        image_url VARCHAR(255) NOT NULL DEFAULT '',
        bio VARCHAR(255) NOT NULL DEFAULT '',
        favourites VARCHAR(255) NOT NULL DEFAULT '',
        password VARCHAR(64) NOT NULL DEFAULT '',
        verified BOOLEAN NOT NULL DEFAULT FALSE,
        CONSTRAINT user_unique_username UNIQUE (username),
        CONSTRAINT user_unique_email UNIQUE (email)
    );

CREATE TABLE
    IF NOT EXISTS socials (
        type INT PRIMARY KEY NOT NULL AUTO_INCREMENT,
        name VARCHAR(64) NOT NULL DEFAULT ''
    );

CREATE TABLE
    IF NOT EXISTS user_socials (
        id VARCHAR(36) PRIMARY KEY NOT NULL DEFAULT (uuid ()),
        user_id VARCHAR(36) NOT NULL DEFAULT '',
        social_type INT NOT NULL,
        value VARCHAR(255) NOT NULL DEFAULT '',
        label VARCHAR(255) NOT NULL DEFAULT '',
        CONSTRAINT fk_users_user_socials FOREIGN KEY (user_id) REFERENCES users (id),
        CONSTRAINT fk_socials_user_socials FOREIGN KEY (social_type) REFERENCES socials (type)
    );

CREATE TABLE
    IF NOT EXISTS email_verification (
        id VARCHAR(36) PRIMARY KEY NOT NULL DEFAULT (uuid ()),
        token VARCHAR(128) NOT NULL DEFAULT '',
        expires DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        user_id VARCHAR(36) NOT NULL DEFAULT '',
        CONSTRAINT fk_users_email_verification FOREIGN KEY (user_id) REFERENCES users (id),
        CONSTRAINT email_verification_unique_user_id UNIQUE (user_id)
    );

CREATE TABLE
    IF NOT EXISTS refresh_tokens (
        id VARCHAR(36) PRIMARY KEY NOT NULL DEFAULT (uuid ()),
        token VARCHAR(255) NOT NULL DEFAULT '',
        expires DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        user_id VARCHAR(36) NOT NULL DEFAULT '',
        CONSTRAINT fk_users_refresh_tokens FOREIGN KEY (user_id) REFERENCES users (id),
    );

CREATE TABLE
    IF NOT EXISTS password_change (
        id VARCHAR(36) PRIMARY KEY NOT NULL DEFAULT (uuid ()),
        token VARCHAR(128) NOT NULL DEFAULT '',
        expires DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        user_id VARCHAR(36) NOT NULL DEFAULT '',
        CONSTRAINT fk_users_password_change FOREIGN KEY (user_id) REFERENCES users (id),
        CONSTRAINT password_change_unique_user_id UNIQUE (user_id)
    );

CREATE TABLE
    IF NOT EXISTS email_change (
        id VARCHAR(36) PRIMARY KEY NOT NULL DEFAULT (uuid ()),
        token VARCHAR(128) NOT NULL DEFAULT '',
        new_email VARCHAR(255) NOT NULL DEFAULT '',
        expires DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        user_id VARCHAR(36) NOT NULL DEFAULT '',
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
        title VARCHAR(36) NOT NULL DEFAULT '',
        content LONGTEXT NOT NULL DEFAULT '',
        tldr VARCHAR(512) NOT NULL DEFAULT '',
        user_id VARCHAR(36) NOT NULL DEFAULT '',
        banner_url VARCHAR(255) NOT NULL DEFAULT '',
        public BOOL NOT NULL DEFAULT FALSE,
        created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        CONSTRAINT fk_users_articles FOREIGN KEY (user_id) REFERENCES users (id)
    );

CREATE TABLE
    IF NOT EXISTS comments (
        id VARCHAR(36) PRIMARY KEY NOT NULL DEFAULT (uuid ()),
        user_id VARCHAR(36) NOT NULL DEFAULT '',
        content VARCHAR(255) NOT NULL DEFAULT '',
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
        discussion_id VARCHAR(36) NOT NULL DEFAULT '',
        comment_id VARCHAR(36) NOT NULL DEFAULT '',
        CONSTRAINT fk_discussions_discussions_comments FOREIGN KEY (discussion_id) REFERENCES discussions (id),
        CONSTRAINT fk_comments_discussions_comments FOREIGN KEY (comment_id) REFERENCES comments (id)
    );

CREATE TABLE
    IF NOT EXISTS articles_reads (
        user_id VARCHAR(36) NOT NULL DEFAULT '',
        article_id VARCHAR(36) NOT NULL DEFAULT '',
        CONSTRAINT fk_users_articles_reads FOREIGN KEY (user_id) REFERENCES users (id),
        CONSTRAINT fk_articles_articles_reads FOREIGN KEY (article_id) REFERENCES articles (id),
        CONSTRAINT articles_reads_unique_user_article UNIQUE (user_id, article_id)
    );

CREATE TABLE
    IF NOT EXISTS articles_ratings (
        user_id VARCHAR(36) NOT NULL DEFAULT '',
        article_id VARCHAR(36) NOT NULL DEFAULT '',
        rating TINYINT UNSIGNED NOT NULL,
        CONSTRAINT fk_users_article_ratings FOREIGN KEY (user_id) REFERENCES users (id),
        CONSTRAINT fk_articles_article_ratings FOREIGN KEY (article_id) REFERENCES articles (id),
    );

CREATE TABLE
    IF NOT EXISTS articles_comments (
        id VARCHAR(36) PRIMARY KEY NOT NULL DEFAULT (uuid ()),
        article_id VARCHAR(36) NOT NULL DEFAULT '',
        comment_id VARCHAR(36) NOT NULL DEFAULT '',
        CONSTRAINT fk_articles_articles_comments FOREIGN KEY (article_id) REFERENCES articles (id),
        CONSTRAINT fk_comments_articles_comments FOREIGN KEY (comment_id) REFERENCES comments (id)
    );

CREATE TABLE
    IF NOT EXISTS articles_discussions (
        id VARCHAR(36) PRIMARY KEY NOT NULL DEFAULT (uuid ()),
        article_id VARCHAR(36) NOT NULL DEFAULT '',
        discussion_id VARCHAR(36) NOT NULL DEFAULT '',
        CONSTRAINT fk_articles_articles_discussions FOREIGN KEY (article_id) REFERENCES articles (id),
        CONSTRAINT fk_discussions_articles_discussions FOREIGN KEY (discussion_id) REFERENCES discussions (id)
    );