-- +migrate Up
ALTER TABLE users_viewers MODIFY viewer_id binary(16);

-- +migrate Up
ALTER TABLE users_viewers DROP COLUMN is_owner;

-- +migrate Up
CREATE TABLE IF NOT EXISTS viewers(
   id binary(16) PRIMARY KEY,
   owner_id binary(16),
   FOREIGN KEY (owner_id) REFERENCES users(id),
   created_at DATETIME NOT NULL DEFAULT NOW(),
   updated_at DATETIME NOT NULL DEFAULT NOW()
);

ALTER TABLE `users_viewers` ADD CONSTRAINT fk_viewer_id FOREIGN KEY (viewer_id) REFERENCES viewers(id);
