CREATE TABLE IF NOT EXISTS users (
    user_id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
	role TEXT NOT NULL,
    created_date DATETIME NOT NULL,
    deleted_date DATETIME
);