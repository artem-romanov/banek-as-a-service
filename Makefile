run:
	go run ./cmd/main.go

build-linux:
	GOOS=linux GOARCH=amd64 go build -o build/server cmd/main.go 
	
build-linux-gt:
	GOEXPERIMENT=greenteagc GOOS=linux GOARCH=amd64 go build -o build/server cmd/main.go 

build:
	go build -o build/server.exe cmd/main.go 

.PHONY: build build-linux build-linux-gt
