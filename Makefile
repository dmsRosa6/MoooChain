build:
	go build -o ./bin/mooochain

run: build
	./bin/mooochain

test:
	go test -v ./...