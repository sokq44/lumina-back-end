SELECT
    articles.title,
    comments.content
FROM
    articles
    JOIN articles_comments ON articles.id = articles_comments.article_id
    JOIN comments ON comments.id = articles_comments.comment_id
WHERE
    articles_comments.comment_id = 'e1446c48-c6f8-42df-9c4e-3357686c3a36';