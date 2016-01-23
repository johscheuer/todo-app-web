# Example web todo list
This is a simple web todo list for presentations and demos. It's written in Golang + Javascript.

## Building the image
### Compile the code

```Bash
go get github.com/codegangsta/negroni
go get github.com/gorilla/mux
go get github.com/xyproto/simpleredis
```

#### On OSX

```Bash
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/todo-app .
```

### Build the Container

```Bash
docker build -t johscheuer/todo-app-web .
# Tag the image if you want
docker tag -f johscheuer/todo-app-web johscheuer/todo-app-web:v1
docker push johscheuer/todo-app-web
```
