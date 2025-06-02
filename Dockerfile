# Use Go 1.24.3 bookworm as base image
FROM golang:1.24.3-bookworm AS base

# Move to working directory /build
WORKDIR /build

# Copy the entire source code into the container
COPY . .

# Build microservice_go_postgres
RUN go build -o transaction-crud-svc-go-postgres .

# Document the port that may need to be published
EXPOSE 3000

# Start the application
CMD ["/build/transaction-crud-svc-go-postgres"]