# Use a small base image
FROM docker.arvancloud.ir/golang:1.24-alpine AS builder

ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /app

# Install git and certs for Go module resolution
RUN apk add --no-cache git ca-certificates && update-ca-certificates

# Copy go.mod and go.sum and download deps
COPY go.mod go.sum ./
RUN go mod download

# Copy the whole project
COPY . .

# Build the binary from cmd/main.go
RUN go build -o app ./cmd/main.go

# Use a tiny final image
FROM docker.arvancloud.ir/alpine:latest
WORKDIR /root/

COPY --from=builder /app/app .

EXPOSE 8123

CMD ["./app"]
