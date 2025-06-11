FROM golang:1.24-alpine

WORKDIR /app/src

# Copy go.work files at the src level
COPY go.work go.work.sum ./

# Copy all modules
COPY shared/ shared/
COPY api/ api/
COPY listener/ listener/

# Set working directory to the api module
WORKDIR /app/src/api

# Download dependencies
RUN go mod download

# Build the application
RUN go build -o main ./cmd/main.go

# Run the application
CMD ["./main"]