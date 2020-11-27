-- +migrate Up
ALTER TABLE user_settings RENAME COLUMN settings TO data;
ALTER TABLE product_details RENAME COLUMN details TO data;
ALTER TABLE user_assets RENAME COLUMN refs TO data;
ALTER TABLE product_assets RENAME COLUMN refs TO data;
