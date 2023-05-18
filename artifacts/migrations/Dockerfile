FROM golang:1.18

# Set the working directory inside the container
WORKDIR /app

# Copy the source code to the working directory
COPY . .

# Download Go module dependencies
RUN go mod download

# Build the Go application
RUN go build -o migrations cmd/migrations/main.go

# Install curl
RUN apt-get update && apt-get install -y curl

# Install Redis tools
RUN apt-get update && apt-get install -y redis-tools

# Install the MariaDB client
RUN apt-get update && apt-get install -y mariadb-client

# Make the migrations binary executable
RUN chmod +x migrations