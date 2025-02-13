DROP TABLE IF EXISTS status_checks;
DROP TABLE IF EXISTS submitted_flags;
DROP TABLE IF EXISTS flags;
DROP TABLE IF EXISTS endpoints;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS teams;
DROP TABLE IF EXISTS services;
DROP TABLE IF EXISTS ticks;

CREATE TABLE IF NOT EXISTS ticks (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	timestamp TEXT
);
CREATE TABLE IF NOT EXISTS services (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	service_name TEXT UNIQUE NOT NULL
);
CREATE TABLE IF NOT EXISTS teams (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT NOT NULL,
	key TEXT NOT NULL UNIQUE,
	is_admin BOOLEAN NOT NULL
);
CREATE TABLE IF NOT EXISTS users (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	user_name TEXT NOT NULL,
	team_id INTEGER,
	FOREIGN KEY (team_id) REFERENCES teams(id)
);
CREATE TABLE IF NOT EXISTS endpoints (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	team_id INTEGER NOT NULL,
	service_id INT NOT NULL,
	hostname TEXT NOT NULL,
	FOREIGN KEY (service_id) REFERENCES services(id)
);
CREATE TABLE IF NOT EXISTS flags (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	flag_identifier TEXT NOT NULL,
	flag TEXT UNIQUE NOT NULL,
	endpoint_id INTEGER NOT NULL,
	tick INTEGER NOT NULL,
	expiration TEXT NOT NULL,
	FOREIGN KEY (endpoint_id) REFERENCES endpoints(id)
);
CREATE TABLE IF NOT EXISTS submitted_flags (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	flag_id TEXT NOT NULL,
	team_id INTEGER NOT NULL,
	-- service_id INTEGER NOT NULL,
	-- tick INTEGER NOT NULL,
	timestamp TEXT NOT NULL,
	FOREIGN KEY (flag_id) REFERENCES flags(id),
	FOREIGN KEY (team_id) REFERENCES teams(id)
);
CREATE TABLE IF NOT EXISTS status_checks (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	tick INTEGER NOT NULL,
	endpoint_id INTEGER NOT NULL,
	status TEXT NOT NULL,
	timestamp TEXT NOT NULL,
	FOREIGN KEY (endpoint_id) REFERENCES endpoints(id),
	FOREIGN KEY (tick) REFERENCES ticks(id)
);


INSERT INTO ticks (timestamp) VALUES
	(datetime('now'));
INSERT INTO services (service_name) VALUES
	("gobrr"),
	("flightfinder"),
	("passgen");
INSERT INTO teams (id, name, key, is_admin) VALUES
	( 0, "Admin", "test", 1),
	( 1, "Team Awesome", "X7IC2UZD3DUS", 0),
	( 2, "Winning", "RWVMD7N2CY46", 0);
INSERT INTO users (user_name, team_id) VALUES
	("admin", 0);
INSERT INTO endpoints (team_id, service_id, hostname) VALUES
	( 1, 1, "service1.ad.mctf.io:8001"),
	( 1, 2, "service2.ad.mctf.io:8001"),
	( 1, 3, "service3.ad.mctf.io:8001"),
	( 2, 1, "service1.ad.mctf.io:8002"),
	( 2, 2, "service2.ad.mctf.io:8002"),
	( 2, 3, "service3.ad.mctf.io:8002");
INSERT INTO flags (flag_identifier, flag, endpoint_id, tick, expiration)VALUES
	('flag001', "mctf{I2P857iTQKvKYUks}",  1, 1, datetime('now', '+1 hour')),
	('flag002', "mctf{VGiXvmN2iJ/Fwzai}",  2, 1, datetime('now', '+1 hour')),
	('flag003', "mctf{TgT6h0Nb1/2MFz6J}",  3, 1, datetime('now', '+1 hour')),
	('flag004', "mctf{y3kEJSN5zl1IMS/7}",  4, 1, datetime('now', '+1 hour')),
	('flag005', "mctf{y9xhqRKQwbkJmEWO}",  5, 1, datetime('now', '+1 hour')),
	('flag006', "mctf{0Ipuq2xCWjDb83zH}",  6, 1, datetime('now', '+1 hour'));
INSERT INTO submitted_flags (flag_id, team_id, timestamp) VALUES
	(1, 1, unixepoch());
INSERT INTO status_checks (tick, endpoint_id, status, timestamp) VALUES
	( 1, 1, "success", datetime('now')),
	( 1, 2, "success", datetime('now')),
	( 1, 3, "success", datetime('now')),
	( 1, 4, "failure", datetime('now')),
	( 1, 5, "success", datetime('now')),
	( 1, 6, "unknown", datetime('now'));
