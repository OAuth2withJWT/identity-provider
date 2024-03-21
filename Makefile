.PHONY: build
build:
	go build -o ./build/server cmd/server/main.go

.PHONY: test
test:
	go test ./...

.PHONY: clean
clean:
	rm -rf ./build

.PHONY: run
run:
	go run cmd/server/main.go