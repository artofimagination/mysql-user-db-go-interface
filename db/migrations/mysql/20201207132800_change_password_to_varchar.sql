-- +migrate Up
ALTER TABLE users MODIFY password VARCHAR(1024);
