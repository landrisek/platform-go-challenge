FROM golang:1.18 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the source code to the working directory
COPY . .

# Download Go module dependencies
RUN go mod download

# Build the Go application
RUN go build -o user cmd/user/main.go

# Start a new stage to create a lightweight image
FROM golang:1.18

# Set the working directory inside the container
WORKDIR /app

# Copy the built user executable from the previous stage
COPY --from=builder /app/user .

# Expose the port that the user listens on
EXPOSE 8082

# Set the entry point command to start the user
CMD ["./user"]