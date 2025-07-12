SELECT
    articles.title,
    comments.content
FROM
    articles
    JOIN articles_comments ON articles.id = articles_comments.article_id
    JOIN comments ON comments.id = articles_comments.comment_id
WHERE
    articles_comments.comment_id = 'e1446c48-c6f8-42df-9c4e-3357686c3a36';

/* All comments for an article. */
SELECT
    comments.id,
    comments.user_id,
    comments.content,
    comments.created_at,
    comments.last_modified
FROM
    comments
    JOIN articles_comments ON comments.id = articles_comments.comment_id
WHERE
    articles_comments.article_id = "db7c7311-0e82-4ef7-b296-18d7b8acd068";