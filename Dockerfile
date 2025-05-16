# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /build

# Install make
RUN apk add --no-cache make

# Copy only what's needed for building
COPY go.mod go.sum ./
COPY vendor/ ./vendor/
COPY cmd/ ./cmd/
COPY tools/ ./tools/
COPY Makefile ./

# Build the executable
RUN make build

# Final stage - using a minimal Alpine image
FROM alpine:latest

WORKDIR /app

# Copy only the executable from the build stage
COPY --from=builder /build/bin/server ./bin/server

# Make the binary executable
RUN chmod +x ./bin/server

# Expose port that your server uses
EXPOSE 8080

# No entrypoint here - will be specified in docker-compose.yaml
