FROM golang:alpine AS builder

# Copy and download dependency using go mod
COPY go.mod .
COPY go.sum .
RUN go mod download
RUN go mod verify

# Copy the code into the container
COPY . .

# Run test
RUN go test ./...

# Build the application
RUN go build -o tinyclientbuild .

# Command to run the executable
ENTRYPOINT ["./tinyclientbuild"]