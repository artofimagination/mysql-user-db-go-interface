-- +migrate Up
DROP TABLE projects;

-- +migrate Up
CREATE TABLE IF NOT EXISTS project_details(
   id binary(16) PRIMARY KEY,
   data json NOT NULL,
   created_at DATETIME NOT NULL DEFAULT NOW(),
   updated_at DATETIME NOT NULL DEFAULT NOW()
);

-- +migrate Up
CREATE TABLE IF NOT EXISTS project_assets(
   id binary(16) PRIMARY KEY,
   data json NOT NULL,
   created_at DATETIME NOT NULL DEFAULT NOW(),
   updated_at DATETIME NOT NULL DEFAULT NOW()
);

-- +migrate Up
CREATE TABLE IF NOT EXISTS projects(
   id binary(16) NOT NULL PRIMARY KEY,
   products_id binary(16) NOT NULL,
   FOREIGN KEY (products_id) REFERENCES products(id),
   project_details_id binary(16),
   FOREIGN KEY (project_details_id) REFERENCES project_details(id),
   project_assets_id binary(16),
   FOREIGN KEY (project_assets_id) REFERENCES project_assets(id),
   created_at DATETIME NOT NULL DEFAULT NOW(),
   updated_at DATETIME NOT NULL DEFAULT NOW()
);

-- +migrate Up
CREATE TABLE IF NOT EXISTS users_projects(
   projects_id binary(16),
   FOREIGN KEY (projects_id) REFERENCES projects(id),
   users_id binary(16),
   FOREIGN KEY (users_id) REFERENCES users(id),
   privileges_id tinyint,
   created_at DATETIME NOT NULL DEFAULT NOW(),
   updated_at DATETIME NOT NULL DEFAULT NOW()
);

-- +migrate Up
CREATE TABLE IF NOT EXISTS users_viewers(
   users_id binary(16),
   FOREIGN KEY (users_id) REFERENCES users(id),
   viewer_id bigint,
   projects_id binary(16),
   FOREIGN KEY (projects_id) REFERENCES projects(id),
   is_owner bool,
   created_at DATETIME NOT NULL DEFAULT NOW(),
   updated_at DATETIME NOT NULL DEFAULT NOW()
);