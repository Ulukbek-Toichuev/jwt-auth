PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXISTS todos (
    todo_id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    is_done BOOLEAN,
    created_date DATETIME NOT NULL,
    deleted_date DATETIME,
    FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
) 