-- +migrate Up
CREATE TABLE IF NOT EXISTS projects(
   id binary(16) PRIMARY KEY,
   user_id binary(16) REFERENCES users(id),
   features_id integer REFERENCES features(id),
   name varchar(128) UNIQUE NOT NULL,
   config json,
   created_at DATETIME NOT NULL DEFAULT NOW(),
   updated_at DATETIME NOT NULL DEFAULT NOW()
);