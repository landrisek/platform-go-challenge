FROM golang:1.18 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the source code to the working directory
COPY . .

# Download Go module dependencies
RUN go mod download

# Build the Go application
RUN go build -o migrations cmd/migrations/main.go

# Final stage
FROM alpine:latest

# Set the working directory inside the container
WORKDIR /app

# Install curl
RUN apk add --no-cache curl

# Copy the built Go binary from the builder stage
COPY --from=builder /app/migrations /app/

# Copy the provisioning script
COPY scripts/provision-vault.sh /app/provision-vault.sh

# Make the script executable
RUN chmod +x /app/migrations
RUN chmod +x /app/provision-vault.sh

# Set the entrypoint command to run the provisioning script and migrations
ENTRYPOINT ["sh", "/app/provision-vault.sh"]
#CMD ["./migrations"]
