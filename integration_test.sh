#!/bin/bash

docker-compose up -d

go test -v -tags integration

# Clean up
docker-compose kill
docker-compose rm -f
