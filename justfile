set dotenv-load := true

# Platform list for multi-architecture builds
platforms := "linux/amd64,linux/arm64,linux/arm/v7"

# Default recipe - show available commands
default:
  @just --list

# Run the GraphQL server locally for development
run:
  @go run main.go

# Generate GraphQL code from schema (resolvers, models, etc.)
gen:
  @go generate ./...

# Build a single-platform Docker image for the current architecture
docker-build:
  @docker build -t kluzz/graphql-demo .

# Run the Docker container locally on port 8080
docker-run:
  @docker run --rm -p 8080:8080 kluzz/graphql-demo

# Push the locally built single-platform image to Docker Hub
docker-push:
  @docker push kluzz/graphql-demo:latest

# Build and push multi-platform Docker images (AMD64 + ARM64 + ARMv7) to Docker Hub
docker-pushx:
  @docker buildx build --platform {{platforms}} -t kluzz/graphql-demo:latest --push .

# Clean up buildkit containers and Docker resources
docker-cleanup:
  @docker buildx prune -f
  @docker stop $(docker ps -q --filter "name=buildx_buildkit") 2>/dev/null || true
  @docker rm $(docker ps -aq --filter "name=buildx_buildkit") 2>/dev/null || true
  @docker system prune -f


