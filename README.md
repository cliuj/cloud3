# cloud3

## About
cloud3 is a **proof-of-concept** program that allows users to sync files between remote systems via a shared directory.
Files placed inside this shared directory is automatically sync'd with other clients connected to the server.

## Demo

#### Setup
The following `docker compose` command deploys 3 docker containers (1 server, 2 clients):
```bash
# From root of the repository
docker compose -f docker/docker-compose.yml up -d --build
```
```
[+] Running 3/3
 ✔ Container docker-client2-1  Started                                                                                         0.1s
 ✔ Container docker-server-1   Started                                                                                         0.1s
 ✔ Container docker-client1-1  Started                                                                                         0.1s
```

#### Monitoring
Users can monitor the changes of the shared directories via the provided `scripts/peak-docker.sh`:
```bash
chmod u+x ./scripts/peak-docker.sh
watch ./scripts/peak-docker.sh
```

```
Every 2.0s: ./peak-docker.sh  tp-p14s: Sun Oct 27 23:40:37 2024

+ docker exec docker-server-1 ls /tmp/cloud3/server
+ docker exec docker-client1-1 ls /tmp/cloud3/client
+ docker exec docker-client2-1 ls /tmp/cloud3/client
```

#### Usage
Exec into any client container and add files to the `/tmp/cloud3/client/` directory
to the added files reflected on the other client.

###### Example:
```bash
docker exec -it docker-client1-1 bash
cd /tmp/cloud3/client/
echo "Hello world!" > test.txt
```
