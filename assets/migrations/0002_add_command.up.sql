ALTER TABLE message ADD COLUMN command VARCHAR(128);
CREATE INDEX IF NOT EXISTS idx_message_command ON message(command);
