FROM golang:1.12.5-stretch as Builder
WORKDIR /go/src/github.com/johscheuer/todo-app-web/
COPY ${HOME}/ /go/src/github.com/johscheuer/todo-app-web/
RUN GO111MODULE=on CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app/todo-app .

COPY public /app/public

WORKDIR /app
CMD ["./todo-app"]
EXPOSE 3000
