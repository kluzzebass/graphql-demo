# Use the official Go image as the base image
FROM golang:alpine AS builder

# Set the working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Use a minimal alpine image for the final stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

# Set the working directory
WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /app/main .

# Copy README.md for Docker Hub repository overview
COPY --from=builder /app/README.md /README.md

# Add labels for Docker Hub
LABEL org.opencontainers.image.title="GraphQL Demo"
LABEL org.opencontainers.image.description="A GraphQL server demo built with Go and gqlgen"
LABEL org.opencontainers.image.source="https://github.com/kluzz/graphql-demo"
LABEL org.opencontainers.image.documentation="/README.md"

# Expose port 8080 (adjust if your app uses a different port)
EXPOSE 8080

# Run the application
CMD ["./main"]
