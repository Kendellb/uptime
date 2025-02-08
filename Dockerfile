# Step 1: Build the Go application
FROM golang:1.20 AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod  ./

# Download dependencies (if any)
RUN go mod tidy

# Copy the source code into the container
COPY . .

# Build the Go app
RUN go build -o main .

# Step 2: Run the Go application
FROM debian:bookworm-slim

# Install necessary dependencies (such as certificates) for HTTPS support
RUN apt-get update && apt-get install -y ca-certificates procps

# Set the Current Working Directory inside the container
WORKDIR /root/

# Copy the compiled binary from the builder container
COPY --from=builder /app/main .

# Copy the static directory to serve static files
COPY --from=builder /app/static /root/static

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./main"]

