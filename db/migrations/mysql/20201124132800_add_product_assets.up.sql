-- +migrate Up
CREATE TABLE IF NOT EXISTS product_assets(
   id binary(16) PRIMARY KEY,
   refs json NOT NULL,
   created_at DATETIME NOT NULL DEFAULT NOW(),
   updated_at DATETIME NOT NULL DEFAULT NOW()
);

-- +migrate Up
ALTER TABLE products 
ADD COLUMN product_assets_id binary(16) AFTER public, 
ADD FOREIGN KEY (product_assets_id) REFERENCES product_assets(id);

-- +migrate Up
ALTER TABLE users ADD UNIQUE (name);
