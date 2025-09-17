# syntax=docker/dockerfile:1.4

# Multi-stage build for Gjallarhorn
FROM --platform=$BUILDPLATFORM node:18-alpine AS frontend-builder

WORKDIR /app/frontend

# Copy package files
COPY package*.json ./
RUN npm ci

# Copy frontend source
COPY . .

# Build frontend
RUN npm run build

# Go build stage
FROM --platform=$BUILDPLATFORM golang:1.21-alpine AS backend-builder

ARG TARGETOS
ARG TARGETARCH

WORKDIR /app

# Install git for go modules
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Copy built frontend from previous stage
COPY --from=frontend-builder /app/frontend/dist ./dist

# Build the application for the target platform
RUN CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH:-amd64} go build -a -installsuffix cgo -o main .

# Final stage
FROM --platform=$TARGETPLATFORM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# Copy the binary
COPY --from=backend-builder /app/main .

# Expose port
EXPOSE 8080

# Run the application
CMD ["./main"]
