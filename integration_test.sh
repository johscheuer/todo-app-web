#!/bin/bash
docker-compose up -d

go test -v -tags integration

[ $? -eq 0 ] || exit $?;
# Clean up
docker-compose kill
docker-compose rm -f
