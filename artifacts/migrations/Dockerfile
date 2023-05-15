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

# Make the migrations binary executable
RUN chmod +x migrations

# Set the entrypoint command to run the migrations
#ENTRYPOINT ["./migrations"]
