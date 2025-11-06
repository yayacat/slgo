# syntax=docker/dockerfile:1

# Build stage
FROM golang:1.20-alpine AS build

WORKDIR /src

# Install git for fetching dependencies
RUN apk add --no-cache git

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 go build -o /bin/server ./example

# Final stage
FROM alpine:latest

# Copy the binary from the build stage
COPY --from=build /bin/server /bin/server

# Copy the stations.json file
COPY example/stations.json /stations.json

# Expose the ports
EXPOSE 18000 18080

# Run the application
CMD ["/bin/server"]
