CREATE DATABASE IF NOT EXISTS bookmarks;

USE "bookmarks";

CREATE TABLE IF NOT EXISTS bookmarks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    bookmark_id TEXT NOT NULL,
    url TEXT NOT NULL,
    title TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

/* add index for url */
CREATE UNIQUE INDEX IF NOT EXISTS bookmarks_url_idx ON bookmarks (url);

/* add index for bookmark_id */
CREATE UNIQUE INDEX IF NOT EXISTS bookmarks_bookmark_id_idx ON bookmarks (bookmark_id);
