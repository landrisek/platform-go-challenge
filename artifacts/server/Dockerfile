FROM golang:1.18 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files
COPY ../go.mod ../go.sum ./

# Download Go module dependencies
RUN go mod download

# Copy the source code to the working directory
COPY . .

# Build the Go application
RUN go build -o server cmd/server/main.go

# Start a new stage to create a lightweight image
FROM golang:1.18

# Set the working directory inside the container
WORKDIR /app

# Copy the built server executable from the previous stage
COPY --from=builder /app/server .

# Expose the port that the server listens on
EXPOSE 8080

# Set the entry point command to start the server
CMD ["./server"]