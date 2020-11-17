-- +migrate Up
CREATE TABLE IF NOT EXISTS features(
   id serial PRIMARY KEY,
   name varchar(128) UNIQUE NOT NULL,
   config json,
   created_at DATETIME NOT NULL DEFAULT NOW(),
   updated_at DATETIME NOT NULL DEFAULT NOW()
);

INSERT INTO features (name, config) values ('Survival', '{"config":{"Food Count": "500", "Creature count": "50"}}');
INSERT INTO features (name, config) values ('Image Processing', '{"config":{"Food Count": "500", "Creature count": "150"}}');