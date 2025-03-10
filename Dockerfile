# Use a newer Go runtime as the base image
FROM golang:1.23.2

# Set the working directory inside the container
WORKDIR /app

# Copy the Go application source code into the container
COPY . .

#   Download the Go modules
RUN go mod tidy

# Build the Go application
RUN go build -o binary .

# Expose the port your application will run on
EXPOSE 8080

# Command to run the application
CMD ["./binary"]