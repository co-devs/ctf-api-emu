CREATE TABLE IF NOT EXISTS albums (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    artist TEXT NOT NULL,
    price REAL NOT NULL
);
INSERT INTO albums (title, artist, price)
VALUES
	("Blue Train", "John Coltrane", 56.99),
	("Jeru", "Gerry Mulligan", 17.99),
	("Sarah Vaughan and Clifford Brown", "Sarah Vaughan", 39.99);
