#!/bin/bash
function clean_up {
    docker-compose kill && docker-compose rm -f
}

clean_up
docker-compose up -d

go test -v -tags integration
RET=$?
clean_up

exit ${RET}
