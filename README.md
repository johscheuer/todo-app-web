[![Build Status](https://travis-ci.org/johscheuer/todo-app-web.svg?branch=master)](https://travis-ci.org/johscheuer/todo-app-web)

# Example web todo list

This is a simple web todo list for presentations and demos. It's written in Golang + Javascript.

## Building the image

### Compile the code

```bash
go get -u github.com/johscheuer/todo-app-web
```

#### On OSX

```bash
CGO_ENABLED=0 GOOS=linux go build -ldflags "-s -w -X main.appVersion=$(git symbolic-ref -q --short HEAD || git describe --tags --exact-match)" -a -installsuffix cgo -o bin/todo-app .
```

### Build the Container

```bash
$ docker build -t johscheuer/todo-app-web .
# Tag the image if you want
docker tag -f johscheuer/todo-app-web johscheuer/todo-app-web:<tag>
docker push johscheuer/todo-app-web
```

## Testing

```bash
./integration_test.sh
```

## Usage

```bash
Usage of bin/todo-app:
  -health-check int
           Period to check all connections (default 15)
  -master string
           The connection string to the Redis master as <hostname/ip>:<port> (default "redis-master:6379")
  -master-password string
           The password used to connect to the master
  -slave string
           The connection string to the Redis slave as <hostname/ip>:<port> (default "redis-slave:6379")
  -slave-password string
           The password used to connect to the slave
  -version
           Shows the version
```
