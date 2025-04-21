CREATE TABLE IF NOT EXISTS bookmarks
    (
        id INTEGER PRIMARY KEY,
        title TEXT NOT NULL,
        url TEXT UNIQUE NOT NULL,
        tags TEXT DEFAULT "[]",
        created_at DATETIME NOT NULL,
        updated_at DATETIME NOT NULL
    );
