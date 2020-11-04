# MYSQL example

This example meant to test basic functionalities of MYSQL. The example allows to create,query and delete users and belonging user settings.
The documents look as follows:

- Documents in users collection
  ```
  users {
      _id -> objectID,
      name -> string,
      email -> string,
      password -> bcrypt blob
      _settings_id -> objectID
  }
- Documents in user_settings collection
  ```
  user_settings {
      _id -> objectID,
      2steps_on -> bool,
  }
# Usage
# Build

- Run ```docker-compose up --build --force-recreate -d main-server``` to generate and start all containers.

- In order to access the db run: ```docker exec -it user-db bash -c "mysql -uroot -p123secure user_database```

- To run bootstrap in your code, the migration files need to be copied manually from the db folder. The destination shall be the root of the golang source.

- Example Dockerfile command: ```RUN git clone https://github.com/artofimagination/mysql-user-db-go-interface $GOPATH/src/user-db-mysql && cp -r $GOPATH/src/user-db-mysql/db $GOPATH/src/load-tester/```

## Execution examples

Use the browser or curl command to execute the following:
- add new user that automatically generates a belonging settings document: ```http://localhost:8080/insert?name=testName&email=testEmail&password=testPass```
- get user with specified email: ```http://localhost:8080/get?email=testEmail```
- delete user specified by the email: ```http://localhost:8080/delete?email=testEmail```
- check user password and email: ```http://localhost:8080/check?email=testEmail&password=testPass```
- get user settings belonging to the user with specified email: ```http://localhost:8080/get-settings?email=testEmail```
- delete user settings belonging to the user with specified email, this test shall fail since the user has the settigns foreign key: ```http://localhost:8080/delete-settings?email=testEmail```
