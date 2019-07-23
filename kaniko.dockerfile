FROM gcr.io/distroless/base
COPY /cache/todo-app /app/todo-app
COPY public /app/public
WORKDIR /app
CMD ["./todo-app"]
EXPOSE 3000
