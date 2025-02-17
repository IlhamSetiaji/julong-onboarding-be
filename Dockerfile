# Use the official Golang image as the build stage
FROM golang:1.23 AS builder

# Set environment variables to ensure `go install` puts `swag` in a known location
ENV GOPATH=/go
ENV GOBIN=$GOPATH/bin
ENV PATH=$GOBIN:$PATH

# Set the working directory inside the container
WORKDIR /app

# Install Swag CLI
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Verify that Swag CLI is installed
RUN swag --version

# Copy the go.mod and go.sum files
COPY go.mod go.sum ./

# Download the Go module dependencies
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the Go application
RUN go build -ldflags "-s -w" -o main .

# Use the official Golang image as the base image for the final stage
FROM golang:1.23

# Install necessary runtime dependencies
RUN apt-get update && apt-get install -y gettext-base && rm -rf /var/lib/apt/lists/*

# Set the working directory inside the container
WORKDIR /app

# Copy the built Go application from the builder stage
COPY --from=builder /app/main .
COPY config.template.json /app/config.template.json
COPY init-config.sh /app/init-config.sh

# Copy the views directory
# COPY views /app/views

# Copy the static directory
# COPY public /app/public

# Make the initialization script executable
RUN chmod +x /app/init-config.sh

# Expose the port on which the application will run
EXPOSE 8000

# Command to run the initialization script
CMD ["/app/init-config.sh"]