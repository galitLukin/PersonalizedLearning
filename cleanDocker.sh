if ! test -z "$(docker ps -q)"; then
  docker kill $(docker ps -q)
fi
docker rm $(docker ps -a -q)
for id in $(docker images -q);
do
  docker rmi $id -f;
done
docker build . -t  lukingalit/personalized-learning:latest
docker run -d -p 80:80 lukingalit/personalized-learning:latest
