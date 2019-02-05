if ! test -z "$(docker ps -q)"; then
  docker kill $(docker ps -q)
fi
docker rm $(docker ps -a -q)
for id in $(docker images -q)
do
  if "$id" -ne "7e70d62cd078"; then
    docker rmi $id
  fi
done
docker build .
