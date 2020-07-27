build:
	go mod download
	go build -o bot

run: build 
	./bot

upgrade:
	go mod download
	go get -u -v
	go mod tidy
	go mod verify

test:
	go test

default: build
