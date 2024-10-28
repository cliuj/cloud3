#!/bin/sh

set -x

docker exec docker-server-1 ls /tmp/cloud3/server
docker exec docker-client1-1 ls /tmp/cloud3/client
docker exec docker-client2-1 ls /tmp/cloud3/client
