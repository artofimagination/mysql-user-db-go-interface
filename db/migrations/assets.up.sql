-- +migrate Up
CREATE TABLE IF NOT EXISTS user_assets(
   id binary(16) PRIMARY KEY,
   refs json NOT NULL,
   created_at DATETIME NOT NULL DEFAULT NOW(),
   updated_at DATETIME NOT NULL DEFAULT NOW()
);

-- +migrate Up
ALTER TABLE users 
ADD COLUMN user_assets_id binary(16) AFTER user_settings_id, 
ADD FOREIGN KEY (user_assets_id)  REFERENCES user_assets(id);
