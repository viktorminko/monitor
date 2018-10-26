#!/bin/bash
#Current dir
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

echo "Load compilation configuration environment variables"
source <(sed -E -n 's/[^#]+/export &/ p' $DIR/compile.env)

DOCKER_WORK_DIR="/usr/local/go/src/github.com/viktorminko/monitor"

echo "Compile application in docker golang container"
docker run \
-v $SRC_DIR/cmd:$DOCKER_WORK_DIR/cmd \
-v $SRC_DIR/pkg:$DOCKER_WORK_DIR/pkg \
-v $SRC_DIR/bin:$DOCKER_WORK_DIR/bin \
-w $DOCKER_WORK_DIR/cmd \
golang \
/bin/bash -c "go get -v && CGO_ENABLED=0 go build -v -o $DOCKER_WORK_DIR/bin/monitor"