.PHONY: deps clean build

deps:
	go get -u ./...

clean:
	rm -rf ./hello-world/hello-world

build:
	GOOS=linux GOARCH=amd64 go build -o main

zip:
	zip main.zip main
