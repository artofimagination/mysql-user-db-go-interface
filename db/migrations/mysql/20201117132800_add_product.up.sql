-- +migrate Up
CREATE TABLE IF NOT EXISTS products(
   id binary(16) PRIMARY KEY,
   name varchar(255) NOT NULL,
   public bool NOT NULL DEFAULT false,
   details json NOT NULL,
   created_at DATETIME NOT NULL DEFAULT NOW(),
   updated_at DATETIME NOT NULL DEFAULT NOW()
);

-- +migrate Up
CREATE TABLE IF NOT EXISTS privileges(
   id tinyint PRIMARY KEY AUTO_INCREMENT,
   name varchar(255) NOT NULL,
   description varchar(512) NOT NULL,
   created_at DATETIME NOT NULL DEFAULT NOW(),
   updated_at DATETIME NOT NULL DEFAULT NOW()
);

-- +migrate Up
INSERT INTO privileges (name, description) VALUES ('Owner', 'The user owns the product, has the right to change or delete it. If there are multiple owners delete/change can happen if all the owners aggree.');
INSERT INTO privileges (name, description) VALUES ('User', 'The user acquired the product. They can create projects based on the product, but cannot modify or delete the product itself.');
INSERT INTO privileges (name, description) VALUES ('Partner', 'The user can change the product, but requires the owners approval.');

-- +migrate Up
CREATE TABLE IF NOT EXISTS users_products(
   products_id binary(16),
   FOREIGN KEY (products_id) REFERENCES products(id),
   users_id binary(16),
   FOREIGN KEY (users_id) REFERENCES users(id),
   privileges_id tinyint,
   created_at DATETIME NOT NULL DEFAULT NOW(),
   updated_at DATETIME NOT NULL DEFAULT NOW()
);
