# Base image for building the Go application
FROM 1.24.5-alpine3.22 as build

# Set the working directory for the build process
WORKDIR /build

# Copy the entire source code into the container
COPY . .

# Download and tidy up Go module dependencies
RUN go mod tidy

# Build the Go application binary
RUN go build -o app .

# Base image for running the compiled application
FROM alpine:3.22.1

# Copy the built application binary from the build stage
COPY --from=build /build/app .

# Expose the application port to allow external access
EXPOSE 8080

# Create a directory for SQLite database files
RUN mkdir db

# Define a volume for persistent database storage
VOLUME /db

# Set the entry point to run the compiled application binary
ENTRYPOINT ["./app"]
