-- +migrate Down
ALTER TABLE user_settings DROP COLUMN two_steps_verif;

-- +migrate Up
ALTER TABLE user_settings ADD settings json NOT NULL;