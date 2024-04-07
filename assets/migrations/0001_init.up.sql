CREATE TABLE IF NOT EXISTS message (
	id SERIAL,
	service VARCHAR(32) NOT NULL,
	username VARCHAR(64),
	password VARCHAR(64),
	client_ip VARCHAR(64),
	created_at INT NOT NULL,

	PRIMARY KEY (id),
	FOREIGN KEY (service) REFERENCES service(name) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_message_created_at ON message(created_at);
CREATE INDEX IF NOT EXISTS idx_message_service ON message(service);
CREATE INDEX IF NOT EXISTS idx_message_username ON message(username);
CREATE INDEX IF NOT EXISTS idx_message_client_ip ON message(client_ip);
