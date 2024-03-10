BINARY_NAME=ddns-update-server

build:
	GOARCH=arm64 GOOS=darwin go build -o bin/${BINARY_NAME}-darwin-arm64 main.go
	GOARCH=amd64 GOOS=linux go build -o bin/${BINARY_NAME}-linux-amd64 main.go

clean:
	go clean
	rm bin/${BINARY_NAME}-darwin-arm64
	rm bin/${BINARY_NAME}-linux-amd64
