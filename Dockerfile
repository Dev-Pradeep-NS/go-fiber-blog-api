# Builder Stage
FROM golang:1.23-alpine AS builder

# Set the working directory
WORKDIR /app

# Install Air for live reloading
RUN go install github.com/air-verse/air@latest

# Copy go.mod and go.sum files
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of your application code
COPY . .

# Build the application
RUN go build -o main ./cmd/server

# Final Stage
FROM golang:1.23-alpine

# Set the working directory
WORKDIR /app

# Copy the Go environment from the builder stage
COPY --from=builder /go /go

# Copy the Air binary from the builder stage
COPY --from=builder /go/bin/air /usr/local/bin/air

# Copy the built application from the builder stage
COPY --from=builder /app/main /app/main

# Copy the Air configuration file
COPY .air.toml .

# Run Air with the specified configuration file
CMD ["air", "-c", ".air.toml"]