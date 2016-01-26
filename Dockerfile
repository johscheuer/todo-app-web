FROM alpine:3.3
MAINTAINER Johannes M. Scheuermann <joh.scheuer@gmail.com>

COPY ./bin/todo-app /app/todo-app
COPY ./public /app/public

WORKDIR /app
CMD ["./todo-app"]
EXPOSE 3000
