set dotenv-load := true



run:
  @go run main.go

gen:
  @go generate ./...