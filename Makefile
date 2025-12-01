build:
	go build -o ./bin/mooochain

run: build
	./bin/mooochain

test:
	go test -v ./...

PB_OUT=core/pb

PROTO_FILES=$(wildcard core/*.proto)

all: proto

proto:
	@echo "Generating protobuf code..."
	@mkdir -p $(PB_OUT)
	protoc --go_out=. $(PROTO_FILES)
	@echo "Done."

clean:
	@echo "Cleaning generated protobuf files..."
	rm -f $(PB_OUT)/*.pb.go
	@echo "Done."

fmt:
	go fmt ./...

tidy:
	go mod tidy
