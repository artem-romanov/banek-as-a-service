APP=app
BUILD_DIR=build
CMD_APP=./cmd/app

run:
	go run $(CMD_APP)

build-linux:
	GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o $(BUILD_DIR)/$(APP) $(CMD_APP)

build:
	go build -o build/server.exe -ldflags "-s -w" $(CMD_APP)

compress:
	upx -8 $(BUILD_DIR)/$(APP)

release-linux: build-linux compress


.PHONY: build build-linux build-linux-gt
