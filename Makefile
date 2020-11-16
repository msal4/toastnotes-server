dev:
	go run main.go
prod:
	./start.sh -b
test:
	GIN_MODE=release go test -v -cover ./...
