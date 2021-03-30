#pip3 install -r test/requirements.txt
docker-compose down
docker stop $(docker ps -aq)
docker rm $(docker ps -aq)
docker system prune -f
docker-compose up --build --force-recreate -d main-server
status=$?; 
if [[ $status != 0 ]]; then 
  exit $status; 
fi
python3 -m pytest -v tests