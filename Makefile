dev:
	go run main.go
prod:
	./start.sh
test:
	GIN_MODE=release go test -v -cover ./...
