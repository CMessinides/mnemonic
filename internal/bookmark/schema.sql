CREATE TABLE IF NOT EXISTS bookmarks
    (
        id INTEGER PRIMARY KEY,
        title TEXT NOT NULL,
        url TEXT UNIQUE NOT NULL,
        tags TEXT DEFAULT "[]",
        created_at DATETIME NOT NULL,
        updated_at DATETIME NOT NULL,
        archived_at DATETIME
    );

CREATE VIEW IF NOT EXISTS active_bookmarks
    AS SELECT id, title, url, tags, created_at, updated_at, (archived_at IS NOT NULL) archived
    FROM bookmarks
    WHERE archived_at IS NULL;

CREATE VIEW IF NOT EXISTS all_bookmarks
    AS SELECT id, title, url, tags, created_at, updated_at, (archived_at IS NOT NULL) archived
    FROM bookmarks;
