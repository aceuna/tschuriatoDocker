# Step 1: Build the Go application
FROM golang:1.21-bullseye as builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download the dependencies
RUN go mod tidy

# Copy the source code into the container
COPY . .

# Build the Go app
RUN go build -o stock-publisher .

# Step 2: Create the final image using a smaller runtime image
FROM debian:bullseye-slim  

# Set the Current Working Directory inside the container
WORKDIR /root/

# Install necessary dependencies for running the app
RUN apt-get update && apt-get install -y ca-certificates
# Install netcat (nc) to check for open ports
RUN apt-get update && apt-get install -y netcat

# Copy the pre-built binary from the builder image
COPY --from=builder /app/stock-publisher .

# Expose the port the app will run on
EXPOSE 8080

# Command to run the binary
#CMD ["./stock-publisher"]
