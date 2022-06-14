# Start by building the application.
FROM golang:1.18-bullseye as build

# Now copy it into our base image.
FROM gcr.io/distroless/base-debian11
COPY ./todo.db todo.db
COPY ./go-rest-api /app
CMD ["/app"]