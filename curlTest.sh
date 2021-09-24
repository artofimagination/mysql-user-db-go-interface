curl 0.0.0.0:8181/index
curl --header "Content-Type: application/json"   --request POST   --data '{"email": "testEmailOwnerGetProject2", "password": "dGVzdFBhc3N3b3Jk", "username": "testUserOwnerGetProject2"}'   http://0.0.0.0:8181/add-user
curl --header "Content-Type: application/json"   --request POST   --data '{"password": "dGVzdFBhc3N3b3Jk", "username": "testUserOwnerGetProject2"}'   http://0.0.0.0:8181/add-user
pip3 install -r tests/requirements.txt && pytest -v tests
docker ps && docker logs user-db-server