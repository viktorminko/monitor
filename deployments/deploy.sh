#!/bin/bash
#current directory
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

echo "Load deployments configuration environment variables"
source <(sed -E -n 's/[^#]+/export &/ p' $DIR/.env)

#compile monitor
echo "---- Compile monitor"
$DIR/../scripts/compile.sh

echo "---- Create deployments package"
mkdir -p $DIR/package && mkdir -p $DIR/package/config && mkdir -p $DIR/package/bin
cp $DIR/../bin/monitor $DIR/package/bin/
cp -r $CONFIG_PATH/* $DIR/package/config
tar -C $DIR/package/ -czvf package.tar.gz ./

echo "---- Delete package folder"
rm -rf $DIR/package

echo "---- Set up docker-machine environment to deployments droplet"
eval $(docker-machine env $DOCKER_MACHINE)

echo "---- Stop monitor container"
docker-compose --project-name $PROJECT_NAME stop

echo "---- Create monitor directory on remote machine"
docker-machine ssh $DOCKER_MACHINE rm -rf $SERVER_PATH
docker-machine ssh $DOCKER_MACHINE mkdir -p $SERVER_PATH

echo "---- Copy package archive to remote machine"
docker-machine scp -r $DIR/package.tar.gz $DOCKER_MACHINE:$SERVER_PATH

echo "---- Extract package archive"
docker-machine ssh $DOCKER_MACHINE tar -xf $SERVER_PATH/package.tar.gz -C $SERVER_PATH

echo "---- Delete package archive on remote machine"
docker-machine ssh $DOCKER_MACHINE rm -f $SERVER_PATH/package.tar.gz

echo "---- Delete package archive on localhost"
rm -f $DIR/package.tar.gz

docker-machine ssh $DOCKER_MACHINE cd $SERVER_PATH

echo "---- Build monitor container"
docker-compose --project-name $PROJECT_NAME build

echo "---- Run monitor container"
docker-compose --project-name $PROJECT_NAME up -d


