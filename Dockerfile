# Use Go 1.23 bookworm as base image
FROM golang:1.23-bookworm AS base

# Node.js stage for frontend build
# =============================================================================
FROM node:22-bookworm AS frontend-builder

# Set working directory for frontend
WORKDIR /frontend

# Copy frontend package files
COPY web/package*.json ./

# Install dependencies
RUN npm ci

# Copy frontend source code
COPY web/ ./

# Build frontend assets
RUN npm run build

# Development stage
# =============================================================================
# Create a development stage based on the "base" image
FROM base AS development

# Change the working directory to /app
WORKDIR /app

ENV GOFLAGS=-buildvcs=false

# Install the air CLI for auto-reloading
RUN go install github.com/air-verse/air@latest

# Copy the go.mod and go.sum files to the /app directory
COPY go.mod go.sum ./

# Install dependencies
RUN go mod download

# Start air for live reloading
CMD ["air"]

# Builder stage
# =============================================================================
# Create a builder stage based on the "base" image
FROM base AS builder

# Move to working directory /build
WORKDIR /build

# Copy built frontend assets from frontend-builder stage
COPY --from=frontend-builder /frontend/dist ./web/dist/

# Copy the go.mod and go.sum files to the /build directory
COPY go.mod go.sum ./

# Install dependencies
RUN go mod download

# Copy the entire source code into the container
COPY . .

# Build the application
# Turn off CGO to ensure static binaries
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo \
    -o thomas .

# Production stage
# =============================================================================
# Create a production stage to run the application binary
FROM alpine:latest AS production

# Move to working directory /prod
WORKDIR /prod

# Copy binary from builder stage
COPY --from=builder /build/thomas ./

# Document the port that may need to be published
EXPOSE 7008

# Start the application
CMD ["/prod/thomas"]
