#!/bin/bash
ls -lah ./bin

docker-compose up -d

docker-compose ps

docker-compose logs todoApp

go test -v -tags integration

[ $? -eq 0 ] || exit $?;
# Clean up
docker-compose kill
docker-compose rm -f
