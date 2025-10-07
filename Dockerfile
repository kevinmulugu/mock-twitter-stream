FROM golang:1.23-alpine AS build
RUN apk add --no-cache git ca-certificates

WORKDIR /app
COPY go.* ./
RUN go mod download
COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
  -ldflags='-w -s -extldflags "-static"' \
  -a \
  -o mock-twitter-stream

FROM scratch
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs
COPY --from=build /app/mock-twitter-stream /app/mock-twitter-stream
USER 65534:65534
EXPOSE 8080
CMD ["/app/mock-twitter-stream", "-addr", ":8080"]