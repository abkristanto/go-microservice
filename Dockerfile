FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install git if needed for modules
RUN apk add --no-cache git

# Cache go.mod / go.sum first
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the code
COPY . .

# Build the binary
RUN go build -o service ./cmd/service

# Final image
FROM alpine:3.20

WORKDIR /app

COPY --from=builder /app/service .

CMD ["./service"]
