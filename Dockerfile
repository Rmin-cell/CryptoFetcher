# Use the official Golang image as the base image
FROM golang:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the entire project directory into the container
COPY . .

# Build the Go application inside the container
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o crypto-tracker .

# Command to run the application
CMD ["./main.go"]
