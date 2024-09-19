CREATE TABLE IF NOT EXISTS message (
	id         integer PRIMARY KEY AUTOINCREMENT,
	created_at TIMESTAMP,
	client_ip  VARCHAR(64),
	service    VARCHAR(32),
	username   VARCHAR(64),
	password   VARCHAR(64),
	command    TEXT
);

CREATE INDEX IF NOT EXISTS idx_message_service   ON message (service);
CREATE INDEX IF NOT EXISTS idx_message_username  ON message (username);
CREATE INDEX IF NOT EXISTS idx_message_client_ip ON message (client_ip);
