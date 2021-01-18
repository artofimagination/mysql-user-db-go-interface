# Example implementation for user/product/project and high volume data store using MYSQL and Golang

NOTE: This implementation is heavily under development, but is in a useful state. Feel free to try, raise tickets and fix bugs.

This example implementation provides a fundamental handling of a user database. The users can have products, projects and project associated data viewers.
All content of each object type (user/product/project) is fully customizable using two json columns.

The database can be used as a store for webshops, platforms that allow users to create projects or applications that require high volume data collection.

- Users table
  ```
  +------------------+---------------+------+-----+-------------------+-------------------+
  | Field            | Type          | Null | Key | Default           | Extra             |
  +------------------+---------------+------+-----+-------------------+-------------------+
  | id               | binary(16)    | NO   | PRI | NULL              |                   |
  | name             | varchar(50)   | NO   | UNI | NULL              |                   |
  | email            | varchar(300)  | NO   | UNI | NULL              |                   |
  | password         | varchar(1024) | YES  |     | NULL              |                   |
  | user_settings_id | binary(16)    | YES  | MUL | NULL              |                   |
  | user_assets_id   | binary(16)    | YES  | MUL | NULL              |                   |
  | created_at       | datetime      | NO   |     | CURRENT_TIMESTAMP | DEFAULT_GENERATED |
  | updated_at       | datetime      | NO   |     | CURRENT_TIMESTAMP | DEFAULT_GENERATED |
  +------------------+---------------+------+-----+-------------------+-------------------+
- Products table
  ```
  +--------------------+--------------+------+-----+-------------------+-------------------+
  | Field              | Type         | Null | Key | Default           | Extra             |
  +--------------------+--------------+------+-----+-------------------+-------------------+
  | id                 | binary(16)   | NO   | PRI | NULL              |                   |
  | name               | varchar(255) | NO   | UNI | NULL              |                   |
  | product_details_id | binary(16)   | YES  | MUL | NULL              |                   |
  | product_assets_id  | binary(16)   | YES  | MUL | NULL              |                   |
  | created_at         | datetime     | NO   |     | CURRENT_TIMESTAMP | DEFAULT_GENERATED |
  | updated_at         | datetime     | NO   |     | CURRENT_TIMESTAMP | DEFAULT_GENERATED |
  +--------------------+--------------+------+-----+-------------------+-------------------+
- Projects table
  ```
  +--------------------+------------+------+-----+-------------------+-------------------+
  | Field              | Type       | Null | Key | Default           | Extra             |
  +--------------------+------------+------+-----+-------------------+-------------------+
  | id                 | binary(16) | NO   | PRI | NULL              |                   |
  | products_id        | binary(16) | NO   | MUL | NULL              |                   |
  | project_details_id | binary(16) | YES  | MUL | NULL              |                   |
  | project_assets_id  | binary(16) | YES  | MUL | NULL              |                   |
  | created_at         | datetime   | NO   |     | CURRENT_TIMESTAMP | DEFAULT_GENERATED |
  | updated_at         | datetime   | NO   |     | CURRENT_TIMESTAMP | DEFAULT_GENERATED |
  +--------------------+------------+------+-----+-------------------+-------------------+
  
# Usage
It is recommended to call only the dbcontrollers functions using it in third party code.

## Build
- Run ```docker-compose up --build --force-recreate -d main-server``` to generate and start all containers.
- In order to access the db run: ```docker exec -it user-db bash -c "mysql -uroot -p123secure user_database```
- To run bootstrap/migration in your code, the migration files need to be copied manually from the db folder. The destination shall be the root of the golang source.
  Example docker file code:
  ```
  RUN git clone https://github.com/artofimagination/mysql-user-db-go-interface /tmp/mysql-user-db-go-interface && \
  cp -r /tmp/mysql-user-db-go-interface/db $GOPATH/src/my-app && \
  rm -fr /tmp/mysql-user-db-go-interface
- .env.example contains an example docker config that is required to run the code as intended. Rename it to .env and customize as needed.

## Running the example code
To run functional testing using the example code run ```./runFunctionalTest.sh```

Once the example main-server is running the user can do the following using the curl command:
- add new user: ```http://localhost:8080/insert?name=testName&email=testEmail&password=testPass```
- get user with specified email: ```http://localhost:8080/get?email=testEmail```
- delete user specified by the email: ```http://localhost:8080/delete?email=testEmail```
- check user password and email: ```http://localhost:8080/check?email=testEmail&password=testPass```
- get user settings belonging to the user with specified email: ```http://localhost:8080/get-settings?email=testEmail```
- delete user settings belonging to the user with specified email, this test shall fail since the user has the settigns foreign key: ```http://localhost:8080/delete-settings?email=testEmail```

# Database
## Entity relation
[Entity relation](docs/DBRelations.jpg)
## UML
[UML](docs/UML.jpg)
## dbcontrollers use cases
[Use cases](docs/UseCase.jpg)


