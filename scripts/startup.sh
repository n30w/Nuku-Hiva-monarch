#!/bin/bash

docker run --env-file=andthensome_env_vars --publish 4000:4000 -d ghcr.io/n30w/andthensome:master
docker run -d \
-v /var/run/docker.sock:/var/run/docker.sock \
containrrr/watchtower
