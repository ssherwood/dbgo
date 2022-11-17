BINARY_NAME=dbgo

all: build test

build:
	go build -o ${BINARY_NAME} cmd/dbgo/dbgo.go

test:
	go test -v ./...

run:
	go build -o ${BINARY_NAME} cmd/dbgo/dbgo.go
	./${BINARY_NAME}

clean:
	go clean
	rm ${BINARY_NAME}