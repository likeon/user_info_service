# build
FROM docker.io/golang:1.23-alpine AS builder

# needed for go-sqlite3
RUN apk add --no-cache gcc musl-dev

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# CGO required for go-sqlite3
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o /app/bin/server .

# stage
FROM docker.io/alpine:latest

COPY --from=builder /app/bin/server /usr/local/bin/server

EXPOSE 8080

ENTRYPOINT ["/usr/local/bin/server"]
