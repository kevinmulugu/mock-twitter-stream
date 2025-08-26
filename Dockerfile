FROM golang:1.23-alpine AS build
WORKDIR /app
COPY . .
RUN go build -o mock-twitter-stream

FROM alpine:latest
WORKDIR /app
COPY --from=build /app/mock-twitter-stream /usr/local/bin/mock-twitter-stream
EXPOSE 8080
CMD ["/usr/local/bin/mock-twitter-stream", "-addr", ":8080"]