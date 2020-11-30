-- +migrate Up
CREATE TABLE IF NOT EXISTS product_details(
   id binary(16) PRIMARY KEY,
   details json NOT NULL,
   created_at DATETIME NOT NULL DEFAULT NOW(),
   updated_at DATETIME NOT NULL DEFAULT NOW()
);

ALTER TABLE products DROP COLUMN details;

-- +migrate Up
ALTER TABLE products 
ADD COLUMN product_details_id binary(16) AFTER public, 
ADD FOREIGN KEY (product_details_id) REFERENCES product_details(id);

-- +migrate Up
ALTER TABLE products ADD UNIQUE (name);
ALTER TABLE privileges ADD UNIQUE (name);
