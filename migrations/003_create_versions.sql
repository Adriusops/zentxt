CREATE TABLE versions (
    id TEXT PRIMARY KEY, --UUID
    file_id TEXT REFERENCES files(id),
    version_number INTEGER NOT NULL,
    path TEXT NOT NULL,
    author TEXT NULL,
    message TEXT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
