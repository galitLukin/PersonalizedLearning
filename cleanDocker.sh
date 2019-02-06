if ! test -z "$(docker ps -q)"; then
  docker kill $(docker ps -q)
fi
docker rm $(docker ps -a -q)
for id in $(docker images -q);
do
  if [ "$id" != "$c" ]; then
    docker rmi $id -f;
  fi
done
docker build . -t  arieg419/quiz-app:latest