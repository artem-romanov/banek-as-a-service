run:
	go run ./cmd/main.go

build-linux:
	GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o build/server cmd/main.go 

build:
	go build -o build/server.exe -ldflags "-s -w" cmd/main.go 

.PHONY: build build-linux build-linux-gt
