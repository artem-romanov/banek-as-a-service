run:
	go run ./cmd/main.go

build-linux:
	GOOS=linux GOARCH=amd64 go build -o build/server cmd/main.go 