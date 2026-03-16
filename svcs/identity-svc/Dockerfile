# Build stage
FROM golang:1.26.1-alpine AS builder

WORKDIR /app

# Copy source code
COPY . .

# Install buf
RUN go install github.com/bufbuild/buf/cmd/buf@latest

# Generate Buf code
RUN buf generate svcs/identity-svc --include-imports --template svcs/identity-svc/buf.gen.yaml --output svcs/identity-svc

# Build the application
# Using CGO_ENABLED=0 for a static binary that works with distroless/static
RUN CGO_ENABLED=0 GOOS=linux go build -o /server svcs/identity-svc/cmd/server/main.go

# Final stage
FROM gcr.io/distroless/static-debian12

WORKDIR /

# Copy the binary from the builder stage
COPY --from=builder /server /server

# Expose the default port
EXPOSE 8080

# Run the application
ENTRYPOINT ["/server"]