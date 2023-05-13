FROM golang:1.18

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files
COPY ../go.mod ../go.sum ./

# Download Go module dependencies
RUN go mod download

# Copy the source code to the working directory
COPY . .

# Build the Go application
RUN go build -o migrations cmd/migrations/main.go

# Set the entry point command to run the application
CMD ["./migrations"]