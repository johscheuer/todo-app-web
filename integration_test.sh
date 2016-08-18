#!/bin/bash

docker-compose up -d

docker-compose ps

curl http://localhost:3000/health

go test -v -tags integration

[ $? -eq 0 ] || exit $?;
# Clean up
docker-compose kill
docker-compose rm -f
