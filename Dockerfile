# Use a minimal base image
FROM docker.arvancloud.ir/golang:1.21-alpine

# Set environment variables
ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Set working directory
WORKDIR /app

# Install Git (needed for go mod) and ca-certificates
RUN apk add --no-cache git ca-certificates && update-ca-certificates

# Copy go.mod and go.sum
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build the Go binary
RUN go build -o app .

# Run the compiled binary
CMD ["./app"]
