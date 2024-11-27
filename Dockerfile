# Use the official Golang image as a build stage
FROM golang:1.23 AS builder

# Set the working directory inside the container
WORKDIR /consumer

# Copy the Go module files into the container
COPY go.mod go.sum ./

# Download the module dependencies
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go consumerlication with CGO enabled for Kafka support
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -v -o main .

# Use an Ubuntu-based image for the final stage to ensure compatibility with CGO
FROM ubuntu:22.04

# Install necessary certificates and dependencies for CGO (like libc)
RUN apt-get update && apt-get install -y ca-certificates libc6

# Set the working directory in the final image
WORKDIR /root/

# Copy the compiled Go binary from the builder stage
COPY --from=builder /consumer/main .

# Ensure the binary is executable
RUN chmod +x /root/main

# Expose the port your consumer runs on
EXPOSE 9090

# Command to run the executable
CMD ["./main"]
