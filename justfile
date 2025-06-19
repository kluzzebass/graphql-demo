set dotenv-load := true



run:
  @go run main.go

gen:
  @go generate ./...

docker-build:
  @docker build -t kluzz/graphql-demo .

docker-run:
  @docker run --rm -p 8080:8080 kluzz/graphql-demo

docker-push:
  @docker push kluzz/graphql-demo:latest

