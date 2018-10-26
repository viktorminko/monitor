#!/bin/bash
#current dir
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

$DIR/compile.sh

echo "Load run configuration environment variables"
source <(sed -E -n 's/[^#]+/export &/ p' $DIR/run.env)

#run monitor in docker container
docker build -t $CONTAINER_NAME . && docker run -v $DIR/../bin/:/home/bin -v $CONFIG_DIR:/home/config $CONTAINER_NAME sh -c \
"exec /home/bin/monitor -workdir /home/config/"