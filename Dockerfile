FROM golang:1.12.5-stretch as Builder
WORKDIR /go/src/github.com/johscheuer/todo-app-web/
COPY ./ /go/src/github.com/johscheuer/todo-app-web/
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o todo-app .

FROM gcr.io/distroless/base
COPY --from=builder /go/src/github.com/johscheuer/todo-app-web/todo-app /app/todo-app
COPY ./public /app/public

WORKDIR /app
CMD ["./todo-app"]
EXPOSE 3000
