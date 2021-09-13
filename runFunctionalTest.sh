pip3 install -r tests/requirements.txt

cp ./tests/.env.functional_test .env

docker-compose down
docker stop $(docker ps -aq)
docker rm $(docker ps -aq)
docker system prune -f

docker-compose up --build --force-recreate -d user-db-server
status=$?; 
if [[ $status != 0 ]]; then 
  exit $status; 
fi
status=$?; 
if [[ $status != 0 ]]; then 
  exit $status; 
fi
python3 -m pytest -v tests