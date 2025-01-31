# Use the official Golang image as the builder
FROM golang:1.21-alpine AS builder

# Install build dependencies (gcc, musl-dev, etc.)
RUN apk add --no-cache gcc musl-dev

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY main.go .

# Copy the SQLite database file
COPY data.db .

# Build the Go application with CGO enabled
RUN CGO_ENABLED=1 GOOS=linux go build -o main .

# Use a minimal Alpine image for the final stage
FROM alpine:latest

# Install SQLite (required to use the SQLite database)
RUN apk --no-cache add sqlite

# Set the working directory
WORKDIR /app

# Copy the compiled binary from the builder stage
COPY --from=builder /app/main .
# Copy the SQLite database file
COPY --from=builder /app/data.db .

# Expose the port your app runs on (replace 8080 with your app's port if different)
EXPOSE 8080

# Run the application
CMD ["./main"]