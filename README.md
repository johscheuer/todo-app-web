# Example web todo list

This is a simple web todo list for presentations and demos. It's written in Golang + Javascript.

## Building the image

### Compile the code

```bash
go get -u github.com/johscheuer/todo-app-web
```

#### On OSX

```bash
CGO_ENABLED=0 GOOS=linux go build -ldflags "-s -w" -a -installsuffix cgo -o bin/todo-app .
```

### Build the Container

```bash
docker build -t johscheuer/todo-app-web .
# Tag the image if you want
docker tag -f johscheuer/todo-app-web johscheuer/todo-app-web:v1
docker push johscheuer/todo-app-web
```

## Testing

```bash
docker-compose up -d
```
