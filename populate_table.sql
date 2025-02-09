DROP TABLE albums;
DROP TABLE api_keys;

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
INSERT INTO api_keys(key)
VALUES
	("test"),
	("a961e4b54588049c8ec3490b69fd94543332ff069a70b9b59e217db04398416e");
