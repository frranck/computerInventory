# Dockerfile
FROM golang:1.24-alpine

RUN apk add --no-cache gcc musl-dev

WORKDIR /app

# Copy go.mod and go.sum first for caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the Go app
RUN CGO_ENABLED=1 go build -o computerInventory ./cmd/main.go

# Set the binary as the entrypoint
CMD ["./computerInventory"]
