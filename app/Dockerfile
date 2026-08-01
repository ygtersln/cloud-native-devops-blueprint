# Build Stage
FROM golang:1.22-alpine AS builder

# Set environment variables for cross-compilation and security
ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /build

# Copy and download dependencies (if any)
COPY go.mod ./
# COPY go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build a statically linked binary
RUN go build -ldflags="-w -s" -o blueprint-api .

# Final Stage (Distroless / Scratch for maximum security & minimal size)
FROM gcr.io/distroless/static-debian12:nonroot

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /build/blueprint-api .

# Expose application port
EXPOSE 8080

# Run as non-root user for security (UID 65532 is standard in distroless)
USER 65532:65532

# Command to run
CMD ["/app/blueprint-api"]
