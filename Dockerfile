# Start from a base image with Go and Redis installed
FROM golang:latest

# Set the working directory inside the container
WORKDIR /app

# Copy go modules and download dependencies
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go application
RUN go build -o main .

# Install Redis
RUN apt-get update && apt-get install -y redis-server

# Expose the port on which your Go application will run
EXPOSE 8080 8081 8000 6379

# Set environment variables (modify as needed)
# ENV INFURA_API_KEY="YOUR_INFURA_API_KEY"
# ENV REDIS_ADDR="redis:6379"
# ENV REDIS_PASSWORD=""

# Run the Go application and redis-server
CMD ["sh", "-c", "service redis-server start && ./main"]
