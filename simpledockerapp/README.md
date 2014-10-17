Run app using container from [docker hub](https://registry.hub.docker.com/_/golang/)
===

### Build docker image
docker build -t simpledockerapp .

### Run container
docker run -it --rm --name simpledockerappinstance simpledockerapp

### Run this in a different tmux panel to inspect the container and run bash shell within the container
container_id=$(docker ps -lq)
docker inspect $container_id
docker exec -it $container_id bash

### AppG
- Runs a bunch of workers (Can use -workers to configure) which echo with a pause interval (Can use -sleep to configure) until we quit the app 
- Uses a channel to tell workers to quit
- Waits on a sigint to quit the app
- Uses a waitgroup to allow all workers complete

