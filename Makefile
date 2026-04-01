run:
	go run ./cmd/app

build-linux:
	GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o build/server cmd/app 

build:
	go build -o build/server.exe -ldflags "-s -w" cmd/app

.PHONY: build build-linux build-linux-gt
